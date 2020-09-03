package main

import (
	"time"
	"enigma/config"
	"enigma/storage"
	"fmt"
	"log"
	"net/http"
)

const ( 
	awaitStorageRetries = 10 // wait resties storage
	awaitStorageSleep = 1 // sleep in seconds betwin wait attempts
)

// Config ...
var Config = config.New()

// GetSecret ...
func GetSecret(s storage.SecretStorage, key string) (string, error) {
	secret, err := s.Get(key)
	if err != nil {
		return "", err
	}
	if err := s.Delete(key); err != nil {
		return "", err
	}
	return secret, nil
}

// SaveSecret ...
func SaveSecret(s storage.SecretStorage, message string, dues int) (string, error) {
	token, err := GenerateUniqToken(s, Config.UniqKeyRetries)
	if err != nil {
		return "", nil
	}
	if err := s.Save(token, message, dues); err != nil {
		return "" , err
	}
	return token, nil
}

func awaitSecretStorage(s storage.SecretStorage) error {
	for i:=0; i<awaitStorageRetries; i++ {
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

func main() {
	if err := awaitSecretStorage(Config.SecretStorage); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", methodMiddleware(indexHandler, "GET"))
	http.HandleFunc("/post/", methodMiddleware(createHandler, "POST"))
	http.HandleFunc("/get/", methodMiddleware(viewHandler, "GET"))

	log.Printf(
		"service started port: %d response_address: %s token_bytes: %d\n",
		Config.ListenPort, Config.ResponseAddress, Config.TokenBytes,
		)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", Config.ListenPort), nil); err != nil {
		log.Fatal(err)
	}
}
