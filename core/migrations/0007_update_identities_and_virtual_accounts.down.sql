ALTER TABLE identities ADD COLUMN kyc_tier TEXT NOT NULL DEFAULT 'tier_1' CHECK (kyc_tier IN ('tier_1', 'tier_2', 'tier_3'));

ALTER TABLE virtual_accounts DROP COLUMN expected_amount_kobo;
ALTER TABLE virtual_accounts DROP COLUMN is_expired;
