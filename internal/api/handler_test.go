package serviceapi

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/nika-piotrowska/tor-exit-nodes-api/internal/generated/oas"
)

type stubRedisPinger struct {
	err error
}

func (s stubRedisPinger) Ping(_ context.Context) error {
	return s.err
}

func TestGetHealthz(t *testing.T) {
	handler := &Handler{}

	response, err := handler.GetHealthz(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if response.Data.Type != "health" {
		t.Fatalf("expected type health, got %q", response.Data.Type)
	}

	if response.Data.Attributes.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Data.Attributes.Status)
	}
}

func TestGetReadyzReturnsReadyWhenRedisIsAvailable(t *testing.T) {
	handler := &Handler{
		RedisPinger: stubRedisPinger{},
	}

	response, err := handler.GetReadyz(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	readyResponse, ok := response.(*oas.ReadinessDocument)
	if !ok {
		t.Fatalf("expected *oas.ReadinessDocument, got %T", response)
	}

	if readyResponse.Data.Type != "readiness" {
		t.Fatalf("expected type readiness, got %q", readyResponse.Data.Type)
	}

	if readyResponse.Data.Attributes.Status != "ready" {
		t.Fatalf("expected status ready, got %q", readyResponse.Data.Attributes.Status)
	}

	if readyResponse.Data.Attributes.Redis != "ok" {
		t.Fatalf("expected redis ok, got %q", readyResponse.Data.Attributes.Redis)
	}
}

func TestGetReadyzReturnsServiceUnavailableWhenRedisFails(t *testing.T) {
	handler := &Handler{
		RedisPinger: stubRedisPinger{err: errors.New("redis down")},
	}

	response, err := handler.GetReadyz(context.Background())
	if response != nil {
		t.Fatalf("expected nil response, got %T", response)
	}

	var statusErr *oas.ErrorResponseStatusCode
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected *oas.ErrorResponseStatusCode, got %T", err)
	}

	if statusErr.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected status code %d, got %d", http.StatusServiceUnavailable, statusErr.StatusCode)
	}

	if len(statusErr.Response.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(statusErr.Response.Errors))
	}

	errObj := statusErr.Response.Errors[0]

	if errObj.Status != "503" {
		t.Fatalf("expected error status 503, got %q", errObj.Status)
	}

	if errObj.Title != "Service Unavailable" {
		t.Fatalf("expected title Service Unavailable, got %q", errObj.Title)
	}

	if !errObj.Detail.Set {
		t.Fatal("expected detail to be set")
	}

	if errObj.Detail.Value != "Redis is unavailable" {
		t.Fatalf("expected detail Redis is unavailable, got %q", errObj.Detail.Value)
	}
}
