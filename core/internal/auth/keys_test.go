package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCredentials_Live(t *testing.T) {
	apiKey, rawSecret, hashedSecret, err := GenerateCredentials(true)
	assert.NoError(t, err)

	assert.True(t, strings.HasPrefix(apiKey, "kobo_live_pk_"))
	assert.True(t, strings.HasPrefix(rawSecret, "kobo_live_sk_"))

	// Verify Hash
	expectedHashBytes := sha256.Sum256([]byte(rawSecret))
	expectedHash := hex.EncodeToString(expectedHashBytes[:])

	assert.Equal(t, expectedHash, hashedSecret)
	assert.Equal(t, expectedHash, HashSecret(rawSecret))
}

func TestGenerateCredentials_Sandbox(t *testing.T) {
	apiKey, rawSecret, hashedSecret, err := GenerateCredentials(false)
	assert.NoError(t, err)

	assert.True(t, strings.HasPrefix(apiKey, "kobo_test_pk_"))
	assert.True(t, strings.HasPrefix(rawSecret, "kobo_test_sk_"))

	// Verify Hash
	expectedHashBytes := sha256.Sum256([]byte(rawSecret))
	expectedHash := hex.EncodeToString(expectedHashBytes[:])

	assert.Equal(t, expectedHash, hashedSecret)
}

func TestGenerateCredentials_Uniqueness(t *testing.T) {
	apiKey1, rawSecret1, _, err1 := GenerateCredentials(true)
	apiKey2, rawSecret2, _, err2 := GenerateCredentials(true)

	assert.NoError(t, err1)
	assert.NoError(t, err2)

	assert.NotEqual(t, apiKey1, apiKey2)
	assert.NotEqual(t, rawSecret1, rawSecret2)
}
