DROP INDEX IF EXISTS idx_items_relevance_score;
DROP INDEX IF EXISTS idx_items_crawl_status;

ALTER TABLE items
DROP COLUMN IF EXISTS validated_at,
DROP COLUMN IF EXISTS relevance_score,
DROP COLUMN IF EXISTS crawl_error,
DROP COLUMN IF EXISTS crawl_status;
