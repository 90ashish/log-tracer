package producer

import (
	"context"
	"log"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/logger"

	"github.com/IBM/sarama"
)

// KafkaProducer is a wrapper around Sarama's SyncProducer to manage the Kafka producer connection
type KafkaProducer struct {
	Producer sarama.SyncProducer
	Config   *config.ProducerConfig
}

// NewKafkaProducer initializes a new KafkaProducer with the provided configuration
func NewKafkaProducer(config *config.ProducerConfig) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	// Set Kafka producer configurations
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.KafkaProducerConfig.Acks)
	saramaConfig.Producer.Retry.Max = config.KafkaProducerConfig.Retries
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Compression = sarama.CompressionGZIP
	saramaConfig.Producer.Partitioner = NewCustomPartitioner

	// Create a new SyncProducer instance
	log.Println("Connecting to Kafka brokers:", config.KafkaProducerConfig.Brokers)
	producer, err := sarama.NewSyncProducer(config.KafkaProducerConfig.Brokers, saramaConfig)
	if err != nil {
		logger.Error("Failed to create Kafka producer", err)
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Config:   config,
	}, nil
}

// SendMessage sends a message to the specified Kafka topic
func (p *KafkaProducer) SendMessage(ctx context.Context, key, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.Config.KafkaProducerConfig.Topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	// Send the message and log the result
	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		logger.Error("Failed to send message to Kafka broker:", err)
		return err
	}

	logger.Info("Message sent", "partition", partition, "offset", offset)
	return nil
}

// Close closes the Kafka producer connection
func (p *KafkaProducer) Close() error {
	return p.Producer.Close()
}
