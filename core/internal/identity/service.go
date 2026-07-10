package identity

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

var ErrIdentityConflict = errors.New("identity with external reference already exists")

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, integratorID uuid.UUID, externalRef, displayName string, metadata json.RawMessage) (*Identity, error) {
	id := uuid.New()

	metaBytes, err := json.Marshal(metadata)
	if err != nil || len(metaBytes) == 0 || string(metaBytes) == "null" {
		metaBytes = []byte("{}")
	}

	arg := sqlc.CreateIdentityParams{
		ID:                id,
		IntegratorID:      integratorID,
		ExternalReference: externalRef,
		DisplayName:       displayName,
		Metadata:          metaBytes,
	}

	row, err := s.repo.CreateIdentity(ctx, arg)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrIdentityConflict
		}
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	_, err = s.repo.InsertIdentityEvent(ctx, sqlc.InsertIdentityEventParams{
		ID:            uuid.New(),
		IdentityID:    id,
		EventType:     "created",
		PreviousState: pgtype.Text{Valid: false},
		NewState:      pgtype.Text{String: "pending", Valid: true},
		Detail:        []byte("{}"),
	})
	if err != nil {
		// Just logging would go here in a real app
	}

	ident := s.mapSQLCToIdentity(row)
	s.populateVirtualAccount(ctx, ident)
	return ident, nil
}

func (s *Service) Delete(ctx context.Context, id, integratorID uuid.UUID) error {
	return s.repo.DeleteIdentityCascade(ctx, sqlc.DeleteIdentityCascadeParams{
		ID:           id,
		IntegratorID: integratorID,
	})
}

func (s *Service) Get(ctx context.Context, id, integratorID uuid.UUID) (*Identity, error) {
	row, err := s.repo.GetIdentityByID(ctx, sqlc.GetIdentityByIDParams{
		ID:           id,
		IntegratorID: integratorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}
	
	ident := s.mapSQLCToIdentity(row)
	s.populateVirtualAccount(ctx, ident)
	return ident, nil
}

func (s *Service) List(ctx context.Context, integratorID uuid.UUID, state *string, limit, offset int32) ([]*Identity, error) {
	var stateParam pgtype.Text
	if state != nil {
		stateParam = pgtype.Text{String: *state, Valid: true}
	}
	
	rows, err := s.repo.ListIdentities(ctx, sqlc.ListIdentitiesParams{
		IntegratorID: integratorID,
		State:        stateParam,
		Limit:        limit,
		Offset:       offset,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list identities: %w", err)
	}

	identities := make([]*Identity, 0, len(rows))
	for _, row := range rows {
		ident := s.mapSQLCToIdentity(row)
		s.populateVirtualAccount(ctx, ident)
		identities = append(identities, ident)
	}
	return identities, nil
}

func (s *Service) UpdateProfile(ctx context.Context, id, integratorID uuid.UUID, displayName *string, metadata json.RawMessage) (*Identity, error) {
	var metaBytes []byte
	var err error
	if metadata != nil {
		metaBytes, err = json.Marshal(metadata)
		if err != nil {
			return nil, fmt.Errorf("invalid metadata")
		}
	}

	var dName pgtype.Text
	if displayName != nil {
		dName = pgtype.Text{String: *displayName, Valid: true}
	}

	row, err := s.repo.UpdateIdentityProfile(ctx, sqlc.UpdateIdentityProfileParams{
		ID:           id,
		IntegratorID: integratorID,
		DisplayName:  dName,
		Metadata:     metaBytes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update identity: %w", err)
	}

	_, _ = s.repo.InsertIdentityEvent(ctx, sqlc.InsertIdentityEventParams{
		ID:            uuid.New(),
		IdentityID:    id,
		EventType:     "renamed",
		PreviousState: pgtype.Text{Valid: false},
		NewState:      pgtype.Text{Valid: false},
		Detail:        []byte("{}"), // Ideally would store diff here
	})

	ident := s.mapSQLCToIdentity(row)
	s.populateVirtualAccount(ctx, ident)
	return ident, nil
}

func (s *Service) populateVirtualAccount(ctx context.Context, ident *Identity) {
	if ident.State == "pending" || ident.State == "failed" {
		return
	}

	va, err := s.repo.GetActiveVirtualAccountByIdentityID(ctx, ident.ID)
	if err != nil {
		// Log error if needed, but return normally
		return
	}

	if !va.AccountNumber.Valid || !va.BankName.Valid || !va.AccountName.Valid {
		return
	}

	var expectedAmountKobo *int64
	if va.ExpectedAmountKobo.Valid {
		expectedAmountKobo = &va.ExpectedAmountKobo.Int64
	}

	ident.VirtualAccount = &VirtualAccountSummary{
		AccountNumber:      va.AccountNumber.String,
		BankName:           va.BankName.String,
		AccountName:        va.AccountName.String,
		ExpectedAmountKobo: expectedAmountKobo,
		IsExpired:          va.IsExpired,
	}
}

func (s *Service) mapSQLCToIdentity(row sqlc.Identity) *Identity {
	return &Identity{
		ID:                row.ID,
		IntegratorID:      row.IntegratorID,
		ExternalReference: row.ExternalReference,
		DisplayName:       row.DisplayName,
		Metadata:          row.Metadata,
		State:             row.State,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}
