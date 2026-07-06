-- 0004_add_restricted_key_fields.up.sql

ALTER TABLE public.api_credentials
    ADD COLUMN IF NOT EXISTS allowed_ips TEXT[] NOT NULL DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS scopes TEXT[] NOT NULL DEFAULT '{}';
