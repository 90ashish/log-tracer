package consumer

import (
	"log-tracer/internal/pkg/logger"

	"github.com/IBM/sarama"
)

// MessageHandler implements the ConsumerGroupHandler interface to process messages
type MessageHandler struct{}

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
		logger.Info("Message received", "partition", message.Partition, "offset", message.Offset, "key", string(message.Key), "value", string(message.Value))
		session.MarkMessage(message, "")
	}
	return nil
}
