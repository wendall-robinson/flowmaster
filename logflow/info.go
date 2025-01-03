package logflow

import (
	"net/http"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// LEGEND:
// Info - INFO level log entry
// M - Message - A log entry with a user provided message
// R - Raw - A log entry with only the error message, no required attributes will be added to the log entry
// Sys - System - A log entry with system attributes attached
// T - Trace - A log entry with a trace ID attached if found in the provided context
// Fatal - A log entry that will panic the system after logging
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Info creates a new Info level LogEntry with key/value pairs
func (l *Logger) Info(message string, args ...string) {
	entry := l.baseEntryLog(INFO, message, args...)

	// Process the log entry
	l.processLog(entry)
}

// InfoR (Info Raw) creates a new Info level LogEntry with a only a message
//
// # If you have initialized the logger with required log event attributes, they will not be added to the log entry using this method.
//
// Example Usage:
//
//	logger.InfoR("This is an INFO message")
//
// Parameters:
//
//   - message: The log message describing the event.
func (l *Logger) InfoR(message string) {
	entry := l.omitRequiredAttrsLog(INFO, message)

	// Process the log entry
	l.processLog(entry)
}

// InfoSys creates a new Info level LogEntry with system attributes and key/value pairs
func (l *Logger) InfoSys(message string, args ...string) {
	entry := l.baseEntryLog(INFO, message, args...)

	// Add system attributes
	l.findSystemAttributes(entry)

	// Process the log entry (optional chaining can modify it further)
	l.processLog(entry)
}

// InfoHTTP creates a new Info level LogEntry with HTTP request/response metadata and key/value pairs
func (l *Logger) InfoHTTP(message string, req *http.Request, statusCode int, args ...string) {
	entry := l.baseEntryLog(INFO, message, args...)

	// Add HTTP metadata using the HTTP struct
	entry.Context["http"] = extractHTTPMetadata(req, statusCode)

	// Process the log entry
	l.processLog(entry)
}

// InfoCPU creates a new Info level LogEntry with a snapshot of CPU usage INFOrmation.
//
// This method captures and logs the current CPU usage, optionally allowing developers
// to add custom key/value pairs for additional context.
//
// Example Usage:
//
//	logger.InfoCPU("CPU usage", "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) InfoCPU(message string, args ...string) {
	const FIELD = "CPUUsagePercent"

	entry := l.baseEntryLog(INFO, message, args...)

	attributes, exists := systemAttribuesExists(entry.Context)
	switch exists {
	case true:
		// If attributes exist but CPU is not set, populate it
		ok := systemAttribuesFieldExists(attributes, FIELD)
		if !ok {
			attributes.CPUUsagePercent = l.getCPUUsage()
		}

	case false:
		// Add new SystemAttributes to the context
		attributes.CPUUsagePercent = l.getCPUUsage()
		entry.Context["system_attributes"] = attributes
	}

	// Process the log entry
	l.processLog(entry)
}

// InfoMem creates a new Info level LogEntry with a snapshot of Memory usage INFOrmation.
//
// This method captures and logs the current memory usage, optionally allowing developers
// to add custom key/value pairs for additional context.
//
// Example Usage:
//
//	logger.InfoMem("Memory usage", "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) InfoMem(message string, args ...string) {
	const FIELD = "MemoryUsage"

	entry := l.baseEntryLog(INFO, message, args...)

	attributes, exists := systemAttribuesExists(entry.Context)
	switch exists {
	case true:
		// If the attributes exist but the specific field is not set, populate it
		ok := systemAttribuesFieldExists(attributes, FIELD)
		if !ok {
			attributes.MemoryUsage = l.getMemoryUsage()
		}

	case false:
		// Add new SystemAttributes to the context
		attributes.MemoryUsage = l.getMemoryUsage()
		entry.Context["system_attributes"] = attributes
	}

	// Process the log entry
	l.processLog(entry)
}

// InfoDisk creates a new Info level LogEntry with a snapshot of Disk usage INFOrmation.
//
// This method captures and logs the current disk usage, optionally allowing developers
// to add custom key/value pairs for additional context.
//
// Example Usage:
//
//	logger.InfoDisk("Disk usage", "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) InfoDisk(message string, args ...string) {
	const (
		FIELD = "DiskUsage"
		PATH  = "/"
	)

	entry := l.baseEntryLog(INFO, message, args...)

	attributes, exists := systemAttribuesExists(entry.Context)
	switch exists {
	case true:
		// If the attributes exist but the specific field is not set, populate it
		ok := systemAttribuesFieldExists(attributes, FIELD)
		if !ok {
			attributes.DiskUsage = l.getDiskUsage(PATH)
		}

	case false:
		// Add new SystemAttributes to the context
		attributes.DiskUsage = l.getDiskUsage(PATH)
		entry.Context["system_attributes"] = attributes
	}

	// Process the log entry
	l.processLog(entry)
}

// InfoOS creates a new Info level LogEntry with a snapshot of the operating system.
//
// This method captures and logs the operating system, optionally allowing developers
// to add custom key/value pairs for additional context.
//
// Example Usage:
//
//	logger.InfoOS("Operating System", "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) InfoOS(message string, args ...string) {
	const FIELD = "OperatingSystem"

	entry := l.baseEntryLog(INFO, message, args...)

	attributes, exists := systemAttribuesExists(entry.Context)
	switch exists {
	case true:
		// If the attributes exist but the specific field is not set, populate it
		ok := systemAttribuesFieldExists(attributes, FIELD)
		if !ok {
			attributes.OperatingSystem = l.getOperatingSystem()
		}

	case false:
		// Add new SystemAttributes to the context
		attributes.OperatingSystem = l.getOperatingSystem()
		entry.Context["system_attributes"] = attributes
	}

	// Process the log entry
	l.processLog(entry)
}

// InfoProc creates a new Info level LogEntry with a snapshot of the running processes.
//
// This method captures and logs the current running processes, optionally allowing developers
// to add custom key/value pairs for additional context.
//
// Example Usage:
//
//	logger.InfoProc("Running processes", "key1", "value1","key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - args: Optional variadic key/value pairs to include additional context in the log.
//
// InfoProcesses creates a new Info level LogEntry with running processes INFOrmation
func (l *Logger) InfoProc(message string, args ...string) {
	const FIELD = "RunningProcesses"

	entry := l.baseEntryLog(INFO, message, args...)

	attributes, exists := systemAttribuesExists(entry.Context)
	switch exists {
	case true:
		// If the attributes exist but the specific field is not set, populate it
		ok := systemAttribuesFieldExists(attributes, FIELD)
		if !ok {
			attributes.RunningProcesses = l.getRunningProcesses()
		}

	case false:
		// Add new SystemAttributes to the context
		attributes.RunningProcesses = l.getRunningProcesses()
		entry.Context["system_attributes"] = attributes
	}

	// Process the log entry
	l.processLog(entry)
}
