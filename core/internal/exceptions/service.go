package exceptions

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

func (s *Service) ListOpen(ctx context.Context, integratorID uuid.UUID, limit, offset int32) ([]sqlc.Exception, error) {
	return s.q.ListOpenExceptions(ctx, sqlc.ListOpenExceptionsParams{
		IntegratorID: integratorID,
		Limit:        limit,
		Offset:       offset,
	})
}
