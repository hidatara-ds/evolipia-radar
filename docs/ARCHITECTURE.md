# 📐 Evolipia Radar - Comprehensive System Architecture

This document provides in-depth technical specifications for the system architecture of **Evolipia Radar**, covering the Go backend services, background ingestion worker, crawler pipeline, real-time SSE progress streaming, LLM integration, and Next.js frontend UI.

---

## 🏗️ High-Level System Architecture

Below is the overall architecture and end-to-end data flow within the Evolipia Radar platform:

```mermaid
graph TD
    subgraph External Sources
        RSS[RSS Feeds]
        ArXiv[ArXiv AI Papers]
        HN[HackerNews API]
        Reddit[Reddit AI/ML Subreddits]
    end

    subgraph Go Backend Services
        Worker[Scheduled Ingestion Worker\ncmd/worker]
        Server[HTTP API & SSE Server\ncmd/server]
        Orchestrator[Crawler Orchestrator\ninternal/crawler]
        Validator[Data Validator & Scorer\ninternal/crawler/validator.go]
        ScoringEngine[Scoring Engine\ninternal/scoring]
        SSEBroadcaster[SSE Progress Broadcaster\ninternal/api/sse.go]
        LLMClient[LLM Gateway / Client\ninternal/llm]
    end

    subgraph Storage & Cloud
        Postgres[(Neon.tech PostgreSQL\n+ pgvector)]
        OpenRouter[OpenRouter / Gemini API]
    end

    subgraph Frontend Next.js UI
        WebUI[Next.js 15 Web App\napp/page.tsx]
        FilterBar[Advanced FilterBar Component]
        CrawlProgress[CrawlProgress Component\n(SSE Listener)]
        DataFreshness[DataFreshness Badge]
    end

    %% Flow Connections
    Worker -->|Triggers Schedule| Orchestrator
    Server -->|Manual POST /api/crawl| Orchestrator
    Orchestrator -->|Fetches Data| RSS
    Orchestrator -->|Fetches Data| ArXiv
    Orchestrator -->|Fetches Data| HN
    Orchestrator -->|Fetches Data| Reddit

    Orchestrator --> Validator
    Validator --> ScoringEngine
    Validator -->|Broadcasting Step Events| SSEBroadcaster
    SSEBroadcaster -->|SSE Stream /api/crawl/progress| CrawlProgress

    Validator -->|Upsert Valid Items| Postgres
    ScoringEngine -->|Calculates Multi-tier Score| Postgres

    Server -->|Query GET /api/items| Postgres
    WebUI -->|HTTP Requests| Server
    LLMClient -->|Generate TLDR / Chat| OpenRouter
```

---

## 📦 Core Component Breakdown

### 1. HTTP API Server (`cmd/server/main.go`)
- **Primary Function**: Serves REST API endpoints for searching and verifying news, handles manual crawl triggers (`POST /api/crawl`), and streams live crawl progress via Server-Sent Events (`GET /api/crawl/progress`).
- **Framework**: `gin-gonic/gin`.
- **Logger**: Structured JSON logging via `log/slog`.
- **Graceful Shutdown**: Listens for `SIGINT` / `SIGTERM` signals using `context.WithTimeout` (10s) to cleanly close database pools and HTTP connections.

### 2. Scheduled Ingestion Worker (`cmd/worker/main.go`)
- **Primary Function**: Standalone background worker process that automatically executes periodic news ingestion.
- **Scheduler**: Powered by `github.com/robfig/cron/v3` with configurable cron expressions via `WORKER_CRON` (default: `*/10 * * * *`).
- **Initial Execution**: Performs a single ingestion run upon startup before adhering to regular cron schedules.

### 3. Crawler Pipeline & Agents (`internal/crawler/`)
- **Orchestrator (`orchestrator.go`)**: Coordinates individual crawl agents, enforces rate limits, and tracks source health.
- **RSS Agent (`rss_agent.go`)**: Fetches and parses XML RSS feeds (TechCrunch AI, VentureBeat, etc.).
- **Trending Agent (`trending_agent.go`)**: Scrapes popular entries from HackerNews, Reddit ML, and ArXiv papers.
- **Validator & Relevance Scorer (`validator.go`)**:
  - Validates URL protocols (`http://` or `https://`).
  - Ensures minimum title length (10 chars) and non-future publication dates.
  - Generates SHA-256 hashes of Title + URL to prevent duplicate database insertions.
  - Performs keyword matching relevance scoring (0-100) based on `TOPICS_KEYWORDS`, rejecting candidates below `MIN_RELEVANCE_SCORE` (default: 30).
- **Retry Runner (`retry.go`)**: Implements exponential backoff retries up to `MAX_CRAWL_RETRIES` (default: 3) to handle transient network issues.

### 4. Real-Time SSE Broadcaster (`internal/api/sse.go`)
- **Primary Function**: Broadcasts real-time step progress events to frontend clients over Server-Sent Events (SSE).
- **Thread Safety**: Uses `sync.RWMutex` to safely manage concurrent active client connections.
- **Progress Stepper Steps**:
  1. `Initializing crawler...` (10%)
  2. `Scanning sources (1/N)...` (20-70%)
  3. `Validating data & scoring...` (80%)
  4. `Saving to database...` (90%)
  5. `Done! X items processed` (100%)

### 5. Multi-Tier Scoring & LLM Engine (`internal/scoring/`, `internal/llm/`)
- **Scoring Algorithm**:
  - **Relevance**: Topic keyword alignment with AI/ML domains.
  - **Impact**: Domain reputation and reach weighting.
  - **Engineering Value**: Technical content density (code, papers, benchmarks, tutorials).
  - **Novelty**: Recency and breakthrough announcement factor.
  - **Final Score**: Weighted combination of all four sub-scores.
- **LLM Client (`internal/llm/client.go`)**: Communicates with **OpenRouter API** or **Google Gemini API** (`google/gemini-flash-1.5`) to produce automated concise summaries (`TLDR`) and key takeaways (`Why it matters`).

### 6. Next.js Frontend Architecture (`app/`, `src/components/`)
- **Framework**: Next.js 15 App Router + React + Vanilla/Tailwind CSS.
- **Key UI Components**:
  - **`FilterBar.tsx`**: Renders real-time debounced search (300ms), date range pickers, source & category multi-selects, relevance slider (0-100%), and saved preset management backed by `localStorage`. Synchronized bi-directionally with URL query parameters.
  - **`CrawlProgress.tsx`**: Renders real-time progress modals from the `/api/crawl/progress` SSE stream, progress percentage bar, estimated time remaining, and current source indicator.
  - **`DataFreshness.tsx`**: Displays "Last crawled X minutes ago" status with color-coded fresh badges (Green < 6h, Yellow 6-24h, Red > 24h).

---

## 🔄 Data Lifecycle (Pipeline Flow)

```
[Raw Source Feed / Web API]
         │
         ▼
[Crawler Agent Ingestion] ──► Fetch content excerpt & metadata
         │
         ▼
[SHA-256 Hash Normalization] ──► Check duplicate content_hash in DB
         │
         ▼
[Data Validation Layer] ──► Check URL protocol, date bounds, text lengths
         │
         ▼
[Keyword Relevance Scoring] ──► Calculate 0-100 score (Reject if < threshold)
         │
         ▼
[Database Persistence] ──► Upsert to PostgreSQL (`sources`, `items`, `scores`)
         │
         ▼
[LLM Summarization Job] ──► Generate TLDR, why_it_matters, tags, & vector embedding
         │
         ▼
[Frontend UI Rendering] ──► Display on Next.js UI via `/api/items` query
```
