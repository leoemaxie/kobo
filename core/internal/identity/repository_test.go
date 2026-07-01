package identity

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/stretchr/testify/assert"
)

type mockQuerier struct {
	sqlc.Querier
	CreateIdentityFunc                 func(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error)
	GetIdentityByIDFunc                func(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error)
	GetIdentityByExternalReferenceFunc func(ctx context.Context, arg sqlc.GetIdentityByExternalReferenceParams) (sqlc.Identity, error)
	UpdateIdentityProfileFunc          func(ctx context.Context, arg sqlc.UpdateIdentityProfileParams) (sqlc.Identity, error)
}

func (m *mockQuerier) CreateIdentity(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error) {
	if m.CreateIdentityFunc != nil {
		return m.CreateIdentityFunc(ctx, arg)
	}
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockQuerier) GetIdentityByID(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error) {
	if m.GetIdentityByIDFunc != nil {
		return m.GetIdentityByIDFunc(ctx, arg)
	}
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockQuerier) GetIdentityByExternalReference(ctx context.Context, arg sqlc.GetIdentityByExternalReferenceParams) (sqlc.Identity, error) {
	if m.GetIdentityByExternalReferenceFunc != nil {
		return m.GetIdentityByExternalReferenceFunc(ctx, arg)
	}
	return sqlc.Identity{}, errors.New("unimplemented")
}
func (m *mockQuerier) UpdateIdentityProfile(ctx context.Context, arg sqlc.UpdateIdentityProfileParams) (sqlc.Identity, error) {
	if m.UpdateIdentityProfileFunc != nil {
		return m.UpdateIdentityProfileFunc(ctx, arg)
	}
	return sqlc.Identity{}, errors.New("unimplemented")
}

func TestRepository_CreateIdentity(t *testing.T) {
	mq := &mockQuerier{
		CreateIdentityFunc: func(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error) {
			assert.Equal(t, "ext-123", arg.ExternalReference)
			return sqlc.Identity{
				ID:                arg.ID,
				ExternalReference: arg.ExternalReference,
			}, nil
		},
	}
	repo := NewRepository(mq)

	id := uuid.New()
	ident, err := repo.CreateIdentity(context.Background(), sqlc.CreateIdentityParams{
		ID:                id,
		ExternalReference: "ext-123",
	})
	assert.NoError(t, err)
	assert.Equal(t, id, ident.ID)
	assert.Equal(t, "ext-123", ident.ExternalReference)
}

func TestRepository_GetIdentityByID(t *testing.T) {
	identID := uuid.New()
	integratorID := uuid.New()

	mq := &mockQuerier{
		GetIdentityByIDFunc: func(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error) {
			assert.Equal(t, identID, arg.ID)
			assert.Equal(t, integratorID, arg.IntegratorID)
			return sqlc.Identity{
				ID:           arg.ID,
				IntegratorID: arg.IntegratorID,
			}, nil
		},
	}
	repo := NewRepository(mq)

	ident, err := repo.GetIdentityByID(context.Background(), sqlc.GetIdentityByIDParams{
		ID:           identID,
		IntegratorID: integratorID,
	})
	assert.NoError(t, err)
	assert.Equal(t, identID, ident.ID)
	assert.Equal(t, integratorID, ident.IntegratorID)
}
