// Package traceflow provides a simple wrapper around OpenTelemetry to make it easier to create and manage traces.
package traceflow

import (
	"context"
	"fmt"
	"reflect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Trace is a struct that holds the context, tracer, span, and attributes for a trace
type Trace struct {
	ctx          context.Context
	service      string
	tracer       trace.Tracer
	span         trace.Span
	parentSpanID string
	attrs        []attribute.KeyValue
	options      []trace.SpanStartOption
	spanKind     *SpanKind
	links        []trace.Link
}

// New creates a new Trace object using the specified tracer from the OpenTelemetry provider.
// If context is nil, a valid context is created for the trace. This context is used for
// trace propagation and management. Users can pass variadic options to customize the trace,
// including automatically adding system-related attributes, custom attributes, or any other
// predefined behaviors.
//
// New also automatically propagates the trace context from
// the provided context. This ensures that the trace is linked to any existing parent trace.
// If no trace exists in the provided context, it starts a fresh trace.
//
// Example usage:
//
//	// Create a new Trace with default settings
//	trace := traceflow.New(ctx, "my-service")
//
//	// Create a new Trace with system information
//	trace := traceflow.New(ctx, "my-service", traceflow.WithSystemInfo())
//
//	// Create a new Trace with custom attributes
//	trace := traceflow.New(ctx, "my-service", traceflow.WithAttributes(
//	    attribute.String("user_id", "123"),
//	    attribute.Int("request_count", 5),
//	))
//
// Notes:
// - The options allow flexibility in configuring the Trace object during initialization.
// - You can create multiple options to fit various use cases and simplify tracing setup.
func New(ctx context.Context, spanName string, opts ...Option) *Trace {
	var traceCtx context.Context

	if ctx == nil {
		traceCtx = context.Background()
	} else {
		traceCtx = ctx
	}

	var parentSpanID string
	// If there's an existing span in the context, preserve its parent SpanContext
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		parentSpanID = span.SpanContext().SpanID().String()
	}

	// Create the Trace object
	t := &Trace{
		ctx:          traceCtx,
		service:      spanName,
		tracer:       otel.GetTracerProvider().Tracer(spanName),
		parentSpanID: parentSpanID,
		attrs:        []attribute.KeyValue{},
		options:      []trace.SpanStartOption{},
		spanKind:     &SpanKind{option: trace.WithSpanKind(trace.SpanKindInternal)},
	}

	// Apply variadic options
	for _, opt := range opts {
		opt(t)
	}

	return t
}

// NewWithoutPropagation creates a new Trace object without propagating the trace context
// from the provided context. This treats the current trace as a parent, but does not copy
// the context. Use this method when you want the new trace to start independently.
//
// Example usage:
//
//	// Create a new trace without propagating the existing context
//	trace := traceflow.NewWithoutPropagation(ctx, "my-service")
//
//	// The new trace will not be linked to the parent trace.
func NewWithoutPropagation(ctx context.Context, spanName string, opts ...Option) *Trace {
	traceCtx := context.Background()

	t := &Trace{
		ctx:      traceCtx,
		service:  spanName,
		tracer:   otel.GetTracerProvider().Tracer(spanName),
		attrs:    []attribute.KeyValue{},
		options:  []trace.SpanStartOption{},
		spanKind: &SpanKind{},
	}

	// Apply variadic options (if any)
	for _, opt := range opts {
		opt(t)
	}

	return t
}

// Start creates a new span within the existing trace using the provided name.
// It includes any attributes, links, and options that have been set on the trace.
// After the span is created, attributes, links, and options are cleared to avoid
// accidental reuse in future spans.
//
// If a span kind (e.g., server, client) has been set, it will also be applied to
// the span during creation.
//
// Example usage:
//
//	trace := traceflow.New(ctx, "my-service")
//	defer trace.Start("operation_name").End()
//
// Notes:
// - This method formats the operation name as "<service>.<name>".
// - Once a span is started, it must be ended using the End method.
func (t *Trace) Start(name string) *Trace {
	// Ensure a valid tracer exists
	if t.tracer == nil {
		t.tracer = otel.GetTracerProvider().Tracer(t.service)
	}

	// Apply attributes if they exist
	if len(t.attrs) > 0 {
		t.options = append(t.options, trace.WithAttributes(t.attrs...))
	}

	// Apply links if they exist
	if len(t.links) > 0 {
		t.options = append(t.options, trace.WithLinks(t.links...))
	}

	// Apply span kind if it exists
	t.options = append(t.options, t.spanKind.option)

	// Start the span
	operation := fmt.Sprintf("%s.%s", t.service, name)
	t.ctx, t.span = t.tracer.Start(trace.ContextWithSpan(t.ctx, trace.SpanFromContext(t.ctx)), operation, t.options...)

	// Clear attributes, links, and options after starting the span to avoid re-use
	t.attrs = nil
	t.links = nil
	t.options = nil

	return t
}

// Now creates and starts a trace with options immediately.
func Now(ctx context.Context, name, operation string, opts ...Option) *Trace {
	trace := New(ctx, name, opts...)

	trace.Start(operation)

	return trace
}

// AddAttribute appends one or more OpenTelemetry attributes to the current trace.
// This method accepts variadic attribute.KeyValue arguments, allowing the caller
// to add single or multiple attributes in a single call. It supports both OpenTelemetry
// predefined attributes (e.g., String, Int, Bool) and custom attributes formatted
// as attribute.KeyValue objects.
//
// Example usage:
//
//	// Add a single attribute
//	trace.AddAttribute(attribute.String("user_id", "123"))
//
//	// Add multiple attributes
//	trace.AddAttribute(
//	    attribute.String("user_id", "123"),
//	    attribute.Int("http_status", 200),
//	    attribute.Bool("success", true),
//	)
//
// AddAttribute accepts one or more custom TraceFlow attributes and appends them to the trace.
func (t *Trace) AddAttribute(attrs ...Attribute) *Trace {
	for _, attr := range attrs {
		t.attrs = append(t.attrs, attr.otelAttr)
	}

	return t
}

// AddAttributeIf conditionally adds an attribute to the trace based on a boolean condition.
// If the condition (cond) is true, the attribute specified by the key and value is added to
// the trace. The method automatically determines the correct OpenTelemetry attribute type
// (e.g., string, int, float, bool) based on the value provided.
//
// Supported types for the value parameter include:
// - string
// - int, int32, int64
// - uint, uint32, uint64
// - float32, float64
// - bool
//
// If the value is of an unsupported type, the attribute is not added.
//
// Example usage:
//
//	// Conditionally add an attribute only if the user ID is valid
//	trace.AddAttributeIf(isValidUser, "user_id", "12345")
//
//	// Conditionally add a numeric attribute
//	trace.AddAttributeIf(isSuccess, "response_time", 150)
//
// This method is particularly useful when attributes should only be included in the trace
// under specific conditions (e.g., based on business logic or performance metrics).
//
// If the condition is false, no attribute is added, and the trace remains unchanged.
func (t *Trace) AddAttributeIf(cond bool, key string, value interface{}) *Trace {
	if !cond {
		return t
	}

	var attr attribute.KeyValue
	switch v := value.(type) {
	case string:
		attr = attribute.String(key, v)
	case int, int32, int64:
		attr = attribute.Int64(key, reflect.ValueOf(v).Int())
	case uint, uint32, uint64:
		attr = attribute.Int64(key, int64(reflect.ValueOf(v).Uint()))
	case float32, float64:
		attr = attribute.Float64(key, reflect.ValueOf(v).Float())
	case bool:
		attr = attribute.Bool(key, v)
	default:
		return t
	}

	t.attrs = append(t.attrs, attr)

	return t
}

// AddLink adds a link to another span within the current traceflow span.
// Span links are used to connect spans that are related but do not have a direct
// parent-child relationship. This is useful when spans from different traces or
// parts of the same trace are logically related and should be connected.
//
// The linked span can come from either the same trace or a different trace.
// The link helps trace analyzers understand the relationships between otherwise
// unrelated spans and allows them to be visualized as part of a larger operation.
//
// Example usage:
//
//	// Link another span's context to the current span
//	otherSpanContext := traceflow.NewSpanContext(externalOtelSpanContext)
//	trace.AddLink(otherSpanContext)
//
// This method is particularly useful in scenarios like batch processing, where a
// single span may be related to multiple spans that are processed together, but
// do not have direct hierarchical relationships.
//
// Notes:
//   - The linked span is represented by its traceflow.SpanContext, which wraps the
//     OpenTelemetry span context (trace.SpanContext).
//   - This method returns the Trace object, allowing chaining of additional methods.
//   - Users interact with traceflow's custom SpanContext type to avoid the need to directly
//     import or use OpenTelemetry types, making tracing integration easier.
func (t *Trace) AddLink(spanContext SpanContext) *Trace {
	link := trace.Link{SpanContext: spanContext.otelSpanContext}
	t.links = append(t.links, link)

	return t
}

// SetStatus sets the status of the current span with a custom code and message.
// This is useful for recording the outcome of the operation represented by the span,
// providing context on whether the operation succeeded, failed, or encountered an error.
//
// The status code should be chosen from OpenTelemetry's predefined status codes (codes.Code),
// which include options like:
// - codes.Ok (indicating success)
// - codes.Error (indicating failure)
//
// The message should provide additional context or details about the status.
//
// Example usage:
//
//	// Set the span's status to indicate an error
//	trace.SetStatus(codes.Error, "database connection failed")
//
// Notes:
// - Ensure the span is properly started before setting its status.
// - This method allows the trace to capture both the status code and a descriptive message.
func (t *Trace) SetStatus(code codes.Code, message string) {
	t.span.SetStatus(code, message)
}

// SetSuccess marks the current span as successful by setting its status to codes.Ok,
// along with a custom success message. This is useful for marking the span as completed
// without any errors and providing a message that reflects the success.
//
// Example usage:
//
//	// Mark the span as successful with a custom message
//	trace.SetSuccess("operation completed successfully")
//
// This method is a shorthand for calling SetStatus with codes.Ok, simplifying the
// process of marking successful spans. It is particularly useful when you want to
// standardize how success is recorded in your traces.
func (t *Trace) SetSuccess(message string) {
	t.SetStatus(codes.Ok, message)
}

// RecordFailure records an error to the current span and marks the span's status as Error
// with a custom message. This method is useful for handling failure scenarios where
// both the error itself and a custom message need to be captured in the trace.
//
// The error is recorded using the RecordError method, and the span's status is set to
// codes.Error to reflect that the span represents a failed operation. A custom message
// is also provided to give additional context on the nature of the failure.
//
// Example usage:
//
//	// Record a failure in the span with an error and a custom message
//	trace.RecordFailure(err, "failed to process user request")
//
// This method is a convenient way to handle both error reporting and span status setting
// in failure cases, ensuring that the trace contains both the error details and the
// status information.
//
// Notes:
//   - Ensure the span is properly started before recording failures.
//   - The recorded error and message will be part of the trace and can be viewed in
//     trace analysis tools for debugging and diagnostics.
func (t *Trace) RecordFailure(err error, message string) {
	t.RecordError(err)
	t.SetStatus(codes.Error, message)
}

// GetTraceID returns the unique trace ID of the current span.
// The trace ID is part of the span's context and is used to identify the
// trace in a distributed system.
//
// If the span's context is valid, the trace ID is returned as a string.
// Otherwise, an empty string is returned.
//
// Example usage:
//
//	traceID := trace.GetTraceID()
//	fmt.Println("Current Trace ID:", traceID)
//
// Notes:
//   - The trace ID is useful for tracking and correlating traces across multiple
//     services in distributed systems.
func (t *Trace) GetTraceID() string {
	if t.span != nil {
		sc := t.span.SpanContext()
		if sc.IsValid() {
			return sc.TraceID().String()
		}
	}

	return ""
}

// FindTraceID returns the unique trace ID for a provided context.
//
// If the context contains a valid trace, the trace ID is returned as a string.
func FindTraceID(ctx context.Context) string {
	// finds the TraceID in a provided context
	sc := trace.SpanContextFromContext(ctx)
	if sc.IsValid() {
		return sc.TraceID().String()
	}

	return ""
}

// GetParentID returns the parent span ID of the current trace, if it exists.
// This ID represents the span from which the current span is derived, allowing
// the trace to establish relationships between parent and child spans.
//
// If the parent span ID is available, it is returned as a string. If no parent
// span exists, an empty string is returned.
//
// Example usage:
//
//	parentID := trace.GetParentID()
//	fmt.Println("Parent Span ID:", parentID)
//
// Notes:
//   - The parent span ID is important for understanding the hierarchy of spans
//     within a distributed trace.
func (t *Trace) GetParentID() string {
	return t.parentSpanID
}

// GetContext returns the context associated with the current span.
// The context carries metadata, including trace and span information,
// which is used for trace propagation across service boundaries.
//
// Example usage:
//
//	ctx := trace.GetContext()
//	// Use the context in subsequent operations
//
// This method is particularly useful when you need to pass the context
// to downstream services or operations that require trace propagation.
//
// Notes:
//   - Ensure that the context is valid and has been properly initialized before
//     passing it to other functions or services.
func (t *Trace) GetContext() context.Context {
	return t.ctx
}

// End marks the completion of the current span, signaling the end of the operation
// being traced. This method should be called after the span's operation has completed,
// allowing the trace to accurately record the duration and any final status or attributes
// of the span.
//
// Example usage:
//
//	// Start a span
//	trace.Start("operation")
//
//	// Perform some operations...
//
//	// End the span once the operation is complete
//	trace.End()
//
// Best Practices:
//
//	// Ensure that the span is always ended, even in the case of errors, by using defer:
//	defer trace.Start("operation").End()
//
// The End method is a critical part of the span lifecycle, as it ensures the span is
// properly closed and its data is recorded in the trace. If the span is nil, the method
// is a no-op, meaning it will do nothing.
//
// Notes:
//   - It is important to ensure that spans are always ended, either explicitly or
//     using defer to guarantee they are closed, even in the case of errors.
//   - Once a span has ended, no additional attributes or status can be set on it.
func (t *Trace) End() {
	if t.span != nil {
		t.span.End()
	}
}
