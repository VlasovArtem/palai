// Package httpapi wires Palai's HTTP routes using the standard library mux.
package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// NewHandler returns the root HTTP handler with all routes registered.
func NewHandler(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealthz(logger))
	return mux
}

type healthResponse struct {
	Status string `json:"status"`
}

func handleHealthz(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, r, logger, http.StatusOK, healthResponse{Status: "ok"})
	}
}

// writeJSON encodes v as the response body. Encoding failures are logged
// because the status line has already been sent.
func writeJSON(w http.ResponseWriter, r *http.Request, logger *slog.Logger, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		logger.ErrorContext(r.Context(), "write json response",
			slog.String("path", r.URL.Path),
			slog.Any("error", err),
		)
	}
}
