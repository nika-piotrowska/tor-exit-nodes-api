package responses

import "testing"

func TestRedisUnavailable(t *testing.T) {
	response := RedisUnavailable()

	if len(response.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(response.Errors))
	}

	errObj := response.Errors[0]

	if errObj.Status != "503" {
		t.Fatalf("expected status 503, got %q", errObj.Status)
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
