package main

import (
	"enigma/config"
	"enigma/storage"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	awaitStorageRetries = 10 // wait resties storage
	awaitStorageSleep   = 1  // sleep in seconds betwin wait attempts
)

// Config ... global config variable for main pakage
var Config = config.New()

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
