package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

// RequestTelemetryMiddleware logs API requests to the database
func RequestTelemetryMiddleware(q *sqlc.Queries) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap the response writer to capture status code
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			latencyMs := time.Since(start).Milliseconds()
			reqID := middleware.GetReqID(r.Context())

			// We only want to log requests that have an integrator_id in the context (authenticated API requests)
			ctx := r.Context()
			integratorIDStr, ok := ctx.Value("integrator_id").(string)
			if ok && integratorIDStr != "" {
				if integratorID, err := uuid.Parse(integratorIDStr); err == nil {
					// We do this in a goroutine so it doesn't block the request response
					go func() {
						// Using context.Background() because the request context is canceled when request ends
						_, _ = q.CreateRequestLog(context.Background(), sqlc.CreateRequestLogParams{
							IntegratorID: integratorID,
							Method:       r.Method,
							Path:         r.URL.Path,
							StatusCode:   int32(ww.Status()),
							LatencyMs:    int32(latencyMs),
							RequestID:    reqID,
						})
					}()
				}
			}
		})
	}
}
