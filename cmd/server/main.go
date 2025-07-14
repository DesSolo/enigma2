package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"enigma/internal/api"
	"enigma/internal/api/service"
	"enigma/internal/config"
	"enigma/internal/storage"
	"enigma/internal/storage/memory"
	"enigma/internal/storage/redis"
)

const (
	indexTemplateFileName      = "index.html"
	viewSecretTemplateFileName = "view_secret.html"
)

func loadServerConfig() (*config.ServerConfig, error) {
	return config.NewServerConfigFromFile() // nolint:wrapcheck
}

func loadSecretStorage(c *config.ServerConfig) (storage.SecretStorage, error) {
	var st storage.SecretStorage

	switch c.Secrets.Storage.Type {
	case "memory":
		st = memory.NewStorage()
	case "redis":
		st = redis.NewStorage(
			c.Redis.Address,
			c.Redis.Password,
			c.Redis.Database,
		)
	default:
		return nil, fmt.Errorf("storage type: %s not supported", c.Secrets.Storage.Type)
	}

	for i := 0; i < c.Secrets.Storage.Await.Retries; i++ {
		ready, err := st.IsReady()
		if ready {
			return st, nil
		}
		if err != nil {
			log.Printf("storage error: %s", err.Error())
		}

		log.Printf("await storage attempt: %d/%d sleep: %s", i+1, c.Secrets.Storage.Await.Retries, c.Secrets.Storage.Await.Interval)
		time.Sleep(c.Secrets.Storage.Await.Interval)
	}

	return nil, fmt.Errorf("could not connect to storage after %d attempts", c.Secrets.Storage.Await.Retries)
}

func loadAPIServer(c *config.ServerConfig, s storage.SecretStorage) (*api.Server, error) {
	secretService := service.NewSecretService(s, c.Secrets.Token.Length, c.Secrets.Token.SaveRetries)
	server := api.NewServer(secretService)

	indexTemplate, err := os.ReadFile(
		path.Join(c.Server.TemplatesPath, indexTemplateFileName),
	)
	if err != nil {
		return nil, fmt.Errorf("fault read index template err: %w", err)
	}

	viewSecretTemplate, err := os.ReadFile(
		path.Join(c.Server.TemplatesPath, viewSecretTemplateFileName),
	)
	if err != nil {
		return nil, fmt.Errorf("fault read view secret template err: %w", err)
	}

	server.LoadHandlers(indexTemplate, viewSecretTemplate, c.Server.ExternalURL)

	return server, nil
}

func main() {
	serverConfig, err := loadServerConfig()
	if err != nil {
		log.Fatalf("fault load config err: %s", err.Error())
	}

	secretStorage, err := loadSecretStorage(serverConfig)
	if err != nil {
		log.Fatalf("fault load secret storage err: %s", err.Error())
	}

	apiServer, err := loadAPIServer(serverConfig, secretStorage)
	if err != nil {
		log.Fatalf("fault load api server err: %s", err.Error())
	}

	if err := apiServer.Run(serverConfig.Server.Bind); err != nil {
		log.Fatalf("fault run api server err: %s", err.Error())
	}
}
