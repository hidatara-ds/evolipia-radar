# Evolipia Radar: AI-Powered News Intelligence

> AI Research Intelligence Platform -- Go + Next.js + PostgreSQL

[![Go](https://img.shields.io/badge/Go-1.24.1-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-Scrape%20Pipeline-2088FF?logo=github-actions&logoColor=white)](https://github.com/hidatara-ds/evolipia-radar/actions)
[![Code Security](https://img.shields.io/badge/Code%20Security-CodeQL-2EA44F?logo=github&logoColor=white)](https://github.com/hidatara-ds/evolipia-radar/security)

Evolipia Radar is an automated intelligence pipeline that discovers AI/ML news, enriches and scores it with concurrent Go workers, then applies LLM-based processing for summarization and structured insights. The system is designed to run continuously through GitHub Actions and produce reliable JSON outputs for downstream API/UI consumption even when production database access is unavailable in CI runners.

## Executive Summary

This repository delivers a production-oriented news intelligence workflow:

- Collects curated AI/ML updates from multiple sources using high-throughput Go crawlers.
- Processes article text with LLM providers (OpenRouter or Gemini Flash) to generate concise summaries and tags.
- Stores outputs in PostgreSQL and JSON artifacts to support both API serving and static delivery workflows.
- Runs on scheduled GitHub Actions to keep content fresh with minimal operational overhead.

## Data Pipeline Architecture

```mermaid
flowchart LR
    A[Target URL Sources] --> B[Go Scraper Worker]
    B --> C[LLM Processing<br/>Gemini / OpenRouter]
    C --> D[JSON Storage<br/>data/news.json]
    C --> E[PostgreSQL Storage]
```

## Key Features

- Go concurrency for fast parallel scraping, normalization, and scoring.
- Dynamic LLM integration with provider/model configuration via environment variables.
- Automated and resilient CI workflow that keeps producing JSON output during restricted network conditions.
- Multi-target output strategy (`data/news.json` and API mirror) for frontend and serverless compatibility.

## Repository Layout

- `cmd/worker-json`: scraping + export worker used by scheduled automation.
- `pkg/`: shared business logic (crawler, AI, DB, scoring, services).
- `api/`: serverless endpoints consumed by deployed frontend.
- `app/`: Next.js interface.
- `data/`: generated JSON artifacts (`news.json`, `hybrid_result.json`, `sem_result.json`, `text_result.json`).

## Local Setup

1. Clone and enter repository:

```bash
git clone https://github.com/hidatara-ds/evolipia-radar.git
cd evolipia-radar
```

2. Install dependencies:

```bash
go mod download
npm install
```

3. Configure environment:

```bash
cp .env.example .env.local
```

Set required variables in `.env.local`:

```env
DATABASE_URL=postgresql://user:password@host:5432/dbname?sslmode=require
LLM_API_KEY=your_key
LLM_PROVIDER=openrouter
LLM_MODEL=google/gemini-flash-1.5
LLM_ENABLED=true
JSON_OUTPUT_PATH=data/news.json
```

Optional CI-safe variables:

```env
CI=true
SKIP_DB=true
```

When `CI=true` or `SKIP_DB=true`, database ping is bypassed and the worker can still emit JSON output without crashing.

4. Run worker export locally:

```bash
go run ./cmd/worker-json
```

5. Run frontend locally:

```bash
npm run dev
```

## CI/CD Workflow Notes

The scheduled workflow lives in `.github/workflows/scrape.yml` and runs every 30 minutes.

Pipeline behavior:

- Runs `go run ./cmd/worker-json` with JSON output enabled.
- Uses `CI=true` and `SKIP_DB=true` to avoid failing on database network restrictions in GitHub-hosted runners.
- Syncs `data/news.json` to `api/news.json` for deployment compatibility.
- Commits updated data artifacts when changes are detected.

This approach keeps automation healthy while preserving production database usage for non-CI runtime environments.

## Commands

- Scraper worker JSON export: `go run ./cmd/worker-json`
- API server: `go run ./cmd/api/main.go`
- Go tests: `go test ./...`
- Frontend dev server: `npm run dev`

## License

MIT License. See `LICENSE.md`.
