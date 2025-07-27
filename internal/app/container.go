package app

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/flosch/pongo2/v6"
	goredis "github.com/redis/go-redis/v9"

	"enigma/internal/api"
	"enigma/internal/config"
	"enigma/internal/pkg/hasher"
	"enigma/internal/pkg/providers/secrets"
	"enigma/internal/pkg/storage"
	"enigma/internal/pkg/storage/memory"
	"enigma/internal/pkg/storage/redis"
	"enigma/pkg/closer"
)

type container struct {
	config          *config.ServerConfig
	secretStorage   storage.SecretStorage
	hasher          hasher.Hasher
	secretsProvider *secrets.Provider
	apiServer       *api.Server

	ctx context.Context
}

func newContainer(ctx context.Context) *container {
	return &container{
		ctx: ctx,
	}
}

// Config ...
func (c *container) Config() *config.ServerConfig {
	if c.config == nil {
		configFilePath := os.Getenv("CONFIG_FILE_PATH")
		if configFilePath == "" {
			configFilePath = "/etc/enigma/config.yml"
		}

		cfg, err := config.NewServerConfigFromFile(configFilePath)
		if err != nil {
			log.Fatal("fail to load config file", configFilePath)
		}

		c.config = cfg
	}

	return c.config
}

// SecretStorage ...
func (c *container) SecretStorage() storage.SecretStorage {
	if c.secretStorage == nil {
		options := c.Config().Secrets.Storage

		var st storage.SecretStorage

		switch options.Kind {
		case "memory":
			st = memory.NewStorage()
		case "redis":
			st = redis.NewStorage(
				goredis.NewClient(&goredis.Options{
					Addr:     options.Redis.Address,
					Password: options.Redis.Password,
					DB:       options.Redis.Database,
				}),
			)
		default:
			log.Fatalf("unsupported storage kind %s", options.Kind)
		}

		for i := 0; i < options.Await.Retries; i++ {
			ready, err := st.IsReady(c.ctx)
			if ready {
				c.secretStorage = st
				closer.Add(st.Close)
				return c.secretStorage
			}

			if err != nil {
				slog.Warn("storage not ready", "err", err)
			}

			slog.Warn("await storage",
				"attempt", i+1,
				"max", options.Await.Retries,
				"await", options.Await.Interval,
			)
			time.Sleep(options.Await.Interval)
		}

		log.Fatalf("could not connect to storage after %d attempts", options.Await.Retries)
	}

	return c.secretStorage
}

// Hasher ...
func (c *container) Hasher() hasher.Hasher {
	if c.hasher == nil {
		options := c.Config().Secrets.Hasher

		switch options.Kind {
		case "aes256":
			aes, err := hasher.NewAESHasher([]byte(options.AES256.Key))
			if err != nil {
				log.Fatalf("could not create hasher: %s", err.Error())
			}

			c.hasher = aes

		default:
			log.Fatalf("hasher kind %s not supported", options.Kind)
		}
	}

	return c.hasher
}

// SecretsProvider ...
func (c *container) SecretsProvider() *secrets.Provider {
	if c.secretsProvider == nil {
		options := c.Config().Secrets

		c.secretsProvider = secrets.New(c.SecretStorage(), c.Hasher(),
			secrets.WithTokenLength(options.Token.Length),
			secrets.WithTokenSaveRetries(options.Token.SaveRetries),
		)
	}

	return c.secretsProvider
}

// APIServer ...
func (c *container) APIServer() *api.Server {
	if c.apiServer == nil {
		options := c.Config().Server

		loader, err := pongo2.NewLocalFileSystemLoader(options.TemplatesPath)
		if err != nil {
			log.Fatalf("failed to load templates path: %s", err.Error())
		}

		templateSet := pongo2.NewSet("templates", loader)

		server := api.NewServer(c.SecretsProvider(), templateSet, options.ExternalURL)

		c.apiServer = server
	}

	return c.apiServer
}
