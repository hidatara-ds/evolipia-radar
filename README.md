<div align="center">

# 🎯 EVOLIPIA-RADAR

[![CI](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/ci.yml/badge.svg)](https://github.com/hidatara-ds/evolipia-radar/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/hidatara-ds/evolipia-radar)](https://goreportcard.com/report/github.com/hidatara-ds/evolipia-radar)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**An AI/ML tech news aggregator with MLOps best practices**

</div>

---

# evolipia-radar

Engineering Verified Overview of Latest Insights, Priorities, Impact & Analytics

An AI/ML tech news aggregator backend that ranks and summarizes notable items from multiple sources.

## Features

- **Multi-source aggregation**: Hacker News, RSS/Atom feeds, arXiv, and custom JSON APIs
- **Intelligent ranking**: Combines popularity, relevance, credibility, and novelty signals
- **Automatic summarization**: Extractive summaries with AI/ML engineer-focused insights
- **Deduplication**: Prevents duplicate items across sources
- **RESTful API**: Clean endpoints for feeds, search, and source management
- **Security**: SSRF protection, rate limiting, and input validation

## Architecture

### Components

- **API Server** (`cmd/api`): Serves REST endpoints for feeds, search, and source management
- **Worker** (`cmd/worker`): Scheduled ingestion, scoring, and summarization
- **Database**: PostgreSQL with proper indexes for performance
- **Scoring**: Configurable weights for popularity, relevance, credibility, and novelty

### Design Principles

This project follows **SOLID principles** and **Clean Architecture**:

1. **Separation of Concerns**: Each layer has a single, well-defined responsibility
   - HTTP handlers only handle HTTP concerns
   - Services contain business logic
   - Repositories handle data access
   - DTOs separate data transfer from domain models

2. **Single Responsibility Principle**: Each package/struct has one reason to change
   - Handlers: HTTP request/response handling
   - Services: Business logic orchestration
   - Repositories: Data persistence
   - Connectors: External data fetching

3. **Dependency Inversion**: High-level modules don't depend on low-level modules
   - Handlers depend on Services (abstraction)
   - Services depend on Repositories (abstraction)
   - No direct database access from handlers

4. **Configuration Management**: Hardcoded values moved to config structures
   - Scoring configs in `internal/scoring/config.go`
   - Summarizer configs in `internal/summarizer/config.go`
   - Application configs in `internal/config/`

## Quick Start

### Prerequisites

- Go 1.21+ (tested with Go 1.24.1)
- PostgreSQL 15+
- Docker & Docker Compose (optional)
- [migrate](https://github.com/golang-migrate/migrate) CLI tool (for database migrations)

### Dependencies

This project uses Go modules. Dependencies are managed in `go.mod` and `go.sum`.

**Main dependencies:**
- `github.com/gin-gonic/gin` - HTTP web framework
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `github.com/robfig/cron/v3` - Cron job scheduling

**Install dependencies:**
```bash
go mod download
# or
go mod tidy
```

**View all dependencies:**
```bash
go list -m all
```

### Local Development

1. **Start PostgreSQL**:
   ```bash
   docker-compose up -d postgres
   ```

2. **Run migrations**:
   ```bash
   make migrate-up
   ```
   **Windows (no make):** install [golang-migrate](https://github.com/golang-migrate/migrate) CLI (e.g. `scoop install migrate` or download from releases), then:
   ```bash
   migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up
   ```

3. **Set environment variables** (optional):
   ```bash
   export DATABASE_URL="postgres://postgres:postgres@localhost:5432/radar?sslmode=disable"
   export PORT=8080
   export WORKER_CRON="*/10 * * * *"  # Every 10 minutes
   ```

4. **Run API server** (in one terminal):
   ```bash
   make run-api
   # or (Windows / no make):
   go run ./cmd/api
   ```

5. **Run worker** (in another terminal):
   ```bash
   make run-worker
   # or (Windows / no make):
   go run ./cmd/worker
   ```

### Web UI (mobile-first)

Setelah API jalan, buka di browser:

- **http://localhost:8080/** — tampilan utama (dibuat untuk ukuran HP)

UI berisi: **Feed**, **Rising**, **Cari**, **Chat AI**, **Sumber**, **Pengaturan**. Ketuk item untuk detail.

- **Chat AI:** Terhubung ke OpenRouter (`arcee-ai/trinity-large-preview:free`) langsung dari frontend untuk demo/chat seputar berita. Untuk produksi, ganti endpoint & API key di `web/index.html` agar menggunakan akun/model milikmu sendiri.
- **Feed/Rising/Search:** Memanggil API yang sama dengan klien mobile nanti, sehingga web UI ini bisa kamu pakai sebagai “preview” UX dan sebagai playground API.

**Cek tampilan ukuran HP:** DevTools (F12) → device toolbar (Ctrl+Shift+M) → pilih perangkat atau lebar ~375px.

### API Endpoints

- `GET /healthz` - Health check
- `GET /v1/feed?date=today&topic=llm` - Top 20 daily feed
- `GET /v1/rising?window=2h` - Rising items in last 2 hours
- `GET /v1/items/{id}` - Item details with scores and summary
- `GET /v1/search?q=rag&topic=llm` - Search items
- `GET /v1/sources` - List all sources
- `POST /v1/sources` - Create new source
- `POST /v1/sources/test` - Test source connection
- `PATCH /v1/sources/{id}/enable` - Enable/disable source

### Example: Add RSS Source

```bash
curl -X POST http://localhost:8080/v1/sources/test \
  -H "Content-Type: application/json" \
  -d '{
    "type": "rss_atom",
    "category": "news",
    "url": "https://openai.com/blog/rss.xml"
  }'

curl -X POST http://localhost:8080/v1/sources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "OpenAI Blog",
    "type": "rss_atom",
    "category": "news",
    "url": "https://openai.com/blog/rss.xml"
  }'

curl -X PATCH http://localhost:8080/v1/sources/{id}/enable \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

### Example: Add JSON API Source

```bash
curl -X POST http://localhost:8080/v1/sources/test \
  -H "Content-Type: application/json" \
  -d '{
    "type": "json_api",
    "category": "news",
    "url": "https://api.example.com/news",
    "mapping_json": {
      "items_path": "data.articles",
      "title_path": "title",
      "url_path": "link",
      "published_at_path": "published_date",
      "summary_path": "excerpt"
    }
  }'
```

## Environment Variables

- `DATABASE_URL` - PostgreSQL connection string (default: `postgres://postgres:postgres@localhost:5432/radar?sslmode=disable`)
- `PORT` - API server port (default: `8080`)
- `CACHE_TTL_SECONDS` - Cache TTL for feed responses (default: `60`)
- `WORKER_CRON` - Cron schedule for worker (default: `*/10 * * * *` - every 10 minutes)
- `MAX_FETCH_BYTES` - Maximum response size in bytes (default: `2000000` - 2MB)
- `FETCH_TIMEOUT_SECONDS` - Request timeout in seconds (default: `8`)

## Docker Deployment

### Build images:
```bash
docker build -f Dockerfile.api -t radar-api .
docker build -f Dockerfile.worker -t radar-worker .
```

### Run:
```bash
docker-compose up -d postgres
# Run migrations
docker run --rm --network host -v $(pwd)/migrations:/migrations migrate/migrate \
  -path /migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up

docker run -d --name radar-api --network host \
  -e DATABASE_URL="postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" \
  radar-api

docker run -d --name radar-worker --network host \
  -e DATABASE_URL="postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" \
  radar-worker
```

## MLOps & Observability (opsional tapi siap pakai)

Branch ini menambahkan beberapa tambahan MLOps/infra yang bisa kamu aktifkan sesuai kebutuhan:

- **Observability stack** (Prometheus, Grafana, Jaeger):
  - File: `docker-compose.observability.yml`
  - Makefile targets:
    - `make obs-up` — start stack (Grafana: `http://localhost:3000`, Prometheus: `http://localhost:9090`, Jaeger: `http://localhost:16686`)
    - `make obs-down` — stop stack
- **Local CI helper**:
  - `make ci` — jalankan `go vet`, `go test ./...`, dan build API/worker secara lokal.
- **ML stack (eksperimen)**:
  - File: `docker-compose.ml.yml` — untuk menjalankan komponen ML tambahan (bisa kamu kembangkan bertahap).
- **CI di GitHub Actions**:
  - Workflow utama: `.github/workflows/ci.yml` (lint, test, build).
  - Workflow tambahan: security, CD, dan ML-pipeline (lihat isi folder `.github/workflows/`).


## Project Structure

```
.
├── cmd/
│   ├── api/                  # API server entry point
│   └── worker/               # Worker entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── db/                   # Database connection and repositories (data access layer)
│   ├── dto/                  # Data Transfer Objects (DTOs for API/connector boundaries)
│   ├── models/               # Domain models (pure data structures)
│   ├── http/                 # HTTP handlers (presentation layer)
│   │   └── handlers/         # HTTP request handlers
│   ├── services/             # Business logic services (application layer)
│   ├── connectors/           # Source connectors (HN, RSS, arXiv, JSON API)
│   ├── scoring/              # Ranking and scoring algorithms
│   │   └── config.go         # Scoring configuration (credibility, relevance keywords)
│   ├── summarizer/           # Extractive summarization
│   │   └── config.go         # Summarizer configuration (topic keywords)
│   ├── normalizer/           # URL normalization and deduplication
│   ├── security/             # SSRF protection
│   ├── workflows/            # Temporal workflows (ML pipeline orchestration - scaffold)
│   ├── activities/           # Temporal activities for each ML pipeline step (scaffold)
│   └── mlpipeline/           # Shared ML pipeline types used by workflows/activities
├── web/                      # Mobile-first web UI + Chat AI (OpenRouter)
├── migrations/               # Database migrations
├── configs/                  # Configuration files
├── .github/                  # GitHub Actions (CI/CD, security, ML pipeline)
├── k8s/                      # Base manifests + Helm chart (optional K8s deploy)
├── terraform/                # Terraform module for EKS (optional infra-as-code)
├── archive/                  # Backups & legacy integration scripts (not used at runtime)
└── docker-compose.yml
```

### Architecture Layers

- **Presentation Layer** (`http/handlers`): HTTP request/response handling, only uses services
- **Application Layer** (`services`): Business logic orchestration, uses repositories and domain services
- **Domain Layer** (`models`, `dto`): Domain models and data transfer objects
- **Data Access Layer** (`db`): Repository pattern for database operations
- **Infrastructure Layer** (`connectors`, `scoring`, `summarizer`, `normalizer`, `security`): External integrations and utilities

## Scoring Formula

Final score = `w1*popularity + w2*relevance + w3*credibility + w4*novelty`

Default weights:
- `w1` (popularity): 0.55 - Based on HN points/comments with recency decay
- `w2` (relevance): 0.25 - AI/ML keyword matching and topic classification
- `w3` (credibility): 0.15 - Domain whitelist/blacklist scoring
- `w4` (novelty): 0.05 - Recency-based scoring

## Security Features

- **SSRF Protection**: Blocks localhost, private IPs, and link-local addresses
- **Rate Limiting**: Per-source fetch limits
- **Input Validation**: URL and mapping validation
- **Size Limits**: Response size caps (2MB default)
- **Timeouts**: Configurable request timeouts (8s default)

## Development

### Setup Development Environment

1. **Install Go dependencies:**
   ```bash
   go mod download
   ```

2. **Install development tools** (optional but recommended):
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Install migrate CLI
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

### Running Tests
```bash
go test ./...
```

### Linting
```bash
golangci-lint run
```

### Database Migrations
```bash
# Create new migration
migrate create -ext sql -dir migrations -seq migration_name

# Up
make migrate-up

# Down
make migrate-down
```

### Code Organization Principles

This project follows **Separation of Concerns** and **Single Responsibility Principle**:

- **DTOs** (`internal/dto`): Data Transfer Objects for API boundaries and external connectors
- **Models** (`internal/models`): Pure domain models representing database entities
- **Repositories** (`internal/db`): Data access layer, no business logic
- **Services** (`internal/services`): Business logic layer, orchestrates repositories and domain services
- **Handlers** (`internal/http/handlers`): HTTP layer, only uses services (no direct repository access)
- **Configs**: Hardcoded configurations moved to dedicated config files for maintainability

## Documentation

Semua dokumentasi tambahan (setup lokal, enhancement plan, dependencies, dll) ada di folder **[docs/](docs/)**. Indeks lengkap: [docs/README.md](docs/README.md).

## License

MIT
