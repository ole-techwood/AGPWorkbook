package internal

import (
	"fmt"
	"os"
)

// Logger writes user-facing messages to the console.
type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Info(format string, args ...any) {
	fmt.Printf(format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}
