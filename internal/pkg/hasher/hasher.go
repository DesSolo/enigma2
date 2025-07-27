//go:generate mockery --case=snake --with-expecter --name Hasher

package hasher

// Hasher ...
type Hasher interface {
	Encrypt(string) (string, error)
	Decrypt(string) (string, error)
}
