-- name: GetConsoleSession :one
SELECT s.user_id, u.integrator_id, u.role
FROM console.sessions s
JOIN console.users u ON s.user_id = u.id
WHERE s.id = $1 AND s.expires_at > now() AND s.revoked_at IS NULL;
