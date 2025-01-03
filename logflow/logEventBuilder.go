package logflow

import "time"

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
	ErrorMsg    string                 `json:"error_msg,omitempty"`   // Error message (if applicable)
	Hostname    string                 `json:"hostname,omitempty"`    // Hostname or server name
	LogID       string                 `json:"log_id,omitempty"`      // Unique log identifier
}

// baseEventLog creates a new LogEntry with the defined log level, message and user provided args
//
// if logging was initialized with required attributes, baseEventLog will ensure they are added to the log entry
func (l *Logger) baseEntryLog(level, message string, args ...string) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		LogID:     generateLogID(),
		Context:   make(map[string]interface{}),
	}

	// add any attributes requested for all logs
	l.logRequired(entry)

	// Process the key/value pairs into the context
	addArgPairs(entry, args)

	return entry
}

// omitRequiredAttrsLog creates a new LogEntry with the defined log level, message
//
// # Any required attributes set at initialization will not be added to the log entry
func (l *Logger) omitRequiredAttrsLog(level, message string) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		LogID:     generateLogID(),
		Context:   make(map[string]interface{}),
	}

	return entry
}

// addArgPairs processes variadic string arg pairs and adds them to the context
func addArgPairs(entry *LogEntry, keyValuePairs []string) {
	for i := 0; i < len(keyValuePairs); i += 2 {
		key := keyValuePairs[i]

		// Check if there's a matching value
		if i+1 < len(keyValuePairs) {
			value := keyValuePairs[i+1]
			entry.Context[key] = value
		} else {
			// If no matching value, add key with an empty string
			entry.Context[key] = ""
		}
	}
}
