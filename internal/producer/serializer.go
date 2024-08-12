package producer

import (
	"encoding/json"
	"log-tracer/internal/pkg/logger"
)

// SerializeToJson serializes the given data into JSON format
func SerializeToJson(data interface{}) ([]byte, error) {
	serializedData, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to serialize data", err)
		return nil, err
	}
	return serializedData, nil
}
