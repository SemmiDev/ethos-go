package logger

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel/trace"
)

// EventMiddlewareConfig holds configuration for the event middleware
type EventMiddlewareConfig struct {
	ServiceName string
	Version     string
	Environment string
	Logger      Logger
	Sampler     *Sampler
}

// EventMiddleware creates an HTTP middleware that implements the Canonical Log Line pattern.
// It initializes an Event at request start, makes it available via context,
// and emits a single comprehensive log event at request end.
func EventMiddleware(cfg EventMiddlewareConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			// Initialize the event with request context
			event := &Event{
				Timestamp:    startTime,
				Service:      cfg.ServiceName,
				Version:      cfg.Version,
				Environment:  cfg.Environment,
				Method:       r.Method,
				Path:         r.URL.Path,
				Query:        r.URL.RawQuery,
				ClientIP:     getClientIP(r),
				UserAgent:    r.UserAgent(),
				FeatureFlags: make(map[string]bool),
				Custom:       make(map[string]any),
			}

			// Extract request ID from context (set by chi middleware.RequestID)
			if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
				if id, ok := reqID.(string); ok {
					event.RequestID = id
				}
			}

			// Extract request ID from header if not in context
			if event.RequestID == "" {
				event.RequestID = r.Header.Get("X-Request-Id")
			}

			// Extract trace context if available
			span := trace.SpanFromContext(r.Context())
			if span.SpanContext().IsValid() {
				event.TraceID = span.SpanContext().TraceID().String()
				event.SpanID = span.SpanContext().SpanID().String()
			}

			// Store event in context for handlers to enrich
			ctx := WithEvent(r.Context(), event)
			r = r.WithContext(ctx)

			// Wrap ResponseWriter to capture status code and bytes
			wrw := &wrapResponseWriter{ResponseWriter: w, statusCode: 200}

			// Execute handler with panic recovery
			func() {
				defer func() {
					if rec := recover(); rec != nil {
						event.StatusCode = 500
						event.Outcome = "error"
						event.Error = &ErrorContext{
							Type:    "PanicError",
							Code:    "panic",
							Message: "Internal server error (panic recovered)",
						}
						wrw.WriteHeader(500)
					}
				}()

				next.ServeHTTP(wrw, r)
			}()

			// Finalize event after request completes
			event.StatusCode = wrw.statusCode
			event.BytesSent = wrw.bytesWritten
			event.DurationMs = time.Since(startTime).Milliseconds()

			// Set outcome if not already set by error handling
			if event.Outcome == "" {
				if event.StatusCode >= 200 && event.StatusCode < 400 {
					event.Outcome = "success"
				} else {
					event.Outcome = "error"
				}
			}

			// Apply tail sampling - only emit if sampler says yes
			if cfg.Sampler == nil || cfg.Sampler.ShouldSample(event) {
				// Emit the single, comprehensive event
				cfg.Logger.Info(ctx, "request_completed",
					Field{Key: "event", Value: event},
				)
			}
		})
	}
}

// Deprecated: WideEventMiddlewareConfig is deprecated, use EventMiddlewareConfig instead.
type WideEventMiddlewareConfig = EventMiddlewareConfig

// Deprecated: WideEventMiddleware is deprecated, use EventMiddleware instead.
func WideEventMiddleware(cfg WideEventMiddlewareConfig) func(http.Handler) http.Handler {
	return EventMiddleware(cfg)
}

// wrapResponseWriter captures status code and bytes written
type wrapResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int64
	wroteHeader  bool
}

func (w *wrapResponseWriter) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.statusCode = statusCode
		w.wroteHeader = true
		w.ResponseWriter.WriteHeader(statusCode)
	}
}

func (w *wrapResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.wroteHeader = true
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytesWritten += int64(n)
	return n, err
}

// Unwrap returns the original ResponseWriter for compatibility with http.Flusher, etc.
func (w *wrapResponseWriter) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// getClientIP extracts the real client IP from request headers
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (common for proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return r.RemoteAddr
}
