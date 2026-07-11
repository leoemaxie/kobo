package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type LogsHandler struct {
	q *sqlc.Queries
}

func NewLogsHandler(q *sqlc.Queries) *LogsHandler {
	return &LogsHandler{q: q}
}

type PaginatedLogsResponse struct {
	Data []LogEntry     `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type PaginationMeta struct {
	Total      int `json:"total"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalPages int `json:"totalPages"`
}

func (h *LogsHandler) GetLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := middleware.GetConsoleSessionContext(ctx)
	integratorUUID := session.IntegratorID
	if integratorUUID == uuid.Nil {
		errors.WriteError(w, http.StatusUnauthorized, "unauthorized", "Integrator ID missing from context")
		return
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	methodParam := r.URL.Query().Get("method")
	statusParam := r.URL.Query().Get("status_code")

	page := 1
	if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
		page = p
	}

	limit := 50
	if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
		limit = l
	}

	offset := (page - 1) * limit

	var methodFilter pgtype.Text
	if methodParam != "" {
		methodFilter = pgtype.Text{String: methodParam, Valid: true}
	}

	var statusFilter pgtype.Int4
	if s, err := strconv.Atoi(statusParam); err == nil && s > 0 {
		statusFilter = pgtype.Int4{Int32: int32(s), Valid: true}
	}

	countParams := sqlc.CountRequestLogsParams{
		IntegratorID: integratorUUID,
		Method:       methodFilter,
		StatusCode:   statusFilter,
	}

	totalLogs, err := h.q.CountRequestLogs(ctx, countParams)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to fetch logs count")
		return
	}

	paginatedParams := sqlc.GetPaginatedRequestLogsParams{
		IntegratorID: integratorUUID,
		Limit:        int32(limit),
		Offset:       int32(offset),
		Method:       methodFilter,
		StatusCode:   statusFilter,
	}

	logsResult, err := h.q.GetPaginatedRequestLogs(ctx, paginatedParams)
	if err != nil {
		errors.WriteError(w, http.StatusInternalServerError, "internal_error", "Failed to fetch paginated logs")
		return
	}

	data := []LogEntry{}
	for _, rl := range logsResult {
		data = append(data, LogEntry{
			Method: rl.Method,
			Path:   rl.Path,
			Status: int(rl.StatusCode),
			Ms:     int(rl.LatencyMs),
			ID:     rl.RequestID,
			Time:   timeSince(rl.CreatedAt),
		})
	}

	totalPages := int((totalLogs + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	res := PaginatedLogsResponse{
		Data: data,
		Meta: PaginationMeta{
			Total:      int(totalLogs),
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}
