package decorator

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

// ApplyQueryDecorators wraps a query handler with tracing, logging and metrics.
// Decorator chain (outermost to innermost): tracing -> logging -> metrics -> handler
// - Tracing: Creates OpenTelemetry spans for distributed tracing
// - Logging: Enriches wide events with query context
// - Metrics: Records query counts and durations (legacy counter interface)
func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], _ logger.Logger, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryTracingDecorator[H, R]{
		base: queryLoggingDecorator[H, R]{
			base: queryMetricsDecorator[H, R]{
				base:   handler,
				client: metricsClient,
			},
		},
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
