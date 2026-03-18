# Architecture Guide

## System Overview

Evolipia Radar is a distributed system designed for real-time AI/ML news aggregation and analysis.

## High-Level Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                     GitHub Actions                            │
│  ┌────────────────────────────────────────────────────────┐  │
│  │  Scheduled Worker (every 30 minutes)                   │  │
│  │  - Fetch from 40+ sources                              │  │
│  │  - Extract & normalize content                         │  │
│  │  - AI-powered summarization                            │  │
│  │  - Auto-tagging & scoring                              │  │
│  └─────────────────────┬──────────────────────────────────┘  │
└────────────────────────┼─────────────────────────────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │      Neon.tech PostgreSQL     │
         │  ┌─────────────────────────┐  │
         │  │  Tables:                │  │
         │  │  - sources              │  │
         │  │  - items                │  │
         │  │  - summaries            │  │
         │  │  - scores               │  │
         │  │  - signals              │  │
         │  └─────────────────────────┘  │
         └───────────────┬───────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │      Vercel Serverless        │
         │  ┌─────────────────────────┐  │
         │  │  Go Functions:          │  │
         │  │  - /api/news            │  │
         │  │  - /api/trending        │  │
         │  │  - /api/search          │  │
         │  │  - /metrics             │  │
         │  └─────────────────────────┘  │
         └───────────────┬───────────────┘
                         │
                         ▼
         ┌───────────────────────────────┐
         │      Next.js Frontend         │
         │  ┌─────────────────────────┐  │
         │  │  - React Dashboard      │  │
         │  │  - Topic Filters        │  │
         │  │  - Real-time Updates    │  │
         │  │  - Responsive Design    │  │
         │  └─────────────────────────┘  │
         └───────────────────────────────┘
```

## Components

### 1. Data Ingestion Layer (GitHub Actions)

**Purpose:** Automated news discovery and processing

**Components:**
- `cmd/worker-json/main.go` - Main worker process
- `pkg/services/worker.go` - Ingestion orchestration
- `pkg/crawler/` - HTTP fetching and parsing
- `pkg/tagging/auto_tagger.go` - AI-powered categorization

**Flow:**
1. Fetch from configured sources (RSS, APIs, web scraping)
2. Extract title, URL, content, published date
3. Generate content hash for deduplication
4. AI summarization (TLDR, why it matters)
5. Auto-tagging based on content analysis
6. Relevance scoring (hot score, credibility, novelty)
7. Store in PostgreSQL

**Scheduling:** Runs every 30 minutes via GitHub Actions cron

### 2. Data Storage Layer (Neon.tech)

**Purpose:** Persistent storage with high availability

**Database:** PostgreSQL 15+ with connection pooling

**Key Tables:**

```sql
sources       - News sources configuration
items         - Raw news articles
summaries     - AI-generated summaries and tags
scores        - Multi-factor relevance scores
signals       - Engagement metrics
fetch_runs    - Ingestion job history
```

**Indexes:**
- `items.published_at` - Time-based queries
- `items.content_hash` - Deduplication
- `summaries.tags` - GIN index for tag filtering
- `scores.final` - Ranking queries

**Connection Pooling:**
- Max connections: 3 (cold start optimized)
- Idle timeout: 30 seconds
- Max lifetime: 5 minutes

### 3. API Layer (Vercel Serverless)

**Purpose:** RESTful API for frontend consumption

**Runtime:** Go 1.24.1 on Vercel serverless

**Endpoints:**

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/news` | GET | Get latest news (with topic filter) |
| `/api/trending` | GET | Get trending articles |
| `/api/search` | GET | Search articles by keyword |
| `/metrics` | GET | System metrics |
| `/healthz` | GET | Health check |

**Features:**
- Auto-retry for cold starts
- 15-second timeout
- CORS enabled
- JSON responses
- Error handling with proper HTTP codes

**Cold Start Optimization:**
- Minimal dependencies (database/sql + lib/pq)
- Connection pooling
- Database ping with timeout
- Frontend auto-retry mechanism

### 4. Presentation Layer (Next.js)

**Purpose:** User interface and experience

**Framework:** Next.js 14 with App Router

**Features:**
- Server-side rendering (SSR)
- Static asset optimization
- Responsive design (mobile-first)
- Real-time updates (30s polling)
- Topic filtering
- Loading states & error handling

**UI Components:**
- Dashboard metrics cards
- Topic filter bar
- News feed with cards
- Empty states
- Error states with retry
- Loading skeletons

## Data Flow

### News Ingestion Flow

```
1. GitHub Actions Trigger (cron: */30 * * * *)
   ↓
2. Worker Process Starts
   ↓
3. Fetch Enabled Sources (parallel)
   ↓
4. For Each Article:
   - Extract metadata
   - Generate content hash
   - Check if exists (dedup)
   - If new:
     * AI summarization
     * Auto-tagging
     * Score calculation
     * Insert to DB
   ↓
5. Update fetch_runs table
   ↓
6. Worker Process Ends
```

### User Request Flow

```
1. User Opens Dashboard
   ↓
2. Frontend Fetches /api/news
   ↓
3. Vercel Function Starts (cold start if needed)
   ↓
4. Connect to Neon.tech DB
   ↓
5. Execute SQL Query:
   SELECT items + summaries + scores
   WHERE published_at >= NOW() - 7 days
   ORDER BY score DESC
   LIMIT 20
   ↓
6. Return JSON Response
   ↓
7. Frontend Renders News Cards
   ↓
8. Auto-refresh every 30s
```

## Scoring Algorithm

Multi-factor scoring system for content relevance:

```go
final_score = (
    hot_score * 0.4 +           // Recency and engagement
    relevance_score * 0.3 +     // AI/ML topic relevance
    credibility_score * 0.2 +   // Source reputation
    novelty_score * 0.1         // Uniqueness
)
```

**Hot Score:** Based on HackerNews algorithm
```
hot = (points - 1) / (age_hours + 2)^1.8
```

**Relevance Score:** AI-powered content analysis
- Keyword matching (LLM, transformer, neural network, etc.)
- Domain reputation (arxiv.org, openai.com, etc.)
- Content quality signals

**Credibility Score:** Source-based
- Academic sources: 1.0
- Company blogs: 0.8
- Tech news: 0.7
- General news: 0.5

**Novelty Score:** Content uniqueness
- First mention: 1.0
- Similar content exists: 0.5-0.8
- Duplicate: 0.0

## Auto-Tagging System

AI-powered categorization into topics:

**Tags:**
- `llm` - Large Language Models
- `vision` - Computer Vision
- `data` - Data Science
- `security` - AI Security
- `rl` - Reinforcement Learning
- `robotics` - Robotics
- `ide` - Developer Tools
- `free-credits` - Student Programs
- `tools` - General Tools
- `research` - Research Papers
- `general_ai` - General AI News

**Algorithm:**
1. Extract keywords from title + content
2. Match against tag patterns
3. Apply domain-based rules
4. Assign multiple tags if applicable
5. Default to `general_ai` if no match

## Scalability Considerations

### Current Scale
- 40+ news sources
- ~100 articles/day
- ~3,000 articles/month
- <100 concurrent users

### Bottlenecks & Solutions

**Database Connections:**
- Current: 3 max connections (Vercel limit)
- Solution: Connection pooling, query optimization

**Cold Starts:**
- Current: 2-5 seconds first request
- Solution: Auto-retry, keep-alive pings, minimal dependencies

**API Rate Limits:**
- Current: Vercel free tier limits
- Solution: Caching, CDN, upgrade to Pro

### Future Scaling

**For 1,000+ concurrent users:**
1. Add Redis caching layer
2. Implement CDN for static assets
3. Database read replicas
4. Horizontal scaling with load balancer

**For 10,000+ articles/day:**
1. Message queue (RabbitMQ/SQS)
2. Distributed workers
3. Elasticsearch for search
4. Time-series database for metrics

## Security

### Authentication
- Currently: Public read-only API
- Future: API keys for write operations

### Data Protection
- HTTPS only
- Environment variables for secrets
- No PII collection
- CORS restrictions

### Database Security
- SSL/TLS connections
- Connection pooling
- Prepared statements (SQL injection prevention)
- Row-level security (RLS) ready

## Monitoring & Observability

### Metrics
- Request count
- Response time (p50, p95, p99)
- Error rate
- Database query time
- Scraper success rate

### Logging
- Vercel function logs
- GitHub Actions logs
- Database slow query logs

### Alerts
- Scraper failures
- API errors >5%
- Database connection issues
- High response times >5s

## Disaster Recovery

### Backup Strategy
- Neon.tech automatic backups (daily)
- Point-in-time recovery (7 days)
- GitHub repository backup

### Recovery Procedures
1. Database restore from Neon.tech backup
2. Redeploy from GitHub main branch
3. Re-run scraper to catch up on missed articles

### RTO/RPO
- Recovery Time Objective: <1 hour
- Recovery Point Objective: <24 hours

## Technology Decisions

### Why Go?
- Fast compilation
- Excellent concurrency
- Small binary size (Vercel compatible)
- Strong standard library

### Why Next.js?
- Server-side rendering
- Excellent developer experience
- Vercel integration
- React ecosystem

### Why Neon.tech?
- Serverless PostgreSQL
- Auto-scaling
- Generous free tier
- Excellent performance

### Why Vercel?
- Zero-config deployment
- Edge network
- Serverless functions
- GitHub integration

## References

- [Next.js Documentation](https://nextjs.org/docs)
- [Vercel Go Runtime](https://vercel.com/docs/functions/runtimes/go)
- [Neon.tech Docs](https://neon.tech/docs)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
