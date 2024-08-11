package producer

import (
	"context"
	"log-tracer/internal/pkg/logger"
)

// LogHandler manages the creation and sending of log messages to Kafka
type LogHandler struct {
	Producer *KafkaProducer
}

// NewLogHandler initializes a new LogHandler with the provided KafkaProducer
func NewLogHandler(producer *KafkaProducer) *LogHandler {
	return &LogHandler{
		Producer: producer,
	}
}

// HandleLogMessage creates and sends a log message to Kafka
func (h *LogHandler) HandleLogMessage(ctx context.Context, key string, logMessage interface{}) error {
	// Serialize the log message to JSON
	serializedMessage, err := SerializeToJson(logMessage)
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

	logger.Info("Log message handled successfully")
	return nil
}
