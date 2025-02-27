package traceflow

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/segmentio/kafka-go"
	"github.com/streadway/amqp"
	"github.com/wendall-robinson/flowmaster/traceflow/internal/kafkacarrier"
	"github.com/wendall-robinson/flowmaster/traceflow/internal/rabbitmqcarrier"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// PropagateNats injects the trace context into NATS headers.
func PropagateNats(ctx context.Context, headers nats.Header) {
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(headers)
	propagator.Inject(ctx, carrier)
}

// PropagateKafka injects the trace context into Kafka headers.
func PropagateKafka(ctx context.Context, headers *[]kafka.Header) {
	propagator := otel.GetTextMapPropagator()
	carrier := kafkacarrier.New(headers)
	propagator.Inject(ctx, carrier)
}

// PropagateRabbitMQ injects the trace context into RabbitMQ headers.
func PropagateRabbitMQ(ctx context.Context, headers amqp.Table) {
	propagator := otel.GetTextMapPropagator()
	carrier := rabbitmqcarrier.New(headers)
	propagator.Inject(ctx, carrier)
}

// ExtractNats extracts the trace context from NATS headers.
func ExtractNats(ctx context.Context, headers nats.Header) context.Context {
	propagator := otel.GetTextMapPropagator()
	carrier := propagation.HeaderCarrier(headers)

	return propagator.Extract(ctx, carrier)
}

// ExtractKafka extracts the trace context from Kafka headers.
func ExtractKafka(ctx context.Context, headers []kafka.Header) context.Context {
	propagator := otel.GetTextMapPropagator()
	carrier := kafkacarrier.New(&headers)

	return propagator.Extract(ctx, carrier)
}

// ExtractRabbitMQ extracts the trace context from RabbitMQ headers.
func ExtractRabbitMQ(ctx context.Context, headers amqp.Table) context.Context {
	propagator := otel.GetTextMapPropagator()
	carrier := rabbitmqcarrier.New(headers)

	return propagator.Extract(ctx, carrier)
}
