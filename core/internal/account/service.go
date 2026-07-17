package account

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type MonnifyClient interface {
	CreateVirtualAccount(ctx context.Context, accountRef, accountName, bvn string) (MonnifyAccountResponse, error)
	ExpireVirtualAccount(ctx context.Context, identifier string) (bool, error)
}

type MonnifyAccountResponse struct {
	AccountNumber   string
	BankName        string
	BankAccountName string
}

type Service struct {
	repo        Repository
	monnifyClient MonnifyClient
}

func NewService(repo Repository, monnifyClient MonnifyClient) *Service {
	return &Service{repo: repo, monnifyClient: monnifyClient}
}

// Provision handles the PENDING -> ACTIVE transition
func (s *Service) Provision(ctx context.Context, identityID, integratorID uuid.UUID) error {
	// 1. Fetch Identity
	ident, err := s.repo.GetIdentityByID(ctx, sqlc.GetIdentityByIDParams{
		ID:           identityID,
		IntegratorID: integratorID,
	})
	if err != nil {
		return fmt.Errorf("identity not found: %w", err)
	}

	// 2. State transition check
	newState, err := ValidTransition(State(ident.State), EventProvisionSuccess)
	if err != nil {
		return fmt.Errorf("invalid state transition: %w", err)
	}

	// 3. Deactivate old accounts if reprovisioning
	_ = s.repo.DeactivateVirtualAccount(ctx, identityID)

	// 4. Generate new accountRef
	accountRef := uuid.New().String()

	// 5. Create PENDING virtual account in DB
	va, err := s.repo.CreateVirtualAccount(ctx, sqlc.CreateVirtualAccountParams{
		ID:              uuid.New(),
		IdentityID:      identityID,
		MonnifyAccountRef: accountRef,
		IsActive:        true,
	})
	if err != nil {
		return fmt.Errorf("failed to create virtual account record: %w", err)
	}

	// 6. Call Monnify API
	resp, monnifyErr := s.monnifyClient.CreateVirtualAccount(ctx, accountRef, ident.DisplayName, "")

	if monnifyErr != nil {
		return fmt.Errorf("monnify provisioning failed: %w", monnifyErr)
	}

	// 7. Update Virtual Account with Monnify details
	var acctNum pgtype.Text
	acctNum.String = resp.AccountNumber
	acctNum.Valid = true

	var bankName pgtype.Text
	bankName.String = resp.BankName
	bankName.Valid = true

	var bankAcctName pgtype.Text
	bankAcctName.String = resp.BankAccountName
	bankAcctName.Valid = true

	_, err = s.repo.UpdateVirtualAccountProvisioning(ctx, sqlc.UpdateVirtualAccountProvisioningParams{
		ID:            va.ID,
		AccountNumber: acctNum,
		BankName:      bankName,
		AccountName:   bankAcctName,
	})
	if err != nil {
		return fmt.Errorf("failed to update virtual account with provisioned details: %w", err)
	}

	// 8. Update Identity State
	_, err = s.repo.UpdateIdentityState(ctx, sqlc.UpdateIdentityStateParams{
		ID:            identityID,
		IntegratorID:  integratorID,
		State:         string(newState),
		FailureReason: pgtype.Text{Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to update identity state: %w", err)
	}

	_, _ = s.repo.InsertIdentityEvent(ctx, sqlc.InsertIdentityEventParams{
		ID:            uuid.New(),
		IdentityID:    identityID,
		EventType:     "activated",
		PreviousState: pgtype.Text{String: ident.State, Valid: true},
		NewState:      pgtype.Text{String: string(newState), Valid: true},
		Detail:        []byte("{}"),
	})

	return nil
}

func (s *Service) Reopen(ctx context.Context, identityID, integratorID uuid.UUID) error {
	ident, err := s.repo.GetIdentityByID(ctx, sqlc.GetIdentityByIDParams{
		ID:           identityID,
		IntegratorID: integratorID,
	})
	if err != nil {
		return fmt.Errorf("identity not found: %w", err)
	}

	newState, err := ValidTransition(State(ident.State), EventReopenInitiated)
	if err != nil {
		return fmt.Errorf("invalid state transition: %w", err)
	}

	// For reopen, we transition immediately to the new state (Active) after provisioning.
	// We'll call Provision again, but since it's CLOSED, Provision would reject PENDING->ACTIVE.
	// Actually, wait, ValidTransition(CLOSED, EventReopenInitiated) returns ACTIVE.
	// So we can do the provisioning steps here.

	_ = s.repo.DeactivateVirtualAccount(ctx, identityID)
	accountRef := uuid.New().String()

	va, err := s.repo.CreateVirtualAccount(ctx, sqlc.CreateVirtualAccountParams{
		ID:              uuid.New(),
		IdentityID:      identityID,
		MonnifyAccountRef: accountRef,
		IsActive:        true,
	})
	if err != nil {
		return fmt.Errorf("failed to create virtual account record: %w", err)
	}

	resp, monnifyErr := s.monnifyClient.CreateVirtualAccount(ctx, accountRef, ident.DisplayName, "")
	if monnifyErr != nil {
		return fmt.Errorf("monnify provisioning failed on reopen: %w", monnifyErr)
	}

	var acctNum pgtype.Text
	acctNum.String = resp.AccountNumber
	acctNum.Valid = true

	var bankName pgtype.Text
	bankName.String = resp.BankName
	bankName.Valid = true

	var bankAcctName pgtype.Text
	bankAcctName.String = resp.BankAccountName
	bankAcctName.Valid = true

	_, err = s.repo.UpdateVirtualAccountProvisioning(ctx, sqlc.UpdateVirtualAccountProvisioningParams{
		ID:            va.ID,
		AccountNumber: acctNum,
		BankName:      bankName,
		AccountName:   bankAcctName,
	})
	if err != nil {
		return fmt.Errorf("failed to update virtual account: %w", err)
	}

	_, err = s.repo.UpdateIdentityState(ctx, sqlc.UpdateIdentityStateParams{
		ID:            identityID,
		IntegratorID:  integratorID,
		State:         string(newState),
		FailureReason: pgtype.Text{Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to update identity state: %w", err)
	}

	_, _ = s.repo.InsertIdentityEvent(ctx, sqlc.InsertIdentityEventParams{
		ID:            uuid.New(),
		IdentityID:    identityID,
		EventType:     "reopened",
		PreviousState: pgtype.Text{String: ident.State, Valid: true},
		NewState:      pgtype.Text{String: string(newState), Valid: true},
		Detail:        []byte("{}"),
	})

	return nil
}

func (s *Service) Close(ctx context.Context, identityID, integratorID uuid.UUID, reason string, sweepDestination []byte) error {
	ident, err := s.repo.GetIdentityByID(ctx, sqlc.GetIdentityByIDParams{
		ID:           identityID,
		IntegratorID: integratorID,
	})
	if err != nil {
		return fmt.Errorf("identity not found: %w", err)
	}

	newState, err := ValidTransition(State(ident.State), EventCloseInitiated)
	if err != nil {
		return fmt.Errorf("invalid state transition: %w", err)
	}

	// Move to CLOSING
	_, err = s.repo.UpdateIdentityState(ctx, sqlc.UpdateIdentityStateParams{
		ID:            identityID,
		IntegratorID:  integratorID,
		State:         string(newState),
		FailureReason: pgtype.Text{Valid: false},
	})
	if err != nil {
		return fmt.Errorf("failed to update identity state: %w", err)
	}

	// Fetch active virtual account and expire it
	va, vaErr := s.repo.GetActiveVirtualAccountByIdentityID(ctx, identityID)
	if vaErr == nil && va.MonnifyAccountRef != "" {
		_, _ = s.monnifyClient.ExpireVirtualAccount(ctx, va.MonnifyAccountRef)
	}

	if len(sweepDestination) == 0 {
		sweepDestination = []byte(`{}`)
	}

	detail := fmt.Sprintf(`{"reason": "%s", "sweep_destination": %s}`, reason, string(sweepDestination))

	_, _ = s.repo.InsertIdentityEvent(ctx, sqlc.InsertIdentityEventParams{
		ID:            uuid.New(),
		IdentityID:    identityID,
		EventType:     "closing_started",
		PreviousState: pgtype.Text{String: ident.State, Valid: true},
		NewState:      pgtype.Text{String: string(newState), Valid: true},
		Detail:        []byte(detail),
	})

	return nil
}
