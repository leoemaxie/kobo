package integrator

import (
	"context"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/auth"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type Service struct {
	q sqlc.Querier
}

func NewService(q sqlc.Querier) *Service {
	return &Service{q: q}
}

type ProvisionResult struct {
	Integrator sqlc.ApiIntegrator
	RawSecret  string // Must be shown to the user exactly once
}

func (s *Service) ProvisionIntegrator(ctx context.Context, name string, isLive bool) (ProvisionResult, error) {
	apiKey, rawSecret, hashedSecret, err := auth.GenerateCredentials(isLive)
	if err != nil {
		return ProvisionResult{}, err
	}

	integratorID := uuid.New()

	integrator, err := s.q.CreateApiIntegrator(ctx, sqlc.CreateApiIntegratorParams{
		ID:            integratorID,
		Name:          name,
		ApiKey:        apiKey,
		ApiSecretHash: hashedSecret,
		IsSandbox:     !isLive,
	})
	if err != nil {
		return ProvisionResult{}, err
	}

	return ProvisionResult{
		Integrator: integrator,
		RawSecret:  rawSecret,
	}, nil
}
