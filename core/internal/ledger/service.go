package ledger

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetStatements(ctx context.Context, identityID uuid.UUID, limit, offset int32) ([]LedgerEntry, error) {
	va, err := s.repo.GetActiveVirtualAccountByIdentityID(ctx, identityID)
	if err != nil {
		return nil, err
	}

	return s.repo.ListLedgerEntriesByAccount(ctx, va.ID, limit, offset)
}

func (s *Service) GetActiveVirtualAccount(ctx context.Context, identityID uuid.UUID) (sqlc.VirtualAccount, error) {
	return s.repo.GetActiveVirtualAccountByIdentityID(ctx, identityID)
}
