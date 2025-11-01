// Package logger provides logging abstraction.
package logger

import (
	"log"
	"os"
)

// Logger defines logging interface.
type Logger interface {
	Info(msg string)
	Infof(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
	Fatal(msg string)
	Fatalf(format string, args ...interface{})
}

// StdLogger implements Logger using stdlib log.
type StdLogger struct {
	info  *log.Logger
	error *log.Logger
}

// NewStdLogger creates standard logger.
func NewStdLogger() Logger {
	return &StdLogger{
		info:  log.New(os.Stdout, "INFO: ", log.LstdFlags),
		error: log.New(os.Stderr, "ERROR: ", log.LstdFlags),
	}
}

// Info logs info message.
func (l *StdLogger) Info(msg string) {
	l.info.Println(msg)
}

// Infof logs formatted info message.
func (l *StdLogger) Infof(format string, args ...interface{}) {
	l.info.Printf(format, args...)
}

// Error logs error message.
func (l *StdLogger) Error(msg string) {
	l.error.Println(msg)
}

// Errorf logs formatted error message.
func (l *StdLogger) Errorf(format string, args ...interface{}) {
	l.error.Printf(format, args...)
}

// Fatal logs fatal message and exits.
func (l *StdLogger) Fatal(msg string) {
	l.error.Fatal(msg)
}

// Fatalf logs formatted fatal message and exits.
func (l *StdLogger) Fatalf(format string, args ...interface{}) {
	l.error.Fatalf(format, args...)
}
