package config

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const defaultConfigFilePath = "/etc/enigma/config.yml"

// ServerConfig ...
type ServerConfig struct {
	Server struct {
		Bind        string `yaml:"bind"`
		ExternalURL string `yaml:"external_url"`
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
		Token struct {
			Length      int `yaml:"lenght"`
			SaveRetries int `yaml:"save_retries"`
		} `yaml:"token"`
	} `yaml:"secrets"`
	Redis struct {
		Address  string `yaml:"address"`
		Password string `yaml:"password"`
		Database int    `yaml:"database"`
	} `yaml:"redis"`
}

func NewServerConfigFromFile() (*ServerConfig, error) {
	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		configFilePath = defaultConfigFilePath
	}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	var cfg ServerConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
