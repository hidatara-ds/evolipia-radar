.PHONY: help migrate-up migrate-down run-api run-worker docker-up docker-down

help:
	@echo "Available targets:"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  run-api       - Run API server"
	@echo "  run-worker    - Run worker"
	@echo "  docker-up     - Start PostgreSQL with docker-compose"
	@echo "  docker-down   - Stop PostgreSQL"

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