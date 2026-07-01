-- name: InsertException :one
INSERT INTO exceptions (id, integrator_id, type, amount_kobo, nomba_reference, related_account_id, status, detected_at)
VALUES ($1, $2, $3, $4, $5, $6, 'open', now())
RETURNING *;

-- name: GetExceptionByID :one
SELECT * FROM exceptions
WHERE id = $1 AND integrator_id = $2;

-- name: ResolveException :one
UPDATE exceptions
SET status = 'resolved',
    resolution_action = $3,
    resolution_notes = $4,
    successor_identity_id = $5,
    resolved_at = now()
WHERE id = $1 AND integrator_id = $2 AND status = 'open'
RETURNING *;

-- name: ListOpenExceptions :many
SELECT * FROM exceptions
WHERE integrator_id = $1 AND status = 'open'
ORDER BY detected_at ASC
LIMIT $2 OFFSET $3;

-- name: ListExceptionsByStatus :many
SELECT * FROM exceptions
WHERE integrator_id = $1 AND status = $2
ORDER BY detected_at DESC
LIMIT $3 OFFSET $4;
