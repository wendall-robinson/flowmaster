package logflow

import "time"

// Error creates a new Error level LogEntry with an error and key/value pairs
func (l *Logger) Error(err error, keyValuePairs ...string) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		LogID:     generateLogID(),
		Context:   make(map[string]interface{}),
	}

	// Add error message and other context
	entry.Context["error"] = err.Error()

	// Process key/value pairs into the context
	addKVPairsToContext(entry, keyValuePairs)

	// Process the log entry
	processLog(entry)

	return entry
}

// ErrorSys creates a new Error level LogEntry with system attributes, an error, and key/value pairs
func (l *Logger) ErrorSys(err error, keyValuePairs ...string) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		LogID:     generateLogID(),
		Context:   make(map[string]interface{}),
	}

	// Add system attributes
	l.findSystemAttributes(entry)

	// Add error message and other context
	entry.Context["error"] = err.Error()

	// Process key/value pairs into the context
	addKVPairsToContext(entry, keyValuePairs)

	// Process the log entry
	processLog(entry)

	return entry
}
