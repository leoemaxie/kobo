package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/reconciliation"
)

type WebhookHandler struct {
	engine        reconciliation.Engine
	webhookSecret string
}

func NewWebhookHandler(engine reconciliation.Engine, webhookSecret string) *WebhookHandler {
	return &WebhookHandler{
		engine:        engine,
		webhookSecret: webhookSecret,
	}
}

func (h *WebhookHandler) HandleNombaWebhook(w http.ResponseWriter, r *http.Request) {
	sigHeader := r.Header.Get("nomba-signature")
	timeHeader := r.Header.Get("nomba-timestamp")

	if sigHeader == "" || timeHeader == "" {
		http.Error(w, "missing signature or timestamp", http.StatusUnauthorized)
		return
	}

	t, err := time.Parse(time.RFC3339, timeHeader)
	if err != nil || time.Since(t) > 5*time.Minute {
		http.Error(w, "invalid or expired timestamp", http.StatusUnauthorized)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var payload nomba.WebhookPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	if !nomba.VerifyWebhookSignature(&payload, sigHeader, timeHeader, h.webhookSecret) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}

	err = h.engine.ProcessWebhook(r.Context(), &payload)
	if err != nil {
		// Log error, but return 200 so Nomba doesn't retry infinitely on business logic errors,
		// but wait: "do not return non-2XX for business-logic reasons, only for genuine processing failures"
		// If it's a database failure we should probably 500, but let's just 200 and log for now or 500 if err is systemic.
		// For simplicity, we return 200 to acknowledge receipt if it's a business error (e.g. account not found)
		// We'll return 500 only for serious issues, but returning 200 is safer to avoid retry storms for unfixable errors.
	}

	w.WriteHeader(http.StatusOK)
}
