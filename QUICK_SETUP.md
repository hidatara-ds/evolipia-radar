# Quick Setup Guide - Deploy Supabase

Setup Evolipia Radar dengan Supabase dalam 5 langkah mudah.

## 📋 Prerequisites
- Account Supabase (gratis di [supabase.com](https://supabase.com))
- GitHub repository dengan branch `deploy-supabase`

---

## 🚀 5-Minute Setup

### 1️⃣ Buat Supabase Project (2 menit)

1. Login ke [supabase.com](https://supabase.com)
2. Klik **New Project**
3. Isi:
   - Name: `evolipia-radar`
   - Database Password: Generate & simpan!
   - Region: Singapore
4. Tunggu project ready (~2 menit)

### 2️⃣ Copy Connection String (30 detik)

1. Di Supabase Dashboard: **Settings** → **Database**
2. Scroll ke **Connection string** → Tab **URI**
3. Copy connection string:
   ```
   postgresql://postgres:[YOUR-PASSWORD]@db.xxx.supabase.co:5432/postgres
   ```
4. Replace `[YOUR-PASSWORD]` dengan password dari step 1

### 3️⃣ Run Migrations Manual (3 menit)

**Karena password punya special characters, pakai cara manual:**

1. Di Supabase Dashboard: **SQL Editor**
2. Buka file `docs/MANUAL_MIGRATION.md`
3. Copy-paste SQL satu per satu:
   - Migration 1: Init Schema (sources, items, signals, scores, summaries, fetch_runs)
   - Migration 2: Scrape Logs Table
   - Migration 3: Data Retention Function
   - Setup RLS & Helper Functions
   - Seed Default Sources (PENTING!)
4. Klik **Run** untuk setiap SQL

**Verification:**
```sql
-- Check semua table ada
SELECT table_name FROM information_schema.tables 
WHERE table_schema = 'public' AND table_type = 'BASE TABLE'
ORDER BY table_name;

-- Expected: fetch_runs, items, scores, scrape_logs, signals, sources, summaries

-- Check sources ada data
SELECT name, type, enabled FROM sources;
-- Expected: 3 sources (Hacker News, Hugging Face, LMSYS)
```

### 4️⃣ Setup GitHub Secret (1 menit)

1. Di GitHub repo: **Settings** → **Secrets and variables** → **Actions**
2. Klik **New repository secret**
3. Name: `SUPABASE_DB_URL`
4. Value: Connection string dari step 2
   ```
   postgresql://postgres:XY76VP@zbEFcW#v@db.axxbcfnrlfnjyfanaogb.supabase.co:5432/postgres
   ```
5. Klik **Add secret**

### 5️⃣ Push & Test (2 menit)

```bash
# Push branch ke GitHub
git push -u origin deploy-supabase

# Test manual trigger:
# 1. Buka GitHub → Actions tab
# 2. Pilih "Scheduled News Scraper"
# 3. Klik "Run workflow" → "Run workflow"
# 4. Tunggu selesai (~2-3 menit)
```

**Verify di Supabase SQL Editor:**
```sql
-- Check scrape logs
SELECT * FROM scrape_logs ORDER BY started_at DESC LIMIT 5;

-- Check items
SELECT COUNT(*) FROM items;
SELECT title, url, domain FROM items ORDER BY created_at DESC LIMIT 10;
```

---

## ✅ Done!

Worker sekarang jalan otomatis 3x/hari:
- **07:00 WIB** (00:00 UTC) - Pagi
- **12:00 WIB** (05:00 UTC) - Siang
- **19:00 WIB** (12:00 UTC) - Malam

---

## 📊 Monitoring

### Check Scrape Logs
```sql
SELECT 
    started_at,
    completed_at,
    items_processed,
    items_new,
    status
FROM scrape_logs
ORDER BY started_at DESC
LIMIT 10;
```

### Check Database Size
```sql
SELECT pg_size_pretty(pg_database_size(current_database()));
```

### Check Recent Items
```sql
SELECT COUNT(*) FROM items WHERE created_at >= CURRENT_DATE;
```

---

## 🔧 Troubleshooting

### Worker Failed di GitHub Actions
1. Check logs di Actions tab
2. Verify `SUPABASE_DB_URL` secret correct
3. Test connection:
   ```sql
   -- Di Supabase SQL Editor
   SELECT 1;  -- Should return 1
   ```

### No Items Scraped
1. Check sources enabled:
   ```sql
   SELECT name, enabled, status FROM sources;
   ```
2. Update sources jika perlu:
   ```sql
   UPDATE sources SET enabled = true WHERE enabled = false;
   ```

### Database Size Warning
```sql
-- Run cleanup manual
SELECT cleanup_old_data();

-- Check size after
SELECT pg_size_pretty(pg_database_size(current_database()));
```

---

## 📚 Full Documentation

- **Manual Migration**: `docs/MANUAL_MIGRATION.md`
- **Complete Setup**: `docs/SUPABASE_SETUP.md`
- **Architecture**: `docs/README.md`
- **Changelog**: `CHANGELOG.md`

---

## 🎯 Next Steps

### Optional: Test Worker Locally

Buat file `.env`:
```bash
DATABASE_URL=postgresql://postgres:XY76VP@zbEFcW#v@db.axxbcfnrlfnjyfanaogb.supabase.co:5432/postgres
MAX_FETCH_BYTES=2000000
FETCH_TIMEOUT_SECONDS=8
```

Run:
```bash
go run ./cmd/worker
```

### Build Flutter App (Separate Repo)

1. Create Flutter project
2. Add `supabase_flutter` dependency
3. Configure Supabase:
   ```dart
   await Supabase.initialize(
     url: 'https://xxx.supabase.co',
     anonKey: 'your-anon-key',
   );
   ```
4. Query items:
   ```dart
   final items = await Supabase.instance.client
     .from('items')
     .select()
     .order('created_at', ascending: false)
     .limit(20);
   ```

---

## 💡 Tips

- **Free Tier Limits**: 500MB database, 2GB bandwidth/month
- **Data Retention**: Auto-cleanup 45 hari (configurable)
- **GitHub Actions**: ~6 minutes/day (3 runs × 2 min)
- **Monitoring**: Supabase Dashboard → Database → Usage

---

## 🆘 Need Help?

- Check `docs/MANUAL_MIGRATION.md` untuk detailed SQL
- Check GitHub Actions logs untuk error details
- Check Supabase Dashboard → Logs untuk database errors
- Open GitHub Issue jika ada masalah
