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

	limitStr := r.URL.Query().Get("limit")
	limit := 50
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	entries, err := h.svc.ListOpen(r.Context(), integratorID, int32(limit), 0)
	if err != nil {
		http.Error(w, "failed to get exceptions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func (h *ExceptionsHandler) Resolve(w http.ResponseWriter, r *http.Request) {
	// Exception resolution logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "resolved"}`))
}
