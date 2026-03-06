package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakePinger struct {
	err error
}

func (f fakePinger) Ping(_ context.Context) error {
	return f.err
}

func TestHealthz_JSONAPI(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != jsonAPIContentType {
		t.Fatalf("expected content-type %q, got %q", jsonAPIContentType, ct)
	}

	var payload jsonAPISuccess
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	if payload.Data.Type != "health" {
		t.Fatalf("expected data.type %q, got %q", "health", payload.Data.Type)
	}

	if payload.Data.ID != serviceID {
		t.Fatalf("expected data.id %q, got %q", serviceID, payload.Data.ID)
	}

	if payload.Data.Attributes["status"] != "ok" {
		t.Fatalf("expected status %q, got %#v", "ok", payload.Data.Attributes["status"])
	}
}

func TestReadyz_WhenRedisOK_Returns200_JSONAPI(t *testing.T) {
	handler := NewHandler(fakePinger{err: nil})

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != jsonAPIContentType {
		t.Fatalf("expected content-type %q, got %q", jsonAPIContentType, ct)
	}

	var payload jsonAPISuccess
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	if payload.Data.Type != "readiness" {
		t.Fatalf("expected data.type %q, got %q", "readiness", payload.Data.Type)
	}

	if payload.Data.ID != serviceID {
		t.Fatalf("expected data.id %q, got %q", serviceID, payload.Data.ID)
	}

	if payload.Data.Attributes["status"] != "ready" {
		t.Fatalf("expected status %q, got %#v", "ready", payload.Data.Attributes["status"])
	}

	if payload.Data.Attributes["redis"] != "ok" {
		t.Fatalf("expected redis %q, got %#v", "ok", payload.Data.Attributes["redis"])
	}
}

func TestReadyz_WhenRedisDown_Returns503_JSONAPIErrors(t *testing.T) {
	handler := NewHandler(fakePinger{err: errors.New("down")})

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}

	if ct := rec.Header().Get("Content-Type"); ct != jsonAPIContentType {
		t.Fatalf("expected content-type %q, got %q", jsonAPIContentType, ct)
	}

	var payload jsonAPIErrors
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	if len(payload.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(payload.Errors))
	}

	if payload.Errors[0].Status != "503" {
		t.Fatalf("expected errors[0].status %q, got %q", "503", payload.Errors[0].Status)
	}

	if payload.Errors[0].Title != "Service Unavailable" {
		t.Fatalf("expected errors[0].title %q, got %q", "Service Unavailable", payload.Errors[0].Title)
	}

	if payload.Errors[0].Detail != "Redis is unavailable" {
		t.Fatalf("expected errors[0].detail %q, got %q", "Redis is unavailable", payload.Errors[0].Detail)
	}
}

func TestReadyz_WhenRedisNotConfigured_Returns503_JSONAPIErrors(t *testing.T) {
	handler := NewHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, rec.Code)
	}

	var payload jsonAPIErrors
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("failed to parse json: %v", err)
	}

	if len(payload.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(payload.Errors))
	}

	if payload.Errors[0].Detail != "Redis client not configured" {
		t.Fatalf("expected errors[0].detail %q, got %q", "Redis client not configured", payload.Errors[0].Detail)
	}
}
