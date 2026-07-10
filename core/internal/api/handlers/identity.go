package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/account"
	"github.com/leoemaxie/kobo/internal/api/dto"
	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/billing"
	"github.com/leoemaxie/kobo/internal/identity"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type IdentityHandler struct {
	svc        *identity.Service
	accountSvc *account.Service
	recorder   *billing.UsageRecorder
}

func NewIdentityHandler(svc *identity.Service, accountSvc *account.Service, recorder *billing.UsageRecorder) *IdentityHandler {
	return &IdentityHandler{svc: svc, accountSvc: accountSvc, recorder: recorder}
}

func (h *IdentityHandler) Create(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())

	var req dto.CreateIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.LogAndWriteError(w, http.StatusBadRequest, "invalid_request", "failed to decode request", err)
		return
	}

	if req.ExternalReference == "" || req.DisplayName == "" {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", "external_reference and display_name are required")
		return
	}

	ident, err := h.svc.Register(r.Context(), integratorID, req.ExternalReference, req.DisplayName, req.Metadata)
	if err != nil {
		if errors.Is(err, identity.ErrIdentityConflict) {
			apierrors.LogAndWriteError(w, http.StatusConflict, "duplicate_external_reference", "identity with external reference already exists", err)
			return
		}
		apierrors.LogAndWriteError(w, http.StatusInternalServerError, "internal_error", "failed to create identity", err)
		return
	}

	// Synchronously provision the account (Serverless Option 1)
	if err := h.accountSvc.Provision(r.Context(), ident.ID, integratorID); err != nil {
		log.Printf("failed to provision account synchronously for identity %s: %v. Deleting identity to avoid failed state.", ident.ID, err)

		_ = h.svc.Delete(r.Context(), ident.ID, integratorID)

		apierrors.LogAndWriteError(w, http.StatusBadGateway, "provisioning_failed", "failed to create virtual account with provider", err)
		return
	}

	// Re-fetch the identity to return the updated state and virtual account details to the client
	if updatedIdent, err := h.svc.Get(r.Context(), ident.ID, integratorID); err == nil {
		ident = updatedIdent
	}

	env := sqlc.ConsoleEnvironmentProduction
	if middleware.GetIntegratorContext(r.Context()).IsSandbox {
		env = sqlc.ConsoleEnvironmentSandbox
	}
	h.recorder.RecordAsync(integratorID, env, "account_provisioned", ident.ID.String(), 5000) // ₦50 in kobo

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ident)
}

func (h *IdentityHandler) Get(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusBadRequest, "invalid_id", "invalid identity ID", err)
		return
	}

	ident, err := h.svc.Get(r.Context(), id, integratorID)
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusNotFound, "not_found", "identity not found", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ident)
}

func (h *IdentityHandler) List(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())

	limit := int32(50)
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.ParseInt(l, 10, 32); err == nil {
			limit = int32(parsed)
		}
	}

	offset := int32(0)
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.ParseInt(o, 10, 32); err == nil {
			offset = int32(parsed)
		}
	}

	var state *string
	if s := r.URL.Query().Get("state"); s != "" {
		state = &s
	}

	identities, err := h.svc.List(r.Context(), integratorID, state, limit, offset)
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusInternalServerError, "internal_error", "failed to list identities", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(identities)
}

func (h *IdentityHandler) Update(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusBadRequest, "invalid_id", "invalid identity ID", err)
		return
	}

	var req dto.UpdateIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.LogAndWriteError(w, http.StatusBadRequest, "invalid_request", "failed to decode request", err)
		return
	}

	var displayName *string
	if req.DisplayName != "" {
		displayName = &req.DisplayName
	}

	ident, err := h.svc.UpdateProfile(r.Context(), id, integratorID, displayName, req.Metadata)
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusInternalServerError, "internal_error", "failed to update identity", err)
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
		apierrors.LogAndWriteError(w, http.StatusBadRequest, "invalid_id", "invalid identity ID", err)
		return
	}

	var req dto.CloseIdentityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.LogAndWriteError(w, http.StatusBadRequest, "invalid_request", "failed to decode request", err)
		return
	}

	if req.Reason == "" || len(req.SweepDestination) == 0 || string(req.SweepDestination) == "null" {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", "reason and sweep_destination are required")
		return
	}

	err = h.accountSvc.Close(r.Context(), id, integratorID, req.Reason, req.SweepDestination)
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusInternalServerError, "internal_error", "failed to close identity", err)
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
		apierrors.LogAndWriteError(w, http.StatusBadRequest, "invalid_id", "invalid identity ID", err)
		return
	}

	err = h.accountSvc.Reopen(r.Context(), id, integratorID)
	if err != nil {
		apierrors.LogAndWriteError(w, http.StatusInternalServerError, "internal_error", "failed to reopen identity", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "reopened"}`))
}
