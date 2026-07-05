package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/leoemaxie/kobo"
	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/leoemaxie/kobo/internal/reconciliation"
	"gopkg.in/yaml.v3"
)

func NewRouter(q *sqlc.Queries, healthHandler *handlers.HealthHandler, identityHandler *handlers.IdentityHandler, ledgerHandler *handlers.LedgerHandler, exceptionsHandler *handlers.ExceptionsHandler, adminHandler *handlers.AdminHandler, engine reconciliation.Engine, webhookSecret string) *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		apierrors.WriteError(w, http.StatusNotFound, "not_found", "The requested resource was not found")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		apierrors.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "The requested method is not allowed for this resource")
	})

	var openapiDoc map[string]interface{}
	_ = yaml.Unmarshal(kobo.OpenAPI, &openapiDoc)
	openapiJSON, _ := json.Marshal(openapiDoc)

	r.Get("/healthz", healthHandler.HealthCheck)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(openapiJSON)
		})

		// Admin routes
		r.Post("/admin/integrators", adminHandler.ProvisionIntegrator)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(chimiddleware.RequestID)
			r.Use(chimiddleware.ClientIPFromXFF())
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
