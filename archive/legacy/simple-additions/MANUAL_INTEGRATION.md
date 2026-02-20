# 🚀 Simple Manual Integration Guide

Since the automated script had issues, here's the manual copy-paste approach:

---

## Step 1: Create GitHub Actions Workflow

Create folder `.github/workflows/` in your repo, then create file `ci.yml` with this content:

```yaml
name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.24.1'

    - name: Test
      run: go test -v ./...

    - name: Build API
      run: go build -o bin/api ./cmd/api

    - name: Build Worker
      run: go build -o bin/worker ./cmd/worker

  docker:
    runs-on: ubuntu-latest
    needs: test
    steps:
    - uses: actions/checkout@v4

    - name: Build Docker images
      run: |
        docker build -f Dockerfile.api -t radar-api:latest .
        docker build -f Dockerfile.worker -t radar-worker:latest .
```

---

## Step 2: Add Observability Stack

Create file `docker-compose.observability.yml`:

```yaml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus:v2.48.0
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.retention.time=15d'

  grafana:
    image: grafana/grafana:10.2.3
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin
    depends_on:
      - prometheus

  jaeger:
    image: jaegertracing/all-in-one:1.50
    ports:
      - "16686:16686"
```

Create file `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'evolipia-radar-api'
    static_configs:
      - targets: ['host.docker.internal:8080']
    metrics_path: /metrics
```

---

## Step 3: Update Makefile

Add these targets to your existing Makefile (copy-paste at the end):

```makefile
# ============================================
# MLOps Additions
# ============================================

.PHONY: obs-up
obs-up: ## Start observability stack (Prometheus, Grafana, Jaeger)
	@echo "Starting observability stack..."
	docker-compose -f docker-compose.observability.yml up -d
	@echo "Grafana: http://localhost:3000 (admin/admin)"
	@echo "Prometheus: http://localhost:9090"
	@echo "Jaeger: http://localhost:16686"

.PHONY: obs-down
obs-down: ## Stop observability stack
	docker-compose -f docker-compose.observability.yml down

.PHONY: obs-logs
obs-logs: ## View observability stack logs
	docker-compose -f docker-compose.observability.yml logs -f

.PHONY: ci
ci: ## Run CI checks locally
	@echo "Running CI checks..."
	go mod tidy
	go vet ./...
	go test -v ./...
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker
	@echo "✅ CI checks passed!"

.PHONY: docker-build
docker-build: ## Build Docker images
	docker build -f Dockerfile.api -t radar-api:latest .
	docker build -f Dockerfile.worker -t radar-worker:latest .
```

---

## Step 4: Add Kubernetes Config (Optional)

Create folder `k8s/` and file `deployment.yml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: evolipia-radar-api
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
      - name: api
        image: radar-api:latest
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          value: "postgres://postgres:postgres@postgres:5432/radar?sslmode=disable"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 10
        readinessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: evolipia-radar-api
spec:
  selector:
    app: api
  ports:
  - port: 80
    targetPort: 8080
  type: LoadBalancer
```

---

## Step 5: Commit and Push

```bash
git add .
git commit -m "feat: Add CI/CD and observability stack"
git push origin mlops-improvements
```

---

## Usage

```bash
# Start observability stack
make obs-up

# View logs
make obs-logs

# Run CI checks locally
make ci

# Build Docker images
make docker-build
```

---

## Access Dashboards

After running `make obs-up`:

| Service | URL | Credentials |
|---------|-----|-------------|
| Grafana | http://localhost:3000 | admin/admin |
| Prometheus | http://localhost:9090 | - |
| Jaeger | http://localhost:16686 | - |

---

## What's Next?

1. **Add Prometheus metrics endpoint** to your API (in `cmd/api/main.go`):
   ```go
   import "github.com/prometheus/client_golang/prometheus/promhttp"

   // Add this route
   router.GET("/metrics", gin.WrapH(promhttp.Handler()))
   ```

2. **Run**: `go get github.com/prometheus/client_golang/prometheus`

3. **Test**: Start your API, then visit http://localhost:9090 to see metrics
