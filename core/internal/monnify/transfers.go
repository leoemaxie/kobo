package monnify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BankInfo struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	NipCode string `json:"nipCode"`
	Logo    string `json:"logo"`
}

type LookupBankAccountResponse struct {
	AccountNumber string `json:"accountNumber"`
	AccountName   string `json:"accountName"`
}

type TransferToBankRequest struct {
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"accountNumber"`
	AccountName   string  `json:"accountName"`
	BankCode      string  `json:"bankCode"`
	MerchantTxRef string  `json:"merchantTxRef"`
	SenderName    string  `json:"senderName"`
	Narration     string  `json:"narration"`
}

type TransferToBankResponse struct {
	MonnifyTransferID string
	Status          string
	FeeNaira        float64
	IsAsync         bool
}

// ListBanks returns all Nigerian banks and their codes.
// Calls: GET /v1/transfers/banks
func (c *Client) ListBanks(ctx context.Context) ([]BankInfo, error) {
	var resp []BankInfo
	path := "/transfers/banks"
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp, ""); err != nil {
		return nil, err
	}
	return resp, nil
}

// LookupBankAccount resolves an account number to the account holder's name.
// Calls: POST /v1/transfers/bank/lookup
func (c *Client) LookupBankAccount(ctx context.Context, accountNumber, bankCode string) (LookupBankAccountResponse, error) {
	reqBody := map[string]string{
		"accountNumber": accountNumber,
		"bankCode":      bankCode,
	}
	var data LookupBankAccountResponse
	path := "/transfers/bank/lookup"
	if err := c.doRequest(ctx, http.MethodPost, path, reqBody, &data, ""); err != nil {
		return LookupBankAccountResponse{}, err
	}
	return data, nil
}

// TransferToBank initiates a bank transfer from Kobo's Monnify wallet.
// Calls: POST /v2/transfers/bank or POST /v2/transfers/bank/{subAccountId}
func (c *Client) TransferToBank(ctx context.Context, req TransferToBankRequest) (TransferToBankResponse, error) {
	token, err := c.tokenManager.GetValidToken(ctx)
	if err != nil {
		return TransferToBankResponse{}, fmt.Errorf("failed to get valid token: %w", err)
	}

	b, _ := json.Marshal(req)

	// Ensure we call v2
	baseURL := c.baseURL
	if strings.HasSuffix(baseURL, "/v1") {
		baseURL = strings.TrimSuffix(baseURL, "/v1") + "/v2"
	} else if strings.Contains(baseURL, "/v1/") {
		baseURL = strings.Replace(baseURL, "/v1/", "/v2/", 1)
	}

	path := fmt.Sprintf("/transfers/bank/%s", c.subAccountID)
	url := fmt.Sprintf("%s%s", baseURL, path)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return TransferToBankResponse{}, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("accountId", c.accountID)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return TransferToBankResponse{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return TransferToBankResponse{}, err
	}

	// 201 means async processing
	if resp.StatusCode == http.StatusCreated {
		return TransferToBankResponse{
			IsAsync: true,
		}, nil
	}

	var baseResp struct {
		Code        string          `json:"code"`
		Description string          `json:"description"`
		Data        json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(bodyBytes, &baseResp); err != nil {
		return TransferToBankResponse{}, fmt.Errorf("failed to decode response: %v, body: %s", err, string(bodyBytes))
	}

	// Accept "00" (V1 success code), "200" (V2 success code), and "201" (V2 processing code)
	if baseResp.Code != "00" && baseResp.Code != "200" && baseResp.Code != "201" {
		return TransferToBankResponse{}, fmt.Errorf("monnify API error %s: %s | RAW BODY: %s", baseResp.Code, baseResp.Description, string(bodyBytes))
	}

	if baseResp.Code == "201" {
		return TransferToBankResponse{
			IsAsync: true,
		}, nil
	}

	var transferData struct {
		ID     string  `json:"id"`
		Status string  `json:"status"`
		Fee    float64 `json:"fee"`
	}

	if len(baseResp.Data) > 0 {
		if err := json.Unmarshal(baseResp.Data, &transferData); err != nil {
			return TransferToBankResponse{}, fmt.Errorf("failed to decode data: %v", err)
		}
	}

	isAsync := transferData.Status == "PENDING_BILLING" || transferData.Status == "PROCESSING"

	return TransferToBankResponse{
		MonnifyTransferID: transferData.ID,
		Status:          transferData.Status,
		FeeNaira:        transferData.Fee,
		IsAsync:         isAsync,
	}, nil
}
