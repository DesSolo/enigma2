//go:generate mockery --case=snake --with-expecter --name SecretStorage

package storage

import (
	"context"
	"time"
)

// SecretStorage interface
type SecretStorage interface {
	IsReady(ctx context.Context) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Save(ctx context.Context, key string, message string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	IsUniq(ctx context.Context, key string) (bool, error)
}
