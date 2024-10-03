package traceflow

import (
	"reflect"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// AddException adds exception information to the trace.
func (t *Trace) AddException(err error, stackTrace string) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("exception.type", reflect.TypeOf(err).String()),
		attribute.String("exception.message", err.Error()),
		attribute.String("exception.stacktrace", stackTrace),
	)

	return t
}

// AddError adds error information to the trace.
func (t *Trace) AddError(err error) *Trace {
	t.attrs = append(t.attrs,
		attribute.String("error.type", reflect.TypeOf(err).String()),
		attribute.String("error.message", err.Error()),
	)

	return t
}

// RecordError records an error to the span and sets the span status to Error.
func (t *Trace) RecordError(err error) {
	if err != nil {
		t.span.RecordError(err)
		t.span.SetStatus(codes.Error, err.Error())
	}
}
