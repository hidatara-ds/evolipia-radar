# Deployment Guide

Complete guide for deploying Evolipia Radar to production.

## Prerequisites

- GitHub account
- Vercel account (free tier works)
- Neon.tech account (free tier works)
- OpenRouter API key (optional, for AI features)

## Quick Deploy

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https://github.com/hidatara-ds/evolipia-radar)

## Manual Deployment

### Step 1: Database Setup (Neon.tech)

1. **Create Neon.tech Account**
   - Go to [neon.tech](https://neon.tech)
   - Sign up with GitHub

2. **Create New Project**
   - Click "New Project"
   - Name: `evolipia-radar`
   - Region: Choose closest to your users
   - PostgreSQL version: 15+

3. **Get Connection String**
   - Copy the connection string
   - Format: `postgresql://user:pass@host/dbname?sslmode=require`

4. **Run Migrations**
   ```bash
   # Install psql if needed
   # macOS: brew install postgresql
   # Ubuntu: sudo apt-get install postgresql-client
   
   # Apply schema
   psql "your-connection-string" < migrations/001_initial_schema.sql
   ```

5. **Verify Tables**
   ```bash
   psql "your-connection-string" -c "\dt"
   ```
   
   Should show: sources, items, summaries, scores, signals, fetch_runs

### Step 2: Vercel Deployment

1. **Install Vercel CLI**
   ```bash
   npm install -g vercel
   ```

2. **Login to Vercel**
   ```bash
   vercel login
   ```

3. **Link Project**
   ```bash
   vercel link
   ```
   
   - Create new project: Yes
   - Project name: `evolipia-radar`
   - Directory: `./` (current)

4. **Set Environment Variables**
   ```bash
   # Required
   vercel env add DATABASE_URL production
   # Paste your Neon.tech connection string
   
   # Optional (for AI features)
   vercel env add LLM_API_KEY production
   vercel env add LLM_PROVIDER production
   vercel env add LLM_MODEL production
   vercel env add LLM_ENABLED production
   ```

5. **Deploy**
   ```bash
   vercel --prod
   ```

6. **Verify Deployment**
   - Open the provided URL
   - Check `/api/news` endpoint
   - Verify dashboard loads

### Step 3: GitHub Actions Setup

1. **Add Repository Secrets**
   - Go to GitHub repository → Settings → Secrets and variables → Actions
   - Add secrets:
     - `DATABASE_URL`: Your Neon.tech connection string
     - `LLM_API_KEY`: Your OpenRouter API key (optional)

2. **Enable GitHub Actions**
   - Go to Actions tab
   - Enable workflows if disabled

3. **Verify Scraper**
   - Go to Actions → Scrape News
   - Click "Run workflow"
   - Wait for completion
   - Check database for new articles

### Step 4: Populate Initial Data

Run the sample data script to populate initial articles:

```bash
# Set DATABASE_URL
export DATABASE_URL="your-connection-string"

# Run populate script
go run scripts/populate_sample_data.go
```

Or manually trigger the scraper:
```bash
curl -X POST https://your-domain.vercel.app/v2/crawl/trigger
```

### Step 5: Configure Custom Domain (Optional)

1. **Add Domain in Vercel**
   - Go to Project Settings → Domains
   - Add your domain
   - Follow DNS configuration instructions

2. **Update Environment Variables**
   ```bash
   vercel env add NEXT_PUBLIC_API_BASE_URL production
   # Value: https://your-domain.com
   ```

3. **Redeploy**
   ```bash
   vercel --prod
   ```

## Environment Variables Reference

### Required

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | `postgresql://user:pass@host/db` |

### Optional

| Variable | Description | Default |
|----------|-------------|---------|
| `LLM_API_KEY` | OpenRouter API key | - |
| `LLM_PROVIDER` | LLM provider | `openrouter` |
| `LLM_MODEL` | Model to use | `google/gemini-flash-1.5` |
| `LLM_ENABLED` | Enable AI features | `false` |
| `LLM_MAX_TOKENS` | Max tokens for generation | `500` |
| `LLM_TEMPERATURE` | Generation temperature | `0.7` |
| `NEXT_PUBLIC_API_BASE_URL` | API base URL | `` (relative) |

## Vercel Configuration

The `vercel.json` file configures:

```json
{
  "version": 2,
  "builds": [
    {
      "src": "package.json",
      "use": "@vercel/next"
    },
    {
      "src": "api/**/index.go",
      "use": "@vercel/go"
    }
  ],
  "routes": [
    {
      "src": "/api/news",
      "dest": "/api/news/index.go"
    }
  ]
}
```

## GitHub Actions Configuration

The `.github/workflows/scrape.yml` file configures:

- **Schedule:** Every 30 minutes (`*/30 * * * *`)
- **Manual trigger:** Available via Actions tab
- **Permissions:** Write access to commit JSON files

## Monitoring

### Vercel Dashboard

1. **Logs**
   - Go to Deployments → Select deployment → Logs
   - Filter by function: `/api/news`
   - Check for errors

2. **Analytics**
   - Go to Analytics tab
   - Monitor request count, response time
   - Check error rate

3. **Functions**
   - Go to Functions tab
   - Monitor invocations, duration
   - Check cold start frequency

### Neon.tech Dashboard

1. **Metrics**
   - Go to project → Metrics
   - Monitor connections, queries
   - Check storage usage

2. **Query Performance**
   - Go to Monitoring → Queries
   - Identify slow queries
   - Optimize indexes

### GitHub Actions

1. **Workflow Runs**
   - Go to Actions tab
   - Check scraper success rate
   - Monitor run duration

2. **Logs**
   - Click on workflow run
   - Check "Run scraper" step
   - Verify articles inserted

## Troubleshooting

### API Returns "Database configuration missing"

**Cause:** DATABASE_URL not set in Vercel

**Solution:**
```bash
vercel env add DATABASE_URL production
# Paste connection string
vercel --prod
```

### API Returns "Failed to connect to database"

**Cause:** Invalid connection string or database down

**Solution:**
1. Verify connection string format
2. Check Neon.tech dashboard for database status
3. Test connection locally:
   ```bash
   psql "your-connection-string" -c "SELECT 1"
   ```

### Cold Start Timeout

**Cause:** First request takes >15 seconds

**Solution:**
- Already implemented: Auto-retry in frontend
- Keep function warm: Set up uptime monitor (e.g., UptimeRobot)
- Upgrade Vercel plan for faster cold starts

### GitHub Actions Scraper Fails

**Cause:** DATABASE_URL secret not set

**Solution:**
1. Go to repository Settings → Secrets
2. Add `DATABASE_URL` secret
3. Re-run workflow

### No Articles Showing

**Cause:** Database empty or scraper not running

**Solution:**
1. Check database:
   ```bash
   psql "your-connection-string" -c "SELECT COUNT(*) FROM items"
   ```
2. Run populate script:
   ```bash
   go run scripts/populate_sample_data.go
   ```
3. Manually trigger scraper:
   ```bash
   curl -X POST https://your-domain.vercel.app/v2/crawl/trigger
   ```

### Build Fails on Vercel

**Cause:** Missing dependencies or build errors

**Solution:**
1. Check build logs in Vercel dashboard
2. Verify `go.mod` and `package.json` are committed
3. Test build locally:
   ```bash
   npm run build
   go build ./api/news/index.go
   ```

## Performance Optimization

### Database

1. **Add Indexes**
   ```sql
   CREATE INDEX idx_items_published ON items(published_at DESC);
   CREATE INDEX idx_summaries_tags ON summaries USING GIN(tags);
   CREATE INDEX idx_scores_final ON scores(final DESC);
   ```

2. **Connection Pooling**
   - Already configured in `api/news/index.go`
   - Max connections: 3
   - Idle timeout: 30s

### API

1. **Caching**
   - Add Redis for response caching
   - Cache duration: 30 seconds
   - Invalidate on new articles

2. **CDN**
   - Vercel automatically uses CDN
   - Static assets cached at edge

### Frontend

1. **Code Splitting**
   - Already implemented with Next.js
   - Dynamic imports for heavy components

2. **Image Optimization**
   - Use Next.js Image component
   - Lazy loading enabled

## Scaling

### Current Limits (Free Tier)

- **Vercel:**
  - 100 GB bandwidth/month
  - 100 GB-hours compute/month
  - 6,000 build minutes/month

- **Neon.tech:**
  - 3 GB storage
  - 1 compute unit
  - 100 hours compute/month

### Upgrade Path

**For 1,000+ users:**
1. Upgrade Vercel to Pro ($20/month)
2. Add Redis caching
3. Optimize database queries

**For 10,000+ users:**
1. Upgrade Neon.tech to Scale plan
2. Add read replicas
3. Implement CDN caching
4. Consider dedicated infrastructure

## Backup & Recovery

### Database Backups

Neon.tech provides:
- Automatic daily backups
- 7-day retention
- Point-in-time recovery

**Manual Backup:**
```bash
pg_dump "your-connection-string" > backup.sql
```

**Restore:**
```bash
psql "your-connection-string" < backup.sql
```

### Code Backups

- GitHub repository is the source of truth
- All deployments are versioned
- Rollback via Vercel dashboard

### Recovery Procedure

1. **Database Failure:**
   - Restore from Neon.tech backup
   - Re-run scraper to catch up

2. **Deployment Failure:**
   - Rollback in Vercel dashboard
   - Or redeploy from GitHub

3. **Complete Disaster:**
   - Create new Neon.tech database
   - Restore from backup
   - Redeploy from GitHub
   - Re-run scraper

## Security Checklist

- [ ] DATABASE_URL stored as secret (not in code)
- [ ] LLM_API_KEY stored as secret
- [ ] HTTPS enabled (automatic with Vercel)
- [ ] CORS configured properly
- [ ] SQL injection prevention (prepared statements)
- [ ] Rate limiting considered
- [ ] Error messages don't leak sensitive info
- [ ] Dependencies regularly updated

## Post-Deployment

1. **Test All Endpoints**
   ```bash
   curl https://your-domain.vercel.app/api/news
   curl https://your-domain.vercel.app/metrics
   curl https://your-domain.vercel.app/healthz
   ```

2. **Monitor for 24 Hours**
   - Check Vercel logs
   - Monitor error rate
   - Verify scraper runs

3. **Set Up Alerts**
   - Vercel: Enable email notifications
   - Neon.tech: Enable usage alerts
   - GitHub: Watch repository for issues

4. **Document Custom Configuration**
   - Update README with your domain
   - Document any custom changes
   - Share with team

## Support

- **Documentation:** [GitHub Docs](https://github.com/hidatara-ds/evolipia-radar/tree/main/docs)
- **Issues:** [GitHub Issues](https://github.com/hidatara-ds/evolipia-radar/issues)
- **Vercel Support:** [Vercel Help](https://vercel.com/help)
- **Neon.tech Support:** [Neon Docs](https://neon.tech/docs)
