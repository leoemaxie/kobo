package main

import (
	"context"
	"log"
	"github.com/leoemaxie/kobo/internal/billing"
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
	idemRepo := reconciliation.NewIdempotencyRepository(q)

	nombaClient := nomba.NewClient(
		cfg.NombaBaseURL,
		cfg.NombaClientID,
		cfg.NombaClientSecret,
		cfg.NombaAccountID,
		cfg.NombaSubAccountID,
		nil,
	)

	sweeper := reconciliation.NewSweeper(q, idemRepo, nombaClient)
	closureSweeper := reconciliation.NewClosureSweeper(q)
	invoiceJob := billing.NewInvoiceJob(q, nombaClient)

	log.Println("Starting Kobo one-off background sweep...")

	if err := sweeper.RunSweep(ctx); err != nil {
		log.Printf("Error running sweep: %v", err)
	}
	
	if err := closureSweeper.Run(ctx); err != nil {
		log.Printf("Error running closure sweep: %v", err)
	}
	
	if err := invoiceJob.Run(ctx); err != nil {
		log.Printf("Error running invoice job: %v", err)
	}

	log.Println("Running KYC checks (placeholder)...")

	log.Println("Sweep completed successfully.")
}
