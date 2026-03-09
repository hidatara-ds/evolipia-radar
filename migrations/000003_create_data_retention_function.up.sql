-- Function untuk cleanup data lama (45 hari)
CREATE OR REPLACE FUNCTION cleanup_old_data() RETURNS void AS $$
BEGIN
    -- Delete news items older than 45 days
    DELETE FROM items 
    WHERE published_at < NOW() - INTERVAL '45 days';
    
    -- Delete scrape logs older than 60 days (keep longer for audit)
    DELETE FROM scrape_logs 
    WHERE started_at < NOW() - INTERVAL '60 days';
    
    -- Vacuum akan di-handle otomatis oleh Supabase
END;
$$ LANGUAGE plpgsql;
