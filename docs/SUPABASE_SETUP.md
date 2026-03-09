# Supabase Setup Guide

Panduan lengkap untuk setup Evolipia Radar dengan Supabase.

## Prerequisites

- Account Supabase (gratis di [supabase.com](https://supabase.com))
- GitHub repository dengan branch `deploy-supabase`
- golang-migrate CLI (untuk run migrations)

## Step 1: Buat Supabase Project

1. Login ke [supabase.com](https://supabase.com)
2. Klik "New Project"
3. Isi:
   - Name: `evolipia-radar` (atau nama lain)
   - Database Password: Generate strong password (simpan!)
   - Region: Pilih yang terdekat (Singapore untuk Indonesia)
4. Tunggu project selesai dibuat (~2 menit)

## Step 2: Get Database Connection String

1. Di Supabase Dashboard, buka: Settings → Database
2. Scroll ke "Connection string"
3. Pilih tab "URI"
4. Copy connection string:
   ```
   postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:5432/postgres
   ```
5. Replace `[YOUR-PASSWORD]` dengan password yang dibuat di Step 1

## Step 3: Install golang-migrate

### Windows (PowerShell)
```powershell
# Download binary
Invoke-WebRequest -Uri "https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.windows-amd64.exe" -OutFile "migrate.exe"

# Move to PATH (optional)
Move-Item migrate.exe C:\Windows\System32\migrate.exe
```

### macOS/Linux
```bash
# Using brew (macOS)
brew install golang-migrate

# Using curl (Linux)
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

## Step 4: Run Migrations

```bash
# Set DATABASE_URL
export DATABASE_URL="postgresql://postgres:[password]@db.xxx.supabase.co:5432/postgres"

# Run migrations
migrate -path migrations -database "$DATABASE_URL" up

# Verify
migrate -path migrations -database "$DATABASE_URL" version
```

Expected output:
```
000003 (dirty=false)
```

## Step 5: Setup RLS & Functions

1. Di Supabase Dashboard, buka: SQL Editor
2. Copy isi file `configs/supabase-setup.sql`
3. Paste dan Run di SQL Editor
4. Verify:
   ```sql
   -- Check RLS enabled
   SELECT tablename, rowsecurity FROM pg_tables WHERE tablename = 'items';
   
   -- Check function exists
   SELECT proname FROM pg_proc WHERE proname = 'cleanup_old_data';
   ```

## Step 6: Setup GitHub Secrets

1. Di GitHub repo, buka: Settings → Secrets and variables → Actions
2. Klik "New repository secret"
3. Tambahkan:
   - Name: `SUPABASE_DB_URL`
   - Value: Connection string dari Step 2
4. Save

## Step 7: Test Worker Locally (Optional)

```bash
# Clone repo
git clone <your-repo>
cd evolipia-radar
git checkout deploy-supabase

# Set env
export DATABASE_URL="postgresql://..."

# Run worker
go run ./cmd/worker
```

Expected output:
```
========================================
Worker started at: 2026-03-09T22:30:00Z
========================================
Starting ingestion...
Ingestion completed successfully
Items processed: 150
New items: 25
========================================
Worker finished at: 2026-03-09T22:32:15Z
Duration: 2m15s
Status: SUCCESS
========================================
```

## Step 8: Verify GitHub Actions

1. Di GitHub repo, buka: Actions tab
2. Pilih workflow "Scheduled News Scraper"
3. Klik "Run workflow" → "Run workflow" (manual trigger)
4. Tunggu selesai (~2-3 menit)
5. Check logs untuk verify success

## Step 9: Monitor

### Check Scrape Logs
Di Supabase Dashboard → SQL Editor:
```sql
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
```

### Check Items
```sql
SELECT COUNT(*) as total_items FROM items;
SELECT COUNT(*) as today_items FROM items WHERE published_at >= CURRENT_DATE;
```

### Check Database Size
```sql
SELECT pg_size_pretty(pg_database_size(current_database()));
```

## Troubleshooting

### Migration Failed
```bash
# Check current version
migrate -path migrations -database "$DATABASE_URL" version

# Force version (if dirty)
migrate -path migrations -database "$DATABASE_URL" force 2

# Try again
migrate -path migrations -database "$DATABASE_URL" up
```

### Worker Error: "connection refused"
- Verify DATABASE_URL correct
- Check Supabase project status (not paused)
- Verify IP not blocked (Supabase allows all by default)

### GitHub Actions Failed
- Check secret `SUPABASE_DB_URL` exists
- Verify connection string format correct
- Check Actions logs untuk detail error

## Maintenance

### Manual Cleanup
```bash
# Via GitHub Actions
Actions → Weekly Data Cleanup → Run workflow

# Via SQL Editor
SELECT cleanup_old_data();
```

### Check Schedule
Workflows jalan otomatis:
- 07:00 WIB (00:00 UTC) - Pagi
- 12:00 WIB (05:00 UTC) - Siang
- 19:00 WIB (12:00 UTC) - Malam
- Minggu 09:00 WIB (02:00 UTC) - Weekly cleanup

## Next Steps

1. Setup Flutter app (terpisah) untuk consume data
2. Configure RLS policies sesuai kebutuhan
3. Monitor database size (free tier: 500MB)
4. Setup alerts di Supabase (optional)
