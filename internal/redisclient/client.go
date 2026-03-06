package redisclient

import (
	"context"
	"errors"
	"os"

	redis "github.com/redis/go-redis/v9"
)

type Pinger struct {
	Client *redis.Client
}

func (p Pinger) Ping(ctx context.Context) error {
	return p.Client.Ping(ctx).Err()
}

func URLFromEnv() (string, error) {
	raw := os.Getenv("REDIS_URL")
	if raw == "" {
		return "", errors.New("REDIS_URL is required")
	}

	return raw, nil
}

func New(raw string) (*redis.Client, error) {
	opts, err := redis.ParseURL(raw)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opts), nil
}
