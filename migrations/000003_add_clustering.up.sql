-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Table to store AI-generated insight clusters
CREATE TABLE IF NOT EXISTS clusters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    summary TEXT NOT NULL,
    embedding vector(1536), -- Assuming typical OpenAI text-embedding-3/ada-002 dimension size
    score FLOAT8 DEFAULT 0.0, -- Used for cluster ranking
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Index for fast cosine similarity search
CREATE INDEX IF NOT EXISTS clusters_embedding_idx ON clusters USING hnsw (embedding vector_cosine_ops);

-- Table mapping raw articles (sources) to their synthesized cluster
CREATE TABLE IF NOT EXISTS cluster_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cluster_id UUID NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    article_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE, -- Assuming 'items' is the existing article table
    assigned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(cluster_id, article_id)
);

-- Index for looking up sources by cluster
CREATE INDEX IF NOT EXISTS idx_cluster_sources_cluster_id ON cluster_sources(cluster_id);
CREATE INDEX IF NOT EXISTS idx_cluster_sources_article_id ON cluster_sources(article_id);
