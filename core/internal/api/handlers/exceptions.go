package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	entries, err := h.svc.ListOpen(r.Context(), integratorID, limit, offset)
	if err != nil {
		http.Error(w, "failed to get exceptions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}
