package dto

import (
	"github.com/google/uuid"
	"time"
)

// LedgerEntryResponse represents a single transaction in the ledger.
type LedgerEntryResponse struct {
	ID               uuid.UUID `json:"id"`
	VirtualAccountID uuid.UUID `json:"virtual_account_id"`
	IdentityID       uuid.UUID `json:"identity_id"`
	AmountKobo       int64     `json:"amount_kobo"`
	Direction        string    `json:"direction"`
	Status           string    `json:"status"`
	NombaReference   string    `json:"nomba_reference"`
	Source           string    `json:"source"`
	Narration        string    `json:"narration,omitempty"`
	SenderName       string    `json:"sender_name,omitempty"`
	OccurredAt       time.Time `json:"occurred_at"`
	CreatedAt        time.Time `json:"created_at"`
}
