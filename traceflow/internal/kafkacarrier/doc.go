// Package kafkacarrier provides an implementation of the OpenTelemetry TextMapCarrier interface
// for Kafka message headers. It facilitates injecting and extracting trace context into Kafka
// message headers, enabling distributed tracing over Kafka messaging systems.
//
// This package defines the KafkaHeaderCarrier type, which adapts the []kafka.Header type
// to implement the TextMapCarrier interface used by OpenTelemetry propagators.
//
// # Usage
//
// ## Injecting trace context into Kafka headers before publishing a message:
//
//	// Initialize Kafka headers and carrier
//	var headers []kafka.Header
//	carrier := &kafkacarrier.KafkaHeaderCarrier{Headers: headers}
//
//	// Propagate the trace context
//	traceflow.PropagateTraceContext(ctx, carrier)
//
//	// Create and publish the Kafka message with headers
//	message := kafka.Message{
//	    Topic:   "topic",
//	    Value:   messageData,
//	    Headers: carrier.Headers(),
//	}
//	err := kafkaWriter.WriteMessages(ctx, message)
//
// ## Extracting trace context from Kafka headers upon receiving a message:
//
//	func handleMessage(msg kafka.Message) {
//	    // Create carrier from Kafka headers
//	    carrier := &kafkacarrier.KafkaHeaderCarrier{Headers: msg.Headers}
//
//	    // Extract the trace context
//	    ctx := traceflow.ExtractTraceContext(context.Background(), carrier)
//
//	    // Start a new span with the extracted context
//	    span := traceflow.New(ctx, "kafka").Start("ConsumeEvent")
//	    defer span.End()
//
//	    // Process the message
//	    // ...
//	}
package kafkacarrier
