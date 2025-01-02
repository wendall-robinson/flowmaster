package logflow

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LoggerOption defines a function that configures the Logger
type LoggerOption func(*Logger)

// NewLogger creates a new Logger with optional configurations
func NewLogger(options ...LoggerOption) *Logger {
	logger := &Logger{
		core:               log.Default(),
		timestampFmt:       time.RFC3339, // Default timestamp format
		includeHostname:    false,
		includeEnvironment: false,
		includeK8s:         false,
		includeYANG:        false,
	}

	// Apply all options
	for _, option := range options {
		option(logger)
	}

	return logger
}

// WithTimestampFormat configures the Logger to use a custom timestamp format.
//
// Example Usage:
//
// logger := logflow.NewLogger(
//
//	WithTimestampFormat("2006-01-02T15:04:05Z07:00"), // ISO 8601 format
//
// )
//
// logEntry := logger.InfoSys("System initialized")
// fmt.Println(logEntry)
//
// Output:
//
//	{
//	    "timestamp": "2025-01-02T15:04:05Z",
//	    "level": "INFO",
//	    "message": "System initialized",
//	    "log_id": "873460012673882214"
//	}
func WithTimeStampFormat(format string) LoggerOption {
	return func(l *Logger) {
		// Use a sample time for validation
		sampleTime := time.Now()
		formatted := sampleTime.Format(format)

		// If the format produces an empty result, assume it's invalid
		if formatted == "" {
			panic(fmt.Sprintf("Invalid timestamp format: %s. Please use a valid Go time format string.", format))
		}

		l.timestampFmt = format
	}
}

// WithEnvHostname configures the Logger to include the hostname in log entries
func WithEnvHostname() LoggerOption {
	return func(l *Logger) {
		l.includeHostname = true

		hostname, err := os.Hostname()
		if err != nil {
			panic(fmt.Sprintf("Failed to get hostname: %v", err))
		}

		l.hostname = hostname
	}
}

// WithEnvironment configures the Logger to include the environment in log entries
func WithEnvironment() LoggerOption {
	return func(l *Logger) {
		l.includeEnvironment = true

		env := os.Getenv("ENVIRONMENT")
		if env == "" {
			panic("Environment not set")
		}

		l.environment = env
	}
}

// WithK8sMetadata adds Kubernetes metadata to the Logger
func WithK8sMetadata(metadata map[string]string) LoggerOption {
	return func(l *Logger) {
		l.includeK8s = true
		l.k8sMetadata = make(map[string]string)

		for key, value := range metadata {
			l.k8sMetadata[key] = value
		}
	}
}

// WithYangConfig adds YANG-derived configuration to the Logger
func WithYangConfig(config map[string]string) LoggerOption {
	return func(l *Logger) {
		l.includeYANG = true
		l.yangConfig = make(map[string]string)

		for key, value := range config {
			l.yangConfig[key] = value
		}
	}
}

func WithSystemInfo() LoggerOption {
	return func(l *Logger) {
		l.includeSystem = true
	}
}
