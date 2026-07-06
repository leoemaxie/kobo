package main

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/config"
	"github.com/leoemaxie/kobo/internal/platform/db"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	pool, err := db.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	q := sqlc.New(pool)
	ctx := context.Background()

	integratorID := uuid.New()
	_, err = q.CreateApiIntegrator(ctx, sqlc.CreateApiIntegratorParams{
		ID:   integratorID,
		Name: "Sandbox Test Integrator",
	})
	if err != nil {
		log.Fatalf("failed to seed integrator: %v", err)
	}

	_, err = q.CreateApiCredential(ctx, sqlc.CreateApiCredentialParams{
		ID:           uuid.New(),
		IntegratorID: integratorID,
		Environment:  sqlc.ConsoleEnvironmentSandbox,
		KeyID:        "kobo_test_seeded_api_key",
		SecretHash:   "seeded_secret_hash_not_real",
	})
	if err != nil {
		log.Fatalf("failed to seed integrator: %v", err)
	}

	log.Printf("Successfully seeded sandbox integrator with ID: %s", integratorID)
}
