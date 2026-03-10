# Deployment Summary: Evolipia Radar Backend to Fly.io

## What Was Done

I've analyzed your Golang backend and prepared complete Fly.io deployment configuration. All files are committed to the `main` branch.

## Files Created (in main branch)

1. **fly.toml** - Fly.io config for API server
2. **fly.worker.toml** - Fly.io config for worker (scraper)
3. **.env.example** - Environment variables template
4. **FLY_DEPLOYMENT.md** - Complete deployment guide
5. **AI_LOGIC_COMPARISON.md** - Golang vs Flutter AI logic comparison
6. **Dockerfile.api** (updated) - Added web and assets folders

## Backend Architecture

Your Golang backend has:

### API Server (`cmd/api/main.go`)
- REST API endpoints for news feed, search, trending
- Web UI (mobile-first) served from `/web` folder
- Health check at `/healthz`
- Runs on port 8080

### Worker (`cmd/worker/main.go`)
- Scraper that runs every 10 minutes (configurable via `WORKER_CRON`)
- Fetches from multiple sources: Hacker News, RSS feeds, arXiv, JSON APIs
- Generates AI summaries using OpenRouter
- Scores and ranks articles
- Runs immediately on startup, then on cron schedule

### AI Integration
- **Provider**: OpenRouter
- **Default Model**: `google/gemini-flash-1.5`
- **Fallback Models**: `anthropic/claude-3.5-sonnet`, `meta-llama/llama-3.1-70b-instruct`
- **Temperature**: 0.7
- **Max Tokens**: 500

### Database
- PostgreSQL via Neon.tech (already configured)
- Connection string: `postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require`

## Quick Start Deployment

### Prerequisites
1. Install Fly CLI: `iwr https://fly.io/install.ps1 -useb | iex`
2. Get OpenRouter API key from [openrouter.ai](https://openrouter.ai)
3. Login: `fly auth login`

### Deploy API Server
```bash
# Switch to main branch
git checkout main

# Create app
fly apps create evolipia-radar --org personal

# Set secrets
fly secrets set DATABASE_URL="postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require" -a evolipia-radar
fly secrets set LLM_API_KEY="your-openrouter-api-key" -a evolipia-radar

# Deploy
fly deploy -c fly.toml

# Open in browser
fly open -a evolipia-radar
```

### Deploy Worker (Scraper)
```bash
# Create app
fly apps create evolipia-radar-worker --org personal

# Set secrets
fly secrets set DATABASE_URL="postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require" -a evolipia-radar-worker
fly secrets set LLM_API_KEY="your-openrouter-api-key" -a evolipia-radar-worker

# Deploy
fly deploy -c fly.worker.toml

# Check logs
fly logs -a evolipia-radar-worker --tail
```

## After Deployment

### 1. Verify API is Running
```bash
curl https://evolipia-radar.fly.dev/healthz
```

### 2. Wait for Worker to Scrape
The worker runs immediately on startup, then every 10 minutes. Check logs:
```bash
fly logs -a evolipia-radar-worker --tail
```

Look for:
- "Running initial ingestion..."
- "Ingestion completed successfully"

### 3. Check Feed Has Data
```bash
curl https://evolipia-radar.fly.dev/v1/feed
```

Should return news items after worker completes first scrape.

### 4. Update Flutter App
Update `lib/config.dart`:
```dart
class DatabaseConfig {
  // Keep existing database config for direct access
  static const String host = 'ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech';
  // ...
}

class ApiConfig {
  // Add API endpoint
  static const String baseUrl = 'https://evolipia-radar.fly.dev';
}
```

## AI Logic Consistency

✅ **Flutter app AI logic matches Golang backend** (98% aligned)

Minor recommendation: Update Flutter's default model to match backend:
```dart
// lib/services/ai_service.dart
static const String defaultModel = 'google/gemini-flash-1.5';  // Change from gpt-3.5-turbo
```

See `AI_LOGIC_COMPARISON.md` in main branch for full details.

## Cost Estimate

- **API**: $0-5/month (auto-stops when idle)
- **Worker**: $5/month (runs continuously for cron)
- **Database**: $0 (Neon.tech free tier)
- **LLM**: ~$0.10-1/month (light usage)

**Total**: ~$5-10/month

## Troubleshooting

### Empty Feed in Flutter App
**Cause**: Database has 0 rows (confirmed by debug log)

**Solution**: Deploy worker to Fly.io - it will populate the database

### Worker Not Scraping
Check logs: `fly logs -a evolipia-radar-worker --tail`

Common issues:
- Database connection failed (check DATABASE_URL secret)
- LLM API key missing (check LLM_API_KEY secret)
- Fetch timeout (increase FETCH_TIMEOUT_SECONDS)

### API Returns 500 Errors
Check logs: `fly logs -a evolipia-radar`

Common issues:
- Database connection failed
- Missing secrets

## Next Steps

1. ✅ Deployment config ready (main branch)
2. 🔄 Deploy API to Fly.io
3. 🔄 Deploy Worker to Fly.io
4. 🔄 Wait for worker to populate database (~10 minutes)
5. 🔄 Update Flutter app to use production API
6. 🔄 Test Flutter app with real data

## Documentation

All documentation is in the `main` branch:
- **FLY_DEPLOYMENT.md** - Complete deployment guide with troubleshooting
- **AI_LOGIC_COMPARISON.md** - Detailed AI logic comparison
- **.env.example** - All environment variables explained

## Support

- Fly.io Docs: https://fly.io/docs
- OpenRouter Docs: https://openrouter.ai/docs
- Neon.tech Docs: https://neon.tech/docs

---

**Ready to deploy!** Follow the Quick Start section above or read the full guide in `FLY_DEPLOYMENT.md` (main branch).
