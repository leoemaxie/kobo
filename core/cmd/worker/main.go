package main

import (
	"context"
	"log"
	"time"

	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/config"
	"github.com/leoemaxie/kobo/internal/platform/db"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/leoemaxie/kobo/internal/reconciliation"
)

func main() {
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
		nil,
	)

	sweeper := reconciliation.NewSweeper(q, idemRepo, nombaClient)

	sweepTicker := time.NewTicker(30 * time.Minute)
	defer sweepTicker.Stop()

	kycTicker := time.NewTicker(1 * time.Hour) // KYC background checks
	defer kycTicker.Stop()

	log.Println("Starting Kobo background worker...")

	// Initial run
	go func() {
		if err := sweeper.RunSweep(ctx); err != nil {
			log.Printf("Error running sweep: %v", err)
		}
	}()

	for {
		select {
		case <-sweepTicker.C:
			if err := sweeper.RunSweep(ctx); err != nil {
				log.Printf("Error running sweep: %v", err)
			}
		case <-kycTicker.C:
			// Run KYC check
			log.Println("Running KYC checks (placeholder)...")
		}
	}
}
