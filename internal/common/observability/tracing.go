package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "github.com/semmidev/ethos-go"

// Tracer returns a tracer for the given name
func Tracer(name string) trace.Tracer {
	return otel.Tracer(instrumentationName + "/" + name)
}

// StartSpan starts a new span with the given name
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return otel.Tracer(instrumentationName).Start(ctx, name, opts...)
}

// SpanFromContext returns the current span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddSpanEvent adds an event to the current span
func AddSpanEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetSpanError records an error on the current span
func SetSpanError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// SetSpanAttributes sets attributes on the current span
func SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// Common attribute keys
var (
	AttrUserID      = attribute.Key("user.id")
	AttrSessionID   = attribute.Key("session.id")
	AttrHabitID     = attribute.Key("habit.id")
	AttrOperation   = attribute.Key("operation")
	AttrComponent   = attribute.Key("component")
	AttrDBOperation = attribute.Key("db.operation")
	AttrDBTable     = attribute.Key("db.table")
	AttrHTTPMethod  = attribute.Key("http.method")
	AttrHTTPPath    = attribute.Key("http.path")
	AttrHTTPStatus  = attribute.Key("http.status_code")
)
