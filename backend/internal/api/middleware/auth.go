package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type contextKey string
const IntegratorIDKey contextKey = "integratorID"

func AuthMiddleware(q *sqlc.Queries) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// In a real application, we would parse Basic Auth or Bearer Token and 
			// verify against api_integrators.api_key_hash using bcrypt.
			// For scaffolding purposes, we simulate checking the token.
			
			// We can assume for local dev the token is just the integrator ID directly for now
			token := strings.TrimPrefix(authHeader, "Bearer ")
			integratorID, err := uuid.Parse(token)
			if err != nil {
				// Fallback to basic auth user if bearer is not UUID
				username, _, ok := r.BasicAuth()
				if ok {
					integratorID, _ = uuid.Parse(username)
				}
			}

			if integratorID == uuid.Nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Context passing
			ctx := context.WithValue(r.Context(), IntegratorIDKey, integratorID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetIntegratorID(ctx context.Context) uuid.UUID {
	id, _ := ctx.Value(IntegratorIDKey).(uuid.UUID)
	return id
}
