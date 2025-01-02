package logflow

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Logger is a logging interface with configurable defaults
type Logger struct {
	core               *log.Logger       // Core logger
	timestampFmt       string            // Timestamp format (e.g., RFC3339, Unix)
	includeHostname    bool              // Include hostname
	includeEnvironment bool              // Include environment
	includeK8s         bool              // Include Kubernetes metadata
	includeYANG        bool              // Include YANG-derived config
	includeSystem      bool              // Include system metadata
	hostname           string            // Hostname
	environment        string            // Environment
	k8sMetadata        map[string]string // Kubernetes metadata
	yangConfig         map[string]string // YANG-derived config
}

// logRequested will apply user-requested fields to the log entry set at the logger initialization
func (l *Logger) logRequested(entry *LogEntry) *LogEntry {
	// Format the timestamp based on the logger's configuration
	if l.timestampFmt != "" {
		entry.Timestamp, _ = time.Parse(time.RFC3339Nano, time.Now().Format(l.timestampFmt))
	} else {
		entry.Timestamp = time.Now()
	}

	// Include default hostname if configured
	if l.includeHostname {
		if hostname, err := os.Hostname(); err == nil {
			entry.Hostname = hostname
		} else {
			handleLoggingError(fmt.Errorf("failed to get hostname: %w", err))
		}
	}

	// Include environment if configured
	if l.includeEnvironment {
		entry.Environment = os.Getenv("ENVIRONMENT")

		if entry.Environment == "" {
			entry.Environment = "ENVIROMENT_NOT_SET"
		}
	}

	// Include Kubernetes metadata if configured
	if l.k8sMetadata != nil {
		for key, value := range l.k8sMetadata {
			if entry.Context == nil {
				entry.Context = make(map[string]interface{})
			}

			entry.Context[key] = value
		}
	}

	// Include YANG configuration metadata if available
	if l.yangConfig != nil {
		for key, value := range l.yangConfig {
			if entry.Context == nil {
				entry.Context = make(map[string]interface{})
			}
			entry.Context[key] = value
		}
	}

	// Include system metadata if configured
	if l.includeSystem {
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
		return "ENVIROMENT_NOT_SET"
	}

	return environment
}

// Log processes the LogEntry
func (e *LogEntry) Log() {
	if err := processLog(e); err != nil {
		handleLoggingError(err)
	}
}
