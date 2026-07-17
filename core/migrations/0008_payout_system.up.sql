-- 0008_payout_system.up.sql
-- Payout system: allows integrators to withdraw their Kobo ledger balance
-- to a registered Nigerian bank account via Monnify's transfer API.
--
-- Conventions:
--   - Tables live in the `console` schema (integrator-facing dashboard data).
--   - IDs are gen_random_uuid() (migration-layer convention post-0003).
--   - All monetary amounts are BIGINT in kobo.
--   - status is a TEXT CHECK constraint (not a PG ENUM) for easier future extension.

-- ---------------------------------------------------------------------------
-- console.payout_bank_accounts
-- One active bank account per integrator at any time.
-- Account name is resolved via Monnify lookup at save time and stored for audit.
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS console.payout_bank_accounts (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id   UUID NOT NULL REFERENCES public.api_integrators(id) ON DELETE CASCADE,
    account_number  TEXT NOT NULL,
    account_name    TEXT NOT NULL,  -- resolved by Monnify /v1/transfers/bank/lookup; stored for audit
    bank_code       TEXT NOT NULL,
    bank_name       TEXT NOT NULL,  -- human-readable e.g. "Guaranty Trust Bank"
    is_active       BOOLEAN NOT NULL DEFAULT TRUE,
    created_by      UUID NOT NULL REFERENCES console.users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Only one active bank account per integrator at a time.
-- Partial unique index mirrors the pattern used in virtual_accounts.
CREATE UNIQUE INDEX IF NOT EXISTS uq_payout_bank_account_active
    ON console.payout_bank_accounts(integrator_id)
    WHERE is_active = TRUE;

CREATE INDEX IF NOT EXISTS idx_payout_bank_accounts_integrator
    ON console.payout_bank_accounts(integrator_id);

-- ---------------------------------------------------------------------------
-- console.payouts
-- One row per payout attempt. State machine:
--   pending → processing → successful
--   pending → failed
--   processing → successful | failed
--
-- Balance is deducted from api_integrators.wallet_balance_kobo atomically
-- inside a DB transaction before this row is created. On failure, the
-- deduction is reversed via CreditIntegratorBalance.
-- ---------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS console.payouts (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id               UUID NOT NULL REFERENCES public.api_integrators(id),
    bank_account_id             UUID NOT NULL REFERENCES console.payout_bank_accounts(id),

    -- Amount fields (all in kobo)
    requested_amount_kobo       BIGINT NOT NULL CHECK (requested_amount_kobo > 0),
    platform_fee_kobo           BIGINT NOT NULL DEFAULT 0,     -- Kobo's cut; reserved for future pricing
    transfer_fee_buffer_kobo    BIGINT NOT NULL DEFAULT 5000,  -- 50 NGN held back to cover Monnify's fee
    actual_transfer_fee_kobo    BIGINT,                        -- populated from Monnify response; NULL until resolved
    net_amount_kobo             BIGINT NOT NULL,               -- amount actually sent to bank = requested - platform_fee

    -- State machine
    status          TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'processing', 'successful', 'failed')),
    failure_reason  TEXT,

    -- merchantTxRef sent to Monnify — doubles as our idempotency key.
    -- Always prefixed 'payout_' so webhook handler can route it correctly.
    merchant_tx_ref TEXT NOT NULL UNIQUE,

    -- Populated after Monnify confirms the transfer (200 sync or webhook).
    monnify_transfer_id   TEXT,

    initiated_by    UUID NOT NULL REFERENCES console.users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_payouts_integrator
    ON console.payouts(integrator_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_payouts_merchant_tx_ref
    ON console.payouts(merchant_tx_ref);

-- Cheap scan for in-progress payouts (used to block concurrent requests).
CREATE INDEX IF NOT EXISTS idx_payouts_in_progress
    ON console.payouts(integrator_id)
    WHERE status IN ('pending', 'processing');
