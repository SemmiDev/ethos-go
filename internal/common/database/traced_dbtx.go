package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "github.com/semmidev/ethos-go/database"

// TracedDBTX wraps a DBTX with OpenTelemetry tracing.
// This implements the DBTX interface and adds tracing to all database operations.
type TracedDBTX struct {
	db     DBTX
	tracer trace.Tracer
}

// NewTracedDBTX creates a new traced database wrapper.
// Use this to wrap your *sqlx.DB or *sqlx.Tx to get automatic tracing.
//
// Example:
//
//	db, _ := sqlx.Connect("postgres", dsn)
//	tracedDB := database.NewTracedDBTX(db)
//	userRepo := adapters.NewUserPostgresRepository(tracedDB)
func NewTracedDBTX(db DBTX) *TracedDBTX {
	return &TracedDBTX{
		db:     db,
		tracer: otel.Tracer(tracerName),
	}
}

// ExecContext executes a query with tracing
func (t *TracedDBTX) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	ctx, span := t.tracer.Start(ctx, "db.exec",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "exec"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	result, err := t.db.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Float64("db.duration_ms", float64(duration.Milliseconds())))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
		if rowsAffected, raErr := result.RowsAffected(); raErr == nil {
			span.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
		}
	}

	return result, err
}

// GetContext gets a single record with tracing
func (t *TracedDBTX) GetContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	ctx, span := t.tracer.Start(ctx, "db.get",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "get"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	err := t.db.GetContext(ctx, dest, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Float64("db.duration_ms", float64(duration.Milliseconds())))

	if err != nil && err != sql.ErrNoRows {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
		if err == sql.ErrNoRows {
			span.SetAttributes(attribute.Bool("db.not_found", true))
		}
	}

	return err
}

// SelectContext selects multiple records with tracing
func (t *TracedDBTX) SelectContext(ctx context.Context, dest interface{}, query string, args ...any) error {
	ctx, span := t.tracer.Start(ctx, "db.select",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "select"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	err := t.db.SelectContext(ctx, dest, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Float64("db.duration_ms", float64(duration.Milliseconds())))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return err
}

// QueryRowxContext queries a single row with tracing
func (t *TracedDBTX) QueryRowxContext(ctx context.Context, query string, args ...any) *sqlx.Row {
	_, span := t.tracer.Start(ctx, "db.query_row",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	// Note: We can't defer span.End() here because the row scan happens after return.
	// The span will end when context is done. For proper tracing, prefer GetContext.
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "query_row"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	return t.db.QueryRowxContext(ctx, query, args...)
}

// QueryxContext queries and returns rows with tracing
func (t *TracedDBTX) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	ctx, span := t.tracer.Start(ctx, "db.query",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "query"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	rows, err := t.db.QueryxContext(ctx, query, args...)
	duration := time.Since(start)

	span.SetAttributes(attribute.Float64("db.duration_ms", float64(duration.Milliseconds())))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return rows, err
}

// PreparexContext prepares a statement with tracing
func (t *TracedDBTX) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	ctx, span := t.tracer.Start(ctx, "db.prepare",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "prepare"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	stmt, err := t.db.PreparexContext(ctx, query)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return stmt, err
}

// Rebind delegates to the underlying DBTX
func (t *TracedDBTX) Rebind(query string) string {
	return t.db.Rebind(query)
}

// DriverName delegates to the underlying DBTX
func (t *TracedDBTX) DriverName() string {
	return t.db.DriverName()
}

// NamedExecContext executes a named query with tracing
func (t *TracedDBTX) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	ctx, span := t.tracer.Start(ctx, "db.named_exec",
		trace.WithSpanKind(trace.SpanKindClient),
	)
	defer span.End()

	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.operation", "named_exec"),
		attribute.String("db.statement", truncateQuery(query)),
	)

	start := time.Now()
	result, err := t.db.NamedExecContext(ctx, query, arg)
	duration := time.Since(start)

	span.SetAttributes(attribute.Float64("db.duration_ms", float64(duration.Milliseconds())))

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
		if rowsAffected, raErr := result.RowsAffected(); raErr == nil {
			span.SetAttributes(attribute.Int64("db.rows_affected", rowsAffected))
		}
	}

	return result, err
}

// Unwrap returns the underlying DBTX (useful for transactions)
func (t *TracedDBTX) Unwrap() DBTX {
	return t.db
}

func truncateQuery(query string) string {
	const maxLen = 500
	if len(query) > maxLen {
		return query[:maxLen] + "..."
	}
	return query
}

// Compile-time check to ensure TracedDBTX implements DBTX
var _ DBTX = (*TracedDBTX)(nil)
