package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
	"github.com/leoemaxie/kobo"
	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/leoemaxie/kobo/internal/reconciliation"
	"gopkg.in/yaml.v3"
	"log/slog"
)

func NewRouter(q *sqlc.Queries, healthHandler *handlers.HealthHandler, identityHandler *handlers.IdentityHandler, ledgerHandler *handlers.LedgerHandler, exceptionsHandler *handlers.ExceptionsHandler, adminHandler *handlers.AdminHandler, adminBillingHandler *handlers.AdminBillingHandler, payoutHandler *handlers.PayoutHandler, analyticsHandler *handlers.AnalyticsHandler, logsHandler *handlers.LogsHandler, engine reconciliation.Engine, webhookSecret string) *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		apierrors.WriteError(w, http.StatusNotFound, "not_found", "The requested resource was not found")
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, req *http.Request) {
		apierrors.WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "The requested method is not allowed for this resource")
	})

	var openapiDoc struct {
		OpenAPI    string                 `yaml:"openapi" json:"openapi"`
		Info       map[string]interface{} `yaml:"info" json:"info"`
		Servers    []interface{}          `yaml:"servers,omitempty" json:"servers,omitempty"`
		Security   []interface{}          `yaml:"security,omitempty" json:"security,omitempty"`
		Tags       []interface{}          `yaml:"tags,omitempty" json:"tags,omitempty"`
		Paths      map[string]interface{} `yaml:"paths" json:"paths"`
		Components map[string]interface{} `yaml:"components,omitempty" json:"components,omitempty"`
	}
	_ = yaml.Unmarshal(kobo.OpenAPI, &openapiDoc)
	openapiJSON, _ := json.Marshal(openapiDoc)

	r.Get("/healthz", healthHandler.HealthCheck)

	r.Route("/v1", func(r chi.Router) {
		r.Use(chimiddleware.RequestID)
		r.Use(chimiddleware.ClientIPFromXFF())
		r.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{
			Level: slog.LevelInfo,
		}))
		r.Use(middleware.Recoverer)

		r.Get("/", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write(openapiJSON)
		})

		// Admin routes
		r.Post("/admin/integrators", adminHandler.ProvisionIntegrator)
		r.Post("/admin/billing/checkout", adminBillingHandler.CreateCheckout)
		r.Post("/admin/billing/verify", adminBillingHandler.VerifyCheckout)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware(q))
			r.Use(middleware.RequestTelemetryMiddleware(q))

			r.Post("/identities", identityHandler.Create)
			r.Get("/identities", identityHandler.List)
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
		webhookHandler := handlers.NewWebhookHandler(engine, payoutHandler.Svc(), webhookSecret)
		r.Post("/webhooks/monnify", webhookHandler.HandleMonnifyWebhook)
	})

	// Console-only routes
	r.Route("/console", func(r chi.Router) {
		r.Use(chimiddleware.RequestID)
		r.Use(chimiddleware.ClientIPFromXFF())
		r.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{
			Level: slog.LevelInfo,
		}))
		r.Use(middleware.Recoverer)
		r.Use(middleware.ConsoleAuthMiddleware(q))
		r.Get("/analytics", analyticsHandler.GetAnalytics)
		r.Get("/logs", logsHandler.GetLogs)

		r.Route("/payouts", func(r chi.Router) {
			r.Use(middleware.ConsoleOwnerAuthMiddleware())
			r.Get("/banks", payoutHandler.ListBanks)
			r.Post("/bank-account/lookup", payoutHandler.LookupBankAccount)
			r.Post("/bank-account", payoutHandler.SaveBankAccount)
			r.Get("/bank-account", payoutHandler.GetBankAccount)
			r.Post("/request", payoutHandler.RequestPayout)
			r.Get("/", payoutHandler.ListPayouts)
		})
	})

	return r
}
