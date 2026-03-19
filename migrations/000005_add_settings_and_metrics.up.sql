-- Migration to add settings and enhance metrics persistence
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Enhance global_metrics to persist clustering stats
ALTER TABLE global_metrics ADD COLUMN IF NOT EXISTS clusters_count INT DEFAULT 0;
ALTER TABLE global_metrics ADD COLUMN IF NOT EXISTS avg_cluster_score DOUBLE PRECISION DEFAULT 0.0;
ALTER TABLE global_metrics ADD COLUMN IF NOT EXISTS top_cluster_titles JSONB DEFAULT '[]';
ALTER TABLE global_metrics ADD COLUMN IF NOT EXISTS last_closeness_update TIMESTAMP WITH TIME ZONE;
