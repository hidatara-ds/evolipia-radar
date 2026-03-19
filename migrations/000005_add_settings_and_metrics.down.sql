DROP TABLE IF EXISTS settings;
ALTER TABLE global_metrics DROP COLUMN IF EXISTS clusters_count;
ALTER TABLE global_metrics DROP COLUMN IF EXISTS avg_cluster_score;
ALTER TABLE global_metrics DROP COLUMN IF EXISTS last_closeness_update;
