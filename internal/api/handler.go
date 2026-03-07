// Package serviceapi contains the ogen handler implementation for the service API.
package serviceapi

import (
	"context"
	"net/http"

	"github.com/nika-piotrowska/tor-exit-nodes-api/internal/api/responses"
	"github.com/nika-piotrowska/tor-exit-nodes-api/internal/generated/oas"
)

var _ oas.Handler = (*Handler)(nil)

type redisPinger interface {
	Ping(ctx context.Context) error
}

// Handler implements the ogen-generated API handler interface.
type Handler struct {
	oas.UnimplementedHandler
	RedisPinger redisPinger
}

// GetHealthz returns the service liveness status.
func (h *Handler) GetHealthz(_ context.Context) (*oas.HealthDocument, error) {
	return responses.HealthOK(), nil
}

// GetReadyz returns the service readiness status based on Redis availability.
func (h *Handler) GetReadyz(ctx context.Context) (oas.GetReadyzRes, error) {
	if err := h.RedisPinger.Ping(ctx); err != nil {
		return nil, &oas.ErrorResponseStatusCode{
			StatusCode: http.StatusServiceUnavailable,
			Response:   responses.RedisUnavailable(),
		}
	}

	return responses.ReadyOK(), nil
}
