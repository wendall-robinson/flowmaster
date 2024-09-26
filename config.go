package traceflow

import (
	"os"
	"runtime"

	"go.opentelemetry.io/otel/attribute"
)

// Option defines a function signature for modifying the Trace object
type Option func(*Trace)

// WithSystemInfo adds system-related attributes: CPU, Memory, Disk
func WithSystemInfo() Option {
	return func(t *Trace) {
		t.AddCPUInfo()
		t.AddMemoryInfo()
		t.AddDiskInfo()
	}
}

// WithAttributes allows users to provide custom attributes to be added
// during the Trace object initialization.
func WithAttributes(attrs ...attribute.KeyValue) Option {
	return func(t *Trace) {
		t.attrs = append(t.attrs, attrs...)
	}
}

// WithConcurrencyInfo adds the number of goroutines to the Trace attributes
func WithConcurrencyInfo() Option {
	return func(t *Trace) {
		t.attrs = append(t.attrs,
			attribute.Int("goroutine.count", runtime.NumGoroutine()),
		)
	}
}

// WithEnVars retrieves environment variables specified in the keys slice and adds them
// as attributes to the Trace. If an environment variable is not set or is empty, a warning is logged,
// and that attribute is not added.
//
// Example usage:
//
//	envKeys := []string{"SERVICE_NAME", "DEPLOYMENT_ENV", "CLOUD_REGION", "SERVICE_VERSION"}
//	trace := traceflow.New(ctx, "my-service", traceflow.WithEnVars(envKeys))
//
// Notes:
// - Only environment variables that are set and non-empty are added to the trace.
func WithEnVars(keys []string) Option {
	return func(t *Trace) {
		for _, key := range keys {
			value, found := os.LookupEnv(key)
			if found && value != "" {
				t.attrs = append(t.attrs, attribute.String(key, value))
			}
		}
	}
}
