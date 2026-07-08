package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DatabaseURL        string
	NombaBaseURL       string
	NombaClientID      string
	NombaClientSecret  string
	NombaAccountID     string
	NombaSubAccountID  string
	NombaWebhookSecret string
	KYCMaxCeiling      int64 // Stored in kobo
}

func Load() (*Config, error) {
	_ = godotenv.Load() // Loads .env if it exists

	dbUrl, err := getEnvOrError("DATABASE_URL")
	if err != nil {
		return nil, err
	}
	clientID, err := getEnvOrError("NOMBA_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	clientSecret, err := getEnvOrError("NOMBA_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}
	accountID, err := getEnvOrError("NOMBA_ACCOUNT_ID")
	if err != nil {
		return nil, err
	}
	webhookSecret, err := getEnvOrError("NOMBA_WEBHOOK_SECRET")
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:               getEnvOrDefault("PORT", "8080"),
		DatabaseURL:        dbUrl,
		NombaBaseURL:       getEnvOrDefault("NOMBA_BASE_URL", "https://sandbox.api.nomba.com/v1"),
		NombaClientID:      clientID,
		NombaClientSecret:  clientSecret,
		NombaAccountID:     accountID,
		NombaSubAccountID:  getEnvOrDefault("NOMBA_SUB_ACCOUNT_ID", accountID),
		NombaWebhookSecret: webhookSecret,
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
