# Local Testing Guide

## Setup

### 1. Start Backend API Server
```bash
go run test_frontend.go
```
This starts the Go API server on `http://localhost:8080` with endpoints:
- `/api/news` - News feed
- `/metrics` - System metrics

### 2. Start Frontend Dev Server
```bash
npm run dev
```
This starts Next.js on `http://localhost:3000`

### 3. Configure API Base URL

Create `.env.local` file (already created):
```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

### 4. Open Browser
```
http://localhost:3000
```

## Troubleshooting

### Issue: 404 errors for `/api/news` and `/metrics`

**Cause:** Browser cache or environment variable not loaded

**Solution:**
1. **Hard refresh browser:** Press `Ctrl + Shift + R` (Windows/Linux) or `Cmd + Shift + R` (Mac)
2. **Clear browser cache:** Open DevTools (F12) → Network tab → Check "Disable cache"
3. **Restart dev server:** Stop (`Ctrl+C`) and run `npm run dev` again
4. **Check console:** Open DevTools → Console tab → Look for "🔧 API Base URL: http://localhost:8080"

### Issue: CORS errors

**Cause:** API server not allowing requests from `localhost:3000`

**Solution:** The `test_frontend.go` server already has CORS enabled via `api.EnableCORS(w)` in handlers.

### Issue: "Failed to load news"

**Possible causes:**
1. Backend API server not running → Start `go run test_frontend.go`
2. `data/news.json` file missing → Run scraper to generate it
3. Wrong API URL → Check `.env.local` file

## Verify Setup

### Check Backend API
```bash
curl http://localhost:8080/api/news
```
Should return JSON with news items.

### Check Frontend
1. Open http://localhost:3000
2. Open DevTools (F12) → Console
3. Look for: `🔧 API Base URL: http://localhost:8080`
4. Check Network tab → Should see requests to `localhost:8080/api/news`

## Expected Behavior

✅ **Metrics cards** show: Sources Processed, Filtered, Active Clusters, Avg Score
✅ **Topic filter** buttons: All, LLM, Vision, Data, Security, RL, Robotics
✅ **News feed** displays 20 articles with:
   - Rank badge (#1, #2, etc.)
   - Title (clickable)
   - Domain and timestamp
   - Tags (colored pills)
   - Summary text
   - Score indicator

## Production Deployment

For production (Vercel), the API endpoints are serverless functions:
- Frontend: `https://your-domain.vercel.app`
- API: `https://your-domain.vercel.app/api/news` (same domain)

No need for `NEXT_PUBLIC_API_BASE_URL` in production - it uses relative paths.
