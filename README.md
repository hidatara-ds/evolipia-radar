# Evolipia Radar: AI-Powered News Intelligence & Auto-Crawler Platform

> AI Research Intelligence Platform -- Go Backend + Next.js Frontend + PostgreSQL

[![Go](https://img.shields.io/badge/Go-1.24.1-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-15.5-000000?logo=next.js&logoColor=white)](https://nextjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org/)

Evolipia Radar is an automated news intelligence platform that discovers AI/ML news, validates and scores content relevance, streams real-time crawl progress using Server-Sent Events (SSE), and provides advanced search and filtering capabilities.

---

## 🚀 Main Features & Recent Upgrades

### 1. AutoCrawl & Auto-Scheduler
- **Cron Scheduler**: Automatic background crawling executed via `github.com/robfig/cron/v3` on server boot.
- **Configurable Interval**: Default run interval of every 6 hours (configurable via `CRAWL_INTERVAL`).
- **Emergency Trigger**: Endpoint `POST /api/crawl` for manual crawl execution.
- **Concurrency Guard**: Atomic lock prevents concurrent double-runs.
- **Graceful Shutdown**: Stops cron workers cleanly on SIGINT/SIGTERM using `sync.WaitGroup`.

### 2. Real-Time Crawl Progress & UI Indicators
- **SSE Stream**: Server-Sent Events endpoint `GET /api/crawl/progress` broadcasting step progress.
- **Step Progress Indicator**:
  1. `Initializing crawler...`
  2. `Scanning sources (1/N)...`
  3. `Parsing content from [source_name]...`
  4. `Validating data...`
  5. `Saving to database...`
  6. `Done! X items processed`
- **UI Components**: `CrawlProgress.tsx` displaying step stepper, percentage bar, source status, estimated remaining time, and toast notifications.
- **Data Freshness Badge**: `DataFreshness.tsx` showing "Last crawled: X minutes ago" with color coding (Green < 6h, Yellow 6-24h, Red > 24h).

### 3. Data Quality, Validation & Retry Mechanism
- **Validation Layer**: Rejects invalid candidates lacking required title length (min 10 chars), URL protocol (`http://` or `https://`), non-future publication date, or excerpt length (min 50 chars).
- **Duplicate Detection**: Hashes normalized title and URL to prevent duplicate inserts.
- **Relevance Scoring**: Topic-based keyword matching algorithm (0-100 score) rejecting low-relevance items (< threshold, default 30).
- **Exponential Backoff Retries**: Retries failed sources up to 3 times (1s, 2s, 4s) before marking as unhealthy.
- **Database Schema**: New migration `migrations/000008_add_crawl_fields.up.sql` adding `crawl_status`, `crawl_error`, `relevance_score`, and `validated_at`.

### 4. Advanced Sort & Filter Bar
- **Debounced Search**: 300ms real-time search across titles, content excerpts, and domains.
- **Date Range Selector**: Today, Last 7 Days, Last 30 Days, or Custom Range.
- **Source & Category Multi-Select**: Checkboxes and category tag filters.
- **Relevance Slider**: Dynamic min-max relevance threshold filter (0-100%).
- **Saved Presets**: Save custom filter configurations to `localStorage` with custom names.
- **URL Synchronization**: Two-way state sync with browser URL query parameters.

---

## 🛠️ Environment Variables List

| Variable Name | Required | Default Value | Description |
|---|---|---|---|
| `PORT` | Optional | `8080` | Backend API server port |
| `DATABASE_URL` | Optional | `postgres://postgres:postgres@localhost:5432/radar?sslmode=disable` | PostgreSQL connection string |
| `CRAWL_INTERVAL` | Optional | `@every 6h` | Auto-scheduler interval (`@every 6h`, `0 */6 * * *`) |
| `MIN_RELEVANCE_SCORE` | Optional | `30` | Minimum score threshold (0-100) to save items |
| `TOPICS_KEYWORDS` | Optional | `llm,agents,vision,open source,infra,robotics,security` | Comma-separated keyword list for relevance scoring |
| `MAX_CRAWL_RETRIES` | Optional | `3` | Maximum retry attempts for failed sources |
| `LLM_API_KEY` | Optional | `""` | OpenRouter / LLM API key |
| `LLM_PROVIDER` | Optional | `openrouter` | LLM provider identifier |
| `LLM_MODEL` | Optional | `google/gemini-flash-1.5` | Default LLM model string |

---

## 📂 Backend & Frontend Directory Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go          # Main entry point with auto-scheduler, SSE stream, and API routes
├── internal/
│   ├── crawler/
│   │   ├── crawler.go       # Crawler core logic
│   │   ├── scheduler.go     # robfig/cron/v3 auto-scheduler & concurrency lock
│   │   ├── validator.go     # Item validation & 0-100 relevance scoring
│   │   └── retry.go         # Exponential backoff retry & source health tracker
│   ├── api/
│   │   ├── handler.go       # HTTP handlers
│   │   ├── items.go         # GET /api/items advanced filter & pagination handler
│   │   ├── middleware.go    # CORS, logging, recovery
│   │   └── sse.go           # SSE progress stream broadcaster
│   ├── models/
│   │   └── item.go          # Item DB models & Progress DTOs
│   └── config/
│       └── config.go        # Env config loading and validation
├── pkg/
│   └── utils/
│       └── utils.go         # SHA256 hashing, normalization, URL validation
└── go.mod

frontend/
├── src/
│   ├── components/
│   │   ├── CrawlProgress.tsx    # Real-time progress UI & step stepper
│   │   ├── FilterBar.tsx        # Advanced filter bar & preset manager
│   │   └── DataFreshness.tsx    # Freshness badge indicator
│   ├── hooks/
│   │   ├── useCrawlProgress.ts  # SSE hook for /api/crawl/progress
│   │   └── useFilters.ts        # Filter state & URL query sync hook
│   └── api/
│       └── client.ts            # API client wrapper
├── app/
│   └── page.tsx                 # Next.js main dashboard
└── package.json
```

---

## 💻 How to Run & Test

### 1. Run Backend Server
```bash
go run ./cmd/server
```
Backend API will start on `http://localhost:8080`.

### 2. Run Backend Formatting & Unit Tests
```bash
gofmt -w .
go vet ./...
go test ./...
```

### 3. Run Frontend Development Server
```bash
npm run dev
```
Access UI at `http://localhost:3000`.

### 4. Build Frontend for Production
```bash
npm run build
```

---

## 📜 License
MIT License. See `LICENSE.md`.
