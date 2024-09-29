package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wendall-robinson/traceflow"
)

func main() {
	// Initialize the OpenTelemetry tracing system
	ctx := context.Background()

	// Initialize OpenTelemetry
	ctx, shutdown, err := traceflow.Init(ctx, "basic-service", traceflow.WithOLTP("otel:4317"))
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	fmt.Println("Application is running...")

	trace := traceflow.New(ctx, "basic-service", traceflow.WithSystemInfo()).AddAttribute(
		traceflow.AddString("example-attribute", "example-value"),
	)

	defer trace.Start("test-span").End()
	fmt.Printf("Trace: %+v\n", trace)
	trace.End()

	// Prevent the app from exiting immediately
	for {
		time.Sleep(30 * time.Minute)
	}
}
