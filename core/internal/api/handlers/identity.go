package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/account"
	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/dto"
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

func (h *IdentityHandler) Create(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())

	var req dto.CreateIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	ident, err := h.svc.Register(r.Context(), integratorID, req.ExternalReference, req.DisplayName, req.KYCTier, req.Metadata)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error())
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
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_id", "invalid identity ID")
		return
	}

	ident, err := h.svc.Get(r.Context(), id, integratorID)
	if err != nil {
		apierrors.WriteError(w, http.StatusNotFound, "not_found", "identity not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ident)
}

func (h *IdentityHandler) Update(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_id", "invalid identity ID")
		return
	}

	var req dto.UpdateIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	var displayName *string
	if req.DisplayName != "" {
		displayName = &req.DisplayName
	}

	ident, err := h.svc.UpdateProfile(r.Context(), id, integratorID, displayName, req.Metadata)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ident)
}

func (h *IdentityHandler) Close(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_id", "invalid identity ID")
		return
	}

	var req dto.CloseIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	err = h.accountSvc.Close(r.Context(), id, integratorID, req.Reason, req.SweepDestination)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", err.Error())
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
