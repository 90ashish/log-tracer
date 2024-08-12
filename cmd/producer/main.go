package main

import (
	"context"
	"fmt"
	"log"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/config_watcher"
	"log-tracer/internal/pkg/logger"
	"log-tracer/internal/producer"
	"sync"
	"time"
)

func main() {
	fmt.Println("Producer service started")

	// Load initial configuration
	configPath := "configs/producer/config.yaml"
	cfg, err := config.LoadProducerConfig(configPath)
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
	var mu sync.Mutex // Mutex to protect configuration reloads

	// Define reload function to update log sources dynamically
	reloadFunc := func(newConfig *config.ProducerConfig) {
		mu.Lock()
		defer mu.Unlock()
		fmt.Println("Reloading configuration...")
		logHandler.UpdateConfig(newConfig)
		fmt.Printf("New Kafka Config: %+v\n", newConfig.KafkaProducerConfig)
	}

	// Start watching the configuration file for changes
	configWatcher := config_watcher.NewConfigWatcher(configPath, reloadFunc)
	go func() {
		if err := configWatcher.WatchConfig(); err != nil {
			log.Fatalf("Failed to watch config file: %v", err)
		}
	}()

	// Simulate log messages being sent
	for {
		err = logHandler.HandleLogMessage(context.Background(), "test_key", "test_message")
		if err != nil {
			logger.Error("Failed to handle log message", err)
			return
		}
		fmt.Println("Log message sent successfully")
		// Sleep or wait to simulate ongoing logging
		time.Sleep(5 * time.Second)
	}
}
