package storage

// SecretStorage interface
type SecretStorage interface {
	Get(key string) (string, error)
	Save(key string, message string, dues int) error
	Delete(key string) error
	IsUniq(key string) (bool, error)
}
