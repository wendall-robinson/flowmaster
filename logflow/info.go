package logflow

import (
	"net/http"
)

// Info creates a new Info level LogEntry with key/value pairs
func (l *Logger) Info(message string, args ...string) *LogEntry {
	entry := l.baseLog(&logBuilder{
		message: message,
		level:   "INFO",
		args:    args,
	})

	// Process the log entry
	processLog(entry)

	return entry
}

// InfoSys creates a new Info level LogEntry with system attributes and key/value pairs
func (l *Logger) InfoSys(message string, args ...string) *LogEntry {
	entry := l.baseLog(&logBuilder{
		message: message,
		level:   "INFO",
		args:    args,
	})

	// Add system attributes
	l.findSystemAttributes(entry)

	// Process the log entry (optional chaining can modify it further)
	processLog(entry)

	return entry
}

// InfoHttp creates a new Info level LogEntry with HTTP request/response metadata and key/value pairs
func (l *Logger) InfoHttp(message string, req *http.Request, statusCode int, args ...string) *LogEntry {
	entry := l.baseLog(&logBuilder{
		message: message,
		level:   "INFO",
		args:    args,
	})

	// Add HTTP metadata using the HTTP struct
	entry.Context["http"] = extractHttpMetadata(req, statusCode)

	// Process the log entry
	processLog(entry)

	return entry
}
