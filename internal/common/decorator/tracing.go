package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/semmidev/ethos-go/internal/common/observability"
)

const tracerName = "github.com/semmidev/ethos-go/decorator"

// commandTracingDecorator creates OpenTelemetry spans for command execution
type commandTracingDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandTracingDecorator[C]) Handle(ctx context.Context, cmd C) error {
	handlerName := generateActionName(cmd)
	tracer := otel.Tracer(tracerName)

	ctx, span := tracer.Start(ctx, fmt.Sprintf("Command/%s", handlerName),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("command.type", handlerName),
		attribute.String("component", "command_handler"),
	)

	start := time.Now()
	err := d.base.Handle(ctx, cmd)
	duration := time.Since(start)

	// Record metrics using OTEL metrics
	recordCommandMetrics(ctx, handlerName, err, duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Bool("command.success", false))
	} else {
		span.SetStatus(codes.Ok, "")
		span.SetAttributes(attribute.Bool("command.success", true))
	}

	span.SetAttributes(attribute.Float64("command.duration_ms", float64(duration.Milliseconds())))

	return err
}

// commandResultTracingDecorator creates OpenTelemetry spans for command-with-result execution
type commandResultTracingDecorator[C any, R any] struct {
	base CommandHandlerWithResult[C, R]
}

func (d commandResultTracingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	handlerName := generateActionName(cmd)
	tracer := otel.Tracer(tracerName)

	ctx, span := tracer.Start(ctx, fmt.Sprintf("Command/%s", handlerName),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("command.type", handlerName),
		attribute.String("component", "command_handler"),
	)

	start := time.Now()
	result, err = d.base.Handle(ctx, cmd)
	duration := time.Since(start)

	// Record metrics using OTEL metrics
	recordCommandMetrics(ctx, handlerName, err, duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Bool("command.success", false))
	} else {
		span.SetStatus(codes.Ok, "")
		span.SetAttributes(attribute.Bool("command.success", true))
	}

	span.SetAttributes(attribute.Float64("command.duration_ms", float64(duration.Milliseconds())))

	return result, err
}

// queryTracingDecorator creates OpenTelemetry spans for query execution
type queryTracingDecorator[Q any, R any] struct {
	base QueryHandler[Q, R]
}

func (d queryTracingDecorator[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	handlerName := generateActionName(query)
	tracer := otel.Tracer(tracerName)

	ctx, span := tracer.Start(ctx, fmt.Sprintf("Query/%s", handlerName),
		trace.WithSpanKind(trace.SpanKindInternal),
	)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("query.type", handlerName),
		attribute.String("component", "query_handler"),
	)

	start := time.Now()
	result, err = d.base.Handle(ctx, query)
	duration := time.Since(start)

	// Record metrics using OTEL metrics
	recordQueryMetrics(ctx, handlerName, err, duration)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.SetAttributes(attribute.Bool("query.success", false))
	} else {
		span.SetStatus(codes.Ok, "")
		span.SetAttributes(attribute.Bool("query.success", true))
	}

	span.SetAttributes(attribute.Float64("query.duration_ms", float64(duration.Milliseconds())))

	return result, err
}

// recordCommandMetrics records command execution metrics using OTEL
func recordCommandMetrics(ctx context.Context, command string, err error, duration time.Duration) {
	metrics := observability.GetMetrics()
	if metrics == nil {
		return
	}

	status := "success"
	if err != nil {
		status = "error"
	}

	metrics.RecordCommand(ctx, strings.ToLower(command), status, duration)
}

// recordQueryMetrics records query execution metrics using OTEL
func recordQueryMetrics(ctx context.Context, query string, err error, duration time.Duration) {
	metrics := observability.GetMetrics()
	if metrics == nil {
		return
	}

	status := "success"
	if err != nil {
		status = "error"
	}

	metrics.RecordQuery(ctx, strings.ToLower(query), status, duration)
}
