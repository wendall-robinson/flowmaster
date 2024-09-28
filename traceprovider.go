package traceflow

import (
	"context"
	"io"
	"log"
	"os"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// InitOption defines a functional option for customizing the Init process
type InitOption func(*TelemetryBuilder)

// noopSpanExporter is a custom no-op exporter to prevent traces from being exported.
type noopSpanExporter struct{}

func (n *noopSpanExporter) ExportSpans(context.Context, []sdktrace.ReadOnlySpan) error {
	// Do nothing
	return nil
}

func (n *noopSpanExporter) Shutdown(context.Context) error {
	// Do nothing
	return nil
}

// TelemetryBuilder holds configuration for OTEL setup
type TelemetryBuilder struct {
	ctx      context.Context
	logger   *log.Logger
	exporter sdktrace.SpanExporter
	filePath string
}

// Init initializes OpenTelemetry with a custom context, service name, and returns
// the initialized context and a shutdown function for cleanup.
func Init(ctx context.Context, serviceName string, opts ...InitOption) (context.Context, func(context.Context), error) {
	// Create a default TelemetryBuilder
	builder := &TelemetryBuilder{
		ctx:    ctx,
		logger: log.New(os.Stdout, "", log.LstdFlags), // Default logger writes to stdout
	}

	// Apply the options
	for _, opt := range opts {
		opt(builder)
	}

	// Set up the exporter based on configuration
	if builder.exporter == nil {
		// Default to stdout exporter if none is specified
		builder.exporter, _ = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}

	// Create SpanProcessor using the exporter
	spanProcessor := sdktrace.NewBatchSpanProcessor(builder.exporter)

	// Create and set the TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(spanProcessor),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Set the Global TracerProvider
	otel.SetTracerProvider(tp)

	// Set global propagator for context propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Shutdown function for cleanup
	shutdown := func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			builder.logger.Printf("Error shutting down tracer provider: %v", err)
		}
	}

	// Log success and return
	builder.logger.Println("OpenTelemetry initialized successfully")

	return builder.ctx, shutdown, nil
}

// WithOLTP sets the OLTP exporter to send traces to an OpenTelemetry collector.
func WithOLTP() InitOption {
	return func(tb *TelemetryBuilder) {
		tb.logger.Println("Using OTLP exporter")

		exp, err := otlptracegrpc.New(tb.ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint("localhost:4317"))
		if err != nil {
			tb.logger.Printf("Failed to create OLTP exporter: %v", err)

			return
		}

		tb.exporter = exp
	}
}

// WithFileLogging sets up trace exporting to a file
func WithFileLogging(filePath string) InitOption {
	return func(tb *TelemetryBuilder) {
		const filemode = 0o644

		tb.filePath = filePath

		// Open or create the file for appending trace logs
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filemode)
		if err != nil {
			tb.logger.Printf("Error creating file exporter: %v", err)

			return
		}

		// Create the exporter to write to the file
		exporter, err := stdouttrace.New(stdouttrace.WithWriter(file))
		if err != nil {
			tb.logger.Printf("Error creating file exporter: %v", err)
			return
		}

		tb.exporter = exporter
	}
}

// WithSilentLogger sets a no-op logger and no-op span exporter, useful for testing.
func WithSilentLogger() InitOption {
	return func(tb *TelemetryBuilder) {
		tb.logger = log.New(io.Discard, "", 0) // Silence the logger
		// Use the custom no-op span exporter to suppress trace output
		tb.exporter = &noopSpanExporter{}
	}
}

// WithLogger allows users to provide a custom logger.
func WithLogger(logger *log.Logger) InitOption {
	return func(tb *TelemetryBuilder) {
		tb.logger = logger
	}
}
