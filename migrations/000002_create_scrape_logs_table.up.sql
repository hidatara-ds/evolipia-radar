CREATE TABLE IF NOT EXISTS scrape_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMP,
    items_processed INTEGER DEFAULT 0,
    items_new INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'running', -- running, success, failed
    error_message TEXT,
    trigger_source VARCHAR(50) DEFAULT 'github_actions' -- github_actions, manual, etc
);

CREATE INDEX idx_scrape_logs_started_at ON scrape_logs(started_at DESC);
