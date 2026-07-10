// Package kobo provides a zero-dependency Go client for the Kobo API v1.
//
// Authentication uses HTTP Basic Auth (API key as username, API secret as
// password). All monetary amounts are integers in kobo (1/100 Naira).
// All timestamps are ISO 8601 / RFC 3339 in UTC.
package kobo

import "time"

// ─── Enums ────────────────────────────────────────────────────────────────────

// IdentityState is the lifecycle state of a virtual account.
type IdentityState string

const (
	IdentityStatePending IdentityState = "pending"
	IdentityStateActive  IdentityState = "active"
	IdentityStateLimited IdentityState = "limited"
	IdentityStateClosing IdentityState = "closing"
	IdentityStateClosed  IdentityState = "closed"
	IdentityStateFailed  IdentityState = "failed"
)



// TransactionDirection is the direction of a transaction (v1 only supports inbound).
type TransactionDirection string

const TransactionDirectionInbound TransactionDirection = "inbound"

// TransactionStatus is the reconciliation status of a transaction.
type TransactionStatus string

const (
	TransactionStatusMatched     TransactionStatus = "matched"
	TransactionStatusPartial     TransactionStatus = "partial"
	TransactionStatusOverpayment TransactionStatus = "overpayment"
)

// ExceptionType is the category of a flagged exception.
type ExceptionType string

const (
	ExceptionTypePaymentToClosedAccount   ExceptionType = "payment_to_closed_account"
	ExceptionTypePaymentToUnknownAccount  ExceptionType = "payment_to_unknown_account"
	ExceptionTypePaymentDuringClosing     ExceptionType = "payment_during_closing"
)

// ExceptionStatus indicates whether an exception has been resolved.
type ExceptionStatus string

const (
	ExceptionStatusOpen     ExceptionStatus = "open"
	ExceptionStatusResolved ExceptionStatus = "resolved"
)

// ResolutionAction is the action taken to resolve an exception.
type ResolutionAction string

const (
	ResolutionActionReturnToSender      ResolutionAction = "return_to_sender"
	ResolutionActionRedirectToSuccessor ResolutionAction = "redirect_to_successor"
	ResolutionActionManualOverride      ResolutionAction = "manual_override"
)

// SweepDestinationType indicates where funds should be swept on account closure.
type SweepDestinationType string

const (
	SweepDestinationRefundToSource      SweepDestinationType = "refund_to_source"
	SweepDestinationIntegratorAccount   SweepDestinationType = "integrator_account"
	SweepDestinationSuccessorIdentity   SweepDestinationType = "successor_identity"
)

// ─── Core Models ──────────────────────────────────────────────────────────────

// VirtualAccountSummary contains the bank details for an identity's virtual account.
type VirtualAccountSummary struct {
	AccountNumber      string `json:"account_number"`
	BankName           string `json:"bank_name"`
	AccountName        string `json:"account_name"`
	ExpectedAmountKobo *int64 `json:"expected_amount_kobo,omitempty"`
	IsExpired          bool   `json:"is_expired"`
}

// Identity is the root resource in Kobo.
type Identity struct {
	ID                string                 `json:"id"`
	ExternalReference string                 `json:"external_reference"`
	DisplayName       string                 `json:"display_name"`
	State             IdentityState          `json:"state"`
	VirtualAccount    *VirtualAccountSummary `json:"virtual_account,omitempty"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	FailureReason     *string                `json:"failure_reason,omitempty"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
}

// Transaction is a single reconciled ledger entry.
type Transaction struct {
	ID             string               `json:"id"`
	AccountID      string               `json:"account_id"`
	// AmountKobo is the transaction amount in kobo (1/100 Naira). Always positive.
	AmountKobo     int64                `json:"amount_kobo"`
	Direction      TransactionDirection `json:"direction"`
	Status         TransactionStatus    `json:"status"`
	NombaReference string               `json:"nomba_reference"`
	OccurredAt     time.Time            `json:"occurred_at"`
}

// Statement is a structured monthly account statement.
type Statement struct {
	AccountID           string        `json:"account_id"`
	Period              string        `json:"period"`
	OpeningBalanceKobo  int64         `json:"opening_balance_kobo"`
	ClosingBalanceKobo  int64         `json:"closing_balance_kobo"`
	TotalInflowKobo     int64         `json:"total_inflow_kobo"`
	Transactions        []Transaction `json:"transactions"`
}

// ExceptionResolution contains the details of a resolved exception.
type ExceptionResolution struct {
	Action ResolutionAction `json:"action"`
	Notes  *string          `json:"notes,omitempty"`
}

// Exception is a misdirected or unmatched payment requiring manual resolution.
type Exception struct {
	ID               string               `json:"id"`
	Type             ExceptionType        `json:"type"`
	AmountKobo       int64                `json:"amount_kobo"`
	NombaReference   string               `json:"nomba_reference"`
	RelatedAccountID *string              `json:"related_account_id,omitempty"`
	Status           ExceptionStatus      `json:"status"`
	Resolution       *ExceptionResolution `json:"resolution,omitempty"`
	DetectedAt       time.Time            `json:"detected_at"`
	ResolvedAt       *time.Time           `json:"resolved_at,omitempty"`
}

// ─── Request Bodies ───────────────────────────────────────────────────────────

// CreateIdentityRequest is the payload to register a new identity.
type CreateIdentityRequest struct {
	ExternalReference string                 `json:"external_reference"`
	DisplayName       string                 `json:"display_name"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateIdentityRequest is the payload to update display profile fields.
// At least one field must be set.
type UpdateIdentityRequest struct {
	DisplayName *string                `json:"display_name,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// SweepDestination describes where account funds should go on closure.
type SweepDestination struct {
	Type                        SweepDestinationType `json:"type"`
	SuccessorIdentityID         *string              `json:"successor_identity_id,omitempty"`
	IntegratorAccountReference  *string              `json:"integrator_account_reference,omitempty"`
}

// CloseIdentityRequest initiates the account closure process.
type CloseIdentityRequest struct {
	SweepDestination SweepDestination `json:"sweep_destination"`
	Reason           *string          `json:"reason,omitempty"`
}

// ResolveExceptionRequest applies a resolution to a flagged exception.
type ResolveExceptionRequest struct {
	Action              ResolutionAction `json:"action"`
	SuccessorIdentityID *string          `json:"successor_identity_id,omitempty"`
	Notes               *string          `json:"notes,omitempty"`
}

// ─── Paginated Responses ──────────────────────────────────────────────────────

// Page is a cursor-paginated response envelope.
type Page[T any] struct {
	Data       []T     `json:"data"`
	NextCursor *string `json:"next_cursor"`
}

// TransactionPage is a paginated list of transactions.
type TransactionPage = Page[Transaction]

// ExceptionPage is a paginated list of exceptions.
type ExceptionPage = Page[Exception]

// ─── Query Options ────────────────────────────────────────────────────────────

// PaginationOptions carries cursor-based pagination parameters.
type PaginationOptions struct {
	Cursor *string
	Limit  *int
}

// ListExceptionsOptions carries filtering and pagination for exceptions.
type ListExceptionsOptions struct {
	PaginationOptions
	Status *ExceptionStatus
}

// GetStatementOptions carries the optional period parameter.
type GetStatementOptions struct {
	// Period is a YYYY-MM formatted string. Defaults to the current month.
	Period *string
}

// ListTransactionsOptions carries pagination for transactions.
type ListTransactionsOptions struct {
	PaginationOptions
}

// ListIdentitiesOptions carries filtering and pagination for identities.
type ListIdentitiesOptions struct {
	State  *IdentityState
	Limit  *int
	Offset *int
}

// ─── Health ───────────────────────────────────────────────────────────────────

// HealthResponse is the response from the /healthz endpoint.
type HealthResponse struct {
	Status string `json:"status"`
	DB     string `json:"db"`
}

// ─── Errors ───────────────────────────────────────────────────────────────────

// APIError represents a structured error returned by the Kobo API.
// Branch on Code, not Message, as Message is not stable.
// See ErrorCode constants (e.g., ErrorCodeIdentityNotFound) for stable codes.
type APIError struct {
	// HTTPStatus is the HTTP status code.
	HTTPStatus int
	// Code is the stable machine-readable error code.
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return "kobo: " + e.Code + ": " + e.Message
}
