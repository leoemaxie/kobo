package middleware

import (
	"context"
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/auth"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type contextKey string

const IntegratorContextKey contextKey = "integratorContext"

type IntegratorContext struct {
	ID        uuid.UUID
	Name      string
	IsSandbox bool
}

func AuthMiddleware(q sqlc.Querier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey, apiSecret, ok := r.BasicAuth()
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Fast fail on mismatched prefixes
			if strings.HasPrefix(apiKey, "kobo_live_") && !strings.HasPrefix(apiSecret, "kobo_live_") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			if strings.HasPrefix(apiKey, "kobo_test_") && !strings.HasPrefix(apiSecret, "kobo_test_") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			integrator, err := q.GetApiIntegratorByKey(r.Context(), apiKey)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hashedProvidedSecret := auth.HashSecret(apiSecret)

			if subtle.ConstantTimeCompare([]byte(integrator.ApiSecretHash), []byte(hashedProvidedSecret)) != 1 {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			integratorCtx := IntegratorContext{
				ID:        integrator.ID,
				Name:      integrator.Name,
				IsSandbox: integrator.IsSandbox,
			}

			ctx := context.WithValue(r.Context(), IntegratorContextKey, integratorCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetIntegratorContext(ctx context.Context) IntegratorContext {
	val, _ := ctx.Value(IntegratorContextKey).(IntegratorContext)
	return val
}

func GetIntegratorID(ctx context.Context) uuid.UUID {
	return GetIntegratorContext(ctx).ID
}
