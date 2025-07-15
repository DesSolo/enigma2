package secrets

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"enigma/internal/pkg/hasher"
	"enigma/internal/pkg/storage"
)

const (
	defaultTokenLength      = 32
	defaultTokenSaveRetries = 3
)

// Provider ...
type Provider struct {
	storage storage.SecretStorage
	hasher  hasher.Hasher

	tokenLength      int
	tokenSaveRetries int
}

// New ...
func New(storage storage.SecretStorage, hasher hasher.Hasher, options ...OptionFunc) *Provider {
	p := &Provider{
		storage: storage,
		hasher:  hasher,

		tokenLength:      defaultTokenLength,
		tokenSaveRetries: defaultTokenSaveRetries,
	}

	for _, option := range options {
		option(p)
	}

	return p
}

// SaveSecret ...
func (p *Provider) SaveSecret(ctx context.Context, message string, dues int) (string, error) {
	token, err := p.generateUniqToken(ctx, p.tokenLength, p.tokenSaveRetries)
	if err != nil {
		return "", fmt.Errorf("s.GenerateUniqToken: %w", err)
	}

	encrypted, err := p.hasher.Encrypt(message)
	if err != nil {
		return "", fmt.Errorf("s.Encrypt: %w", err)
	}

	ttl := time.Duration(dues) * (24 * time.Hour)
	if err := p.storage.Save(ctx, token, encrypted, ttl); err != nil {
		return "", fmt.Errorf("storage.Save: %w", err)
	}

	return token, nil
}

// GetSecret ...
func (p *Provider) GetSecret(ctx context.Context, key string) (string, error) {
	secret, err := p.storage.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("storage.Get: %w", err)
	}

	decrypted, err := p.hasher.Decrypt(secret)
	if err != nil {
		return "", fmt.Errorf("p.hasher.Decrypt: %w", err)
	}

	if err := p.storage.Delete(ctx, key); err != nil {
		return "", fmt.Errorf("storage.Delete: %w", err)
	}

	return decrypted, nil
}

func (p *Provider) generateUniqToken(ctx context.Context, length int, retries int) (string, error) {
	for i := 0; i < retries; i++ {
		token, err := generateToken(length)
		if err != nil {
			return "", fmt.Errorf("generateToken: %w", err)
		}

		uniq, err := p.storage.IsUniq(ctx, token)
		if err != nil {
			return "", fmt.Errorf("storage.IsUniq: %w", err)
		}

		if uniq {
			return token, nil
		}
	}

	return "", errors.New("maximum retries save")
}

func generateToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}

	return fmt.Sprintf("%x", b), nil
}
