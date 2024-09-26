package traceflow

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

// AddHTTPRequest adds HTTP request details as attributes to the trace.
func (t *Trace) AddHTTPRequest(req *http.Request) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("http.method", req.Method),
		attribute.String("http.url", req.URL.String()),
		attribute.String("http.user_agent", req.UserAgent()),
		attribute.String("http.client_ip", req.RemoteAddr),
	)

	return t
}

// AddHTTPResponse adds HTTP response details as attributes to the trace.
func (t *Trace) AddHTTPResponse(statusCode int, contentLength int64) *Trace {
	t.attrs = append(t.attrs,
		attribute.Int("http.status_code", statusCode),
		attribute.Int64("http.content_length", contentLength),
	)

	return t
}

// AddHTTPHeaders adds HTTP headers as attributes to the trace.
func (t *Trace) AddHTTPHeaders(headers http.Header) *Trace {
	for key, values := range headers {
		for _, value := range values {
			t.attrs = append(t.attrs, attribute.String("http.header."+key, value))
		}
	}

	return t
}

// InjectHTTPContext injects the trace context into the headers of an HTTP request.
// This is useful for propagating trace information across service boundaries in
// distributed systems, ensuring that the context of a trace is preserved as the
// request flows through multiple services.
//
// If the trace context (t.ctx) is nil, the method is a no-op, meaning it does not
// modify the request. This may occur if the Trace was improperly initialized or if
// context handling was bypassed.
//
// OpenTelemetry uses a "propagator" to encode the trace context as HTTP headers.
// The method leverages the global TextMapPropagator to inject the current trace context
// into the provided HTTP request's headers. The context is serialized and transmitted
// in a format compatible with OpenTelemetry's distributed tracing standards.
//
// Example usage:
//
//	// Create a new HTTP request
//	req, _ := http.NewRequest("GET", "http://example.com", nil)
//
//	// Inject the trace context into the request's headers
//	trace.InjectHTTPContext(req)
//
//	// The request now contains the trace context in its headers, which will
//	// be used by downstream services to continue the trace.
//
// This method is particularly useful in microservice architectures, where services
// need to pass trace information along with the request to maintain end-to-end
// visibility into the trace across multiple services.
//
// Notes:
// - Ensure that the trace context is properly initialized before calling this method.
// - The trace context will be injected into the request's headers using the W3C Trace Context format by default.
// - The method has no effect if the context is missing or invalid.
func (t *Trace) InjectHTTPContext(req *http.Request) *Trace {
	if t.ctx == nil {
		return t
	}

	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(req.Header)
	propagator.Inject(t.ctx, carrier)

	return t
}

// ExtractHTTPContext extracts the trace context from the headers of an HTTP request
// and updates the current Trace's context (t.ctx). This is used to continue a trace
// that was initiated by an upstream service, ensuring that the trace context flows
// through distributed systems without being lost.
//
// OpenTelemetry's TextMapPropagator is used to extract the trace context from the
// HTTP headers, where the upstream service has injected the trace information. This
// allows the current service to join the same trace by updating its context with
// the extracted trace information.
//
// The method updates t.ctx with the context extracted from the HTTP request's headers.
// This is necessary for services that are part of a distributed architecture where
// trace information needs to be propagated and shared between services.
//
// Example usage:
//
//	// Extract the trace context from an incoming HTTP request
//	trace := traceflow.NewTrace(ctx, "my-service")
//	trace.ExtractHTTPContext(req)
//
//	// The trace context is now updated, allowing the service to continue the trace.
//	trace.Start("operation").End()
//
// This method is useful in scenarios where your service is receiving HTTP requests
// from other services in the same distributed system. It ensures that the trace context
// initiated by an upstream service is properly propagated to downstream services.
//
// Notes:
//   - The method uses OpenTelemetry's global TextMapPropagator to extract the trace context.
//   - It expects the trace context to be present in the incoming HTTP request headers
//     in a format compatible with OpenTelemetry's distributed tracing standards (W3C Trace Context by default).
//   - After extraction, t.ctx is updated with the extracted context, and the service can
//     continue the trace from the point where the upstream service left off.
func (t *Trace) ExtractHTTPContext(req *http.Request) *Trace {
	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(t.ctx, propagation.HeaderCarrier(req.Header))
	t.ctx = ctx

	return t
}
