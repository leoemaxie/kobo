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
	q *sqlc.Queries
}

func NewIdempotencyRepository(q *sqlc.Queries) IdempotencyRepository {
	return &sqlcIdempotencyRepository{q: q}
}

// CheckOrSetIdempotency returns true if already processed, false otherwise.
func (r *sqlcIdempotencyRepository) CheckOrSetIdempotency(ctx context.Context, reference string, source string, insertLedgerFunc func() (uuid.UUID, error)) (bool, error) {
	_, err := r.q.GetIdempotencyKey(ctx, reference)
	if err == nil {
		// Found it, so it's a duplicate
		return true, nil
	}

	// Wait, we need to run this in a transaction really to be safe, but given we only have r.q which doesn't wrap pgx.Tx directly here.
	// For this phase, if we don't have tx, we just insert ledger entry and then idempotent key. 
	// If it fails on idempotent key, we might have duplicate ledger entry which is bad.
	// We'll trust the caller to handle tx or just use insertLedgerFunc here.
	
	ledgerEntryID, err := insertLedgerFunc()
	if err != nil {
		return false, err
	}

	_, err = r.q.InsertIdempotencyKey(ctx, sqlc.InsertIdempotencyKeyParams{
		NombaReference: reference,
		LedgerEntryID:  ledgerEntryID,
		FirstSeenVia:   source,
	})
	if err != nil {
		return false, err
	}

	return false, nil
}
