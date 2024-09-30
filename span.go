package traceflow

import "go.opentelemetry.io/otel/trace"

// SpanContext wraps the OpenTelemetry span context type.
type SpanContext struct {
	otelSpanContext trace.SpanContext
}

// SpanKind allows fluent setting of the span kind.
type SpanKind struct {
	trace  *Trace
	option trace.SpanStartOption
}

// NewSpanContext creates a new SpanContext from OpenTelemetry's span context.
func NewSpanContext(sc trace.SpanContext) SpanContext {
	return SpanContext{otelSpanContext: sc}
}

// Server sets the span kind to server and returns the Trace object for chaining.
func (t *Trace) Server() *Trace {
	if t.spanKind == nil {
		t.spanKind = &SpanKind{trace: t}
	}

	t.spanKind.option = trace.WithSpanKind(trace.SpanKindServer)

	return t
}

// Client sets the span kind to client and returns the Trace object for chaining.
func (t *Trace) Client() *Trace {
	if t.spanKind == nil {
		t.spanKind = &SpanKind{trace: t}
	}

	t.spanKind.option = trace.WithSpanKind(trace.SpanKindClient)

	return t
}

// Producer sets the span kind to producer and returns the Trace object for chaining.
func (t *Trace) Producer() *Trace {
	if t.spanKind == nil {
		t.spanKind = &SpanKind{trace: t}
	}

	t.spanKind.option = trace.WithSpanKind(trace.SpanKindProducer)

	return t
}

// Consumer sets the span kind to consumer and returns the Trace object for chaining.
func (t *Trace) Consumer() *Trace {
	if t.spanKind == nil {
		t.spanKind = &SpanKind{trace: t}
	}

	t.spanKind.option = trace.WithSpanKind(trace.SpanKindConsumer)

	return t
}
