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
	Host     string `yaml:"host"`
	GrpcPort int64  `yaml:"grpc_port"`
	HttpPort int64  `yaml:"http_port"`
}

type DB struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
}

type Kafka struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	OrderTopic string `yaml:"order_topic"`
	Brokers    string `yaml:"brokers"`
}

type Config struct {
	Service Service `yaml:"service"`
	DB      DB      `yaml:"db_master"`
	Kafka   Kafka   `yaml:"kafka"`
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
