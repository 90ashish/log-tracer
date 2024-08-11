package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// KafkaConfig holds the Kafka configuration parameters
type KafkaConfig struct {
	Brokers         []string `yaml:"brokers"`
	Topic           string   `yaml:"topic"`
	Acks            int      `yaml:"acks"` // Changed to int
	Retries         int      `yaml:"retries"`
	BatchSize       int      `yaml:"batch_size"`
	LingerMS        int      `yaml:"linger_ms"`
	CompressionType string   `yaml:"compression_type"`
}

// Config is the main configuration structure holding all configurations
type Config struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

// LoadConfig loads the configuration from the given file path
func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
