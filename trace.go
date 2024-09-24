// Package gotraceit provides a simple wrapper around OpenTelemetry to make it easier to create and manage traces.
package gotraceit

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Tracex is a struct that holds the context, tracer, span, and attributes for a trace
type Tracex struct {
	ctx          context.Context
	service      string
	tracer       trace.Tracer
	span         trace.Span
	parentSpanID string
	attrs        []attribute.KeyValue
	options      []trace.SpanStartOption
	spanKind     *SpanKindSetter
	links        []trace.Link
}

// NewTracex creates a new trace object using the specified tracer from the OpenTelemetry provider.
// If context is nil, as a valid context is created for the trace. This context is used for
// trace propagation and management.
// NewTracex creates a new trace object and captures the parent span ID if available.
func NewTracex(ctx context.Context, spanName string) *Tracex {
	var (
		traceCtx     context.Context
		parentSpanID string
	)

	if ctx == nil {
		traceCtx = context.Background()
	} else {
		traceCtx = ctx
		if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
			parentSpanID = span.SpanContext().SpanID().String()
		}
	}

	return &Tracex{
		ctx:          traceCtx,
		service:      spanName,
		tracer:       otel.GetTracerProvider().Tracer(spanName),
		parentSpanID: parentSpanID,
		attrs:        []attribute.KeyValue{},
		options:      []trace.SpanStartOption{},
	}
}

// Start creates a new span within the existing trace with and a name
func (t *Tracex) Start(name string) *Tracex {
	if len(t.attrs) > 0 {
		t.options = append(t.options, trace.WithAttributes(t.attrs...))
	}

	if len(t.links) > 0 {
		t.options = append(t.options, trace.WithLinks(t.links...))
	}

	operation := fmt.Sprintf("%s.%s", t.service, name)
	t.ctx, t.span = t.tracer.Start(t.ctx, operation, t.options...)

	return t
}

// AddAttributes adds attributes to the trace
func (t *Tracex) AddAttributes(attrs ...attribute.KeyValue) *Tracex {
	t.attrs = append(t.attrs, attrs...)

	return t
}

// AddAttributeIf adds an attribute to the trace if the condition is true
func (t *Tracex) AddAttributeIf(cond bool, key string, value interface{}) *Tracex {
	if cond {
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
	}

	return t
}

// AddKeyValue adds attributes to the trace
func (t *Tracex) AddKeyValue(key, value string) *Tracex {
	attr := attribute.String(key, value)
	t.attrs = append(t.attrs, attr)

	return t
}

// AddJSON adds JSON to the span attributes as a string
func (t *Tracex) AddJSON(payload json.RawMessage) *Tracex {
	jsonString := string(payload)
	jsonAttr := attribute.String("payload", jsonString)
	t.attrs = append(t.attrs, jsonAttr)

	return t
}

// AddLink adds a link to another span to this Trace's span.
// The linked span can be from the same or different trace.
func (t *Tracex) AddLink(spanContext trace.SpanContext) *Tracex {
	link := trace.Link{SpanContext: spanContext}
	t.links = append(t.links, link)

	return t
}

// InjectHTTPContext injects the trace context into an HTTP request's headers.
// This method is a no-op if t.ctx is nil, reflecting situations where the Trace
// was improperly instantiated or if context handling was bypassed.
func (t *Tracex) InjectHTTPContext(req *http.Request) *Tracex {
	if t.ctx == nil {
		return t
	}

	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(req.Header)
	propagator.Inject(t.ctx, carrier)

	return t
}

// ExtractHTTPContext extracts the trace context from an HTTP request and updates the Trace's context.
func (t *Tracex) ExtractHTTPContext(req *http.Request) *Tracex {
	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(t.ctx, propagation.HeaderCarrier(req.Header))
	t.ctx = ctx

	return t
}

// RecordError records an error to the span and sets the span status to Error.
func (t *Tracex) RecordError(err error) {
	if err != nil {
		t.span.RecordError(err)
		t.span.SetStatus(codes.Error, err.Error())
	}
}

// SetStatus sets the span status with a custom code and message.
func (t *Tracex) SetStatus(code codes.Code, message string) {
	t.span.SetStatus(code, message)
}

// SetSuccess marks the span status as Ok and sets a custom success message.
func (t *Tracex) SetSuccess(message string) {
	t.span.SetStatus(codes.Ok, message)
}

// RecordFailure records an error to the span, sets the span status to Error, and sets a custom message.
func (t *Tracex) RecordFailure(err error, message string) {
	t.RecordError(err)
	t.SetStatus(codes.Error, message)
}

// SpanKind sets the span kind to the provided value
func (t *Tracex) SpanKind() *SpanKindSetter {
	if t.spanKind == nil {
		t.spanKind = &SpanKindSetter{trace: t}
	}

	return t.spanKind
}

// GetTraceID returns the trace ID of the current span
func (t *Tracex) GetTraceID() string {
	if t.span != nil {
		sc := t.span.SpanContext()
		if sc.IsValid() {
			return sc.TraceID().String()
		}
	}

	return ""
}

// GetParentID returns the parent ID of the current span
func (t *Tracex) GetParentID() string {
	return t.parentSpanID
}

// GetContext returns the context of the current span
func (t *Tracex) GetContext() context.Context {
	return t.ctx
}

// End ends the current span
func (t *Tracex) End() {
	if t.span != nil {
		t.span.End()
	}
}
