package logflow

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

const (
	// INFO level log
	INFO = "INFO"
	// WARN level log
	WARN = "WARN"
	// DEBUG level log
	DEBUG = "DEBUG"
	// ERROR level log
	ERROR = "ERROR"
)

var (
	errNoLogLevel      = errors.New("log level must be specified")
	errTimeStampFormat = errors.New("timestamp format must be specified")
)

// LoggerConfig holds configuration options for the logger
type LoggerConfig struct {
	OutputToStdout     bool              // Whether to print logs to stdout
	LogCollector       LogCollector      // Interface for sending logs to a centralized collector
	LogLevel           string            // Minimum log level to output (e.g., INFO, DEBUG)
	TimestampFmt       string            // Format for timestamps (e.g., RFC3339, Unix)
	IncludeHostname    bool              // Include hostname in logs
	IncludeEnvironment bool              // Include environment (e.g., dev, prod)
	IncludeSystem      bool              // Include system metadata in all logs
	IncludeK8s         bool              // Include Kubernetes metadata in logs
	IncludeYANG        bool              // Include YANG-based configuration metadata
	YANGConfig         map[string]string // YANG-based configuration metadata
	K8sConfig          map[string]string // Kubernetes metadata
}

// LogCollector defines an interface for sending logs to a centralized logging system
type LogCollector interface {
	Send(data []byte) error
}

// DefaultLoggerConfig returns a LoggerConfig with default settings
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		OutputToStdout:     true,
		LogLevel:           INFO,
		TimestampFmt:       time.RFC3339,
		IncludeHostname:    false,
		IncludeEnvironment: false,
		IncludeSystem:      false,
		IncludeK8s:         false,
		IncludeYANG:        false,
	}
}

// SetDefaults populates default values for missing configuration options
func (c *LoggerConfig) SetDefaults() {
	if c.TimestampFmt == "" {
		c.TimestampFmt = time.RFC3339
	}
	if c.LogLevel == "" {
		c.LogLevel = INFO
	}
}

// Validate checks if the LoggerConfig is valid
func (c *LoggerConfig) Validate() error {
	if c.LogLevel == "" {
		return errNoLogLevel
	}

	if c.TimestampFmt == "" {
		return errTimeStampFormat
	}

	return nil
}

func (e LogEntry) String() string {
	// Convert the log entry to JSON or a readable format
	b, _ := json.Marshal(e)

	return string(b)
}
