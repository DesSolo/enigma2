package config

import (
	"enigma/storage"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config ...
type Config struct {
	ListenPort int
	TokenBytes int
	UniqKeyRetries int
	ResponseAddress string
	SecretStorage storage.SecretStorage
}

// GetEnv ...
func GetEnv(key string, fault string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fault
	}
	return value

}

// GetEnvInt ...
func GetEnvInt(key string, fault int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fault
	}
	v, err := strconv.Atoi(key)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

// GetEnvStorage ...
func GetEnvStorage(key string, fault string) storage.SecretStorage {
	value := os.Getenv(key)
	if len(value) == 0 {
		value = fault
	}
	switch value {
	case "Redis":
		{
			st := storage.NewRedisStorage(
				GetEnv("REDIS_ADDRESS", "localhost:6379"),
				GetEnv("REDIS_PASSWORD", ""),
				GetEnvInt("REDIS_DATABASE", 0),
			)
			return &st
		}
	default:
		{
			st := storage.NewMemoryStorage()
			return &st
		}
	}
}

// New ...
func New() *Config {
	ListenPort := GetEnvInt("LISTEN_PORT", 9000)
	return &Config{
		ListenPort,
		GetEnvInt("TOKEN_BYTES", 20),
		GetEnvInt("UNIQ_KEY_RETRIES", 3),
		GetEnv("RESPONSE_ADDRESS", fmt.Sprintf("http://127.0.0.1:%d", ListenPort)),
		GetEnvStorage("SECRET_STORAGE", "Memory"),
	}
}