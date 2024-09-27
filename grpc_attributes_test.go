package traceflow

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

func init() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
}

func TestInjectGRPCContext(t *testing.T) {
	// Create a context and initialize it with some metadata
	initialMD := metadata.New(map[string]string{
		"initial-key": "initial-value",
	})

	ctx := metadata.NewOutgoingContext(context.Background(), initialMD)

	// Create a new trace and start a span to ensure trace context is populated
	trace := New(ctx, "test-service")
	trace.Start("test-operation") // Start a span so that the trace context is populated

	// Inject the trace context into the gRPC context
	injectedCtx := trace.InjectGRPCContext(ctx)

	// Extract metadata from the outgoing context
	md, ok := metadata.FromOutgoingContext(injectedCtx)
	if !ok {
		t.Fatalf("Expected metadata to be present in outgoing context")
	}

	// Ensure that the metadata is non-empty
	if len(md) == 0 {
		t.Fatalf("Expected metadata to be non-empty")
	}

	// Check if the traceparent is present in the metadata
	if traceparent, ok := md["traceparent"]; !ok || len(traceparent) == 0 {
		t.Fatalf("Expected 'traceparent' header to be present in metadata")
	}
}
