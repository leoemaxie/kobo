-- internal/platform/db/queries/payouts.sql
-- SQLC queries for the payout system.
-- All tables live in the console schema. Balance lives on public.api_integrators.

-- ---------------------------------------------------------------------------
-- Bank account management
-- ---------------------------------------------------------------------------

-- name: GetActivePayoutBankAccount :one
SELECT * FROM console.payout_bank_accounts
WHERE integrator_id = $1 AND is_active = TRUE
LIMIT 1;

-- name: DeactivatePayoutBankAccounts :exec
-- Marks all active bank accounts for an integrator as inactive.
-- Always called inside a transaction before InsertPayoutBankAccount.
UPDATE console.payout_bank_accounts
SET is_active = FALSE, updated_at = now()
WHERE integrator_id = $1 AND is_active = TRUE;

-- name: InsertPayoutBankAccount :one
INSERT INTO console.payout_bank_accounts (
    integrator_id, account_number, account_name, bank_code, bank_name, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- ---------------------------------------------------------------------------
-- Payout lifecycle
-- ---------------------------------------------------------------------------

-- name: CreatePayout :one
INSERT INTO console.payouts (
    integrator_id,
    bank_account_id,
    requested_amount_kobo,
    platform_fee_kobo,
    transfer_fee_buffer_kobo,
    net_amount_kobo,
    merchant_tx_ref,
    initiated_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: UpdatePayoutStatus :one
UPDATE console.payouts
SET
    status                   = $2,
    failure_reason           = $3,
    nomba_transfer_id        = $4,
    actual_transfer_fee_kobo = $5,
    updated_at               = now()
WHERE id = $1
RETURNING *;

-- name: GetPayoutByMerchantTxRef :one
SELECT * FROM console.payouts
WHERE merchant_tx_ref = $1;

-- name: GetPayoutByID :one
SELECT * FROM console.payouts
WHERE id = $1;

-- name: ListPayoutsForIntegrator :many
SELECT * FROM console.payouts
WHERE integrator_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountInProgressPayouts :one
-- Returns the number of payouts currently in 'pending' or 'processing' state.
-- Used to prevent an integrator from initiating a second concurrent payout.
SELECT COUNT(*) FROM console.payouts
WHERE integrator_id = $1 AND status IN ('pending', 'processing');

-- ---------------------------------------------------------------------------
-- Balance operations on public.api_integrators
-- These supplement the existing UpdateIntegratorWalletBalance query in billing.sql.
-- ---------------------------------------------------------------------------

-- name: LockIntegratorRow :one
-- Acquires a row-level lock on the integrator for the duration of the caller's
-- transaction. MUST be called inside BEGIN/COMMIT. This prevents concurrent
-- payout requests from racing to deduct the same balance.
SELECT wallet_balance_kobo FROM public.api_integrators
WHERE id = $1
FOR UPDATE;

-- name: DeductIntegratorBalance :one
-- Atomically deducts the given amount from wallet_balance_kobo.
-- The WHERE clause ensures we never go negative: if balance < amount,
-- zero rows are updated and the caller must roll back.
UPDATE public.api_integrators
SET
    wallet_balance_kobo = wallet_balance_kobo - $2,
    updated_at          = now()
WHERE id = $1
  AND wallet_balance_kobo >= $2
RETURNING wallet_balance_kobo;

-- name: CreditIntegratorBalance :exec
-- Reversal: re-credits the balance when a payout fails after deduction.
-- Also used by the webhook handler to reverse a failed transfer.
UPDATE public.api_integrators
SET
    wallet_balance_kobo = wallet_balance_kobo + $2,
    updated_at          = now()
WHERE id = $1;
