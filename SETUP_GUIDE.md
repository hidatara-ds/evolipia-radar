# Setup Guide - Evolipia Radar (Supabase Deploy)

Complete setup guide untuk deploy Evolipia Radar dengan Supabase dan GitHub Actions.

## 📋 Prerequisites

- Account Supabase (gratis di [supabase.com](https://supabase.com))
- GitHub repository dengan branch `deploy-supabase`
- **Tidak perlu local testing** (langsung deploy ke GitHub Actions)

---

## 🚀 Setup Steps

### Step 1: Buat Supabase Project (2 menit)

1. Login ke [supabase.com](https://supabase.com)
2. Klik **New Project**
3. Isi form:
   - **Name**: `evolipia-radar`
   - **Database Password**: Generate strong password (SIMPAN!)
   - **Region**: Singapore (untuk Indonesia)
4. Tunggu project ready (~2 menit)

### Step 2: Copy Connection String (30 detik)

1. Di Supabase Dashboard: **Settings** → **Database**
2. Scroll ke **Connection string**
3. Pilih tab **URI**
4. Copy connection string:
   ```
   postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:5432/postgres
   ```
5. Replace `[YOUR-PASSWORD]` dengan password dari Step 1

**⚠️ SIMPAN connection string ini!** Akan digunakan di Step 4.

### Step 3: Run Database Migrations (5 menit)

**Karena golang-migrate CLI ribet dengan special characters, kita pakai cara manual via SQL Editor.**

#### 3.1 Buka SQL Editor

1. Di Supabase Dashboard: **SQL Editor** (sidebar kiri)
2. Klik **New query**

#### 3.2 Run Migration 1: Init Schema

Copy-paste SQL berikut, lalu klik **Run**:

```sql
-- Migration 1: Init Schema
-- Creates: sources, items, signals, scores, summaries, fetch_runs tables

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

**Expected**: "Success. No rows returned"

#### 3.3 Run Migration 2: Scrape Logs

Copy-paste SQL berikut, lalu klik **Run**:

```sql
-- Migration 2: Scrape Logs Table
-- For monitoring worker execution

CREATE TABLE IF NOT EXISTS scrape_logs (
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

**Expected**: "Success. No rows returned"

#### 3.4 Run Migration 3: Data Retention Function

Copy-paste SQL berikut, lalu klik **Run**:

```sql
-- Migration 3: Data Retention Function
-- Auto cleanup old data (45 days)

CREATE OR REPLACE FUNCTION cleanup_old_data() RETURNS void AS $$
BEGIN
    DELETE FROM items 
    WHERE published_at < NOW() - INTERVAL '45 days';
    
    DELETE FROM scrape_logs 
    WHERE started_at < NOW() - INTERVAL '60 days';
END;
$$ LANGUAGE plpgsql;
```

**Expected**: "Success. No rows returned"

#### 3.5 Setup RLS & Helper Functions

Copy-paste SQL berikut, lalu klik **Run**:

```sql
-- Setup Row Level Security & Helper Functions

-- Enable RLS for items
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

-- Policy: Allow read for all (for Flutter app)
CREATE POLICY "Allow anonymous read" ON items
    FOR SELECT USING (true);

-- View for daily feed
CREATE OR REPLACE VIEW daily_feed AS
SELECT 
    i.id,
    i.title,
    i.url,
    i.published_at,
    i.domain,
    i.category,
    s.final as score,
    sm.tldr as summary
FROM items i
LEFT JOIN scores s ON i.id = s.item_id
LEFT JOIN summaries sm ON i.id = sm.item_id
WHERE i.published_at >= CURRENT_DATE - INTERVAL '1 day'
ORDER BY s.final DESC NULLS LAST;

-- Function to get daily feed
CREATE OR REPLACE FUNCTION get_daily_feed(limit_count INTEGER DEFAULT 20)
RETURNS TABLE (
    id UUID,
    title TEXT,
    url TEXT,
    published_at TIMESTAMPTZ,
    domain TEXT,
    category TEXT,
    score DOUBLE PRECISION,
    summary TEXT
) AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM daily_feed
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;
```

**Expected**: "Success. No rows returned"

#### 3.6 Seed Default Sources (PENTING!)

Copy-paste SQL berikut, lalu klik **Run**:

```sql
-- Seed default sources
-- Worker needs these sources to scrape data

INSERT INTO sources (name, type, category, url, enabled, status) VALUES
('Hacker News', 'hackernews', 'tech', 'https://news.ycombinator.com', true, 'active'),
('Hugging Face Papers', 'huggingface', 'ai', 'https://huggingface.co/papers', true, 'active'),
('LMSYS Chatbot Arena', 'lmsys', 'ai', 'https://chat.lmsys.org', true, 'active');

-- Verify
SELECT id, name, type, enabled FROM sources;
```

**Expected**: 3 rows returned dengan sources yang di-insert

#### 3.7 Verify Setup

Copy-paste SQL berikut untuk verify semua OK:

```sql
-- Verify all tables exist
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- Expected: fetch_runs, items, scores, scrape_logs, signals, sources, summaries

-- Verify sources
SELECT COUNT(*) as source_count FROM sources WHERE enabled = true;
-- Expected: 3

-- Verify functions
SELECT proname FROM pg_proc 
WHERE pronamespace = 'public'::regnamespace
AND proname IN ('cleanup_old_data', 'get_daily_feed');
-- Expected: 2 functions
```

✅ **Migrations Complete!**

### Step 4: Setup GitHub Secret (1 menit)

1. Buka GitHub repository: `https://github.com/hidatara-ds/evolipia-radar`
2. Klik **Settings** → **Secrets and variables** → **Actions**
3. Klik **New repository secret**
4. Isi:
   - **Name**: `SUPABASE_DB_URL`
   - **Secret**: Connection string dari Step 2
     ```
     postgresql://postgres:eRmnQG8QCxiblkWe@db.xxx.supabase.co:5432/postgres
     ```
5. Klik **Add secret**

### Step 5: Deploy & Test (2 menit)

#### 5.1 Push Branch (jika belum)

```bash
git push origin deploy-supabase
```

#### 5.2 Manual Trigger Workflow

1. Buka GitHub repository
2. Klik tab **Actions**
3. Pilih workflow **"Scheduled News Scraper"**
4. Klik **Run workflow** (dropdown)
5. Klik **Run workflow** (button hijau)
6. Tunggu workflow selesai (~2-3 menit)

#### 5.3 Check Logs

1. Klik pada workflow run yang baru
2. Klik job **"scrape"**
3. Expand **"Run Worker"**
4. Verify output:
   ```
   ========================================
   Worker started at: 2026-03-09T...
   ========================================
   Starting ingestion...
   Ingestion completed successfully
   Items processed: X
   New items: Y
   ========================================
   Worker finished at: 2026-03-09T...
   Status: SUCCESS
   ========================================
   ```

#### 5.4 Verify Data di Supabase

Di Supabase SQL Editor:

```sql
-- Check scrape logs
SELECT * FROM scrape_logs ORDER BY started_at DESC LIMIT 5;

-- Check items
SELECT COUNT(*) FROM items;
SELECT title, url, domain FROM items ORDER BY created_at DESC LIMIT 10;

-- Check scores
SELECT COUNT(*) FROM scores;
```

✅ **Setup Complete!**

---

## 🎉 Done!

Worker sekarang jalan otomatis 3x/hari:
- **07:00 WIB** (00:00 UTC)
- **12:00 WIB** (05:00 UTC)
- **19:00 WIB** (12:00 UTC)

---

## 📱 Next: Build Flutter App

Flutter app akan query langsung ke Supabase (tidak perlu backend API).

### Get Supabase Credentials

1. Di Supabase Dashboard: **Settings** → **API**
2. Copy:
   - **Project URL**: `https://xxx.supabase.co`
   - **anon public key**: `eyJhbGc...`

### Flutter Setup

```dart
// pubspec.yaml
dependencies:
  supabase_flutter: ^2.0.0

// main.dart
await Supabase.initialize(
  url: 'https://axxbcfnrlfnjyfanaogb.supabase.co',
  anonKey: 'your-anon-key',
);

// Query items
final items = await Supabase.instance.client
  .from('items')
  .select('*, scores(*), summaries(*)')
  .order('created_at', ascending: false)
  .limit(20);
```

---

## 🔧 Maintenance

### Manual Cleanup

```sql
SELECT cleanup_old_data();
```

### Check Database Size

```sql
SELECT pg_size_pretty(pg_database_size(current_database()));
```

### Disable/Enable Sources

```sql
-- Disable source
UPDATE sources SET enabled = false WHERE name = 'Hacker News';

-- Enable source
UPDATE sources SET enabled = true WHERE name = 'Hacker News';
```

---

## 🆘 Troubleshooting

### Worker Failed

1. Check GitHub Actions logs
2. Verify `SUPABASE_DB_URL` secret correct
3. Check Supabase project not paused

### No Items Scraped

```sql
-- Check sources enabled
SELECT name, enabled, status FROM sources;

-- Check fetch_runs for errors
SELECT * FROM fetch_runs ORDER BY fetched_at DESC LIMIT 10;
```

### Database Full (500MB limit)

```sql
-- Run cleanup
SELECT cleanup_old_data();

-- Check size
SELECT pg_size_pretty(pg_database_size(current_database()));
```

---

## 📚 Additional Resources

- [Supabase Documentation](https://supabase.com/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Flutter Supabase SDK](https://pub.dev/packages/supabase_flutter)
