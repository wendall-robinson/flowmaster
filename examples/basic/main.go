package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/wendall-robinson/traceflow"
)

func main() {
	// Initialize the OpenTelemetry tracing system
	ctx := context.Background()

	// Initialize OpenTelemetry with OTLP export to a collector running on localhost:4317
	ctx, shutdown, err := traceflow.Init(ctx, "basic-service", traceflow.WithOLTP("otel:4317"))
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	fmt.Println("Application is running...")

	// Create a new traceflow instance with the context and service name, we add an optional system info attribute here
	trace := traceflow.New(ctx, "basic-service", traceflow.WithSystemInfo())

	// add additional attributes to the trace, you can also add these when calling .New()
	trace.AddAttribute(
		traceflow.AddString("example-attribute", "example-value"),
	)

	// start the span, don't forget to call .End() when the span is complete
	trace.Start("main-span")

	// when calling a downstream function or service that you want to trace,
	// pass in the context from the current trace to propagate the trace context
	for _, test := range []bool{true, false} {
		traceIf(trace.GetContext(), test)
	}

	// Test the HTTP request span, be sure to pass in the trace context to propagate the trace
	sendHTTPRequest(trace.GetContext())

	// End the main span when all other spans are complete
	trace.End()

	// Prevent the app from exiting immediately
	for {
		time.Sleep(30 * time.Minute)
	}
}

func traceIf(ctx context.Context, test bool) {
	trace := traceflow.New(ctx, "basic-service").AddAttributeIf(test, "test", "true")

	defer trace.Start("AddAttributeIf-Span").End()
}

func sendHTTPRequest(ctx context.Context) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("User-Agent", "Go-http-client/1.1")
	req.RemoteAddr = "127.0.0.1"

	trace := traceflow.New(ctx, "mock-http-client").Client()
	trace.AddHTTPRequest(req)

	defer trace.Start("http-request").End()

	// once the trace has started, inject the trace context into the HTTP request
	trace.InjectHTTPContext(req)

	_ = receiveHTTPResponse(req)
}

func receiveHTTPResponse(req *http.Request) bool {
	ctx := context.Background()

	trace := traceflow.New(ctx, "mock-remote-http-service").Server()

	// extract the trace context from the HTTP request
	trace.ExtractHTTPContext(req).AddHTTPResponse(200, 1234)

	defer trace.Start("http-response").End()

	return true
}
