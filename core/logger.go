package core

import (
	"log"
	"os"
)

var (
	// AccessLog holds a global access LogObject
	AccessLog LoggerInterface

	// SystemLog holds a global system LogObject
	SystemLog LoggerInterface
)

func setupLogger(output string) LoggerInterface {
	switch output {
	case "stdout":
		return &StdLogger{}
	case "null":
		return &NullLogger{}
	default:
		fd, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatalf("Error opening log output %s: %s", output, err.Error())
		}
		return &Logger{log.New(fd, "", log.LstdFlags)}
	}
}

// LoggerInterface specifies an interface that can log different message levels
type LoggerInterface interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

// StdLogger implements LoggerInterface to log to output using regular log
type StdLogger struct{}

// Info logs to log.Logger with info level prefix
func (l *StdLogger) Info(fmt string, args ...interface{}) {
	log.Printf(":: I :: "+fmt, args...)
}

// Error logs to log.Logger with error level prefix
func (l *StdLogger) Error(fmt string, args ...interface{}) {
	log.Printf(":: E :: "+fmt, args...)
}

// Fatal logs to standard log with fatal prefix and terminates program
func (l *StdLogger) Fatal(fmt string, args ...interface{}) {
	log.Fatalf(":: F :: "+fmt, args...)
}

// Logger implements LoggerInterface to log to output using underlying log.Logger
type Logger struct {
	logger *log.Logger
}

// Info logs to log.Logger with info level prefix
func (l *Logger) Info(fmt string, args ...interface{}) {
	l.logger.Printf("I :: "+fmt, args...)
}

// Error logs to log.Logger with error level prefix
func (l *Logger) Error(fmt string, args ...interface{}) {
	l.logger.Printf("E :: "+fmt, args...)
}

// Fatal logs to log.Logger with fatal prefix and terminates program
func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.logger.Fatalf("F :: "+fmt, args...)
}

// NullLogger implements LoggerInterface to do absolutely fuck-all
type NullLogger struct {
	LoggerInterface
}

// Info does nothing
func (l *NullLogger) Info(fmt string, args ...interface{}) {
	// do nothing
}

// Error does nothing
func (l *NullLogger) Error(fmt string, args ...interface{}) {
	// do nothing
}

// Fatal simply terminates the program
func (l *NullLogger) Fatal(fmt string, args ...interface{}) {
	os.Exit(1)
}
