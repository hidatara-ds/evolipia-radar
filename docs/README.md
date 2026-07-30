# 📚 Evolipia Radar - Documentation Hub & AI Agent Blueprint

> **Central Information & Master Documentation for Evolipia Radar AI Agents and Software Engineers.**

Welcome to the official documentation repository for **Evolipia Radar**. This documentation suite is specifically structured so that **AI Assistants (LLMs)** and **Developers** can rapidly understand the system architecture, data models, API endpoints, development workflows, and deployment procedures without scanning every line of source code.

---

## 🗺️ Documentation Navigation Map

Each documentation file has a clear, dedicated purpose and scope:

| Documentation File | Scope & Purpose |
|---|---|
| 📐 [**ARCHITECTURE.md**](file:///e:/evolipia-radar-1/docs/ARCHITECTURE.md) | **Complete System Architecture**: Comprehensive breakdown of Go backend (`cmd/`, `internal/`), Next.js frontend (`app/`, `src/`), crawler pipeline, background worker, SSE real-time stream, and scoring engine. |
| 🔌 [**API_REFERENCE.md**](file:///e:/evolipia-radar-1/docs/API_REFERENCE.md) | **REST API & SSE Specification**: Detailed specifications for HTTP endpoints (`GET /api/items`, `POST /api/crawl`, `POST /v2/chat`, `POST /v2/summarize`), query parameters, SSE event payload formats (`/api/crawl/progress`), and Prometheus `/metrics`. |
| 🗄️ [**DATABASE_SCHEMA.md**](file:///e:/evolipia-radar-1/docs/DATABASE_SCHEMA.md) | **Database Schema & ERD**: ERD diagram, Neon.tech PostgreSQL table structures (`sources`, `items`, `scores`, `summaries`, `signals`, `settings`), `pgvector` vector embedding support, and migration history `000001` - `000008`. |
| 🛠️ [**DEVELOPMENT_GUIDE.md**](file:///e:/evolipia-radar-1/docs/DEVELOPMENT_GUIDE.md) | **Setup & Dev Workflow**: Prerequisites, step-by-step local setup (Windows/Linux/macOS), Environment Variables reference table, running backend/worker/frontend, and unit testing. |
| 🚀 [**DEPLOYMENT_GUIDE.md**](file:///e:/evolipia-radar-1/docs/DEPLOYMENT_GUIDE.md) | **Production & Infrastructure**: Production deployment instructions for Go backend & worker on **Fly.io**, frontend on **Vercel**, Docker Compose, Kubernetes, and Terraform. |
| 📋 [**openapi.yaml**](file:///e:/evolipia-radar-1/docs/openapi.yaml) | **OpenAPI 3.0 Specification**: Machine-readable REST API specification in OpenAPI YAML format for SDK generators and Swagger UI. |
| 🤝 [**CONTRIBUTING.md**](file:///e:/evolipia-radar-1/docs/CONTRIBUTING.md) | **Contribution Guidelines**: Code style standards (Go & TypeScript), git flow, and pull request procedures. |

---

## ⚡ Quick System Overview

**Evolipia Radar** is an automated AI/ML research intelligence platform built with **Go 1.24+** (Backend & Ingestion Worker), **Next.js 15+** (React Frontend UI), and **PostgreSQL (Neon.tech)** with vector search capabilities (`pgvector`).

### Key Component Summary:
1. **Go API Server (`cmd/server`)**: Provides REST API endpoints, handles database connections, broadcasts real-time crawling progress via Server-Sent Events (SSE), and triggers manual crawl executions.
2. **Scheduled Worker (`cmd/worker`)**: Standalone background ingestion service running `robfig/cron/v3` on configurable schedules (default: every 6 hours) to harvest news from RSS & Trending agents.
3. **Crawler & Scoring Engine (`internal/crawler`, `internal/scoring`)**: Fetches content, validates URLs & text length, prevents duplicate entries via SHA-256 content hashing, and computes multi-tier relevance scores (Relevance, Impact, Engineering Value, Novelty).
4. **LLM Gateway Service (`internal/ai`, `internal/llm`)**: Provides a REST API proxy to OpenRouter / Gemini LLMs for interactive chat and auto-summarization (`tldr`, `why_it_matters`).
5. **Next.js Frontend (`app/`, `src/`)**: Modern UI featuring a dynamic Filter Bar (search, date picker, relevance slider, saved presets), live SSE Crawl Progress stepper modal, and Data Freshness indicators.

---

## 📂 Repository Directory Structure

```
evolipia-radar/
├── cmd/
│   ├── server/          # Main HTTP API Server & SSE Broadcaster
│   ├── worker/          # Background Scheduled Ingestion Worker
│   └── api/             # Standalone API runner
├── internal/
│   ├── ai/              # AI Gateway service & HTTP handlers (v2)
│   ├── crawler/         # Orchestrator, RSS agent, Trending agent, Validator, Retries
│   ├── db/              # GORM / SQL Database connection & helpers
│   ├── llm/             # OpenRouter & Gemini LLM HTTP Client
│   ├── models/          # Data Models, GORM structs, & DTOs
│   ├── scoring/         # Keyword & Multi-factor Relevance Scoring Engine
│   └── services/        # Business logic services (Feed, Source, Worker)
├── app/                 # Next.js App Router (pages, layout, globals.css, route proxies)
├── src/
│   └── components/      # React Components (FilterBar, CrawlProgress, DataFreshness)
├── migrations/          # SQL Migration files (000001 - 000008)
└── docs/                # Modular documentation suite (README, ARCHITECTURE, API, DB, DEV, DEPLOY)
```

---

## 🤖 Guidance for Future AI Sessions

If you are an AI assistant opening this repository:
- **Do not** scan all codebase files line-by-line to understand functional behavior.
- Read [ARCHITECTURE.md](file:///e:/evolipia-radar-1/docs/ARCHITECTURE.md) to understand end-to-end component interactions.
- Read [API_REFERENCE.md](file:///e:/evolipia-radar-1/docs/API_REFERENCE.md) for REST API contracts and SSE event structures.
- Read [DATABASE_SCHEMA.md](file:///e:/evolipia-radar-1/docs/DATABASE_SCHEMA.md) to understand PostgreSQL entity relationships.
