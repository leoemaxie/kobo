package exceptions

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListOpen(ctx context.Context, integratorID uuid.UUID, limit, offset int32) ([]Exception, error) {
	return s.repo.ListOpenExceptions(ctx, integratorID, limit, offset)
}

func (s *Service) Resolve(ctx context.Context, id uuid.UUID, integratorID uuid.UUID, action string, notes string, successorIdentityID *uuid.UUID) error {
	_, err := s.repo.ResolveException(ctx, id, integratorID, action, notes, successorIdentityID)
	return err
}
