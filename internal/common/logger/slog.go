package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	slogmulti "github.com/samber/slog-multi"
	"github.com/semmidev/ethos-go/config"
	"go.opentelemetry.io/otel/trace"
	"gopkg.in/natefinch/lumberjack.v2"
)

type slogLogger struct {
	handler slog.Handler
}

////////////////////////////////////////////////////////////////
// PUBLIC CONSTRUCTOR
////////////////////////////////////////////////////////////////

func New(cfg *config.Config) (Logger, error) {
	var handlers []slog.Handler

	level := parseLevel(cfg.LoggerLevel)

	rotator := &lumberjack.Logger{
		// Lokasi & nama file log utama. ->/var/log/app.log
		Filename: cfg.LoggerFile,
		// Ukuran file log sebelum rotasi. -> 10MB
		MaxSize: cfg.LoggerMaxSize,
		// Mengaktifkan kompresi file log. -> true
		Compress: cfg.LoggerCompress,
		// Tanpa MaxBackups: log bisa banyak sekali dalam waktu singkat.
		// Batas Jumlah file log lama yang disimpan, Lebih dari itu yang tertua dihapus. -> 3
		MaxBackups: cfg.LoggerMaxBackups,
		// Tanpa MaxAge: log lama bisa tersimpan selamanya walau jarang dipakai.
		MaxAge: cfg.LoggerMaxAge, // Usia file log sebelum dihapus. -> 28 hari
	}

	for output := range strings.SplitSeq(cfg.LoggerOutput, "|") {
		switch strings.TrimSpace(output) {
		case "stdout":
			handlers = append(handlers, slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level:       level,
				AddSource:   true,
				ReplaceAttr: replaceTime,
			}))
		case "file":
			handlers = append(handlers, slog.NewJSONHandler(rotator, &slog.HandlerOptions{
				Level:       level,
				AddSource:   true,
				ReplaceAttr: replaceTime,
			}))
		}
	}

	if len(handlers) == 0 {
		return nil, fmt.Errorf("no valid logger output configured")
	}

	var combined slog.Handler
	if len(handlers) == 1 {
		combined = handlers[0]
	} else {
		combined = slogmulti.Fanout(handlers...)
	}

	// Note: Log sampling is handled by the EventMiddleware's Sampler for HTTP requests.
	// This logger outputs all logs passed to it; the sampling decision happens
	// at the request level in the middleware, not at the individual log level.
	handler := slogmulti.
		Pipe(newRequestIDMiddleware()).
		Handler(combined)

	return &slogLogger{handler: handler}, nil
}

////////////////////////////////////////////////////////////////
// LOGGER METHODS (REAL SOURCE)
////////////////////////////////////////////////////////////////

func (l *slogLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelDebug, msg, nil, fields...)
}

func (l *slogLogger) Info(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelInfo, msg, nil, fields...)
}

func (l *slogLogger) Warn(ctx context.Context, msg string, fields ...Field) {
	l.log(ctx, slog.LevelWarn, msg, nil, fields...)
}

func (l *slogLogger) Error(ctx context.Context, err error, msg string, fields ...Field) {
	l.log(ctx, slog.LevelError, msg, err, fields...)
}

func (l *slogLogger) With(fields ...Field) Logger {
	return &slogLogger{
		handler: l.handler.WithAttrs(toAttrs(fields)),
	}
}

func (l *slogLogger) log(
	ctx context.Context,
	level slog.Level,
	msg string,
	err error,
	fields ...Field,
) {
	if err != nil {
		fields = append(fields, Field{
			Key:   "error",
			Value: err.Error(),
		})
	}

	r := slog.NewRecord(
		time.Now(),
		level,
		msg,
		callerPC(), // REAL CALL SITE
	)

	r.AddAttrs(toAttrs(fields)...)

	_ = l.handler.Handle(ctx, r)
}

// REQUEST ID MIDDLEWARE

type requestIDHandler struct {
	next slog.Handler
}

func newRequestIDMiddleware() func(slog.Handler) slog.Handler {
	return func(next slog.Handler) slog.Handler {
		return &requestIDHandler{next: next}
	}
}

func (h *requestIDHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *requestIDHandler) Handle(ctx context.Context, r slog.Record) error {
	if id, ok := ctx.Value(middleware.RequestIDKey).(string); ok && id != "" {
		r.AddAttrs(slog.String("request_id", id))
	}

	// Add TraceID and SpanID if available in context
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	return h.next.Handle(ctx, r)
}

func (h *requestIDHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &requestIDHandler{next: h.next.WithAttrs(attrs)}
}

func (h *requestIDHandler) WithGroup(name string) slog.Handler {
	return &requestIDHandler{next: h.next.WithGroup(name)}
}

// HELPERS

// callerPC: pc = program counter → alamat instruksi di memory tempat log itu dipanggil.
func callerPC() uintptr {
	// stack:
	// Lewati 3 stack frame, dan ambil lokasi kode yang memanggil logger.
	// app → interface → slogLogger → runtime.Caller
	const skip = 3

	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return 0
	}
	return pc
}

func toAttrs(fields []Field) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(fields))
	for _, f := range fields {
		attrs = append(attrs, slog.Any(f.Key, f.Value))
	}
	return attrs
}

func replaceTime(_ []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		a.Value = slog.StringValue(time.Now().UTC().Format(time.RFC3339))
	}
	return a
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
