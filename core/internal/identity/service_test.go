package identity

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/stretchr/testify/assert"
)

type mockRepository struct {
	CreateIdentityFunc                 func(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error)
	GetIdentityByIDFunc                func(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error)
	GetIdentityByExternalReferenceFunc func(ctx context.Context, arg sqlc.GetIdentityByExternalReferenceParams) (sqlc.Identity, error)
	UpdateIdentityProfileFunc          func(ctx context.Context, arg sqlc.UpdateIdentityProfileParams) (sqlc.Identity, error)
	UpdateIdentityStateFunc            func(ctx context.Context, arg sqlc.UpdateIdentityStateParams) (sqlc.Identity, error)
	ListIdentitiesByStateFunc          func(ctx context.Context, arg sqlc.ListIdentitiesByStateParams) ([]sqlc.Identity, error)
	InsertIdentityEventFunc            func(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error)
	ListIdentityEventsFunc             func(ctx context.Context, identityID uuid.UUID) ([]sqlc.IdentityEvent, error)
}

func (m *mockRepository) CreateIdentity(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error) {
	if m.CreateIdentityFunc != nil {
		return m.CreateIdentityFunc(ctx, arg)
	}
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockRepository) GetIdentityByID(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error) {
	if m.GetIdentityByIDFunc != nil {
		return m.GetIdentityByIDFunc(ctx, arg)
	}
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockRepository) GetIdentityByExternalReference(ctx context.Context, arg sqlc.GetIdentityByExternalReferenceParams) (sqlc.Identity, error) {
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockRepository) UpdateIdentityProfile(ctx context.Context, arg sqlc.UpdateIdentityProfileParams) (sqlc.Identity, error) {
	if m.UpdateIdentityProfileFunc != nil {
		return m.UpdateIdentityProfileFunc(ctx, arg)
	}
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockRepository) UpdateIdentityState(ctx context.Context, arg sqlc.UpdateIdentityStateParams) (sqlc.Identity, error) {
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockRepository) ListIdentitiesByState(ctx context.Context, arg sqlc.ListIdentitiesByStateParams) ([]sqlc.Identity, error) {
	return nil, errors.New("unimplemented")
}
func (m *mockRepository) InsertIdentityEvent(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error) {
	if m.InsertIdentityEventFunc != nil {
		return m.InsertIdentityEventFunc(ctx, arg)
	}
	return sqlc.IdentityEvent{}, nil // Default to no-op
}
func (m *mockRepository) ListIdentityEvents(ctx context.Context, identityID uuid.UUID) ([]sqlc.IdentityEvent, error) {
	return nil, errors.New("unimplemented")
}

func TestService_Register(t *testing.T) {
	repo := &mockRepository{
		CreateIdentityFunc: func(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error) {
			assert.Equal(t, "ext-123", arg.ExternalReference)
			assert.Equal(t, "John Doe", arg.DisplayName)
			return sqlc.Identity{
				ID:                arg.ID,
				IntegratorID:      arg.IntegratorID,
				ExternalReference: arg.ExternalReference,
				DisplayName:       arg.DisplayName,
				Metadata:          arg.Metadata,
				State:             "pending",
			}, nil
		},
		InsertIdentityEventFunc: func(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error) {
			assert.Equal(t, "created", arg.EventType)
			return sqlc.IdentityEvent{}, nil
		},
	}

	svc := NewService(repo)

	integratorID := uuid.New()
	ident, err := svc.Register(context.Background(), integratorID, "ext-123", "John Doe", json.RawMessage(`{"age": 30}`))
	assert.NoError(t, err)
	assert.NotNil(t, ident)
	assert.Equal(t, "pending", ident.State)
	assert.Equal(t, "ext-123", ident.ExternalReference)
	assert.Equal(t, "John Doe", ident.DisplayName)
}

func TestService_Get(t *testing.T) {
	identID := uuid.New()
	integratorID := uuid.New()

	repo := &mockRepository{
		GetIdentityByIDFunc: func(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error) {
			assert.Equal(t, identID, arg.ID)
			assert.Equal(t, integratorID, arg.IntegratorID)
			return sqlc.Identity{
				ID:                arg.ID,
				IntegratorID:      arg.IntegratorID,
				ExternalReference: "ext-123",
				DisplayName:       "John Doe",
				State:             "active",
			}, nil
		},
	}

	svc := NewService(repo)

	ident, err := svc.Get(context.Background(), identID, integratorID)
	assert.NoError(t, err)
	assert.NotNil(t, ident)
	assert.Equal(t, "active", ident.State)
	assert.Equal(t, "ext-123", ident.ExternalReference)
}

func TestService_UpdateProfile(t *testing.T) {
	identID := uuid.New()
	integratorID := uuid.New()

	repo := &mockRepository{
		UpdateIdentityProfileFunc: func(ctx context.Context, arg sqlc.UpdateIdentityProfileParams) (sqlc.Identity, error) {
			assert.Equal(t, identID, arg.ID)
			assert.Equal(t, integratorID, arg.IntegratorID)
			assert.Equal(t, "Jane Doe", arg.DisplayName.String)
			return sqlc.Identity{
				ID:           arg.ID,
				IntegratorID: arg.IntegratorID,
				DisplayName:  arg.DisplayName.String,
				Metadata:     arg.Metadata,
				State:        "active",
			}, nil
		},
		InsertIdentityEventFunc: func(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error) {
			assert.Equal(t, "renamed", arg.EventType)
			return sqlc.IdentityEvent{}, nil
		},
	}

	svc := NewService(repo)

	newName := "Jane Doe"
	ident, err := svc.UpdateProfile(context.Background(), identID, integratorID, &newName, json.RawMessage(`{"foo": "bar"}`))
	assert.NoError(t, err)
	assert.NotNil(t, ident)
	assert.Equal(t, "Jane Doe", ident.DisplayName)
}
