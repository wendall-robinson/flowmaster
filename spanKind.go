package traceflow

import "go.opentelemetry.io/otel/trace"

// SpanKind allows fluent setting of the span kind.
type SpanKind struct {
	trace  *Trace
	option trace.SpanStartOption
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
func (s *SpanKind) Server() *Trace {
	s.option = trace.WithSpanKind(trace.SpanKindServer)
	return s.trace
}

// Client sets the span kind to client
func (s *SpanKind) Client() *Trace {
	s.option = trace.WithSpanKind(trace.SpanKindClient)
	return s.trace
}

// Producer sets the span kind to producer
func (s *SpanKind) Producer() *Trace {
	s.option = trace.WithSpanKind(trace.SpanKindProducer)
	return s.trace
}

// Consumer sets the span kind to consumer
func (s *SpanKind) Consumer() *Trace {
	s.option = trace.WithSpanKind(trace.SpanKindConsumer)
	return s.trace
}
