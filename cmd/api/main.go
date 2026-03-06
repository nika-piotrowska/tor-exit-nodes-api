// Package main starts the HTTP server for the tor-exit-nodes-api service.
// The service exposes endpoints used to determine whether an IP address belongs to a Tor exit node.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	redis "github.com/redis/go-redis/v9"
)

const serviceID = "service"

type pinger interface {
	Ping(ctx context.Context) error
}

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

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := run(ctx)
	stop()

	if err != nil {
		log.Printf("fatal: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	rdb, err := newRedisClient(redisURLFromEnv())
	if err != nil {
		return fmt.Errorf("failed to create redis client: %w", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Printf("failed to close redis client: %v", err)
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
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		return srv.Shutdown(shutdownCtx)

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
