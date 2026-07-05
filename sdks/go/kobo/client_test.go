package kobo

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := New("pk_test", "sk_test")
	if c.baseURL != defaultProductionBase {
		t.Errorf("expected %q, got %q", defaultProductionBase, c.baseURL)
	}
	if c.apiKey != "pk_test" || c.apiSecret != "sk_test" {
		t.Error("credentials not set correctly")
	}

	cs := NewSandbox("pk_test", "sk_test")
	if cs.baseURL != defaultSandboxBase {
		t.Errorf("expected %q, got %q", defaultSandboxBase, cs.baseURL)
	}
}

func TestWithBaseURL(t *testing.T) {
	c := New("pk", "sk", WithBaseURL("http://localhost:8080"))
	if c.baseURL != "http://localhost:8080" {
		t.Errorf("expected localhost, got %q", c.baseURL)
	}
}

func TestHealth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthz" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"status":"ok","db":"ok"}`)
	}))
	defer ts.Close()

	client := New("pk", "sk", WithBaseURL(ts.URL))
	ctx := context.Background()
	res, err := client.Health(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "ok" || res.DB != "ok" {
		t.Errorf("unexpected response: %+v", res)
	}
}

func TestAPIError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, `{"code":"invalid_request","message":"bad parameter"}`)
	}))
	defer ts.Close()

	client := New("pk", "sk", WithBaseURL(ts.URL))
	_, err := client.Health(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.HTTPStatus != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", apiErr.HTTPStatus)
	}
	if apiErr.Code != "invalid_request" {
		t.Errorf("expected invalid_request, got %s", apiErr.Code)
	}
}

func TestIdentitiesCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/identities" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"id":"id_123","state":"pending","display_name":"Test User"}`)
	}))
	defer ts.Close()

	client := New("pk", "sk", WithBaseURL(ts.URL))
	res, err := client.Identities.Create(context.Background(), CreateIdentityRequest{
		DisplayName: "Test User",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != "id_123" {
		t.Errorf("expected id_123, got %s", res.ID)
	}
	if res.DisplayName != "Test User" {
		t.Errorf("expected Test User, got %s", res.DisplayName)
	}
}
