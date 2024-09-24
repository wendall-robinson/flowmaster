package gotraceit

import "go.opentelemetry.io/otel/attribute"

// Attribute wraps the OpenTelemetry attribute type
type Attribute attribute.KeyValue

// StringAttr creates a string OTEL attribute.
func StringAttr(key, value string) Attribute {
	return Attribute(attribute.String(key, value))
}

// StringSliceAttr creates a string slice OTEL attribute.
func StringSliceAttr(key string, value []string) Attribute {
	return Attribute(attribute.StringSlice(key, value))
}

// IntAttr creates an int OTEL attribute.
func IntAttr(key string, value int) Attribute {
	return Attribute(attribute.Int(key, value))
}

// IntSliceAttr creates an int slice OTEL attribute.
func IntSliceAttr(key string, value []int) Attribute {
	return Attribute(attribute.IntSlice(key, value))
}

// FloatAttr creates an int64 OTEL attribute.
func FloatAttr(key string, value float64) Attribute {
	return Attribute(attribute.Float64(key, value))
}

// FloatSliceAttr creates a float64 slice OTEL attribute.
func FloatSliceAttr(key string, value []float64) Attribute {
	return Attribute(attribute.Float64Slice(key, value))
}

// BoolAttr creates a bool OTEL attribute.
func BoolAttr(key string, value bool) Attribute {
	return Attribute(attribute.Bool(key, value))
}

// BoolSliceAttr creates a bool slice OTEL attribute.
func BoolSliceAttr(key string, value []bool) Attribute {
	return Attribute(attribute.BoolSlice(key, value))
}
