-- ============================================
-- Migration 000006: Add pgvector for semantic search
-- ============================================

-- 1. Enable the pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- 2. Add embedding column (1536-dim for text-embedding-3-small compatibility)
--    Nullable so existing rows are unaffected.
ALTER TABLE items ADD COLUMN IF NOT EXISTS embedding vector(1536);

-- 3. Track which model generated the embedding (for future model upgrades)
ALTER TABLE items ADD COLUMN IF NOT EXISTS embedding_model TEXT;

-- 4. HNSW index for fast approximate nearest-neighbor search using cosine distance.
--    m=16: number of bi-directional links per node (higher = more accurate, more memory)
--    ef_construction=64: size of dynamic candidate list during index build
CREATE INDEX IF NOT EXISTS items_embedding_hnsw_idx
    ON items USING hnsw (embedding vector_cosine_ops)
    WITH (m = 16, ef_construction = 64);
