package main

import (
	"context"
	"io"

	"net/http"
	"time"

	log "github.com/wendall-robinson/flowmaster/logflow"
	"github.com/wendall-robinson/flowmaster/traceflow"
)

type App struct {
	ctx    context.Context
	logger log.Logger
}

func main() {
	// Initialize the OpenTelemetry tracing system
	ctx := context.Background()

	// Init the logger
	logger := log.NewLogger()

	app := App{
		ctx:    ctx,
		logger: *logger,
	}

	// Initialize OpenTelemetry with OTLP export to a collector running on localhost:4317
	ctx, shutdown, err := traceflow.Init(ctx, "basic-service", traceflow.WithOLTP("otel:4317"))
	if err != nil {
		// log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
		logger.Info("Failed to initialize OpenTelemetry: %v")
	}

	defer shutdown(ctx)

	// fmt.Println("Application is running...")
	app.logger.Info("Application is running...")

	// app.logger.InfoSys("System initialized")

	// Create a new traceflow instance with the context and service name, we add an optional system info attribute here
	trace := traceflow.New(app.ctx, "basic-service", traceflow.WithSystemInfo())

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
	app.sendHTTPRequest()

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

func (a *App) sendHTTPRequest() {
	// The URL of the endpoint
	url := "http://web:8080/test"

	// Create a new GET request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// log.Fatalf("Failed to create request: %s", err)
		a.logger.Info("Failed to create request")
	}

	// Create a new traceflow instance with the context and service name
	trace := traceflow.New(a.ctx, "http-client").Client()

	// Start the span and inject the trace context into the request
	// be sure to inject the context after starting the span
	defer trace.Start("http-request").InjectHTTPContext(req).End()

	// Perform the request using http.DefaultClient
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		a.logger.Info("Failed to read response body")
	}

	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.logger.Info("Failed to read response body")
	}

	a.logger.InfoHttp("HTTP Request", req, resp.StatusCode, "response", string(body))
}
