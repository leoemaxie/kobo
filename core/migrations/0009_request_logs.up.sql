CREATE TABLE request_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id UUID NOT NULL REFERENCES api_integrators(id) ON DELETE CASCADE,
    method VARCHAR(10) NOT NULL,
    path TEXT NOT NULL,
    status_code INTEGER NOT NULL,
    latency_ms INTEGER NOT NULL,
    request_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_request_logs_integrator_id_created_at ON request_logs(integrator_id, created_at DESC);

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_extension WHERE extname = 'pg_cron') THEN
        PERFORM cron.schedule('cleanup_request_logs', '0 0 * * *', 'DELETE FROM request_logs WHERE created_at < NOW() - INTERVAL ''30 days''');
    END IF;
END $$;
