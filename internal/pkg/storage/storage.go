//go:generate mockery --case=snake --with-expecter --name SecretStorage

package storage

import "context"

// SecretStorage interface
type SecretStorage interface {
	GetInfo(ctx context.Context) string
	IsReady(ctx context.Context) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Save(ctx context.Context, key string, message string, dues int) error
	Delete(ctx context.Context, key string) error
	IsUniq(ctx context.Context, key string) (bool, error)
}
