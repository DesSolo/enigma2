package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"enigma/internal/api"
	"enigma/internal/config"
	"enigma/internal/pkg/hasher"
	"enigma/internal/pkg/providers/secrets"
	"enigma/internal/pkg/storage"
	"enigma/internal/pkg/storage/memory"
	"enigma/internal/pkg/storage/redis"

	goredis "github.com/redis/go-redis/v9"
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
			goredis.NewClient(&goredis.Options{
				Addr:     c.Redis.Address,
				Password: c.Redis.Password,
				DB:       c.Redis.Database,
			}),
		)
	default:
		return nil, fmt.Errorf("storage type: %s not supported", c.Secrets.Storage.Type)
	}

	for i := 0; i < c.Secrets.Storage.Await.Retries; i++ {
		ready, err := st.IsReady(context.Background())
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

func loadHasher(c *config.ServerConfig) (hasher.Hasher, error) {
	switch c.Secrets.Hasher.Kind {
	case "aes256":
		aes, err := hasher.NewAESHasher([]byte(c.Secrets.Hasher.AES256.Key))
		if err != nil {
			return nil, fmt.Errorf("could not create hasher: %s", err.Error())
		}

		return aes, nil

	default:
		return nil, fmt.Errorf("hasher kind %s not supported", c.Secrets.Hasher.Kind)
	}
}

func loadAPIServer(c *config.ServerConfig, h hasher.Hasher, s storage.SecretStorage) (*api.Server, error) {
	secretService := secrets.New(s, h,
		secrets.WithTokenLength(c.Secrets.Token.Length),
		secrets.WithTokenSaveRetries(c.Secrets.Token.SaveRetries),
	)
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

	secretsHasher, err := loadHasher(serverConfig)
	if err != nil {
		log.Fatalf("fault load hasher err: %s", err.Error())
	}

	apiServer, err := loadAPIServer(serverConfig, secretsHasher, secretStorage)
	if err != nil {
		log.Fatalf("fault load api server err: %s", err.Error())
	}

	if err := apiServer.Run(serverConfig.Server.Bind); err != nil {
		log.Fatalf("fault run api server err: %s", err.Error())
	}
}
