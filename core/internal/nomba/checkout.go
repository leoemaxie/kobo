package nomba

import (
	"context"
	"net/http"
)

type CheckoutOrderRequest struct {
	Order               OrderInfo `json:"order"`
	CustomerEmail       string    `json:"customerEmail,omitempty"`
	CallbackUrl         string    `json:"callbackUrl,omitempty"`
	AllowedPaymentMethods []string  `json:"allowedPaymentMethods,omitempty"`
	SaveCard            bool      `json:"saveCard,omitempty"`
}

type OrderInfo struct {
	OrderReference string `json:"orderReference"`
	Amount         string `json:"amount"` // e.g. "10000.00"
	Currency       string `json:"currency"`
}

type CheckoutOrderResponse struct {
	CheckoutLink string `json:"checkoutLink"`
	OrderReference     string `json:"orderReference"`
}

// CreateCheckoutOrder creates a Nomba hosted checkout link.
func (c *Client) CreateCheckoutOrder(ctx context.Context, req CheckoutOrderRequest) (CheckoutOrderResponse, error) {
	var resp struct {
		CheckoutLink string `json:"checkoutLink"`
		OrderReference string `json:"orderReference"`
	}

	if err := c.doRequest(ctx, http.MethodPost, "/checkout/order", req, &resp, req.Order.OrderReference); err != nil {
		return CheckoutOrderResponse{}, err
	}

	return CheckoutOrderResponse{
		CheckoutLink: resp.CheckoutLink,
		OrderReference:     resp.OrderReference,
	}, nil
}

type ChargeTokenRequest struct {
	TokenKey      string    `json:"tokenKey"`
	Amount        string    `json:"amount"` // e.g. "10000.00"
	Currency      string    `json:"currency"`
	OrderReference string    `json:"orderReference"`
	CustomerEmail string    `json:"customerEmail,omitempty"`
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
