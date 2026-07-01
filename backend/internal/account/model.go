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
	ID              uuid.UUID
	IdentityID      uuid.UUID
	NombaAccountRef string
	AccountNumber   *string
	BankName        *string
	AccountName     *string
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
