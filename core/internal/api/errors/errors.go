package errors

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorResponse represents a standardized API error shape.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// WriteError sends a JSON formatted error response.
func WriteError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    code,
		Message: message,
	})
}

// LogAndWriteError logs the underlying error using slog and then sends a JSON formatted error response.
func LogAndWriteError(w http.ResponseWriter, status int, code, message string, err error) {
	if err != nil {
		slog.Error("api request failed", slog.Int("status", status), slog.String("code", code), slog.Any("error", err))
	}
	WriteError(w, status, code, message)
}
