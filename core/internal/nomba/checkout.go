package nomba

import (
	"context"
	"net/http"
)

type CheckoutOrderRequest struct {
	Order        OrderInfo `json:"order"`
	TokenizeCard string    `json:"tokenizeCard,omitempty"`
}

type OrderInfo struct {
	OrderReference        string   `json:"orderReference"`
	Amount                string   `json:"amount"` // e.g. "10000.00"
	Currency              string   `json:"currency"`
	CustomerEmail         string   `json:"customerEmail,omitempty"`
	CustomerId            string   `json:"customerId,omitempty"`
	AccountId             string   `json:"accountId,omitempty"`
	CallbackUrl           string   `json:"callbackUrl,omitempty"`
	AllowedPaymentMethods []string `json:"allowedPaymentMethods,omitempty"`
}

type TransactionDetails struct {
	TransactionDate        string `json:"transactionDate"`
	PaymentReference       string `json:"paymentReference"`
	PaymentVendorReference string `json:"paymentVendorReference"`
	TokenizedCardPayment   string `json:"tokenizedCardPayment"`
	StatusCode             string `json:"statusCode"`
}

type TransferDetails struct {
	SessionId                  string `json:"sessionId"`
	BeneficiaryAccountName     string `json:"beneficiaryAccountName"`
	BeneficiaryAccountNumber   string `json:"beneficiaryAccountNumber"`
	OriginatorAccountName      string `json:"originatorAccountName"`
	OriginatorAccountNumber    string `json:"originatorAccountNumber"`
	Narration                  string `json:"narration"`
	DestinationInstitutionCode string `json:"destinationInstitutionCode"`
	PaymentReference           string `json:"paymentReference"`
}

type CardDetails struct {
	CardPan      string `json:"cardPan"`
	CardType     string `json:"cardType"`
	CardCurrency string `json:"cardCurrency"`
	CardBank     string `json:"cardBank"`
}

type VerifyTransactionResponse struct {
	Success            string             `json:"success"`
	Message            string             `json:"message"`
	Order              OrderInfo          `json:"order"`
	TransactionDetails TransactionDetails `json:"transactionDetails"`
	TransferDetails    TransferDetails    `json:"transferDetails"`
	CardDetails        CardDetails        `json:"cardDetails"`
}

type CheckoutOrderResponse struct {
	CheckoutLink   string `json:"checkoutLink"`
	OrderReference string `json:"orderReference"`
}

func (c *Client) CreateCheckoutOrder(ctx context.Context, req CheckoutOrderRequest) (CheckoutOrderResponse, error) {
	var resp struct {
		CheckoutLink   string `json:"checkoutLink"`
		OrderReference string `json:"orderReference"`
	}

	if err := c.doRequest(ctx, http.MethodPost, "/checkout/order", req, &resp, req.Order.OrderReference); err != nil {
		return CheckoutOrderResponse{}, err
	}

	return CheckoutOrderResponse{
		CheckoutLink:   resp.CheckoutLink,
		OrderReference: resp.OrderReference,
	}, nil
}

type ChargeTokenRequest struct {
	TokenKey       string `json:"tokenKey"`
	Amount         string `json:"amount"` // e.g. "10000.00"
	Currency       string `json:"currency"`
	OrderReference string `json:"orderReference"`
	CustomerEmail  string `json:"customerEmail,omitempty"`
}

type ChargeTokenResponse struct {
	TransactionID  string `json:"transactionId"`
	Status         string `json:"status"`
	OrderReference string `json:"orderReference"`
}

// ChargeToken charges a previously saved card tokenKey.
func (c *Client) ChargeToken(ctx context.Context, req ChargeTokenRequest) (ChargeTokenResponse, error) {
	var resp struct {
		TransactionID  string `json:"transactionId"`
		Status         string `json:"status"`
		OrderReference string `json:"orderReference"`
	}

	if err := c.doRequest(ctx, http.MethodPost, "/payments/token", req, &resp, req.OrderReference); err != nil {
		return ChargeTokenResponse{}, err
	}

	return ChargeTokenResponse{
		TransactionID:  resp.TransactionID,
		Status:         resp.Status,
		OrderReference: resp.OrderReference,
	}, nil
}
