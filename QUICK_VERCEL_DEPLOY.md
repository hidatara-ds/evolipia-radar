# Quick Vercel Deploy (Copy-Paste Commands)

## 1. Setup GitHub Secrets

Go to: `https://github.com/hidatara-ds/evolipia-radar/settings/secrets/actions`

Click "New repository secret" and add:

**DATABASE_URL:**
```
postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require
```

**LLM_API_KEY:**
```
your-openrouter-api-key-here
```

## 2. Push Code

```bash
git checkout main
git add .
git commit -m "Add GitHub Actions + Vercel deployment"
git push origin main
```

## 3. Test GitHub Actions

Go to: `https://github.com/hidatara-ds/evolipia-radar/actions`

Click "Scrape News" → "Run workflow" → "Run workflow"

Wait 2-5 minutes, check if `data/news.json` updated.

## 4. Deploy to Vercel

### Option A: Web UI (Easiest)

1. Go to: https://vercel.com/new
2. Sign in with GitHub
3. Import `hidatara-ds/evolipia-radar`
4. Click "Deploy"
5. Done! Get your URL: `https://evolipia-radar.vercel.app`

### Option B: CLI

```bash
npm install -g vercel
vercel login
vercel --prod
```

## 5. Test API

```bash
curl https://evolipia-radar.vercel.app/healthz
curl https://evolipia-radar.vercel.app/api/news
```

## 6. Update Flutter

In `lib/config.dart`:

```dart
class ApiConfig {
  static const String baseUrl = 'https://evolipia-radar.vercel.app';
}
```

## Done!

- Scraper runs every 30 minutes automatically
- API is live at Vercel
- 100% free (except ~$0.10/month for LLM)

See `VERCEL_DEPLOYMENT.md` for full guide.
