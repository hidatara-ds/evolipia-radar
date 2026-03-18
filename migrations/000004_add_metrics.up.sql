-- Table to store global system metrics (for serverless persistence)
CREATE TABLE IF NOT EXISTS global_metrics (
    id SERIAL PRIMARY KEY,
    articles_processed INT DEFAULT 0,
    filtered_articles INT DEFAULT 0,
    api_hits INT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Initialize the single metrics row
INSERT INTO global_metrics (id, articles_processed, filtered_articles, api_hits)
VALUES (1, 0, 0, 0)
ON CONFLICT (id) DO NOTHING;
