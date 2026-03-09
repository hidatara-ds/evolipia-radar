# Evolipia Radar - Supabase Deploy Branch

Branch ini dioptimalkan untuk deployment dengan database Supabase dan timer-based scraping (3x sehari).

## Architecture

- **Database**: Supabase PostgreSQL (managed, always-on)
- **Worker**: GitHub Actions scheduled (07:00, 12:00, 19:00 WIB)
- **Frontend**: Flutter app (terpisah) langsung query Supabase via SDK
- **No API Server**: Flutter → Supabase directly (lebih simple)

## Schedule

Worker jalan otomatis 3 kali sehari:
- **07:00 WIB** (00:00 UTC) - Pagi
- **12:00 WIB** (05:00 UTC) - Siang  
- **19:00 WIB** (12:00 UTC) - Malam

## Setup

### 1. Supabase Setup

- Buat project di [supabase.com](https://supabase.com)
- Copy Database URL: Settings → Database → Connection string
- Format: `postgresql://postgres:[password]@db.xxx.supabase.co:5432/postgres`
- Tambahkan ke GitHub Secrets: `SUPABASE_DB_URL`

### 2. Run Migrations

```bash
export DATABASE_URL="postgresql://..."
migrate -path migrations -database "$DATABASE_URL" up
```

### 3. Test Worker Local

```bash
export DATABASE_URL="postgresql://..."
go run ./cmd/worker
```

### 4. GitHub Actions

- Workflow otomatis jalan sesuai schedule (3x/hari)
- Manual trigger: Actions tab → Scheduled News Scraper → Run workflow
- Optional: Centang "retention_cleanup" untuk bersihkan data lama

## Data Retention

- Data berita: 45 hari (auto-cleanup tersedia)
- Scrape logs: 60 hari (untuk audit)
- Cleanup manual via workflow_dispatch atau bisa di-schedule terpisah

## File Structure (Clean)

```
.
├── cmd/
│   └── worker/       # One-shot scraper (no API)
├── internal/         # Business logic (cleaned)
│   ├── config/
│   ├── db/
│   ├── dto/
│   ├── models/
│   ├── services/
│   ├── connectors/
│   ├── scoring/
│   ├── summarizer/
│   ├── normalizer/
│   └── security/
├── migrations/       # DB migrations + retention functions
├── configs/          # Config files
└── .github/workflows/# CI/CD simple
```

## Monitoring

- Cek scrape logs di Supabase:
  ```sql
  SELECT * FROM scrape_logs ORDER BY started_at DESC LIMIT 10;
  ```
- GitHub Actions logs untuk detail execution
- Supabase Dashboard untuk database metrics

## Flutter App

Flutter app terpisah akan:
- Pakai `supabase_flutter` SDK
- Direct query ke Supabase (no backend API)
- Row Level Security (RLS) untuk keamanan (setup di Supabase dashboard)

## Notes

- Branch ini tidak punya Web UI dan API server
- Tidak ada Kubernetes/Terraform/Docker (Supabase managed)
- Worker jalan 3x/hari (hemat GitHub Actions minutes)
- Database size monitoring: ~500MB limit, auto-cleanup available
