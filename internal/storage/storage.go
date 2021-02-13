package storage

// SecretStorage interface
type SecretStorage interface {
	IsReady() (bool, error)
	Get(key string) (string, error)
	Save(key string, message string, dues int) error
	Delete(key string) error
	IsUniq(key string) (bool, error)
}
