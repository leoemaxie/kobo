-- name: InsertUsageEvent :exec
INSERT INTO console.usage_events (
    integrator_id, environment, event_type, reference_id, amount_kobo
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: UpsertBillingRecord :exec
INSERT INTO console.billing_records (
    integrator_id, environment, period, 
    accounts_provisioned, transactions_processed, webhook_deliveries, amount_due_kobo
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) ON CONFLICT (integrator_id, period, environment) DO UPDATE SET
    accounts_provisioned = console.billing_records.accounts_provisioned + EXCLUDED.accounts_provisioned,
    transactions_processed = console.billing_records.transactions_processed + EXCLUDED.transactions_processed,
    webhook_deliveries = console.billing_records.webhook_deliveries + EXCLUDED.webhook_deliveries,
    amount_due_kobo = console.billing_records.amount_due_kobo + EXCLUDED.amount_due_kobo,
    synced_at = now();

-- name: InsertPaymentMethod :one
INSERT INTO console.payment_methods (
    integrator_id, nomba_token_key, card_last4, card_brand, is_default
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetDefaultPaymentMethod :one
SELECT * FROM console.payment_methods
WHERE integrator_id = $1 AND is_default = TRUE
LIMIT 1;

-- name: InsertInvoice :one
INSERT INTO console.invoices (
    integrator_id, billing_record_id, period, amount_kobo, status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: UpdateInvoiceStatus :exec
UPDATE console.invoices
SET status = $2, nomba_order_ref = $3, paid_at = $4, retry_count = $5, next_retry_at = $6
WHERE id = $1;

-- name: GetPendingInvoices :many
SELECT * FROM console.invoices
WHERE status = 'open' OR status = 'failed' AND retry_count < 3 AND (next_retry_at IS NULL OR next_retry_at <= now());

-- name: GetIntegratorInvoices :many
SELECT * FROM console.invoices
WHERE integrator_id = $1
ORDER BY created_at DESC;

-- name: GetIntegratorUsageEvents :many
SELECT * FROM console.usage_events
WHERE integrator_id = $1 AND environment = $2 AND occurred_at >= $3 AND occurred_at <= $4
ORDER BY occurred_at DESC;

-- name: GetBillingRecords :many
SELECT * FROM console.billing_records
WHERE integrator_id = $1
ORDER BY synced_at DESC;

-- name: GetIntegratorWalletBalance :one
SELECT wallet_balance_kobo FROM public.api_integrators WHERE id = $1;

-- name: UpdateIntegratorWalletBalance :exec
UPDATE public.api_integrators
SET wallet_balance_kobo = wallet_balance_kobo + $2
WHERE id = $1;

-- name: GetIdentityByVirtualAccountID :one
SELECT i.*, c.environment as credential_environment
FROM public.identities i
JOIN public.virtual_accounts v ON i.id = v.identity_id
JOIN public.api_credentials c ON i.integrator_id = c.integrator_id
WHERE v.id = $1
LIMIT 1;
