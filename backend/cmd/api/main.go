package main

import (
	"context"
	"log"
	"net/http"

	"github.com/leoemaxie/kobo/internal/api"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/identity"
	"github.com/leoemaxie/kobo/internal/platform/config"
	"github.com/leoemaxie/kobo/internal/platform/db"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
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
	identityRepo := identity.NewRepository(q)
	identitySvc := identity.NewService(identityRepo)

	identityHandler := handlers.NewIdentityHandler(identitySvc)
	router := api.NewRouter(q, identityHandler)

	log.Printf("Starting Kobo server on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
