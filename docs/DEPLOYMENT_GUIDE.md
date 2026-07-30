# 🚀 Evolipia Radar - Production Infrastructure & Deployment Guide

This document provides a production deployment guide for **Evolipia Radar**, covering options for hosting on **Fly.io** (Go Backend & Worker), **Vercel** (Next.js Frontend), **Docker Compose**, **Kubernetes**, and **Terraform**.

---

## 🏗️ Production Architecture Overview

```
[Vercel Cloud] ───────────► Next.js Frontend (Edge Network)
                                 │
                                 │ HTTP API & SSE Stream Requests
                                 ▼
[Fly.io Platform] ────────► Go Backend API (`fly.toml`) ──┐
                         ► Scheduled Worker (`fly.worker.toml`) ├─► [Neon.tech Serverless Postgres]
                                                          │
[OpenRouter / Gemini] ◄── LLM Service Proxy ──────────────┘
```

---

## 🚁 Option 1: Deployment on Fly.io (Go Backend & Worker)

The Go backend server and background ingestion worker can be deployed independently on **Fly.io** while sharing the same Neon.tech PostgreSQL database.

### 1. Deploy HTTP API Server (`cmd/server`)
Configuration file `fly.toml`:
```toml
app = "evolipia-radar-api"
primary_region = "sin"

[build]
  dockerfile = "Dockerfile.api"

[env]
  PORT = "8080"
  CRAWL_INTERVAL = "@every 6h"
  MIN_RELEVANCE_SCORE = "30"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    force_https = true
    handlers = ["http", "tls"]
    port = 443
```

- **Set Secret Environment Variables**:
  ```bash
  fly secrets set DATABASE_URL="postgres://user:password@ep-cool-neon.singapore.aws.neon.tech/radar?sslmode=require" LLM_API_KEY="sk-or-v1-your-key"
  ```
- **Deploy Command**:
  ```bash
  fly deploy --config fly.toml
  ```

### 2. Deploy Scheduled Ingestion Worker (`cmd/worker`)
Configuration file `fly.worker.toml`:
```toml
app = "evolipia-radar-worker"
primary_region = "sin"

[build]
  dockerfile = "Dockerfile.worker"

[env]
  WORKER_CRON = "0 */6 * * *"
```

- **Set Secret & Deploy Worker**:
  ```bash
  fly secrets set --app evolipia-radar-worker DATABASE_URL="postgres://user:password@ep-cool-neon.singapore.aws.neon.tech/radar?sslmode=require"
  fly deploy --config fly.worker.toml
  ```

---

## 🔺 Option 2: Deployment on Vercel (Next.js Frontend)

The Next.js frontend is deployed on Vercel with path rewrites configured in `vercel.json`:

```json
{
  "version": 2,
  "buildCommand": "npm run build",
  "outputDirectory": ".next",
  "framework": "nextjs",
  "rewrites": [
    {
      "source": "/api/:path*",
      "destination": "https://evolipia-radar-api.fly.dev/api/:path*"
    }
  ]
}
```

### Steps:
1. Connect the GitHub repository to your Vercel Dashboard.
2. Configure Environment Variables on Vercel:
   - `NEXT_PUBLIC_API_BASE_URL` = `https://evolipia-radar-api.fly.dev`
3. Click **Deploy**. Vercel will build and deploy the Next.js application across its global edge network.

---

## 🐳 Option 3: Containerized Deployment (Docker & Docker Compose)

Utilize Docker to run the complete stack locally or on a single VPS:

### 1. Build Multi-stage Docker Images
- `Dockerfile.api`: Containerized Go HTTP API Server.
- `Dockerfile.worker`: Containerized Go Scheduled Worker.

### 2. Run via Docker Compose
```bash
# Start API Server, Worker, Postgres & Prometheus Observability
docker-compose -f docker-compose.yml -f docker-compose.observability.yml up -d --build
```

---

## ☸️ Option 4: Kubernetes & Terraform Infrastructure (`k8s/`, `terraform/`)

Infrastructure as Code (IaC) manifests are available in the repository:
- **`terraform/`**: GCP resource provisioning, Cloud SQL, and GKE cluster definitions.
- **`k8s/`**: Kubernetes Deployments, Services, ConfigMaps, and Ingress controller manifests.

---

## 🩺 Production Health & Monitoring

- **API Server Health Check**: `GET /api/items?limit=1`
- **Prometheus Telemetry Scrape**: `GET /metrics`
- **Live SSE Crawl Monitoring**: `GET /api/crawl/progress`
