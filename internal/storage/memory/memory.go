package memory

import "time"

type secret struct {
	text   string
	expire time.Time
}

// Storage ...
type Storage struct {
	secrets map[string]*secret
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
	ttl := time.Duration(dues) * (24 * time.Hour)
	s.secrets[key] = &secret{
		text:   message,
		expire: time.Now().Add(ttl),
	}

	return nil
}

// Delete ...
func (s *Storage) Delete(key string) error {
	delete(s.secrets, key)
	return nil
}

// IsUniq ...
func (s *Storage) IsUniq(key string) (bool, error) {
	_, ok := s.secrets[key]
	return !ok, nil
}
