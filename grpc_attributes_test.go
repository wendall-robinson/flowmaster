package traceflow

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc/metadata"
)

func TestInjectGRPCContext(t *testing.T) {
	ctx := context.Background()

	// Initialize OpenTelemetry with silent logger for testing
	ctx, shutdown, err := Init(ctx, "test-service", WithSilentLogger())
	if err != nil {
		t.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	// Start a new trace
	trace := New(ctx, "test-service")
	defer trace.Start("grpc-operation").End()

	// Inject trace context into gRPC metadata
	newCtx := trace.InjectGRPCContext(ctx)

	// Extract metadata from the new context
	md, ok := metadata.FromOutgoingContext(newCtx)
	assert.True(t, ok, "Expected metadata to be present after injecting context")
	assert.NotEmpty(t, md, "Expected non-empty metadata after injecting context")

	// Check that the "traceparent" header is present
	traceparent, exists := md["traceparent"]
	assert.True(t, exists, "Expected 'traceparent' header to be present in metadata")
	assert.NotEmpty(t, traceparent, "Expected 'traceparent' header to be non-empty")

	// Ensure that the TracerProvider is valid and not a noop provider
	tp := otel.GetTracerProvider()
	assert.NotEqual(t, noop.NewTracerProvider(), tp, "Expected a valid TracerProvider to be set")
}
