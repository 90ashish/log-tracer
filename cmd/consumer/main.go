package main

import (
	"context"
	"fmt"
	"log-tracer/internal/config"
	"log-tracer/internal/consumer"
	"log-tracer/internal/pkg/logger"
)

func main() {
	fmt.Println("Consumer service started")

	// Load configuration
	cfg, err := config.LoadConsumerConfig("configs/consumer/config.yaml")
	if err != nil {
		logger.Error("Failed to load config", err)
		return
	}
	fmt.Printf("Loaded Kafka Config: %+v\n", cfg.KafkaConsumer)

	// Initialize Kafka consumer
	kafkaConsumer, err := consumer.NewKafkaConsumer(cfg)
	if err != nil {
		logger.Error("Failed to create Kafka consumer", err)
		return
	}
	defer kafkaConsumer.Close()

	// Start consuming messages
	err = kafkaConsumer.Consume(context.Background())
	if err != nil {
		logger.Error("Failed to consume messages", err)
	}
}
