-- name: InsertLedgerEntry :one
INSERT INTO ledger_entries (id, virtual_account_id, identity_id, amount_kobo, direction, status, nomba_reference, source, narration, sender_name, occurred_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: ListLedgerEntriesByAccount :many
SELECT * FROM ledger_entries
WHERE virtual_account_id = $1
ORDER BY occurred_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLedgerEntriesByIdentityAndPeriod :many
SELECT * FROM ledger_entries
WHERE identity_id = $1
  AND occurred_at >= $2
  AND occurred_at < $3
ORDER BY occurred_at ASC;

-- name: GetLedgerOpeningBalance :one
-- Calculates the sum of all 'inbound' matched transactions before the given time
SELECT COALESCE(SUM(amount_kobo), 0)::BIGINT AS balance
FROM ledger_entries
WHERE identity_id = $1
  AND occurred_at < $2
  AND status = 'matched';
