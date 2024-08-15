package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

// FilteringConfig defines the structure for log filtering configuration
type FilteringConfig struct {
	MinSeverity string `yaml:"min_severity"`
}

// EnrichmentConfig defines the structure for log enrichment configuration
type EnrichmentConfig struct {
	AddTimestamp bool `yaml:"add_timestamp"`
	AddHostInfo  bool `yaml:"add_host_info"`
}

// RetryConfig defines the structure for retry configuration
type RetryConfig struct {
	MaxRetries       int `yaml:"max_retries"`
	InitialBackoffMS int `yaml:"initial_backoff_ms"`
	MaxBackoffMS     int `yaml:"max_backoff_ms"`
}

// KafkaConsumerConfig holds the Kafka configuration parameters for the consumer
type KafkaConsumerConfig struct {
	Brokers         []string         `yaml:"brokers"`
	GroupID         string           `yaml:"group_id"`
	Topic           string           `yaml:"topic"`
	AutoOffsetReset string           `yaml:"auto_offset_reset"`
	MaxWaitTime     int              `yaml:"max_wait_time"`
	MinBytes        int              `yaml:"min_bytes"`
	MaxBytes        int              `yaml:"max_bytes"`
	Filtering       FilteringConfig  `yaml:"filtering"`
	Enrichment      EnrichmentConfig `yaml:"enrichment"`
	Retry           RetryConfig      `yaml:"retry"`
}

// ConsumerConfig is the main configuration structure holding consumer configurations
type ConsumerConfig struct {
	KafkaConsumer KafkaConsumerConfig `yaml:"kafka_consumer"`
}

// LoadConsumerConfig loads the consumer configuration from the given file path
func LoadConsumerConfig(configPath string) (*ConsumerConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var config ConsumerConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Loaded consumer config: %+v\n", config)

	return &config, nil
}
