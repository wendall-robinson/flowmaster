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
// This ensures that the context of a trace is propagated across service boundaries in
// distributed systems.
//
// The trace context is injected using the global propagator, which handles the
// serialization of the trace context as HTTP headers. Users of traceflow do not need to
// import or manage OpenTelemetry propagators directly.
//
// Example usage:
//
//	req, _ := http.NewRequest("GET", "http://example.com", nil)
//	trace.InjectHTTPContext(req)
//
// Notes:
//   - The trace context is injected into the HTTP request's headers using the default
//     W3C Trace Context format.
func (t *Trace) InjectHTTPContext(req *http.Request) *Trace {
	if t.ctx == nil {
		return t
	}

	// Use the internal OpenTelemetry propagator to inject the context
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(req.Header)
	propagator.Inject(t.ctx, carrier)

	return t
}

// ExtractHTTPContext extracts the trace context from the headers of an HTTP request
// and updates the Trace's context. This ensures that the current service can join
// an existing trace initiated by an upstream service.
//
// The trace context is extracted using the global propagator. Users of traceflow do
// not need to interact with OpenTelemetry’s propagators directly.
//
// Example usage:
//
//	trace := traceflow.New(ctx, "my-service")
//	trace.ExtractHTTPContext(req)
//
// Notes:
// - This method updates the Trace's context (t.ctx) with the extracted trace context.
func (t *Trace) ExtractHTTPContext(req *http.Request) *Trace {
	// Use the internal OpenTelemetry propagator to extract the context
	propagator := otel.GetTextMapPropagator()
	ctx := propagator.Extract(t.ctx, propagation.HeaderCarrier(req.Header))
	t.ctx = ctx

	return t
}
