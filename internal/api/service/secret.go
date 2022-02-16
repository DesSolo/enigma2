package service

import (
	"crypto/rand"
	"enigma/internal/storage"
	"errors"
	"fmt"
)

type SecretService struct {
	storage          storage.SecretStorage
	tokenLenght      int
	tokenSaveRetries int
}

func NewSecretService(s storage.SecretStorage, tokenLenght, tokenSaveRetries int) *SecretService {
	return &SecretService{
		storage:          s,
		tokenLenght:      tokenLenght,
		tokenSaveRetries: tokenSaveRetries,
	}
}

func (s *SecretService) Save(message string, dues int) (string, error) {
	token, err := s.GenerateUniqToken(s.tokenLenght, s.tokenSaveRetries)
	if err != nil {
		return "", err
	}

	// todo mkae hash!

	if err := s.storage.Save(token, message, dues); err != nil {
		return "", err
	}

	return token, nil
}

func (s *SecretService) Get(key string) (string, error) {
	secret, err := s.storage.Get(key)
	if err != nil {
		return "", err
	}

	// todo unhash!

	if err := s.storage.Delete(key); err != nil {
		return "", err
	}

	return secret, nil
}

func (s *SecretService) GenerateUniqToken(length int, retries int) (string, error) {
	for i := 0; i < retries; i++ {
		token, err := generateToken(length)
		if err != nil {
			return "", err
		}

		uniq, err := s.storage.IsUniq(token)
		if err != nil {
			return "", err
		}

		if uniq {
			return token, nil
		}
	}

	return "", errors.New("maximum retries save")
}

func generateToken(lenght int) (string, error) {
	b := make([]byte, lenght)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
