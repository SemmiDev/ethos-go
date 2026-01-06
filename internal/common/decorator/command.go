package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

// ApplyCommandDecorators wraps a command handler with wide event enrichment and metrics.
// The logging decorator now enriches the wide event instead of logging separately.
func ApplyCommandDecorators[H any](handler CommandHandler[H], _ logger.Logger, metricsClient MetricsClient) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: commandMetricsDecorator[H]{
			base:   handler,
			client: metricsClient,
		},
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

// ApplyCommandResultDecorators wraps a command-with-result handler with wide event enrichment and metrics.
func ApplyCommandResultDecorators[H any, R any](handler CommandHandlerWithResult[H, R], _ logger.Logger, metricsClient MetricsClient) CommandHandlerWithResult[H, R] {
	return commandResultLoggingDecorator[H, R]{
		base: commandResultMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
	}
}

type CommandHandlerWithResult[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
