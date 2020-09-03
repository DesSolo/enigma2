package storage

import (
	"time"
	"context"
	r"github.com/go-redis/redis"
)

var ctx = context.Background()

// RedisStorage ...
type RedisStorage struct {
	client *r.Client
}

// IsReady ...
func (s *RedisStorage) IsReady() (bool, error) {
	if err := s.client.Ping(ctx).Err(); err != nil {
		return false, err
	}
	return true, nil
}

// Get ...
func (s *RedisStorage) Get(key string) (string, error) {
	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Save ...
func (s *RedisStorage) Save(key string, message string, dues int) error {
	ttl := time.Duration(dues) * (24 * time.Hour)
	if err := s.client.Set(ctx, key, message, ttl).Err(); err != nil {
		return err
	}
	return nil
}

// Delete ...
func (s *RedisStorage) Delete(key string) error {
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return err
	}
	return nil
}

// IsUniq ...
func (s *RedisStorage) IsUniq(key string) (bool, error) {
	val, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if val == 0 {
		return true, nil
	}
	return false, nil
}

// NewRedisStorage ... addr localhost:6379, password "", database 0
func NewRedisStorage(addr, password string, database int) RedisStorage {
	client := r.NewClient(&r.Options{
        Addr:     addr,
        Password: password,
        DB:       database,
    })
	return RedisStorage{client}
}