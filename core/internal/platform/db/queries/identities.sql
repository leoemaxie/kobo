-- name: CreateIdentity :one
INSERT INTO identities (id, integrator_id, external_reference, display_name, state, metadata)
VALUES ($1, $2, $3, $4, 'pending', $5)
RETURNING *;

-- name: GetIdentityByID :one
SELECT * FROM identities
WHERE id = $1 AND integrator_id = $2;

-- name: GetIdentityByExternalReference :one
SELECT * FROM identities
WHERE integrator_id = $1 AND external_reference = $2;

-- name: DeleteIdentityCascade :exec
WITH deleted_va AS (
    DELETE FROM virtual_accounts WHERE identity_id = $1
), deleted_events AS (
    DELETE FROM identity_events WHERE identity_id = $1
)
DELETE FROM identities WHERE identities.id = $1 AND identities.integrator_id = $2;

-- name: UpdateIdentityProfile :one
-- Used for renames and metadata updates only. Never changes state or id.
UPDATE identities
SET display_name = COALESCE(sqlc.narg(display_name), display_name),
    metadata = COALESCE(sqlc.narg(metadata), metadata),
    updated_at = now()
WHERE id = $1 AND integrator_id = $2
RETURNING *;

-- name: UpdateIdentityState :one
-- The only query allowed to mutate `state`. Called exclusively from
-- account/lifecycle.go's Transition() function, never directly from a handler
-- or service method, per the package boundary rules in ARCHITECTURE.md.
UPDATE identities
SET state = $3,
    failure_reason = $4,
    updated_at = now()
WHERE id = $1 AND integrator_id = $2
RETURNING *;

-- name: ListIdentitiesByState :many
SELECT * FROM identities
WHERE integrator_id = $1 AND state = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: InsertIdentityEvent :one
INSERT INTO identity_events (id, identity_id, event_type, previous_state, new_state, detail)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListIdentityEvents :many
SELECT * FROM identity_events
WHERE identity_id = $1
ORDER BY created_at ASC;

-- name: ListAllIdentitiesByState :many
SELECT * FROM identities
WHERE state = $1
ORDER BY created_at ASC;

-- name: ListIdentities :many
SELECT * FROM identities
WHERE integrator_id = $1 
  AND (sqlc.narg('state')::text IS NULL OR state = sqlc.narg('state'))
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;
