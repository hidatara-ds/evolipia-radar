# 🔌 Evolipia Radar - Complete API Reference & SSE Specification

This document provides complete API specifications for all HTTP REST endpoints and Server-Sent Events (SSE) streams in **Evolipia Radar**.

---

## 🌐 Overview & Base URL

- **Backend Base URL (Local)**: `http://localhost:8080`
- **Frontend Proxy Base URL**: `http://localhost:3000/api`
- **Content-Type**: `application/json` (Except SSE: `text/event-stream`)

---

## ⚠️ Standard Error Response Format

All error responses from the backend return HTTP status codes `4xx` or `5xx` with a uniform JSON payload:

```json
{
  "success": false,
  "code": "VALIDATION_ERROR",
  "error": "Detailed error message here"
}
```

### Key Error Codes:
- `VALIDATION_ERROR`: Invalid request parameters or malformed JSON payload.
- `NOT_FOUND`: The requested resource was not found.
- `INTERNAL_ERROR`: Internal server or database error.
- `LLM_PROVIDER_ERROR`: Failed response from external LLM provider.

---

## 📌 REST API Endpoints Reference

### 1. GET `/api/items` — Query & Filter Items
Fetches a list of validated and scored AI articles/news with pagination and advanced search & filtering capabilities.

- **HTTP Method**: `GET`
- **Query Parameters**:

| Parameter | Type | Default | Description |
|---|---|---|---|
| `search` | string | `""` | Search query (searches title, raw excerpt, and domain) |
| `date_from` | string | `""` | Start date filter (`YYYY-MM-DD`) |
| `date_to` | string | `""` | End date filter (`YYYY-MM-DD`) |
| `sources` | string | `""` | Comma-separated source IDs or names (`sources[]=id1&sources[]=id2`) |
| `categories` | string | `""` | Comma-separated category list (`llm`, `agents`, `vision`, `infra`) |
| `min_relevance` | integer | `0` | Minimum relevance score threshold (0-100) |
| `status` | string | `"all"` | Crawl status filter (`verified`, `pending`, `done`, `failed`, `all`) |
| `sort_by` | string | `"date"` | Field to sort by (`date`, `relevance`, `credibility`, `impact`) |
| `sort_order` | string | `"desc"` | Sort direction (`asc`, `desc`) |
| `page` | integer | `1` | Page number (1-indexed) |
| `limit` | integer | `20` | Page limit (1 - 100) |

- **Sample Request**:
  ```http
  GET /api/items?search=agents&min_relevance=50&sort_by=relevance&sort_order=desc&page=1&limit=10 HTTP/1.1
  Host: localhost:8080
  ```

- **Sample Response Payload (200 OK)**:
  ```json
  {
    "success": true,
    "data": [
      {
        "id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
        "title": "Latest Advances in Autonomous Agents Frameworks",
        "url": "https://example.com/ai/agents-advances",
        "source_name": "HackerNews RSS",
        "domain": "example.com",
        "category": "agents",
        "raw_excerpt": "Detailed technical breakdown of open-source agent frameworks...",
        "published_at": "2026-07-30T10:00:00Z",
        "crawl_status": "done",
        "relevance_score": 85,
        "scores": {
          "hot": 75.5,
          "relevance": 85.0,
          "credibility": 90.0,
          "novelty": 80.0,
          "impact": 88.0,
          "engineering_value": 92.0,
          "final": 86.5,
          "reasoning": "High relevance to autonomous agent architectures."
        },
        "summary": {
          "tldr": "Multi-agent framework benchmarked with 30% latency reduction.",
          "why_it_matters": "Provides standard design patterns for enterprise agent orchestration.",
          "tags": ["agents", "llm", "framework"]
        }
      }
    ],
    "page": 1,
    "limit": 10,
    "total_count": 125,
    "filtered_count": 45,
    "total_pages": 5
  }
  ```

---

### 2. POST `/api/crawl` — Trigger Manual Crawl Execution
Triggers an immediate background crawling process. Equipped with a **concurrency lock** to prevent concurrent double-runs.

- **HTTP Method**: `POST`
- **Request Body**: Empty JSON `{}` or optional source filter payload.
- **Sample Response (200 OK / 202 Accepted)**:
  ```json
  {
    "success": true,
    "message": "Crawl task triggered successfully",
    "timestamp": "2026-07-30T12:00:00Z"
  }
  ```
- **Conflict Response (409 Conflict)**:
  ```json
  {
    "success": false,
    "code": "CONCURRENCY_LOCK",
    "error": "A crawl execution is already in progress"
  }
  ```

---

### 3. GET `/api/crawl/progress` — Real-Time SSE Event Stream
Server-Sent Events (SSE) stream broadcasting real-time progress updates to frontend clients.

- **HTTP Method**: `GET`
- **Headers**:
  - `Accept: text/event-stream`
  - `Cache-Control: no-cache`
  - `Connection: keep-alive`

- **SSE Event Format Payload**:
  ```text
  event: progress
  data: {"step":2,"message":"Scanning source (1/4): HackerNews RSS...","progress":32,"current_source":"HackerNews RSS","total_sources":4,"processed_items":12,"estimated_remaining_secs":9,"timestamp":"2026-07-30T12:00:05Z","is_complete":false}
  ```

- **SSE Completion Event Payload**:
  ```text
  event: progress
  data: {"step":6,"message":"Done! 48 items processed successfully","progress":100,"processed_items":48,"timestamp":"2026-07-30T12:00:25Z","is_complete":true}
  ```

---

### 4. POST `/v2/chat` — LLM Chat Completion Gateway
Gateway endpoint for interactive chat completions powered by OpenRouter / Gemini API.

- **HTTP Method**: `POST`
- **Request Body**:
  ```json
  {
    "messages": [
      {"role": "system", "content": "You are an AI news analyst."},
      {"role": "user", "content": "Explain recent trends in open-source LLMs."}
    ],
    "model": "google/gemini-flash-1.5",
    "temperature": 0.7
  }
  ```
- **Response Payload (200 OK)**:
  ```json
  {
    "success": true,
    "message": {
      "role": "assistant",
      "content": "Recent trends in open-source LLMs focus on mixture-of-experts (MoE)..."
    },
    "model_used": "google/gemini-flash-1.5",
    "usage": {
      "prompt_tokens": 42,
      "completion_tokens": 128,
      "total_tokens": 170
    }
  }
  ```

---

### 5. POST `/v2/summarize` — LLM Text Summarization
Compresses long article excerpts into concise TLDR and "Why It Matters" sections.

- **HTTP Method**: `POST`
- **Request Body**:
  ```json
  {
    "text": "Full text of the AI news article or paper excerpt...",
    "instructions": "Focus on benchmark performance and architecture changes."
  }
  ```
- **Response Payload (200 OK)**:
  ```json
  {
    "success": true,
    "tldr": "New model architecture improves inference speed by 40%.",
    "why_it_matters": "Lowers latency for edge computing AI deployments.",
    "tags": ["performance", "inference", "edge-ai"]
  }
  ```

---

### 6. GET `/metrics` — Prometheus Telemetry Metrics
Exposes application performance metrics for Prometheus scraping.

- **HTTP Method**: `GET`
- **Content-Type**: `text/plain; version=0.0.4`

---

## 📄 OpenAPI 3.0 Specification

All endpoints listed above are also documented in OpenAPI 3.0 YAML format at:
[docs/openapi.yaml](file:///e:/evolipia-radar-1/docs/openapi.yaml)
