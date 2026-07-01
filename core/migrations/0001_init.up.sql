-- 0001_init.up.sql
-- Core schema for Kobo. Run via golang-migrate.
-- Conventions:
--   - All IDs are UUIDv4, generated application-side (not gen_random_uuid()),
--     so Go code controls ID generation and can set them before insert when needed.
--   - All monetary amounts are BIGINT representing kobo (smallest Naira unit).
--   - All timestamps are TIMESTAMPTZ, stored and read as UTC.
--   - Every table has created_at; mutable tables also have updated_at.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ---------------------------------------------------------------------------
-- Integrators: one row per Kobo API consumer (e.g. the school-fees demo app).
-- Namespaces all other data; nothing below is visible across integrators.
-- ---------------------------------------------------------------------------
CREATE TABLE api_integrators (
    id              UUID PRIMARY KEY,
    name            TEXT NOT NULL,
    api_key_hash    TEXT NOT NULL UNIQUE,      -- store a hash, never the raw key
    api_secret_hash TEXT NOT NULL,             -- used to verify HMAC request signatures
    is_sandbox      BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- Identities: the root object. See ARCHITECTURE.md package boundary rules —
-- identity.id never changes; renames only touch display_name/metadata.
-- ---------------------------------------------------------------------------
CREATE TABLE identities (
    id                  UUID PRIMARY KEY,
    integrator_id       UUID NOT NULL REFERENCES api_integrators(id),
    external_reference  TEXT NOT NULL,          -- integrator's own ID, unique per integrator
    display_name        TEXT NOT NULL,
    kyc_tier            TEXT NOT NULL DEFAULT 'tier_1'
                            CHECK (kyc_tier IN ('tier_1', 'tier_2', 'tier_3')),
    state               TEXT NOT NULL DEFAULT 'pending'
                            CHECK (state IN ('pending', 'active', 'limited', 'closing', 'closed', 'failed')),
    failure_reason      TEXT,                   -- populated only when state = 'failed'
    metadata            JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),

    UNIQUE (integrator_id, external_reference)
);

CREATE INDEX idx_identities_integrator_state ON identities(integrator_id, state);

-- ---------------------------------------------------------------------------
-- Identity lifecycle events: append-only audit log of every state transition
-- and profile change. This is what makes renames and closures auditable,
-- per the "Renames" and lifecycle sections of ARCHITECTURE.md / openapi.yaml.
-- ---------------------------------------------------------------------------
CREATE TABLE identity_events (
    id              UUID PRIMARY KEY,
    identity_id     UUID NOT NULL REFERENCES identities(id),
    event_type      TEXT NOT NULL
                        CHECK (event_type IN (
                            'created', 'provisioned', 'provisioning_failed',
                            'activated', 'limited', 'unlimited',
                            'closing_started', 'closed', 'reopened',
                            'renamed', 'metadata_updated'
                        )),
    previous_state  TEXT,
    new_state       TEXT,
    detail          JSONB NOT NULL DEFAULT '{}'::jsonb,  -- e.g. {"old_name": "...", "new_name": "..."} for renames
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_identity_events_identity ON identity_events(identity_id, created_at);

-- ---------------------------------------------------------------------------
-- Virtual accounts: 1:1 with identities at any point in time. An identity
-- can have more than one row over its lifetime only if Nomba requires a new
-- accountRef on reopen (see NOMBA_INTEGRATION.md open item #2) — modeled as
-- a table rather than a single column on identities so that history survives
-- a reprovision.
-- ---------------------------------------------------------------------------
CREATE TABLE virtual_accounts (
    id                  UUID PRIMARY KEY,
    identity_id         UUID NOT NULL REFERENCES identities(id),
    nomba_account_ref   TEXT NOT NULL,           -- the accountRef Kobo sent to Nomba (= identity.id in the common case)
    account_number      TEXT,                    -- Nomba's bankAccountNumber, null until provisioning succeeds
    bank_name            TEXT,
    account_name        TEXT,
    is_active           BOOLEAN NOT NULL DEFAULT TRUE,  -- false once superseded by a reprovisioned account
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Only one active virtual account per identity at a time.
CREATE UNIQUE INDEX idx_virtual_accounts_identity_active
    ON virtual_accounts(identity_id)
    WHERE is_active = TRUE;

-- Account number lookups happen on every webhook/sweep event, so this must be fast.
CREATE UNIQUE INDEX idx_virtual_accounts_account_number
    ON virtual_accounts(account_number)
    WHERE account_number IS NOT NULL;

-- ---------------------------------------------------------------------------
-- Ledger entries: the reconciled record of every inbound transfer.
-- One row per matched transaction. Statements are derived by summing these,
-- never by re-querying Nomba live.
-- ---------------------------------------------------------------------------
CREATE TABLE ledger_entries (
    id                  UUID PRIMARY KEY,
    virtual_account_id  UUID NOT NULL REFERENCES virtual_accounts(id),
    identity_id         UUID NOT NULL REFERENCES identities(id),  -- denormalized for query convenience
    amount_kobo         BIGINT NOT NULL CHECK (amount_kobo > 0),
    direction           TEXT NOT NULL DEFAULT 'inbound' CHECK (direction = 'inbound'),
    status              TEXT NOT NULL CHECK (status IN ('matched', 'partial', 'overpayment')),
    nomba_reference      TEXT NOT NULL,           -- idempotency key, see idempotency_keys table
    source               TEXT NOT NULL CHECK (source IN ('webhook', 'sweep')),
    narration            TEXT,
    sender_name           TEXT,
    occurred_at          TIMESTAMPTZ NOT NULL,    -- Nomba's transaction time, not our insert time
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_ledger_entries_account ON ledger_entries(virtual_account_id, occurred_at);
CREATE INDEX idx_ledger_entries_identity_period ON ledger_entries(identity_id, occurred_at);

-- ---------------------------------------------------------------------------
-- Idempotency keys: dedup table for both webhook deliveries and sweep
-- backfills. A unique constraint on nomba_reference is what makes
-- reconciliation/idempotency.go's "check before write" logic atomic and safe
-- under concurrent webhook + sweep processing.
-- ---------------------------------------------------------------------------
CREATE TABLE idempotency_keys (
    nomba_reference  TEXT PRIMARY KEY,
    ledger_entry_id  UUID NOT NULL REFERENCES ledger_entries(id),
    first_seen_via   TEXT NOT NULL CHECK (first_seen_via IN ('webhook', 'sweep')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ---------------------------------------------------------------------------
-- Exceptions: misdirected payments and unmatched transfers. See Section 8 of
-- the Kobo concept note and the /v1/exceptions paths in openapi.yaml.
-- ---------------------------------------------------------------------------
CREATE TABLE exceptions (
    id                   UUID PRIMARY KEY,
    integrator_id        UUID NOT NULL REFERENCES api_integrators(id),
    type                 TEXT NOT NULL CHECK (type IN (
                             'payment_to_closed_account',
                             'payment_to_unknown_account',
                             'payment_during_closing'
                         )),
    amount_kobo          BIGINT NOT NULL,
    nomba_reference       TEXT NOT NULL,
    related_account_id   UUID REFERENCES virtual_accounts(id),  -- nullable: unknown-account case has no match
    status               TEXT NOT NULL DEFAULT 'open' CHECK (status IN ('open', 'resolved')),
    resolution_action    TEXT CHECK (resolution_action IN ('return_to_sender', 'redirect_to_successor', 'manual_override')),
    resolution_notes     TEXT,
    successor_identity_id UUID REFERENCES identities(id),
    detected_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    resolved_at          TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_exceptions_integrator_status ON exceptions(integrator_id, status, detected_at);

-- nomba_reference should not be double-flagged as an exception once resolved
-- and reprocessed; enforce at the application layer in exceptions/service.go
-- (not a DB constraint, since a reference could legitimately appear once as
-- an exception and, after a redirect resolution, later as a normal ledger
-- entry under a different identity).
