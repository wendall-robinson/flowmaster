package traceflow

import (
	"context"
	"testing"
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// TestAddEvent tests the AddEvent method.
func TestAddEvent(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")

	// Use a fixed time for the test
	eventTime := time.Date(2024, time.January, 1, 10, 0, 0, 0, time.UTC)
	trace.AddEvent("test-event", eventTime)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("event.name", "test-event") {
		t.Errorf("Expected event.name to be 'test-event', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("event.timestamp", eventTime.String()) {
		t.Errorf("Expected event.timestamp to be '%s', got %v", eventTime.String(), trace.attrs[1])
	}
}

// TestAddTaskInfo tests the AddTaskInfo method.
func TestAddTaskInfo(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddTaskInfo("task-123", "example-task", 3)

	if len(trace.attrs) != 3 {
		t.Fatalf("Expected 3 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("task.id", "task-123") {
		t.Errorf("Expected task.id to be 'task-123', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("task.name", "example-task") {
		t.Errorf("Expected task.name to be 'example-task', got %v", trace.attrs[1])
	}
	if trace.attrs[2] != attribute.Int("task.retries", 3) {
		t.Errorf("Expected task.retries to be 3, got %v", trace.attrs[2])
	}
}

// TestAddUser tests the AddUser method.
func TestAddUser(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddUser("user-456", "testuser")

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("user.id", "user-456") {
		t.Errorf("Expected user.id to be 'user-456', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.String("user.username", "testuser") {
		t.Errorf("Expected user.username to be 'testuser', got %v", trace.attrs[1])
	}
}

// TestAddCustomMetric tests the AddCustomMetric method.
func TestAddCustomMetric(t *testing.T) {
	ctx := context.TODO()
	trace := New(ctx, "test-service")
	trace.AddCustomMetric("response_time", 1.23)

	if len(trace.attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(trace.attrs))
	}
	if trace.attrs[0] != attribute.String("metric.name", "response_time") {
		t.Errorf("Expected metric.name to be 'response_time', got %v", trace.attrs[0])
	}
	if trace.attrs[1] != attribute.Float64("metric.value", 1.23) {
		t.Errorf("Expected metric.value to be 1.23, got %v", trace.attrs[1])
	}
}
