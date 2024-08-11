package logger

import (
	"log"
)

func Error(message string, err error) {
	log.Printf("[ERROR] %s: %v", message, err)
}

func Info(message string, keysAndValues ...interface{}) {
	log.Printf("[INFO] %s: %v", message, keysAndValues)
}
