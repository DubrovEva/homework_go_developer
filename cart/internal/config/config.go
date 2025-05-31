package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	configFileEnv = "CONFIG_FILE"
)

type Service struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type Products struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Token      string `yaml:"token"`
	Schema     string `yaml:"schema"`
	RetryCount uint   `yaml:"retry_count"`
	RPS        int    `yaml:"rps"`
}

type Loms struct {
	Host string `yaml:"host"`
	Port int64  `yaml:"port"`
}

type Config struct {
	Service  Service  `yaml:"service"`
	Products Products `yaml:"product_service"`
	Loms     Loms     `yaml:"loms_service"`
}

func LoadConfig() (*Config, error) {
	path := filepath.Clean(os.Getenv(configFileEnv))

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %w", err)
	}

	defer f.Close()

	return readConfig(f)
}

func readConfig(file io.Reader) (*Config, error) {
	config := &Config{}
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("can't decode config: %w", err)
	}

	return config, nil
}
