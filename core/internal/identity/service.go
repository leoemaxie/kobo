package identity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, integratorID uuid.UUID, externalRef, displayName, kycTier string, metadata json.RawMessage) (*Identity, error) {
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
		KycTier:           kycTier,
		Metadata:          metaBytes,
	}

	row, err := s.repo.CreateIdentity(ctx, arg)
	if err != nil {
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

	return s.mapSQLCToIdentity(row), nil
}

func (s *Service) Get(ctx context.Context, id, integratorID uuid.UUID) (*Identity, error) {
	row, err := s.repo.GetIdentityByID(ctx, sqlc.GetIdentityByIDParams{
		ID:           id,
		IntegratorID: integratorID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get identity: %w", err)
	}
	return s.mapSQLCToIdentity(row), nil
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

	return s.mapSQLCToIdentity(row), nil
}

func (s *Service) mapSQLCToIdentity(row sqlc.Identity) *Identity {
	var failureReason *string
	if row.FailureReason.Valid {
		failureReason = &row.FailureReason.String
	}
	return &Identity{
		ID:                row.ID,
		IntegratorID:      row.IntegratorID,
		ExternalReference: row.ExternalReference,
		Profile: DisplayProfile{
			DisplayName: row.DisplayName,
			Metadata:    row.Metadata,
		},
		KYCTier:       row.KycTier,
		State:         row.State,
		FailureReason: failureReason,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
