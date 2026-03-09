# Manual Migration Guide - Supabase

Panduan untuk run migrations secara manual via Supabase SQL Editor (tanpa golang-migrate CLI).

## Kenapa Manual?

- Password dengan karakter khusus (`#`, `@`, dll) sulit di-escape di bash
- Lebih simple, tidak perlu install golang-migrate
- Langsung via Supabase Dashboard (copy-paste SQL)

## Step-by-Step Manual Migration

### 1. Login ke Supabase Dashboard

1. Buka [supabase.com](https://supabase.com)
2. Login dan pilih project `evolipia-radar` (atau nama project Anda)
3. Di sidebar kiri, klik **SQL Editor**

### 2. Run Migration 1: Init Schema

Copy SQL berikut dan paste di SQL Editor, lalu klik **Run**:

```sql
-- Migration 000001: Init Schema
-- File: migrations/000001_init_schema.up.sql

-- Create sources table
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

-- Create items table
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

-- Create signals table
CREATE TABLE signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    points INT NULL,
    comments INT NULL,
    rank_pos INT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_signals_item_fetched ON signals(item_id, fetched_at DESC);

-- Create scores table
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

-- Create summaries table
CREATE TABLE summaries (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    tldr TEXT NOT NULL,
    why_it_matters TEXT NOT NULL,
    tags JSONB NOT NULL,
    method TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_summaries_tags_gin ON summaries USING GIN (tags);

-- Create fetch_runs table
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

**Expected Result**: "Success. No rows returned"

### 3. Run Migration 2: Scrape Logs Table

Copy SQL berikut dan paste di SQL Editor, lalu klik **Run**:

```sql
-- Migration 000002: Create Scrape Logs Table
-- File: migrations/000002_create_scrape_logs_table.up.sql

CREATE TABLE IF NOT EXISTS scrape_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    items_processed INTEGER DEFAULT 0,
    items_new INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'running', -- running, success, failed
    error_message TEXT,
    trigger_source VARCHAR(50) DEFAULT 'github_actions' -- github_actions, manual, etc
);

CREATE INDEX idx_scrape_logs_started_at ON scrape_logs(started_at DESC);
```

**Expected Result**: "Success. No rows returned"

### 4. Run Migration 3: Data Retention Function

Copy SQL berikut dan paste di SQL Editor, lalu klik **Run**:

```sql
-- Migration 000003: Create Data Retention Function
-- File: migrations/000003_create_data_retention_function.up.sql

-- Function untuk cleanup data lama (45 hari)
CREATE OR REPLACE FUNCTION cleanup_old_data() RETURNS void AS $$
BEGIN
    -- Delete news items older than 45 days
    DELETE FROM items 
    WHERE published_at < NOW() - INTERVAL '45 days';
    
    -- Delete scrape logs older than 60 days (keep longer for audit)
    DELETE FROM scrape_logs 
    WHERE started_at < NOW() - INTERVAL '60 days';
    
    -- Vacuum akan di-handle otomatis oleh Supabase
END;
$$ LANGUAGE plpgsql;
```

**Expected Result**: "Success. No rows returned"

### 5. Setup RLS & Helper Functions

Copy SQL berikut dan paste di SQL Editor, lalu klik **Run**:

```sql
-- Setup RLS (Row Level Security) untuk keamanan
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

-- Policy: Allow read untuk semua (anonymous) - untuk Flutter app
CREATE POLICY "Allow anonymous read" ON items
    FOR SELECT USING (true);

-- View untuk Flutter app (simplified)
CREATE OR REPLACE VIEW daily_feed AS
SELECT 
    id,
    title,
    url,
    summary,
    score,
    published_at,
    source,
    topic
FROM items
WHERE published_at >= CURRENT_DATE - INTERVAL '1 day'
ORDER BY score DESC;

-- Function untuk get feed (bisa dipanggil dari Flutter via RPC)
CREATE OR REPLACE FUNCTION get_daily_feed(limit_count INTEGER DEFAULT 20)
RETURNS SETOF items AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM items
    WHERE published_at >= CURRENT_DATE - INTERVAL '1 day'
    ORDER BY score DESC
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;
```

**Expected Result**: "Success. No rows returned"

### 6. Verify Migrations

Copy SQL berikut untuk verify semua table dan function sudah dibuat:

```sql
-- Check tables
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
  AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- Expected: fetch_runs, items, scores, scrape_logs, signals, sources, summaries

-- Check indexes
SELECT indexname 
FROM pg_indexes 
WHERE schemaname = 'public'
ORDER BY indexname;

-- Check functions
SELECT proname 
FROM pg_proc 
WHERE pronamespace = 'public'::regnamespace
  AND proname IN ('cleanup_old_data', 'get_daily_feed');

-- Check RLS enabled
SELECT tablename, rowsecurity 
FROM pg_tables 
WHERE tablename = 'items';
-- Expected: rowsecurity = true

-- Check policies
SELECT policyname, tablename 
FROM pg_policies 
WHERE tablename = 'items';
-- Expected: "Allow anonymous read"
```

### 7. Seed Default Sources (Important!)

Worker membutuhkan sources untuk scraping. Insert default sources:

```sql
-- Insert default sources dari configs/default_sources.yaml
INSERT INTO sources (name, type, category, url, enabled, status) VALUES
('Hacker News', 'hackernews', 'tech', 'https://news.ycombinator.com', true, 'active'),
('Hugging Face Papers', 'huggingface', 'ai', 'https://huggingface.co/papers', true, 'active'),
('LMSYS Chatbot Arena', 'lmsys', 'ai', 'https://chat.lmsys.org', true, 'active');

-- Verify
SELECT id, name, type, category, enabled FROM sources;
```

### 8. Test Data (Optional)

Insert test data untuk verify everything works:

```sql
-- Get a source_id first
SELECT id FROM sources LIMIT 1;
-- Copy the UUID, replace [SOURCE_ID] below

-- Insert test item
INSERT INTO items (source_id, title, url, published_at, content_hash, domain, category)
VALUES (
    '[SOURCE_ID]'::uuid,  -- Replace with actual source_id
    'Test Article',
    'https://example.com/test-' || gen_random_uuid(),  -- Unique URL
    NOW(),
    md5(random()::text),  -- Random hash
    'example.com',
    'technology'
);

-- Insert test scrape log
INSERT INTO scrape_logs (started_at, completed_at, items_processed, items_new, status)
VALUES (
    NOW() - INTERVAL '5 minutes',
    NOW(),
    1,
    1,
    'success'
);

-- Verify
SELECT * FROM items LIMIT 1;
SELECT * FROM scrape_logs LIMIT 1;

-- Test cleanup function (won't delete anything yet, data too new)
SELECT cleanup_old_data();

-- Test daily feed view
SELECT * FROM daily_feed LIMIT 5;

-- Test get_daily_feed function
SELECT * FROM get_daily_feed(5);
```

## ✅ Verification Checklist

Setelah run semua SQL di atas, verify:

- [ ] Table `sources` exists
- [ ] Table `items` exists
- [ ] Table `signals` exists
- [ ] Table `scores` exists
- [ ] Table `summaries` exists
- [ ] Table `fetch_runs` exists
- [ ] Table `scrape_logs` exists
- [ ] Function `cleanup_old_data()` exists
- [ ] Function `get_daily_feed()` exists
- [ ] View `daily_feed` exists
- [ ] RLS enabled on `items` table
- [ ] Policy "Allow anonymous read" exists
- [ ] Default sources inserted (3 sources)
- [ ] Test data inserted successfully (optional)

## 🎯 Next Step: Setup GitHub Actions

Sekarang database sudah ready! Next step:

### 1. Add GitHub Secret

1. Di GitHub repo, buka: **Settings** → **Secrets and variables** → **Actions**
2. Klik **New repository secret**
3. Name: `SUPABASE_DB_URL`
4. Value: `postgresql://postgres:XY76VP@zbEFcW#v@db.axxbcfnrlfnjyfanaogb.supabase.co:5432/postgres`
   (GitHub Actions akan handle special characters dengan benar)
5. Klik **Add secret**

### 2. Test Worker Locally (Optional)

Jika mau test worker di local, buat file `.env`:

```bash
# File: .env
DATABASE_URL=postgresql://postgres:[your-password]@db.axxbcfnrlfnjyfanaogb.supabase.co:5432/postgres
MAX_FETCH_BYTES=2000000
FETCH_TIMEOUT_SECONDS=8
```

Lalu run:

```bash
# Load .env
source .env  # atau: export $(cat .env | xargs)

# Run worker
go run ./cmd/worker
```

### 3. Push Branch & Test GitHub Actions

```bash
# Push branch ke GitHub
git push -u origin deploy-supabase

# Lalu di GitHub:
# 1. Buka Actions tab
# 2. Pilih "Scheduled News Scraper"
# 3. Klik "Run workflow" → "Run workflow"
# 4. Tunggu selesai (~2-3 menit)
# 5. Check logs
```

### 4. Monitor

Di Supabase SQL Editor, check hasil scraping:

```sql
-- Check scrape logs
SELECT 
    started_at,
    completed_at,
    items_processed,
    items_new,
    status,
    trigger_source
FROM scrape_logs
ORDER BY started_at DESC
LIMIT 10;

-- Check items
SELECT COUNT(*) as total_items FROM items;
SELECT * FROM items ORDER BY created_at DESC LIMIT 10;
```

## 🔧 Troubleshooting

### Error: "relation already exists"
Artinya table sudah dibuat sebelumnya. Skip migration tersebut atau drop table dulu:
```sql
DROP TABLE IF EXISTS scrape_logs CASCADE;
-- Lalu run migration lagi
```

### Error: "function already exists"
```sql
DROP FUNCTION IF EXISTS cleanup_old_data();
-- Lalu run migration lagi
```

### Test Connection dari Local
Jika mau test connection (tanpa special char issue), buat file `test-connection.go`:

```go
package main

import (
    "context"
    "fmt"
    "os"
    "github.com/jackc/pgx/v5"
)

func main() {
    // Hardcode connection string (untuk testing)
    connStr := "postgresql://postgres:XY76VP@zbEFcW#v@db.axxbcfnrlfnjyfanaogb.supabase.co:5432/postgres"
    
    conn, err := pgx.Connect(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to connect: %v\n", err)
        os.Exit(1)
    }
    defer conn.Close(context.Background())
    
    var result int
    err = conn.QueryRow(context.Background(), "SELECT 1").Scan(&result)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Println("✅ Connection successful!")
    fmt.Printf("Result: %d\n", result)
}
```

Run:
```bash
go run test-connection.go
```

## 📚 Summary

Manual migration via Supabase SQL Editor:
1. ✅ Lebih mudah (copy-paste SQL)
2. ✅ Tidak perlu install golang-migrate
3. ✅ Tidak ada masalah dengan special characters
4. ✅ Visual feedback langsung di dashboard

Setelah migration selesai, GitHub Actions akan handle scraping otomatis 3x/hari!
