package billing

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

// UsageRecorder writes usage_events to the DB.
// All writes are fire-and-forget (non-blocking) to avoid adding latency to API calls.
type UsageRecorder struct {
	q *sqlc.Queries
}

func NewUsageRecorder(q *sqlc.Queries) *UsageRecorder {
	return &UsageRecorder{q: q}
}

func (r *UsageRecorder) RecordAsync(integratorID uuid.UUID, environment sqlc.ConsoleEnvironment, eventType string, referenceID string, amountKobo int64) {
	// Execute in a goroutine so the calling API handler isn't blocked by the DB insert
	go func() {
		ctx := context.Background() // Use background context since the request context might be canceled
		err := r.q.InsertUsageEvent(ctx, sqlc.InsertUsageEventParams{
			IntegratorID: integratorID,
			Environment:  environment,
			EventType:    eventType,
			ReferenceID:  referenceID,
			AmountKobo:   amountKobo,
		})
		if err != nil {
			// Just log the error. In a production system we'd send this to Sentry/DataDog
			log.Printf("failed to record billing usage event: %v", err)
		}
	}()
}
