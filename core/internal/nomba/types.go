package nomba

// WebhookPayload represents the Kobo-relevant subset of the Nomba webhook payload.
type WebhookPayload struct {
	EventType string `json:"event_type"`
	RequestID string `json:"requestId"`
	Data      struct {
		Merchant struct {
			WalletID string `json:"walletId"`
			UserID   string `json:"userId"`
		} `json:"merchant"`
		Transaction struct {
			AliasAccountNumber string      `json:"aliasAccountNumber"`
			Type               string      `json:"type"`
			TransactionID      string      `json:"transactionId"`
			ResponseCode       interface{} `json:"responseCode"`
			TransactionAmount  float64     `json:"transactionAmount"`
			Narration          string      `json:"narration"`
			Time               string      `json:"time"`
			AliasAccountType   string      `json:"aliasAccountType"`
		} `json:"transaction"`
		Customer struct {
			SenderName    string `json:"senderName"`
			AccountNumber string `json:"accountNumber"`
			BankName      string `json:"bankName"`
		} `json:"customer"`
		Order struct {
			Amount         float64 `json:"amount"`
			OrderReference string  `json:"orderReference"`
			PaymentMethod  string  `json:"paymentMethod"`
			CustomerId     string  `json:"customerId"`
		} `json:"order"`
		TokenizedCardData *TokenizedCardData `json:"tokenizedCardData,omitempty"`
	} `json:"data"`
}

type TokenizedCardData struct {
	TokenKey         string `json:"tokenKey"`
	CardType         string `json:"cardType"`
	TokenExpiryYear  string `json:"tokenExpiryYear"`
	TokenExpiryMonth string `json:"tokenExpiryMonth"`
	CardPan          string `json:"cardPan"`
}
