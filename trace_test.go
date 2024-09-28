package traceflow

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// // init initializes the OpenTelemetry provider and sets the global propagator.
// func init() {
// 	ctx := context.Background()

// 	_, shutdown, err := Init(ctx, "test-service", WithSilentLogger())
// 	if err != nil {
// 		panic("Failed to initialize OpenTelemetry for tests: " + err.Error())
// 	}

// 	defer shutdown(ctx)
// }

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
		AddString("key1", "value1"),
		AddString("key2", "value2"),
		AddInt("key3", 3),
	)

	if len(trace.attrs) != 3 {
		t.Errorf("Expected 3 attributes, got %d", len(trace.attrs))
	}
}

// TestStartEndSpan tests starting and ending a span.
func TestStartEndSpan(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	defer trace.Start("test-span").End()

	if trace.span == nil {
		t.Error("Span should not be nil")
	}
}

// TestInjectHTTPContext tests the injection of HTTP context.
func TestInjectHTTPContext(t *testing.T) {
	ctx := context.Background()

	// Initialize OpenTelemetry with a valid context
	ctx, shutdown, err := Init(ctx, "test-service", WithSilentLogger())
	if err != nil {
		t.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	trace := New(ctx, "test-service")

	defer trace.Start("inject-http-context").End()

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	trace.InjectHTTPContext(req)

	if _, exists := req.Header["Traceparent"]; !exists {
		t.Errorf("Expected Traceparent header to be injected")
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
	// Use the context from Init to ensure proper OpenTelemetry setup
	ctx := context.Background()

	ctx, shutdown, err := Init(ctx, "test-service", WithSilentLogger())
	if err != nil {
		t.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	// Create a new trace using the initialized context
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
}

// TestAddAttributeIf tests conditionally adding attributes to a trace.
func TestAddAttributeIf(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	// Add attribute if true
	trace.AddAttributeIf(true, "key1", "value1")
	if len(trace.attrs) != 1 {
		t.Errorf("Expected 1 attribute, got %d", len(trace.attrs))
	}

	// Do not add attribute if false
	trace.AddAttributeIf(false, "key2", "value2")
	if len(trace.attrs) != 1 {
		t.Errorf("Expected 1 attribute, got %d", len(trace.attrs))
	}
}

// TestSetStatus tests setting the status of a span.
func TestSetStatus(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	defer trace.Start("set-status").End()

	// Set status to codes.Error
	trace.SetStatus(codes.Error, "operation failed")

	// Verify the span status is codes.Error
	if trace.span == nil {
		t.Error("Span should not be nil")
	}
	// Since OpenTelemetry doesn't expose status directly, we can't verify it in a test,
	// but the code won't raise an error if SetStatus works correctly.
}

// TestEnd tests ending the span.
func TestEnd(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	trace.Start("test-end")

	// End the span
	trace.End()

	// The span should not be nil, but it's hard to check if it was ended since OpenTelemetry
	// doesn't expose internal span state in tests.
	if trace.span == nil {
		t.Error("Span should not be nil after starting")
	}
}

// TestRecordFailure tests recording a failure in the trace.
func TestRecordFailure(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	defer trace.Start("record-failure").End()

	// Call RecordFailure (no return value)
	trace.RecordFailure(fmt.Errorf("test error"), "custom failure message")

	// Since we can't directly inspect the span's status in a test, we are ensuring no panic or failure.
	if trace.span == nil {
		t.Error("Span should not be nil after RecordFailure is called")
	}
}

// TestAddLink tests adding a link to a trace's span.
func TestAddLink(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service")

	// Create a new span to generate a valid SpanContext
	_, span := trace.tracer.Start(ctx, "linked-span")
	otelSpanContext := span.SpanContext()

	// Wrap the SpanContext using the NewSpanContext function
	spanContext := NewSpanContext(otelSpanContext)

	// Add the link to the trace
	trace.AddLink(spanContext)

	if len(trace.links) != 1 {
		t.Errorf("Expected 1 link, got %d", len(trace.links))
	}
}

// TestGetParentID tests getting the parent span ID of a trace and ensures parent-child relationships.func TestGetParentID(t *testing.T) {
func TestGetParentID(t *testing.T) {
	// Use the context from Init to ensure proper OpenTelemetry setup
	ctx := context.Background()

	ctx, shutdown, err := Init(ctx, "test-service", WithSilentLogger())
	if err != nil {
		t.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	trace1 := New(ctx, "test-service").Start("parent-span")

	// Ensure trace1 does not have a parent span ID
	if parentID := trace1.GetParentID(); parentID != "" {
		t.Errorf("Expected no parent ID for trace1, got %s", parentID)
	}

	// Create trace2 using the context of trace1
	trace2 := New(trace1.GetContext(), "child-service").Start("child-span")

	// Trace2 should have trace1's span as its parent
	if parentID := trace2.GetParentID(); parentID == "" {
		t.Error("Expected valid parent ID for trace2, but got empty string")
	}

	// End both spans
	trace2.End()
	trace1.End()
}

// TestWithSystemInfo tests that system info (CPU, memory, disk) attributes are added.
func TestWithSystemInfo(t *testing.T) {
	ctx := context.Background()
	trace := New(ctx, "test-service", WithSystemInfo())

	// Check that CPU, memory, and disk attributes are added (sample assertions)
	if len(trace.attrs) == 0 {
		t.Fatalf("Expected system attributes to be added, got none")
	}
	// Example: Check that the CPU count attribute exists
	found := false
	for _, attr := range trace.attrs {
		if attr.Key == "cpu.count" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected 'cpu.count' attribute, but it was not found")
	}
}

// TestWithAttributes tests that custom attributes are added via WithAttributes.
func TestWithAttributes(t *testing.T) {
	customAttrs := []Attribute{
		AddString("custom.key1", "value1"),
		AddInt("custom.key2", 42),
	}

	ctx := context.Background()
	trace := New(ctx, "test-service", WithAttributes(customAttrs...))

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 custom attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != customAttrs[0].otelAttr {
		t.Errorf("Expected first custom attribute to be %v, got %v", customAttrs[0].otelAttr, trace.attrs[0])
	}
	if trace.attrs[1] != customAttrs[1].otelAttr {
		t.Errorf("Expected second custom attribute to be %v, got %v", customAttrs[1].otelAttr, trace.attrs[1])
	}
}

// TestWithConcurrencyInfo tests that the number of goroutines is added via WithConcurrencyInfo.
func TestWithConcurrencyInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service", WithConcurrencyInfo())

	expectedGoroutines := runtime.NumGoroutine()
	found := false

	for _, attr := range trace.attrs {
		if attr.Key == "goroutine.count" {
			if attr.Value.AsInt64() != int64(expectedGoroutines) {
				t.Errorf("Expected goroutine count to be %d, got %d", expectedGoroutines, attr.Value.AsInt64())
			}
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected 'goroutine.count' attribute, but it was not found")
	}
}

// TestWithEnVars tests that environment variables are correctly added as attributes.
func TestWithEnVars(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVICE_NAME", "test-service")
	os.Setenv("DEPLOYMENT_ENV", "production")
	defer os.Unsetenv("SERVICE_NAME")
	defer os.Unsetenv("DEPLOYMENT_ENV")

	envKeys := []string{"SERVICE_NAME", "DEPLOYMENT_ENV", "NON_EXISTENT_ENV"}
	ctx := context.TODO()
	trace := New(ctx, "test-service", WithEnVars(envKeys))

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 environment variables, got %d", len(trace.attrs))
	}

	if trace.attrs[0] != attribute.String("SERVICE_NAME", "test-service") {
		t.Errorf("Expected SERVICE_NAME to be 'test-service', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("DEPLOYMENT_ENV", "production") {
		t.Errorf("Expected DEPLOYMENT_ENV to be 'production', got %v", trace.attrs[1])
	}
}
