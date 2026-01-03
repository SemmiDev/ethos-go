package decorator

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], log logger.Logger, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
		logger: log,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
