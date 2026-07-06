-- 0005_billing.down.sql

DROP TABLE IF EXISTS console.invoices;
DROP TABLE IF EXISTS console.payment_methods;
DROP TABLE IF EXISTS console.usage_events;

ALTER TABLE public.api_integrators
DROP COLUMN IF EXISTS wallet_balance_kobo;
