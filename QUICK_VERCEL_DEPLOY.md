# Quick Vercel Deployment Guide

## Status: ✅ READY TO DEPLOY

The source type issue has been **FIXED**! The worker now supports:
- `hackernews` (database) → `hacker_news` (connector) ✅
- `huggingface` (database) → `huggingface` (connector) ✅  
- `lmsys` (database) → `lmsys` (connector) ✅

## Deploy to Vercel

1. **Connect Repository**
   ```bash
   # Install Vercel CLI
   npm i -g vercel
   
   # Login and deploy
   vercel login
   vercel --prod
   ```

2. **Set Environment Variables**
   In Vercel dashboard, add:
   ```
   DATABASE_URL=postgresql://radar_owner:npg_ntTN8wojqf3R@ep-rough-darkness-a5qvqhqr.us-east-2.aws.neon.tech/radar?sslmode=require
   LLM_API_KEY=your_openrouter_key
   ```

3. **GitHub Actions**
   - Runs every 30 minutes automatically
   - Scrapes news from all 3 sources
   - Updates `data/news.json`
   - Vercel serves the JSON via API endpoints

## API Endpoints

- `GET /api/news` - All news items
- `GET /api/news?topic=ai` - Filter by topic
- `GET /api/news?date=today` - Today's news
- `GET /api/trending` - Trending items (last 2 hours, score > 0.5)
- `GET /api/search?q=openai` - Search news
- `GET /healthz` - Health check

## Test Locally

```bash
# Test the scraper
go run ./cmd/worker-json

# Test API endpoints
vercel dev
```

## 100% Free Solution ✅

- **GitHub Actions**: Free 2000 minutes/month
- **Vercel**: Free hosting + serverless functions
- **Neon.tech**: Free PostgreSQL database
- **OpenRouter**: Pay-per-use (very cheap)

No credit card verification required!