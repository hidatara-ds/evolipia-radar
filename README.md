# Evolipia Radar - Supabase Deploy Branch

Branch ini dioptimalkan untuk deployment dengan database Supabase dan timer-based scraping (3x sehari via GitHub Actions).

## 🚀 Quick Start (5 Menit)

### 1. Setup Supabase Database

1. Buat project di [supabase.com](https://supabase.com)
2. Copy connection string dari Settings → Database → URI
3. Run migrations manual via SQL Editor (lihat [SETUP_GUIDE.md](SETUP_GUIDE.md))

### 2. Setup GitHub Actions

1. Add GitHub Secret:
   - Name: `SUPABASE_DB_URL`
   - Value: Connection string dari Supabase
2. Push branch: `git push origin deploy-supabase`
3. Manual trigger: Actions → "Scheduled News Scraper" → Run workflow

### 3. Done! 🎉

Worker akan jalan otomatis 3x/hari:
- **07:00 WIB** (00:00 UTC)
- **12:00 WIB** (05:00 UTC)
- **19:00 WIB** (12:00 UTC)

## 📚 Documentation

- **[SETUP_GUIDE.md](SETUP_GUIDE.md)** - Complete setup guide (START HERE!)
- **[docs/MANUAL_MIGRATION.md](docs/MANUAL_MIGRATION.md)** - Database migrations
- **[DEPLOYMENT_CHECKLIST.md](DEPLOYMENT_CHECKLIST.md)** - Verification checklist
- **[CHANGELOG.md](CHANGELOG.md)** - Version history

## 🏗️ Architecture

```
┌─────────────────────────────────────┐
│      GitHub Actions (Scheduler)      │
│      3x/day: 07:00, 12:00, 19:00    │
│                                      │
│  ┌────────────────────────────┐    │
│  │  Worker (One-shot)          │    │
│  │  - Fetch from sources       │    │
│  │  - Process & score          │    │
│  │  - Save to DB               │    │
│  │  - Exit                     │    │
│  └────────────────────────────┘    │
└──────────────┬──────────────────────┘
               │
               ▼
    ┌──────────────────────┐
    │  Supabase PostgreSQL  │
    │  (Managed Database)   │
    │                       │
    │  - items              │
    │  - sources            │
    │  - scores             │
    │  - summaries          │
    │  - scrape_logs        │
    └──────────┬────────────┘
               │
               ▼
    ┌──────────────────────┐
    │   Flutter App         │
    │   (Supabase SDK)      │
    │                       │
    │  - Direct query       │
    │  - No backend API     │
    └───────────────────────┘
```

## ⚠️ Important Notes

### IPv6 Requirement

Supabase database requires IPv6 connectivity. For local testing:
- **Windows**: See [ENABLE_IPV6_WINDOWS.md](ENABLE_IPV6_WINDOWS.md)
- **Production**: GitHub Actions supports IPv6 by default ✅

### No Local Testing Required

You can skip local testing and deploy directly to GitHub Actions:
1. Setup Supabase database
2. Add GitHub Secret
3. Run workflow
4. Monitor via GitHub Actions logs

## 🎯 Key Features

- **Serverless**: No 24/7 server maintenance
- **Cost-effective**: Free tier Supabase + minimal GitHub Actions minutes
- **Simple**: Direct database access, no API layer
- **Scalable**: Supabase handles database scaling
- **Maintainable**: Auto cleanup (45 days retention)

## 📊 Database Schema

### Core Tables
- `sources` - News sources configuration
- `items` - Scraped news items
- `signals` - Engagement signals (points, comments)
- `scores` - Calculated relevance scores
- `summaries` - AI-generated summaries
- `scrape_logs` - Worker execution logs

### Helper Functions
- `cleanup_old_data()` - Auto cleanup (45 days)
- `get_daily_feed()` - Get recent items

## 🔧 Development

### Build Worker
```bash
go build -o worker ./cmd/worker
```

### Run Tests
```bash
go test ./...
```

### Verify Code
```bash
go vet ./...
```

## 📱 Flutter App Integration

Flutter app uses Supabase SDK (not direct PostgreSQL):

```dart
// Initialize
await Supabase.initialize(
  url: 'https://xxx.supabase.co',
  anonKey: 'your-anon-key',
);

// Query items
final items = await Supabase.instance.client
  .from('items')
  .select()
  .order('created_at', ascending: false)
  .limit(20);
```

## 🆘 Troubleshooting

### Worker Failed in GitHub Actions
- Check GitHub Actions logs
- Verify `SUPABASE_DB_URL` secret is set correctly
- Check Supabase project is not paused

### No Items Scraped
- Verify sources are enabled in database
- Check scrape_logs table for errors
- Verify migrations are run correctly

### Database Size Warning
```sql
-- Check size
SELECT pg_size_pretty(pg_database_size(current_database()));

-- Run cleanup
SELECT cleanup_old_data();
```

## 📈 Monitoring

### Check Scrape Logs
```sql
SELECT * FROM scrape_logs 
ORDER BY started_at DESC 
LIMIT 10;
```

### Check Recent Items
```sql
SELECT COUNT(*) FROM items 
WHERE created_at >= CURRENT_DATE;
```

### Database Size
```sql
SELECT pg_size_pretty(pg_database_size(current_database()));
```

## 🔗 Links

- [Supabase Dashboard](https://supabase.com/dashboard)
- [GitHub Actions](https://github.com/hidatara-ds/evolipia-radar/actions)
- [Supabase Docs](https://supabase.com/docs)

## 📄 License

See [LICENSE.md](LICENSE.md)

## 🤝 Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md)
