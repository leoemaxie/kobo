package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Repository interface {
	GetActiveVirtualAccountByIdentityID(ctx context.Context, identityID uuid.UUID) (sqlc.VirtualAccount, error)
	ListLedgerEntriesByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int32) ([]LedgerEntry, error)
}

type repository struct {
	q *sqlc.Queries
}

func NewRepository(q *sqlc.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) GetActiveVirtualAccountByIdentityID(ctx context.Context, identityID uuid.UUID) (sqlc.VirtualAccount, error) {
	return r.q.GetActiveVirtualAccountByIdentityID(ctx, identityID)
}

func (r *repository) ListLedgerEntriesByAccount(ctx context.Context, accountID uuid.UUID, limit, offset int32) ([]LedgerEntry, error) {
	rows, err := r.q.ListLedgerEntriesByAccount(ctx, sqlc.ListLedgerEntriesByAccountParams{
		VirtualAccountID: accountID,
		Limit:            limit,
		Offset:           offset,
	})
	if err != nil {
		return nil, err
	}

	entries := make([]LedgerEntry, len(rows))
	for i, row := range rows {
		entries[i] = LedgerEntry{
			ID:               row.ID,
			VirtualAccountID: row.VirtualAccountID,
			IdentityID:       row.IdentityID,
			AmountKobo:       row.AmountKobo,
			Direction:        row.Direction,
			Status:           row.Status,
			MonnifyReference:   row.MonnifyReference,
			Source:           row.Source,
			Narration:        row.Narration.String,
			SenderName:       row.SenderName.String,
			OccurredAt:       row.OccurredAt,
			CreatedAt:        row.CreatedAt,
		}
	}
	return entries, nil
}
