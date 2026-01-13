package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

// ApplyCommandDecorators wraps a command handler with tracing, logging and metrics.
// Decorator chain (outermost to innermost): tracing -> logging -> metrics -> handler
// - Tracing: Creates OpenTelemetry spans for distributed tracing
// - Logging: Enriches wide events with command context
// - Metrics: Records command counts and durations (legacy counter interface)
func ApplyCommandDecorators[H any](handler CommandHandler[H], _ logger.Logger, metricsClient MetricsClient) CommandHandler[H] {
	return commandTracingDecorator[H]{
		base: commandLoggingDecorator[H]{
			base: commandMetricsDecorator[H]{
				base:   handler,
				client: metricsClient,
			},
		},
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

// ApplyCommandResultDecorators wraps a command-with-result handler with tracing, logging and metrics.
// Decorator chain (outermost to innermost): tracing -> logging -> metrics -> handler
func ApplyCommandResultDecorators[H any, R any](handler CommandHandlerWithResult[H, R], _ logger.Logger, metricsClient MetricsClient) CommandHandlerWithResult[H, R] {
	return commandResultTracingDecorator[H, R]{
		base: commandResultLoggingDecorator[H, R]{
			base: commandResultMetricsDecorator[H, R]{
				base:   handler,
				client: metricsClient,
			},
		},
	}
}

type CommandHandlerWithResult[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
