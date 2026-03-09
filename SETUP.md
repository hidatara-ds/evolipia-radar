# Setup Guide

## 1. Create Neon Database

1. Sign up: https://neon.tech
2. Create project: `evolipia-radar`
3. Copy connection string from dashboard

## 2. Run Migrations

Open Neon SQL Editor and run these in order:

### Migration 1: Core Tables

```sql
CREATE TABLE sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    category TEXT NOT NULL,
    url TEXT NOT NULL,
    mapping_json JSONB NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT NOT NULL DEFAULT 'pending',
    last_test_status TEXT NULL,
    last_test_message TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sources_enabled ON sources(enabled);
CREATE UNIQUE INDEX idx_sources_unique_url ON sources(url);

CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    content_hash TEXT NOT NULL,
    domain TEXT NOT NULL,
    category TEXT NOT NULL,
    raw_excerpt TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_items_dedup ON items(content_hash);
CREATE INDEX idx_items_published_at ON items(published_at DESC);
CREATE INDEX idx_items_domain ON items(domain);
CREATE INDEX idx_items_source_id ON items(source_id);

CREATE TABLE signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    points INT NULL,
    comments INT NULL,
    rank_pos INT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_signals_item_fetched ON signals(item_id, fetched_at DESC);

CREATE TABLE scores (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    hot DOUBLE PRECISION NOT NULL,
    relevance DOUBLE PRECISION NOT NULL,
    credibility DOUBLE PRECISION NOT NULL,
    novelty DOUBLE PRECISION NOT NULL,
    final DOUBLE PRECISION NOT NULL,
    computed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_scores_final ON scores(final DESC);

CREATE TABLE summaries (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    tldr TEXT NOT NULL,
    why_it_matters TEXT NOT NULL,
    tags JSONB NOT NULL,
    method TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_summaries_tags_gin ON summaries USING GIN (tags);

CREATE TABLE fetch_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL,
    error TEXT NULL,
    items_fetched INT NOT NULL DEFAULT 0,
    items_inserted INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_fetch_runs_source_time ON fetch_runs(source_id, fetched_at DESC);
```

### Migration 2: Monitoring

```sql
CREATE TABLE scrape_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    items_processed INTEGER DEFAULT 0,
    items_new INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'running',
    error_message TEXT,
    trigger_source VARCHAR(50) DEFAULT 'github_actions'
);

CREATE INDEX idx_scrape_logs_started_at ON scrape_logs(started_at DESC);
```

### Migration 3: Cleanup Function

```sql
CREATE OR REPLACE FUNCTION cleanup_old_data() RETURNS void AS $$
BEGIN
    DELETE FROM items WHERE published_at < NOW() - INTERVAL '45 days';
    DELETE FROM scrape_logs WHERE started_at < NOW() - INTERVAL '60 days';
END;
$$ LANGUAGE plpgsql;
```

### Seed Data

```sql
INSERT INTO sources (name, type, category, url, enabled, status) VALUES
('Hacker News', 'hackernews', 'tech', 'https://news.ycombinator.com', true, 'active'),
('Hugging Face Papers', 'huggingface', 'ai', 'https://huggingface.co/papers', true, 'active'),
('LMSYS Chatbot Arena', 'lmsys', 'ai', 'https://chat.lmsys.org', true, 'active');
```

## 3. GitHub Actions

1. Go to: https://github.com/YOUR_USERNAME/evolipia-radar/settings/secrets/actions
2. Add secret:
   - Name: `DATABASE_URL`
   - Value: Your Neon connection string
3. Push to trigger workflow

## 4. Verify

```sql
SELECT COUNT(*) FROM sources;  -- Should be 3
SELECT COUNT(*) FROM items;    -- Should have data after first run
SELECT * FROM scrape_logs ORDER BY started_at DESC LIMIT 5;
```

Done!
