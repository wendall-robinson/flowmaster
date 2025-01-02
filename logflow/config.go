package logflow

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	info    = "INFO"
	warning = "WARNING"
	debug   = "DEBUG"
	prod    = "PROD"
	dev     = "DEV"
	stage   = "STAGE"
	err     = "ERROR"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp,omitempty"`   // When the event occurred
	Level       string                 `json:"level,omitempty"`       // Log level (INFO, ERROR, etc.)
	Message     string                 `json:"message,omitempty"`     // Log message
	Context     map[string]interface{} `json:"context,omitempty"`     // Additional metadata
	Source      string                 `json:"source,omitempty"`      // Source module/system
	Environment string                 `json:"environment,omitempty"` // Environment (dev, prod, etc.)
	TraceID     string                 `json:"trace_id,omitempty"`    // Trace ID for distributed tracing
	SpanID      string                 `json:"span_id,omitempty"`     // Span ID for distributed tracing
	ErrorCode   string                 `json:"error_code,omitempty"`  // Error code (if applicable)
	Hostname    string                 `json:"hostname,omitempty"`    // Hostname or server name
	LogID       string                 `json:"log_id,omitempty"`      // Unique log identifier
}

// SetDefaults sets default values for the log entry
func (e *LogEntry) SetDefaults() {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	if e.Level == "" {
		e.Level = info
	}
}

func (e LogEntry) Validate() error {
	// List of required fields
	if e.Timestamp.IsZero() {
		return fmt.Errorf("missing required field: timestamp")
	}

	if e.Level == "" {
		return fmt.Errorf("missing required field: level")
	}

	if e.Message == "" {
		return fmt.Errorf("missing required field: message")
	}

	return nil
}

func (e LogEntry) String() string {
	// Convert the log entry to JSON or a readable format
	b, _ := json.Marshal(e)

	return string(b)
}
