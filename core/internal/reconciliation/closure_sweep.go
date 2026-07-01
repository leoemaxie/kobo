package reconciliation

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/account"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type ClosureSweeper struct {
	q *sqlc.Queries
}

func NewClosureSweeper(q *sqlc.Queries) *ClosureSweeper {
	return &ClosureSweeper{q: q}
}

func (s *ClosureSweeper) Run(ctx context.Context) error {
	log.Println("Running closure sweep...")

	identities, err := s.q.ListAllIdentitiesByState(ctx, "closing")
	if err != nil {
		return err
	}

	for _, ident := range identities {
		// Mock sweeping logic:
		// 1. We would check the Nomba virtual account balance.
		// 2. If > 0, we'd trigger a Nomba payout to the sweep_destination.
		// 3. Since this is Kobo v1 and outbound API is not integrated, we simulate
		//    that the funds were swept by transitioning the identity to CLOSED.

		log.Printf("Simulating sweep for identity %s", ident.ID)

		newState, err := account.ValidTransition(account.State(ident.State), account.EventZeroBalance)
		if err != nil {
			log.Printf("Invalid state transition for identity %s: %v", ident.ID, err)
			continue
		}

		_, err = s.q.UpdateIdentityState(ctx, sqlc.UpdateIdentityStateParams{
			ID:            ident.ID,
			IntegratorID:  ident.IntegratorID,
			State:         string(newState),
			FailureReason: pgtype.Text{Valid: false},
		})
		if err != nil {
			log.Printf("Failed to update identity state for %s: %v", ident.ID, err)
			continue
		}

		// Look up the closing_started event to find the destination
		events, _ := s.q.ListIdentityEvents(ctx, ident.ID)
		var dest json.RawMessage
		for i := len(events) - 1; i >= 0; i-- {
			if events[i].EventType == "closing_started" {
				var parsed struct {
					SweepDestination json.RawMessage `json:"sweep_destination"`
				}
				if json.Unmarshal(events[i].Detail, &parsed) == nil {
					dest = parsed.SweepDestination
				}
				break
			}
		}
		if len(dest) == 0 {
			dest = []byte(`{}`)
		}

		detail := map[string]interface{}{
			"sweep_destination": dest,
			"simulated":         true,
		}
		detailBytes, _ := json.Marshal(detail)

		_, _ = s.q.InsertIdentityEvent(ctx, sqlc.InsertIdentityEventParams{
			ID:            ident.ID,
			IdentityID:    ident.ID,
			EventType:     "closed",
			PreviousState: pgtype.Text{String: ident.State, Valid: true},
			NewState:      pgtype.Text{String: string(newState), Valid: true},
			Detail:        detailBytes,
		})

		log.Printf("Identity %s successfully closed", ident.ID)
	}

	return nil
}
