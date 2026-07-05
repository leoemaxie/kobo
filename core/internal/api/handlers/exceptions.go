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
		limit, _ = strconv.Atoi(limitStr)
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

	if err := h.svc.Resolve(r.Context(), id, integratorID, req.ResolutionAction, req.ResolutionNotes, req.SuccessorIdentityID); err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to resolve exception")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "resolved"}`))
}
