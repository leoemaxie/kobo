package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/leoemaxie/kobo/internal/integrator"
)

type AdminHandler struct {
	svc *integrator.Service
}

func NewAdminHandler(svc *integrator.Service) *AdminHandler {
	return &AdminHandler{svc: svc}
}

type ProvisionIntegratorRequest struct {
	Name   string `json:"name"`
	IsLive bool   `json:"is_live"`
}

type ProvisionIntegratorResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
	IsSandbox bool   `json:"is_sandbox"`
}

func (h *AdminHandler) ProvisionIntegrator(w http.ResponseWriter, r *http.Request) {
	var req ProvisionIntegratorRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	result, err := h.svc.ProvisionIntegrator(r.Context(), req.Name, req.IsLive)
	if err != nil {
		http.Error(w, "failed to provision integrator", http.StatusInternalServerError)
		return
	}

	resp := ProvisionIntegratorResponse{
		ID:        result.Integrator.ID.String(),
		Name:      result.Integrator.Name,
		APIKey:    result.Credential.KeyID,
		APISecret: result.RawSecret,
		IsSandbox: result.Credential.Environment == "sandbox", // Ensure we treat "sandbox" as true
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
