# Enhancements Quick Start Guide

## What's New

Evolipia Radar has been significantly enhanced with:
- **10+ new AI/ML data sources** (HuggingFace, Papers with Code, LMSYS Arena, etc.)
- **LLM-powered summarization** via OpenRouter
- **Modern PWA UI** with Tailwind CSS and dark mode
- **Phase 2 & 3 scaffolds** for advanced features

## Quick Start (5 minutes)

### 1. Update Dependencies
```bash
go mod tidy
```

### 2. Set Up LLM Summarization (Optional)
```bash
# Get API key from https://openrouter.ai/
export LLM_ENABLED=true
export LLM_API_KEY=your_openrouter_key_here
export LLM_MODEL=google/gemini-flash-1.5
```

### 3. Start the System
```bash
# Terminal 1: Start PostgreSQL
docker-compose up -d postgres

# Terminal 2: Run migrations
make migrate-up

# Terminal 3: Start API server
go run ./cmd/api

# Terminal 4: Start worker
go run ./cmd/worker
```

### 4. Access the UI
Open http://localhost:8080 in your browser

**New Features:**
- 📱 PWA support - "Add to Home Screen"
- 🌙 Dark mode toggle in Settings
- 🤖 AI Chat with OpenRouter (set API key in Settings)
- 📰 10+ new AI/ML sources

## New Data Sources

### Automatically Enabled
These sources are configured in `configs/default_sources.yaml`:

1. **HuggingFace Trending** - Top 50 trending models
2. **Papers with Code** - Latest research papers
3. **LMSYS Chatbot Arena** - LLM leaderboard
4. **OpenAI Status** - API status updates
5. **GitHub Trending** - Trending AI/ML repos
6. **Anthropic Blog** - Claude updates
7. **DeepMind Blog** - Research announcements
8. **Hugging Face Blog** - Platform updates

### Manual Addition
Add more sources via API:

```bash
# Example: Add Vercel AI SDK blog
curl -X POST http://localhost:8080/v1/sources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Vercel AI SDK Blog",
    "type": "rss_atom",
    "category": "news",
    "url": "https://vercel.com/blog/category/ai/rss.xml",
    "enabled": true
  }'
```

## LLM Summarization

### How It Works
1. Worker fetches new items
2. If `LLM_ENABLED=true`, uses OpenRouter for abstractive summaries
3. Falls back to extractive summarization on error
4. Summaries include "Why it matters" insights

### Cost Estimation
- Gemini Flash: Free tier available (default)
- Claude 3.5 Sonnet: ~$0.001 per summary
- 1000 items/day = ~$0-1/day depending on model

### Supported Models
```bash
# Default (Gemini Flash - Free tier)
LLM_MODEL=google/gemini-flash-1.5

# Premium (best quality)
LLM_MODEL=anthropic/claude-3.5-sonnet

# Balanced
LLM_MODEL=google/gemini-flash-1.5

# Budget (free tier)
LLM_MODEL=meta-llama/llama-3.1-8b-instruct:free
```

## PWA Features

### Install as App
1. Open http://localhost:8080 in Chrome/Edge
2. Click "Install" icon in address bar
3. App opens in standalone window

### Offline Support
- Service worker caches static assets
- Works offline after first visit
- Syncs when connection restored

### Dark Mode
1. Open Settings panel
2. Click "Dark Mode" toggle
3. Preference persists across sessions

## Testing New Features

### Test New Connectors
```bash
# Check worker logs
go run ./cmd/worker

# Look for:
# "Processing source: HuggingFace Trending"
# "Fetched X items from Papers with Code"
```

### Test LLM Summarization
```bash
# Enable LLM
export LLM_ENABLED=true
export LLM_API_KEY=your_key

# Run worker
go run ./cmd/worker

# Check database
psql $DATABASE_URL -c "SELECT method, COUNT(*) FROM summaries GROUP BY method;"

# Should show:
#   method   | count
# -----------+-------
#  extractive|   100
#  llm       |    50
```

### Test PWA
1. Open DevTools > Application > Service Workers
2. Verify "evolipia-radar-v1" is registered
3. Go offline (DevTools > Network > Offline)
4. Refresh page - should still load

## Configuration Reference

### Environment Variables

```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/radar?sslmode=disable

# API Server
PORT=8080
CACHE_TTL_SECONDS=60

# Worker
WORKER_CRON="*/10 * * * *"  # Every 10 minutes

# Fetching
MAX_FETCH_BYTES=2000000      # 2MB
FETCH_TIMEOUT_SECONDS=8

# LLM (Optional)
LLM_ENABLED=false
LLM_PROVIDER=openrouter
LLM_MODEL=google/gemini-flash-1.5
LLM_FALLBACK_MODELS=anthropic/claude-3.5-sonnet,meta-llama/llama-3.1-70b-instruct
LLM_API_KEY=
LLM_MAX_TOKENS=500
LLM_TEMPERATURE=0.7
```

### Source Types

| Type | Description | Example |
|------|-------------|---------|
| `hacker_news` | Hacker News API | N/A (no URL needed) |
| `rss_atom` | RSS/Atom feeds | `https://openai.com/blog/rss.xml` |
| `arxiv` | arXiv papers | N/A (auto-queries AI/ML categories) |
| `json_api` | Custom JSON APIs | Requires `mapping_json` |
| `huggingface_trending` | HF trending models | N/A |
| `papers_with_code` | Papers with Code API | N/A |
| `lmsys_arena` | LMSYS leaderboard | N/A |
| `openai_status` | OpenAI status RSS | N/A |
| `anthropic_docs` | Anthropic docs scraper | N/A |
| `github_trending` | GitHub trending | N/A |

## Troubleshooting

### Worker Not Fetching New Sources
```bash
# Check if sources are enabled
psql $DATABASE_URL -c "SELECT name, type, enabled FROM sources;"

# Enable a source
curl -X PATCH http://localhost:8080/v1/sources/{id}/enable \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

### LLM Summarization Not Working
```bash
# Check environment variables
echo $LLM_ENABLED
echo $LLM_API_KEY

# Test API key
curl https://openrouter.ai/api/v1/models \
  -H "Authorization: Bearer $LLM_API_KEY"

# Check worker logs for errors
go run ./cmd/worker 2>&1 | grep -i "llm\|summarization"
```

### PWA Not Installing
- Requires HTTPS in production (localhost is OK for dev)
- Check manifest.json is accessible: http://localhost:8080/web/manifest.json
- Check service worker: DevTools > Application > Service Workers

### Dark Mode Not Persisting
- Check localStorage: `localStorage.getItem('dark_mode')`
- Clear browser cache and try again
- Check browser console for errors

## Next Steps

### Phase 2: Advanced Features
See `docs/PHASE2_SCAFFOLD.md` for:
- Intelligent crawler with JavaScript rendering
- Vector search with pgvector
- Real-time WebSocket updates

### Phase 3: Production Polish
See `docs/PHASE3_PLAN.md` for:
- Admin dashboard
- Personalization engine
- Mobile apps (Capacitor)
- Email digests
- Slack/Discord integration

## API Examples

### Get Feed with New Sources
```bash
# Get today's top items
curl http://localhost:8080/v1/feed?date=today

# Filter by category
curl http://localhost:8080/v1/feed?date=today&topic=llm

# Get rising items
curl http://localhost:8080/v1/rising?window=2h
```

### Search
```bash
# Full-text search
curl http://localhost:8080/v1/search?q=RAG+systems

# Search with topic filter
curl http://localhost:8080/v1/search?q=transformers&topic=research
```

### Manage Sources
```bash
# List all sources
curl http://localhost:8080/v1/sources

# Test a source before adding
curl -X POST http://localhost:8080/v1/sources/test \
  -H "Content-Type: application/json" \
  -d '{
    "type": "rss_atom",
    "url": "https://example.com/feed.xml"
  }'

# Add source
curl -X POST http://localhost:8080/v1/sources \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Blog",
    "type": "rss_atom",
    "category": "news",
    "url": "https://example.com/feed.xml",
    "enabled": true
  }'
```

## Performance Tips

### Optimize Worker Schedule
```bash
# More frequent updates (every 5 minutes)
export WORKER_CRON="*/5 * * * *"

# Less frequent (every hour)
export WORKER_CRON="0 * * * *"

# Only during business hours (9 AM - 5 PM)
export WORKER_CRON="0 9-17 * * *"
```

### LLM Summarization
```bash
# Use cheaper model (default - free tier)
export LLM_MODEL=google/gemini-flash-1.5

# Or premium model
export LLM_MODEL=anthropic/claude-3.5-sonnet

# Reduce token limit
export LLM_MAX_TOKENS=300

# Only summarize high-scored items
# (Modify worker.go to check score before summarizing)
```

### Cache Optimization
```bash
# Increase cache TTL for less frequent updates
export CACHE_TTL_SECONDS=300  # 5 minutes

# Decrease for real-time feel
export CACHE_TTL_SECONDS=30   # 30 seconds
```

## Support

### Documentation
- `docs/PHASE1_IMPLEMENTATION.md` - What's implemented
- `docs/PHASE2_SCAFFOLD.md` - Next features
- `docs/PHASE3_PLAN.md` - Future roadmap
- `docs/LOCAL_SETUP.md` - Detailed setup guide
- `docs/DEPENDENCIES.md` - Dependency management

### Issues
- Check existing issues on GitHub
- Include logs and environment details
- Provide steps to reproduce

### Contributing
See `CONTRIBUTING.md` for guidelines.
