package reconciliation

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type IdempotencyRepository interface {
	CheckOrSetIdempotency(ctx context.Context, reference string, source string, insertLedgerFunc func() (uuid.UUID, error)) (bool, error)
}

type sqlcIdempotencyRepository struct {
	q sqlc.Querier
}

func NewIdempotencyRepository(q sqlc.Querier) IdempotencyRepository {
	return &sqlcIdempotencyRepository{q: q}
}

// CheckOrSetIdempotency returns true if already processed, false otherwise.
func (r *sqlcIdempotencyRepository) CheckOrSetIdempotency(ctx context.Context, reference string, source string, insertLedgerFunc func() (uuid.UUID, error)) (bool, error) {
	_, err := r.q.GetIdempotencyKey(ctx, reference)
	if err == nil {
		// Found it, so it's a duplicate
		return true, nil
	}

	ledgerEntryID, err := insertLedgerFunc()
	if err != nil {
		return false, err
	}

	_, err = r.q.InsertIdempotencyKey(ctx, sqlc.InsertIdempotencyKeyParams{
		MonnifyReference: reference,
		LedgerEntryID:  ledgerEntryID,
		FirstSeenVia:   source,
	})
	if err != nil {
		return false, err
	}

	return false, nil
}
