# Evolipia Radar Architecture

Evolipia Radar is a scalable AI system built to intelligently discover, cluster, and synthesize high-signal news from the tech ecosystem, operating strictly within zero-cost AI budgets.

## 1. High-Level Pipeline
1. **Ingestion (Phase 3 Orchestrator)**: RSS feeds, Algolia/HackerNews trending APIs.
2. **Evaluation & Deduplication**: High-speed noise filtering, URL caching, content-length filters.
3. **Semantic Brain (Phase 2 ClusterService)**: Vector embeddings (`pgvector`), threshold-based grouping.
4. **Insights (Phase 2 LLM)**: Summarization, semantic titling, guardrails against LLM hallucination.
5. **Gateway (Phase 1 Vercel Adapters)**: REST endpoints exposed for dashboards.

## 2. Infrastructure Layer
- **Compute**: Vercel Serverless Functions (Golang runtime `api/*.go`)
- **Frontend**: Next.js App Router providing a minimalist read-only dashboard.
- **Database**: Supabase PostgreSQL with `pgvector`.
- **LLM Sub-system**: OpenRouter API wrapped by an internal `TokenTracker` middleware to actively block overages and implement graceful degradation (raw text mapping instead of AI summaries).
