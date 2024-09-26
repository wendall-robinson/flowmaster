package traceflow

import (
	"context"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

// TestAddHTTPRequest tests that HTTP request details are correctly added as attributes.
func TestAddHTTPRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	req.Header.Set("User-Agent", "Go-http-client/1.1")
	req.RemoteAddr = "127.0.0.1"

	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddHTTPRequest(req)

	if len(trace.attrs) != 4 {
		t.Fatalf("Expected 4 attributes, got %d", len(trace.attrs))
	}

	expectedAttrs := []attribute.KeyValue{
		attribute.String("http.method", "GET"),
		attribute.String("http.url", "http://example.com"),
		attribute.String("http.user_agent", "Go-http-client/1.1"),
		attribute.String("http.client_ip", "127.0.0.1"),
	}

	for i, expectedAttr := range expectedAttrs {
		if trace.attrs[i] != expectedAttr {
			t.Errorf("Expected attribute %v, got %v", expectedAttr, trace.attrs[i])
		}
	}
}

// TestAddHTTPResponse tests that HTTP response details are correctly added as attributes.
func TestAddHTTPResponse(t *testing.T) {
	statusCode := 200
	contentLength := int64(1234)

	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddHTTPResponse(statusCode, contentLength)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}

	if trace.attrs[0] != attribute.Int("http.status_code", 200) {
		t.Errorf("Expected http.status_code to be 200, got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.Int64("http.content_length", 1234) {
		t.Errorf("Expected http.content_length to be 1234, got %v", trace.attrs[1])
	}
}

// TestAddHTTPHeaders tests that HTTP headers are correctly added as attributes.
func TestAddHTTPHeaders(t *testing.T) {
	headers := http.Header{
		"X-Custom-Header":   []string{"CustomValue"},
		"X-Multiple-Header": []string{"Value1", "Value2"},
	}

	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddHTTPHeaders(headers)

	if len(trace.attrs) != 3 { // Two headers, but one has multiple values
		t.Fatalf("Expected 3 attributes, got %d", len(trace.attrs))
	}

	expectedAttrs := []attribute.KeyValue{
		attribute.String("http.header.X-Custom-Header", "CustomValue"),
		attribute.String("http.header.X-Multiple-Header", "Value1"),
		attribute.String("http.header.X-Multiple-Header", "Value2"),
	}

	for i, expectedAttr := range expectedAttrs {
		if trace.attrs[i] != expectedAttr {
			t.Errorf("Expected attribute %v, got %v", expectedAttr, trace.attrs[i])
		}
	}
}
