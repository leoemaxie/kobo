package dto

import "encoding/json"

// CreateIdentityRequest represents the payload for creating a new identity.
type CreateIdentityRequest struct {
	ExternalReference string          `json:"external_reference" validate:"required"`
	DisplayName       string          `json:"display_name" validate:"required"`
	Metadata          json.RawMessage `json:"metadata"`
}

// UpdateIdentityRequest represents the payload for updating an identity.
type UpdateIdentityRequest struct {
	DisplayName string          `json:"display_name"`
	Metadata    json.RawMessage `json:"metadata"`
}

// CloseIdentityRequest represents the payload for closing an identity.
type CloseIdentityRequest struct {
	Reason           string          `json:"reason" validate:"required"`
	SweepDestination json.RawMessage `json:"sweep_destination" validate:"required"`
}
