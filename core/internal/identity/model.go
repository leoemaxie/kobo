package identity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type VirtualAccountSummary struct {
	AccountNumber      string `json:"account_number"`
	BankName           string `json:"bank_name"`
	AccountName        string `json:"account_name"`
	ExpectedAmountKobo *int64 `json:"expected_amount_kobo"`
	IsExpired          bool   `json:"is_expired"`
}

type Identity struct {
	ID                uuid.UUID              `json:"id"`
	IntegratorID      uuid.UUID              `json:"-"`
	ExternalReference string                 `json:"external_reference"`
	DisplayName       string                 `json:"display_name"`
	Metadata          json.RawMessage        `json:"metadata"`
	State             string                 `json:"state"`
	VirtualAccount    *VirtualAccountSummary `json:"virtual_account"`
	FailureReason     *string                `json:"failure_reason,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}
