package nomba

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
	httpClient   *http.Client
}

func NewClient(baseURL, clientID, clientSecret, accountID string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	tm := NewTokenManager(baseURL, clientID, clientSecret, accountID, httpClient)
	return &Client{
		tokenManager: tm,
		baseURL:      baseURL,
		accountID:    accountID,
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

	var baseResp struct {
		Code        string          `json:"code"`
		Description string          `json:"description"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&baseResp); err != nil {
		return err
	}

	if baseResp.Code != "00" {
		return fmt.Errorf("nomba API error %s: %s", baseResp.Code, baseResp.Description)
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

func (c *Client) CreateVirtualAccount(ctx context.Context, accountRef, accountName, bvn, kycTier string) (account.NombaAccountResponse, error) {
	reqBody := map[string]interface{}{
		"accountRef":  accountRef,
		"accountName": accountName,
	}
	if bvn != "" {
		reqBody["bvn"] = bvn
	}

	idempotencyKey := accountRef 

	var data createVAData
	if err := c.doRequest(ctx, http.MethodPost, "/accounts/virtual", reqBody, &data, idempotencyKey); err != nil {
		return account.NombaAccountResponse{}, err
	}

	return account.NombaAccountResponse{
		AccountNumber:   data.BankAccountNumber,
		BankName:        data.BankName,
		BankAccountName: data.BankAccountName,
	}, nil
}
