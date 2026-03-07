// Package redisclient provides helpers for configuring and pinging Redis.
package redisclient

import (
	"context"
	"errors"
	"os"

	redis "github.com/redis/go-redis/v9"
)

// Pinger wraps a Redis client and exposes a ping operation.
type Pinger struct {
	Client *redis.Client
}

// Ping checks whether Redis is reachable.
func (p Pinger) Ping(ctx context.Context) error {
	return p.Client.Ping(ctx).Err()
}

// URLFromEnv returns the Redis connection URL from the REDIS_URL environment variable.
func URLFromEnv() (string, error) {
	raw := os.Getenv("REDIS_URL")
	if raw == "" {
		return "", errors.New("REDIS_URL is required")
	}

	return raw, nil
}

// New builds a Redis client from the provided Redis connection URL.
func New(raw string) (*redis.Client, error) {
	opts, err := redis.ParseURL(raw)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}
