package traceflow

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InjectGRPCContext injects the trace context into the gRPC metadata.
// This is useful for propagating trace information across gRPC service boundaries.
// The trace context allows downstream services to continue the trace, providing
// full visibility into the flow of operations in a distributed system.
//
// This method uses OpenTelemetry's propagator to inject the trace context into
// gRPC metadata, ensuring that trace information is propagated with outgoing
// gRPC requests.
//
// After injection, the trace context is stored in the gRPC metadata and appended
// to the outgoing context, which will be used in the gRPC client to send the trace
// context along with the request.
//
// Example usage:
//
//	// Create a new trace in the gRPC client
//	trace := traceflow.New(ctx, "grpc-client")
//
//	// Inject the trace context into the outgoing gRPC request
//	ctx = trace.InjectGRPCContext(ctx)
//
//	// Perform a gRPC call with the injected trace context
//	response, err := client.SomeRPC(ctx, &pb.Request{})
//	if err != nil {
//		log.Fatalf("Failed to call gRPC service: %v", err)
//	}
//
//	// Continue trace logic if needed
//	trace.Start("some-operation").End()
//
// Notes:
//   - This method should be used in the gRPC client to propagate trace context to
//     downstream services.
//   - The trace context is appended to the context as gRPC metadata, using the W3C
//     Trace Context format by default.
//   - If the trace context (`t.ctx`) is nil or not properly initialized, this method
//     is a no-op and does not modify the outgoing context.
func (t *Trace) InjectGRPCContext(ctx context.Context) context.Context {
	if t.ctx == nil {
		return ctx
	}

	// Retrieve or create outgoing metadata from the provided context
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	// Create a map carrier and inject the trace context into it
	mdMap := map[string]string{}

	for k, v := range md {
		if len(v) > 0 {
			mdMap[k] = v[0]
		}
	}

	carrier := propagation.MapCarrier(mdMap)
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(t.ctx, carrier)

	for k, v := range mdMap {
		md.Set(k, v)
	}

	return metadata.NewOutgoingContext(ctx, md)
}

// UnaryClientInterceptor is a gRPC client interceptor that injects the trace context into
// outgoing gRPC requests. This allows the gRPC client to propagate the trace context to
// downstream services, ensuring that traces can be linked across distributed services.
//
// The interceptor intercepts every gRPC client call, extracts the current trace context
// from the client's context, and injects it into the outgoing gRPC metadata. This is
// essential for maintaining end-to-end visibility in distributed systems, allowing the
// trace to continue across multiple service boundaries.
//
// Example usage:
//
//	opts := []grpc.DialOption{
//	    grpc.WithUnaryInterceptor(traceflow.UnaryClientInterceptor()),
//	}
//	conn, err := grpc.Dial("localhost:50051", opts...)
//	if err != nil {
//	    log.Fatalf("Failed to connect: %v", err)
//	}
//	defer conn.Close()
//
//	client := pb.NewMyServiceClient(conn)
//	// Now all outgoing requests will have trace context injected
//
// Notes:
//   - This interceptor is designed for unary RPCs. For streaming RPCs, a different
//     interceptor (e.g., `StreamClientInterceptor`) is required.
//   - The trace context is injected using OpenTelemetry's propagator, and the trace context
//     is transmitted in a format that follows the W3C Trace Context standard.
//   - This interceptor should be included as part of the gRPC client options.
//
// Returns:
//   - A gRPC `grpc.UnaryClientInterceptor` function that can be added to the gRPC client
//     configuration to enable automatic trace context propagation.
func UnaryClientInterceptor(trace *Trace) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		// Inject the trace context into gRPC metadata
		ctx = trace.InjectGRPCContext(ctx)

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
