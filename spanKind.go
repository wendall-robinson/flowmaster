package traceflow

import "go.opentelemetry.io/otel/trace"

// SpanKind allows fluent setting of the span kind.
type SpanKind struct {
	trace *Trace
}

// SpanKind returns a SpanKind object that allows the caller to set the kind of the span.
// Span kinds define the role of the span in a distributed trace and categorize the span
// as one of the following:
// - Server: Indicates that the span represents a server-side operation.
// - Client: Indicates that the span represents a client-side operation.
// - Producer: Indicates that the span represents a message producer.
// - Consumer: Indicates that the span represents a message consumer.
//
// Example usage:
//
//	trace.SpanKind().Server().Start("operation")
//
// This method returns a SpanKind object, which can be used to set the appropriate
// span kind. If the SpanKind object has not been previously initialized, it is created
// and linked to the current trace.
//
// Notes:
//   - Setting the correct span kind is important for accurate tracing and categorization
//     of operations in distributed systems.
//   - Ensure that the span kind is set before the span is started to properly classify it.
func (t *Trace) SpanKind() *SpanKind {
	if t.spanKind == nil {
		t.spanKind = &SpanKind{trace: t}
	}

	return t.spanKind
}

// Server sets the span kind to server
func (sk *SpanKind) Server() *Trace {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindServer))
	return sk.trace
}

// Client sets the span kind to client
func (sk *SpanKind) Client() *Trace {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindClient))
	return sk.trace
}

// Producer sets the span kind to producer
func (sk *SpanKind) Producer() *Trace {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindProducer))
	return sk.trace
}

// Consumer sets the span kind to consumer
func (sk *SpanKind) Consumer() *Trace {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindConsumer))
	return sk.trace
}
