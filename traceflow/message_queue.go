package traceflow

import (
	"context"

	"go.opentelemetry.io/otel"
)

// MessageCarrier defines the interface for injecting and extracting trace context.
type MessageCarrier interface {
	Get(key string) string
	Set(key, value string)
	Keys() []string
}

// PropagateTraceContext injects the trace context from the given context into the carrier.
func PropagateTraceContext(ctx context.Context, carrier MessageCarrier) {
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(ctx, carrier)
}

// ExtractTraceContext extracts the trace context from the carrier and returns a new context.
func ExtractTraceContext(ctx context.Context, carrier MessageCarrier) context.Context {
	propagator := otel.GetTextMapPropagator()
	return propagator.Extract(ctx, carrier)
}
