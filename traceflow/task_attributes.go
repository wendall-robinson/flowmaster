package traceflow

import (
	"time"

	"go.opentelemetry.io/otel/attribute"
)

// AddEvent adds an event name and timestamp as an attribute to the trace.
func (t *Trace) AddEvent(eventName string, timestamp time.Time) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("event.name", eventName),
		attribute.String("event.timestamp", timestamp.String()),
	)

	return t
}

// AddTaskInfo adds task-related information to the trace.
func (t *Trace) AddTaskInfo(taskID, taskName string, retries int) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("task.id", taskID),
		attribute.String("task.name", taskName),
		attribute.Int("task.retries", retries),
	)

	return t
}

// AddUser adds user-related attributes to the trace.
func (t *Trace) AddUser(userID, username string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("user.id", userID),
		attribute.String("user.username", username),
	)

	return t
}

// AddCustomMetric adds a custom metric to the trace.
func (t *Trace) AddCustomMetric(metricName string, value float64) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("metric.name", metricName),
		attribute.Float64("metric.value", value),
	)

	return t
}
