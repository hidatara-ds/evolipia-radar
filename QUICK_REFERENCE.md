# Quick Reference - Evolipia Radar Enhancements

## 🚀 Quick Start (30 seconds)

```bash
# Start everything
docker-compose up -d postgres
make migrate-up
go run ./cmd/api &
go run ./cmd/worker &

# Open browser
open http://localhost:8080
```

## 📦 New Features

### 10+ New Data Sources
- HuggingFace Trending Models
- Papers with Code
- LMSYS Chatbot Arena
- OpenAI Status Updates
- GitHub Trending AI/ML
- Anthropic/DeepMind/HF Blogs

### LLM Summarization
```bash
export LLM_ENABLED=true
export LLM_API_KEY=sk-or-v1-...
export LLM_MODEL=google/gemini-flash-1.5
```

### PWA Features
- Offline support
- "Add to Home Screen"
- Dark mode toggle

## 🔧 Configuration

### Environment Variables
```bash
# Database
DATABASE_URL=postgres://postgres:postgres@localhost:5432/radar?sslmode=disable

# API
PORT=8080
CACHE_TTL_SECONDS=60

# Worker
WORKER_CRON="*/10 * * * *"

# LLM (Optional)
LLM_ENABLED=false
LLM_API_KEY=
LLM_MODEL=google/gemini-flash-1.5
LLM_MAX_TOKENS=500
LLM_TEMPERATURE=0.7
```

### Source Types
| Type | Description |
|------|-------------|
| `hacker_news` | Hacker News API |
| `rss_atom` | RSS/Atom feeds |
| `arxiv` | arXiv papers |
| `json_api` | Custom JSON APIs |
| `huggingface_trending` | HF trending models |
| `papers_with_code` | Papers with Code API |
| `lmsys_arena` | LMSYS leaderboard |
| `openai_status` | OpenAI status RSS |
| `anthropic_docs` | Anthropic docs |
| `github_trending` | GitHub trending |

## 📡 API Endpoints

### Feed
```bash
GET /v1/feed?date=today&topic=llm
GET /v1/rising?window=2h
GET /v1/items/:id
GET /v1/search?q=RAG+systems
```

### Sources
```bash
GET /v1/sources
POST /v1/sources
POST /v1/sources/test
PATCH /v1/sources/:id/enable
```

## 🧪 Testing

### Test New Connectors
```bash
go run ./cmd/worker
# Check logs for "Processing source: HuggingFace Trending"
```

### Test LLM Summarization
```bash
export LLM_ENABLED=true
export LLM_API_KEY=your_key
go run ./cmd/worker
psql $DATABASE_URL -c "SELECT method, COUNT(*) FROM summaries GROUP BY method;"
```

### Test PWA
1. Open http://localhost:8080
2. DevTools > Application > Service Workers
3. Verify registration
4. Test offline mode

## 📁 File Structure

```
evolipia-radar/
├── internal/
│   ├── llm/              # ✅ LLM client
│   ├── connectors/       # ✅ New connectors
│   │   ├── huggingface.go
│   │   └── lmsys.go
│   ├── crawler/          # 🏗️ Phase 2 scaffold
│   ├── search/           # 🏗️ Phase 2 scaffold
│   └── realtime/         # 🏗️ Phase 2 scaffold
├── web/
│   ├── manifest.json     # ✅ PWA manifest
│   ├── sw.js             # ✅ Service worker
│   └── index.html        # ✅ Updated UI
├── docs/
│   ├── PHASE1_IMPLEMENTATION.md
│   ├── PHASE2_SCAFFOLD.md
│   ├── PHASE3_PLAN.md
│   └── ENHANCEMENTS_QUICKSTART.md
└── configs/
    └── default_sources.yaml  # ✅ New sources
```

## 🐛 Troubleshooting

### Worker not fetching new sources
```bash
# Check enabled sources
psql $DATABASE_URL -c "SELECT name, type, enabled FROM sources;"

# Enable a source
curl -X PATCH http://localhost:8080/v1/sources/{id}/enable \
  -H "Content-Type: application/json" \
  -d '{"enabled": true}'
```

### LLM not working
```bash
# Check env vars
echo $LLM_ENABLED
echo $LLM_API_KEY

# Test API key
curl https://openrouter.ai/api/v1/models \
  -H "Authorization: Bearer $LLM_API_KEY"
```

### PWA not installing
- Requires HTTPS (localhost OK for dev)
- Check manifest: http://localhost:8080/web/manifest.json
- Check service worker in DevTools

## 💰 Cost Estimates

### LLM Summarization
- Gemini Flash: Free tier available
- Claude 3.5 Sonnet: ~$0.001/summary
- 1000 items/day = ~$0-1/day

### Production (10k items/day)
- LLM: ~$10/day
- Infrastructure: ~$50/month
- Total: ~$320/month

## 📚 Documentation

- **Quick Start:** `docs/ENHANCEMENTS_QUICKSTART.md`
- **Phase 1:** `docs/PHASE1_IMPLEMENTATION.md`
- **Phase 2:** `docs/PHASE2_SCAFFOLD.md`
- **Phase 3:** `docs/PHASE3_PLAN.md`
- **Summary:** `ENHANCEMENTS_SUMMARY.md`
- **Complete:** `IMPLEMENTATION_COMPLETE.md`

## 🎯 Next Steps

### Phase 2 (Next 1-2 months)
- Intelligent crawler with JS rendering
- Vector search with pgvector
- Real-time WebSocket updates

### Phase 3 (Next 3-6 months)
- Admin dashboard
- Personalization engine
- Mobile apps
- Email digests
- Slack/Discord integration

## 🔗 Useful Commands

```bash
# Build
go build ./cmd/api
go build ./cmd/worker

# Test
go test ./...

# Lint
golangci-lint run

# Migrations
make migrate-up
make migrate-down

# Docker
docker-compose up -d
docker-compose down

# Database
psql $DATABASE_URL
```

## 📞 Support

1. Check `docs/` directory
2. Review code comments
3. Test with examples
4. Open GitHub issue

---

**Status:** ✅ Phase 1 Complete  
**Version:** 2.0.0  
**Last Updated:** March 8, 2026
