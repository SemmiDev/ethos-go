package decorator

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

// ApplyQueryDecorators wraps a query handler with wide event enrichment and metrics.
// The logging decorator now enriches the wide event instead of logging separately.
func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], _ logger.Logger, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
