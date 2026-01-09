package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	"github.com/semmidev/ethos-go/internal/habits/domain/habitlog"
)

// DeleteHabitLog command encapsulates habit log deletion
type DeleteHabitLog struct {
	LogID  string `validate:"uuid"`
	UserID string `validate:"uuid"`
}

// DeleteHabitLogHandler processes habit log deletion commands
type DeleteHabitLogHandler decorator.CommandHandler[DeleteHabitLog]

type deleteHabitLogHandler struct {
	repo      habitlog.Repository
	validator *validator.Validator
}

// NewDeleteHabitLogHandler creates a new handler with decorators
func NewDeleteHabitLogHandler(
	repo habitlog.Repository,
	validator *validator.Validator,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) DeleteHabitLogHandler {
	if repo == nil {
		panic("nil habit log repository")
	}

	return decorator.ApplyCommandDecorators(
		deleteHabitLogHandler{
			repo:      repo,
			validator: validator,
		},
		log,
		metricsClient,
	)
}

func (h deleteHabitLogHandler) Handle(ctx context.Context, cmd DeleteHabitLog) error {
	// Validate input
	if err := h.validator.Validate(cmd); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			details := make(map[string]interface{})
			for _, ve := range validationErrors {
				details[ve.Field] = ve.Message
			}
			return apperror.ValidationFailedWithDetails("validation failed", details)
		}
		return apperror.ValidationFailed(err.Error())
	}

	return h.repo.DeleteHabitLog(ctx, cmd.LogID, cmd.UserID)
}
