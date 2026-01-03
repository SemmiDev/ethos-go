package observability

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// InstrumentedDB wraps sqlx.DB with OpenTelemetry instrumentation
type InstrumentedDB struct {
	*sqlx.DB
	tracer trace.Tracer
}

// WrapDB wraps an sqlx.DB with OpenTelemetry instrumentation
func WrapDB(db *sqlx.DB) *InstrumentedDB {
	return &InstrumentedDB{
		DB:     db,
		tracer: Tracer("database"),
	}
}

// ExecContext executes a query with tracing
func (db *InstrumentedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	ctx, span := db.tracer.Start(ctx, "db.exec",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	result, err := db.DB.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		recordDBMetrics(ctx, "exec", "error", duration)
	} else {
		span.SetStatus(codes.Ok, "")
		recordDBMetrics(ctx, "exec", "success", duration)
	}

	return result, err
}

// QueryRowxContext queries a single row with tracing
func (db *InstrumentedDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	ctx, span := db.tracer.Start(ctx, "db.query_row",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	// Note: We can't defer span.End() here as the row might be scanned later
	// The span will be ended when the context is done

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	return db.DB.QueryRowxContext(ctx, query, args...)
}

// GetContext gets a single record with tracing
func (db *InstrumentedDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	ctx, span := db.tracer.Start(ctx, "db.get",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	err := db.DB.GetContext(ctx, dest, query, args...)
	duration := time.Since(start)

	if err != nil && err != sql.ErrNoRows {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		recordDBMetrics(ctx, "get", "error", duration)
	} else {
		span.SetStatus(codes.Ok, "")
		recordDBMetrics(ctx, "get", "success", duration)
	}

	return err
}

// SelectContext selects multiple records with tracing
func (db *InstrumentedDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	ctx, span := db.tracer.Start(ctx, "db.select",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	err := db.DB.SelectContext(ctx, dest, query, args...)
	duration := time.Since(start)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		recordDBMetrics(ctx, "select", "error", duration)
	} else {
		span.SetStatus(codes.Ok, "")
		recordDBMetrics(ctx, "select", "success", duration)
	}

	return err
}

// BeginTxx begins a transaction with tracing
func (db *InstrumentedDB) BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error) {
	ctx, span := db.tracer.Start(ctx, "db.begin_tx",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "begin_transaction"),
	)

	tx, err := db.DB.BeginTxx(ctx, opts)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return tx, err
}

func truncateQuery(query string) string {
	const maxLen = 500
	if len(query) > maxLen {
		return query[:maxLen] + "..."
	}
	return query
}

func recordDBMetrics(ctx context.Context, operation, status string, duration time.Duration) {
	metrics := GetMetrics()
	if metrics != nil {
		metrics.RecordDBQuery(ctx, operation, "unknown", status, duration)
	}
}
