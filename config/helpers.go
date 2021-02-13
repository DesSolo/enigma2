package config

import (
	"log"
	"os"
	"strconv"
)

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
