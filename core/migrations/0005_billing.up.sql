-- 0005_billing.up.sql

ALTER TABLE public.api_integrators
ADD COLUMN wallet_balance_kobo BIGINT NOT NULL DEFAULT 0;

-- Granular usage events for auditing (ring-buffer, pruned after 90 days)
CREATE TABLE console.usage_events (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id  UUID NOT NULL REFERENCES public.api_integrators(id),
    environment    console.environment NOT NULL,
    event_type     TEXT NOT NULL CHECK (event_type IN (
                       'account_provisioned', 'transaction_processed', 'webhook_delivered'
                   )),
    reference_id   TEXT NOT NULL,   -- identity_id / ledger_entry_id / webhook delivery ID
    amount_kobo    BIGINT NOT NULL DEFAULT 0,
    occurred_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX ON console.usage_events(integrator_id, occurred_at);

-- Saved payment methods (Monnify card tokens)
CREATE TABLE console.payment_methods (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id  UUID NOT NULL REFERENCES public.api_integrators(id),
    monnify_token_key TEXT NOT NULL,    -- returned by Monnify after first checkout
    card_last4     TEXT,
    card_brand     TEXT,              -- 'visa' / 'mastercard'
    is_default     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX ON console.payment_methods(integrator_id) WHERE is_default = TRUE;

-- Invoices
CREATE TABLE console.invoices (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id      UUID NOT NULL REFERENCES public.api_integrators(id),
    billing_record_id  UUID NOT NULL REFERENCES console.billing_records(id),
    period             TEXT NOT NULL,             -- '2026-10'
    amount_kobo        BIGINT NOT NULL,
    status             TEXT NOT NULL DEFAULT 'open'
                           CHECK (status IN ('open','paid','failed','void')),
    monnify_order_ref    TEXT,                      -- Monnify checkout/charge reference
    paid_at            TIMESTAMPTZ,
    retry_count        INT NOT NULL DEFAULT 0,
    next_retry_at      TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now()
);
