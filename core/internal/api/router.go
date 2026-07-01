package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/leoemaxie/kobo/internal/reconciliation"
)

func NewRouter(q *sqlc.Queries, identityHandler *handlers.IdentityHandler, ledgerHandler *handlers.LedgerHandler, exceptionsHandler *handlers.ExceptionsHandler, adminHandler *handlers.AdminHandler, engine reconciliation.Engine, webhookSecret string) *chi.Mux {
	r := chi.NewRouter()

	// Admin routes
	r.Post("/admin/integrators", adminHandler.ProvisionIntegrator)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(q))

		r.Post("/identities", identityHandler.Create)
		r.Get("/identities/{id}", identityHandler.Get)
		r.Get("/identities/{id}/statements", ledgerHandler.GetStatements)
		r.Get("/exceptions", exceptionsHandler.ListOpen)
	})

	// Public Webhooks
	webhookHandler := handlers.NewWebhookHandler(engine, webhookSecret)
	r.Post("/webhooks/nomba", webhookHandler.HandleNombaWebhook)

	return r
}
