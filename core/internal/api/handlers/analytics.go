package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"

	"github.com/leoemaxie/kobo/internal/api/errors"
	"github.com/leoemaxie/kobo/internal/api/middleware"
	"github.com/leoemaxie/kobo/internal/platform/db/sqlc"
)

type AnalyticsHandler struct {
	q *sqlc.Queries
}

func NewAnalyticsHandler(q *sqlc.Queries) *AnalyticsHandler {
	return &AnalyticsHandler{
		q: q,
	}
}

type Metric struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	Delta string `json:"delta,omitempty"`
	Sub   string `json:"sub,omitempty"`
	Trend string `json:"trend,omitempty"`
	Bar   int    `json:"bar,omitempty"`
}

type LogEntry struct {
	Method string `json:"method"`
	Path   string `json:"path"`
	Status int    `json:"status"`
	Ms     int    `json:"ms"`
	ID     string `json:"id"`
	Time   string `json:"time"`
}

type AnalyticsResponse struct {
	Metrics []Metric   `json:"metrics"`
	Logs    []LogEntry `json:"logs"`
}

func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session := middleware.GetConsoleSessionContext(ctx)
	integratorUUID := session.IntegratorID
	if integratorUUID == uuid.Nil {
		errors.WriteError(w, http.StatusUnauthorized, "unauthorized", "Integrator ID missing from context")
		return
	}

	// Fetch Total API Requests
	totalRequests, err := h.q.GetTotalApiRequests(ctx, integratorUUID)
	if err != nil {
		totalRequests = 0
	}

	// Fetch Virtual Accounts
	virtualAccounts, err := h.q.CountVirtualAccountsByIntegrator(ctx, integratorUUID)
	if err != nil {
		virtualAccounts = 0
	}

	// Fetch Error Rate
	errorRate, err := h.q.GetErrorRate(ctx, integratorUUID)
	if err != nil {
		errorRate = 0
	}

	// Fetch P99 Latency
	p99Latency, err := h.q.GetP99Latency(ctx, integratorUUID)
	if err != nil {
		p99Latency = 0
	}

	metrics := []Metric{
		{
			Key:   "api_requests",
			Label: "API Requests",
			Value: strconv.FormatInt(totalRequests, 10),
			Sub:   "Last 30 days",
			Trend: "neutral",
			Delta: "0%",
			Bar:   0,
		},
		{
			Key:   "virtual_accounts",
			Label: "Virtual Accounts",
			Value: strconv.FormatInt(virtualAccounts, 10),
			Sub:   "Total",
			Trend: "neutral",
			Delta: "0",
			Bar:   0,
		},
		{
			Key:   "error_rate",
			Label: "Error Rate",
			Value: strconv.FormatFloat(errorRate, 'f', 2, 64) + "%",
			Sub:   "Last 30 days",
			Trend: "neutral",
			Delta: "0%",
			Bar:   0,
		},
		{
			Key:   "p99_latency",
			Label: "p99 Latency",
			Value: strconv.FormatFloat(p99Latency, 'f', 0, 64) + "ms",
			Sub:   "Last 30 days",
			Trend: "neutral",
			Delta: "0ms",
			Bar:   0,
		},
	}

	// Fetch Recent Logs
	recentLogs, err := h.q.GetRecentRequestLogs(ctx, integratorUUID)
	logs := []LogEntry{}
	if err == nil {
		for _, rl := range recentLogs {
			logs = append(logs, LogEntry{
				Method: rl.Method,
				Path:   rl.Path,
				Status: int(rl.StatusCode),
				Ms:     int(rl.LatencyMs),
				ID:     rl.RequestID,
				Time:   timeSince(rl.CreatedAt),
			})
		}
	}

	res := AnalyticsResponse{
		Metrics: metrics,
		Logs:    logs,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func timeSince(t time.Time) string {
	d := time.Since(t)
	if d.Hours() > 24 {
		return strconv.Itoa(int(d.Hours()/24)) + "d ago"
	}
	if d.Hours() > 1 {
		return strconv.Itoa(int(d.Hours())) + "h ago"
	}
	if d.Minutes() > 1 {
		return strconv.Itoa(int(d.Minutes())) + "m ago"
	}
	return "just now"
}
