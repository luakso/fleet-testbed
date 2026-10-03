package main

import (
	"io"
	"net/http"
)

func init() {
	register("GET /healthz", handleHealthz)
}

// handleHealthz replies 200 with body "ok" so callers can tell the server is up.
func handleHealthz(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "ok")
	}
}
