package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/wendall-robinson/flowmaster/traceflow"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	trace := traceflow.New(r.Context(), "http-handler").Server()
	trace.ExtractHTTPContext(r)

	defer trace.Start("testing-endpoint").End()

	fmt.Fprintf(w, "hello world!")

}

func main() {
	// Initialize the OpenTelemetry tracing system
	ctx := context.Background()

	// Initialize OpenTelemetry with OTLP export to a collector running on localhost:4317
	ctx, shutdown, err := traceflow.Init(ctx, "web-service", traceflow.WithOLTP("otel:4317"))
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	fmt.Println("Web Service is running...")

	// Set up the /test endpoint
	http.HandleFunc("/test", testHandler)

	// Run the server on port 8080
	fmt.Println("Server running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
