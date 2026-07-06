-- 0004_add_restricted_key_fields.down.sql

ALTER TABLE public.api_credentials
    DROP COLUMN IF EXISTS allowed_ips,
    DROP COLUMN IF EXISTS scopes;
