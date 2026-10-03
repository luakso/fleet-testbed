package main

import (
	"net/http"
	"testing"
)

func TestHealthz(t *testing.T) {
	h := NewHandler(NewStore())
	code, body := do(t, h, "GET", "/healthz", "")
	if code != http.StatusOK || body != "ok" {
		t.Fatalf("GET /healthz: got %d %q, want 200 %q", code, body, "ok")
	}
}
