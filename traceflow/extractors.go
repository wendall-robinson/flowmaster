package traceflow

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// ExtractGRPCContext extracts the trace context from gRPC metadata.
// This is useful in a gRPC server to continue a trace initiated by an upstream service,
// ensuring that the trace context flows through the distributed system as part of the
// service request lifecycle.
//
// This method uses OpenTelemetry's propagator to extract the trace context from gRPC metadata
// present in the incoming request. The extracted trace context is then used to update the
// current Trace's context (t.ctx), allowing the service to join the existing trace and
// continue the tracing process.
//
// Example usage:
//
//		func (s *server) SomeRPC(ctx context.Context, req *pb.Request) (*pb.Response, error) {
//	     //Extract trace context from incoming gRPC request
//		    newCtx := traceflow.ExtractGRPCContext(ctx)
//
//		    trace := traceflow.New(newCtx, "grpc-server")
//
//
//		    // Continue the trace
//		    defer trace.Start("processing-request").End()
//
//		    // Handle request
//		    return &pb.Response{}, nil
//		}
//
// This method is particularly useful in distributed architectures where services need to
// propagate trace context with each request to maintain full trace visibility.
//
// Notes:
//   - The trace context is expected to be present in the incoming gRPC metadata in a format
//     compatible with OpenTelemetry's propagation standards (W3C Trace Context by default).
//   - If no metadata is found in the context, or the trace context is missing, the method
//     returns the original context unmodified.
func ExtractGRPCContext(ctx context.Context) context.Context {
	// Extract incoming metadata from the gRPC context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	// Convert gRPC metadata into a simple map for propagation purposes
	mdMap := make(map[string]string)

	for k, v := range md {
		if len(v) > 0 {
			mdMap[k] = v[0]
		}
	}

	// Extract the trace context from the metadata map
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.MapCarrier(mdMap)

	newCtx := propagator.Extract(ctx, carrier)

	return newCtx
}
