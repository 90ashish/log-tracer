package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// LogSourceConfig defines the structure for each log source
type LogSourceConfig struct {
	Name          string `yaml:"name"`
	Environment   string `yaml:"environment"`
	SeverityLevel string `yaml:"severity_level"`
}

// KafkaProducerConfig holds the Kafka configuration parameters for the producer
type KafkaProducerConfig struct {
	Brokers         []string          `yaml:"brokers"`
	Topic           string            `yaml:"topic"`
	Acks            int               `yaml:"acks"`
	Retries         int               `yaml:"retries"`
	BatchSize       int               `yaml:"batch_size"`
	LingerMS        int               `yaml:"linger_ms"`
	CompressionType string            `yaml:"compression_type"`
	Sources         []LogSourceConfig `yaml:"sources"`
}

// ProducerConfig is the main configuration structure holding producer configurations
type ProducerConfig struct {
	KafkaProducerConfig KafkaProducerConfig `yaml:"kafka_producer"`
}

// LoadProducerConfig loads the producer configuration from the given file path
func LoadProducerConfig(configPath string) (*ProducerConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config ProducerConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Loaded producer config: %+v\n", config)

	return &config, nil
}
