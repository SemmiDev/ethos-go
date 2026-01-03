# Build Stage for Frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Build Stage for Backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Copy built frontend assets to the expected location for embedding
COPY --from=frontend-builder /app/frontend/dist ./internal/web/dist

# Install migrate using go install (or download binary)
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Build the binaries
# CGO_ENABLED=0 for static binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ethos-api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ethos-worker ./cmd/worker

# Final Stage
FROM gcr.io/distroless/static-debian12 AS final

COPY --from=backend-builder /app/ethos-api /ethos-api
COPY --from=backend-builder /app/ethos-worker /ethos-worker
COPY --from=backend-builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=backend-builder /app/migrations /migrations

# Default to running the API
ENTRYPOINT ["/ethos-api"]
