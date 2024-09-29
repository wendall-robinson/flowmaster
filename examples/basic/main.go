package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	tf "github.com/wendall-robinson/traceflow"
)

func main() {
	// Set up a context for tracing
	ctx := context.Background()

	// Initialize OpenTelemetry with tracing
	ctx, shutdown, err := tf.Init(ctx, "example-app")
	if err != nil {
		log.Fatalf("Failed to initialize OpenTelemetry: %v", err)
	}

	defer shutdown(ctx)

	// Start a trace
	trace := tf.New(ctx, "example-operation")
	defer trace.Start("example-span").End()

	// Add a sample attribute
	trace.AddAttribute(
		tf.AddString("example-key", "example-value"),
	)

	fmt.Println("Trace started, and the app is holding for further traces...")

	// Keep the application running (simulates a server or long-running process)
	select {
	case <-time.After(10 * time.Minute):
		fmt.Println("Exiting after 10 minutes")
	case <-interrupt():
		fmt.Println("App interrupted, shutting down")
	}
}

// interrupt listens for termination signals (e.g., CTRL+C)
func interrupt() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	// Listen for interrupt signal
	return ch
}
