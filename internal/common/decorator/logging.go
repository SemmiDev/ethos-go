package decorator

import (
	"context"
	"fmt"

	"github.com/semmidev/ethos-go/internal/common/logger"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger logger.Logger
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	handlerType := generateActionName(cmd)

	log := d.logger.With(
		logger.Field{Key: "command", Value: handlerType},
		logger.Field{Key: "command_body", Value: fmt.Sprintf("%#v", cmd)},
	)

	log.Debug(ctx, "Executing command")
	defer func() {
		if err == nil {
			log.Info(ctx, "Command executed successfully")
		} else {
			log.Error(ctx, err, "Failed to execute command")
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type commandResultLoggingDecorator[C any, R any] struct {
	base   CommandHandlerWithResult[C, R]
	logger logger.Logger
}

func (d commandResultLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	handlerType := generateActionName(cmd)

	log := d.logger.With(
		logger.Field{Key: "command", Value: handlerType},
		logger.Field{Key: "command_body", Value: fmt.Sprintf("%#v", cmd)},
	)

	log.Debug(ctx, "Executing command")
	defer func() {
		if err == nil {
			log.Info(ctx, "Command executed successfully")
		} else {
			log.Error(ctx, err, "Failed to execute command")
		}
	}()

	return d.base.Handle(ctx, cmd)
}

type queryLoggingDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	logger logger.Logger
}

func (d queryLoggingDecorator[C, R]) Handle(ctx context.Context, cmd C) (result R, err error) {
	log := d.logger.With(
		logger.Field{Key: "query", Value: generateActionName(cmd)},
		logger.Field{Key: "query_body", Value: fmt.Sprintf("%#v", cmd)},
	)

	log.Debug(ctx, "Executing query")
	defer func() {
		if err == nil {
			log.Info(ctx, "Query executed successfully")
		} else {
			log.Error(ctx, err, "Failed to execute query")
		}
	}()

	return d.base.Handle(ctx, cmd)
}
