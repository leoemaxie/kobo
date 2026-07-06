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
	Credential sqlc.ApiCredential
	RawSecret  string // Must be shown to the user exactly once
}

func (s *Service) ProvisionIntegrator(ctx context.Context, name string, isLive bool) (ProvisionResult, error) {
	apiKey, rawSecret, hashedSecret, err := auth.GenerateCredentials(isLive)
	if err != nil {
		return ProvisionResult{}, err
	}

	integratorID := uuid.New()

	integrator, err := s.q.CreateApiIntegrator(ctx, sqlc.CreateApiIntegratorParams{
		ID:   integratorID,
		Name: name,
	})
	if err != nil {
		return ProvisionResult{}, err
	}

	env := sqlc.ConsoleEnvironmentSandbox
	if isLive {
		env = sqlc.ConsoleEnvironmentProduction
	}

	// CreateApiCredential createdBy is nullable, so we pass an invalid/empty UUID.
	// We handle this correctly by passing the default uuid.UUID{} or sql.NullString/pgtype.UUID depending on sqlc settings.
	// Since sqlc usually maps nullable UUIDs to uuid.NullUUID, let's use that.
	cred, err := s.q.CreateApiCredential(ctx, sqlc.CreateApiCredentialParams{
		ID:           uuid.New(),
		IntegratorID: integratorID,
		Environment:  env,
		KeyID:        apiKey,
		SecretHash:   hashedSecret,
	})
	if err != nil {
		return ProvisionResult{}, err
	}

	return ProvisionResult{
		Integrator: integrator,
		Credential: cred,
		RawSecret:  rawSecret,
	}, nil
}
