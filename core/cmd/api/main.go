package main

import (
	"context"
	"log"
	"net/http"

	"github.com/leoemaxie/kobo/internal/account"
	"github.com/leoemaxie/kobo/internal/api"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/billing"
	"github.com/leoemaxie/kobo/internal/exceptions"
	"github.com/leoemaxie/kobo/internal/identity"
	"github.com/leoemaxie/kobo/internal/integrator"
	"github.com/leoemaxie/kobo/internal/ledger"
	"github.com/leoemaxie/kobo/internal/monnify"
	"github.com/leoemaxie/kobo/internal/payout"
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
	monnifyClient := monnify.NewClient(cfg.MonnifyBaseURL, cfg.MonnifyClientID, cfg.MonnifyClientSecret, cfg.MonnifyAccountID, cfg.MonnifySubAccountID, nil)
	accountSvc := account.NewService(accountRepo, monnifyClient)

	ledgerRepo := ledger.NewRepository(q)
	ledgerSvc := ledger.NewService(ledgerRepo)

	exceptionsRepo := exceptions.NewRepository(q)
	exceptionsSvc := exceptions.NewService(exceptionsRepo)
	integratorSvc := integrator.NewService(q)

	idemRepo := reconciliation.NewIdempotencyRepository(q)
	usageRecorder := billing.NewUsageRecorder(q)
	reconEngine := reconciliation.NewEngine(q, idemRepo, usageRecorder, monnifyClient)
	payoutSvc := payout.NewService(pool, q, monnifyClient)

	healthHandler := handlers.NewHealthHandler(pool)
	identityHandler := handlers.NewIdentityHandler(identitySvc, accountSvc, usageRecorder)
	ledgerHandler := handlers.NewLedgerHandler(ledgerSvc, monnifyClient)
	exceptionsHandler := handlers.NewExceptionsHandler(exceptionsSvc)
	adminHandler := handlers.NewAdminHandler(integratorSvc)
	adminBillingHandler := handlers.NewAdminBillingHandler(monnifyClient, q)
	payoutHandler := handlers.NewPayoutHandler(payoutSvc, monnifyClient)

	analyticsHandler := handlers.NewAnalyticsHandler(q)
	logsHandler := handlers.NewLogsHandler(q)

	router := api.NewRouter(q, healthHandler, identityHandler, ledgerHandler, exceptionsHandler, adminHandler, adminBillingHandler, payoutHandler, analyticsHandler, logsHandler, reconEngine, cfg.MonnifyWebhookSecret)

	log.Printf("Starting Kobo server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
