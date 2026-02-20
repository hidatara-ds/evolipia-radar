# =============================================================================
# EVOLIPIA-RADAR Makefile
# =============================================================================

# Variables
APP_NAME := evolipia-radar
GO_VERSION := 1.24.1
DOCKER_REGISTRY := ghcr.io/hidatara-ds
API_IMAGE := $(DOCKER_REGISTRY)/$(APP_NAME)/api
WORKER_IMAGE := $(DOCKER_REGISTRY)/$(APP_NAME)/worker

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := gofmt
GOLINT := golangci-lint

# Build flags
LDFLAGS := -ldflags "-s -w -X main.version=$(shell git describe --tags --always) -X main.buildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"

# =============================================================================
# Default Target
# =============================================================================

.PHONY: all
all: clean lint test build

# =============================================================================
# Development Setup
# =============================================================================

.PHONY: setup
setup: ## Setup development environment
	@echo "Setting up development environment..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/cosmtrek/air@latest
	@echo "Setup complete!"

.PHONY: setup-hooks
setup-hooks: ## Setup pre-commit hooks
	@echo "Installing pre-commit hooks..."
	pip install pre-commit
	pre-commit install
	pre-commit install --hook-type commit-msg

# =============================================================================
# Build Targets
# =============================================================================

.PHONY: build
build: build-api build-worker ## Build all binaries

.PHONY: build-api
build-api: ## Build API server binary
	@echo "Building API server..."
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) $(LDFLAGS) -o bin/api ./cmd/api

.PHONY: build-worker
build-worker: ## Build worker binary
	@echo "Building worker..."
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) $(LDFLAGS) -o bin/worker ./cmd/worker

.PHONY: build-all
build-all: ## Build for all platforms
	@echo "Building for all platforms..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/api-linux-amd64 ./cmd/api
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/api-linux-arm64 ./cmd/api
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o bin/api-darwin-amd64 ./cmd/api
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o bin/api-darwin-arm64 ./cmd/api

# =============================================================================
# Testing Targets
# =============================================================================

.PHONY: test
test: test-unit test-integration ## Run all tests

.PHONY: test-unit
test-unit: ## Run unit tests
	@echo "Running unit tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./internal/...
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GOTEST) -v -tags=integration ./tests/integration/...

.PHONY: test-e2e
test-e2e: ## Run end-to-end tests
	@echo "Running e2e tests..."
	$(GOTEST) -v -tags=e2e ./tests/e2e/...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -func=coverage.out

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests
	@echo "Running benchmark tests..."
	$(GOTEST) -bench=. -benchmem -run=^$$ ./... | tee benchmark.txt

# =============================================================================
# Linting & Formatting
# =============================================================================

.PHONY: lint
lint: ## Run all linters
	@echo "Running linters..."
	$(GOLINT) run --config=.golangci.yml ./...

.PHONY: lint-fix
lint-fix: ## Run linters with auto-fix
	@echo "Running linters with auto-fix..."
	$(GOLINT) run --config=.golangci.yml --fix ./...

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	$(GOFMT) -s -w .

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOCMD) vet ./...

.PHONY: tidy
tidy: ## Tidy go modules
	@echo "Tidying go modules..."
	$(GOMOD) tidy
	$(GOMOD) verify

# =============================================================================
# Database Operations
# =============================================================================

.PHONY: migrate-up
migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	migrate -path migrations -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	migrate -path migrations -database "$(DATABASE_URL)" down

.PHONY: migrate-create
migrate-create: ## Create a new migration (usage: make migrate-create name=add_users_table)
	@echo "Creating migration: $(name)"
	migrate create -ext sql -dir migrations -seq $(name)

.PHONY: migrate-version
migrate-version: ## Show current migration version
	@echo "Current migration version:"
	migrate -path migrations -database "$(DATABASE_URL)" version

.PHONY: migrate-force
migrate-force: ## Force migration version (usage: make migrate-force version=1)
	@echo "Forcing migration to version: $(version)"
	migrate -path migrations -database "$(DATABASE_URL)" force $(version)

.PHONY: db-reset
db-reset: migrate-down migrate-up ## Reset database

.PHONY: db-seed
db-seed: ## Seed database with test data
	@echo "Seeding database..."
	$(GOCMD) run scripts/seed/main.go

# =============================================================================
# Docker Operations
# =============================================================================

.PHONY: docker-build
docker-build: docker-build-api docker-build-worker ## Build all Docker images

.PHONY: docker-build-api
docker-build-api: ## Build API Docker image
	@echo "Building API Docker image..."
	docker build -f Dockerfile.api -t $(API_IMAGE):latest -t $(API_IMAGE):$(shell git describe --tags --always) .

.PHONY: docker-build-worker
docker-build-worker: ## Build Worker Docker image
	@echo "Building Worker Docker image..."
	docker build -f Dockerfile.worker -t $(WORKER_IMAGE):latest -t $(WORKER_IMAGE):$(shell git describe --tags --always) .

.PHONY: docker-push
docker-push: ## Push Docker images to registry
	@echo "Pushing Docker images..."
	docker push $(API_IMAGE):latest
	docker push $(API_IMAGE):$(shell git describe --tags --always)
	docker push $(WORKER_IMAGE):latest
	docker push $(WORKER_IMAGE):$(shell git describe --tags --always)

.PHONY: docker-compose-up
docker-compose-up: ## Start services with docker-compose
	@echo "Starting services..."
	docker-compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop services with docker-compose
	@echo "Stopping services..."
	docker-compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## View docker-compose logs
	@echo "Viewing logs..."
	docker-compose logs -f

# =============================================================================
# Kubernetes Operations
# =============================================================================

.PHONY: k8s-apply
k8s-apply: ## Apply Kubernetes manifests
	@echo "Applying Kubernetes manifests..."
	kubectl apply -k k8s/overlays/development

.PHONY: k8s-delete
k8s-delete: ## Delete Kubernetes resources
	@echo "Deleting Kubernetes resources..."
	kubectl delete -k k8s/overlays/development

.PHONY: helm-install
helm-install: ## Install Helm chart
	@echo "Installing Helm chart..."
	helm upgrade --install $(APP_NAME) ./k8s/helm/$(APP_NAME) 		--namespace $(APP_NAME) 		--create-namespace 		--values ./k8s/helm/$(APP_NAME)/values.yaml

.PHONY: helm-uninstall
helm-uninstall: ## Uninstall Helm chart
	@echo "Uninstalling Helm chart..."
	helm uninstall $(APP_NAME) --namespace $(APP_NAME)

.PHONY: helm-template
helm-template: ## Render Helm templates
	@echo "Rendering Helm templates..."
	helm template $(APP_NAME) ./k8s/helm/$(APP_NAME) 		--namespace $(APP_NAME) 		--values ./k8s/helm/$(APP_NAME)/values.yaml

# =============================================================================
# Development Server
# =============================================================================

.PHONY: run-api
run-api: ## Run API server locally
	@echo "Starting API server..."
	$(GOCMD) run ./cmd/api

.PHONY: run-worker
run-worker: ## Run worker locally
	@echo "Starting worker..."
	$(GOCMD) run ./cmd/worker

.PHONY: run-dev
run-dev: ## Run with hot reload (requires air)
	@echo "Starting with hot reload..."
	air -c .air.toml

# =============================================================================
# ML Pipeline Operations
# =============================================================================

.PHONY: ml-train
ml-train: ## Train ML model
	@echo "Training ML model..."
	cd ml/training && python train_ranking_model.py

.PHONY: ml-evaluate
ml-evaluate: ## Evaluate ML model
	@echo "Evaluating ML model..."
	cd ml/evaluation && python evaluate_model.py

.PHONY: ml-data-quality
ml-data-quality: ## Run data quality checks
	@echo "Running data quality checks..."
	cd ml/data_quality && great_expectations checkpoint run content_items_checkpoint

# =============================================================================
# Security Scanning
# =============================================================================

.PHONY: security-scan
security-scan: security-trivy security-nancy ## Run all security scans

.PHONY: security-trivy
security-trivy: ## Run Trivy security scan
	@echo "Running Trivy scan..."
	trivy fs --severity HIGH,CRITICAL .

.PHONY: security-nancy
security-nancy: ## Run Nancy dependency scan
	@echo "Running Nancy scan..."
	$(GOCMD) list -json -deps ./... | nancy sleuth

.PHONY: security-gosec
security-gosec: ## Run GoSec security scan
	@echo "Running GoSec scan..."
	gosec -fmt sarif -out gosec.sarif ./...

# =============================================================================
# Documentation
# =============================================================================

.PHONY: docs-generate
docs-generate: ## Generate API documentation
	@echo "Generating API documentation..."
	swag init -g cmd/api/main.go -o docs/api

.PHONY: docs-serve
docs-serve: ## Serve documentation locally
	@echo "Serving documentation..."
	cd docs && mkdocs serve

# =============================================================================
# Utilities
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f benchmark.txt
	docker system prune -f

.PHONY: deps
deps: ## Download and verify dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) verify

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy

.PHONY: generate
generate: ## Run go generate
	@echo "Running go generate..."
	$(GOCMD) generate ./...

.PHONY: proto
proto: ## Generate protobuf files
	@echo "Generating protobuf files..."
	protoc --go_out=. --go_opt=paths=source_relative 		--go-grpc_out=. --go-grpc_opt=paths=source_relative 		api/proto/*.proto

# =============================================================================
# CI/CD Helpers
# =============================================================================

.PHONY: ci
ci: lint test build ## Run CI pipeline locally

.PHONY: release
release: ## Create a new release (usage: make release version=v1.0.0)
	@echo "Creating release $(version)..."
	git tag -a $(version) -m "Release $(version)"
	git push origin $(version)

# =============================================================================
# Help
# =============================================================================

.PHONY: help
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
