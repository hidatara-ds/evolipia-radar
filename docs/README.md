# Evolipia Radar Documentation

Documentation untuk branch `deploy-supabase`.

## Quick Links

- [Supabase Setup Guide](SUPABASE_SETUP.md) - Panduan lengkap setup dari awal
- [Local Setup](LOCAL_SETUP.md) - Setup untuk development lokal
- [Dependencies](DEPENDENCIES.md) - Daftar dependencies yang digunakan
- [Quick Start](QUICK_START.md) - Quick start guide

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                    GitHub Actions                        │
│  (Scheduled 3x/day: 07:00, 12:00, 19:00 WIB)           │
│                                                          │
│  ┌──────────────────────────────────────────┐          │
│  │  Worker (One-shot execution)              │          │
│  │  - Fetch from sources                     │          │
│  │  - Process & score                        │          │
│  │  - Save to DB                             │          │
│  │  - Exit                                   │          │
│  └──────────────────────────────────────────┘          │
└─────────────────────┬───────────────────────────────────┘
                      │
                      ▼
         ┌────────────────────────┐
         │   Supabase PostgreSQL   │
         │   (Managed Database)    │
         │                         │
         │  - items table          │
         │  - scrape_logs table    │
         │  - RLS policies         │
         │  - Functions            │
         └────────────┬────────────┘
                      │
                      ▼
         ┌────────────────────────┐
         │   Flutter App           │
         │   (Direct Query)        │
         │                         │
         │  - supabase_flutter SDK │
         │  - No backend API       │
         └─────────────────────────┘
```

## Key Features

- **Timer-based Scraping**: 3x sehari via GitHub Actions (hemat resources)
- **Serverless**: Tidak perlu maintain server 24/7
- **Simple Architecture**: Flutter → Supabase directly (no API layer)
- **Auto Cleanup**: Data retention 45 hari (otomatis)
- **Cost Effective**: Free tier Supabase + minimal GitHub Actions minutes

## File Structure

```
.
├── cmd/
│   └── worker/              # One-shot scraper
│       └── main.go
├── internal/
│   ├── config/              # Configuration
│   ├── db/                  # Database connection
│   ├── models/              # Data models
│   ├── services/            # Business logic
│   ├── connectors/          # Source connectors
│   ├── scoring/             # Scoring algorithm
│   ├── summarizer/          # Text summarization
│   ├── normalizer/          # Data normalization
│   └── security/            # Security utilities
├── migrations/              # Database migrations
│   ├── 000001_init_schema.up.sql
│   ├── 000002_create_scrape_logs_table.up.sql
│   └── 000003_create_data_retention_function.up.sql
├── configs/
│   ├── default_sources.yaml # Default news sources
│   └── supabase-setup.sql   # RLS & functions setup
├── .github/workflows/
│   ├── scheduled-scraper.yml # Main scraper workflow
│   ├── pr-check.yml          # PR validation
│   └── weekly-cleanup.yml    # Data cleanup
└── docs/                    # Documentation
```

## Workflows

### Scheduled Scraper
- **File**: `.github/workflows/scheduled-scraper.yml`
- **Schedule**: 3x/day (07:00, 12:00, 19:00 WIB)
- **Manual**: Bisa trigger manual dengan option cleanup

### PR Check
- **File**: `.github/workflows/pr-check.yml`
- **Trigger**: Setiap PR ke main/deploy-supabase
- **Actions**: Build, vet, test

### Weekly Cleanup
- **File**: `.github/workflows/weekly-cleanup.yml`
- **Schedule**: Minggu 09:00 WIB
- **Actions**: Cleanup data lama, check DB size

## Database Schema

### items table
```sql
- id (UUID)
- title (TEXT)
- url (TEXT)
- summary (TEXT)
- score (FLOAT)
- published_at (TIMESTAMP)
- source (VARCHAR)
- topic (VARCHAR)
- created_at (TIMESTAMP)
```

### scrape_logs table
```sql
- id (UUID)
- started_at (TIMESTAMP)
- completed_at (TIMESTAMP)
- items_processed (INTEGER)
- items_new (INTEGER)
- status (VARCHAR) -- running, success, failed
- error_message (TEXT)
- trigger_source (VARCHAR)
```

## Environment Variables

```bash
# Required
DATABASE_URL=postgresql://postgres:[password]@db.xxx.supabase.co:5432/postgres

# Optional
MAX_FETCH_BYTES=2000000
FETCH_TIMEOUT_SECONDS=8
LOG_LEVEL=debug
TRIGGER_SOURCE=github_actions
```

## Monitoring

### Scrape Logs
```sql
SELECT * FROM scrape_logs ORDER BY started_at DESC LIMIT 10;
```

### Database Size
```sql
SELECT pg_size_pretty(pg_database_size(current_database()));
```

### Recent Items
```sql
SELECT COUNT(*) FROM items WHERE published_at >= CURRENT_DATE;
```

## Maintenance

### Data Retention
- Items: 45 hari (auto-cleanup via function)
- Scrape logs: 60 hari
- Manual cleanup: Via GitHub Actions workflow

### Database Limits
- Supabase Free Tier: 500MB
- Monitor via SQL: `SELECT pg_size_pretty(pg_database_size(current_database()));`
- Auto-cleanup helps stay within limits

## Development

### Local Testing
```bash
# Setup
export DATABASE_URL="postgresql://..."
go mod download

# Run worker
go run ./cmd/worker

# Build
make build

# Test
make test
```

### Adding New Sources
1. Edit `configs/default_sources.yaml`
2. Add connector di `internal/connectors/`
3. Register di `internal/services/worker.go`
4. Test locally
5. Deploy via git push

## Support

- Issues: GitHub Issues
- Docs: `/docs` folder
- Supabase: [supabase.com/docs](https://supabase.com/docs)
