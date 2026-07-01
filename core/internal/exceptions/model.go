package exceptions

import (
	"time"

	"github.com/google/uuid"
)

type Exception struct {
	ID                  uuid.UUID
	IntegratorID        uuid.UUID
	Type                string
	AmountKobo          int64
	NombaReference      string
	RelatedAccountID    *uuid.UUID
	Status              string
	ResolutionAction    *string
	ResolutionNotes     *string
	SuccessorIdentityID *uuid.UUID
	DetectedAt          time.Time
	ResolvedAt          *time.Time
	CreatedAt           time.Time
}
