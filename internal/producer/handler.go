package producer

import (
	"context"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/logger"
	"sync"
	"time"
)

// LogHandler manages the creation and sending of log messages to Kafka
type LogHandler struct {
	Producer *KafkaProducer
	Config   *config.ProducerConfig
	mu       sync.RWMutex // Mutex to protect configuration reloads
}

// NewLogHandler initializes a new LogHandler with the provided KafkaProducer
func NewLogHandler(producer *KafkaProducer) *LogHandler {
	return &LogHandler{
		Producer: producer,
		Config:   producer.Config,
	}
}

// UpdateConfig updates the LogHandler's configuration dynamically
func (h *LogHandler) UpdateConfig(newConfig *config.ProducerConfig) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Config = newConfig
}

// HandleLogMessage creates and sends log messages to Kafka for each source
func (h *LogHandler) HandleLogMessage(ctx context.Context, key string, logMessage interface{}) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, source := range h.Config.KafkaProducerConfig.Sources {
		// Add contextual information to the log message
		contextualLogMessage := map[string]interface{}{
			"timestamp":      time.Now().Format(time.RFC3339),
			"service_name":   source.Name,
			"environment":    source.Environment,
			"severity_level": source.SeverityLevel,
			"message":        logMessage,
		}

		// Serialize the log message to JSON
		serializedMessage, err := SerializeToJson(contextualLogMessage)
		if err != nil {
			logger.Error("Failed to serialize log message", err)
			return err
		}

		// Send the serialized log message to Kafka
		err = h.Producer.SendMessage(ctx, []byte(key), serializedMessage)
		if err != nil {
			logger.Error("Failed to send log message", err)
			return err
		}

		logger.Info("Log message handled successfully", "service_name", source.Name)
	}
	return nil
}
