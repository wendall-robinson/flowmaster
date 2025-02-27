// Package rabbitmqcarrier provides an implementation of the OpenTelemetry TextMapCarrier interface
// for RabbitMQ message headers. It enables injecting and extracting trace context into RabbitMQ
// message headers (amqp.Table), allowing for distributed tracing over RabbitMQ messaging systems.
//
// This package defines the RabbitMQHeaderCarrier type, which wraps the amqp.Table type
// to implement the TextMapCarrier interface required by OpenTelemetry propagators.
//
// # Usage
//
// ## Injecting trace context into RabbitMQ headers before publishing a message:
//
//	// Initialize RabbitMQ headers and carrier
//	headers := amqp.Table{}
//	carrier := &rabbitmqcarrier.RabbitMQHeaderCarrier{Headers: headers}
//
//	// Propagate the trace context
//	traceflow.PropagateTraceContext(ctx, carrier)
//
//	// Publish the message with headers
//	err := channel.Publish(
//	    exchangeName,
//	    routingKey,
//	    false,
//	    false,
//	    amqp.Publishing{
//	        Headers:     headers,
//	        Body:        messageData,
//	        ContentType: "application/json",
//	    },
//	)
//
// ## Extracting trace context from RabbitMQ headers upon receiving a message:
//
//	func handleMessage(msg amqp.Delivery) {
//	    // Create carrier from RabbitMQ headers
//	    carrier := &rabbitmqcarrier.RabbitMQHeaderCarrier{Headers: msg.Headers}
//
//	    // Extract the trace context
//	    ctx := traceflow.ExtractTraceContext(context.Background(), carrier)
//
//	    // Start a new span with the extracted context
//	    span := traceflow.New(ctx, "rabbitmq").Start("ConsumeEvent")
//	    defer span.End()
//
//	    // Process the message
//	    // ...
//	}
package rabbitmqcarrier
