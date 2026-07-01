package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	apierrors "github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/ledger"
)

type LedgerHandler struct {
	svc *ledger.Service
}

func NewLedgerHandler(svc *ledger.Service) *LedgerHandler {
	return &LedgerHandler{svc: svc}
}

func (h *LedgerHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Pagination parameters
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"data": [], "next_cursor": null}`))
}

func (h *LedgerHandler) GetStatement(w http.ResponseWriter, r *http.Request) {
	// Integrator access check could go here if needed
	_ = middleware.GetIntegratorID(r.Context())

	idStr := chi.URLParam(r, "accountId")
	id, err := uuid.Parse(idStr)
	if err != nil {
		apierrors.WriteError(w, http.StatusBadRequest, "invalid_id", "invalid account ID")
		return
	}

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

	entries, err := h.svc.GetStatements(r.Context(), id, limit, offset)
	if err != nil {
		apierrors.WriteError(w, http.StatusInternalServerError, "internal_error", "failed to get statements")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
