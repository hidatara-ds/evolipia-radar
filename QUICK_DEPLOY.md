# Quick Deploy Commands

Copy-paste these commands to deploy Evolipia Radar to Fly.io.

## Prerequisites

```powershell
# Install Fly CLI (Windows PowerShell)
iwr https://fly.io/install.ps1 -useb | iex

# Login
fly auth login

# Switch to main branch
git checkout main
```

## Deploy API Server

```powershell
# Create app
fly apps create evolipia-radar --org personal

# Set secrets (replace YOUR_OPENROUTER_KEY with actual key)
fly secrets set DATABASE_URL="postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require" -a evolipia-radar

fly secrets set LLM_API_KEY="YOUR_OPENROUTER_KEY" -a evolipia-radar

# Deploy
fly deploy -c fly.toml

# Check status
fly status -a evolipia-radar

# View logs
fly logs -a evolipia-radar

# Open in browser
fly open -a evolipia-radar
```

## Deploy Worker (Scraper)

```powershell
# Create app
fly apps create evolipia-radar-worker --org personal

# Set secrets (replace YOUR_OPENROUTER_KEY with actual key)
fly secrets set DATABASE_URL="postgresql://evolipia-radar_owner:npg_ntTN8wojqf3R@ep-quiet-butterfly-a1qlqxqy.ap-southeast-1.aws.neon.tech/evolipia-radar?sslmode=require" -a evolipia-radar-worker

fly secrets set LLM_API_KEY="YOUR_OPENROUTER_KEY" -a evolipia-radar-worker

# Deploy
fly deploy -c fly.worker.toml

# Check logs (wait for scraping to complete)
fly logs -a evolipia-radar-worker --tail
```

## Verify Deployment

```powershell
# Test API health
curl https://evolipia-radar.fly.dev/healthz

# Check feed (should have data after worker runs)
curl https://evolipia-radar.fly.dev/v1/feed
```

## Update Flutter App

```powershell
# Switch back to mobile-app branch
git checkout mobile-app
```

Then update `lib/config.dart`:
```dart
class ApiConfig {
  static const String baseUrl = 'https://evolipia-radar.fly.dev';
}
```

## Monitor

```powershell
# API logs
fly logs -a evolipia-radar --tail

# Worker logs
fly logs -a evolipia-radar-worker --tail

# Status
fly status -a evolipia-radar
fly status -a evolipia-radar-worker
```

## Redeploy After Changes

```powershell
git checkout main
git pull

# Redeploy API
fly deploy -c fly.toml

# Redeploy worker
fly deploy -c fly.worker.toml
```

---

**That's it!** Your backend will be live at `https://evolipia-radar.fly.dev`

For detailed guide, see `FLY_DEPLOYMENT.md` in main branch.
