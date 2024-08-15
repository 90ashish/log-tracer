package consumer

import (
	"encoding/json"
	"log"
	"log-tracer/internal/config"
	"log-tracer/internal/pkg/logger"
	"os"
	"time"

	"github.com/IBM/sarama"
)

// MessageHandler implements the ConsumerGroupHandler interface to process messages
type MessageHandler struct {
	config *config.ConsumerConfig
}

// NewMessageHandler creates a new MessageHandler with the provided configuration
func NewMessageHandler(cfg *config.ConsumerConfig) *MessageHandler {
	return &MessageHandler{
		config: cfg,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *MessageHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *MessageHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from Kafka
func (h *MessageHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logMessage := string(message.Value)

		// Debugging: Print the raw log message
		logger.Info("Raw log message received", "logMessage", logMessage)

		if !h.shouldProcessLog(logMessage) {
			continue // Skip the message if it doesn't meet the filtering criteria
		}

		logMessage = h.enrichLog(logMessage)

		logger.Info("Processed log message", "logMessage", logMessage)

		session.MarkMessage(message, "")
	}
	return nil
}

// shouldProcessLog checks if the log should be processed based on severity
func (h *MessageHandler) shouldProcessLog(logMessage string) bool {
	logSeverity := h.extractSeverity(logMessage)
	return h.isSeverityAllowed(logSeverity)
}

// extractSeverity extracts the severity level from the log message
func (h *MessageHandler) extractSeverity(logMessage string) string {
	// Assume logMessage is a JSON string
	var logData map[string]interface{}
	err := json.Unmarshal([]byte(logMessage), &logData)
	if err != nil {
		log.Println("Failed to unmarshal log message:", err)
		return "UNKNOWN"
	}

	if severity, ok := logData["severity_level"].(string); ok {
		return severity
	}
	return "UNKNOWN"
}

// isSeverityAllowed checks if the log severity meets the minimum severity level
func (h *MessageHandler) isSeverityAllowed(logSeverity string) bool {
	severityOrder := map[string]int{
		"DEBUG": 1,
		"INFO":  2,
		"WARN":  3,
		"ERROR": 4,
	}
	return severityOrder[logSeverity] >= severityOrder[h.config.KafkaConsumer.Filtering.MinSeverity]
}

// enrichLog enriches the log with additional metadata
func (h *MessageHandler) enrichLog(logMessage string) string {
	if h.config.KafkaConsumer.Enrichment.AddTimestamp {
		logMessage = h.addTimestamp(logMessage)
	}
	if h.config.KafkaConsumer.Enrichment.AddHostInfo {
		logMessage = h.addHostInfo(logMessage)
	}
	return logMessage
}

// addTimestamp adds a timestamp to the message value
func (h *MessageHandler) addTimestamp(value string) string {
	timestamp := time.Now().Format(time.RFC3339)
	return value + " " + timestamp
}

// addHostInfo adds host information to the message value
func (h *MessageHandler) addHostInfo(value string) string {
	host, err := os.Hostname()
	if err != nil {
		log.Println("Failed to get hostname:", err)
		host = "unknown-host"
	}
	return value + " " + host
}
