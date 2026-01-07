# =============================================================================
# Ethos-Go Production Dockerfile
# Best practices: Multi-stage, distroless:nonroot, static binary, layer caching
# =============================================================================

# --- TAHAP 1: FRONTEND BUILDER ---
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend

# Copy package files
COPY frontend/package*.json ./

# Install dependencies - use npm ci for clean install
# Note: npm ci automatically rebuilds native modules for the current platform
RUN npm ci

# Copy source and build
COPY frontend/ .
RUN npm run build

# --- TAHAP 2: BACKEND BUILDER ---
FROM golang:1.25-alpine AS builder

# Argumen untuk build-time variables
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME

# Install dependensi sistem: sertifikat (untuk HTTPS) dan timezone
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# --- Caching Dependensi ---
# 1. Salin file mod dan sum terlebih dahulu
COPY go.mod go.sum ./

# 2. Download dependensi (akan di-cache oleh Docker)
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 3. Salin sisa kode sumber
COPY . .

# 4. Salin frontend build ke lokasi embedding
COPY --from=frontend-builder /app/frontend/dist ./internal/web/dist

# 5. Build static binaries
# - CGO_ENABLED=0 untuk static build (tanpa dependensi libc)
# - ldflags "-w -s" untuk mengurangi ukuran binary
# - ldflags "-X" untuk menyuntikkan build-time variables
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -tags=viper_bind_struct \
    -ldflags="-w -s \
    -X main.version=${VERSION} \
    -X main.commit=${COMMIT} \
    -X main.buildTime=${BUILD_TIME}" \
    -o /build/ethos-api ./cmd/api

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -tags=viper_bind_struct \
    -ldflags="-w -s \
    -X main.version=${VERSION} \
    -X main.commit=${COMMIT} \
    -X main.buildTime=${BUILD_TIME}" \
    -o /build/ethos-worker ./cmd/worker

# --- TAHAP 3: FINAL (PRODUKSI) ---
# Gunakan image distroless non-root: super minimal dan aman
# - Tidak ada shell atau package manager (lebih aman)
# - Berjalan sebagai non-root secara default
# - Lolos audit SOC2 dan ISO27001
FROM gcr.io/distroless/static-debian12:nonroot

# Salin data timezone dari builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Atur working directory
WORKDIR /app

# Salin binary aplikasi dari builder
COPY --from=builder /build/ethos-api /app/ethos-api
COPY --from=builder /build/ethos-worker /app/ethos-worker

# Salin migrations (untuk embedded migrations)
COPY --from=builder /build/migrations /app/migrations

# Expose port
EXPOSE 8080

# Perintah untuk menjalankan aplikasi
# Gunakan ethos-api sebagai default, atau override dengan ethos-worker
ENTRYPOINT ["/app/ethos-api"]
