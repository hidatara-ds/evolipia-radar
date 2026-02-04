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

- **API Server** (`cmd/api`): Serves REST endpoints for feeds, search, and source management
- **Worker** (`cmd/worker`): Scheduled ingestion, scoring, and summarization
- **Database**: PostgreSQL with proper indexes for performance
- **Scoring**: Configurable weights for popularity, relevance, credibility, and novelty

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### Local Development

1. **Start PostgreSQL**:
   ```bash
   docker-compose up -d postgres
   ```

2. **Run migrations**:
   ```bash
   make migrate-up
   ```
   Or manually:
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
   # or
   go run ./cmd/api
   ```

5. **Run worker** (in another terminal):
   ```bash
   make run-worker
   # or
   go run ./cmd/worker
   ```

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

## Project Structure

```
.
├── cmd/
│   ├── api/          # API server entry point
│   └── worker/       # Worker entry point
├── internal/
│   ├── config/       # Configuration management
│   ├── db/           # Database connection and repositories
│   ├── models/       # Data models
│   ├── http/         # HTTP handlers
│   ├── services/     # Business logic services
│   ├── connectors/   # Source connectors (HN, RSS, arXiv, JSON API)
│   ├── scoring/      # Ranking and scoring algorithms
│   ├── summarizer/   # Extractive summarization
│   ├── normalizer/   # URL normalization and deduplication
│   └── security/     # SSRF protection
├── migrations/       # Database migrations
├── configs/          # Configuration files
└── docker-compose.yml
```

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

## License

MIT
