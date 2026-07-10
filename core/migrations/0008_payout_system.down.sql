-- 0008_payout_system.down.sql

DROP INDEX IF EXISTS console.idx_payouts_in_progress;
DROP INDEX IF EXISTS console.idx_payouts_merchant_tx_ref;
DROP INDEX IF EXISTS console.idx_payouts_integrator;
DROP TABLE IF EXISTS console.payouts;

DROP INDEX IF EXISTS console.uq_payout_bank_account_active;
DROP INDEX IF EXISTS console.idx_payout_bank_accounts_integrator;
DROP TABLE IF EXISTS console.payout_bank_accounts;
