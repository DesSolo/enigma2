package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Storage ...
type Storage struct {
	client *redis.Client
}

// NewStorage ... addr localhost:6379, password "", database 0
func NewStorage(addr, password string, database int) *Storage {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})

	return &Storage{
		client: client,
	}
}

// GetInfo ...
func (s *Storage) GetInfo(_ context.Context) string {
	return fmt.Sprintf("Redis addr: %s db: %d", s.client.Options().Addr, s.client.Options().DB)
}

// IsReady ...
func (s *Storage) IsReady(ctx context.Context) (bool, error) {
	if err := s.client.Ping(ctx).Err(); err != nil {
		return false, fmt.Errorf("client.Ping: %w", err)
	}

	return true, nil
}

// Get ...
func (s *Storage) Get(ctx context.Context, key string) (string, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("client.Get: %w", err)
	}

	return val, nil
}

// Save ...
func (s *Storage) Save(ctx context.Context, key string, message string, ttl time.Duration) error {
	if err := s.client.Set(ctx, key, message, ttl).Err(); err != nil {
		return fmt.Errorf("client.Set: %w", err)
	}

	return nil
}

// Delete ...
func (s *Storage) Delete(ctx context.Context, key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("client.Del: %w", err)
	}

	return nil
}

// IsUniq ...
func (s *Storage) IsUniq(ctx context.Context, key string) (bool, error) {
	val, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("client.Exists: %w", err)
	}

	if val == 0 {
		return true, nil
	}

	return false, nil
}
