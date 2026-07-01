package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/leoemaxie/kobo/internal/api/handlers"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

func NewRouter(q *sqlc.Queries, identityHandler *handlers.IdentityHandler) *chi.Mux {
	r := chi.NewRouter()

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(q))

		r.Post("/identities", identityHandler.Create)
		r.Get("/identities/{id}", identityHandler.Get)
	})

	// Public Webhooks
	r.Post("/webhooks/nomba", func(w http.ResponseWriter, r *http.Request) {
		// To be implemented in Phase 6
	})

	return r
}
