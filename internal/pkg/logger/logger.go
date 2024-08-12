package logger

import (
	"fmt"
	"log"
	"os"
)

func Info(args ...interface{}) {
	log.Output(2, "INFO: "+fmt.Sprintln(args...))
}

func Error(args ...interface{}) {
	log.Output(2, "ERROR: "+fmt.Sprintln(args...))
}

func Debug(args ...interface{}) {
	log.Output(2, "DEBUG: "+fmt.Sprintln(args...))
}

func Warn(args ...interface{}) {
	log.Output(2, "WARN: "+fmt.Sprintln(args...))
}

func Fatal(args ...interface{}) {
	log.Output(2, "FATAL: "+fmt.Sprintln(args...))
	os.Exit(1)
}
