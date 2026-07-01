package nomba

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret"
	timestamp := "2026-07-01T12:00:00Z"

	payload := &WebhookPayload{
		EventType: "payment_success",
		RequestID: "req-123",
	}
	payload.Data.Merchant.UserID = "user-123"
	payload.Data.Merchant.WalletID = "wallet-123"
	payload.Data.Transaction.TransactionID = "txn-123"
	payload.Data.Transaction.Type = "vact_transfer"
	payload.Data.Transaction.Time = "2026-07-01T12:00:00Z"
	payload.Data.Transaction.ResponseCode = "00"

	// Create valid signature
	hashingPayload := fmt.Sprintf(
		"%s:%s:%s:%s:%s:%s:%s:%s:%s",
		payload.EventType,
		payload.RequestID,
		payload.Data.Merchant.UserID,
		payload.Data.Merchant.WalletID,
		payload.Data.Transaction.TransactionID,
		payload.Data.Transaction.Type,
		payload.Data.Transaction.Time,
		"00",
		timestamp,
	)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(hashingPayload))
	validSig := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	t.Run("Valid Signature", func(t *testing.T) {
		assert.True(t, VerifyWebhookSignature(payload, validSig, timestamp, secret))
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		assert.False(t, VerifyWebhookSignature(payload, "invalid-sig", timestamp, secret))
	})

	t.Run("Null Response Code handling", func(t *testing.T) {
		payloadNull := *payload
		payloadNull.Data.Transaction.ResponseCode = "null"
		
		hashingPayloadNull := fmt.Sprintf(
			"%s:%s:%s:%s:%s:%s:%s:%s:%s",
			payloadNull.EventType,
			payloadNull.RequestID,
			payloadNull.Data.Merchant.UserID,
			payloadNull.Data.Merchant.WalletID,
			payloadNull.Data.Transaction.TransactionID,
			payloadNull.Data.Transaction.Type,
			payloadNull.Data.Transaction.Time,
			"", // "null" becomes ""
			timestamp,
		)
		macNull := hmac.New(sha256.New, []byte(secret))
		macNull.Write([]byte(hashingPayloadNull))
		validSigNull := base64.StdEncoding.EncodeToString(macNull.Sum(nil))

		assert.True(t, VerifyWebhookSignature(&payloadNull, validSigNull, timestamp, secret))
	})
}
