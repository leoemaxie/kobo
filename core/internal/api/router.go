package api

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/leoemaxie/kobo/internal/reconciliation"
)

func NewRouter(q *sqlc.Queries, identityHandler *handlers.IdentityHandler, ledgerHandler *handlers.LedgerHandler, exceptionsHandler *handlers.ExceptionsHandler, adminHandler *handlers.AdminHandler, engine reconciliation.Engine, webhookSecret string) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/healthz", handlers.HealthCheck)

	r.Route("/v1", func(r chi.Router) {
		// Admin routes
		r.Post("/admin/integrators", adminHandler.ProvisionIntegrator)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(chimiddleware.RequestID)
			r.Use(chimiddleware.RealIP)
			r.Use(middleware.RequestLogger)
			r.Use(middleware.Recoverer)
			r.Use(middleware.AuthMiddleware(q))

			r.Post("/identities", identityHandler.Create)
			r.Get("/identities/{id}", identityHandler.Get)
			r.Patch("/identities/{id}", identityHandler.Update)
			r.Post("/identities/{id}/close", identityHandler.Close)
			r.Post("/identities/{id}/reopen", identityHandler.Reopen)

			r.Get("/accounts/{accountId}/transactions", ledgerHandler.GetTransactions)
			r.Get("/accounts/{accountId}/statement", ledgerHandler.GetStatement)

			r.Get("/exceptions", exceptionsHandler.ListOpen)
			r.Post("/exceptions/{exceptionId}/resolve", exceptionsHandler.Resolve)
		})

		// Public Webhooks
		webhookHandler := handlers.NewWebhookHandler(engine, webhookSecret)
		r.Post("/webhooks/nomba", webhookHandler.HandleNombaWebhook)
	})

	return r
}
