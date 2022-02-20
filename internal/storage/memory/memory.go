package memory

import (
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
	mux    sync.RWMutex
}

// NewStorage ...
func NewStorage() *Storage {
	return &Storage{
		secrets: make(map[string]*secret),
	}
}

// GetInfo ...
func (s *Storage) GetInfo() string {
	return "Memory"
}

// IsReady ...
func (s *Storage) IsReady() (bool, error) {
	return true, nil
}

// Get ...
func (s *Storage) Get(key string) (string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	
	secret, ok := s.secrets[key]
	if !ok {
		return "", nil
	}

	if secret.expire.Before(time.Now()) {
		delete(s.secrets, key)
		return "", nil
	}

	return secret.text, nil
}

// Save ...
func (s *Storage) Save(key string, message string, dues int) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	ttl := time.Duration(dues) * (24 * time.Hour)
	s.secrets[key] = &secret{
		text:   message,
		expire: time.Now().Add(ttl),
	}

	return nil
}

// Delete ...
func (s *Storage) Delete(key string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.secrets, key)
	return nil
}

// IsUniq ...
func (s *Storage) IsUniq(key string) (bool, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	_, ok := s.secrets[key]
	return !ok, nil
}
