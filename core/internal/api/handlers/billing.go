package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/nomba"
)

type AdminBillingHandler struct {
	nombaClient *nomba.Client
}

func NewAdminBillingHandler(nombaClient *nomba.Client) *AdminBillingHandler {
	return &AdminBillingHandler{nombaClient: nombaClient}
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

	orderRef := "ref_" + uuid.New().String()

	amount := "100.00" // Default for save card (authorization hold)
	if req.Type == "topup" && req.Amount != "" {
		amount = req.Amount
	}

	checkoutReq := nomba.CheckoutOrderRequest{
		Order: nomba.OrderInfo{
			OrderReference: orderRef,
			Amount:         amount,
			Currency:       "NGN",
		},
		CallbackUrl: req.CallbackUrl,
		CustomerEmail: req.Email,
		AllowedPaymentMethods: []string{"card"},
		SaveCard: req.Type == "save_card",
	}

	resp, err := h.nombaClient.CreateCheckoutOrder(r.Context(), checkoutReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Additional handlers like GetUsage, GetInvoices, etc., can be queried directly via Drizzle in Console
// Since Console has direct DB access, we only expose the Nomba proxy endpoints here.
