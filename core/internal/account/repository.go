package account

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Repository interface {
	CreateVirtualAccount(ctx context.Context, arg sqlc.CreateVirtualAccountParams) (sqlc.VirtualAccount, error)
	GetActiveVirtualAccountByIdentityID(ctx context.Context, identityID uuid.UUID) (sqlc.VirtualAccount, error)
	GetVirtualAccountByAccountNumber(ctx context.Context, accountNumber *string) (sqlc.VirtualAccount, error)
	UpdateVirtualAccountProvisioning(ctx context.Context, arg sqlc.UpdateVirtualAccountProvisioningParams) (sqlc.VirtualAccount, error)
	DeactivateVirtualAccount(ctx context.Context, identityID uuid.UUID) error

	// Lifecycle operations
	GetIdentityByID(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error)
	UpdateIdentityState(ctx context.Context, arg sqlc.UpdateIdentityStateParams) (sqlc.Identity, error)
	InsertIdentityEvent(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error)
}

type sqlcRepository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &sqlcRepository{q: q}
}

func (r *sqlcRepository) CreateVirtualAccount(ctx context.Context, arg sqlc.CreateVirtualAccountParams) (sqlc.VirtualAccount, error) {
	return r.q.CreateVirtualAccount(ctx, arg)
}

func (r *sqlcRepository) GetActiveVirtualAccountByIdentityID(ctx context.Context, identityID uuid.UUID) (sqlc.VirtualAccount, error) {
	return r.q.GetActiveVirtualAccountByIdentityID(ctx, identityID)
}

func (r *sqlcRepository) GetVirtualAccountByAccountNumber(ctx context.Context, accountNumber *string) (sqlc.VirtualAccount, error) {
	var pgAccountNumber pgtype.Text
	if accountNumber != nil {
		pgAccountNumber = pgtype.Text{String: *accountNumber, Valid: true}
	}
	return r.q.GetVirtualAccountByAccountNumber(ctx, pgAccountNumber)
}

func (r *sqlcRepository) UpdateVirtualAccountProvisioning(ctx context.Context, arg sqlc.UpdateVirtualAccountProvisioningParams) (sqlc.VirtualAccount, error) {
	return r.q.UpdateVirtualAccountProvisioning(ctx, arg)
}

func (r *sqlcRepository) DeactivateVirtualAccount(ctx context.Context, identityID uuid.UUID) error {
	return r.q.DeactivateVirtualAccount(ctx, identityID)
}

func (r *sqlcRepository) GetIdentityByID(ctx context.Context, arg sqlc.GetIdentityByIDParams) (sqlc.Identity, error) {
	return r.q.GetIdentityByID(ctx, arg)
}

func (r *sqlcRepository) UpdateIdentityState(ctx context.Context, arg sqlc.UpdateIdentityStateParams) (sqlc.Identity, error) {
	return r.q.UpdateIdentityState(ctx, arg)
}

func (r *sqlcRepository) InsertIdentityEvent(ctx context.Context, arg sqlc.InsertIdentityEventParams) (sqlc.IdentityEvent, error) {
	return r.q.InsertIdentityEvent(ctx, arg)
}
