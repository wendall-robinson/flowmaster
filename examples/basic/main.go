package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/wendall-robinson/flowmaster/traceflow"
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
	// The URL of the endpoint
	url := "http://web:8080/test"

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %s", err)
	}

	// Create a new traceflow instance with the context and service name
	trace := traceflow.New(ctx, "http-client").Client()

	// Start the span and inject the trace context into the request
	// be sure to inject the context after starting the span
	defer trace.Start("http-request").InjectHTTPContext(req).End()

	// Perform the request using http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request: %s", err)
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %s", err)
	}

	// Print the response
	fmt.Println(string(body))
}
