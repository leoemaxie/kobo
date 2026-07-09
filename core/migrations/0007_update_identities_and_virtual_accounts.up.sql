ALTER TABLE identities DROP COLUMN kyc_tier;

ALTER TABLE virtual_accounts ADD COLUMN expected_amount_kobo BIGINT;
ALTER TABLE virtual_accounts ADD COLUMN is_expired BOOLEAN NOT NULL DEFAULT FALSE;
