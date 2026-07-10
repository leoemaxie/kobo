// Package kobo provides a zero-dependency Go client for the Kobo API v1.
// This file contains the HTTP transport layer.
package kobo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	defaultProductionBase = "https://api.kobo.triumphsystems.tech/v1"
	defaultSandboxBase    = "https://sandbox.api.kobo.triumphsystems.tech/v1"
	defaultTimeout        = 30 * time.Second
	sdkVersion            = "0.1.0"
)

// ─── Client Configuration ─────────────────────────────────────────────────────

// Option is a functional option for configuring the Client.
type Option func(*Client)

// WithHTTPClient replaces the default HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithBaseURL overrides the base URL (useful for testing against a local server).
func WithBaseURL(base string) Option {
	return func(c *Client) { c.baseURL = base }
}

// WithTimeout sets the HTTP request timeout (default: 30s).
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// ─── Client ───────────────────────────────────────────────────────────────────

// Client is the Kobo API client.
// Create one instance and reuse it across goroutines; it is safe for concurrent use.
type Client struct {
	apiKey     string
	apiSecret  string
	baseURL    string
	httpClient *http.Client

	// Resource-scoped sub-clients
	Identities  *IdentitiesService
	Accounts    *AccountsService
	Exceptions  *ExceptionsService
}

// New creates a new Kobo production client.
//
//	client := kobo.New("kobo_live_pk_...", "kobo_live_sk_...")
func New(apiKey, apiSecret string, opts ...Option) *Client {
	return newClient(apiKey, apiSecret, defaultProductionBase, opts...)
}

// NewSandbox creates a new Kobo sandbox client.
//
//	client := kobo.NewSandbox("kobo_test_pk_...", "kobo_test_sk_...")
func NewSandbox(apiKey, apiSecret string, opts ...Option) *Client {
	return newClient(apiKey, apiSecret, defaultSandboxBase, opts...)
}

func newClient(apiKey, apiSecret, base string, opts ...Option) *Client {
	c := &Client{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   base,
		httpClient: &http.Client{Timeout: defaultTimeout},
	}
	for _, o := range opts {
		o(c)
	}
	c.Identities = &IdentitiesService{c: c}
	c.Accounts = &AccountsService{c: c}
	c.Exceptions = &ExceptionsService{c: c}
	return c
}

// Health calls the unauthenticated /healthz endpoint.
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	var out HealthResponse
	if err := c.do(ctx, http.MethodGet, "/healthz", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ─── Internal HTTP helpers ────────────────────────────────────────────────────

func (c *Client) do(ctx context.Context, method, path string, query url.Values, body, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("kobo: marshal request: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	rawURL := c.baseURL + path
	if len(query) > 0 {
		rawURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return fmt.Errorf("kobo: build request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "kobo-go/"+sdkVersion)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("kobo: request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("kobo: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr APIError
		if jsonErr := json.Unmarshal(respBytes, &apiErr); jsonErr != nil {
			apiErr.Code = "unknown_error"
			apiErr.Message = string(respBytes)
		}
		apiErr.HTTPStatus = resp.StatusCode
		return &apiErr
	}

	if out != nil && len(respBytes) > 0 {
		if err := json.Unmarshal(respBytes, out); err != nil {
			return fmt.Errorf("kobo: decode response: %w", err)
		}
	}
	return nil
}

// ─── Pagination Helper ────────────────────────────────────────────────────────

func paginationQuery(opts PaginationOptions) url.Values {
	q := url.Values{}
	if opts.Cursor != nil {
		q.Set("page[cursor]", *opts.Cursor)
	}
	if opts.Limit != nil {
		q.Set("page[limit]", strconv.Itoa(*opts.Limit))
	}
	return q
}

// ─── Identities Service ───────────────────────────────────────────────────────

// IdentitiesService groups all /identities operations.
type IdentitiesService struct{ c *Client }

// Create registers a new identity and provisions its virtual account.
// The identity is returned immediately in PENDING state; poll Get or listen
// for the identity.activated webhook to confirm provisioning.
func (s *IdentitiesService) Create(ctx context.Context, req CreateIdentityRequest) (*Identity, error) {
	var out Identity
	if err := s.c.do(ctx, http.MethodPost, "/identities", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// List returns a list of identities for the integrator.
func (s *IdentitiesService) List(ctx context.Context, opts ListIdentitiesOptions) ([]Identity, error) {
	q := url.Values{}
	if opts.State != nil {
		q.Set("state", string(*opts.State))
	}
	if opts.Limit != nil {
		q.Set("limit", strconv.Itoa(*opts.Limit))
	}
	if opts.Offset != nil {
		q.Set("offset", strconv.Itoa(*opts.Offset))
	}
	
	var out []Identity
	if err := s.c.do(ctx, http.MethodGet, "/identities", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Get fetches a single identity record.
func (s *IdentitiesService) Get(ctx context.Context, identityID string) (*Identity, error) {
	var out Identity
	if err := s.c.do(ctx, http.MethodGet, "/identities/"+identityID, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Update patches display profile fields (display_name, metadata).
// At least one field must be present in the request.
func (s *IdentitiesService) Update(ctx context.Context, identityID string, req UpdateIdentityRequest) (*Identity, error) {
	var out Identity
	if err := s.c.do(ctx, http.MethodPatch, "/identities/"+identityID, nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Close initiates the closure of an identity's virtual account.
// Returns 202 Accepted immediately; closure completes asynchronously.
func (s *IdentitiesService) Close(ctx context.Context, identityID string, req CloseIdentityRequest) (*Identity, error) {
	var out Identity
	if err := s.c.do(ctx, http.MethodPost, "/identities/"+identityID+"/close", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Reopen transitions a CLOSED identity back to ACTIVE without re-running provisioning.
func (s *IdentitiesService) Reopen(ctx context.Context, identityID string) (*Identity, error) {
	var out Identity
	if err := s.c.do(ctx, http.MethodPost, "/identities/"+identityID+"/reopen", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ─── Accounts Service ─────────────────────────────────────────────────────────

// AccountsService groups all /accounts operations.
type AccountsService struct{ c *Client }

// ListTransactions returns a cursor-paginated list of reconciled transactions
// for the given account.
func (s *AccountsService) ListTransactions(ctx context.Context, accountID string, opts ListTransactionsOptions) (*TransactionPage, error) {
	q := paginationQuery(opts.PaginationOptions)
	var out TransactionPage
	if err := s.c.do(ctx, http.MethodGet, "/accounts/"+accountID+"/transactions", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetStatement returns the structured statement for the given period.
// Period is YYYY-MM; defaults to the current calendar month when nil.
func (s *AccountsService) GetStatement(ctx context.Context, accountID string, opts GetStatementOptions) (*Statement, error) {
	q := url.Values{}
	if opts.Period != nil {
		q.Set("period", *opts.Period)
	}
	var out Statement
	if err := s.c.do(ctx, http.MethodGet, "/accounts/"+accountID+"/statement", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ─── Exceptions Service ───────────────────────────────────────────────────────

// ExceptionsService groups all /exceptions operations.
type ExceptionsService struct{ c *Client }

// List returns a cursor-paginated list of exceptions.
func (s *ExceptionsService) List(ctx context.Context, opts ListExceptionsOptions) (*ExceptionPage, error) {
	q := paginationQuery(opts.PaginationOptions)
	if opts.Status != nil {
		q.Set("status", string(*opts.Status))
	}
	var out ExceptionPage
	if err := s.c.do(ctx, http.MethodGet, "/exceptions", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Resolve applies a resolution action to a flagged exception.
func (s *ExceptionsService) Resolve(ctx context.Context, exceptionID string, req ResolveExceptionRequest) (*Exception, error) {
	var out Exception
	if err := s.c.do(ctx, http.MethodPost, "/exceptions/"+exceptionID+"/resolve", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
