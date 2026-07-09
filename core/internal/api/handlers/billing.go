package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/nomba"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type AdminBillingHandler struct {
	nombaClient *nomba.Client
	q           sqlc.Querier
}

func NewAdminBillingHandler(nombaClient *nomba.Client, q sqlc.Querier) *AdminBillingHandler {
	return &AdminBillingHandler{nombaClient: nombaClient, q: q}
}

type CreateCheckoutRequest struct {
	IntegratorID string `json:"integrator_id"`
	Amount       string `json:"amount,omitempty"` // For top-up
	CallbackUrl  string `json:"callback_url"`
	Email        string `json:"email,omitempty"`
	Type         string `json:"type"` // "save_card" or "topup"
}

func (h *AdminBillingHandler) CreateCheckout(w http.ResponseWriter, r *http.Request) {
	var req CreateCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.IntegratorID == "" {
		http.Error(w, "integrator_id is required", http.StatusBadRequest)
		return
	}
	if req.CallbackUrl == "" {
		http.Error(w, "callback_url is required", http.StatusBadRequest)
		return
	}
	if req.Type != "save_card" && req.Type != "topup" {
		http.Error(w, "type must be 'save_card' or 'topup'", http.StatusBadRequest)
		return
	}

	orderRef := "ref_" + uuid.New().String()

	var amount float64 = 100.00 // Default for save card (authorization hold)
	if req.Type == "topup" && req.Amount != "" {
		parsedAmount, err := strconv.ParseFloat(req.Amount, 64)
		if err != nil {
			http.Error(w, "invalid amount format", http.StatusBadRequest)
			return
		}
		amount = parsedAmount
	}

	tokenizeCard := false
	if req.Type == "save_card" {
		tokenizeCard = true
	}

	checkoutReq := nomba.CheckoutOrderRequest{
		Order: nomba.OrderInfo{
			OrderReference: orderRef,
			Amount:         amount,
			Currency:       "NGN",
			CustomerEmail:  req.Email,
			CustomerId:     req.IntegratorID,
			CallbackUrl:    req.CallbackUrl,
		},
		TokenizeCard: tokenizeCard,
	}

	b, _ := json.Marshal(checkoutReq)
	slog.Info("Sending to Nomba", "payload", string(b))

	resp, err := h.nombaClient.CreateCheckoutOrder(r.Context(), checkoutReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type VerifyCheckoutRequest struct {
	IntegratorID string `json:"integrator_id"`
	OrderRef     string `json:"order_ref"`
}

func (h *AdminBillingHandler) VerifyCheckout(w http.ResponseWriter, r *http.Request) {
	var req VerifyCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.IntegratorID == "" || req.OrderRef == "" {
		http.Error(w, "integrator_id and order_ref are required", http.StatusBadRequest)
		return
	}

	_, err := uuid.Parse(req.IntegratorID)
	if err != nil {
		http.Error(w, "invalid integrator_id format", http.StatusBadRequest)
		return
	}

	resp, err := h.nombaClient.VerifyTransaction(r.Context(), req.OrderRef)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to verify transaction with Nomba: %v", err), http.StatusInternalServerError)
		return
	}

	if !resp.Success {
		http.Error(w, fmt.Sprintf("transaction not successful: %s", resp.Message), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"success"}`))
}
