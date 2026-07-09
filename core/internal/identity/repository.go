package identity

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Repository interface {
	CreateIdentity(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error)
	GetIdentityByID(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error)
	GetIdentityByExternalReference(ctx context.Context, arg sqlc.GetIdentityByExternalReferenceParams) (sqlc.Identity, error)
	UpdateIdentityProfile(ctx context.Context, arg sqlc.UpdateIdentityProfileParams) (sqlc.Identity, error)
	UpdateIdentityState(ctx context.Context, arg sqlc.UpdateIdentityStateParams) (sqlc.Identity, error)
	ListIdentitiesByState(ctx context.Context, arg sqlc.ListIdentitiesByStateParams) ([]sqlc.Identity, error)
	InsertIdentityEvent(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error)
	ListIdentityEvents(ctx context.Context, identityID uuid.UUID) ([]sqlc.IdentityEvent, error)
	GetActiveVirtualAccountByIdentityID(ctx context.Context, identityID uuid.UUID) (sqlc.VirtualAccount, error)
}

type sqlcRepository struct {
	q sqlc.Querier
}

func NewRepository(q sqlc.Querier) Repository {
	return &sqlcRepository{q: q}
}

func (r *sqlcRepository) CreateIdentity(ctx context.Context, arg sqlc.CreateIdentityParams) (sqlc.Identity, error) {
	return r.q.CreateIdentity(ctx, arg)
}

func (r *sqlcRepository) GetIdentityByID(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error) {
	return r.q.GetIdentityByID(ctx, arg)
}

func (r *sqlcRepository) GetIdentityByExternalReference(ctx context.Context, arg sqlc.GetIdentityByExternalReferenceParams) (sqlc.Identity, error) {
	return r.q.GetIdentityByExternalReference(ctx, arg)
}

func (r *sqlcRepository) UpdateIdentityProfile(ctx context.Context, arg sqlc.UpdateIdentityProfileParams) (sqlc.Identity, error) {
	return r.q.UpdateIdentityProfile(ctx, arg)
}

func (r *sqlcRepository) UpdateIdentityState(ctx context.Context, arg sqlc.UpdateIdentityStateParams) (sqlc.Identity, error) {
	return r.q.UpdateIdentityState(ctx, arg)
}

func (r *sqlcRepository) ListIdentitiesByState(ctx context.Context, arg sqlc.ListIdentitiesByStateParams) ([]sqlc.Identity, error) {
	return r.q.ListIdentitiesByState(ctx, arg)
}

func (r *sqlcRepository) InsertIdentityEvent(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error) {
	return r.q.InsertIdentityEvent(ctx, arg)
}

func (r *sqlcRepository) ListIdentityEvents(ctx context.Context, identityID uuid.UUID) ([]sqlc.IdentityEvent, error) {
	return r.q.ListIdentityEvents(ctx, identityID)
}

func (r *sqlcRepository) GetActiveVirtualAccountByIdentityID(ctx context.Context, identityID uuid.UUID) (sqlc.VirtualAccount, error) {
	return r.q.GetActiveVirtualAccountByIdentityID(ctx, identityID)
}
