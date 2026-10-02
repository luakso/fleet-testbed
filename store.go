package main

import "sync"

// Store is an in-memory key-value store that is safe for concurrent use.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// NewStore returns an empty Store.
func NewStore() *Store {
	return &Store{data: make(map[string]string)}
}

// Get returns the value for key and whether it was present.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

// Put sets key to value, replacing any previous value.
func (s *Store) Put(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}
