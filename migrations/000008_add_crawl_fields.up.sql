-- Add crawl status, error, relevance score, and validation timestamp to items table
ALTER TABLE items
ADD COLUMN IF NOT EXISTS crawl_status TEXT NOT NULL DEFAULT 'done',
ADD COLUMN IF NOT EXISTS crawl_error TEXT NULL,
ADD COLUMN IF NOT EXISTS relevance_score INT NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS validated_at TIMESTAMPTZ NULL;

CREATE INDEX IF NOT EXISTS idx_items_crawl_status ON items(crawl_status);
CREATE INDEX IF NOT EXISTS idx_items_relevance_score ON items(relevance_score DESC);
