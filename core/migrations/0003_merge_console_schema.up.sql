-- 0003_merge_console_schema.up.sql

CREATE SCHEMA IF NOT EXISTS console;

DO $$ BEGIN
    CREATE TYPE console.webhook_status AS ENUM ('active', 'disabled');
EXCEPTION WHEN duplicate_object THEN null; END $$;

DO $$ BEGIN
    CREATE TYPE console.environment AS ENUM ('sandbox', 'production');
EXCEPTION WHEN duplicate_object THEN null; END $$;

DO $$ BEGIN
    CREATE TYPE console.integrator_status AS ENUM ('active', 'suspended');
EXCEPTION WHEN duplicate_object THEN null; END $$;

DO $$ BEGIN
    CREATE TYPE console.user_role AS ENUM ('owner', 'member', 'superadmin');
EXCEPTION WHEN duplicate_object THEN null; END $$;

DO $$ BEGIN
    CREATE TYPE console.admin_action AS ENUM (
        'integrator_suspended',
        'integrator_reinstated',
        'production_access_granted',
        'credential_force_revoked',
        'billing_adjustment_applied',
        'user_session_revoked'
    );
EXCEPTION WHEN duplicate_object THEN null; END $$;

-- 1. Enhance public.api_integrators
ALTER TABLE public.api_integrators
    ADD COLUMN plan TEXT NOT NULL DEFAULT 'pay_as_you_go',
    ADD COLUMN status console.integrator_status NOT NULL DEFAULT 'active',
    ADD COLUMN production_access_granted BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN production_access_granted_at TIMESTAMPTZ,
    ADD COLUMN production_access_granted_by UUID;

-- 2. Drop legacy auth fields from api_integrators
ALTER TABLE public.api_integrators DROP COLUMN api_key;
ALTER TABLE public.api_integrators DROP COLUMN api_secret_hash;
ALTER TABLE public.api_integrators DROP COLUMN is_sandbox;

-- 3. Create public.api_credentials for auth
CREATE TABLE IF NOT EXISTS public.api_credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id UUID NOT NULL REFERENCES public.api_integrators(id),
    environment console.environment NOT NULL,
    key_id TEXT NOT NULL UNIQUE,
    secret_hash TEXT NOT NULL,
    label TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by UUID,
    rotated_at TIMESTAMPTZ,
    revoked_at TIMESTAMPTZ,
    revoked_by UUID,
    revoked_reason TEXT
);
CREATE INDEX IF NOT EXISTS idx_api_credentials_integrator_env ON public.api_credentials(integrator_id, environment);

-- 4. Create Console schema tables
CREATE TABLE IF NOT EXISTS console.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id UUID REFERENCES public.api_integrators(id),
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role console.user_role NOT NULL DEFAULT 'member',
    email_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Now we can add the FK from api_credentials.created_by to users.id
DO $$ BEGIN
    ALTER TABLE public.api_credentials
        ADD CONSTRAINT fk_api_credentials_created_by
        FOREIGN KEY (created_by) REFERENCES console.users(id);
EXCEPTION WHEN duplicate_object THEN null; END $$;

CREATE TABLE IF NOT EXISTS console.sessions (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES console.users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    revoked_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_sessions_user ON console.sessions(user_id);

CREATE TABLE IF NOT EXISTS console.email_verification_tokens (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES console.users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS console.password_reset_tokens (
    id TEXT PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES console.users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS console.billing_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id UUID NOT NULL REFERENCES public.api_integrators(id),
    environment console.environment NOT NULL,
    period TEXT NOT NULL,
    accounts_provisioned BIGINT NOT NULL DEFAULT 0,
    transactions_processed BIGINT NOT NULL DEFAULT 0,
    webhook_deliveries BIGINT NOT NULL DEFAULT 0,
    amount_due_kobo BIGINT NOT NULL DEFAULT 0,
    adjustment_kobo BIGINT NOT NULL DEFAULT 0,
    adjustment_reason TEXT,
    synced_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_billing_integrator_period_env ON console.billing_records(integrator_id, period, environment);

CREATE TABLE IF NOT EXISTS console.webhooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    integrator_id UUID NOT NULL REFERENCES public.api_integrators(id),
    environment console.environment NOT NULL,
    url TEXT NOT NULL,
    secret TEXT NOT NULL,
    status console.webhook_status NOT NULL DEFAULT 'active',
    events JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS console.admin_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_user_id UUID NOT NULL REFERENCES console.users(id),
    action console.admin_action NOT NULL,
    target_integrator_id UUID REFERENCES public.api_integrators(id),
    target_user_id UUID REFERENCES console.users(id),
    target_credential_id UUID REFERENCES public.api_credentials(id),
    detail JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
