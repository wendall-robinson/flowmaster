package traceflow

import (
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
)

// Attribute wraps OpenTelemetry's KeyValue type
type Attribute struct {
	otelAttr attribute.KeyValue
}

// AddString creates a string OTEL attribute.
func AddString(key, value string) Attribute {
	return Attribute{
		otelAttr: attribute.String(key, value),
	}
}

// AddStringSlice creates a string slice OTEL attribute.
func AddStringSlice(key string, value []string) Attribute {
	return Attribute{
		otelAttr: attribute.StringSlice(key, value),
	}
}

// AddInt creates an int OTEL attribute.
func AddInt(key string, value int) Attribute {
	return Attribute{
		otelAttr: attribute.Int(key, value),
	}
}

// AddIntSlice creates an int slice OTEL attribute.
func AddIntSlice(key string, value []int) Attribute {
	return Attribute{
		otelAttr: attribute.IntSlice(key, value),
	}
}

// AddFloat creates an int64 OTEL attribute.
func AddFloat(key string, value float64) Attribute {
	return Attribute{
		otelAttr: attribute.Float64(key, value),
	}
}

// AddFloatSlice creates a float64 slice OTEL attribute.
func AddFloatSlice(key string, value []float64) Attribute {
	return Attribute{
		otelAttr: attribute.Float64Slice(key, value),
	}
}

// AddBool creates a bool OTEL attribute.
func AddBool(key string, value bool) Attribute {
	return Attribute{
		otelAttr: attribute.Bool(key, value),
	}
}

// AddBoolSlice creates a bool slice OTEL attribute.
func AddBoolSlice(key string, value []bool) Attribute {
	return Attribute{
		otelAttr: attribute.BoolSlice(key, value),
	}
}

// AddJSON adds a JSON payload as a string attribute to the trace.
// The JSON payload is passed as a json.RawMessage, which is then converted
// to a string and added to the trace as an OpenTelemetry string attribute
// with the key "payload". This method is useful when you want to include
// JSON-encoded data as part of the trace's attributes.
//
// Example usage:
//
//	jsonPayload := json.RawMessage(`{"key":"value"}`)
//	trace.AddJSON(jsonPayload)
//
// The resulting trace attribute will include the JSON data as a string:
//
//	"payload": "{\"key\":\"value\"}"
//
// This method is particularly useful when you need to include structured data
// (such as API responses or request bodies) in your traces, but want to store
// it as a single string attribute.
//
// Note that the JSON data is not parsed or validated, and is added directly
// as a string. Be mindful of the size of the JSON payload, as OpenTelemetry
// attributes have practical size limits that should not be exceeded.
func (t *Trace) AddJSON(payload json.RawMessage) *Trace {
	jsonString := string(payload)
	jsonAttr := attribute.String("payload", jsonString)
	t.attrs = append(t.attrs, jsonAttr)

	return t
}
