# Database Schema

Complete database schema documentation for Evolipia Radar.

## Overview

- **Database:** PostgreSQL 15+
- **Hosting:** Neon.tech (serverless PostgreSQL)
- **Connection Pooling:** Enabled (max 3 connections)
- **Migrations:** Located in `migrations/` directory

## Entity Relationship Diagram

```
┌─────────────┐
│   sources   │
└──────┬──────┘
       │ 1
       │
       │ N
┌──────▼──────┐      ┌─────────────┐
│    items    │◄─────┤  summaries  │
└──────┬──────┘ 1:1  └─────────────┘
       │ 1
       │
       ├─────────────┐
       │ N           │ N
┌──────▼──────┐ ┌───▼────────┐
│   signals   │ │   scores   │
└─────────────┘ └────────────┘

┌─────────────┐
│ fetch_runs  │
└──────┬──────┘
       │ N
       │
       │ 1
┌──────▼──────┐
│   sources   │
└─────────────┘
```

## Tables

### sources

News sources configuration.

```sql
CREATE TABLE sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- 'rss', 'api', 'scrape'
    category VARCHAR(100),
    url TEXT NOT NULL UNIQUE,
    mapping_json JSONB,
    enabled BOOLEAN DEFAULT true,
    status VARCHAR(50) DEFAULT 'active',
    last_test_status VARCHAR(50),
    last_test_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sources_enabled ON sources(enabled);
CREATE INDEX idx_sources_type ON sources(type);
```

**Columns:**
- `id` - Unique identifier
- `name` - Human-readable source name
- `type` - Source type (rss, api, scrape)
- `category` - Content category
- `url` - Source URL
- `mapping_json` - Field mapping configuration
- `enabled` - Whether source is active
- `status` - Current status
- `last_test_status` - Last health check result
- `last_test_message` - Health check message
- `created_at` - Creation timestamp
- `updated_at` - Last update timestamp

**Example:**
```sql
INSERT INTO sources (name, type, category, url, enabled)
VALUES ('Hacker News', 'api', 'tech', 'https://news.ycombinator.com/best', true);
```

---

### items

Raw news articles.

```sql
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    content_hash VARCHAR(64) NOT NULL UNIQUE,
    domain VARCHAR(255),
    category VARCHAR(100),
    raw_excerpt TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_items_source ON items(source_id);
CREATE INDEX idx_items_published ON items(published_at DESC);
CREATE INDEX idx_items_hash ON items(content_hash);
CREATE INDEX idx_items_domain ON items(domain);
```

**Columns:**
- `id` - Unique identifier
- `source_id` - Reference to source
- `title` - Article title
- `url` - Article URL
- `published_at` - Publication timestamp
- `content_hash` - SHA-256 hash for deduplication
- `domain` - Source domain
- `category` - Content category
- `raw_excerpt` - Raw content excerpt
- `created_at` - Creation timestamp

**Deduplication:**
Content hash is generated from `title + url + published_at` to prevent duplicates.

---

### summaries

AI-generated summaries and tags.

```sql
CREATE TABLE summaries (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    tldr TEXT,
    why_it_matters TEXT,
    tags JSONB DEFAULT '[]'::jsonb,
    method VARCHAR(50) DEFAULT 'ai',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_summaries_tags ON summaries USING GIN(tags);
```

**Columns:**
- `item_id` - Reference to item (1:1)
- `tldr` - Brief summary
- `why_it_matters` - Significance explanation
- `tags` - Array of topic tags (JSONB)
- `method` - Generation method (ai, manual)
- `created_at` - Creation timestamp

**Tags Format:**
```json
["llm", "tools", "research"]
```

**Querying by Tag:**
```sql
-- Find items with 'llm' tag
SELECT i.* FROM items i
JOIN summaries s ON s.item_id = i.id
WHERE s.tags @> '["llm"]'::jsonb;
```

---

### scores

Multi-factor relevance scores.

```sql
CREATE TABLE scores (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    hot FLOAT DEFAULT 0,
    relevance FLOAT DEFAULT 0,
    credibility FLOAT DEFAULT 0,
    novelty FLOAT DEFAULT 0,
    final FLOAT DEFAULT 0,
    computed_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_scores_final ON scores(final DESC);
CREATE INDEX idx_scores_hot ON scores(hot DESC);
```

**Columns:**
- `item_id` - Reference to item (1:1)
- `hot` - Recency and engagement score (0-1)
- `relevance` - Topic relevance score (0-1)
- `credibility` - Source credibility score (0-1)
- `novelty` - Content uniqueness score (0-1)
- `final` - Weighted final score (0-1)
- `computed_at` - Computation timestamp

**Score Calculation:**
```
final = (hot * 0.4) + (relevance * 0.3) + (credibility * 0.2) + (novelty * 0.1)
```

---

### signals

Engagement metrics (points, comments, rank).

```sql
CREATE TABLE signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    points INT DEFAULT 0,
    comments INT DEFAULT 0,
    rank_pos INT,
    fetched_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_signals_item ON signals(item_id);
CREATE INDEX idx_signals_fetched ON signals(fetched_at DESC);
```

**Columns:**
- `id` - Unique identifier
- `item_id` - Reference to item
- `points` - Upvotes/points
- `comments` - Comment count
- `rank_pos` - Position in ranking
- `fetched_at` - Fetch timestamp

**Usage:**
Track engagement over time for trending detection.

---

### fetch_runs

Scraper job history.

```sql
CREATE TABLE fetch_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    error TEXT,
    items_fetched INT DEFAULT 0,
    items_inserted INT DEFAULT 0,
    fetched_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_fetch_runs_source ON fetch_runs(source_id);
CREATE INDEX idx_fetch_runs_fetched ON fetch_runs(fetched_at DESC);
```

**Columns:**
- `id` - Unique identifier
- `source_id` - Reference to source
- `status` - Run status (success, error)
- `error` - Error message if failed
- `items_fetched` - Total items fetched
- `items_inserted` - New items inserted
- `fetched_at` - Run timestamp

---

## Migrations

### Running Migrations

```bash
# Apply all migrations
psql $DATABASE_URL < migrations/001_initial_schema.sql

# Verify
psql $DATABASE_URL -c "\dt"
```

### Creating New Migrations

1. Create file: `migrations/XXX_description.sql`
2. Add migration SQL
3. Test locally
4. Apply to production

**Example:**
```sql
-- migrations/002_add_user_bookmarks.sql
CREATE TABLE bookmarks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, item_id)
);

CREATE INDEX idx_bookmarks_user ON bookmarks(user_id);
```

---

## Queries

### Common Queries

**Get Latest News:**
```sql
SELECT 
    i.id,
    i.title,
    i.url,
    i.domain,
    i.published_at,
    i.category,
    COALESCE(sc.final, 0.5) as score,
    COALESCE(su.tldr, '') as tldr,
    COALESCE(su.why_it_matters, '') as why_it_matters,
    COALESCE(su.tags, '[]'::jsonb) as tags
FROM items i
LEFT JOIN scores sc ON i.id = sc.item_id
LEFT JOIN summaries su ON i.id = su.item_id
WHERE i.published_at >= NOW() - INTERVAL '7 days'
ORDER BY COALESCE(sc.final, 0) DESC, i.published_at DESC
LIMIT 20;
```

**Get News by Topic:**
```sql
SELECT i.*, su.tags
FROM items i
JOIN summaries su ON i.id = su.item_id
WHERE su.tags @> '["llm"]'::jsonb
  AND i.published_at >= NOW() - INTERVAL '7 days'
ORDER BY i.published_at DESC
LIMIT 20;
```

**Get Trending Articles:**
```sql
SELECT i.*, sc.hot
FROM items i
JOIN scores sc ON i.id = sc.item_id
WHERE i.published_at >= NOW() - INTERVAL '24 hours'
ORDER BY sc.hot DESC
LIMIT 10;
```

**Search Articles:**
```sql
SELECT i.*, su.tldr
FROM items i
LEFT JOIN summaries su ON i.id = su.item_id
WHERE i.title ILIKE '%transformer%'
   OR su.tldr ILIKE '%transformer%'
ORDER BY i.published_at DESC
LIMIT 20;
```

---

## Performance Optimization

### Indexes

Current indexes provide optimal performance for:
- Time-based queries (published_at)
- Tag filtering (GIN index on tags)
- Score-based ranking (final score)
- Deduplication (content_hash)

### Query Optimization

**Use EXPLAIN ANALYZE:**
```sql
EXPLAIN ANALYZE
SELECT * FROM items
WHERE published_at >= NOW() - INTERVAL '7 days'
ORDER BY published_at DESC
LIMIT 20;
```

**Optimize with Covering Indexes:**
```sql
-- If frequently querying title + published_at
CREATE INDEX idx_items_published_title 
ON items(published_at DESC, title);
```

### Connection Pooling

Configured in `api/news/index.go`:
```go
db.SetMaxOpenConns(3)
db.SetMaxIdleConns(1)
db.SetConnMaxLifetime(5 * time.Minute)
db.SetConnMaxIdleTime(30 * time.Second)
```

---

## Backup & Recovery

### Automated Backups

Neon.tech provides:
- Daily automated backups
- 7-day retention
- Point-in-time recovery

### Manual Backup

```bash
# Full backup
pg_dump $DATABASE_URL > backup_$(date +%Y%m%d).sql

# Schema only
pg_dump --schema-only $DATABASE_URL > schema.sql

# Data only
pg_dump --data-only $DATABASE_URL > data.sql
```

### Restore

```bash
# Full restore
psql $DATABASE_URL < backup_20260318.sql

# Restore specific table
pg_restore -t items backup.sql
```

---

## Monitoring

### Key Metrics

- **Connection count:** Should stay < 3
- **Query duration:** p95 < 100ms
- **Table sizes:** Monitor growth
- **Index usage:** Ensure indexes are used

### Monitoring Queries

**Connection Count:**
```sql
SELECT count(*) FROM pg_stat_activity;
```

**Slow Queries:**
```sql
SELECT query, mean_exec_time, calls
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

**Table Sizes:**
```sql
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

**Index Usage:**
```sql
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

---

## Security

### Best Practices

- ✅ Use SSL/TLS connections
- ✅ Store credentials in environment variables
- ✅ Use prepared statements (SQL injection prevention)
- ✅ Implement row-level security (RLS) if needed
- ✅ Regular security updates
- ✅ Monitor for suspicious queries

### Row-Level Security (Future)

```sql
-- Enable RLS
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

-- Create policy
CREATE POLICY items_select_policy ON items
FOR SELECT
USING (true);  -- Public read

CREATE POLICY items_insert_policy ON items
FOR INSERT
WITH CHECK (auth.uid() = source_id);  -- Only source owner can insert
```

---

## Troubleshooting

### Connection Issues

**Error:** `connection refused`
- Check DATABASE_URL format
- Verify Neon.tech database is active
- Check network/firewall

**Error:** `too many connections`
- Reduce max_connections in code
- Check for connection leaks
- Monitor active connections

### Performance Issues

**Slow Queries:**
1. Run EXPLAIN ANALYZE
2. Check if indexes are used
3. Add missing indexes
4. Optimize query structure

**High CPU:**
1. Identify expensive queries
2. Add indexes
3. Optimize application logic
4. Consider caching

---

## References

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Neon.tech Docs](https://neon.tech/docs)
- [PostgreSQL Performance Tips](https://wiki.postgresql.org/wiki/Performance_Optimization)
