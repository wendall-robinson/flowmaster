package logflow

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
)

// Logger is a logging interface with configurable defaults
type Logger struct {
	core         *log.Logger       // Core logger
	config       LoggerConfig      // Logger configuration
	timestampFmt string            // Timestamp format (e.g., RFC3339, Unix)
	hostname     string            // Hostname
	environment  string            // Environment
	k8sMetadata  map[string]string // Kubernetes metadata
	yangConfig   map[string]string // YANG-derived config
}

// New initializes and returns a new Logger with the given configuration options.
func New(options ...Option) (*Logger, error) {
	config := DefaultLoggerConfig()

	// Apply the provided options
	for _, option := range options {
		option(&config)
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid logger configuration: %w", err)
	}

	// Create the core logger
	l := &Logger{
		config: config,
	}

	// Set the timestamp format
	l.timestampFmt = config.TimestampFmt

	// Find the requested metadata
	l.findMetaData()

	return l, nil
}

// findMetaData will find the requested metadata for the logger
func (l *Logger) findMetaData() {
	if l.config.IncludeHostname {
		l.hostname = l.getHostname()
	}

	if l.config.IncludeEnvironment {
		l.environment = l.getEnvironment()
	}

	if l.config.IncludeK8s {
		l.k8sMetadata = make(map[string]string)
		// TODO: Add Kubernetes metadata retrieval
	}

	if l.config.IncludeYANG {
		l.yangConfig = make(map[string]string)
		// TODO: Add YANG metadata retrieval
	}
}

// logRequired will apply required fields to the log entry set at the logger initialization
func (l *Logger) logRequired(entry *LogEntry) *LogEntry {
	if l.config.IncludeHostname {
		entry.Hostname = l.hostname
	}

	if l.config.IncludeEnvironment {
		entry.Environment = l.environment
	}

	if l.config.IncludeK8s {
		entry.Context["k8s"] = l.k8sMetadata
	}

	if l.config.IncludeYANG {
		entry.Context["yang"] = l.yangConfig
	}

	if l.config.IncludeSystem {
		l.findSystemAttributes(entry)
	}

	return entry
}

// getHostname returns a hostname value from the logger configuration or OS environment variables if set
func (l *Logger) getHostname() string {
	if l.hostname != "" {
		return l.hostname
	}

	hostname, err := os.Hostname()
	if err != nil {
		return "HOSTNAME_NOT_SET"
	}

	return hostname
}

// getEnvironment returns an environment value from the logger configuration or OS environment variables if set
func (l *Logger) getEnvironment() string {
	if l.environment != "" {
		return l.environment
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		return "ENVIRONMENT_NOT_SET"
	}

	return environment
}

// Helper function to generate a unique LogID
func generateLogID() string {
	// Generate a random LogID (replace with a more robust solution if needed)
	return fmt.Sprintf("%d", rand.Int63())
}

// processLog handles processing the log entry, including formatting and sending it to configured sinks.
func (l *Logger) processLog(entry *LogEntry) {
	// Format the log entry as JSON
	data, err := json.Marshal(entry)
	if err != nil {
		l.core.Printf("failed to marshal log entry: %v", err)
	}

	// Send to stdout if enabled
	if l.config.OutputToStdout {
		fmt.Println(string(data))
	}

	// Send to log collector if configured
	if l.config.LogCollector != nil {
		if err := l.config.LogCollector.Send(data); err != nil {
			l.core.Printf("failed to send log to collector: %v", err)
		}
	}
}
