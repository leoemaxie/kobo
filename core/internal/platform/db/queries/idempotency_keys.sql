-- name: InsertIdempotencyKey :one
INSERT INTO idempotency_keys (monnify_reference, ledger_entry_id, first_seen_via)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetIdempotencyKey :one
SELECT * FROM idempotency_keys
WHERE monnify_reference = $1;
