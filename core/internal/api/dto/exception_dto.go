package dto

import "github.com/google/uuid"

// ResolveExceptionRequest represents the payload for resolving an exception.
type ResolveExceptionRequest struct {
	ResolutionAction    string     `json:"resolution_action" validate:"required"`
	ResolutionNotes     string     `json:"resolution_notes,omitempty"`
	SuccessorIdentityID *uuid.UUID `json:"successor_identity_id,omitempty"`
}
