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
	InstanceID string `yaml:"instance_id"`
}

type Kafka struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	OrderTopic      string `yaml:"order_topic"`
	Brokers         string `yaml:"brokers"`
	ConsumerGroupID string `yaml:"consumer_group_id"`
}

type Config struct {
	Service Service `yaml:"service"`
	Kafka   Kafka   `yaml:"kafka"`
}

func LoadConfig() (*Config, error) {
	path := filepath.Clean(os.Getenv(configFileEnv))

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open config file: %w", err)
	}

	defer f.Close()

	config, err := readConfig(f)
	if err != nil {
		return nil, err
	}

	// Override instance ID from environment if provided
	instanceID := os.Getenv("INSTANCE_ID")
	if instanceID != "" {
		config.Service.InstanceID = instanceID
	}

	return config, nil
}

func readConfig(file io.Reader) (*Config, error) {
	config := &Config{}
	if err := yaml.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("can't decode config: %w", err)
	}

	return config, nil
}
