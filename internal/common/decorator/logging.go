package decorator

import (
	"context"
	"time"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

// commandLoggingDecorator enriches the wide event with command context.
// Instead of logging separately, it adds command info to the request's wide event.
type commandLoggingDecorator[C any] struct {
	base CommandHandler[C]
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)
	startTime := time.Now()

	// Add command context to wide event
	logger.SetCustom(ctx, "command", handlerType)
	logger.SetCustom(ctx, "command_payload", cmd)

	err = d.base.Handle(ctx, cmd)

	// Record command duration
	durationMs := time.Since(startTime).Milliseconds()
	logger.SetCustom(ctx, "command_duration_ms", durationMs)

	// If there's an error, enrich the wide event with error context
	if err != nil {
		logger.AddError(ctx, "CommandError", handlerType, err.Error(), false)
	}

	return err
}

// commandResultLoggingDecorator enriches the wide event for commands with results.
type commandResultLoggingDecorator[C any, R any] struct {
	base CommandHandlerWithResult[C, R]
}

func (d commandResultLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	handlerType := generateActionName(cmd)
	startTime := time.Now()

	// Add command context to wide event
	logger.SetCustom(ctx, "command", handlerType)
	logger.SetCustom(ctx, "command_payload", cmd)

	result, err = d.base.Handle(ctx, cmd)

	// Record command duration
	durationMs := time.Since(startTime).Milliseconds()
	logger.SetCustom(ctx, "command_duration_ms", durationMs)

	// If there's an error, enrich the wide event with error context
	if err != nil {
		logger.AddError(ctx, "CommandError", handlerType, err.Error(), false)
	}

	return result, err
}

// queryLoggingDecorator enriches the wide event with query context.
type queryLoggingDecorator[C any, R any] struct {
	base QueryHandler[C, R]
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, query C) (result R, err error) {
	handlerType := generateActionName(query)
	startTime := time.Now()

	// Add query context to wide event
	logger.SetCustom(ctx, "query", handlerType)
	logger.SetCustom(ctx, "query_payload", query)

	result, err = d.base.Handle(ctx, query)

	// Record query duration
	durationMs := time.Since(startTime).Milliseconds()
	logger.SetCustom(ctx, "query_duration_ms", durationMs)

	// If there's an error, enrich the wide event with error context
	if err != nil {
		logger.AddError(ctx, "QueryError", handlerType, err.Error(), false)
	}

	return result, err
}
