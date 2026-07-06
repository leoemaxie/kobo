package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/leoemaxie/kobo/internal/api/dto"
	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/exceptions"
)

type ExceptionsHandler struct {
	svc *exceptions.Service
}

func NewExceptionsHandler(svc *exceptions.Service) *ExceptionsHandler {
	return &ExceptionsHandler{svc: svc}
}

func (h *ExceptionsHandler) ListOpen(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		} else {
			apierrors.WriteError(w, http.StatusBadRequest, "invalid_query", "limit must be between 1 and 1000")
			return
		}
	}

	entries, err := h.svc.ListOpen(r.Context(), integratorID, int32(limit), 0)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to get exceptions")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (h *ExceptionsHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	integratorID := middleware.GetIntegratorID(r.Context())

	idStr := chi.URLParam(r, "exceptionId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_id", "invalid exception ID")
		return
	}

	var req dto.ResolveExceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}

	if req.ResolutionAction == "" {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_request", "resolution_action is required")
		return
	}

	if err := h.svc.Resolve(r.Context(), id, integratorID, req.ResolutionAction, req.ResolutionNotes, req.SuccessorIdentityID); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to resolve exception")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "resolved"}`))
}
