package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func do(t *testing.T, h http.Handler, method, path, body string) (int, string) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, path, strings.NewReader(body)))
	b, _ := io.ReadAll(rec.Result().Body)
	return rec.Code, string(b)
}

func TestGetMissingKey(t *testing.T) {
	h := NewHandler(NewStore())
	if code, _ := do(t, h, "GET", "/kv/nope", ""); code != http.StatusNotFound {
		t.Fatalf("GET missing key: got %d, want 404", code)
	}
}

func TestPutThenGet(t *testing.T) {
	h := NewHandler(NewStore())
	if code, _ := do(t, h, "PUT", "/kv/greeting", "hello"); code != http.StatusNoContent {
		t.Fatalf("PUT: got %d, want 204", code)
	}
	code, body := do(t, h, "GET", "/kv/greeting", "")
	if code != http.StatusOK || body != "hello" {
		t.Fatalf("GET: got %d %q, want 200 %q", code, body, "hello")
	}
}

func TestPutOverwrites(t *testing.T) {
	h := NewHandler(NewStore())
	do(t, h, "PUT", "/kv/k", "one")
	do(t, h, "PUT", "/kv/k", "two")
	if _, body := do(t, h, "GET", "/kv/k", ""); body != "two" {
		t.Fatalf("GET after overwrite: got %q, want %q", body, "two")
	}
}

func TestCheckLoopback(t *testing.T) {
	for addr, ok := range map[string]bool{
		"127.0.0.1:0":    true,
		"localhost:8080": true,
		"[::1]:0":        true,
		"0.0.0.0:8080":   false,
		":8080":          false,
		"10.0.0.1:80":    false,
	} {
		if err := checkLoopback(addr); (err == nil) != ok {
			t.Errorf("checkLoopback(%q) = %v, want ok=%v", addr, err, ok)
		}
	}
}
