package sqlc

import (
	"context"

	"github.com/google/uuid"
)

// Re-add types that sqlc failed to generate cleanly due to cross-schema enums
type ConsoleEnvironment string

const (
	ConsoleEnvironmentSandbox    ConsoleEnvironment = "sandbox"
	ConsoleEnvironmentProduction ConsoleEnvironment = "production"
)

type Invoice = ConsoleInvoice
type PaymentMethod = ConsolePaymentMethod

func (q *Queries) SuspendIntegrator(ctx context.Context, integratorID uuid.UUID) error {
	// Stub
	return nil
}

func (q *Queries) GenerateMonthlyInvoices(ctx context.Context, period string) error {
	rollupQuery := `
		INSERT INTO console.billing_records (
			integrator_id, environment, period,
			accounts_provisioned, transactions_processed, webhook_deliveries, amount_due_kobo
		)
		SELECT 
			integrator_id, environment, $1 as period,
			COUNT(CASE WHEN event_type = 'account_provisioned' THEN 1 END) as accounts_provisioned,
			COUNT(CASE WHEN event_type = 'transaction_processed' THEN 1 END) as transactions_processed,
			COUNT(CASE WHEN event_type = 'webhook_delivery' THEN 1 END) as webhook_deliveries,
			SUM(amount_kobo) as amount_due_kobo
		FROM console.usage_events
		WHERE to_char(occurred_at, 'YYYY-MM') = $1 AND environment = 'production'
		GROUP BY integrator_id, environment
		ON CONFLICT (integrator_id, period, environment) DO UPDATE SET
			accounts_provisioned = EXCLUDED.accounts_provisioned,
			transactions_processed = EXCLUDED.transactions_processed,
			webhook_deliveries = EXCLUDED.webhook_deliveries,
			amount_due_kobo = EXCLUDED.amount_due_kobo,
			synced_at = now();
	`
	_, err := q.db.Exec(ctx, rollupQuery, period)
	if err != nil {
		return err
	}

	invoiceQuery := `
		INSERT INTO console.invoices (
			integrator_id, billing_record_id, period, amount_kobo, status
		)
		SELECT 
			integrator_id, id, period, amount_due_kobo, 'open'
		FROM console.billing_records
		WHERE period = $1 AND environment = 'production' AND amount_due_kobo > 0
		ON CONFLICT DO NOTHING;
	`
	_, err = q.db.Exec(ctx, invoiceQuery, period)
	return err
}
