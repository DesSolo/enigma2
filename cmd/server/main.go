package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
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

var version = "local"

const banner = `
 _______ __   _ _____  ______ _______ _______
 |______ | \  |   |   |  ____ |  |  | |_____|
 |______ |  \_| __|__ |_____| |  |  | |     |

 version: %s

`

func loadServerConfig() (*config.ServerConfig, error) {
	return config.NewServerConfigFromFile() // nolint:wrapcheck
}

func configureLogger(c *config.ServerConfig) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.Level(c.Logging.Level),
	})

	slog.SetDefault(
		slog.New(handler),
	)
}

func loadSecretStorage(c *config.ServerConfig) (storage.SecretStorage, error) {
	var st storage.SecretStorage

	switch c.Secrets.Storage.Kind {
	case "memory":
		st = memory.NewStorage()
	case "redis":
		st = redis.NewStorage(
			goredis.NewClient(&goredis.Options{
				Addr:     c.Secrets.Storage.Redis.Address,
				Password: c.Secrets.Storage.Redis.Password,
				DB:       c.Secrets.Storage.Redis.Database,
			}),
		)
	default:
		return nil, fmt.Errorf("storage kind: %s not supported", c.Secrets.Storage.Kind)
	}

	for i := 0; i < c.Secrets.Storage.Await.Retries; i++ {
		ready, err := st.IsReady(context.Background())
		if ready {
			return st, nil
		}
		if err != nil {
			slog.Warn("storage not ready", "err", err)
		}

		slog.Warn("await storage",
			"attempt", i+1,
			"max", c.Secrets.Storage.Await.Retries,
			"await", c.Secrets.Storage.Await.Interval,
		)
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
		return nil, fmt.Errorf("failed to read index template err: %w", err)
	}

	viewSecretTemplate, err := os.ReadFile(
		path.Join(c.Server.TemplatesPath, viewSecretTemplateFileName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read view secret template err: %w", err)
	}

	server.LoadHandlers(indexTemplate, viewSecretTemplate, c.Server.ExternalURL)

	return server, nil
}

func main() {
	serverConfig, err := loadServerConfig()
	if err != nil {
		log.Fatalf("failed to load config err: %s", err.Error())
	}

	configureLogger(serverConfig)

	secretStorage, err := loadSecretStorage(serverConfig)
	if err != nil {
		log.Fatalf("failed to load secret storage err: %s", err.Error())
	}

	secretsHasher, err := loadHasher(serverConfig)
	if err != nil {
		log.Fatalf("failed to load hasher err: %s", err.Error())
	}

	apiServer, err := loadAPIServer(serverConfig, secretsHasher, secretStorage)
	if err != nil {
		log.Fatalf("failed to load api server err: %s", err.Error())
	}

	fmt.Printf(banner, version)

	slog.Info("server running on", "addr", serverConfig.Server.Bind)

	if err := apiServer.Run(serverConfig.Server.Bind); err != nil {
		log.Fatalf("failed to run api server err: %s", err.Error())
	}
}
