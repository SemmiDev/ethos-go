package logger

import "context"

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
}

// Logger defines the logging interface
type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warn(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, err error, msg string, fields ...Field)
	With(fields ...Field) Logger
}
