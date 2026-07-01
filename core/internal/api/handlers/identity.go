package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/account"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/identity"
)

type IdentityHandler struct {
	svc        *identity.Service
	accountSvc *account.Service
}

func NewIdentityHandler(svc *identity.Service, accountSvc *account.Service) *IdentityHandler {
	return &IdentityHandler{svc: svc, accountSvc: accountSvc}
}

type createIdentityReq struct {
	ExternalReference string          `json:"externalReference" validate:"required"`
	DisplayName       string          `json:"displayName" validate:"required"`
	KYCTier           string          `json:"kycTier" validate:"required"`
	Metadata          json.RawMessage `json:"metadata"`
}

func (h *IdentityHandler) Create(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())

	var req createIdentityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ident, err := h.svc.Register(r.Context(), integratorID, req.ExternalReference, req.DisplayName, req.KYCTier, req.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ident)
}

func (h *IdentityHandler) Get(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid identity ID", http.StatusBadRequest)
		return
	}

	ident, err := h.svc.Get(r.Context(), id, integratorID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ident)
}

type updateIdentityReq struct {
	DisplayName string          `json:"displayName"`
	Metadata    json.RawMessage `json:"metadata"`
}

func (h *IdentityHandler) Update(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid identity ID", http.StatusBadRequest)
		return
	}

	var req updateIdentityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var displayName *string
	if req.DisplayName != "" {
		displayName = &req.DisplayName
	}

	ident, err := h.svc.UpdateProfile(r.Context(), id, integratorID, displayName, req.Metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ident)
}

type closeIdentityReq struct {
	Reason           string          `json:"reason"`
	SweepDestination json.RawMessage `json:"sweep_destination"`
}

func (h *IdentityHandler) Close(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid identity ID", http.StatusBadRequest)
		return
	}

	var req closeIdentityReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.accountSvc.Close(r.Context(), id, integratorID, req.Reason, req.SweepDestination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"status": "closing initiated"}`))
}

func (h *IdentityHandler) Reopen(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid identity ID", http.StatusBadRequest)
		return
	}

	err = h.accountSvc.Reopen(r.Context(), id, integratorID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "reopened"}`))
}
