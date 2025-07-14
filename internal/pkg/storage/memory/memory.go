package memory

import (
	"context"
	"sync"
	"time"
)

type secret struct {
	text   string
	expire time.Time
}

// Storage ...
type Storage struct {
	secrets map[string]*secret
	mux     sync.RWMutex
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		secrets: make(map[string]*secret),
	}
}

// IsReady ...
func (s *Storage) IsReady(_ context.Context) (bool, error) {
	return true, nil
}

// Get ...
func (s *Storage) Get(_ context.Context, key string) (string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	secret, ok := s.secrets[key]
	if !ok {
		return "", nil
	}

	if secret.expire.Before(time.Now().UTC()) {
		delete(s.secrets, key)
		return "", nil
	}

	return secret.text, nil
}

// Save ...
func (s *Storage) Save(_ context.Context, key string, message string, ttl time.Duration) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.secrets[key] = &secret{
		text:   message,
		expire: time.Now().UTC().Add(ttl),
	}

	return nil
}

// Delete ...
func (s *Storage) Delete(_ context.Context, key string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.secrets, key)
	return nil
}

// IsUniq ...
func (s *Storage) IsUniq(_ context.Context, key string) (bool, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	_, ok := s.secrets[key]
	return !ok, nil
}
