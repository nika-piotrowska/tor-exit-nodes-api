// Package main starts the HTTP server for the tor-exit-nodes-api service.
// The service exposes endpoints used to determine whether an IP address belongs to a Tor exit node.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const (
	jsonAPIContentType = "application/vnd.api+json"
	serviceID          = "service"
)

type jsonAPIData struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

type jsonAPISuccess struct {
	Data jsonAPIData `json:"data"`
}

type jsonAPIError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail,omitempty"`
}

type jsonAPIErrors struct {
	Errors []jsonAPIError `json:"errors"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", jsonAPIContentType)
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(v); err != nil {
		log.Printf("failed to encode json response: %v", err)
	}
}

func writeJSONAPIError(w http.ResponseWriter, status int, title, detail string) {
	writeJSON(w, status, jsonAPIErrors{
		Errors: []jsonAPIError{
			{
				Status: strconv.Itoa(status),
				Title:  title,
				Detail: detail,
			},
		},
	})
}

// pinger represents a dependency we can health-check (e.g. Redis).
// Using an interface keeps handlers easy to test.
type pinger interface {
	Ping(ctx context.Context) error
}

// goRedisPinger adapts redis.Client to our pinger interface.
type goRedisPinger struct {
	client *redis.Client
}

func (p goRedisPinger) Ping(ctx context.Context) error {
	return p.client.Ping(ctx).Err()
}

func buildHandler(redisPinger pinger) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", healthHandler)
	mux.HandleFunc("/readyz", readinessHandler(redisPinger))

	return mux
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, jsonAPISuccess{
		Data: jsonAPIData{
			Type: "health",
			ID:   serviceID,
			Attributes: map[string]any{
				"status": "ok",
			},
		},
	})
}

func readinessHandler(redisPinger pinger) http.HandlerFunc {
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

		writeJSON(w, http.StatusOK, jsonAPISuccess{
			Data: jsonAPIData{
				Type: "readiness",
				ID:   serviceID,
				Attributes: map[string]any{
					"status": "ready",
					"redis":  "ok",
				},
			},
		})
	}
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err := run(ctx); err != nil {
		log.Printf("fatal: %v", err)
		stop()
		os.Exit(1)
	}

	stop()
}

func run(ctx context.Context) error {
	rdb, err := newRedisClient(redisURLFromEnv())
	if err != nil {
		return fmt.Errorf("failed to create redis client: %w", err)
	}
	defer func() {
		if closeErr := rdb.Close(); closeErr != nil {
			log.Printf("failed to close redis client: %v", closeErr)
		}
	}()

	srv := &http.Server{
		Addr:              ":3000",
		Handler:           buildHandler(goRedisPinger{client: rdb}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("server listening on %s", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown failed: %w", err)
		}

		return <-serverErr
	case err := <-serverErr:
		return err
	}
}

func redisURLFromEnv() string {
	raw := os.Getenv("REDIS_URL")
	if raw == "" {
		return "redis://localhost:6380/0"
	}

	return raw
}

func newRedisClient(raw string) (*redis.Client, error) {
	opts, err := redis.ParseURL(raw)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}
