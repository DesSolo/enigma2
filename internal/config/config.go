package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultConfigFilePath = "/etc/enigma/config.yml"

// ServerConfig ...
type ServerConfig struct {
	Server struct {
		Bind          string `yaml:"bind"`
		ExternalURL   string `yaml:"external_url"`
		TemplatesPath string `yaml:"templates_path"`
	} `yaml:"server"`
	Secrets struct {
		Storage struct {
			Type  string `yaml:"type"`
			Await struct {
				Retries  int           `yaml:"retries"`
				Interval time.Duration `yaml:"interval"`
			} `yaml:"await"`
		} `yaml:"storage"`
		Hasher struct {
			Kind   string `yaml:"kind"`
			AES256 struct {
				Key string `yaml:"key"`
			} `yaml:"aes256"`
		} `yaml:"hasher"`
		Token struct {
			Length      int `yaml:"length"`
			SaveRetries int `yaml:"save_retries"`
		} `yaml:"token"`
	} `yaml:"secrets"`
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
}

// NewServerConfigFromFile ...
func NewServerConfigFromFile() (*ServerConfig, error) {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}

	data, err := os.ReadFile(configFilePath) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	var cfg ServerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return &cfg, nil
}
