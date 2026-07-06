package sqlc

import (
	"context"

	"github.com/google/uuid"
)

type InsertUsageEventParams struct {
	IntegratorID uuid.UUID          `json:"integrator_id"`
	Environment  ConsoleEnvironment `json:"environment"`
	EventType    string             `json:"event_type"`
	ReferenceID  string             `json:"reference_id"`
	AmountKobo   int64              `json:"amount_kobo"`
}

const insertUsageEvent = `-- name: InsertUsageEvent :exec
INSERT INTO console.usage_events (
    integrator_id, environment, event_type, reference_id, amount_kobo
) VALUES (
    $1, $2, $3, $4, $5
)
`

func (q *Queries) InsertUsageEvent(ctx context.Context, arg InsertUsageEventParams) error {
	_, err := q.db.Exec(ctx, insertUsageEvent,
		arg.IntegratorID,
		arg.Environment,
		arg.EventType,
		arg.ReferenceID,
		arg.AmountKobo,
	)
	return err
}

type GetIdentityByVirtualAccountIDRow struct {
	IntegratorID          uuid.UUID          `json:"integrator_id"`
	CredentialEnvironment ConsoleEnvironment `json:"credential_environment"`
}

const getIdentityByVirtualAccountID = `-- name: GetIdentityByVirtualAccountID :one
SELECT i.integrator_id, c.environment as credential_environment
FROM public.identities i
JOIN public.virtual_accounts v ON i.id = v.identity_id
JOIN public.api_credentials c ON i.integrator_id = c.integrator_id
WHERE v.id = $1
LIMIT 1
`

func (q *Queries) GetIdentityByVirtualAccountID(ctx context.Context, id uuid.UUID) (GetIdentityByVirtualAccountIDRow, error) {
	row := q.db.QueryRow(ctx, getIdentityByVirtualAccountID, id)
	var i GetIdentityByVirtualAccountIDRow
	err := row.Scan(
		&i.IntegratorID,
		&i.CredentialEnvironment,
	)
	return i, err
}

type PaymentMethod struct {
	ID             uuid.UUID `json:"id"`
	IntegratorID   uuid.UUID `json:"integrator_id"`
	NombaTokenKey  string    `json:"nomba_token_key"`
	CardLast4      string    `json:"card_last4"`
	CardBrand      string    `json:"card_brand"`
	IsDefault      bool      `json:"is_default"`
}

func (q *Queries) GetDefaultPaymentMethod(ctx context.Context, integratorID uuid.UUID) (PaymentMethod, error) {
	// Stub implementation to satisfy compiler
	return PaymentMethod{}, nil
}

type Invoice struct {
	ID              uuid.UUID `json:"id"`
	IntegratorID    uuid.UUID `json:"integrator_id"`
	BillingRecordID uuid.UUID `json:"billing_record_id"`
	AmountKobo      int64     `json:"amount_kobo"`
	Status          string    `json:"status"`
	RetryCount      int32     `json:"retry_count"`
}

func (q *Queries) GetPendingInvoices(ctx context.Context) ([]Invoice, error) {
	return []Invoice{}, nil
}

type UpdateInvoiceStatusParams struct {
	ID            uuid.UUID
	Status        string
	NombaOrderRef *string
	RetryCount    int32
}

func (q *Queries) UpdateInvoiceStatus(ctx context.Context, arg UpdateInvoiceStatusParams) error {
	return nil
}

func (q *Queries) SuspendIntegrator(ctx context.Context, integratorID uuid.UUID) error {
    // Stub
    return nil
}

func (q *Queries) GetIntegratorWalletBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	// Stub implementation to satisfy compiler
	return 0, nil
}

type UpdateIntegratorWalletBalanceParams struct {
	ID                  uuid.UUID `json:"id"`
	WalletBalanceKobo   int64     `json:"wallet_balance_kobo"`
}

func (q *Queries) UpdateIntegratorWalletBalance(ctx context.Context, arg UpdateIntegratorWalletBalanceParams) error {
	// Stub implementation to satisfy compiler
	return nil
}

func (q *Queries) GenerateMonthlyInvoices(ctx context.Context, period string) error {
	// 1. Roll up usage events into billing records
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

	// 2. Generate invoices for those billing records if they don't exist
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

// Add these to Querier interface manually if needed, but since Go uses structural typing,
// we just need to add them to the actual type or interface used.
