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

// Add these to Querier interface manually if needed, but since Go uses structural typing,
// we just need to add them to the actual type or interface used.
