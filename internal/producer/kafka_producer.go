package producer

import (
	"context"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/logger"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	Producer sarama.SyncProducer
	Config   *config.Config
}

func NewKafkaProducer(config *config.Config) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.RequiredAcks = sarama.RequiredAcks(config.Kafka.Acks)
	saramaConfig.Producer.Retry.Max = config.Kafka.Retries
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Compression = sarama.CompressionGZIP
	saramaConfig.Producer.Partitioner = NewCustomPartitioner

	producer, err := sarama.NewSyncProducer(config.Kafka.Brokers, saramaConfig)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{
		Producer: producer,
		Config:   config,
	}, nil
}

func (p *KafkaProducer) SendMessage(ctx context.Context, key, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: p.Config.Kafka.Topic,
		Key:   sarama.ByteEncoder(key),
		Value: sarama.ByteEncoder(message),
	}

	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		logger.Error("Failed to send message", err)
		return err
	}

	logger.Info("Message sent", "partition", partition, "offset", offset)
	return nil
}

func (p *KafkaProducer) Close() error {
	return p.Producer.Close()
}
