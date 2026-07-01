package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/auth"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
	"github.com/stretchr/testify/assert"
)

type mockAuthQuerier struct {
	sqlc.Querier
	GetApiIntegratorByKeyFunc func(ctx context.Context, apiKey string) (sqlc.ApiIntegrator, error)
}

func (m *mockAuthQuerier) GetApiIntegratorByKey(ctx context.Context, apiKey string) (sqlc.ApiIntegrator, error) {
	if m.GetApiIntegratorByKeyFunc != nil {
		return m.GetApiIntegratorByKeyFunc(ctx, apiKey)
	}
	return sqlc.ApiIntegrator{}, errors.New("unimplemented GetApiIntegratorByKey")
}

func TestAuthMiddleware_Valid(t *testing.T) {
	apiKey, rawSecret, hashedSecret, _ := auth.GenerateCredentials(false)
	integratorID := uuid.New()

	mq := &mockAuthQuerier{
		GetApiIntegratorByKeyFunc: func(ctx context.Context, key string) (sqlc.ApiIntegrator, error) {
			assert.Equal(t, apiKey, key)
			return sqlc.ApiIntegrator{
				ID:            integratorID,
				Name:          "Test Integrator",
				ApiSecretHash: hashedSecret,
				IsSandbox:     true,
			}, nil
		},
	}

	handler := AuthMiddleware(mq)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctxData := GetIntegratorContext(r.Context())
		assert.Equal(t, integratorID, ctxData.ID)
		assert.Equal(t, "Test Integrator", ctxData.Name)
		assert.True(t, ctxData.IsSandbox)

		// Assert backward compatibility works
		assert.Equal(t, integratorID, GetIntegratorID(r.Context()))
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth(apiKey, rawSecret)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_MismatchedPrefix(t *testing.T) {
	mq := &mockAuthQuerier{} // Should not be called
	handler := AuthMiddleware(mq)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth("kobo_live_pk_123", "kobo_test_sk_456")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	apiKey, _, hashedSecret, _ := auth.GenerateCredentials(false)
	
	mq := &mockAuthQuerier{
		GetApiIntegratorByKeyFunc: func(ctx context.Context, key string) (sqlc.ApiIntegrator, error) {
			return sqlc.ApiIntegrator{
				ID:            uuid.New(),
				ApiSecretHash: hashedSecret,
				IsSandbox:     true,
			}, nil
		},
	}

	handler := AuthMiddleware(mq)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req := httptest.NewRequest("GET", "/", nil)
	req.SetBasicAuth(apiKey, "kobo_test_sk_wrongsecret")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
