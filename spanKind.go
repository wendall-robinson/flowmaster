package gotraceit

import "go.opentelemetry.io/otel/trace"

// SpanKindSetter allows fluent setting of the span kind.
type SpanKindSetter struct {
	trace *Tracex
}

// Server sets the span kind to server
func (sk *SpanKindSetter) Server() *Tracex {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindServer))
	return sk.trace
}

// Client sets the span kind to client
func (sk *SpanKindSetter) Client() *Tracex {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindClient))
	return sk.trace
}

// Producer sets the span kind to producer
func (sk *SpanKindSetter) Producer() *Tracex {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindProducer))
	return sk.trace
}

// Consumer sets the span kind to consumer
func (sk *SpanKindSetter) Consumer() *Tracex {
	sk.trace.options = append(sk.trace.options, trace.WithSpanKind(trace.SpanKindConsumer))
	return sk.trace
}
