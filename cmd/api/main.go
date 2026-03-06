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

	"github.com/joho/godotenv"

	httpapi "github.com/nika-piotrowska/tor-exit-nodes-api/internal/httpapi"
	"github.com/nika-piotrowska/tor-exit-nodes-api/internal/redisclient"
)

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
	if err := loadEnv(); err != nil {
		return err
	}

	redisURL, err := redisclient.URLFromEnv()
	if err != nil {
		return fmt.Errorf("failed to read redis configuration: %w", err)
	}

	rdb, err := redisclient.New(redisURL)
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
		Handler:           httpapi.NewHandler(redisclient.Pinger{Client: rdb}),
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

func loadEnv() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to load .env: %w", err)
	}

	return nil
}
