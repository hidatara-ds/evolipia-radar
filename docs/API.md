# API Documentation

## Base URL

- **Production:** `https://evolipia-radar.vercel.app`
- **Local Development:** `http://localhost:3000`

## Authentication

Currently, all endpoints are public and do not require authentication.

## Endpoints

### Get News

Retrieve latest news articles with optional topic filtering.

```http
GET /api/news?topic={topic}
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `topic` | string | No | Filter by topic (llm, vision, data, security, rl, robotics, ide, free-credits) |

**Response:**

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "title": "Article Title",
        "url": "https://example.com/article",
        "domain": "example.com",
        "published_at": "2026-03-18T12:00:00Z",
        "category": "tech",
        "score": 0.85,
        "tldr": "Brief summary of the article",
        "why_it_matters": "Explanation of significance",
        "tags": ["llm", "tools"]
      }
    ],
    "total_count": 20,
    "last_updated": "2026-03-18T13:00:00Z"
  }
}
```

**Example Requests:**

```bash
# Get all news
curl https://evolipia-radar.vercel.app/api/news

# Get LLM news only
curl https://evolipia-radar.vercel.app/api/news?topic=llm

# Get Vision news only
curl https://evolipia-radar.vercel.app/api/news?topic=vision
```

**Status Codes:**

- `200 OK` - Success
- `500 Internal Server Error` - Server error

---

### Get Trending

Retrieve trending articles based on engagement signals.

```http
GET /api/trending
```

**Response:**

```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "uuid",
        "title": "Trending Article",
        "url": "https://example.com/trending",
        "domain": "example.com",
        "published_at": "2026-03-18T10:00:00Z",
        "category": "tech",
        "score": 0.92,
        "tldr": "Summary",
        "why_it_matters": "Significance",
        "tags": ["llm"]
      }
    ],
    "total_count": 10
  }
}
```

---

### Search Articles

Search articles by keyword.

```http
GET /api/search?q={query}&topic={topic}&limit={limit}&offset={offset}
```

**Query Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `q` | string | Yes | Search query |
| `topic` | string | No | Filter by topic |
| `limit` | integer | No | Results per page (default: 20, max: 100) |
| `offset` | integer | No | Pagination offset (default: 0) |

**Response:**

```json
{
  "success": true,
  "data": {
    "items": [...],
    "total_count": 45,
    "query": "transformer",
    "limit": 20,
    "offset": 0
  }
}
```

**Example:**

```bash
curl "https://evolipia-radar.vercel.app/api/search?q=transformer&topic=llm&limit=10"
```

---

### Get Metrics

Retrieve system metrics and statistics.

```http
GET /metrics
```

**Response:**

```json
{
  "articles_processed": 1250,
  "filtered_articles": 850,
  "api_hits": 5420,
  "clusters": 12,
  "avg_cluster_score": 0.75,
  "top_cluster_titles": [
    "GPT-5 Release",
    "New Vision Model",
    "RL Breakthrough"
  ]
}
```

---

### Health Check

Check API health status.

```http
GET /healthz
```

**Response:**

```json
{
  "status": "healthy",
  "timestamp": "2026-03-18T13:00:00Z",
  "version": "2.0.0"
}
```

---

### Trigger Crawl

Manually trigger a news crawl cycle.

```http
POST /v2/crawl/trigger
```

**Response:**

```json
{
  "success": true,
  "stats": {
    "discovered": 25,
    "inserted": 18,
    "duplicates": 7
  },
  "duration_ms": 12500
}
```

**Note:** This endpoint may have rate limiting in production.

---

## Error Responses

All endpoints return errors in the following format:

```json
{
  "success": false,
  "error": "Error message description"
}
```

**Common Error Messages:**

- `"Database configuration missing"` - DATABASE_URL not set
- `"Failed to connect to database"` - Database connection error
- `"Failed to load news"` - Query execution error
- `"Database connection timeout"` - Connection timeout (cold start)

---

## Rate Limiting

Currently, no rate limiting is enforced. Future versions may implement:

- 100 requests/minute per IP for read endpoints
- 10 requests/minute for write endpoints (trigger crawl)

---

## CORS

CORS is enabled for all origins:

```
Access-Control-Allow-Origin: *
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Headers: Content-Type
```

---

## Data Models

### NewsItem

```typescript
interface NewsItem {
  id: string;              // UUID
  title: string;           // Article title
  url: string;             // Article URL
  domain: string;          // Source domain
  published_at: string;    // ISO 8601 timestamp
  category: string;        // Category (e.g., "tech")
  score: number;           // Relevance score (0-1)
  tldr?: string;           // AI-generated summary
  why_it_matters?: string; // Significance explanation
  tags: string[];          // Topic tags
}
```

### Response

```typescript
interface Response<T> {
  success: boolean;
  data?: T;
  error?: string;
}
```

---

## Best Practices

### Pagination

For large result sets, use `limit` and `offset`:

```bash
# Page 1
curl "/api/search?q=llm&limit=20&offset=0"

# Page 2
curl "/api/search?q=llm&limit=20&offset=20"

# Page 3
curl "/api/search?q=llm&limit=20&offset=40"
```

### Caching

Responses can be cached for up to 30 seconds:

```
Cache-Control: public, max-age=30
```

### Error Handling

Always check the `success` field:

```javascript
const response = await fetch('/api/news');
const data = await response.json();

if (data.success) {
  // Handle success
  console.log(data.data.items);
} else {
  // Handle error
  console.error(data.error);
}
```

### Retry Logic

Implement exponential backoff for failed requests:

```javascript
async function fetchWithRetry(url, retries = 3) {
  for (let i = 0; i < retries; i++) {
    try {
      const response = await fetch(url);
      if (response.ok) return await response.json();
    } catch (error) {
      if (i === retries - 1) throw error;
      await new Promise(r => setTimeout(r, 2000 * (i + 1)));
    }
  }
}
```

---

## SDK Examples

### JavaScript/TypeScript

```typescript
class EvolipiaRadarClient {
  constructor(private baseUrl: string = 'https://evolipia-radar.vercel.app') {}

  async getNews(topic?: string): Promise<NewsItem[]> {
    const url = topic 
      ? `${this.baseUrl}/api/news?topic=${topic}`
      : `${this.baseUrl}/api/news`;
    
    const response = await fetch(url);
    const data = await response.json();
    
    if (!data.success) {
      throw new Error(data.error);
    }
    
    return data.data.items;
  }

  async search(query: string, options?: {
    topic?: string;
    limit?: number;
    offset?: number;
  }): Promise<{ items: NewsItem[]; total: number }> {
    const params = new URLSearchParams({
      q: query,
      ...options
    });
    
    const response = await fetch(`${this.baseUrl}/api/search?${params}`);
    const data = await response.json();
    
    if (!data.success) {
      throw new Error(data.error);
    }
    
    return {
      items: data.data.items,
      total: data.data.total_count
    };
  }
}

// Usage
const client = new EvolipiaRadarClient();
const llmNews = await client.getNews('llm');
const searchResults = await client.search('transformer', { topic: 'llm' });
```

### Python

```python
import requests
from typing import List, Dict, Optional

class EvolipiaRadarClient:
    def __init__(self, base_url: str = "https://evolipia-radar.vercel.app"):
        self.base_url = base_url
    
    def get_news(self, topic: Optional[str] = None) -> List[Dict]:
        url = f"{self.base_url}/api/news"
        params = {"topic": topic} if topic else {}
        
        response = requests.get(url, params=params)
        data = response.json()
        
        if not data["success"]:
            raise Exception(data["error"])
        
        return data["data"]["items"]
    
    def search(self, query: str, topic: Optional[str] = None, 
               limit: int = 20, offset: int = 0) -> Dict:
        url = f"{self.base_url}/api/search"
        params = {
            "q": query,
            "limit": limit,
            "offset": offset
        }
        if topic:
            params["topic"] = topic
        
        response = requests.get(url, params=params)
        data = response.json()
        
        if not data["success"]:
            raise Exception(data["error"])
        
        return {
            "items": data["data"]["items"],
            "total": data["data"]["total_count"]
        }

# Usage
client = EvolipiaRadarClient()
llm_news = client.get_news(topic="llm")
results = client.search("transformer", topic="llm", limit=10)
```

---

## Changelog

### v2.0.0 (2026-03-18)
- Migrated to Neon.tech PostgreSQL
- Added auto-retry for cold starts
- Improved error handling
- Added topic filtering

### v1.0.0 (2026-03-01)
- Initial API release
- Basic news endpoints
- Metrics endpoint

---

## Support

For API issues or questions:
- GitHub Issues: https://github.com/hidatara-ds/evolipia-radar/issues
- Documentation: https://github.com/hidatara-ds/evolipia-radar/tree/main/docs
