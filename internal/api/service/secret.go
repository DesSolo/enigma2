package service

import (
	"crypto/rand"
	"errors"
	"fmt"

	"enigma/internal/storage"
)

// SecretService ...
type SecretService struct {
	storage          storage.SecretStorage
	tokenLength      int
	tokenSaveRetries int
}

// NewSecretService ...
func NewSecretService(s storage.SecretStorage, tokenLength, tokenSaveRetries int) *SecretService {
	return &SecretService{
		storage:          s,
		tokenLength:      tokenLength,
		tokenSaveRetries: tokenSaveRetries,
	}
}

// Save ...
func (s *SecretService) Save(message string, dues int) (string, error) {
	token, err := s.GenerateUniqToken(s.tokenLength, s.tokenSaveRetries)
	if err != nil {
		return "", fmt.Errorf("s.GenerateUniqToken: %w", err)
	}

	// todo make hash!

	if err := s.storage.Save(token, message, dues); err != nil {
		return "", fmt.Errorf("storage.Save: %w", err)
	}

	return token, nil
}

// Get ...
func (s *SecretService) Get(key string) (string, error) {
	secret, err := s.storage.Get(key)
	if err != nil {
		return "", fmt.Errorf("storage.Get: %w", err)
	}

	// todo unhash!

	if err := s.storage.Delete(key); err != nil {
		return "", fmt.Errorf("storage.Delete: %w", err)
	}

	return secret, nil
}

// GenerateUniqToken ...
func (s *SecretService) GenerateUniqToken(length int, retries int) (string, error) {
	for i := 0; i < retries; i++ {
		token, err := generateToken(length)
		if err != nil {
			return "", fmt.Errorf("generateToken: %w", err)
		}

		uniq, err := s.storage.IsUniq(token)
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
