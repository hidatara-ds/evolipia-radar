.PHONY: help migrate-up migrate-down run-api run-worker docker-up docker-down

help:
	@echo "Available targets:"
	@echo ""
	@echo "Development:"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  run-api       - Run API server"
	@echo "  run-worker    - Run worker"
	@echo "  docker-up     - Start PostgreSQL with docker-compose"
	@echo "  docker-down   - Stop PostgreSQL"
	@echo ""
	@echo "Mobile/PWA:"
	@echo "  test-pwa          - Test PWA configuration locally"
	@echo "  deploy-mobile     - Deploy to hosting platform"
	@echo "  build-mobile      - Build optimized production binary"
	@echo "  docker-build-mobile - Build Docker image"
	@echo "  lighthouse        - Run Lighthouse PWA audit"
	@echo ""
	@echo "Observability:"
	@echo "  obs-up        - Start observability stack (Grafana/Prometheus)"
	@echo "  obs-down      - Stop observability stack"
	@echo ""
	@echo "For more info: make mobile-help"

migrate-up:
	@echo "Running migrations..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up

migrate-down:
	@echo "Rolling back migrations..."
	@migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" down

run-api:
	@echo "Starting API server..."
	@go run ./cmd/api

run-worker:
	@echo "Starting worker..."
	@go run ./cmd/worker

docker-up:
	@echo "Starting PostgreSQL..."
	@docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@sleep 5

docker-down:
	@echo "Stopping PostgreSQL..."
	@docker-compose down

	# ============================================
# MLOps Additions
# ============================================

.PHONY: obs-up
obs-up:
	@echo "Starting observability stack..."
	docker-compose -f docker-compose.observability.yml up -d
	@echo "Grafana: http://localhost:3000 (admin/admin)"

.PHONY: obs-down
obs-down:
	docker-compose -f docker-compose.observability.yml down

.PHONY: ci
ci:
	@echo "Running CI checks..."
	go mod tidy
	go vet ./...
	go test -v ./...
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	@echo "✅ CI checks passed!"
# ============================================
# MLOps Additions
# ============================================

.PHONY: docker-compose-obs
docker-compose-obs: ## Start observability stack
	docker-compose -f docker-compose.observability.yml up -d

.PHONY: docker-compose-ml
docker-compose-ml: ## Start ML stack
	docker-compose -f docker-compose.ml.yml up -d

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run --config=.golangci.yml ./...

.PHONY: security-scan
security-scan: ## Run security scans
	trivy fs --severity HIGH,CRITICAL .

# ============================================
# Mobile/PWA Targets
# ============================================

.PHONY: test-pwa
test-pwa:
	@echo "Testing PWA configuration..."
	@chmod +x test-pwa.sh
	@./test-pwa.sh

.PHONY: deploy-mobile
deploy-mobile:
	@echo "Deploying mobile PWA..."
	@chmod +x deploy-mobile.sh
	@./deploy-mobile.sh

.PHONY: build-mobile
build-mobile:
	@echo "Building optimized production binary..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o api ./cmd/api
	@echo "✅ Binary built: ./api"

.PHONY: docker-build-mobile
docker-build-mobile:
	@echo "Building Docker image for mobile deployment..."
	@docker build -f Dockerfile.api -t evolipia-radar-api:latest .
	@echo "✅ Image built: evolipia-radar-api:latest"

.PHONY: lighthouse
lighthouse:
	@echo "Running Lighthouse PWA audit..."
	@echo "Make sure server is running on http://localhost:8080"
	@npx lighthouse http://localhost:8080 --view --preset=pwa

.PHONY: mobile-help
mobile-help:
	@echo "Mobile/PWA targets:"
	@echo "  test-pwa          - Test PWA configuration locally"
	@echo "  deploy-mobile     - Deploy to hosting platform (Fly.io/Render/Railway)"
	@echo "  build-mobile      - Build optimized production binary"
	@echo "  docker-build-mobile - Build Docker image"
	@echo "  lighthouse        - Run Lighthouse PWA audit"
