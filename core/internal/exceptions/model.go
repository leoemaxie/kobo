package exceptions

import (
	"time"

	"github.com/google/uuid"
)

type Exception struct {
	ID                  uuid.UUID  `json:"id"`
	IntegratorID        uuid.UUID  `json:"integrator_id"`
	Type                string     `json:"type"`
	AmountKobo          int64      `json:"amount_kobo"`
	NombaReference      string     `json:"nomba_reference"`
	RelatedAccountID    *uuid.UUID `json:"related_account_id,omitempty"`
	Status              string     `json:"status"`
	ResolutionAction    *string    `json:"resolution_action,omitempty"`
	ResolutionNotes     *string    `json:"resolution_notes,omitempty"`
	SuccessorIdentityID *uuid.UUID `json:"successor_identity_id,omitempty"`
	DetectedAt          time.Time  `json:"detected_at"`
	ResolvedAt          *time.Time `json:"resolved_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}
