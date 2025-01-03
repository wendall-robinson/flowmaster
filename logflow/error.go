package logflow

import (
	"context"

	"github.com/wendall-robinson/flowmaster/traceflow"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// LEGEND:
// Err - ERROR level log entry
// M - MESSAGE - A log entry with a user provided message
// R - RAW - A log entry with only the error message, no required attributes will be added to the log entry
// S - SYSTEM - A log entry with system attributes attached
// T - TRACE - A log entry with a trace ID attached if found in the provided context
// F - FATAL A log entry that will panic the system after logging
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// Err creates a new Error level LogEntry with the user message, error and optional variadic args
//
// If you have initialized the logger with required log event attributes, they will also be added to the log entry using this method.
//
// Example Usage:
//
//	logger.Err("This is an Error", err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - err: The error that occurred.
//   - args: Optional variadic args to include additional context in the log.
func (l *Logger) Err(message string, err error, args ...string) {
	entry := l.baseEntryLog(ERROR, message, args...)

	entry.ErrorMsg = err.Error()

	// Process the log entry
	l.processLog(entry)
}

// ErrR creates a new Error level LogEntry
//
// # This method is a shorthand for Error Raw which will only log the error and the optional variadic args
//
// If you have initialized the logger with required log event attributes, they will not be added to the log entry using this method.
// This gives users the flexibility to opt out of adding any configured required attributes to the log entry when needed
//
// Example Usage:
//
//	logger.ErrR("This is an ERROR RAW message", err.Error(), "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - error: The error that occurred.
//   - args: Optional variadic args to include additional context in the log.
func (l *Logger) ErrR(err error) {
	entry := l.omitRequiredAttrsLog(ERROR, err.Error())

	entry.ErrorMsg = err.Error()

	// Process the log entry
	l.processLog(entry)
}

// ErrM (Error Message) creates a new Error level LogEntry with a message and optional variadic args
//
// # This methods does not require the error to be passed in
//
// Example Usage:
//
//	logger.ErrorM("This is an ERROR message", "key1", "value1",
//		"key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - args: Optional variadic arg pairs to include additional context in the log.
func (l *Logger) ErrM(message string, args ...string) {
	entry := l.baseEntryLog(ERROR, message, args...)

	entry.ErrorMsg = message

	// Process the log entry
	l.processLog(entry)
}

// ErrF (Error Fatal) creates a new Error level LogEntry with an error and variadic args
//
// # This method will panic the system after logging the error
//
// Example Usage:
//
//	logger.ErrF(err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - err: The error that occurred.
//   - args: Optional variadic args pairs to include additional context in the log.
func (l *Logger) ErrF(err error, args ...string) {
	entry := l.baseEntryLog(ERROR, err.Error(), args...)

	entry.ErrorMsg = err.Error()

	// Process the log entry
	l.processLog(entry)

	panic(err)
}

// ErrMF (Fatal Error with Message) creates a new Error level LogEntry with a message, an error, and optional variadic args
//
// # This method will panic the system after logging the error
//
// Example Usage:
//
//	logger.ErrMF("This is a FATAL ERROR message", err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - message: The log message describing the event.
//   - err: The error that occurred.
//   - args: Optional variadic args to include additional context in the log.
func (l *Logger) ErrMF(message string, err error, args ...string) {
	entry := l.baseEntryLog(ERROR, message, args...)

	entry.ErrorMsg = err.Error()

	// Process the log entry
	l.processLog(entry)

	panic(err)
}

// ErrS (Error with System Attributes) creates a new Error level LogEntry with system attributes, an error, and variadic args
//
// This method will log the error and include system attributes in the log entry
// Example Usage:
//
//	logger.ErrS(err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - err: The error that occurred.
//   - args: Optional variadic arg pairs to include additional context in the log.
func (l *Logger) ErrS(err error, args ...string) {
	entry := l.baseEntryLog(ERROR, err.Error(), args...)

	// Add system attributes
	l.findSystemAttributes(entry)

	entry.ErrorMsg = err.Error()

	// Process the log entry
	l.processLog(entry)
}

// ErrSF (Fatal Error with System Attributes) creates a new Error level LogEntry with system attributes, an error, and variadic args
//
// # This method will panic the system after logging the error
//
// Example Usage:
//
//	logger.ErrorSysP(err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - err: The error that occurred.
//   - args: Optional variadic arg pairs to include additional context in the log.
func (l *Logger) ErrSF(err error, args ...string) {
	entry := l.baseEntryLog(ERROR, err.Error(), args...)

	// Add system attributes
	l.findSystemAttributes(entry)

	entry.ErrorMsg = err.Error()

	// Process the log entry
	l.processLog(entry)

	panic(err)
}

// ErrT creates a new Error level LogEntry with a trace ID, an error, and optional variadic args
//
// # This method will attempt to find the trace ID in the span of the context
//
// Example Usage:
//
//	logger.ErrorTrace(ctx, err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - ctx: The context to find the trace ID in the span.
//   - err: The error that occurred.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) ErrT(ctx context.Context, err error, args ...string) {
	entry := l.baseEntryLog(ERROR, err.Error(), args...)

	entry.ErrorMsg = err.Error()

	// use traceflow to find the trace ID in the span of the context if it exists
	entry.TraceID = traceflow.FindTraceID(ctx)

	// Process the log entry
	l.processLog(entry)
}

// ErrTM (Error Tracing with Message) creates a new Error level LogEntry with a trace ID, a message, an error, and optional variadic args
//
// # This method will attempt to find the trace ID in the span of the context
//
// Example Usage:
//
//	logger.ErrorTraceM(ctx, "This is an ERROR message", err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - ctx: The context to find the trace ID in the span.
//   - message: The log message describing the event.
//   - err: The error that occurred.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) ErrTM(ctx context.Context, message string, err error, args ...string) {
	entry := l.baseEntryLog(ERROR, message, args...)

	entry.ErrorMsg = err.Error()

	// use traceflow to find the trace ID in the span of the context if it exists
	entry.TraceID = traceflow.FindTraceID(ctx)

	// Process the log entry
	l.processLog(entry)
}

// ErrTF (Fatal Error with Tracing) creates a new Error level LogEntry with a trace ID, an error, and optional variadic args
//
// # This method will panic the system after logging the error
//
// Example Usage:
//
//	logger.ErrTF(ctx, err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - ctx: The context to find the trace ID in the span.
//   - err: The error that occurred.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) ErrTF(ctx context.Context, err error, args ...string) {
	entry := l.baseEntryLog(ERROR, err.Error(), args...)

	entry.ErrorMsg = err.Error()

	// use traceflow to find the trace ID in the span of the context if it exists
	entry.TraceID = traceflow.FindTraceID(ctx)

	// Process the log entry
	l.processLog(entry)

	panic(err)
}

// ErrTSF (Fatal Error with Tracing and System Attributes) creates a new Error level LogEntry with a trace ID, system attributes, error message, and optional variadic args
//
// # This method will panic the system after logging the error
//
// Example Usage:
//
//	logger.ErrTSF(ctx, err, "key1", "value1", "key2", "value2")
//
// Parameters:
//
//   - ctx: The context to find the trace ID in the span.
//   - err: The error that occurred.
//   - args: Optional variadic key/value pairs to include additional context in the log.
func (l *Logger) ErrTSF(ctx context.Context, err error, args ...string) {
	entry := l.baseEntryLog(ERROR, err.Error(), args...)

	entry.ErrorMsg = err.Error()

	// use traceflow to find the trace ID in the span of the context if it exists
	entry.TraceID = traceflow.FindTraceID(ctx)

	// Add system attributes
	l.findSystemAttributes(entry)

	// Process the log entry
	l.processLog(entry)

	panic(err)
}
