package account

import (
	"time"

	"github.com/google/uuid"
)

type State string

const (
	StatePending State = "pending"
	StateActive  State = "active"
	StateLimited State = "limited"
	StateClosing State = "closing"
	StateClosed  State = "closed"
	StateFailed  State = "failed"
)

type VirtualAccount struct {
	ID              uuid.UUID `json:"id"`
	IdentityID      uuid.UUID `json:"identity_id"`
	NombaAccountRef string    `json:"nomba_account_ref"`
	AccountNumber   *string   `json:"account_number,omitempty"`
	BankName        *string   `json:"bank_name,omitempty"`
	AccountName     *string   `json:"account_name,omitempty"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
