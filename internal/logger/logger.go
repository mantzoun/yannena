package logger

import (
	"log"
	"os"
)

type Logger struct {
	level  int
	prefix string
	log    *log.Logger
}

func NewLogger(prefix string) *Logger {
	var logger Logger

	logger.log = log.New(os.Stdout, "", log.LstdFlags)
	logger.prefix = prefix
	logger.level = 0

	return &logger
}

func (logger *Logger) Debug(message string) {
	logger.log.Printf("%-8s %-8s %s\n", "DEBUG", logger.prefix, message)
}
