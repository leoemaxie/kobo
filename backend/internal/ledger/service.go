package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Service struct {
	q *sqlc.Queries
}

func NewService(q *sqlc.Queries) *Service {
	return &Service{q: q}
}

func (s *Service) GetStatements(ctx context.Context, identityID uuid.UUID, limit, offset int32) ([]sqlc.LedgerEntry, error) {
	// Technically the query requires a virtual account ID, or identity ID.
	// Let's use ListLedgerEntriesByIdentityAndPeriod if we had it, but we only have:
	// ListLedgerEntriesByAccount (needs VA id)
	
	// We should just get the active virtual account and list entries.
	va, err := s.q.GetActiveVirtualAccountByIdentityID(ctx, identityID)
	if err != nil {
		return nil, err
	}

	return s.q.ListLedgerEntriesByAccount(ctx, sqlc.ListLedgerEntriesByAccountParams{
		VirtualAccountID: va.ID,
		Limit:            limit,
		Offset:           offset,
	})
}
