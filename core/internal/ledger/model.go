package ledger

import (
	"time"

	"github.com/google/uuid"
)

type LedgerEntry struct {
	ID               uuid.UUID
	VirtualAccountID uuid.UUID
	IdentityID       uuid.UUID
	AmountKobo       int64
	Direction        string
	Status           string
	NombaReference   string
	Source           string
	Narration        string
	SenderName       string
	OccurredAt       time.Time
	CreatedAt        time.Time
}
