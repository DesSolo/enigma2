package config

import (
	"enigma/internal/storage"
	"enigma/internal/storage/memory"
	"enigma/internal/storage/redis"
	"fmt"
	"os"
)

// ServerConfig ...
type ServerConfig struct {
	ListenPort      int
	TokenBytes      int
	UniqKeyRetries  int
	ResponseAddress string
	SecretStorage   storage.SecretStorage
}

// NewSeverConfig ...
func NewSeverConfig() *ServerConfig {
	ListenPort := GetEnvInt("LISTEN_PORT", 9000)
	return &ServerConfig{
		ListenPort,
		GetEnvInt("TOKEN_BYTES", 20),
		GetEnvInt("UNIQ_KEY_RETRIES", 3),
		GetEnv("RESPONSE_ADDRESS", fmt.Sprintf("http://127.0.0.1:%d", ListenPort)),
		GetEnvStorage("SECRET_STORAGE", "Memory"),
	}
}

// GetEnvStorage ...
func GetEnvStorage(key string, fault string) storage.SecretStorage {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = fault
	}
	switch value {
	case "Redis":
		return redis.NewStorage(
			GetEnv("REDIS_ADDRESS", "localhost:6379"),
			GetEnv("REDIS_PASSWORD", ""),
			GetEnvInt("REDIS_DATABASE", 0),
		)
	default:
		return memory.NewStorage()
	}
}
