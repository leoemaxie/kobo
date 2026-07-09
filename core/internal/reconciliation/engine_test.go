package reconciliation

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/billing"
	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/stretchr/testify/assert"
)

// mockQuerier embeds sqlc.Querier so we only need to override methods we test.
type mockQuerier struct {
	sqlc.Querier
	GetVirtualAccountByAccountNumberFunc func(ctx context.Context, accountNumber pgtype.Text) (sqlc.VirtualAccount, error)
	InsertLedgerEntryFunc                func(ctx context.Context, arg sqlc.InsertLedgerEntryParams) (sqlc.LedgerEntry, error)
	ListActiveVirtualAccountsFunc        func(ctx context.Context) ([]sqlc.VirtualAccount, error)
	GetIdentityByVirtualAccountIDFunc    func(ctx context.Context, id uuid.UUID) (sqlc.GetIdentityByVirtualAccountIDRow, error)
	InsertUsageEventFunc                 func(ctx context.Context, arg sqlc.InsertUsageEventParams) error
	UnsetDefaultPaymentMethodsFunc       func(ctx context.Context, integratorID uuid.UUID) error
	InsertPaymentMethodFunc              func(ctx context.Context, arg sqlc.InsertPaymentMethodParams) (sqlc.ConsolePaymentMethod, error)
}

func (m *mockQuerier) GetVirtualAccountByAccountNumber(ctx context.Context, accountNumber pgtype.Text) (sqlc.VirtualAccount, error) {
	if m.GetVirtualAccountByAccountNumberFunc != nil {
		return m.GetVirtualAccountByAccountNumberFunc(ctx, accountNumber)
	}
	return sqlc.VirtualAccount{}, errors.New("unimplemented GetVirtualAccountByAccountNumber")
}

func (m *mockQuerier) InsertLedgerEntry(ctx context.Context, arg sqlc.InsertLedgerEntryParams) (sqlc.LedgerEntry, error) {
	if m.InsertLedgerEntryFunc != nil {
		return m.InsertLedgerEntryFunc(ctx, arg)
	}
	return sqlc.LedgerEntry{}, errors.New("unimplemented InsertLedgerEntry")
}

func (m *mockQuerier) ListActiveVirtualAccounts(ctx context.Context) ([]sqlc.VirtualAccount, error) {
	if m.ListActiveVirtualAccountsFunc != nil {
		return m.ListActiveVirtualAccountsFunc(ctx)
	}
	return nil, errors.New("unimplemented ListActiveVirtualAccounts")
}

func (m *mockQuerier) GetIdentityByVirtualAccountID(ctx context.Context, id uuid.UUID) (sqlc.GetIdentityByVirtualAccountIDRow, error) {
	if m.GetIdentityByVirtualAccountIDFunc != nil {
		return m.GetIdentityByVirtualAccountIDFunc(ctx, id)
	}
	return sqlc.GetIdentityByVirtualAccountIDRow{}, errors.New("unimplemented GetIdentityByVirtualAccountID")
}

func (m *mockQuerier) InsertUsageEvent(ctx context.Context, arg sqlc.InsertUsageEventParams) error {
	if m.InsertUsageEventFunc != nil {
		return m.InsertUsageEventFunc(ctx, arg)
	}
	return nil // Allow it to pass if not explicitly mocked
}

func (m *mockQuerier) UnsetDefaultPaymentMethods(ctx context.Context, integratorID uuid.UUID) error {
	if m.UnsetDefaultPaymentMethodsFunc != nil {
		return m.UnsetDefaultPaymentMethodsFunc(ctx, integratorID)
	}
	return nil
}

func (m *mockQuerier) InsertPaymentMethod(ctx context.Context, arg sqlc.InsertPaymentMethodParams) (sqlc.ConsolePaymentMethod, error) {
	if m.InsertPaymentMethodFunc != nil {
		return m.InsertPaymentMethodFunc(ctx, arg)
	}
	return sqlc.ConsolePaymentMethod{}, nil
}

// mockIdempotencyRepo
type mockIdempotencyRepo struct {
	CheckOrSetIdempotencyFunc func(ctx context.Context, reference string, source string, insertLedgerFunc func() (uuid.UUID, error)) (bool, error)
}

func (m *mockIdempotencyRepo) CheckOrSetIdempotency(ctx context.Context, reference string, source string, insertLedgerFunc func() (uuid.UUID, error)) (bool, error) {
	if m.CheckOrSetIdempotencyFunc != nil {
		return m.CheckOrSetIdempotencyFunc(ctx, reference, source, insertLedgerFunc)
	}
	return false, errors.New("unimplemented CheckOrSetIdempotency")
}

// mockNombaClient
type mockNombaClient struct {
	FetchSingleTransactionFunc func(ctx context.Context, transactionRef string) (*nomba.TransactionResult, error)
}

func (m *mockNombaClient) FetchSingleTransaction(ctx context.Context, transactionRef string) (*nomba.TransactionResult, error) {
	if m.FetchSingleTransactionFunc != nil {
		return m.FetchSingleTransactionFunc(ctx, transactionRef)
	}
	return &nomba.TransactionResult{Status: "SUCCESS"}, nil
}

func TestProcessWebhook_IgnoreNonPaymentSuccess(t *testing.T) {
	eng := NewEngine(nil, nil, nil, nil)
	payload := &nomba.WebhookPayload{EventType: "some_other_event"}
	err := eng.ProcessWebhook(context.Background(), payload)
	assert.NoError(t, err)
}

func TestProcessWebhook_IgnoreNonVirtual(t *testing.T) {
	eng := NewEngine(nil, nil, nil, nil)
	payload := &nomba.WebhookPayload{EventType: "payment_success"}
	payload.Data.Transaction.AliasAccountType = "CARD"
	err := eng.ProcessWebhook(context.Background(), payload)
	assert.NoError(t, err)
}

func TestProcessWebhook_AccountNotFound(t *testing.T) {
	mq := &mockQuerier{
		GetVirtualAccountByAccountNumberFunc: func(ctx context.Context, accountNumber pgtype.Text) (sqlc.VirtualAccount, error) {
			return sqlc.VirtualAccount{}, errors.New("not found")
		},
	}
	eng := NewEngine(mq, nil, nil, nil)

	payload := &nomba.WebhookPayload{EventType: "payment_success"}
	payload.Data.Transaction.AliasAccountType = "VIRTUAL"
	payload.Data.Transaction.AliasAccountNumber = "12345"

	err := eng.ProcessWebhook(context.Background(), payload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "virtual account not found")
}

func TestProcessWebhook_Success(t *testing.T) {
	accID := uuid.New()
	identID := uuid.New()

	mq := &mockQuerier{
		GetVirtualAccountByAccountNumberFunc: func(ctx context.Context, accountNumber pgtype.Text) (sqlc.VirtualAccount, error) {
			return sqlc.VirtualAccount{
				ID:            accID,
				IdentityID:    identID,
				AccountNumber: pgtype.Text{String: "12345", Valid: true},
			}, nil
		},
		InsertLedgerEntryFunc: func(ctx context.Context, arg sqlc.InsertLedgerEntryParams) (sqlc.LedgerEntry, error) {
			assert.Equal(t, accID, arg.VirtualAccountID)
			assert.Equal(t, int64(15000), arg.AmountKobo) // 150.0 * 100
			assert.Equal(t, "txn-123", arg.NombaReference)
			return sqlc.LedgerEntry{ID: arg.ID}, nil
		},
		GetIdentityByVirtualAccountIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.GetIdentityByVirtualAccountIDRow, error) {
			return sqlc.GetIdentityByVirtualAccountIDRow{
				IntegratorID:          uuid.New(),
				CredentialEnvironment: "sandbox",
			}, nil
		},
		InsertUsageEventFunc: func(ctx context.Context, arg sqlc.InsertUsageEventParams) error {
			return nil
		},
	}

	mIdem := &mockIdempotencyRepo{
		CheckOrSetIdempotencyFunc: func(ctx context.Context, reference string, source string, insertLedgerFunc func() (uuid.UUID, error)) (bool, error) {
			assert.Equal(t, "txn-123", reference)
			_, err := insertLedgerFunc()
			return false, err // not duplicate
		},
	}

	mNomba := &mockNombaClient{
		FetchSingleTransactionFunc: func(ctx context.Context, transactionRef string) (*nomba.TransactionResult, error) {
			return &nomba.TransactionResult{Status: "SUCCESS"}, nil
		},
	}

	recorder := billing.NewUsageRecorder(mq)
	eng := NewEngine(mq, mIdem, recorder, mNomba)

	payload := &nomba.WebhookPayload{EventType: "payment_success"}
	payload.Data.Transaction.AliasAccountType = "VIRTUAL"
	payload.Data.Transaction.AliasAccountNumber = "12345"
	payload.Data.Transaction.TransactionAmount = 150.0
	payload.Data.Transaction.TransactionID = "txn-123"
	payload.Data.Transaction.Time = time.Now().Format(time.RFC3339)

	err := eng.ProcessWebhook(context.Background(), payload)
	assert.NoError(t, err)
}

func TestProcessWebhook_CheckoutSuccess(t *testing.T) {
	integratorID := uuid.New()
	unsetCalled := false
	insertCalled := false

	mq := &mockQuerier{
		UnsetDefaultPaymentMethodsFunc: func(ctx context.Context, id uuid.UUID) error {
			assert.Equal(t, integratorID, id)
			unsetCalled = true
			return nil
		},
		InsertPaymentMethodFunc: func(ctx context.Context, arg sqlc.InsertPaymentMethodParams) (sqlc.ConsolePaymentMethod, error) {
			assert.Equal(t, integratorID, arg.IntegratorID)
			assert.Equal(t, "token_123", arg.NombaTokenKey)
			assert.Equal(t, "1111", arg.CardLast4.String)
			assert.Equal(t, "Visa", arg.CardBrand.String)
			assert.True(t, arg.IsDefault)
			insertCalled = true
			return sqlc.ConsolePaymentMethod{}, nil
		},
	}

	recorder := billing.NewUsageRecorder(mq)
	eng := NewEngine(mq, nil, recorder, nil)

	payload := &nomba.WebhookPayload{EventType: "payment_success"}
	payload.Data.Transaction.Type = "online_checkout"
	payload.Data.Order.CustomerId = integratorID.String()
	payload.Data.TokenizedCardData = &nomba.TokenizedCardData{
		TokenKey: "token_123",
		CardType: "Visa",
		CardPan:  "4***45**** ****1111",
	}

	err := eng.ProcessWebhook(context.Background(), payload)
	assert.NoError(t, err)
	assert.True(t, unsetCalled)
	assert.True(t, insertCalled)
}
