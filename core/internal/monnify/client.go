package monnify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/leoemaxie/kobo/internal/account"
)

type TokenManager struct {
	mu           sync.Mutex
	clientID     string
	clientSecret string
	accountID    string
	baseURL      string
	client       *http.Client

	accessToken  string
	refreshToken string
	expiresAt    time.Time
}

func NewTokenManager(baseURL, clientID, clientSecret, accountID string, client *http.Client) *TokenManager {
	if client == nil {
		client = http.DefaultClient
	}
	return &TokenManager{
		clientID:     clientID,
		clientSecret: clientSecret,
		accountID:    accountID,
		baseURL:      baseURL,
		client:       client,
	}
}

func (tm *TokenManager) GetValidToken(ctx context.Context) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	// Proactively refresh 5 minutes before expiry
	if tm.accessToken != "" && time.Now().Add(5*time.Minute).Before(tm.expiresAt) {
		return tm.accessToken, nil
	}

	if tm.refreshToken != "" {
		if err := tm.refresh(ctx); err == nil {
			return tm.accessToken, nil
		}
	}

	if err := tm.issue(ctx); err != nil {
		return "", err
	}

	return tm.accessToken, nil
}

type authResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Data        struct {
		BusinessID   string    `json:"businessId"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
		ExpiresAt    time.Time `json:"expiresAt"`
	} `json:"data"`
}

func (tm *TokenManager) issue(ctx context.Context) error {
	url := fmt.Sprintf("%s/auth/token/issue", tm.baseURL)
	body := map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     tm.clientID,
		"client_secret": tm.clientSecret,
	}

	return tm.doAuthReq(ctx, url, body, "")
}

func (tm *TokenManager) refresh(ctx context.Context) error {
	url := fmt.Sprintf("%s/auth/token/refresh", tm.baseURL)
	body := map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": tm.refreshToken,
	}

	return tm.doAuthReq(ctx, url, body, tm.accessToken)
}

func (tm *TokenManager) doAuthReq(ctx context.Context, url string, body map[string]string, currentToken string) error {
	b, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accountId", tm.accountID)
	if currentToken != "" {
		req.Header.Set("Authorization", "Bearer "+currentToken)
	}

	resp, err := tm.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("auth request failed with status: %d", resp.StatusCode)
	}

	var authResp authResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return err
	}

	if authResp.Code != "00" {
		return fmt.Errorf("auth request failed with code: %s, description: %s", authResp.Code, authResp.Description)
	}

	tm.accessToken = authResp.Data.AccessToken
	tm.refreshToken = authResp.Data.RefreshToken
	tm.expiresAt = authResp.Data.ExpiresAt

	return nil
}

type Client struct {
	tokenManager *TokenManager
	baseURL      string
	accountID    string
	subAccountID string
	httpClient   *http.Client
}

func NewClient(baseURL, clientID, clientSecret, accountID, subAccountID string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	tm := NewTokenManager(baseURL, clientID, clientSecret, accountID, httpClient)
	return &Client{
		tokenManager: tm,
		baseURL:      baseURL,
		accountID:    accountID,
		subAccountID: subAccountID,
		httpClient:   httpClient,
	}
}

func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, response interface{}, idempotencyKey string) error {
	token, err := c.tokenManager.GetValidToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get valid token: %w", err)
	}

	var reqBody io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewReader(b)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("accountId", c.accountID)
	if idempotencyKey != "" {
		req.Header.Set("X-Idempotent-key", idempotencyKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var baseResp struct {
		Code        string          `json:"code"`
		Description string          `json:"description"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &baseResp); err != nil {
		return fmt.Errorf("failed to decode response: %v, body: %s", err, string(bodyBytes))
	}

	if baseResp.Code != "00" {
		return fmt.Errorf("monnify API error %s: %s | RAW BODY: %s", baseResp.Code, baseResp.Description, string(bodyBytes))
	}

	if response != nil && len(baseResp.Data) > 0 {
		return json.Unmarshal(baseResp.Data, response)
	}

	return nil
}

type createVAData struct {
	BankName          string `json:"bankName"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankAccountName   string `json:"bankAccountName"`
}

func (c *Client) CreateVirtualAccount(ctx context.Context, accountRef, accountName, bvn string) (account.MonnifyAccountResponse, error) {
	reqBody := map[string]interface{}{
		"accountRef":  accountRef,
		"accountName": accountName,
	}
	if bvn != "" {
		reqBody["bvn"] = bvn
	}

	idempotencyKey := accountRef

	path := "/accounts/virtual"
	if c.subAccountID != "" {
		path = fmt.Sprintf("/accounts/virtual/%s", c.subAccountID)
	}

	var data createVAData
	if err := c.doRequest(ctx, http.MethodPost, path, reqBody, &data, idempotencyKey); err != nil {
		return account.MonnifyAccountResponse{}, err
	}

	return account.MonnifyAccountResponse{
		AccountNumber:   data.BankAccountNumber,
		BankName:        data.BankName,
		BankAccountName: data.BankAccountName,
	}, nil
}

func (c *Client) VerifyTransaction(ctx context.Context, orderReference string) (VerifyTransactionResponse, error) {
	var resp VerifyTransactionResponse
	// The endpoint is GET /v1/checkout/transaction?idType=ORDER_REFERENCE&id=xxx
	path := fmt.Sprintf("/checkout/transaction?idType=ORDER_REFERENCE&id=%s", orderReference)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp, ""); err != nil {
		return VerifyTransactionResponse{}, err
	}
	return resp, nil
}

type FetchVirtualAccountResponse struct {
	CreatedAt         string `json:"createdAt"`
	AccountHolderID   string `json:"accountHolderId"`
	AccountRef        string `json:"accountRef"`
	BVN               string `json:"bvn"`
	AccountName       string `json:"accountName"`
	BankName          string `json:"bankName"`
	BankAccountNumber string `json:"bankAccountNumber"`
	BankAccountName   string `json:"bankAccountName"`
	Currency          string `json:"currency"`
	CallbackURL       string `json:"callbackUrl"`
	Expired           bool   `json:"expired"`
}

func (c *Client) FetchVirtualAccount(ctx context.Context, identifier string) (FetchVirtualAccountResponse, error) {
	var data FetchVirtualAccountResponse
	path := fmt.Sprintf("/accounts/virtual/%s", identifier)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &data, ""); err != nil {
		return FetchVirtualAccountResponse{}, err
	}
	return data, nil
}

type UpdateVirtualAccountRequest struct {
	NewAccountRef  string `json:"newAccountRef,omitempty"`
	AccountName    string `json:"accountName,omitempty"`
	ExpectedAmount string `json:"expectedAmount,omitempty"`
}

func (c *Client) UpdateVirtualAccount(ctx context.Context, identifier string, req UpdateVirtualAccountRequest) (bool, error) {
	var resp struct {
		Updated bool `json:"updated"`
	}
	path := fmt.Sprintf("/accounts/virtual/%s", identifier)
	if err := c.doRequest(ctx, http.MethodPut, path, req, &resp, ""); err != nil {
		return false, err
	}
	return resp.Updated, nil
}

func (c *Client) ExpireVirtualAccount(ctx context.Context, identifier string) (bool, error) {
	var resp struct {
		Expired bool `json:"expired"`
	}
	path := fmt.Sprintf("/accounts/virtual/%s", identifier)
	if err := c.doRequest(ctx, http.MethodDelete, path, nil, &resp, ""); err != nil {
		return false, err
	}
	return resp.Expired, nil
}
