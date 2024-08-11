package main

import (
	"context"
	"fmt"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/logger"
	"log-tracer/internal/producer"
)

func main() {
	fmt.Println("Producer service started")

	// Load configuration
	cfg, err := config.LoadProducerConfig("configs/producer/config.yaml")
	if err != nil {
		logger.Error("Failed to load config", err)
		return
	}
	fmt.Printf("Loaded Kafka Config: %+v\n", cfg.KafkaProducerConfig)

	// Initialize Kafka producer
	kafkaProducer, err := producer.NewKafkaProducer(cfg)
	if err != nil {
		logger.Error("Failed to create Kafka producer", err)
		fmt.Printf("Failed to initialize Kafka producer: %v\n", err)
		return
	}
	defer kafkaProducer.Close()

	// Initialize LogHandler
	logHandler := producer.NewLogHandler(kafkaProducer)

	// Test sending a log message
	err = logHandler.HandleLogMessage(context.Background(), "test_key", "test_message")
	if err != nil {
		logger.Error("Failed to handle log message", err)
		return
	}

	fmt.Println("Log message sent successfully")
}
