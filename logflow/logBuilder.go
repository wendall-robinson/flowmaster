package logflow

import "time"

// logBuilder is a helper struct for building log entries
type logBuilder struct {
	message string
	level   string
	args    []string
}

// baseLog creates a new LogEntry with the given log level and message
// if logging was initialized with required attributes, they will be added to the log entry
func (l *Logger) baseLog(builder *logBuilder) *LogEntry {
	entry := &LogEntry{
		Timestamp: time.Now(),
		Level:     builder.level,
		Message:   builder.message,
		LogID:     generateLogID(),
		Context:   make(map[string]interface{}),
	}

	// add any attributes requested for all logs
	l.logRequested(entry)

	// Process the key/value pairs into the context
	addKVPairsToContext(entry, builder.args)

	return entry
}
