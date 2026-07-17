package monnify

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type TransactionResult struct {
	ID                     string `json:"id"`
	Status                 string `json:"status"`
	Amount                 string `json:"amount"` // Note: This is string in JSON? Or float? "100.0" is string in example
	Type                   string `json:"type"`
	TimeCreated            string `json:"timeCreated"`
	PaymentVendorReference string `json:"paymentVendorReference"`
	RecipientAccountNumber string `json:"recipientAccountNumber"`
	RecipientAccountType   string `json:"recipientAccountType"`
	SenderName             string `json:"senderName"`
	EntryType              string `json:"entryType"`
	Narration              string `json:"narration"`
}

type FetchTransactionsResponse struct {
	Results []TransactionResult `json:"results"`
	Cursor  string              `json:"cursor"`
}

func (c *Client) FetchTransactions(ctx context.Context, accountNumber string, dateFrom, dateTo time.Time) ([]TransactionResult, error) {
	// Example uses: dateFrom=<date>&dateTo=<date>

	q := url.Values{}
	q.Set("virtual_account", accountNumber)
	q.Set("dateFrom", dateFrom.Format("2006-01-02T15:04:05Z"))
	q.Set("dateTo", dateTo.Format("2006-01-02T15:04:05Z"))

	path := fmt.Sprintf("/transactions/virtual?%s", q.Encode())

	var data FetchTransactionsResponse
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &data, ""); err != nil {
		return nil, err
	}

	return data.Results, nil
}

func (c *Client) FetchSingleTransaction(ctx context.Context, transactionRef string) (*TransactionResult, error) {
	q := url.Values{}
	q.Set("transactionRef", transactionRef)
	path := fmt.Sprintf("/transactions/accounts/single?%s", q.Encode())

	var data TransactionResult
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &data, ""); err != nil {
		return nil, err
	}

	return &data, nil
}
