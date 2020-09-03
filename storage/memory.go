package storage

import "time"

// MemorySecret ...
type MemorySecret struct {
	text string
	expire time.Time
}

// MemoryStorage ...
type MemoryStorage struct {
	secrets map[string]*MemorySecret
}

// Get ...
func (s *MemoryStorage) Get(key string) (string, error) {
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
func (s *MemoryStorage) Save(key string, message string, dues int) error {
	ttl := time.Duration(dues) * (24 * time.Hour)
	s.secrets[key] = &MemorySecret{
		text: message,
		expire: time.Now().Add(ttl),
	}
	return nil
}

// Delete ...
func (s *MemoryStorage) Delete(key string) error {
	delete(s.secrets, key)
	return nil
}

// IsUniq ...
func (s *MemoryStorage) IsUniq(key string) (bool, error) {
	_, ok := s.secrets[key]
	return !ok, nil

}

// NewMemoryStorage ...
func NewMemoryStorage() MemoryStorage {
	return MemoryStorage{make(map[string]*MemorySecret)}
}
