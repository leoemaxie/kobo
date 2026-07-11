-- name: CreateRequestLog :one
INSERT INTO request_logs (
    integrator_id, method, path, status_code, latency_ms, request_id
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetRecentRequestLogs :many
SELECT *
FROM request_logs
WHERE integrator_id = $1
ORDER BY created_at DESC
LIMIT 5;

-- name: GetErrorRate :one
SELECT 
    COALESCE(
        (COUNT(CASE WHEN status_code >= 500 THEN 1 END)::FLOAT / NULLIF(COUNT(*), 0)) * 100, 
        0
    )::FLOAT AS error_rate
FROM request_logs
WHERE integrator_id = $1 
  AND created_at >= NOW() - INTERVAL '30 days';

-- name: GetP99Latency :one
SELECT 
    COALESCE(
        percentile_cont(0.99) WITHIN GROUP (ORDER BY latency_ms), 
        0
    )::FLOAT AS p99_latency
FROM request_logs
WHERE integrator_id = $1 
  AND created_at >= NOW() - INTERVAL '30 days';

-- name: GetTotalApiRequests :one
SELECT COUNT(*)
FROM request_logs
WHERE integrator_id = $1
  AND created_at >= NOW() - INTERVAL '30 days';

-- name: CountVirtualAccountsByIntegrator :one
SELECT COUNT(*)
FROM virtual_accounts va
JOIN identities i ON va.identity_id = i.id
WHERE i.integrator_id = $1;
