DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_cron') THEN
        PERFORM cron.unschedule('cleanup_request_logs');
    END IF;
END $$;

DROP TABLE IF EXISTS request_logs;
