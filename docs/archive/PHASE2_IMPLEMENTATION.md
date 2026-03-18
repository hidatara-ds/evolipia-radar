# PHASE 2: Add News Sources - Implementation Complete ✅

## Overview
Added 18 high-quality AI/ML news sources from research papers, company blogs, tech news sites, and community aggregators.

## News Sources Added

### 📚 Research Papers (5 sources)
1. **arXiv AI Papers** - `cat:cs.AI` (Artificial Intelligence)
2. **arXiv Machine Learning** - `cat:cs.LG` (Machine Learning)
3. **arXiv Computer Vision** - `cat:cs.CV` (Computer Vision)
4. **arXiv NLP** - `cat:cs.CL` (Computational Linguistics)
5. **Papers with Code** - Latest ML papers with implementations

### 🏢 Company Blogs (5 sources)
6. **HuggingFace Blog** - Models, datasets, and ML tools
7. **Google AI Blog** - Research from Google AI/DeepMind
8. **OpenAI News** - GPT, DALL-E, and safety research
9. **Anthropic News** - Claude and AI safety
10. **DeepMind Blog** - AlphaFold, AlphaGo, and more

### 📰 Tech News Sites (4 sources)
11. **MIT Technology Review AI** - In-depth AI analysis
12. **VentureBeat AI** - AI industry news and startups
13. **The Verge AI** - Consumer AI and tech trends
14. **TechCrunch AI** - AI startup funding and launches

### 🔬 Research & Analysis (3 sources)
15. **The Gradient** - ML research analysis and essays
16. **Distill.pub** - Interactive ML explanations
17. **AI Alignment Forum** - AI safety and alignment research

### 🌐 Community (1 source)
18. **Hacker News** - Tech community discussions (already exists)

## Implementation

### File Created
`scripts/add_news_sources.go` - Script to populate database with sources

### Features
- ✅ Upsert logic (skip if source already exists by URL)
- ✅ Detailed logging with emojis for status
- ✅ Summary statistics (added, skipped, failed)
- ✅ All sources enabled by default
- ✅ Proper categorization (news, research)
- ✅ Support for RSS/Atom feeds and arXiv API

### Source Types
- `hacker_news` - Hacker News API
- `rss_atom` - Generic RSS/Atom feeds
- `arxiv` - arXiv API with category filters

## How to Run

### Prerequisites
1. Database connection configured in `.env` or environment variables
2. `DATABASE_URL` pointing to PostgreSQL database

### Execute Script
```bash
go run scripts/add_news_sources.go
```

### Expected Output
```
✅ Added source: HuggingFace Blog (rss_atom)
✅ Added source: Google AI Blog (rss_atom)
⏭️  Source already exists: Hacker News (https://news.ycombinator.com)
...
============================================================
📊 Summary:
   ✅ Added: 17 sources
   ⏭️  Skipped (already exists): 1 sources
   ❌ Failed: 0 sources
   📦 Total sources in database: 18
============================================================
```

## Database Schema
Sources are stored in the `sources` table:
```sql
CREATE TABLE sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,  -- hacker_news, rss_atom, arxiv
    category VARCHAR(50) NOT NULL,  -- news, research
    url TEXT NOT NULL UNIQUE,
    mapping_json JSONB,
    enabled BOOLEAN DEFAULT true,
    status VARCHAR(50) DEFAULT 'active',  -- active, pending, failed
    last_test_status VARCHAR(50),
    last_test_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## Integration with Worker

### Automatic Scraping
The worker (`cmd/worker-json/main.go`) automatically:
1. Fetches all enabled sources from database
2. Scrapes content based on source type
3. Stores items in `items` table
4. Generates AI summaries (if LLM enabled)
5. Calculates relevance scores
6. Exports top 100 items to `data/news.json`

### Scraping Schedule
- **GitHub Actions**: Every 30 minutes (`.github/workflows/scrape.yml`)
- **Manual Trigger**: Via `/v2/crawl/trigger` API endpoint

## Source Configuration

### RSS/Atom Feeds
```go
{
    Name:     "HuggingFace Blog",
    Type:     "rss_atom",
    Category: "news",
    URL:      "https://huggingface.co/blog/feed.xml",
    Enabled:  true,
    Status:   "active",
}
```

### arXiv API
```go
{
    Name:     "arXiv AI Papers",
    Type:     "arxiv",
    Category: "research",
    URL:      "http://export.arxiv.org/api/query?search_query=cat:cs.AI&sortBy=submittedDate&sortOrder=descending&max_results=50",
    Enabled:  true,
    Status:   "active",
}
```

## Testing

### Verify Sources Added
```sql
SELECT name, type, category, enabled, status 
FROM sources 
ORDER BY created_at DESC;
```

### Test Scraping
```bash
# Run worker to scrape all sources
go run ./cmd/worker-json

# Check output
cat data/news.json | jq '.items | length'
```

### Check Logs
```bash
# Worker logs show which sources were scraped
go run ./cmd/worker-json 2>&1 | grep "Processing source"
```

## Next Steps

### PHASE 3: Improve Scraper (Optional)
- Add auto-tagging based on keywords
- Implement deduplication by URL
- Add relevance scoring algorithm
- Retry logic with exponential backoff

### PHASE 4: UI Improvements (Optional)
- Add source badges in news cards
- Show source statistics
- Filter by source
- Topic velocity charts

## Troubleshooting

### Issue: Script times out
**Cause:** Cannot connect to database

**Solution:** 
1. Check `DATABASE_URL` environment variable
2. Ensure database is accessible
3. For local testing, use Neon.tech connection string from `.env.example`

### Issue: Sources not appearing in feed
**Cause:** Worker hasn't run yet

**Solution:**
1. Run worker manually: `go run ./cmd/worker-json`
2. Or trigger via API: `POST /v2/crawl/trigger`
3. Wait for GitHub Actions to run (every 30 min)

### Issue: Duplicate sources
**Cause:** URL changed or script run multiple times

**Solution:** Script has upsert logic - duplicates are automatically skipped by URL

## Production Deployment

### GitHub Actions
The scraper runs automatically via GitHub Actions:
```yaml
# .github/workflows/scrape.yml
- cron: '*/30 * * * *'  # Every 30 minutes
```

### Vercel
Sources are scraped in GitHub Actions, not on Vercel. Vercel only serves the API endpoints that read from `data/news.json`.

## Notes
- All sources are enabled by default
- Sources can be disabled via database: `UPDATE sources SET enabled = false WHERE name = 'Source Name'`
- New sources can be added by inserting into `sources` table
- Worker automatically picks up new enabled sources
- No code changes needed to add more RSS feeds

## Success Criteria
✅ 18 diverse AI/ML news sources configured
✅ Mix of research papers, company blogs, and tech news
✅ Upsert logic prevents duplicates
✅ All sources enabled and ready to scrape
✅ Integration with existing worker pipeline
✅ Documentation complete
