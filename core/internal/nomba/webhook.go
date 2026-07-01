package nomba

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

// VerifyWebhookSignature verifies the exact signature algorithm
func VerifyWebhookSignature(payload *WebhookPayload, signatureHeader string, timestampHeader string, secret string) bool {
	responseCode := ""
	if str, ok := payload.Data.Transaction.ResponseCode.(string); ok && str != "null" {
		responseCode = str
	}

	hashingPayload := fmt.Sprintf(
		"%s:%s:%s:%s:%s:%s:%s:%s:%s",
		payload.EventType,
		payload.RequestID,
		payload.Data.Merchant.UserID,
		payload.Data.Merchant.WalletID,
		payload.Data.Transaction.TransactionID,
		payload.Data.Transaction.Type,
		payload.Data.Transaction.Time,
		responseCode,
		timestampHeader,
	)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(hashingPayload))
	expectedSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return strings.EqualFold(expectedSig, signatureHeader)
}
