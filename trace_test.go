package traceflow

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

// init initializes the OpenTelemetry provider and sets the global propagator.
func init() {
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
	)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	otel.SetTracerProvider(tp)
}

// TestNewTrace tests the creation of a new Trace instance.
func TestNewTrace(t *testing.T) {
	ctx := context.Background()
	tracerName := "test-service"

	trace := New(ctx, tracerName)

	if trace.tracer == nil {
		t.Error("Tracer should not be nil")
	}
}

// TestAddAttributes tests adding attributes to a trace.
func TestAddAttributes(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	trace.AddAttribute(
		attribute.String("key1", "value1"),
		attribute.String("key2", "value2"),
		attribute.Int("key3", 3),
	)

	if len(trace.attrs) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(trace.attrs))
	}
}

// TestInjectHTTPContext tests the injection of HTTP context.
func TestInjectHTTPContext(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	defer trace.Start("inject-http-context").End()

	//nolint:noctx
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	trace.InjectHTTPContext(req)

	if _, exists := req.Header["Traceparent"]; !exists {
		t.Error("Expected Traceparent header to be injected")
	}
}

// TestExtractHTTPContext tests the extraction of HTTP context.
func TestExtractHTTPContext(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	//nolint:noctx
	req, _ := http.NewRequest("GET", "http://example.com", nil)

	// Simulating a Traceparent header being present
	req.Header.Add("Traceparent", "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	trace.ExtractHTTPContext(req)

	// Testing if the context was properly updated
	if trace.ctx == ctx {
		t.Error("Context should have been updated with extracted values")
	}
}

func TestGetTraceID(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	// Initially, before the span is started, TraceID should be empty
	if traceID := trace.GetTraceID(); traceID != "" {
		t.Errorf("Expected no Trace ID before span is started, got %s", traceID)
	}

	// Start the span and defer its end
	defer trace.Start("get-trace-id").End()

	// After starting the span, check for a valid Trace ID
	if traceID := trace.GetTraceID(); traceID == "" {
		t.Error("Expected valid Trace ID after span is started")
	}

	fmt.Printf("Trace: %+v\n", trace.GetTraceID())
}
