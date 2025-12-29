package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

// LogLevel defines the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

type CustomLogger struct {
	logger   *log.Logger
	logLevel LogLevel
}

var logger *CustomLogger

// InitLogger initializes the custom logger with default INFO level
func InitLogger() {
	logger = &CustomLogger{
		logger:   log.New(os.Stdout, "", 0),
		logLevel: INFO,
	}
	logger.Info("Logger initialized")
}

// InitLoggerWithLevel initializes the custom logger with specified log level
func InitLoggerWithLevel(level LogLevel) {
	logger = &CustomLogger{
		logger:   log.New(os.Stdout, "", 0),
		logLevel: level,
	}
	logger.Info("Logger initialized")
}
func LoggingMiddleware(route string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := log.New(os.Stdout, "", log.LstdFlags)
		logger.Printf("=============================================")
		logger.Printf("ENTRY-------------------> %s %s", r.Method, route)
		logger.Printf("=============================================")
		defer func() {
			logger.Printf("=============================================")
			logger.Printf("EXIT<------------------- %s %s ", r.Method, route)
			logger.Printf("=============================================")
		}()
		handler(w, r)
	}
}

// InitLoggerFromConfig initializes the logger and sets log level from environment
func InitLoggerFromConfig() {
	// Initialize logger with default level (INFO)
	InitLogger()
	logger := GetLogger()

	// Set log level based on environment variable
	// Options: DEBUG, INFO, WARN, ERROR
	// Default is INFO if not set or invalid
	logLevelEnv := os.Getenv("LOG_LEVEL")
	switch logLevelEnv {
	case "DEBUG":
		SetLogLevel(DEBUG)
		logger.Info("Log level set to DEBUG")
	case "WARN":
		SetLogLevel(WARN)
		logger.Info("Log level set to WARN")
	case "ERROR":
		SetLogLevel(ERROR)
		logger.Info("Log level set to ERROR")
	default:
		SetLogLevel(INFO)
		if logLevelEnv != "" {
			logger.Warnf("Invalid LOG_LEVEL '%s', using INFO", logLevelEnv)
		}
	}
}

// SetLogLevel sets the logging level
func SetLogLevel(level LogLevel) {
	if logger != nil {
		logger.logLevel = level
	}
}

// GetLogger returns the singleton logger instance
func GetLogger() *CustomLogger {
	if logger == nil {
		InitLogger()
	}
	return logger
}

// getCallerInfo retrieves the file name and line number of the caller
func getCallerInfo() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown:0"
	}
	// Get only the file name, not the full path
	parts := strings.Split(file, "/")
	fileName := parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", fileName, line)
}

// formatLog formats the log message with timestamp, level, caller info, and message
func (l *CustomLogger) formatLog(level LogLevel, color, levelName, message string) {
	// Check if this log level should be printed
	if level < l.logLevel {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	caller := getCallerInfo()
	l.logger.Printf("%s[%s] [%s] [%s] %s%s\n", color, timestamp, levelName, caller, message, ColorReset)
}

// Info logs an informational message in green
func (l *CustomLogger) Info(message string) {
	l.formatLog(INFO, ColorGreen, "INFO", message)
}

// Debug logs a debug message in yellow
func (l *CustomLogger) Debug(message string) {
	l.formatLog(DEBUG, ColorYellow, "DEBUG", message)
}

// Error logs an error message in red
func (l *CustomLogger) Error(message string) {
	l.formatLog(ERROR, ColorRed, "ERROR", message)
}

// Warn logs a warning message in purple
func (l *CustomLogger) Warn(message string) {
	l.formatLog(WARN, ColorPurple, "WARN", message)
}

// Infof logs an informational message with formatting
func (l *CustomLogger) Infof(format string, args ...interface{}) {
	l.formatLog(INFO, ColorGreen, "INFO", fmt.Sprintf(format, args...))
}

// Debugf logs a debug message with formatting
func (l *CustomLogger) Debugf(format string, args ...interface{}) {
	l.formatLog(DEBUG, ColorYellow, "DEBUG", fmt.Sprintf(format, args...))
}

// Errorf logs an error message with formatting
func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	l.formatLog(ERROR, ColorRed, "ERROR", fmt.Sprintf(format, args...))
}

// Warnf logs a warning message with formatting
func (l *CustomLogger) Warnf(format string, args ...interface{}) {
	l.formatLog(WARN, ColorPurple, "WARN", fmt.Sprintf(format, args...))
}
