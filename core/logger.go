package core

import (
	"log"
	"os"
)

var (
	// AccessLog holds a global access LogObject
	AccessLog loggerInterface

	// SystemLog holds a global system LogObject
	SystemLog loggerInterface
)

func setupLogger(output string) loggerInterface {
	switch output {
	case "stdout":
		return &stdLogger{}
	case "null":
		return &nullLogger{}
	default:
		fd, err := os.OpenFile(output, os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			log.Fatalf(logOutputErrStr, output, err.Error())
		}
		return &logger{log.New(fd, "", log.LstdFlags)}
	}
}

// LoggerInterface specifies an interface that can log different message levels
type loggerInterface interface {
	Info(string, ...interface{})
	Error(string, ...interface{})
	Fatal(string, ...interface{})
}

// StdLogger implements LoggerInterface to log to output using regular log
type stdLogger struct{}

// Info logs to log.Logger with info level prefix
func (l *stdLogger) Info(fmt string, args ...interface{}) {
	log.Printf(":: I :: "+fmt, args...)
}

// Error logs to log.Logger with error level prefix
func (l *stdLogger) Error(fmt string, args ...interface{}) {
	log.Printf(":: E :: "+fmt, args...)
}

// Fatal logs to standard log with fatal prefix and terminates program
func (l *stdLogger) Fatal(fmt string, args ...interface{}) {
	log.Fatalf(":: F :: "+fmt, args...)
}

// logger implements LoggerInterface to log to output using underlying log.Logger
type logger struct {
	logger *log.Logger
}

// Info logs to log.Logger with info level prefix
func (l *logger) Info(fmt string, args ...interface{}) {
	l.logger.Printf("I :: "+fmt, args...)
}

// Error logs to log.Logger with error level prefix
func (l *logger) Error(fmt string, args ...interface{}) {
	l.logger.Printf("E :: "+fmt, args...)
}

// Fatal logs to log.Logger with fatal prefix and terminates program
func (l *logger) Fatal(fmt string, args ...interface{}) {
	l.logger.Fatalf("F :: "+fmt, args...)
}

// nullLogger implements LoggerInterface to do absolutely fuck-all
type nullLogger struct{}

// Info does nothing
func (l *nullLogger) Info(fmt string, args ...interface{}) {
	// do nothing
}

// Error does nothing
func (l *nullLogger) Error(fmt string, args ...interface{}) {
	// do nothing
}

// Fatal simply terminates the program
func (l *nullLogger) Fatal(fmt string, args ...interface{}) {
	os.Exit(1)
}
