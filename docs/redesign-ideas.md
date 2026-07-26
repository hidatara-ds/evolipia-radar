# Evolipia Radar v2 Architecture Audit & Intelligence Platform Migration Plan

> Scope: inspection-only architecture review. This document does not implement code, create migrations, modify schemas, change API endpoints, or refactor existing behavior.

## Executive Summary

Evolipia Radar is currently implemented as an article-first AI news intelligence system. The active persistence model centers on `items`, with `sources`, `signals`, `scores`, `summaries`, `clusters`, `cluster_sources`, `fetch_runs`, `global_metrics`, and `settings` supporting ingestion, scoring, summarization, clustering, metrics, and configuration.

The codebase already contains early concepts named `signals`, but the current `signals` table is not a universal intelligence unit. It is a social-metric time series attached to an `items` row. Therefore, the platform cannot yet treat a GitHub repository, Hugging Face model, arXiv paper, Reddit thread, official announcement, documentation update, and news article as first-class peers.

Recommended migration direction: evolve the schema incrementally by introducing a new canonical `signals_v2` or `signals` domain model behind compatibility views/adapters, while keeping the existing `items` table operational during migration. A big-bang replacement is not recommended because the frontend, API filters, worker, scoring, summaries, and legacy feed endpoints are tightly coupled to the article-shaped `items` model.

## Evidence Map: Files Inspected

The audit is based on these real code areas:

- Database migrations: `migrations/000001_init_schema.up.sql`, `000003_add_clustering.up.sql`, `000004_add_metrics.up.sql`, `000005_add_settings_and_metrics.up.sql`, `000006_add_pgvector.up.sql`, `000007_add_llm_scores.up.sql`, and `000008_add_crawl_fields.up.sql`.
- Database models and repositories: `internal/models/*.go`, `internal/db/repositories.go`, and `internal/db/db.go`.
- Ingestion and crawling: `internal/services/worker.go`, `internal/connectors/*.go`, `internal/crawler/*.go`, and `cmd/worker/main.go`.
- AI/scoring/summarization/search/clustering: `internal/scoring/*.go`, `internal/summarizer/*.go`, `internal/ai/*.go`, `internal/search/vector.go`, and `internal/cluster/*.go`.
- Backend APIs: `cmd/server/main.go`, `cmd/api/main.go`, `internal/api/*.go`, and `internal/http/handlers/handlers.go`.
- Frontend: `app/page.tsx`, `src/api/client.ts`, `src/hooks/*.ts`, and `src/components/*.tsx`.

---

# Phase 1 — Current Architecture Audit

## 1. Database Documentation

### 1.1 Migration Inventory

| Migration | Purpose | Notes |
|---|---|---|
| `000001_init_schema` | Creates `sources`, `items`, `signals`, `scores`, `summaries`, `fetch_runs` | Core article-centric schema. |
| `000003_add_clustering` | Creates `clusters`, `cluster_sources` | `embedding` stored as `TEXT` fallback because pgvector was commented out in this migration. |
| `000004_add_metrics` | Creates singleton-like `global_metrics` | Uses `id SERIAL PRIMARY KEY` and inserts row `id = 1`. |
| `000005_add_settings_and_metrics` | Adds `settings`; extends `global_metrics` | Adds cluster metrics and top cluster titles JSONB. |
| `000006_add_pgvector` | Adds pgvector extension and `items.embedding vector(1536)` | Reintroduces pgvector requirement after migration 000003 commented it out. |
| `000007_add_llm_scores` | Adds `impact`, `engineering_value`, `reasoning` to `scores` | Repository upsert/select currently does not fully persist/read these added fields. |
| `000008_add_crawl_fields` | Adds crawl status/error/relevance/validation fields to `items` | Frontend `/api/items` depends on these fields. |

### 1.2 Current Tables

#### `sources`

Purpose: source/feed configuration.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `name TEXT NOT NULL`
- `type TEXT NOT NULL`
- `category TEXT NOT NULL`
- `url TEXT NOT NULL`
- `mapping_json JSONB NULL`
- `enabled BOOLEAN NOT NULL DEFAULT FALSE`
- `status TEXT NOT NULL DEFAULT 'pending'`
- `last_test_status TEXT NULL`
- `last_test_message TEXT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Indexes and constraints:

- Primary key: `id`
- Unique index: `idx_sources_unique_url ON sources(url)`
- Index: `idx_sources_enabled ON sources(enabled)`

Relationships:

- One source has many `items`.
- One source has many `fetch_runs`.

Limitations:

- `type`, `category`, and `status` are free-form `TEXT` values, not enums or lookup tables.
- `mapping_json` supports JSON API flexibility, but the source taxonomy is not normalized.
- A single `category` column is too narrow for future multi-dimensional taxonomy such as source format, trust tier, ecosystem area, and content type.

#### `items`

Purpose: canonical article/content row for current product.

Columns from migrations and later extensions:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE`
- `title TEXT NOT NULL`
- `url TEXT NOT NULL`
- `published_at TIMESTAMPTZ NOT NULL`
- `content_hash TEXT NOT NULL`
- `domain TEXT NOT NULL`
- `category TEXT NOT NULL`
- `raw_excerpt TEXT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `embedding vector(1536) NULL`
- `embedding_model TEXT NULL`
- `crawl_status TEXT NOT NULL DEFAULT 'done'`
- `crawl_error TEXT NULL`
- `relevance_score INT NOT NULL DEFAULT 0`
- `validated_at TIMESTAMPTZ NULL`

Indexes and constraints:

- Primary key: `id`
- Unique index: `idx_items_dedup ON items(content_hash)`
- Index: `idx_items_published_at ON items(published_at DESC)`
- Index: `idx_items_domain ON items(domain)`
- Index: `idx_items_source_id ON items(source_id)`
- HNSW vector index: `items_embedding_hnsw_idx ON items USING hnsw (embedding vector_cosine_ops)`
- Index: `idx_items_crawl_status ON items(crawl_status)`
- Index: `idx_items_relevance_score ON items(relevance_score DESC)`

Relationships:

- Many items belong to one source.
- One item has zero or many social metric rows in `signals`.
- One item has one score row in `scores`.
- One item has one summary row in `summaries`.
- One item can belong to many clusters through `cluster_sources`.

Limitations:

- The table name and columns assume an article-like resource: `title`, `url`, `published_at`, `domain`, `raw_excerpt`.
- `content_hash` is unique across all items, so cross-source duplicates collapse into one row and lose per-source confirmation unless represented elsewhere.
- `category` duplicates source/category semantics and is not normalized.
- `relevance_score` duplicates score semantics already present in `scores.relevance`.
- Embedding on `items` makes semantic search article-centric.

#### `signals` — current meaning

Purpose: time-series social/engagement metrics for an item.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE`
- `points INT NULL`
- `comments INT NULL`
- `rank_pos INT NULL`
- `fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Indexes:

- `idx_signals_item_fetched ON signals(item_id, fetched_at DESC)`

Important naming issue:

- This table is named `signals`, but it is not the future Signal concept requested for v2. It stores metrics for an `item` rather than being the primary intelligence event/entity.

#### `scores`

Purpose: article-level score dimensions.

Columns:

- `item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE`
- `hot DOUBLE PRECISION NOT NULL`
- `relevance DOUBLE PRECISION NOT NULL`
- `credibility DOUBLE PRECISION NOT NULL`
- `novelty DOUBLE PRECISION NOT NULL`
- `final DOUBLE PRECISION NOT NULL`
- `computed_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `impact DOUBLE PRECISION NOT NULL DEFAULT 0.0`
- `engineering_value DOUBLE PRECISION NOT NULL DEFAULT 0.0`
- `reasoning TEXT NOT NULL DEFAULT ''`

Indexes:

- `idx_scores_final ON scores(final DESC)`

Limitations:

- Primary key is `item_id`, preventing historical score versions.
- `impact`, `engineering_value`, and `reasoning` exist in schema/model, but the repository upsert only inserts `hot`, `relevance`, `credibility`, `novelty`, and `final`, which means LLM score extensions can remain at defaults unless written elsewhere.
- Score dimensions are article-centric and do not represent confidence, momentum, consensus, or phase.

#### `summaries`

Purpose: one summary per item.

Columns:

- `item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE`
- `tldr TEXT NOT NULL`
- `why_it_matters TEXT NOT NULL`
- `tags JSONB NOT NULL`
- `method TEXT NOT NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`

Indexes:

- `idx_summaries_tags_gin ON summaries USING GIN (tags)`

Limitations:

- `tags` is JSONB, flexible but not normalized; topic analytics require JSON containment rather than joins.
- One summary per item prevents multiple model versions, languages, or prompt versions.

#### `fetch_runs`

Purpose: per-source ingestion run history.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE`
- `fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- `status TEXT NOT NULL`
- `error TEXT NULL`
- `items_fetched INT NOT NULL DEFAULT 0`
- `items_inserted INT NOT NULL DEFAULT 0`

Indexes:

- `idx_fetch_runs_source_time ON fetch_runs(source_id, fetched_at DESC)`

Limitations:

- Tracks source-level success/failure and aggregate counts only.
- Does not track per-stage latency, parse failures by reason, source rate limits, or connector version.

#### `clusters`

Purpose: AI-generated clusters of articles.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `title VARCHAR(255) NOT NULL`
- `summary TEXT NOT NULL`
- `embedding TEXT`
- `score FLOAT8 DEFAULT 0.0`
- `created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP`
- `updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP`

Indexes:

- No active vector index; pgvector index is commented out in migration 000003.

Limitations:

- `embedding` is text in this table, while `items.embedding` later uses pgvector. This is inconsistent.
- Cluster membership maps to `items`, not future universal signals.
- No explicit topic, entities, source count, cluster phase, or trend history.

#### `cluster_sources`

Purpose: many-to-many mapping between clusters and items.

Columns:

- `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`
- `cluster_id UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE`
- `article_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE`
- `assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP`

Indexes and constraints:

- Unique constraint: `UNIQUE(cluster_id, article_id)`
- Index: `idx_cluster_sources_cluster_id ON cluster_sources(cluster_id)`
- Index: `idx_cluster_sources_article_id ON cluster_sources(article_id)`

Limitations:

- Column name `article_id` reinforces article-only mental model.
- No confidence or clustering method metadata on assignment.

#### `global_metrics`

Purpose: singleton metrics row for UI/system overview.

Columns:

- `id SERIAL PRIMARY KEY`
- `articles_processed INT DEFAULT 0`
- `filtered_articles INT DEFAULT 0`
- `api_hits INT DEFAULT 0`
- `updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP`
- `clusters_count INT DEFAULT 0`
- `avg_cluster_score DOUBLE PRECISION DEFAULT 0.0`
- `top_cluster_titles JSONB DEFAULT '[]'`
- `last_closeness_update TIMESTAMP WITH TIME ZONE`

Limitations:

- Metrics are globally aggregated and article-named.
- No dimensions by source category, signal type, connector, or date bucket.

#### `settings`

Purpose: key-value configuration.

Columns:

- `key TEXT PRIMARY KEY`
- `value TEXT`
- `updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP`

Limitations:

- No type, scope, encryption marker, or audit trail.

### 1.3 Current ERD — Text Format

```text
sources
  ├─< items
  │    ├─< signals              # current engagement snapshots, not v2 signals
  │    ├─1 scores
  │    ├─1 summaries
  │    └─< cluster_sources >─ clusters
  └─< fetch_runs

global_metrics       # singleton operational metrics
settings             # key/value config
```

### 1.4 Reusable Tables

- `sources`: reusable if expanded with flexible taxonomy, reliability, ownership, and source-level weighting.
- `fetch_runs`: reusable as ingestion run history, but should be extended or paired with run events/stage metrics.
- `clusters`: conceptually reusable, but should point to future signals instead of only article items.
- `settings`: reusable for low-risk feature flags and configuration.

### 1.5 Duplicated or Conflicting Data

- `items.category` and `sources.category` overlap.
- `items.relevance_score` and `scores.relevance` overlap but use different scales/types (`INT 0-100` vs `DOUBLE PRECISION`).
- `summaries.tags` duplicates `dto.ContentItem.Tags` during ingestion, but tags are not normalized into durable topic records.
- `clusters.embedding TEXT` conflicts with `items.embedding vector(1536)`.
- `signals` table name conflicts with the desired future meaning of Signal.

### 1.6 Scalability Issues

- `items.content_hash` as a global unique dedupe key prevents preserving all source occurrences of the same event.
- Search uses `ILIKE` across title/domain/excerpt in `/api/items`, without trigram/full-text indexes.
- Social metrics append to `signals` but there is no rollup table for momentum by time bucket.
- `scores` stores only current score state, making trend scoring and score regression analysis difficult.
- The worker computes scores after every processed source over all eligible recent items, which can become expensive as volume grows.
- Cron workers are process-local; multi-instance deployments need distributed locks.

---

## 2. Crawling and Ingestion Pipeline

There are two ingestion concepts in the codebase: the durable worker pipeline and an experimental multi-agent orchestrator.

### 2.1 Durable Worker Pipeline

Primary code path: `internal/services/worker.go`.

Flow:

```text
cron/manual command
  ↓
Worker.RunIngestion
  ↓
SourceRepository.GetEnabled
  ↓
processSource per source
  ↓
fetchItems by source.Type
  ↓
connectors return []dto.ContentItem
  ↓
NormalizeURL + ContentHash(title, normalizedURL)
  ↓
ItemRepository.GetByContentHash
  ↓
insert item if new
  ↓
GenerateExtractiveSummary + SummaryRepository.Upsert
  ↓
insert engagement snapshot into current signals table if available
  ↓
FetchRunRepository.Create
  ↓
computeScores over recent items needing scoring
```

### 2.2 Connectors and Source Handling

| Source type/function | Source | Parser | Cleaning/normalization | Storage impact |
|---|---|---|---|---|
| `FetchHackerNews` | Hacker News Firebase API | JSON top stories and item endpoint | URL parsed for domain; HN item without URL becomes discussion URL | `items`, `signals` with points/comments/rank-like field |
| `FetchRSSAtom` | Source-configured RSS/Atom feed | String-based RSS/Atom parsing | `fetchWithLimits`; URL later normalized by worker | `items`, summaries, scores |
| `FetchArxiv` | arXiv API | XML decoding | Domain forced to `arxiv.org`; categories copied to DTO tags | Stored as category `news` despite research nature |
| `FetchHuggingFaceTrending` | Hugging Face model API | JSON | Builds model URL/title; likes/downloads mapped to points/comments | Stored as item category `models` |
| `FetchPapersWithCode` | Papers With Code API | JSON | Chooses paper URL, parses date, normalizes domain | Stored as research-shaped item |
| `FetchLMSYSArena` | Chatbot Arena page | Regex over HTML | Produces model ranking-like items or placeholder | Stored as benchmark item |
| `FetchOpenAIStatus` | OpenAI status RSS | RSS reuse | Adds OpenAI/status/API tags and domain | Stored as status item |
| `FetchAnthropicDocs` | Anthropic release notes page | Regex over HTML | Looks for version/date patterns; returns one synthetic item | Stored as docs item |
| `FetchGitHubTrending` | GitHub trending page | Regex over HTML | Extracts owner/repo and stars | Stored as tools item |
| `FetchJSONAPI` | Configured JSON endpoint | Generic mapping_json | Extracts mapped values and dates | Stored as item |

### 2.3 Cleaning and Safety

- `fetchWithLimits` validates outbound URL, requires absolute URLs, blocks credentials, blocks localhost/private IPs, disables redirects, applies HTTP timeout, and limits response bytes.
- Important issue: `validateOutboundURL` only allows `https`, but arXiv constants and seed scripts use `http://export.arxiv.org/...`. This can cause the arXiv connector or sources to fail unless changed upstream or allowed differently.
- RSS/Atom parsing is string-based rather than XML-parser-based for RSS, which is fragile for namespaces, CDATA, escaped content, and feed variants.
- Several HTML parsers use regex placeholders, which is acceptable for prototype signals but not robust for production intelligence.

### 2.4 Deduplication

- Persistent dedupe: `ContentHash(title, normalizedURL)` into `items.content_hash` with unique index.
- Runtime orchestrator dedupe: `CrawlBudget.seenCache` by URL, in memory for 24 hours.
- Limitation: dedupe occurs at item level and loses cross-source evidence. For intelligence, duplicates should be preserved as observations attached to a canonical signal/event.

### 2.5 Experimental Multi-Agent Orchestrator

Primary files: `internal/crawler/orchestrator.go`, `rss_agent.go`, `trending_agent.go`, `agent.go`, `budget.go`.

Flow:

```text
Orchestrator.Start ticker
  ↓
RunCycle
  ↓
RSSAgent + TrendingAgent return []Article
  ↓
CrawlBudget dedupe/rate limit by URL
  ↓
optional in-memory cluster service
  ↓
AI cluster service receives generated UUID article ID
  ↓
optional embedding generation attempted for generated UUID
```

Critical issue:

- The orchestrator generates a fresh UUID for each discovered article instead of inserting/looking up an `items` row first. Downstream cluster mappings reference `items(id)`, so cluster inserts can fail or become logically disconnected unless another path creates the item with the same UUID.

---

## 3. AI Pipeline Audit

### 3.1 Summarization

Current implementation:

```text
models.Item
  ↓
GenerateExtractiveSummary(item)
  ↓
models.Summary{TLDR, WhyItMatters, Tags, Method}
  ↓
summaries table
  ↓
Feed/detail API and UI display
```

- Extractive summary uses `raw_excerpt` when present, otherwise title.
- Tags are keyword-derived from title/excerpt/category.
- LLM summarization exists via `GenerateLLMSummary`, `internal/llm/client.go`, and `internal/ai.Service`, but the durable worker path uses extractive summarization by default.

Database impact:

- One row in `summaries` per item.

Frontend usage:

- The current Next.js dashboard mainly uses `raw_excerpt`, score fields, category, and source fields from `/api/items`.
- Legacy detail/feed handlers return `summary.tldr`, `why_it_matters`, and tags.

### 3.2 Embeddings and Search

Current implementation:

```text
text/title/content
  ↓
OpenRouter embedding request
  ↓
items.embedding vector(1536), items.embedding_model
  ↓
SearchByEmbedding / HybridSearcher
  ↓
search API responses
```

Issues:

- Embeddings are attached to `items`, so semantic search is article-centric.
- `clusters.embedding` is `TEXT`, while `items.embedding` is pgvector.
- OpenRouter embedding code defaults to `text-embedding-3-small`; provider compatibility should be verified before production reliance.

### 3.3 Relevance Scoring

There are two relevance concepts:

1. Validator relevance score in `internal/crawler/validator.go`: integer 0-100 based on topic keyword matches in title/excerpt.
2. Algorithmic score in `internal/scoring/scoring.go`: floating-point score used in `scores.relevance` and final weighted score.

Current scoring flow:

```text
models.Item + latest models.Signal + models.Summary
  ↓
ComputeScore
  ↓
hot, relevance, credibility, novelty, final
  ↓
scores table
  ↓
feed ranking / sorting / UI score display
```

Score dimensions:

- `hot`: points/comments/rank plus time decay.
- `relevance`: keyword matching over title, excerpt, summary, tags, category.
- `credibility`: domain list/config based.
- `novelty`: age decay based on publication time.
- `final`: weighted combination.

Limitations:

- No cross-source confirmation.
- No source-category-specific logic.
- No entity importance.
- No topic velocity or time-series trend phase.
- No confidence score.
- No score versioning/history.

### 3.4 Trending Score

Current trending is mostly item-level:

- `GetRising` selects items with recent rows in current `signals` table.
- `computeHotScore` uses engagement snapshot values and publication age.
- Cluster service increments cluster score when articles are associated.

Limitations:

- Momentum is not normalized by source type.
- There is no time-bucket trend history by topic/entity/source category.
- No consensus metric across independent sources.

### 3.5 Entity Extraction

No durable entity extraction pipeline was found. Entities such as companies, products, models, repos, authors, organizations, and papers are not normalized into tables.

Existing closest approximations:

- Tags in `summaries.tags`.
- Regex extraction in specific connectors, such as model names in LMSYS or repo names in GitHub Trending.
- Topic keyword matching in scoring/summarizer.

### 3.6 Topic Extraction

Current topic extraction is keyword-based:

- `summarizer.extractTags` maps title/excerpt/category text to tags.
- `scoring.computeRelevanceScoreWithConfig` uses configured topic keyword groups.
- Tags live as JSONB in `summaries`, not normalized topic rows.

---

## 4. Backend Audit

### 4.1 Main Runtime Entrypoints

| Entrypoint | Responsibility |
|---|---|
| `cmd/server/main.go` | Gin server for `/api/items`, crawl status/progress, manual crawl, static web assets, scheduler startup/shutdown. |
| `cmd/worker/main.go` | Cron worker that periodically runs ingestion. |
| `cmd/api/main.go` | Alternate API server with static routes, legacy `/api/news`, `/api/trending`, `/api/search`, v1/v2 route groups. |
| `cmd/worker-json/main.go` | Exports latest news-like data to JSON from database or fallback mode. |

### 4.2 Services

| Service/module | Responsibility |
|---|---|
| `Worker` | Enabled source ingestion, connector dispatch, item insert/dedupe, summary creation, signal snapshot insert, score computation. |
| `FeedService` | Builds feed/rising/detail response payloads from repositories. |
| `SourceService` | Source CRUD/test/enable workflow for legacy handlers. |
| `ai.Service` | Thin wrapper around LLM provider for chat/summarize/embed. |
| `HybridSearcher` | Text + semantic search and reciprocal rank fusion. |
| `ClusterService` | Persistent AI cluster creation/assignment/merge logic. |
| `Scheduler` | Cron-based crawl scheduling with process-local concurrency guard and SSE progress callbacks. |

### 4.3 API Endpoints Found

Primary server (`cmd/server/main.go`):

- `GET /healthz`
- `GET /api/items`
- `GET /api/crawl/progress`
- `GET /api/crawl/status`
- `POST /api/crawl`
- static `/` and `/web`

Alternate/legacy server (`cmd/api/main.go` and `internal/http/handlers`):

- `GET /healthz`
- `GET /metrics`
- `GET /api/news`
- `GET /api/trending`
- `GET /api/search`
- v1 routes for feed/rising/item/search/sources/settings depending on registration.
- v2 AI routes: `POST /v2/chat`, `POST /v2/summarize`.

### 4.4 Cron Jobs and Scheduled Work

- `Scheduler` in `internal/crawler/scheduler.go` uses `robfig/cron/v3`, supports cron descriptors and interval expressions, and blocks overlapping crawls in-process.
- `cmd/worker/main.go` creates a separate cron worker using `cfg.WorkerCron`.
- Multi-instance deployments can run duplicate jobs because locking is in memory, not database-backed.

### 4.5 Queues

No durable queue system was found. Work is synchronous inside worker loops or cron callbacks.

---

## 5. Frontend Audit

### 5.1 Article Schema Dependencies

The active Next.js dashboard depends on `NewsItem` from `src/api/client.ts`. It expects fields such as:

- `id`
- `source_id`
- `source_name`
- `title`
- `url`
- `published_at`
- `domain`
- `category`
- `raw_excerpt`
- `crawl_status`
- `crawl_error`
- `relevance_score`
- `validated_at`
- `created_at`
- optional nested `score`

### 5.2 Score Dependencies

- `app/page.tsx` renders metric cards and feed cards using `relevance_score`, `score.final`, `score.credibility`, `score.impact`, and visual heat/scale behavior.
- `FilterBar` uses `min_relevance`, status, source, category, and sort fields.
- `/api/items` supports sort modes `date`, `relevance`, `credibility`, `impact`, and `oldest`.

### 5.3 Summary and Tags Dependencies

- The active `/api/items` dashboard uses excerpts more directly than summary rows.
- Legacy handlers expose `summary.tldr`, `why_it_matters`, and `tags`, so any detail/feed page built on legacy v1 routes depends on summary structure.

### 5.4 Components That Would Break on Schema Changes

If `items` is removed or renamed without compatibility:

- `src/api/client.ts` breaks because `fetchItems` calls `/api/items` and deserializes article-shaped `NewsItem`.
- `app/page.tsx` breaks because it renders item title, URL, published date, domain, category, excerpt, crawl status, and score fields.
- `FilterBar` and `useFilters` break if query parameters or sortable/filterable field names change.
- `CrawlProgress` and `useCrawlProgress` are less schema-coupled, but rely on crawl endpoint payload shape.
- `DataFreshness` relies on crawl status timestamps rather than item schema and is relatively safe.

---

## 6. Current Limitations and Bottlenecks

### 6.1 Article-Centric Core

Evidence:

- `items` is the central table for storage, scoring, summaries, embeddings, API output, and frontend display.
- `scores`, `summaries`, current `signals`, and `cluster_sources` all reference `items(id)`.

Impact:

- Non-article sources must be disguised as items with title/URL/date/domain/category.
- A GitHub repo, model leaderboard movement, status incident, or documentation update cannot preserve source-specific metadata cleanly.

### 6.2 Naming Conflict Around `signals`

Evidence:

- Existing `signals` table means engagement snapshots for items: points/comments/rank/fetched_at.

Impact:

- The desired v2 Signal concept cannot be introduced as `signals` without migration conflict or semantic confusion.

### 6.3 Cross-Source Confirmation Is Lost

Evidence:

- Unique `items.content_hash` dedupes by title and normalized URL.
- There is no observation table representing multiple sources reporting the same event.

Impact:

- Consensus scoring cannot be calculated reliably.
- Duplicate articles can be filtered out before they become evidence.

### 6.4 Duplicated Scoring Concepts

Evidence:

- `items.relevance_score INT` and `scores.relevance DOUBLE PRECISION` both represent relevance.
- Scoring, validation, summaries, and API sorting use related but distinct concepts.

Impact:

- Users can see inconsistent relevance values.
- Future scoring extensions become difficult to reason about.

### 6.5 Source Taxonomy Is Too Flat

Evidence:

- `sources.type` and `sources.category` are free-form text.
- DTO category may also be set by connector and then overwritten by source category in worker item creation.

Impact:

- Adding official/blog/research/social/developer/source-type-specific scoring requires brittle switch statements or text conventions.

### 6.6 Fragile Parser Strategy

Evidence:

- RSS parsing uses string splitting.
- LMSYS, GitHub, and Anthropic docs scraping use regex heuristics or placeholder items.

Impact:

- Production reliability is limited for dynamic or frequently changing pages.

### 6.7 Inconsistent Embedding Architecture

Evidence:

- `items.embedding` is pgvector.
- `clusters.embedding` is text fallback.

Impact:

- Cluster similarity and item similarity cannot share consistent vector operations.

### 6.8 Operational Scalability Concerns

Evidence:

- Cron concurrency guards and crawl budget are in memory.
- Token budget is in memory.
- No durable queue exists.

Impact:

- Multiple servers can duplicate crawls and AI requests.
- Restarting processes resets crawl dedupe, token accounting, and budgets.

---

# Phase 2 — Product Evolution: Article to Signal

## 7. Proposed Target Architecture

Future unit of information: `Signal`.

A signal is a normalized observation from any source category, such as:

- news article
- official announcement
- GitHub repository/release/issue
- Reddit discussion
- Hacker News story
- Hugging Face model/dataset/space
- arXiv paper
- Papers With Code paper/task/result
- company blog post
- release note
- documentation update
- benchmark leaderboard change

### 7.1 Target Data Flow

```text
Source Registry
  ↓
Connector Adapter
  ↓
Raw Observation Capture
  ↓
Signal Normalization
  ↓
Canonical Dedupe / Entity Resolution
  ↓
Entity Extraction
  ↓
Topic Extraction
  ↓
Signal Scoring
  ↓
Trend Detection and Rollups
  ↓
Insight Generation
  ↓
Frontend Intelligence Views
```

### 7.2 Proposed Backend Separation

Recommended domains:

- `source_registry`: source definitions, taxonomy, weights, health.
- `collectors`: connector-specific fetch and parsing.
- `normalization`: maps connector DTOs to canonical signals.
- `identity`: canonical URL/object identifiers and dedupe.
- `enrichment`: entities, topics, embeddings, summaries.
- `scoring`: signal scores, confidence, momentum, consensus.
- `trends`: daily/hourly rollups, cluster history, trend phases.
- `insights`: generated briefings and explanations.
- `compatibility`: adapters that expose old article-shaped `/api/items` until frontend migration completes.

---

# Phase 3 — Database Redesign Proposal

Do not apply this schema yet. This is a proposed v2 target after auditing current implementation.

## 8. Recommended Schema Direction

Recommendation: introduce a v2 signal model alongside the current schema first, then migrate reads/writes gradually.

### 8.1 Source and Taxonomy Tables

#### `source_categories`

```sql
source_categories (
  id UUID PK,
  key TEXT UNIQUE NOT NULL,
  label TEXT NOT NULL,
  description TEXT,
  parent_id UUID NULL FK source_categories(id),
  metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
```

Purpose:

- Allows expandable taxonomy without schema redesign.
- Supports top-level categories such as `official`, `news`, `community`, `developer`, `research`, `social`, `blog`, `documentation`, `video`, `podcast`, `benchmark`, and `model_registry`.

#### Extend or replace `sources`

Preferred: evolve existing `sources` with additive columns rather than replace immediately.

Add later, not now:

- `category_id UUID NULL REFERENCES source_categories(id)`
- `trust_tier TEXT`
- `default_weight DOUBLE PRECISION`
- `metadata JSONB NOT NULL DEFAULT '{}'`
- `health_status TEXT`
- `last_success_at TIMESTAMPTZ`
- `last_failure_at TIMESTAMPTZ`

Keep old `type` and `category` during compatibility period.

### 8.2 Canonical Signal Tables

Because `signals` already exists with a different meaning, use one of these strategies:

- Safer: create `intelligence_signals` or `signals_v2` first.
- Later: rename old `signals` to `item_signal_metrics` and reserve `signals` for canonical v2.

#### `signals_v2`

```sql
signals_v2 (
  id UUID PK DEFAULT gen_random_uuid(),
  source_id UUID NOT NULL REFERENCES sources(id),
  source_category_id UUID NULL REFERENCES source_categories(id),
  signal_type TEXT NOT NULL,
  canonical_key TEXT NOT NULL,
  title TEXT NOT NULL,
  url TEXT,
  summary TEXT,
  content_excerpt TEXT,
  raw_payload JSONB NOT NULL DEFAULT '{}',
  source_published_at TIMESTAMPTZ,
  observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  normalized_at TIMESTAMPTZ,
  language TEXT,
  domain TEXT,
  status TEXT NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(source_id, canonical_key)
)
```

Indexes:

- `(observed_at DESC)`
- `(source_published_at DESC)`
- `(source_id, observed_at DESC)`
- `(signal_type, observed_at DESC)`
- GIN on `raw_payload`
- optional trigram/full-text index on `title` and `content_excerpt`

Why needed:

- Allows each connector to preserve source-native metadata in `raw_payload`.
- `signal_type` distinguishes article, repo, paper, release, docs update, discussion, status incident, model, benchmark, etc.
- `canonical_key` dedupes within source while allowing cross-source observations to remain available.

#### `signal_observations`

```sql
signal_observations (
  id UUID PK DEFAULT gen_random_uuid(),
  signal_id UUID NOT NULL REFERENCES signals_v2(id) ON DELETE CASCADE,
  source_id UUID NOT NULL REFERENCES sources(id),
  observation_key TEXT NOT NULL,
  url TEXT,
  title TEXT,
  raw_payload JSONB NOT NULL DEFAULT '{}',
  observed_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(source_id, observation_key)
)
```

Purpose:

- Preserves cross-source evidence and duplicate/near-duplicate observations.
- Enables consensus scoring.

### 8.3 Entities and Topics

#### `entities`

```sql
entities (
  id UUID PK DEFAULT gen_random_uuid(),
  entity_type TEXT NOT NULL,
  name TEXT NOT NULL,
  canonical_name TEXT NOT NULL,
  aliases JSONB NOT NULL DEFAULT '[]',
  external_ids JSONB NOT NULL DEFAULT '{}',
  importance_score DOUBLE PRECISION DEFAULT 0,
  metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE(entity_type, canonical_name)
)
```

Entity types:

- company, product, model, dataset, repository, paper, author, organization, benchmark, framework, regulation, event.

#### `signal_entities`

```sql
signal_entities (
  signal_id UUID REFERENCES signals_v2(id) ON DELETE CASCADE,
  entity_id UUID REFERENCES entities(id) ON DELETE CASCADE,
  role TEXT,
  confidence DOUBLE PRECISION NOT NULL DEFAULT 0,
  extraction_method TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (signal_id, entity_id, role)
)
```

#### `topics`

```sql
topics (
  id UUID PK DEFAULT gen_random_uuid(),
  key TEXT UNIQUE NOT NULL,
  label TEXT NOT NULL,
  parent_id UUID NULL REFERENCES topics(id),
  keywords JSONB NOT NULL DEFAULT '[]',
  metadata JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
```

#### `signal_topics`

```sql
signal_topics (
  signal_id UUID REFERENCES signals_v2(id) ON DELETE CASCADE,
  topic_id UUID REFERENCES topics(id) ON DELETE CASCADE,
  confidence DOUBLE PRECISION NOT NULL DEFAULT 0,
  extraction_method TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (signal_id, topic_id)
)
```

### 8.4 Scoring and Trend Tables

#### `signal_scores`

```sql
signal_scores (
  signal_id UUID PRIMARY KEY REFERENCES signals_v2(id) ON DELETE CASCADE,
  signal_score DOUBLE PRECISION NOT NULL,
  confidence_score DOUBLE PRECISION NOT NULL,
  momentum DOUBLE PRECISION NOT NULL,
  growth_rate DOUBLE PRECISION NOT NULL,
  consensus_score DOUBLE PRECISION NOT NULL,
  credibility_score DOUBLE PRECISION NOT NULL,
  entity_score DOUBLE PRECISION NOT NULL,
  topic_score DOUBLE PRECISION NOT NULL,
  developer_activity_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  community_engagement_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  research_activity_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  officiality_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  trend_phase TEXT NOT NULL,
  reasoning TEXT NOT NULL DEFAULT '',
  model_version TEXT NOT NULL,
  computed_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
```

#### `signal_score_history`

```sql
signal_score_history (
  id UUID PK DEFAULT gen_random_uuid(),
  signal_id UUID NOT NULL REFERENCES signals_v2(id) ON DELETE CASCADE,
  score_payload JSONB NOT NULL,
  model_version TEXT NOT NULL,
  computed_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
```

#### `trend_history`

```sql
trend_history (
  id UUID PK DEFAULT gen_random_uuid(),
  trend_key TEXT NOT NULL,
  trend_type TEXT NOT NULL,
  bucket_start TIMESTAMPTZ NOT NULL,
  bucket_granularity TEXT NOT NULL,
  signal_count INT NOT NULL DEFAULT 0,
  source_count INT NOT NULL DEFAULT 0,
  weighted_score DOUBLE PRECISION NOT NULL DEFAULT 0,
  momentum DOUBLE PRECISION NOT NULL DEFAULT 0,
  growth_rate DOUBLE PRECISION NOT NULL DEFAULT 0,
  phase TEXT NOT NULL,
  metadata JSONB NOT NULL DEFAULT '{}',
  UNIQUE(trend_key, trend_type, bucket_start, bucket_granularity)
)
```

#### `topic_clusters`

```sql
topic_clusters (
  id UUID PK DEFAULT gen_random_uuid(),
  topic_id UUID NULL REFERENCES topics(id),
  title TEXT NOT NULL,
  summary TEXT NOT NULL,
  embedding vector(1536),
  score DOUBLE PRECISION NOT NULL DEFAULT 0,
  phase TEXT NOT NULL DEFAULT 'emerging',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
)
```

#### `cluster_signals`

```sql
cluster_signals (
  cluster_id UUID REFERENCES topic_clusters(id) ON DELETE CASCADE,
  signal_id UUID REFERENCES signals_v2(id) ON DELETE CASCADE,
  assignment_confidence DOUBLE PRECISION NOT NULL DEFAULT 0,
  assignment_method TEXT NOT NULL,
  assigned_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (cluster_id, signal_id)
)
```

### 8.5 Compatibility Tables/Views

During migration, keep the old tables and add adapters:

- `item_signal_map(item_id UUID, signal_id UUID, created_at TIMESTAMPTZ)`
- Optional view `items_compat` generated from `signals_v2` for legacy frontend/API reads.

---

# Phase 4 — Migration Strategy

## 9. Can the Current Schema Evolve?

Yes, but not by renaming `items` directly at the start. The current schema should evolve additively because:

- API and frontend depend on `/api/items` and `NewsItem` fields.
- Scoring, summaries, embeddings, and clusters reference `items`.
- Existing live data needs continuity.

Recommended approach: additive migration with compatibility, then gradual cutover.

## 10. Migration Order

### Step 0 — Freeze Contract and Add Observability

- Document current `/api/items`, legacy feed/search, and worker contracts.
- Add migration test fixtures in a staging environment.
- Do not change production reads yet.

### Step 1 — Introduce Taxonomy Tables

- Add `source_categories` and `topics`.
- Backfill categories from existing `sources.category`, `sources.type`, and `items.category`.
- Keep old columns.

### Step 2 — Introduce `signals_v2` and `signal_observations`

- Backfill one `signals_v2` row per existing `items` row.
- Create `item_signal_map` to preserve old item ID to new signal ID.
- Convert current `signals` engagement rows into signal metrics/history linked through the map.

### Step 3 — Introduce Entities and Normalized Topics

- Run deterministic keyword/topic extraction first.
- Add LLM/entity extraction later behind feature flags.
- Backfill `signal_topics` from `summaries.tags`.

### Step 4 — Add Signal-Based Scores

- Compute `signal_scores` in parallel with old `scores`.
- Compare score distributions before frontend exposure.
- Keep old `scores` powering old pages until parity is acceptable.

### Step 5 — Add Trends and Cluster Migration

- Create `topic_clusters` and `cluster_signals`.
- Backfill from existing `clusters` and `cluster_sources` through `item_signal_map`.
- Add trend rollups from historical `signals`/engagement rows and source observations.

### Step 6 — API Compatibility Layer

- Keep `/api/items` unchanged.
- Add future `/api/signals`, `/api/trends`, `/api/entities`, and `/api/topics` when ready.
- Optionally power `/api/items` from `items_compat` after validation.

### Step 7 — Frontend Migration

- Add signal-oriented views while preserving article card compatibility.
- Replace article wording with signal wording gradually.
- Add source-type-specific cards for repo, paper, discussion, official update, and article.

### Step 8 — Deprecation

- After production parity and data validation, mark old `items`, current `signals`, old `scores`, and old `summaries` as legacy.
- Only then consider renaming old `signals` to `item_signal_metrics` or archiving old tables.

## 11. Backward Compatibility

Must preserve:

- `/api/items` response shape.
- Existing `NewsItem` frontend model.
- Existing `scores`/`summaries` reads.
- Existing source management routes.

Recommended compatibility layer:

```text
signals_v2 + signal_scores + signal_topics
  ↓
items_compat view or adapter query
  ↓
/api/items unchanged
  ↓
current frontend unchanged during backend migration
```

## 12. Breaking Changes and Risks

| Risk | Impact | Mitigation |
|---|---|---|
| Renaming `signals` too early | Breaks current engagement scoring | Use `signals_v2` first. |
| Replacing `items` too early | Breaks frontend and repositories | Use `item_signal_map` and compatibility view. |
| Score distribution drift | Users lose trust in rankings | Run old/new scoring in parallel. |
| Backfill errors | Duplicate or missing signals | Use deterministic canonical keys and audit queries. |
| pgvector availability | Migrations fail in environments without extension | Keep vector feature optional or environment-gated. |
| Multi-source dedupe mismatch | Consensus inflated or undercounted | Separate canonical signal from observations. |
| LLM extraction cost | Cost spikes during backfill | Batch, cache, deterministic fallback, model/version tracking. |

## 13. Rollback Strategy

- Each additive migration should be independently reversible.
- Keep old writes to `items`, `scores`, `summaries`, and current `signals` until v2 is validated.
- Use feature flags for v2 writes, v2 reads, and frontend signal views.
- For backfills, write idempotent jobs keyed by deterministic IDs/canonical keys.
- Rollback by disabling v2 read flags and returning `/api/items` to old queries.

---

# Phase 5 — Signal-Based Scoring Redesign

## 14. Proposed Score Outputs

### 14.1 `signal_score`

Overall importance score.

Conceptual inputs:

- source credibility
- signal type weight
- recency
- entity importance
- topic importance
- engagement
- research/developer activity
- cross-source consensus
- novelty

### 14.2 `confidence_score`

How reliable the system is that the signal is real and correctly understood.

Conceptual inputs:

- source trust tier
- official source indicator
- duplicate confirmation across independent sources
- parser confidence
- entity/topic extraction confidence
- raw payload completeness

### 14.3 `momentum`

Current velocity of attention/activity.

Conceptual inputs:

- observation count change over recent buckets
- engagement acceleration
- GitHub stars/forks/issues over time
- Reddit/HN comments over time
- paper/code mentions over time

### 14.4 `growth_rate`

Relative increase compared with baseline.

Conceptual inputs:

- current bucket count vs trailing average
- source-category-weighted growth
- entity/topic baseline normalization

### 14.5 `trend_phase`

Lifecycle classification.

Suggested phases:

- `new`
- `emerging`
- `accelerating`
- `mainstream`
- `cooling`
- `stale`

### 14.6 `cross_source_consensus`

Whether independent sources support the same underlying signal.

Conceptual inputs:

- number of independent source categories
- source credibility mix
- semantic similarity among observations
- time proximity
- entity/topic overlap

### 14.7 `emerging_signal_detection`

A boolean or score highlighting early unusual movement.

Conceptual inputs:

- low historical baseline
- high growth rate
- high-quality early source
- entity/topic novelty
- small but diverse source confirmation

## 15. Suggested Formula Sketch

```text
signal_score =
  0.20 * credibility_score +
  0.15 * recency_score +
  0.15 * entity_score +
  0.15 * topic_score +
  0.15 * consensus_score +
  0.10 * engagement_score +
  0.10 * novelty_score
```

Category-specific modifiers:

- Official announcements get higher confidence and officiality weight.
- GitHub/repos get developer activity weight.
- Reddit/HN get community engagement and velocity weight.
- arXiv/Papers With Code get research activity and citation/implementation weight.
- Documentation/release notes get officiality and product-change weight.

---

# Phase 6 — Flexible Source Taxonomy

## 16. Taxonomy Requirements

The taxonomy should support future expansion without schema redesign.

Recommended model:

- A normalized `source_categories` table with hierarchical parent/child categories.
- A `metadata JSONB` field for category-specific parameters.
- Optional many-to-many source-category assignment later if a source fits multiple categories.

Initial top-level categories:

- `official`
- `news`
- `community`
- `developer`
- `research`
- `social`
- `blog`
- `documentation`
- `video`
- `podcast`
- `benchmark`
- `model_registry`
- `status`

Example subcategories:

- `developer.github_repo`
- `developer.release_note`
- `research.arxiv_paper`
- `research.paper_with_code`
- `official.company_blog`
- `official.status_page`
- `community.hacker_news`
- `community.reddit_thread`
- `model_registry.huggingface_model`

---

# Phase 7 — Frontend Impact and Product UI Direction

## 17. Impact by Future Page

### Home Page / Dashboard

Current state:

- Article feed with metrics, crawl progress, filter bar, and item cards.

Future change:

- Replace pure article list with signal stream.
- Add signal-type icons and card variants.
- Add intelligence briefing: what changed, why it matters, confidence, and consensus.

### Trending Page

Current state:

- Existing backend has rising/feed concepts, but Next.js app is mostly single dashboard.

Future change:

- Trend cards grouped by topic/entity.
- Show momentum, growth rate, source diversity, and phase.

### Analytics

Future change:

- Time-series charts for topic velocity, entity mentions, source category mix, and score distributions.

### Article Detail / Signal Detail

Future change:

- Article detail becomes Signal Detail.
- Show raw observations, related entities, topic assignment, score explanation, and timeline.

### Filtering

Current filters:

- Search, date, sources, categories, min relevance, status, sort.

Future filters:

- Signal type, source category, entity, topic, confidence, momentum, consensus, phase, official-only, developer-only, research-only.

### Search

Current search:

- `/api/items` text filters and optional legacy hybrid search.

Future search:

- Search across signals, observations, entities, topics, and clusters.
- Use semantic search over canonical signal text and entity/topic graph filters.

## 18. UI Recommendation: Intelligence Platform, Not News Aggregator

Use these primitives:

- **Signal cards** instead of article cards.
- **Trend clusters** instead of only chronological feed sections.
- **Briefing panel** instead of static hero.
- **Evidence drawer** showing source observations and cross-source confirmation.
- **Score explanation popover** showing credibility, confidence, momentum, and consensus drivers.
- **Entity pages** for companies, models, repos, frameworks, and papers.
- **Topic pages** for agents, inference, robotics, safety, open source, evals, infrastructure, and policy.

---

# Deliverables Summary

## 19. Current Architecture Review

Evolipia Radar is article-first. The durable ingestion worker converts all sources into `dto.ContentItem`, then into `items`, then into summaries, engagement snapshots, and scores. The frontend consumes `/api/items` as article-shaped `NewsItem` records.

## 20. Current Database Documentation

Covered above in section 1. Tables: `sources`, `items`, current `signals`, `scores`, `summaries`, `fetch_runs`, `clusters`, `cluster_sources`, `global_metrics`, and `settings`.

## 21. Current Data Flow Diagram

```text
sources
  ↓
Worker / connectors
  ↓
dto.ContentItem
  ↓
NormalizeURL + ContentHash
  ↓
items
  ├─ summaries
  ├─ current signals engagement snapshots
  ├─ scores
  ├─ embeddings
  └─ cluster_sources → clusters
  ↓
/api/items + legacy feed/search APIs
  ↓
Next.js dashboard
```

## 22. Weakness Analysis

Main weaknesses:

- Article-centric persistence.
- Current `signals` name conflict.
- No canonical cross-source observation layer.
- Duplicated relevance concepts.
- Flat source taxonomy.
- Fragile parsing for some sources.
- Inconsistent vector representation.
- In-memory scheduling/budget controls.
- No durable queues or trend rollups.

## 23. Proposed Architecture

Introduce signal-first architecture with separate source registry, connector adapters, raw observation capture, canonical signal normalization, entity/topic extraction, signal scoring, trend rollups, insight generation, and compatibility adapters.

## 24. Proposed Database Schema

Recommended new concepts:

- `source_categories`
- `signals_v2` or `intelligence_signals`
- `signal_observations`
- `entities`
- `signal_entities`
- `topics`
- `signal_topics`
- `signal_scores`
- `signal_score_history`
- `trend_history`
- `topic_clusters`
- `cluster_signals`
- `item_signal_map`

## 25. Migration Strategy

Use additive migration with old/new parallel operation:

1. Document contracts and add tests.
2. Add taxonomy tables.
3. Add v2 signals and observations.
4. Backfill item-to-signal mapping.
5. Add entities/topics.
6. Run signal scores in parallel.
7. Migrate clusters/trends.
8. Add signal APIs.
9. Migrate frontend.
10. Deprecate old article-only tables after parity.

## 26. Risk Analysis

Highest risks:

- Breaking `/api/items` and the current frontend.
- Misusing old `signals` for new signal semantics.
- Losing live data during dedupe/backfill.
- Score drift and user trust loss.
- AI backfill cost.
- pgvector extension availability.
- Multi-instance duplicate ingestion.

## 27. Estimated Implementation Complexity

| Area | Complexity | Reason |
|---|---:|---|
| Schema additions | Medium | Mostly additive, but needs careful compatibility. |
| Backfill scripts | High | Must preserve existing item identity and dedupe semantics. |
| Connector normalization | Medium-High | Each source type needs source-native payload mapping. |
| Entity/topic extraction | High | Requires NLP/LLM pipeline, confidence, and versioning. |
| Signal scoring engine | High | Needs new dimensions and score history. |
| Trend rollups | High | Requires time buckets, consensus, and source diversity. |
| Frontend migration | Medium-High | Existing cards and filters are article-shaped. |
| Operational hardening | Medium | DB locks/queues/observability needed for scale. |

## 28. Step-by-Step Implementation Roadmap

### Milestone 1 — Stabilize Current System

- Add contract tests for `/api/items`.
- Add integration tests for worker dedupe and score persistence.
- Fix repository mismatch for `scores.impact`, `engineering_value`, and `reasoning` before relying on these fields.
- Decide how to handle arXiv HTTP vs HTTPS fetch validation.

### Milestone 2 — Taxonomy Foundation

- Add source category/topic taxonomy tables.
- Backfill from current source/item categories and summary tags.
- Add admin-visible taxonomy metadata later.

### Milestone 3 — Signal Foundation

- Add `signals_v2`, `signal_observations`, and `item_signal_map`.
- Backfill all current items.
- Dual-write new ingestion into old `items` and new `signals_v2`.

### Milestone 4 — Intelligence Enrichment

- Add deterministic entity/topic extractors.
- Add LLM extractors behind flags.
- Store extraction confidence and method.

### Milestone 5 — Signal Scoring

- Implement signal score computation in parallel.
- Compare old final score vs new signal score.
- Build score history and explanations.

### Milestone 6 — Trend Engine

- Add hourly/daily trend rollups.
- Add momentum/growth/phase classification.
- Add cluster-to-signal migration.

### Milestone 7 — API Expansion

- Add `/api/signals`, `/api/trends`, `/api/entities`, and `/api/topics`.
- Keep `/api/items` stable until frontend fully migrates.

### Milestone 8 — Frontend Evolution

- Introduce signal cards and trend clusters.
- Add evidence drawer and score explanations.
- Add entity/topic pages.
- Gradually replace article terminology.

### Milestone 9 — Decommission Legacy Article Core

- Stop old-only writes.
- Archive old current `signals` engagement table or rename to `item_signal_metrics`.
- Replace item-first clusters with signal-first clusters.
- Keep compatibility views as long as external clients need them.
