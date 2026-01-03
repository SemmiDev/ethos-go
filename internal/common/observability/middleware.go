package observability

import (
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// HTTPMiddleware returns an HTTP middleware that instruments requests with OpenTelemetry
func HTTPMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Wrap with otelhttp for automatic tracing
		handler := otelhttp.NewHandler(next, serviceName,
			otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
				return r.Method + " " + r.URL.Path
			}),
		)

		// Add custom metrics on top
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			// Track active requests
			metrics := GetMetrics()
			if metrics != nil {
				metrics.HTTPRequestsActive.Add(r.Context(), 1)
				defer metrics.HTTPRequestsActive.Add(r.Context(), -1)
			}

			// Serve request
			handler.ServeHTTP(wrapped, r)

			// Record metrics
			duration := time.Since(start)
			if metrics != nil {
				metrics.RecordHTTPRequest(r.Context(), r.Method, r.URL.Path, wrapped.statusCode, duration)
			}

			// Add response attributes to span
			span := trace.SpanFromContext(r.Context())
			span.SetAttributes(
				attribute.Int("http.status_code", wrapped.statusCode),
				attribute.Int64("http.response_content_length", int64(wrapped.bytesWritten)),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// Unwrap returns the original ResponseWriter for compatibility
func (rw *responseWriter) Unwrap() http.ResponseWriter {
	return rw.ResponseWriter
}
