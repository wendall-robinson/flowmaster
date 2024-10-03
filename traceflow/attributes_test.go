package traceflow

import (
	"context"
	"encoding/json"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// TestStringAttr tests the creation of a string OTEL attribute.
func TestStringAttr(t *testing.T) {
	attr := AddString("key", "value")
	expected := attribute.String("key", "value")

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestStringSliceAttr tests the creation of a string slice OTEL attribute.
func TestStringSliceAttr(t *testing.T) {
	attr := AddStringSlice("key", []string{"value1", "value2"})
	expected := attribute.StringSlice("key", []string{"value1", "value2"})

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestIntAttr tests the creation of an int OTEL attribute.
func TestIntAttr(t *testing.T) {
	attr := AddInt("key", 123)
	expected := attribute.Int("key", 123)

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestIntSliceAttr tests the creation of an int slice OTEL attribute.
func TestIntSliceAttr(t *testing.T) {
	attr := AddIntSlice("key", []int{1, 2, 3})
	expected := attribute.IntSlice("key", []int{1, 2, 3})

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestFloatAttr tests the creation of a float64 OTEL attribute.
func TestFloatAttr(t *testing.T) {
	attr := AddFloat("key", 1.23)
	expected := attribute.Float64("key", 1.23)

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestFloatSliceAttr tests the creation of a float64 slice OTEL attribute.
func TestFloatSliceAttr(t *testing.T) {
	attr := AddFloatSlice("key", []float64{1.1, 2.2, 3.3})
	expected := attribute.Float64Slice("key", []float64{1.1, 2.2, 3.3})

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestBoolAttr tests the creation of a bool OTEL attribute.
func TestBoolAttr(t *testing.T) {
	attr := AddBool("key", true)
	expected := attribute.Bool("key", true)

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestBoolSliceAttr tests the creation of a bool slice OTEL attribute.
func TestBoolSliceAttr(t *testing.T) {
	attr := AddBoolSlice("key", []bool{true, false, true})
	expected := attribute.BoolSlice("key", []bool{true, false, true})

	if attr.otelAttr != expected {
		t.Errorf("Expected attribute %v, got %v", expected, attr.otelAttr)
	}
}

// TestAddJSON tests that the AddJSON method correctly adds a JSON payload as a string attribute.
func TestAddJSON(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")

	jsonPayload := json.RawMessage(`{"key":"value"}`)
	trace.AddJSON(jsonPayload)

	if len(trace.attrs) != 1 {
		t.Fatalf("Expected 1 attribute, got %d", len(trace.attrs))
	}

	expectedAttr := attribute.String("payload", `{"key":"value"}`)
	if trace.attrs[0] != expectedAttr {
		t.Errorf("Expected attribute %v, got %v", expectedAttr, trace.attrs[0])
	}
}
