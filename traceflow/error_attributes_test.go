package traceflow

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// TestAddException tests that exception details are correctly added as attributes.
func TestAddException(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	err := errors.New("test exception")
	stackTrace := "example stack trace"

	trace.AddException(err, stackTrace)

	if len(trace.attrs) != 3 {
		t.Fatalf("Expected 3 attributes, got %d", len(trace.attrs))
	}

	expectedAttrs := []attribute.KeyValue{
		attribute.String("exception.type", "*errors.errorString"),
		attribute.String("exception.message", "test exception"),
		attribute.String("exception.stacktrace", stackTrace),
	}

	for i, expectedAttr := range expectedAttrs {
		if trace.attrs[i] != expectedAttr {
			t.Errorf("Expected attribute %v, got %v", expectedAttr, trace.attrs[i])
		}
	}
}

// TestAddError tests that error details are correctly added as attributes.
func TestAddError(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	err := errors.New("test error")

	trace.AddError(err)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}

	expectedAttrs := []attribute.KeyValue{
		attribute.String("error.type", "*errors.errorString"),
		attribute.String("error.message", "test error"),
	}

	for i, expectedAttr := range expectedAttrs {
		if trace.attrs[i] != expectedAttr {
			t.Errorf("Expected attribute %v, got %v", expectedAttr, trace.attrs[i])
		}
	}
}

// TestRecordError tests that an error is recorded to the span and that the span status is set to Error.
func TestRecordError(t *testing.T) {
	mockSpan := &mockSpan{}
	trace := &Trace{
		span: mockSpan,
	}

	err := errors.New("test error")
	trace.RecordError(err)

	if len(mockSpan.recordedErrors) != 1 {
		t.Fatalf("Expected 1 error to be recorded, got %d", len(mockSpan.recordedErrors))
	}

	if mockSpan.recordedErrors[0] != err {
		t.Errorf("Expected recorded error to be %v, got %v", err, mockSpan.recordedErrors[0])
	}

	if mockSpan.statusCode != codes.Error {
		t.Errorf("Expected span status to be Error, got %v", mockSpan.statusCode)
	}

	if mockSpan.statusDescription != "test error" {
		t.Errorf("Expected status description to be 'test error', got %v", mockSpan.statusDescription)
	}
}

// Mock Span for testing RecordError
type mockSpan struct {
	trace.Span
	recordedErrors    []error
	statusCode        codes.Code
	statusDescription string
}

// RecordError mocks the RecordError method of the Span interface
func (m *mockSpan) RecordError(err error, _ ...trace.EventOption) {
	m.recordedErrors = append(m.recordedErrors, err)
}

// SetStatus mocks the SetStatus method of the Span interface
func (m *mockSpan) SetStatus(code codes.Code, description string) {
	m.statusCode = code
	m.statusDescription = description
}
