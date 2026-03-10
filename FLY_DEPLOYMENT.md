# Fly.io Deployment Guide for Evolipia Radar

This guide walks you through deploying the Evolipia Radar backend (API + Worker) to Fly.io.

## Prerequisites

1. **Fly.io Account**: Sign up at [fly.io](https://fly.io)
   - New users get $5 credit
   - After credit: ~$5/month for hobby usage
   
2. **Fly CLI**: Install flyctl
   ```bash
   # Windows (PowerShell)
   iwr https://fly.io/install.ps1 -useb | iex
   
   # macOS/Linux
   curl -L https://fly.io/install.sh | sh
   ```

3. **OpenRouter API Key**: Get from [openrouter.ai](https://openrouter.ai)
   - Sign up and create an API key
   - Add credits ($5-10 recommended for testing)

4. **Neon.tech Database**: Already configured
   - Connection string: `postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require`

## Architecture Overview

The deployment consists of two separate Fly.io apps:

1. **API Server** (`evolipia-radar`): Serves REST API and web UI
2. **Worker** (`evolipia-radar-worker`): Runs scraper every 10 minutes

Both connect to the same Neon.tech PostgreSQL database.

## Step 1: Login to Fly.io

```bash
fly auth login
```

## Step 2: Run Database Migrations

Before deploying, ensure your Neon.tech database has the latest schema:

```bash
# Install golang-migrate if not already installed
# Windows: scoop install migrate
# macOS: brew install golang-migrate
# Or download from: https://github.com/golang-migrate/migrate/releases

# Run migrations
migrate -path migrations -database "postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require" up
```

## Step 3: Deploy API Server

### 3.1 Create the API app

```bash
fly apps create evolipia-radar --org personal
```

### 3.2 Set secrets

```bash
# Database connection
fly secrets set DATABASE_URL="postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require" -a evolipia-radar

# OpenRouter API key (replace with your actual key)
fly secrets set LLM_API_KEY="sk-or-v1-your-api-key-here" -a evolipia-radar
```

### 3.3 Deploy

```bash
fly deploy -c fly.toml
```

### 3.4 Check status

```bash
fly status -a evolipia-radar
fly logs -a evolipia-radar
```

### 3.5 Open in browser

```bash
fly open -a evolipia-radar
```

Your API should be live at: `https://evolipia-radar.fly.dev`

## Step 4: Deploy Worker (Scraper)

### 4.1 Create the worker app

```bash
fly apps create evolipia-radar-worker --org personal
```

### 4.2 Set secrets

```bash
# Database connection
fly secrets set DATABASE_URL="postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require" -a evolipia-radar-worker

# OpenRouter API key (replace with your actual key)
fly secrets set LLM_API_KEY="sk-or-v1-your-api-key-here" -a evolipia-radar-worker
```

### 4.3 Deploy

```bash
fly deploy -c fly.worker.toml
```

### 4.4 Check logs

```bash
fly logs -a evolipia-radar-worker
```

You should see:
- "Worker started with cron schedule: */10 * * * *"
- "Running initial ingestion..."
- "Ingestion completed successfully"

## Step 5: Verify Deployment

### 5.1 Check API health

```bash
curl https://evolipia-radar.fly.dev/healthz
```

Expected response:
```json
{"status":"ok"}
```

### 5.2 Check feed endpoint

```bash
curl https://evolipia-radar.fly.dev/v1/feed
```

Should return news items (may be empty initially, wait for worker to run).

### 5.3 Check worker logs

```bash
fly logs -a evolipia-radar-worker --tail
```

Wait for the next cron run (every 10 minutes) and verify:
- "Starting scheduled ingestion..."
- "Ingestion completed successfully"

## Step 6: Update Flutter App

Update your Flutter app's `lib/config.dart`:

```dart
class ApiConfig {
  static const String baseUrl = 'https://evolipia-radar.fly.dev';
}
```

## Monitoring & Management

### View logs

```bash
# API logs
fly logs -a evolipia-radar

# Worker logs
fly logs -a evolipia-radar-worker

# Follow logs in real-time
fly logs -a evolipia-radar --tail
```

### Check machine status

```bash
fly status -a evolipia-radar
fly status -a evolipia-radar-worker
```

### Scale machines

```bash
# Scale API to 2 machines
fly scale count 2 -a evolipia-radar

# Scale back to 1
fly scale count 1 -a evolipia-radar
```

### Update environment variables

```bash
# Update LLM model
fly secrets set LLM_MODEL="anthropic/claude-3.5-sonnet" -a evolipia-radar
fly secrets set LLM_MODEL="anthropic/claude-3.5-sonnet" -a evolipia-radar-worker

# Disable LLM (use extractive summaries only)
fly secrets set LLM_ENABLED="false" -a evolipia-radar-worker
```

### Redeploy after code changes

```bash
# Commit your changes
git add .
git commit -m "Update feature"

# Redeploy API
fly deploy -c fly.toml

# Redeploy worker
fly deploy -c fly.worker.toml
```

## Cost Optimization

### Auto-stop machines when idle

The API is configured to auto-stop when idle and auto-start on requests:
- `auto_stop_machines = "stop"`
- `auto_start_machines = true`
- `min_machines_running = 0`

This means:
- API sleeps after ~5 minutes of inactivity
- Wakes up automatically on first request (cold start ~10-30s)
- No charges while sleeping

### Worker runs continuously

The worker doesn't auto-stop because it needs to run cron jobs. To reduce costs:

1. **Increase cron interval** (less frequent scraping):
   ```bash
   fly secrets set WORKER_CRON="0 */1 * * *" -a evolipia-radar-worker  # Every hour
   ```

2. **Stop worker when not needed**:
   ```bash
   fly scale count 0 -a evolipia-radar-worker  # Stop
   fly scale count 1 -a evolipia-radar-worker  # Start
   ```

## Troubleshooting

### API returns 500 errors

Check logs:
```bash
fly logs -a evolipia-radar
```

Common issues:
- Database connection failed (check DATABASE_URL secret)
- Missing LLM_API_KEY (check secrets)

### Worker not scraping

Check logs:
```bash
fly logs -a evolipia-radar-worker --tail
```

Common issues:
- Cron not triggering (check WORKER_CRON format)
- Database connection failed
- Fetch timeout (increase FETCH_TIMEOUT_SECONDS)

### Database connection errors

Verify connection string:
```bash
fly secrets list -a evolipia-radar
```

Test connection locally:
```bash
psql "postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require"
```

### Cold start too slow

Increase min_machines_running in `fly.toml`:
```toml
[http_service]
  min_machines_running = 1  # Keep 1 machine always running
```

Then redeploy:
```bash
fly deploy -c fly.toml
```

## AI Logic Consistency

The Golang backend uses:
- **Provider**: OpenRouter
- **Model**: `google/gemini-flash-1.5` (default)
- **Fallback**: `anthropic/claude-3.5-sonnet`, `meta-llama/llama-3.1-70b-instruct`
- **Temperature**: 0.7
- **Max Tokens**: 500

### System Prompt (Summarization)
```
You are an AI/ML news analyst. Provide concise, technical summaries focused on engineering impact.
```

### User Prompt Format
```
Summarize this AI/ML news article:

Title: [title]
Content: [content]

Provide:
1. A 2-sentence summary (TLDR)
2. One sentence explaining why this matters to AI/ML engineers

Format your response as:
TLDR: [your summary]
WHY: [why it matters]
```

**Flutter app should match these exactly** for consistency.

## Next Steps

1. ✅ Deploy API and Worker to Fly.io
2. ✅ Verify scraper is populating database
3. ✅ Test API endpoints
4. 🔄 Update Flutter app to use production API
5. 🔄 Test Flutter app with real data
6. 🔄 Monitor costs and optimize as needed

## Support

- Fly.io Docs: https://fly.io/docs
- Fly.io Community: https://community.fly.io
- OpenRouter Docs: https://openrouter.ai/docs
- Neon.tech Docs: https://neon.tech/docs

## Estimated Monthly Costs

- **API**: $0-5/month (mostly idle with auto-stop)
- **Worker**: $5/month (runs continuously for cron)
- **Database**: $0 (Neon.tech free tier)
- **LLM**: Variable (depends on usage, ~$0.10-1/month for light usage)

**Total**: ~$5-10/month
