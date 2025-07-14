package redis

import (
	"fmt"
	"time"

	r "github.com/go-redis/redis"
)

// Storage ...
type Storage struct {
	client *r.Client
}

// NewStorage ... addr localhost:6379, password "", database 0
func NewStorage(addr, password string, database int) *Storage {
	client := r.NewClient(&r.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})

	return &Storage{
		client: client,
	}
}

// GetInfo ...
func (s *Storage) GetInfo() string {
	return fmt.Sprintf("Redis addr: %s db: %d", s.client.Options().Addr, s.client.Options().DB)
}

// IsReady ...
func (s *Storage) IsReady() (bool, error) {
	if err := s.client.Ping().Err(); err != nil {
		return false, fmt.Errorf("client.Ping: %w", err)
	}

	return true, nil
}

// Get ...
func (s *Storage) Get(key string) (string, error) {
	val, err := s.client.Get(key).Result()
	if err != nil {
		return "", fmt.Errorf("client.Get: %w", err)
	}

	return val, nil
}

// Save ...
func (s *Storage) Save(key string, message string, dues int) error {
	ttl := time.Duration(dues) * (24 * time.Hour)
	if err := s.client.Set(key, message, ttl).Err(); err != nil {
		return fmt.Errorf("client.Set: %w", err)
	}

	return nil
}

// Delete ...
func (s *Storage) Delete(key string) error {
	if err := s.client.Del(key).Err(); err != nil {
		return fmt.Errorf("client.Del: %w", err)
	}

	return nil
}

// IsUniq ...
func (s *Storage) IsUniq(key string) (bool, error) {
	val, err := s.client.Exists(key).Result()
	if err != nil {
		return false, fmt.Errorf("client.Exists: %w", err)
	}

	if val == 0 {
		return true, nil
	}

	return false, nil
}
