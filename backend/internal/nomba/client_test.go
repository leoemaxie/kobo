package nomba

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVirtualAccount(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/token/issue", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"code": "00",
			"description": "Success",
			"data": {
				"access_token": "mock-token",
				"refresh_token": "mock-refresh",
				"expiresAt": "2030-01-01T00:00:00Z"
			}
		}`))
	})
	
	mux.HandleFunc("/accounts/virtual", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{
			"code": "00",
			"description": "Success",
			"data": {
				"bankAccountNumber": "1234567890",
				"bankName": "Nomba MFB",
				"bankAccountName": "Nomba/Test User"
			}
		}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL, "client-id", "client-secret", "account-id", server.Client())
	
	resp, err := client.CreateVirtualAccount(context.Background(), "ref-123", "Test User", "", "tier_1")
	assert.NoError(t, err)
	assert.Equal(t, "1234567890", resp.AccountNumber)
	assert.Equal(t, "Nomba MFB", resp.BankName)
}
