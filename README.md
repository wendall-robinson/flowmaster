# traceflow - Simplified OpenTelemetry Tracing

`traceflow` is a Go package that provides a simplified and fluent interface for integrating OpenTelemetry tracing into Go applications. It abstracts some of the common tasks associated with managing traces, allowing developers to focus on their application logic.

## Features

- **Simple Trace and Span Management**: Easily start and manage traces with minimal boilerplate.
- **Context Propagation**: Tools to inject and extract contexts for HTTP requests to support distributed tracing.
- **Error Handling**: Built-in mechanisms to record errors and set span statuses effectively.
- **Link Spans**: Facilitate linking spans across different traces or services.

## Installation

To install `traceflow`, use the following `go get` command:

```bash
go get -u github.com/wendall-robinson/traceflow
```

## Usage
Below are some examples demonstrating how to use traceflow:

### Creating a New Trace
Start by creating a new trace with a specific service name:

```golang
package main

import (
    "context"
    "log"
    "fmt"

    "github.com/wendall-robinson/traceflow"
)

func main() {
    var (
        mathIsHard = fmt.Errorf("math is hard")
        two = "two"
        three = 3
    )

    ctx := context.Background()

    trace := traceflow.New(ctx, "example-service").
        AddAttributes(
		    attribute.String("param1", two),
		    attribute.Int("param2", three),
        )

    defer traceflow.Start("main-operation").End()


    if two != three {
        trace.RecordFailure(mathIsHard, "error comparing parameters")
    }
}
```

### Adding Attributes and Links
You can add attributes and links to spans easily:

```golang
func operation(ctx context.Context) {
    trace := traceflow.NewTrace(ctx, "example-service")

    trace.AddKeyValue("user_id", "12345")
    trace.AddLink(otherSpanContext) // Assuming otherSpanContext is available

    defer trace.Start("operation").End()
}
```
### Handling HTTP Requests
traceflow can also inject and extract traces from HTTP requests, aiding in distributed tracing across microservices:

```golang
func httpHandler(w http.ResponseWriter, r *http.Request) {
    trace := traceflow.NewTrace(r.Context(), "web-service")
    trace.ExtractHTTPContext(r)

    defer trace.Start("http-request").End()
    // Process the request
}
```
