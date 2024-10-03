package traceflow

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestInit(t *testing.T) {
	// Initialize OpenTelemetry with default settings
	ctx := context.Background()
	serviceName := "test-service"

	// Call the Init function
	newCtx, shutdown, err := Init(ctx, serviceName, WithSilentLogger())

	// Ensure there are no errors
	assert.NoError(t, err, "Expected no error during Init()")

	// Ensure that a valid context is returned
	assert.NotNil(t, newCtx, "Expected context to be initialized")

	// Ensure that the shutdown function is not nil
	assert.NotNil(t, shutdown, "Expected a valid shutdown function")

	// Check that a TracerProvider has been set and is not the default (no-op)
	tp := otel.GetTracerProvider()
	assert.NotEqual(t, noop.NewTracerProvider(), tp, "Expected a valid TracerProvider to be set")

	// Ensure the shutdown function runs without errors
	assert.NotPanics(t, func() {
		shutdown(newCtx)
	}, "Expected shutdown to execute without panicking")
}

func TestInitWithSilentLogger(t *testing.T) {
	// Create a buffer to capture any log output
	var logOutput bytes.Buffer

	// Initialize OpenTelemetry with silent logger
	ctx := context.Background()
	serviceName := "test-service"

	// Pass in the WithSilentLogger option to suppress logging
	newCtx, shutdown, err := Init(ctx, serviceName, WithSilentLogger())

	// Ensure there are no errors
	assert.NoError(t, err, "Expected no error during Init()")

	// Ensure that a valid context is returned
	assert.NotNil(t, newCtx, "Expected context to be initialized")

	// Ensure that the shutdown function is not nil
	assert.NotNil(t, shutdown, "Expected a valid shutdown function")

	// Check that the logger produced no output
	assert.Equal(t, "", logOutput.String(), "Expected no log output when using silent logger")

	// Check that a valid TracerProvider is set (not the noop provider)
	tp := otel.GetTracerProvider()
	assert.NotEqual(t, noop.NewTracerProvider(), tp, "Expected a valid TracerProvider to be set")

	// Ensure the shutdown function runs without errors
	assert.NotPanics(t, func() {
		shutdown(newCtx)
	}, "Expected shutdown to execute without panicking")
}
