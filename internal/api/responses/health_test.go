package responses

import "testing"

func TestHealthOK(t *testing.T) {
	response := HealthOK()

	if response.Data.Type != "health" {
		t.Fatalf("expected type health, got %q", response.Data.Type)
	}

	if response.Data.ID != "service" {
		t.Fatalf("expected id service, got %q", response.Data.ID)
	}

	if response.Data.Attributes.Status != "ok" {
		t.Fatalf("expected status ok, got %q", response.Data.Attributes.Status)
	}
}

func TestReadyOK(t *testing.T) {
	response := ReadyOK()

	if response.Data.Type != "readiness" {
		t.Fatalf("expected type readiness, got %q", response.Data.Type)
	}

	if response.Data.ID != "service" {
		t.Fatalf("expected id service, got %q", response.Data.ID)
	}

	if response.Data.Attributes.Status != "ready" {
		t.Fatalf("expected status ready, got %q", response.Data.Attributes.Status)
	}

	if response.Data.Attributes.Redis != "ok" {
		t.Fatalf("expected redis ok, got %q", response.Data.Attributes.Redis)
	}
}
