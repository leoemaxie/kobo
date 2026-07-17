package billing

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

// MockQuerier is a simple mock implementation of sqlc.Querier for testing
type MockQuerier struct {
	sqlc.Querier
	WalletBalance   int64
	UpdatedBalance  int64
	InvoiceStatus   string
	Suspended       bool
	PendingInvoices []sqlc.Invoice
	DefaultPayment  sqlc.PaymentMethod
	SuspendCalled   bool
	BalanceUpdated  bool
	StatusUpdated   bool
	GenerateCalled  bool
	GenerateError   error
	PendingError    error
	WalletError     error
	PaymentError    error
}

func (m *MockQuerier) RollupUsageEvents(ctx context.Context, period string) error {
	return nil
}

func (m *MockQuerier) GenerateInvoicesForPeriod(ctx context.Context, period string) error {
	m.GenerateCalled = true
	return m.GenerateError
}

func (m *MockQuerier) GetPendingInvoices(ctx context.Context) ([]sqlc.Invoice, error) {
	return m.PendingInvoices, m.PendingError
}

func (m *MockQuerier) GetIntegratorWalletBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	return m.WalletBalance, m.WalletError
}

func (m *MockQuerier) UpdateIntegratorWalletBalance(ctx context.Context, arg sqlc.UpdateIntegratorWalletBalanceParams) error {
	m.BalanceUpdated = true
	m.UpdatedBalance = m.WalletBalance + arg.WalletBalanceKobo
	return nil
}

func (m *MockQuerier) GetDefaultPaymentMethod(ctx context.Context, integratorID uuid.UUID) (sqlc.PaymentMethod, error) {
	return m.DefaultPayment, m.PaymentError
}

func (m *MockQuerier) UpdateInvoiceStatus(ctx context.Context, arg sqlc.UpdateInvoiceStatusParams) error {
	m.StatusUpdated = true
	m.InvoiceStatus = arg.Status
	return nil
}

func (m *MockQuerier) SuspendIntegrator(ctx context.Context, integratorID uuid.UUID) error {
	m.SuspendCalled = true
	m.Suspended = true
	return nil
}

// MockMonnifyClient wraps methods if necessary, though testing against actual HTTP client usually requires httptest.
// We'll skip monnify client execution testing here and just test the DB logic.

func TestInvoiceJob_WalletDeduction_FullCoverage(t *testing.T) {
	mockQ := &MockQuerier{
		WalletBalance: 50000, // ₦500
		PendingInvoices: []sqlc.Invoice{
			{
				ID:           uuid.New(),
				IntegratorID: uuid.New(),
				AmountKobo:   40000, // ₦400 invoice
				Status:       "open",
			},
		},
	}

	job := NewInvoiceJob(mockQ, nil)

	err := job.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockQ.BalanceUpdated {
		t.Errorf("expected wallet balance to be updated")
	}

	if mockQ.UpdatedBalance != 10000 {
		t.Errorf("expected remaining balance to be 10000, got %d", mockQ.UpdatedBalance)
	}

	if mockQ.InvoiceStatus != "paid" {
		t.Errorf("expected invoice status to be paid, got %s", mockQ.InvoiceStatus)
	}
}

func TestInvoiceJob_WalletDeduction_PartialCoverage(t *testing.T) {
	mockQ := &MockQuerier{
		WalletBalance: 20000, // ₦200
		PendingInvoices: []sqlc.Invoice{
			{
				ID:           uuid.New(),
				IntegratorID: uuid.New(),
				AmountKobo:   50000, // ₦500 invoice
				Status:       "open",
			},
		},
		DefaultPayment: sqlc.PaymentMethod{
			MonnifyTokenKey: "", // Simulating missing token to trigger failInvoice
		},
	}

	job := NewInvoiceJob(mockQ, nil)

	err := job.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockQ.BalanceUpdated {
		t.Errorf("expected wallet balance to be updated")
	}

	if mockQ.UpdatedBalance != 0 {
		t.Errorf("expected remaining balance to be 0, got %d", mockQ.UpdatedBalance)
	}

	if mockQ.InvoiceStatus != "failed" {
		t.Errorf("expected invoice status to be failed (due to no payment method), got %s", mockQ.InvoiceStatus)
	}
}

func TestInvoiceJob_RetrySuspension(t *testing.T) {
	mockQ := &MockQuerier{
		WalletBalance: 0,
		PendingInvoices: []sqlc.Invoice{
			{
				ID:           uuid.New(),
				IntegratorID: uuid.New(),
				AmountKobo:   50000,
				Status:       "failed",
				RetryCount:   2, // This will be the 3rd retry
			},
		},
		DefaultPayment: sqlc.PaymentMethod{
			MonnifyTokenKey: "", // Force failure
		},
	}

	job := NewInvoiceJob(mockQ, nil)

	err := job.Run(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockQ.SuspendCalled {
		t.Errorf("expected integrator to be suspended after 3 failed retries")
	}
}
