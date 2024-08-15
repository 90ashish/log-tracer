package consumer

import (
	"context"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

// KafkaConsumer wraps the Sarama ConsumerGroup to manage message consumption
type KafkaConsumer struct {
	ConsumerGroup sarama.ConsumerGroup
	Config        *config.ConsumerConfig
}

// NewKafkaConsumer initializes a new KafkaConsumer with the provided configuration
func NewKafkaConsumer(cfg *config.ConsumerConfig) (*KafkaConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	saramaConfig.Consumer.MaxWaitTime = time.Duration(cfg.KafkaConsumer.MaxWaitTime) * time.Millisecond
	saramaConfig.Consumer.Fetch.Min = int32(cfg.KafkaConsumer.MinBytes)
	saramaConfig.Consumer.Fetch.Max = int32(cfg.KafkaConsumer.MaxBytes)

	consumerGroup, err := sarama.NewConsumerGroup(cfg.KafkaConsumer.Brokers, cfg.KafkaConsumer.GroupID, saramaConfig)
	if err != nil {
		logger.Error("Failed to create Kafka consumer group", err)
		return nil, err
	}

	return &KafkaConsumer{
		ConsumerGroup: consumerGroup,
		Config:        cfg,
	}, nil
}

// Consume starts consuming messages from the Kafka topic
func (kc *KafkaConsumer) Consume(ctx context.Context) error {
	handler := NewMessageHandler(kc.Config)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for {
			if err := kc.ConsumerGroup.Consume(ctx, []string{kc.Config.KafkaConsumer.Topic}, handler); err != nil {
				logger.Error("Error consuming messages", err)
				break
			}
		}
	}()

	// Wait for termination signal
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-sigterm:
		cancel()
		return nil
	}
}

// Close closes the Kafka consumer connection
func (kc *KafkaConsumer) Close() error {
	return kc.ConsumerGroup.Close()
}
