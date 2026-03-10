# GitHub Actions + Vercel Deployment Guide

**100% Free Solution - No Credit Card Required**

This deployment uses:
- **GitHub Actions**: Runs scraper every 30 minutes (free: 2000 min/month)
- **Vercel**: Hosts serverless API (free tier)
- **GitHub Repo**: Stores scraped data in `data/news.json`

## Architecture

```
GitHub Actions (Scraper)
    ↓ every 30 min
Scrapes news → Saves to data/news.json → Commits to repo
    ↓
Vercel API reads data/news.json
    ↓
Flutter app calls Vercel API
```

## Prerequisites

1. **GitHub Account** (you already have this)
2. **Vercel Account**: Sign up at [vercel.com](https://vercel.com) with GitHub
3. **OpenRouter API Key**: Get from [openrouter.ai](https://openrouter.ai)
4. **Neon.tech Database**: Already configured

## Step 1: Setup GitHub Secrets

Go to your repo: `https://github.com/hidatara-ds/evolipia-radar/settings/secrets/actions`

Add these secrets:

1. **DATABASE_URL**
   ```
   postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require
   ```

2. **LLM_API_KEY**
   ```
   your-openrouter-api-key-here
   ```

## Step 2: Push Code to GitHub

```bash
# Make sure you're on main branch
git checkout main

# Add all new files
git add .github/workflows/scrape.yml
git add cmd/worker-json/main.go
git add api/
git add data/news.json
git add vercel.json

# Commit
git commit -m "Add GitHub Actions scraper and Vercel API"

# Push
git push origin main
```

## Step 3: Test GitHub Actions Manually

1. Go to: `https://github.com/hidatara-ds/evolipia-radar/actions`
2. Click on "Scrape News" workflow
3. Click "Run workflow" → "Run workflow"
4. Wait for it to complete (~2-5 minutes)
5. Check if `data/news.json` was updated with new commit

## Step 4: Deploy to Vercel

### Option A: Vercel CLI (Recommended)

```bash
# Install Vercel CLI
npm install -g vercel

# Login
vercel login

# Deploy
vercel --prod
```

### Option B: Vercel Web UI

1. Go to [vercel.com/new](https://vercel.com/new)
2. Click "Import Git Repository"
3. Select `hidatara-ds/evolipia-radar`
4. Configure:
   - **Framework Preset**: Other
   - **Root Directory**: `./` (leave as is)
   - **Build Command**: (leave empty)
   - **Output Directory**: (leave empty)
5. Click "Deploy"

## Step 5: Verify Deployment

Once Vercel deploys, you'll get a URL like: `https://evolipia-radar.vercel.app`

Test the API:

```bash
# Health check
curl https://evolipia-radar.vercel.app/healthz

# Get news feed
curl https://evolipia-radar.vercel.app/api/news

# Get trending
curl https://evolipia-radar.vercel.app/api/trending

# Search
curl https://evolipia-radar.vercel.app/api/search?q=llm
```

## Step 6: Update Flutter App

Update `lib/config.dart`:

```dart
class ApiConfig {
  static const String baseUrl = 'https://evolipia-radar.vercel.app';
}
```

Update `lib/screens/feed_screen.dart` to use API instead of direct database:

```dart
Future<void> loadNews() async {
  setState(() {
    isLoading = true;
    error = null;
  });

  try {
    final response = await http.get(
      Uri.parse('${ApiConfig.baseUrl}/api/news'),
    );

    if (response.statusCode == 200) {
      final data = json.decode(response.body);
      if (data['success']) {
        final itemsJson = data['data']['items'] as List;
        setState(() {
          items = itemsJson.map((json) => NewsItem.fromJson(json)).toList();
          isLoading = false;
        });
      } else {
        throw Exception(data['error']);
      }
    } else {
      throw Exception('Failed to load news');
    }
  } catch (e) {
    setState(() {
      error = e.toString();
      isLoading = false;
    });
  }
}
```

## How It Works

### 1. GitHub Actions Scraper

- Runs every 30 minutes (configurable in `.github/workflows/scrape.yml`)
- Executes `cmd/worker-json/main.go`
- Scrapes news from sources (Hacker News, RSS, etc.)
- Generates AI summaries using OpenRouter
- Saves to `data/news.json`
- Commits and pushes to repo

### 2. Vercel API

- Serverless Go functions in `api/` folder
- Reads `data/news.json` from repo
- Provides REST endpoints:
  - `GET /api/news` - All news (with filters)
  - `GET /api/news/{id}` - Single item
  - `GET /api/trending` - Trending items
  - `GET /api/search?q=query` - Search

### 3. Flutter App

- Calls Vercel API endpoints
- No direct database access needed
- Works on all platforms (Android, iOS, Web)

## Monitoring

### GitHub Actions

View workflow runs:
```
https://github.com/hidatara-ds/evolipia-radar/actions
```

Check logs to see:
- Scraping progress
- Number of items fetched
- AI summary generation
- Commit status

### Vercel

View deployment logs:
```
https://vercel.com/dashboard
```

Check:
- API response times
- Error rates
- Traffic

## Cost Breakdown

- **GitHub Actions**: FREE (2000 min/month, using ~24 hours/month)
- **Vercel**: FREE (100GB bandwidth, 100 serverless invocations/day)
- **GitHub Storage**: FREE (data/news.json is tiny)
- **Neon.tech**: FREE (database)
- **OpenRouter**: ~$0.10-1/month (pay-as-you-go)

**Total: ~$0.10-1/month** (only LLM API costs)

## Customization

### Change Scraper Frequency

Edit `.github/workflows/scrape.yml`:

```yaml
schedule:
  - cron: '0 * * * *'  # Every hour
  # or
  - cron: '0 */2 * * *'  # Every 2 hours
  # or
  - cron: '0 0 * * *'  # Once daily at midnight
```

### Disable AI Summaries (Save LLM Costs)

In GitHub Secrets, set:
```
LLM_ENABLED=false
```

This will use extractive summaries only (free).

### Add More Data Sources

Edit `internal/connectors/` in your Golang code to add:
- More RSS feeds
- Custom JSON APIs
- Reddit
- Twitter/X
- etc.

## Troubleshooting

### GitHub Actions Not Running

1. Check workflow file syntax: `.github/workflows/scrape.yml`
2. Verify secrets are set correctly
3. Check Actions tab for error messages
4. Manually trigger: Actions → Scrape News → Run workflow

### Vercel API Returns Empty Data

1. Check if `data/news.json` exists in repo
2. Verify GitHub Actions ran successfully
3. Check Vercel logs for errors
4. Redeploy: `vercel --prod`

### Flutter App Shows No Data

1. Test API directly: `curl https://your-app.vercel.app/api/news`
2. Check API response format
3. Verify Flutter is using correct base URL
4. Check network logs in Flutter DevTools

### Database Connection Errors

1. Verify DATABASE_URL secret is correct
2. Check Neon.tech database is active
3. Test connection locally: `psql "postgresql://..."`

## Advanced: Custom Domain

1. Buy domain (e.g., from Namecheap, GoDaddy)
2. In Vercel dashboard:
   - Go to your project
   - Settings → Domains
   - Add your domain
   - Update DNS records as instructed
3. Update Flutter app with new domain

## Migration from Direct Database Access

If your Flutter app currently connects directly to Neon.tech:

1. Keep database config for now (backup)
2. Add API config
3. Update screens one by one to use API
4. Test thoroughly
5. Remove direct database access once confirmed working

## Next Steps

1. ✅ Push code to GitHub
2. ✅ Setup GitHub Secrets
3. ✅ Test GitHub Actions workflow
4. ✅ Deploy to Vercel
5. ✅ Update Flutter app
6. 🔄 Test end-to-end
7. 🔄 Monitor for 24 hours
8. 🔄 Optimize as needed

## Support

- GitHub Actions Docs: https://docs.github.com/en/actions
- Vercel Docs: https://vercel.com/docs
- OpenRouter Docs: https://openrouter.ai/docs

---

**You're all set!** This solution is 100% free (except minimal LLM costs) and requires no credit card verification.
