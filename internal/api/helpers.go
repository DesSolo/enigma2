package api

import (
	"crypto/rand"
	"enigma/internal/storage"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

const (
	awaitStorageRetries = 10 // wait resties storage
	awaitStorageSleep   = 1  // sleep in seconds betwin wait attempts
)

func awaitSecretStorage(s storage.SecretStorage) error {
	for i := 0; i < awaitStorageRetries; i++ {
		ready, err := s.IsReady()
		if ready {
			return nil
		}

		if err != nil {
			log.Println("storage error:", err)
		}

		log.Printf("await storage attempt: %d/%d sleeping: %d sec", i+1, awaitStorageRetries, awaitStorageSleep)
		time.Sleep(awaitStorageSleep * time.Second)
	}
	return fmt.Errorf("could not connect to storage after %d attempts", awaitStorageRetries)
}

func genRegexpGetView(tokenBytes int) *regexp.Regexp {
	return regexp.MustCompile(
		fmt.Sprintf("^/get/(\\w{%d})$", tokenBytes*2),
	)
}

func raiseError(w http.ResponseWriter, method int) {
	http.Error(w, http.StatusText(method), method)
}

func loadTemplate(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("error load template file: %s err: %s", file, err)
	}

	return data
}

func generateUniqToken(retries int) (string, error) {
	for i := 0; ; i++ {
		key, err := generateToken(sConfig.TokenBytes)
		if err != nil {
			return "", err
		}

		uniq, err := sStorage.IsUniq(key)
		if err != nil {
			return "", err
		}

		if uniq == true {
			return key, nil
		}

		if i >= (retries - 1) {
			break
		}
	}

	return "", errors.New("maximum retries save")
}

func generateToken(bytes int) (string, error) {
	b := make([]byte, bytes)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func getSecret(s storage.SecretStorage, key string) (string, error) {
	secret, err := s.Get(key)
	if err != nil {
		return "", err
	}

	if err := s.Delete(key); err != nil {
		return "", err
	}

	return secret, nil
}

func saveSecret(s storage.SecretStorage, message string, dues int) (string, error) {
	token, err := generateUniqToken(sConfig.UniqKeyRetries)
	if err != nil {
		return "", err
	}

	if err := s.Save(token, message, dues); err != nil {
		return "", err
	}

	return token, nil
}
