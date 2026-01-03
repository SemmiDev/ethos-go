# ============================================================================
# Ethos-Go Makefile
# ============================================================================

# Variables
APP_NAME := ethos-go
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
LDFLAGS := -ldflags "-w -s -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildTime=$(BUILD_TIME)"

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GORUN := $(GOCMD) run
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

# Docker/Podman
DOCKER := docker
COMPOSE := docker-compose

# Directories
MIGRATIONS_DIR := migrations
CMD_DIR := cmd/api
BUILD_DIR := build

# ============================================================================
# Development Commands
# ============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run the application locally
	@echo "🚀 Running application..."
	@$(GORUN) $(CMD_DIR)/main.go

.PHONY: build-frontend
build-frontend: ## Build the React frontend
	@echo "🎨 Building frontend..."
	@cd frontend && npm install && npm run build
	@echo "📦 Copying frontend build to embed directory..."
	@rm -rf internal/web/dist
	@cp -r frontend/dist internal/web/dist
	@echo "✅ Frontend built and ready for embedding"

.PHONY: build
build: build-frontend ## Build the full application with embedded frontend
	@echo "🔨 Building $(APP_NAME) with embedded frontend..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "✅ Binary created at $(BUILD_DIR)/$(APP_NAME)"
	@echo "📊 Binary size: $$(du -h $(BUILD_DIR)/$(APP_NAME) | cut -f1)"

.PHONY: build-backend
build-backend: ## Build only the Go backend (without frontend)
	@echo "🔨 Building $(APP_NAME) backend only..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME) ./$(CMD_DIR)
	@echo "✅ Binary created at $(BUILD_DIR)/$(APP_NAME)"

.PHONY: build-linux
build-linux: build-frontend ## Build Linux binary with embedded frontend (for Docker)
	@echo "🔨 Building $(APP_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(APP_NAME)-linux ./$(CMD_DIR)
	@echo "✅ Linux binary created at $(BUILD_DIR)/$(APP_NAME)-linux"

# ============================================================================
# Testing
# ============================================================================

.PHONY: test
test: ## Run tests
	@echo "🧪 Running tests..."
	@$(GOTEST) -v -race -cover ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "🧪 Running tests with coverage..."
	@$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

.PHONY: test-short
test-short: ## Run short tests only
	@echo "🧪 Running short tests..."
	@$(GOTEST) -v -short ./...

# ============================================================================
# Code Quality
# ============================================================================

.PHONY: fmt
fmt: ## Format code
	@echo "🎨 Formatting code..."
	@$(GOFMT) ./...
	@echo "✅ Code formatted"

.PHONY: vet
vet: ## Run go vet
	@echo "🔍 Running go vet..."
	@$(GOVET) ./...
	@echo "✅ Vet check passed"

.PHONY: lint
lint: ## Run golangci-lint
	@echo "🔍 Running linter..."
	@golangci-lint run ./...
	@echo "✅ Lint check passed"

.PHONY: check
check: fmt vet ## Run all code quality checks
	@echo "✅ All checks passed"

# ============================================================================
# Dependencies
# ============================================================================

.PHONY: deps
deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@$(GOMOD) download
	@echo "✅ Dependencies downloaded"

.PHONY: deps-tidy
deps-tidy: ## Tidy dependencies
	@echo "📦 Tidying dependencies..."
	@$(GOMOD) tidy
	@echo "✅ Dependencies tidied"

.PHONY: deps-verify
deps-verify: ## Verify dependencies
	@echo "📦 Verifying dependencies..."
	@$(GOMOD) verify
	@echo "✅ Dependencies verified"

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "📦 Updating dependencies..."
	@$(GOMOD) get -u ./...
	@$(GOMOD) tidy
	@echo "✅ Dependencies updated"

# ============================================================================
# Code Generation
# ============================================================================

.PHONY: generate
generate: ## Generate OpenAPI code
	@echo "🔄 Generating OpenAPI code..."
	@$(GORUN) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config api/openapi/habits-cfg.yaml api/openapi/habits.yml
	@$(GORUN) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config api/openapi/auth-cfg.yaml api/openapi/auth.yml
	@$(GORUN) github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config api/openapi/notifications-cfg.yaml api/openapi/notifications.yml
	@echo "✅ Code generated"

.PHONY: generate-mocks
generate-mocks: ## Generate mocks for testing
	@echo "🔄 Generating mocks..."
	@$(GORUN) github.com/vektra/mockery/v2@latest --all
	@echo "✅ Mocks generated"

# ============================================================================
# Database Migrations
# ============================================================================

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=migration_name)
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: name is required. Usage: make migrate-create name=migration_name"; \
		exit 1; \
	fi
	@echo "📝 Creating migration: $(name)..."
	@migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)
	@echo "✅ Migration created"

.PHONY: migrate-up
migrate-up: ## Run all up migrations
	@echo "⬆️  Running migrations up..."
	@migrate -path $(MIGRATIONS_DIR) -database "${DATABASE_URL}" up
	@echo "✅ Migrations applied"

.PHONY: migrate-down
migrate-down: ## Rollback last migration
	@echo "⬇️  Rolling back migration..."
	@migrate -path $(MIGRATIONS_DIR) -database "${DATABASE_URL}" down 1
	@echo "✅ Migration rolled back"

.PHONY: migrate-force
migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then \
		echo "❌ Error: version is required. Usage: make migrate-force version=1"; \
		exit 1; \
	fi
	@echo "🔧 Forcing migration to version $(version)..."
	@migrate -path $(MIGRATIONS_DIR) -database "${DATABASE_URL}" force $(version)
	@echo "✅ Migration forced"

# ============================================================================
# Docker Commands
# ============================================================================

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@$(DOCKER) build -t $(APP_NAME):$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		.
	@$(DOCKER) tag $(APP_NAME):$(VERSION) $(APP_NAME):latest
	@echo "✅ Docker image built: $(APP_NAME):$(VERSION)"

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@$(DOCKER) run --rm -p 8080:8080 --env-file .env $(APP_NAME):latest

.PHONY: docker-push
docker-push: ## Push Docker image to registry
	@echo "🐳 Pushing Docker image..."
	@$(DOCKER) push $(APP_NAME):$(VERSION)
	@$(DOCKER) push $(APP_NAME):latest
	@echo "✅ Docker image pushed"

# ============================================================================
# Docker Compose Commands
# ============================================================================

.PHONY: compose-up
compose-up: ## Start all services with docker-compose
	@echo "🐳 Starting services..."
	@$(COMPOSE) -f compose.dev.yml up -d
	@echo "✅ Services started"

.PHONY: compose-down
compose-down: ## Stop all services
	@echo "🐳 Stopping services..."
	@$(COMPOSE) -f compose.dev.yml down
	@echo "✅ Services stopped"

.PHONY: compose-logs
compose-logs: ## Show logs from all services
	@$(COMPOSE) -f compose.dev.yml logs -f

.PHONY: compose-ps
compose-ps: ## Show running services
	@$(COMPOSE) -f compose.dev.yml ps

.PHONY: compose-restart
compose-restart: compose-down compose-up ## Restart all services

.PHONY: compose-build
compose-build: ## Build and start services
	@echo "🐳 Building and starting services..."
	@export VERSION=$(VERSION) COMMIT=$(COMMIT) BUILD_TIME=$(BUILD_TIME) && \
		$(COMPOSE) -f compose.dev.yml up -d --build
	@echo "✅ Services built and started"

.PHONY: compose-rebuild
compose-rebuild: ## Rebuild services from scratch
	@echo "🐳 Rebuilding services..."
	@$(COMPOSE) -f compose.dev.yml down -v
	@export VERSION=$(VERSION) COMMIT=$(COMMIT) BUILD_TIME=$(BUILD_TIME) && \
		$(COMPOSE) -f compose.dev.yml up -d --build --force-recreate
	@echo "✅ Services rebuilt"

# ============================================================================
# Database Commands (via Docker Compose)
# ============================================================================

.PHONY: db-shell
db-shell: ## Open PostgreSQL shell
	@echo "🗄️  Opening database shell..."
	@$(COMPOSE) -f compose.dev.yml exec postgresql psql -U postgres -d ethosgo

.PHONY: db-reset
db-reset: ## Reset database (WARNING: destroys all data)
	@echo "⚠️  Resetting database..."
	@$(COMPOSE) -f compose.dev.yml down -v
	@$(COMPOSE) -f compose.dev.yml up -d postgresql
	@echo "✅ Database reset"

# ============================================================================
# Utility Commands
# ============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "✅ Cleaned"

.PHONY: version
version: ## Show version information
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"

.PHONY: network
network: ## Create Docker network
	@echo "🌐 Creating network..."
	@$(DOCKER) network create ethosgo-network 2>/dev/null || true
	@echo "✅ Network created"

# ============================================================================
# Quick Start Commands
# ============================================================================

.PHONY: dev
dev: compose-up ## Start development environment
	@echo ""
	@echo "✅ Development environment ready!"
	@echo ""
	@echo "🌐 Services:"
	@echo "   Frontend:    http://localhost:3001"
	@echo "   API:         http://localhost:8080"
	@echo "   Grafana:     http://localhost:3000 (admin/admin)"
	@echo "   Prometheus:  http://localhost:9090"
	@echo "   Jaeger:      http://localhost:16686"
	@echo ""
	@echo "📚 API Docs:"
	@echo "   Auth:        http://localhost:8080/auth/doc"
	@echo "   Habits:      http://localhost:8080/habits/doc"
	@echo ""

.PHONY: stop
stop: compose-down ## Stop development environment

.PHONY: restart
restart: compose-restart ## Restart development environment

.PHONY: logs
logs: compose-logs ## Show application logs

.PHONY: logs-frontend
logs-frontend: ## Show frontend logs only
	@$(COMPOSE) -f compose.dev.yml logs -f frontend

.PHONY: logs-app
logs-app: ## Show backend app logs only
	@$(COMPOSE) -f compose.dev.yml logs -f app

.PHONY: rebuild-frontend
rebuild-frontend: ## Rebuild frontend only
	@echo "🐳 Rebuilding frontend..."
	@$(COMPOSE) -f compose.dev.yml up -d --build --force-recreate frontend
	@echo "✅ Frontend rebuilt"

.PHONY: rebuild-app
rebuild-app: ## Rebuild backend app only
	@echo "🐳 Rebuilding backend..."
	@export VERSION=$(VERSION) COMMIT=$(COMMIT) BUILD_TIME=$(BUILD_TIME) && \
		$(COMPOSE) -f compose.dev.yml up -d --build --force-recreate app
	@echo "✅ Backend rebuilt"

# ============================================================================
# CI/CD Commands
# ============================================================================

.PHONY: ci
ci: deps check test ## Run CI pipeline
	@echo "✅ CI pipeline completed"

.PHONY: pre-commit
pre-commit: fmt vet test-short ## Run pre-commit checks
	@echo "✅ Pre-commit checks passed"

# ============================================================================
# Default Target
# ============================================================================

.DEFAULT_GOAL := help
