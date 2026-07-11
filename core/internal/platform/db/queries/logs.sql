-- name: GetPaginatedRequestLogs :many
SELECT *
FROM request_logs
WHERE integrator_id = $1
  AND (sqlc.narg('method')::text IS NULL OR method = sqlc.narg('method'))
  AND (sqlc.narg('status_code')::int IS NULL OR status_code = sqlc.narg('status_code'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountRequestLogs :one
SELECT COUNT(*)
FROM request_logs
WHERE integrator_id = $1
  AND (sqlc.narg('method')::text IS NULL OR method = sqlc.narg('method'))
  AND (sqlc.narg('status_code')::int IS NULL OR status_code = sqlc.narg('status_code'));
