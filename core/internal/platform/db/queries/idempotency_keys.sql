-- name: InsertIdempotencyKey :one
INSERT INTO idempotency_keys (nomba_reference, ledger_entry_id, first_seen_via)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetIdempotencyKey :one
SELECT * FROM idempotency_keys
WHERE nomba_reference = $1;
