// Package httpapi provides HTTP handlers for health and readiness endpoints.
package httpapi

import (
	"context"
	"net/http"
	"time"
)

// Pinger represents a dependency that can verify connectivity, for example to Redis.
type Pinger interface {
	Ping(ctx context.Context) error
}

// NewHandler builds an HTTP handler with health and readiness endpoints.
func NewHandler(redisPinger Pinger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/readyz", readinessHandler(redisPinger))

	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSONAPISuccess(
		w,
		http.StatusOK,
		"health",
		serviceID,
		map[string]any{
			"status": "ok",
		},
	)
}

func readinessHandler(redisPinger Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if redisPinger == nil {
			writeJSONAPIError(w, http.StatusServiceUnavailable, "Service Unavailable", "Redis client not configured")
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
		defer cancel()

		if err := redisPinger.Ping(ctx); err != nil {
			writeJSONAPIError(w, http.StatusServiceUnavailable, "Service Unavailable", "Redis is unavailable")
			return
		}

		writeJSONAPISuccess(
			w,
			http.StatusOK,
			"readiness",
			serviceID,
			map[string]any{
				"status": "ready",
				"redis":  "ok",
			},
		)
	}
}
