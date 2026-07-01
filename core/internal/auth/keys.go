package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateCredentials creates a new unhashed API Key and a new API Secret.
// It returns the public API Key, the raw API Secret, and the SHA-256 hash of the API Secret.
func GenerateCredentials(isLive bool) (apiKey string, rawSecret string, hashedSecret string, err error) {
	keyBytes := make([]byte, 24)
	if _, err := rand.Read(keyBytes); err != nil {
		return "", "", "", fmt.Errorf("failed to generate secure bytes: %w", err)
	}

	secretBytes := make([]byte, 32)
	if _, err := rand.Read(secretBytes); err != nil {
		return "", "", "", fmt.Errorf("failed to generate secure bytes: %w", err)
	}

	envPrefix := "test"
	if isLive {
		envPrefix = "live"
	}

	apiKey = fmt.Sprintf("kobo_%s_pk_%s", envPrefix, hex.EncodeToString(keyBytes))
	rawSecret = fmt.Sprintf("kobo_%s_sk_%s", envPrefix, hex.EncodeToString(secretBytes))

	hash := sha256.Sum256([]byte(rawSecret))
	hashedSecret = hex.EncodeToString(hash[:])

	return apiKey, rawSecret, hashedSecret, nil
}

// HashSecret computes the SHA-256 hash of a provided secret for comparison.
func HashSecret(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return hex.EncodeToString(hash[:])
}
