-- 0001_init.down.sql
-- Reverses 0001_init.up.sql exactly, in dependency order (children before parents).

DROP INDEX IF EXISTS idx_exceptions_integrator_status;
DROP TABLE IF EXISTS exceptions;

DROP TABLE IF EXISTS idempotency_keys;

DROP INDEX IF EXISTS idx_ledger_entries_identity_period;
DROP INDEX IF EXISTS idx_ledger_entries_account;
DROP TABLE IF EXISTS ledger_entries;

DROP INDEX IF EXISTS idx_virtual_accounts_account_number;
DROP INDEX IF EXISTS idx_virtual_accounts_identity_active;
DROP TABLE IF EXISTS virtual_accounts;

DROP INDEX IF EXISTS idx_identity_events_identity;
DROP TABLE IF EXISTS identity_events;

DROP INDEX IF EXISTS idx_identities_integrator_state;
DROP TABLE IF EXISTS identities;

DROP TABLE IF EXISTS api_integrators;

DROP EXTENSION IF EXISTS "uuid-ossp";
