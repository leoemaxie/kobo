package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type consoleContextKey string

const ConsoleSessionContextKey consoleContextKey = "consoleSessionContext"

type ConsoleSessionContext struct {
	UserID       uuid.UUID
	IntegratorID uuid.UUID
	Role         string
}

func ConsoleAuthMiddleware(q sqlc.Querier) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				apierrors.WriteError(w, http.StatusUnauthorized, "unauthorized", "Missing or invalid authorization header")
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")

			session, err := q.GetConsoleSession(r.Context(), token)
			if err != nil {
				apierrors.WriteError(w, http.StatusUnauthorized, "unauthorized", "Invalid or expired session")
				return
			}

			// Some console users might not have an integrator_id if they are superadmins,
			// but for payouts, we require an integrator_id.
			var integratorID uuid.UUID
			if session.IntegratorID.Valid {
				integratorID = session.IntegratorID.Bytes
			} else {
				apierrors.WriteError(w, http.StatusForbidden, "forbidden", "User is not associated with an integrator")
				return
			}

			var role string
			if session.Role != nil {
				role = session.Role.(string)
			}

			consoleCtx := ConsoleSessionContext{
				UserID:       session.UserID,
				IntegratorID: integratorID,
				Role:         role,
			}

			ctx := context.WithValue(r.Context(), ConsoleSessionContextKey, consoleCtx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ConsoleOwnerAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := GetConsoleSessionContext(r.Context())
			if session.Role != "owner" {
				apierrors.WriteError(w, http.StatusForbidden, "forbidden", "Only workspace owners can perform this action")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetConsoleSessionContext(ctx context.Context) ConsoleSessionContext {
	val, _ := ctx.Value(ConsoleSessionContextKey).(ConsoleSessionContext)
	return val
}
