-- ============================================
-- Rollback Migration 000006: Remove pgvector
-- ============================================

-- Drop index first (depends on column)
DROP INDEX IF EXISTS items_embedding_hnsw_idx;

-- Drop columns
ALTER TABLE items DROP COLUMN IF EXISTS embedding_model;
ALTER TABLE items DROP COLUMN IF EXISTS embedding;

-- Drop extension last
DROP EXTENSION IF EXISTS vector;
