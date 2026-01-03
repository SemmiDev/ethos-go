package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], log logger.Logger, metricsClient MetricsClient) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: commandMetricsDecorator[H]{
			base:   handler,
			client: metricsClient,
		},
		logger: log,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func ApplyCommandResultDecorators[H any, R any](handler CommandHandlerWithResult[H, R], log logger.Logger, metricsClient MetricsClient) CommandHandlerWithResult[H, R] {
	return commandResultLoggingDecorator[H, R]{
		base: commandResultMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
		logger: log,
	}
}

type CommandHandlerWithResult[C any, R any] interface {
	Handle(ctx context.Context, cmd C) (R, error)
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
