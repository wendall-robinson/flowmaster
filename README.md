# TraceFlow - Simplified OpenTelemetry Tracing for Go

`TraceFlow` is a Go package that provides a simple, fluent interface for integrating OpenTelemetry distributed tracing into your Go applications. It abstracts away much of the complexity and boilerplate involved in creating and managing traces, allowing developers to focus on their application logic while still benefiting from powerful distributed tracing features.

## Why Use TraceFlow?

In modern, distributed systems, it's crucial to have observability tools like tracing to monitor and debug complex interactions between microservices, databases, and other components. **OpenTelemetry** is a robust solution for capturing this telemetry data, but its complexity can make initial adoption challenging. This is where TraceFlow comes in.

## Key Features and Benefits:

* **Fluent Interface**: Chain methods together to add attributes, links, and status information to your trace with minimal code.
  * Example: `defer trace.Start("operation").AddAttribute(...).End()`
* **Automatic System Information**: Automatically capture and add CPU, memory, and disk information to your traces without manual setup.
* **Built-in Best Practices**: Enforces best practices such as always ending spans and recording errors and success states, helping you avoid common pitfalls.
* **Extendable**: Easily extend the functionality with your own custom attributes or predefined behavior using variadic options.
* **Error and Exception Handling**: Capture detailed error and exception information, including stack traces and error messages.
* **Context Propagation**: Simplifies passing and extracting trace context across service boundaries (e.g., HTTP requests).

## Who Should Use This Package?

* **Go Developers** working in microservice-based or distributed systems who want to integrate OpenTelemetry tracing without writing verbose boilerplate code.
* Teams and companies seeking enhanced observability to monitor their distributed systems in production environments.
* Developers who want to **add tracing to existing services with minimal code changes and minimal learning curve.**
* Anyone who wants a lightweight, flexible tracing solution that can scale with their system and evolve as new OpenTelemetry features are introduced.

## Installation
To install TraceFlow, use `go get`:

```bash
go get github.com/wendall-robinson/traceflow
```

## Quick Start
**Basic Usage Example**

Here's how you can start tracing operations in your Go application using `TraceFlow`:
```go
package main

import (
    "context"
    "log"

    "github.com/wendall-robinson/traceflow"
)

func main() {
    ctx := context.Background()

    // Create a new trace
    trace := traceflow.New(ctx, "example-service")

    // Start a new operation and add attributes
    defer trace.Start("main-operation").End()

    // Add custom attributes
    trace.AddAttribute(
        traceflow.AddString("user_id", "12345"),
        traceflow.AddInt("response_time", 200),
    )

    // Simulate an error and record it
    err := someOperation()
    if err != nil {
        trace.RecordFailure(err, "failed to complete operation")
    }
}
```
### Key Features Demonstrated:

* **Creating and Starting a Trace:** traceflow.New and trace.Start.
* **Adding Attributes:** trace.AddAttribute to capture custom key-value pairs.
* **Error Handling:** Using trace.RecordFailure to record errors and failures within the trace.

## Default Context Propagation

By default, when you create a new trace using the `New()` method, the trace context is automatically copied from the provided context (if it exists). This ensures that the trace is linked to any existing parent trace, making it easier to maintain the trace chain across distributed services.

**Example usage:**
```go
// Incoming request with an existing trace context
ctx := r.Context() // from an HTTP request

// Create a new trace, linked to the parent trace (if it exists)
trace := traceflow.New(ctx, "my-service")
trace.Start("operation").End()
```

## Starting a Fresh Trace Without Context Propagation

If you want to start a fresh trace and not propagate the existing trace context, use the NewWithoutPropagation() method. This allows you to create an independent trace.

**Example usage:**
```go
// Start a new trace without copying the existing trace context
trace := traceflow.NewWithoutPropagation(ctx, "my-service")
trace.Start("operation").End()
```
## Advanced Features
### Advanced Features: System Information
 * **Adding System Information:** Automatically add CPU, memory, and disk usage to your traces:
    ```go
    trace.AddCpuInfo().AddMemoryInfo().AddDiskInfo()
    ```
* **Add all System Information:** (adds CPU, Memory and Disk info)
    ```go
    trace.WithSystemInfo()
    ```
### Advanced Features:  Injecting Trace Context into HTTP Requests
In distributed systems, it's important to propagate the trace context across service boundaries, allowing each service to continue a trace. This is particularly useful in microservice architectures, where HTTP requests are often used to communicate between services.

With the InjectHTTPContext method, you can easily inject the trace context into the headers of an outgoing HTTP request, ensuring that the trace information is passed along to downstream services.

#### InjectHttpContext
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"log"
	"time"

	"github.com/wendall-robinson/traceflow"
)

func main() {
	// Create a new trace
	ctx := context.Background()
	trace := traceflow.New(ctx, "example-service")

	// Start a new span for the operation
	defer trace.Start("http-outgoing-request").End()

	// Create an HTTP request to be sent to another service
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Inject the trace context into the request headers
	trace.InjectHTTPContext(req)

	// Send the request using an HTTP client
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}

	defer resp.Body.Close()

	fmt.Printf("Response status: %s\n", resp.Status)
}
```
**How It Works:**

* A new trace is started with `traceflow.New()`.
* The `InjectHTTPContext` method automatically injects the trace context into the request headers, making it possible for the receiving service to extract and continue the trace.
* The downstream service can use the corresponding `ExtractHTTPContext` method to continue the trace.

This pattern ensures seamless propagation of trace context between services, providing end-to-end visibility in distributed tracing.

#### ExtractHTTPContext
With the `ExtractHTTPContext` method, you can easily extract the trace context from incoming request headers and continue the trace.
```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"log"
	"github.com/wendall-robinson/traceflow"
)

func main() {
	// Start an HTTP server to receive requests
	http.HandleFunc("/process", func(w http.ResponseWriter, req *http.Request) {
		// Create a new trace, using the incoming request's context
		ctx := context.Background()
		trace := traceflow.New(ctx, "example-service")

		// Extract the trace context from the incoming request headers
		trace.ExtractHTTPContext(req)

		// Start a new span for this operation
		defer trace.Start("process-request").End()

		// Simulate some processing
		fmt.Fprintln(w, "Processing request...")

		log.Println("Request processed, trace continued.")
	})

	// Start the HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```
**How It Works:**

* A new trace is initialized with `traceflow.New()`.
* The `ExtractHTTPContext` method extracts the trace context from the incoming HTTP request's headers. This allows the trace context to be continued from where the upstream service left off.
* A new span is started for the current operation (process-request) and ended after the operation completes.

**Usage:**

When the client service sends a request (such as in the earlier example with InjectHTTPContext), the trace context is passed along in the request headers. The receiving service extracts the context with ExtractHTTPContext and can continue the trace, creating a new span for the current operation.
