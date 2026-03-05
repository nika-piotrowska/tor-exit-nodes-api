package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpEndpoint(t *testing.T) {
	handler := buildHandler()
	req := httptest.NewRequest(http.MethodGet, "/up", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if got := rec.Body.String(); got != "ok" {
		t.Fatalf("expected body %q, got %q", "ok", got)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("expected content-type %q, got %q", "text/plain; charset=utf-8", ct)
	}
}
