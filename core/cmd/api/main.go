package main

import (
	"context"
	"log"
	"net/http"

	"github.com/leoemaxie/kobo/internal/api"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/account"
	"github.com/leoemaxie/kobo/internal/exceptions"
	"github.com/leoemaxie/kobo/internal/identity"
	"github.com/leoemaxie/kobo/internal/integrator"
	"github.com/leoemaxie/kobo/internal/ledger"
	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/config"
	"github.com/leoemaxie/kobo/internal/platform/db"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/leoemaxie/kobo/internal/platform/telemetry"
	"github.com/leoemaxie/kobo/internal/reconciliation"
)

func main() {
	telemetry.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	q := sqlc.New(pool)
	identityRepo := identity.NewRepository(q)
	identitySvc := identity.NewService(identityRepo)

	accountRepo := account.NewRepository(q)
	nombaClient := nomba.NewClient(cfg.NombaBaseURL, cfg.NombaClientID, cfg.NombaClientSecret, cfg.NombaAccountID, nil)
	accountSvc := account.NewService(accountRepo, nombaClient)

	ledgerRepo := ledger.NewRepository(q)
	ledgerSvc := ledger.NewService(ledgerRepo)
	
	exceptionsRepo := exceptions.NewRepository(q)
	exceptionsSvc := exceptions.NewService(exceptionsRepo)
	integratorSvc := integrator.NewService(q)

	identityHandler := handlers.NewIdentityHandler(identitySvc, accountSvc)
	ledgerHandler := handlers.NewLedgerHandler(ledgerSvc)
	exceptionsHandler := handlers.NewExceptionsHandler(exceptionsSvc)
	adminHandler := handlers.NewAdminHandler(integratorSvc)
	
	idemRepo := reconciliation.NewIdempotencyRepository(q)
	reconEngine := reconciliation.NewEngine(q, idemRepo)

	router := api.NewRouter(q, identityHandler, ledgerHandler, exceptionsHandler, adminHandler, reconEngine, cfg.NombaWebhookSecret)

	log.Printf("Starting Kobo server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
