DROP TABLE IF EXISTS cluster_sources;
DROP TABLE IF EXISTS clusters;
-- Note: Reverting the extension is dangerous if other tables use pgvector:
-- DROP EXTENSION IF EXISTS vector;
