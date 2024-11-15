// Package natscarrier provides an implementation of the OpenTelemetry TextMapCarrier interface
// for NATS message headers. It enables injecting and extracting trace context into NATS
// messages, allowing for distributed tracing across services communicating via NATS.
//
// This package defines the NATSHeaderCarrier type, which wraps the nats.Header type
// to implement the TextMapCarrier interface required by OpenTelemetry propagators.
//
// # Usage
//
// ## Injecting trace context into NATS headers before publishing a message:
//
//	// Create NATS headers and carrier
//	headers := nats.Header{}
//	carrier := &natscarrier.NATSHeaderCarrier{Headers: headers}
//
//	// Propagate the trace context
//	traceflow.PropagateTraceContext(ctx, carrier)
//
//	// Publish the message with headers
//	err := natsConn.PublishMsg(&nats.Msg{
//	    Subject: "subject",
//	    Data:    messageData,
//	    Header:  headers,
//	})
//
// ## Extracting trace context from NATS headers upon receiving a message:
//
//	func handleMessage(msg *nats.Msg) {
//	    // Create carrier from NATS headers
//	    carrier := &natscarrier.NATSHeaderCarrier{Headers: msg.Header}
//
//	    // Extract the trace context
//	    ctx := traceflow.ExtractTraceContext(context.Background(), carrier)
//
//	    // Start a new span with the extracted context
//	    span := traceflow.New(ctx, "nats").Start("ConsumeEvent")
//	    defer span.End()
//
//	    // Process the message
//	    // ...
//	}
package natscarrier
