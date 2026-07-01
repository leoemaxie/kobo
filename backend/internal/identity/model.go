package identity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DisplayProfile struct {
	DisplayName string
	Metadata    json.RawMessage
}

type Identity struct {
	ID                uuid.UUID
	IntegratorID      uuid.UUID
	ExternalReference string
	Profile           DisplayProfile
	KYCTier           string
	State             string
	FailureReason     *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
