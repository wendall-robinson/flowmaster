// Package errors is a custom package that holds error definitions for traceflow.
package errors

import "fmt"

// ErrContextIsNil is returned when a trace is not found
var ErrContextIsNil = fmt.Errorf("context is nil")

// ErrTraceExporterCreation is returned when a trace exporter cannot be created
var ErrTraceExporterCreation = fmt.Errorf("failed to create trace exporter")

// ErrStdOutExporter is returned when a stdout trace exporter cannot be created
var ErrStdOutExporter = fmt.Errorf("failed to create stdout trace exporter")
