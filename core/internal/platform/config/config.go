package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	MonnifyBaseURL       string
	MonnifyClientID      string
	MonnifyClientSecret  string
	MonnifyAccountID     string
	MonnifySubAccountID  string
	MonnifyWebhookSecret string
	KYCMaxCeiling      int64 // Stored in kobo
}

func Load() (*Config, error) {
	_ = godotenv.Load() // Loads .env if it exists

	dbUrl, err := getEnvOrError("DATABASE_URL")
	if err != nil {
		return nil, err
	}
	clientID, err := getEnvOrError("MONNIFY_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := getEnvOrError("MONNIFY_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	accountID, err := getEnvOrError("MONNIFY_ACCOUNT_ID")
	if err != nil {
		return nil, err
	}
	webhookSecret, err := getEnvOrError("MONNIFY_WEBHOOK_SECRET")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:               getEnvOrDefault("PORT", "8080"),
		DatabaseURL:        dbUrl,
		MonnifyBaseURL:       getEnvOrDefault("MONNIFY_BASE_URL", "https://sandbox.api.monnify.com/v1"),
		MonnifyClientID:      clientID,
		MonnifyClientSecret:  clientSecret,
		MonnifyAccountID:     accountID,
		MonnifySubAccountID:  getEnvOrDefault("MONNIFY_SUB_ACCOUNT_ID", accountID),
		MonnifyWebhookSecret: webhookSecret,
		KYCMaxCeiling:      1000000 * 100, // 1 million Naira in kobo as default
	}
	return cfg, nil
}

func getEnvOrDefault(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}

func getEnvOrError(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("missing required environment variable: %s", key)
	}
	return val, nil
}
