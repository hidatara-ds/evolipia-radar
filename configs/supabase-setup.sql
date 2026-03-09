-- Run this in Supabase SQL Editor setelah migration

-- Enable RLS (Row Level Security) untuk keamanan
ALTER TABLE items ENABLE ROW LEVEL SECURITY;

-- Policy: Allow read untuk semua (anonymous) - untuk Flutter app
CREATE POLICY "Allow anonymous read" ON items
    FOR SELECT USING (true);

-- Policy: Allow insert hanya dari specific role (worker)
-- Note: Worker pakai service_role key, bypass RLS
-- Atau bisa buat policy berdasarkan JWT claims

-- View untuk Flutter app (simplified)
CREATE OR REPLACE VIEW daily_feed AS
SELECT 
    id,
    title,
    url,
    summary,
    score,
    published_at,
    source,
    topic
FROM items
WHERE published_at >= CURRENT_DATE - INTERVAL '1 day'
ORDER BY score DESC;

-- Function untuk get feed (bisa dipanggil dari Flutter via RPC)
CREATE OR REPLACE FUNCTION get_daily_feed(limit_count INTEGER DEFAULT 20)
RETURNS SETOF items AS $$
BEGIN
    RETURN QUERY
    SELECT * FROM items
    WHERE published_at >= CURRENT_DATE - INTERVAL '1 day'
    ORDER BY score DESC
    LIMIT limit_count;
END;
$$ LANGUAGE plpgsql;
