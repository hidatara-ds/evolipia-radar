# 🛠️ Evolipia Radar - Development & Local Setup Guide

This document provides step-by-step instructions for developers to configure environment variables, manage dependencies, run the Go backend API server, background worker, PostgreSQL database, and Next.js frontend UI in a local environment.

---

## 📋 Prerequisites

Before starting, ensure you have the following software installed:

1. **Go 1.24+**: Download from [go.dev](https://go.dev/dl/). Verify with `go version`.
2. **Node.js 18+ & npm**: Download from [nodejs.org](https://nodejs.org/). Verify with `node -v` and `npm -v`.
3. **PostgreSQL 15+** or **Neon.tech Account**: Use local Docker or Neon.tech cloud.
4. **golang-migrate CLI** *(optional, for running manual SQL migrations)*:
   ```bash
   go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
   ```

---

## ⚙️ Environment Variables Reference Table

Create a `.env` file in the project root based on `.env.example`:

| Variable Name | Required | Default Value | Description |
|---|---|---|---|
| `PORT` | No | `8080` | Port for the Go HTTP API Server |
| `DATABASE_URL` | No | `postgres://postgres:postgres@localhost:5432/radar?sslmode=disable` | Connection string for PostgreSQL / Neon.tech |
| `CRAWL_INTERVAL` | No | `@every 6h` | Auto-scheduler interval (`@every 6h` or cron expression `0 */6 * * *`) |
| `WORKER_CRON` | No | `*/10 * * * *` | Cron schedule for standalone worker (`cmd/worker`) |
| `MIN_RELEVANCE_SCORE` | No | `30` | Minimum score (0-100) required to save items to DB |
| `TOPICS_KEYWORDS` | No | `llm,agents,vision,open source,infra,robotics,security,ai` | Comma-separated keyword list for relevance scoring |
| `MAX_CRAWL_RETRIES` | No | `3` | Maximum retry attempts for failed HTTP source fetches |
| `FETCH_TIMEOUT_SECONDS` | No | `8` | HTTP client timeout for news retrieval (seconds) |
| `MAX_FETCH_BYTES` | No | `2000000` | Maximum allowed payload size for source responses (bytes) |
| `LLM_ENABLED` | No | `false` | Enable/disable LLM integration features |
| `LLM_PROVIDER` | No | `openrouter` | LLM Provider (`openrouter`, `gemini`, `openai`) |
| `LLM_MODEL` | No | `google/gemini-flash-1.5` | Default primary LLM model identifier |
| `LLM_API_KEY` | No | `""` | API Key for OpenRouter / LLM Provider |

---

## 🚀 Step-by-Step Local Execution Workflow

### 1. Clone & Install Dependencies
```bash
git clone https://github.com/hidatara-ds/evolipia-radar.git
cd evolipia-radar

# Download Go backend dependencies
go mod download

# Install Frontend npm packages
npm install
```

### 2. Setup Local Database (Via Docker)
If not using Neon.tech cloud, spin up local PostgreSQL using Docker Compose:
```bash
docker-compose up -d postgres
```

### 3. Run Database Migrations
Execute SQL migrations to create database schema and tables:
```bash
# Using migrate CLI:
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/radar?sslmode=disable" up
```

### 4. Run Go Backend API Server (`cmd/server`)
Starts the HTTP API server on `http://localhost:8080`:
```bash
go run cmd/server/main.go
```
*Output*:
```text
{"time":"...","level":"INFO","msg":"Starting Evolipia Radar Server..."}
{"time":"...","level":"INFO","msg":"Database connected successfully"}
{"time":"...","level":"INFO","msg":"HTTP Server listening on port 8080"}
```

### 5. Run Scheduled Ingestion Worker (`cmd/worker`) — *Optional*
To run a dedicated background worker process for automated crawling:
```bash
go run cmd/worker/main.go
```

### 6. Run Next.js Frontend UI
Start Next.js development server on `http://localhost:3000`:
```bash
npm run dev
```

---

## 🧪 Testing & Verification Commands

### Run Go Backend Unit Tests:
```bash
go test -v ./internal/... ./pkg/...
```

### Run Specific Test Suite:
```bash
go test -v ./internal/crawler -run TestValidator
```

### Run Frontend ESLint & Type-Check:
```bash
npm run lint
```

---

## 🪟 Windows Specific Notes

For Windows PowerShell users:
- If script execution is blocked by policies, run: `Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass`
- Execute automated setup scripts if available: `./setup-windows.ps1`
