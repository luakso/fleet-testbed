package main

import (
	"io"
	"net/http"
)

// routes maps a ServeMux pattern to a handler constructor. Features that live in
// their own file add their endpoint from an init function with register, so
// independent changes never edit the same lines.
var routes = map[string]func(*Store) http.HandlerFunc{}

func register(pattern string, h func(*Store) http.HandlerFunc) {
	routes[pattern] = h
}

func init() {
	register("GET /kv/{key}", handleGet)
	register("PUT /kv/{key}", handlePut)
}

// NewHandler returns the HTTP handler serving every registered route from s.
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	for pattern, h := range routes {
		mux.HandleFunc(pattern, h(s))
	}
	return mux
}

// handleGet writes the stored value, or 404 when the key is absent.
func handleGet(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		v, ok := s.Get(r.PathValue("key"))
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		io.WriteString(w, v)
	}
}

// handlePut stores the request body as the key's value and replies 204.
func handlePut(s *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "cannot read body", http.StatusBadRequest)
			return
		}
		s.Put(r.PathValue("key"), string(body))
		w.WriteHeader(http.StatusNoContent)
	}
}
