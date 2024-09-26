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
