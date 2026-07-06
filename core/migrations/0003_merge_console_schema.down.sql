-- 0003_merge_console_schema.down.sql

DROP TABLE IF EXISTS console.admin_audit_log;
DROP TABLE IF EXISTS console.webhooks;
DROP TABLE IF EXISTS console.billing_records;
DROP TABLE IF EXISTS console.password_reset_tokens;
DROP TABLE IF EXISTS console.email_verification_tokens;
DROP TABLE IF EXISTS console.sessions;

ALTER TABLE public.api_credentials DROP CONSTRAINT IF EXISTS fk_api_credentials_created_by;
DROP TABLE IF EXISTS console.users;
DROP TABLE IF EXISTS public.api_credentials;

ALTER TABLE public.api_integrators DROP COLUMN IF EXISTS plan;
ALTER TABLE public.api_integrators DROP COLUMN IF EXISTS status;
ALTER TABLE public.api_integrators DROP COLUMN IF EXISTS production_access_granted;
ALTER TABLE public.api_integrators DROP COLUMN IF EXISTS production_access_granted_at;
ALTER TABLE public.api_integrators DROP COLUMN IF EXISTS production_access_granted_by;

-- Note: Cannot restore dropped data, but we can restore the schema
ALTER TABLE public.api_integrators ADD COLUMN api_key TEXT NOT NULL DEFAULT '';
ALTER TABLE public.api_integrators ADD COLUMN api_secret_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE public.api_integrators ADD COLUMN is_sandbox BOOLEAN NOT NULL DEFAULT true;

DROP TYPE IF EXISTS console.admin_action;
DROP TYPE IF EXISTS console.user_role;
DROP TYPE IF EXISTS console.integrator_status;
DROP TYPE IF EXISTS console.environment;
DROP TYPE IF EXISTS console.webhook_status;

DROP SCHEMA IF EXISTS console CASCADE;
