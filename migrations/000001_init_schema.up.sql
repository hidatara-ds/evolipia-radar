-- Create sources table
CREATE TABLE sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    category TEXT NOT NULL,
    url TEXT NOT NULL,
    mapping_json JSONB NULL,
    enabled BOOLEAN NOT NULL DEFAULT FALSE,
    status TEXT NOT NULL DEFAULT 'pending',
    last_test_status TEXT NULL,
    last_test_message TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sources_enabled ON sources(enabled);
CREATE UNIQUE INDEX idx_sources_unique_url ON sources(url);

-- Create items table
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    url TEXT NOT NULL,
    published_at TIMESTAMPTZ NOT NULL,
    content_hash TEXT NOT NULL,
    domain TEXT NOT NULL,
    category TEXT NOT NULL,
    raw_excerpt TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_items_dedup ON items(content_hash);
CREATE INDEX idx_items_published_at ON items(published_at DESC);
CREATE INDEX idx_items_domain ON items(domain);
CREATE INDEX idx_items_source_id ON items(source_id);

-- Create signals table
CREATE TABLE signals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    points INT NULL,
    comments INT NULL,
    rank_pos INT NULL,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_signals_item_fetched ON signals(item_id, fetched_at DESC);

-- Create scores table
CREATE TABLE scores (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    hot DOUBLE PRECISION NOT NULL,
    relevance DOUBLE PRECISION NOT NULL,
    credibility DOUBLE PRECISION NOT NULL,
    novelty DOUBLE PRECISION NOT NULL,
    final DOUBLE PRECISION NOT NULL,
    computed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_scores_final ON scores(final DESC);

-- Create summaries table
CREATE TABLE summaries (
    item_id UUID PRIMARY KEY REFERENCES items(id) ON DELETE CASCADE,
    tldr TEXT NOT NULL,
    why_it_matters TEXT NOT NULL,
    tags JSONB NOT NULL,
    method TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_summaries_tags_gin ON summaries USING GIN (tags);

-- Create fetch_runs table
CREATE TABLE fetch_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES sources(id) ON DELETE CASCADE,
    fetched_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    status TEXT NOT NULL,
    error TEXT NULL,
    items_fetched INT NOT NULL DEFAULT 0,
    items_inserted INT NOT NULL DEFAULT 0
);

CREATE INDEX idx_fetch_runs_source_time ON fetch_runs(source_id, fetched_at DESC);
