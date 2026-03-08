# Phase 2 Scaffold - Intelligent Crawling & Real-time Features

## Overview
Phase 2 adds advanced crawling capabilities, vector search, and real-time updates. Scaffold files are created with interfaces and TODOs.

## Scaffolded Components

### 1. Intelligent Crawler (`internal/crawler/crawler.go`)

#### Features to Implement
- **Headless Browser Support**
  - Use [rod](https://github.com/go-rod/rod) or [chromedp](https://github.com/chromedp/chromedp)
  - Render JavaScript-heavy sites (OpenAI blog, etc.)
  - Wait for dynamic content loading

- **Content Extraction**
  - Integrate [go-readability](https://github.com/go-shiori/go-readability)
  - Extract clean article text
  - Remove ads, navigation, footers
  - Preserve article structure

- **Smart Features**
  - Respect robots.txt
  - Adaptive rate limiting per domain
  - Proxy rotation support
  - Exponential backoff retry
  - Fingerprint randomization

#### Implementation Steps
```bash
# 1. Add dependencies
go get github.com/go-rod/rod
go get github.com/go-shiori/go-readability

# 2. Implement Crawl() method
# - Launch headless browser
# - Navigate and wait for content
# - Extract using readability
# - Return CrawlResult

# 3. Add to worker
# - Use crawler for full-text extraction
# - Enhance summaries with full content
```

#### Usage Example
```go
crawler := crawler.NewCrawler(crawler.Config{
    Headless:         true,
    MaxConcurrent:    5,
    RespectRobotsTxt: true,
    Timeout:          30 * time.Second,
})

result, err := crawler.Crawl(ctx, "https://openai.com/blog/...")
if err != nil {
    log.Printf("Crawl failed: %v", err)
}

// Use result.Content for better summarization
```

### 2. Vector Search (`internal/search/vector.go`)

#### Features to Implement
- **pgvector Integration**
  - Add vector column to items table
  - Create IVFFlat index for fast similarity search
  - Store 1536-dimensional embeddings (OpenAI compatible)

- **Embedding Generation**
  - Use OpenRouter embeddings API
  - Batch processing for efficiency
  - Cache embeddings

- **Search Capabilities**
  - Semantic search by query
  - Find similar articles
  - Topic clustering

#### Implementation Steps
```bash
# 1. Install pgvector extension
# In PostgreSQL:
CREATE EXTENSION vector;

# 2. Add migration
migrate create -ext sql -dir migrations -seq add_vector_search

# 3. Migration up:
ALTER TABLE items ADD COLUMN embedding vector(1536);
CREATE INDEX ON items USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

# 4. Implement IndexItem()
# - Call OpenRouter embeddings API
# - Store in embedding column

# 5. Implement Search()
# - Generate query embedding
# - Use <=> operator for cosine similarity
# - Return top K results
```

#### Usage Example
```go
vs := search.NewVectorSearch()

// Index new item
err := vs.IndexItem(ctx, itemID, item.Title + " " + item.Content)

// Semantic search
results, err := vs.Search(ctx, "RAG systems with LLMs", 10)

// Find similar articles
similar, err := vs.FindSimilar(ctx, itemID, 5)
```

#### API Endpoints to Add
```go
// GET /v1/search/semantic?q=query&limit=10
// GET /v1/items/:id/similar?limit=5
```

### 3. Real-time Updates (`internal/realtime/websocket.go`)

#### Features to Implement
- **WebSocket Server**
  - Use [gorilla/websocket](https://github.com/gorilla/websocket)
  - Hub pattern for connection management
  - Topic-based subscriptions

- **Event Types**
  - `new_item`: New high-scored item added
  - `rising_item`: Item momentum increasing
  - `trending_topic`: Topic gaining traction

- **Client Features**
  - Subscribe to specific topics
  - Reconnection logic
  - Heartbeat/ping-pong

#### Implementation Steps
```bash
# 1. Add dependency
go get github.com/gorilla/websocket

# 2. Implement Hub.Run()
# - Handle registration/unregistration
# - Broadcast to subscribed clients
# - Topic filtering

# 3. Add WebSocket endpoint
# GET /v1/ws - Upgrade to WebSocket

# 4. Integrate with worker
# - Broadcast when new items added
# - Broadcast when rising items detected
```

#### Usage Example

**Server:**
```go
hub := realtime.NewHub()
go hub.Run(ctx)

// In worker after adding item:
if score.Final > 0.8 {
    hub.BroadcastNewItem(item.ID, item.Title, score.Final)
}
```

**Client (JavaScript):**
```javascript
const ws = new WebSocket('ws://localhost:8080/v1/ws');

ws.onopen = () => {
    ws.send(JSON.stringify({
        type: 'subscribe',
        topic: 'new_items'
    }));
};

ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.type === 'new_item') {
        showNotification(msg.payload.title);
    }
};
```

## Database Migrations for Phase 2

### Migration: Add Vector Search
```sql
-- Up
CREATE EXTENSION IF NOT EXISTS vector;

ALTER TABLE items ADD COLUMN embedding vector(1536);
ALTER TABLE items ADD COLUMN full_text TEXT;
ALTER TABLE items ADD COLUMN reading_time_minutes INT;

CREATE INDEX items_embedding_idx ON items USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
CREATE INDEX items_full_text_idx ON items USING gin(to_tsvector('english', full_text));

-- Down
DROP INDEX IF EXISTS items_embedding_idx;
DROP INDEX IF EXISTS items_full_text_idx;

ALTER TABLE items DROP COLUMN IF EXISTS embedding;
ALTER TABLE items DROP COLUMN IF EXISTS full_text;
ALTER TABLE items DROP COLUMN IF EXISTS reading_time_minutes;
```

### Migration: Add Real-time Tracking
```sql
-- Up
CREATE TABLE IF NOT EXISTS realtime_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_type VARCHAR(50) NOT NULL,
    item_id UUID REFERENCES items(id),
    payload JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX realtime_events_created_at_idx ON realtime_events(created_at DESC);
CREATE INDEX realtime_events_type_idx ON realtime_events(event_type);

-- Down
DROP TABLE IF EXISTS realtime_events;
```

## Configuration Updates

Add to `internal/config/config.go`:
```go
type Config struct {
    // ... existing fields ...
    
    // Crawler
    CrawlerEnabled      bool
    CrawlerHeadless     bool
    CrawlerMaxConcurrent int
    
    // Vector Search
    VectorSearchEnabled bool
    EmbeddingsModel     string
    
    // Real-time
    WebSocketEnabled    bool
    WebSocketPort       string
}
```

Environment variables:
```bash
# Crawler
CRAWLER_ENABLED=true
CRAWLER_HEADLESS=true
CRAWLER_MAX_CONCURRENT=5

# Vector Search
VECTOR_SEARCH_ENABLED=true
EMBEDDINGS_MODEL=text-embedding-3-small

# Real-time
WEBSOCKET_ENABLED=true
WEBSOCKET_PORT=8081
```

## Testing Phase 2

### Test Crawler
```go
func TestCrawler(t *testing.T) {
    crawler := crawler.NewCrawler(crawler.Config{
        Headless: true,
        Timeout:  30 * time.Second,
    })
    
    result, err := crawler.Crawl(context.Background(), "https://openai.com/blog/gpt-4")
    assert.NoError(t, err)
    assert.NotEmpty(t, result.Content)
    assert.NotEmpty(t, result.Title)
}
```

### Test Vector Search
```go
func TestVectorSearch(t *testing.T) {
    vs := search.NewVectorSearch()
    
    // Index test items
    err := vs.IndexItem(ctx, itemID1, "GPT-4 is a large language model")
    assert.NoError(t, err)
    
    // Search
    results, err := vs.Search(ctx, "language models", 5)
    assert.NoError(t, err)
    assert.Greater(t, len(results), 0)
}
```

### Test WebSocket
```bash
# Use wscat to test
npm install -g wscat
wscat -c ws://localhost:8080/v1/ws

# Send subscription
> {"type":"subscribe","topic":"new_items"}

# Should receive messages when new items added
```

## Performance Considerations

### Crawler
- **Concurrency**: Limit to 5-10 concurrent crawls
- **Memory**: Headless Chrome uses ~100MB per instance
- **Rate Limiting**: 1 request per second per domain
- **Timeout**: 30 seconds max per page

### Vector Search
- **Indexing**: Batch process 100 items at a time
- **Storage**: 1536 floats × 4 bytes = 6KB per item
- **Query Time**: <50ms with IVFFlat index
- **Embeddings Cost**: ~$0.0001 per item (OpenAI)

### WebSocket
- **Connections**: Support 10,000+ concurrent connections
- **Memory**: ~10KB per connection
- **Broadcast**: <10ms to all clients
- **Heartbeat**: Ping every 30 seconds

## Dependencies to Add

```bash
go get github.com/go-rod/rod
go get github.com/go-shiori/go-readability
go get github.com/gorilla/websocket
go get github.com/pgvector/pgvector-go
```

## Next Steps

1. **Implement Crawler**
   - Start with rod integration
   - Add readability extraction
   - Test on OpenAI blog

2. **Implement Vector Search**
   - Run pgvector migration
   - Integrate embeddings API
   - Add search endpoints

3. **Implement WebSocket**
   - Add gorilla/websocket
   - Create hub and client management
   - Update UI to connect

4. **Integration**
   - Use crawler in worker for full-text
   - Index items with vector search
   - Broadcast new items via WebSocket

See `PHASE3_PLAN.md` for admin dashboard, personalization, and mobile apps.
