package observability

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Metrics holds all application metrics
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal   metric.Int64Counter
	HTTPRequestDuration metric.Float64Histogram
	HTTPRequestsActive  metric.Int64UpDownCounter

	// Database metrics
	DBQueryDuration   metric.Float64Histogram
	DBQueriesTotal    metric.Int64Counter
	DBConnectionsOpen metric.Int64UpDownCounter

	// Application metrics
	CommandsTotal   metric.Int64Counter
	CommandDuration metric.Float64Histogram
	QueriesTotal    metric.Int64Counter
	QueryDuration   metric.Float64Histogram

	// Auth metrics
	AuthAttemptsTotal metric.Int64Counter
	ActiveSessions    metric.Int64UpDownCounter

	// Business metrics
	HabitsCreated   metric.Int64Counter
	HabitLogsTotal  metric.Int64Counter
	UsersRegistered metric.Int64Counter
}

var globalMetrics *Metrics

// InitMetrics initializes all application metrics
func InitMetrics(ctx context.Context) (*Metrics, error) {
	meter := otel.Meter(instrumentationName)

	m := &Metrics{}
	var err error

	// HTTP metrics
	m.HTTPRequestsTotal, err = meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	m.HTTPRequestDuration, err = meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return nil, err
	}

	m.HTTPRequestsActive, err = meter.Int64UpDownCounter(
		"http_requests_active",
		metric.WithDescription("Number of active HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	// Database metrics
	m.DBQueryDuration, err = meter.Float64Histogram(
		"db_query_duration_seconds",
		metric.WithDescription("Database query duration in seconds"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5),
	)
	if err != nil {
		return nil, err
	}

	m.DBQueriesTotal, err = meter.Int64Counter(
		"db_queries_total",
		metric.WithDescription("Total number of database queries"),
		metric.WithUnit("{query}"),
	)
	if err != nil {
		return nil, err
	}

	m.DBConnectionsOpen, err = meter.Int64UpDownCounter(
		"db_connections_open",
		metric.WithDescription("Number of open database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	// Application metrics
	m.CommandsTotal, err = meter.Int64Counter(
		"commands_total",
		metric.WithDescription("Total number of commands executed"),
		metric.WithUnit("{command}"),
	)
	if err != nil {
		return nil, err
	}

	m.CommandDuration, err = meter.Float64Histogram(
		"command_duration_seconds",
		metric.WithDescription("Command execution duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	m.QueriesTotal, err = meter.Int64Counter(
		"app_queries_total",
		metric.WithDescription("Total number of application queries executed"),
		metric.WithUnit("{query}"),
	)
	if err != nil {
		return nil, err
	}

	m.QueryDuration, err = meter.Float64Histogram(
		"app_query_duration_seconds",
		metric.WithDescription("Application query execution duration in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	// Auth metrics
	m.AuthAttemptsTotal, err = meter.Int64Counter(
		"auth_attempts_total",
		metric.WithDescription("Total number of authentication attempts"),
		metric.WithUnit("{attempt}"),
	)
	if err != nil {
		return nil, err
	}

	m.ActiveSessions, err = meter.Int64UpDownCounter(
		"active_sessions",
		metric.WithDescription("Number of active user sessions"),
		metric.WithUnit("{session}"),
	)
	if err != nil {
		return nil, err
	}

	// Business metrics
	m.HabitsCreated, err = meter.Int64Counter(
		"habits_created_total",
		metric.WithDescription("Total number of habits created"),
		metric.WithUnit("{habit}"),
	)
	if err != nil {
		return nil, err
	}

	m.HabitLogsTotal, err = meter.Int64Counter(
		"habit_logs_total",
		metric.WithDescription("Total number of habit logs created"),
		metric.WithUnit("{log}"),
	)
	if err != nil {
		return nil, err
	}

	m.UsersRegistered, err = meter.Int64Counter(
		"users_registered_total",
		metric.WithDescription("Total number of users registered"),
		metric.WithUnit("{user}"),
	)
	if err != nil {
		return nil, err
	}

	globalMetrics = m
	return m, nil
}

// GetMetrics returns the global metrics instance
func GetMetrics() *Metrics {
	return globalMetrics
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("http.method", method),
		attribute.String("http.path", path),
		attribute.Int("http.status_code", statusCode),
	}

	m.HTTPRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.HTTPRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// RecordDBQuery records database query metrics
func (m *Metrics) RecordDBQuery(ctx context.Context, operation, table, status string, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("db.operation", operation),
		attribute.String("db.table", table),
		attribute.String("status", status),
	}

	m.DBQueriesTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.DBQueryDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// RecordCommand records command execution metrics
func (m *Metrics) RecordCommand(ctx context.Context, command, status string, duration time.Duration) {
	attrs := []attribute.KeyValue{
		attribute.String("command", command),
		attribute.String("status", status),
	}

	m.CommandsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.CommandDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs...))
}

// RecordAuthAttempt records authentication attempt
func (m *Metrics) RecordAuthAttempt(ctx context.Context, authType, status string) {
	attrs := []attribute.KeyValue{
		attribute.String("type", authType),
		attribute.String("status", status),
	}

	m.AuthAttemptsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}
