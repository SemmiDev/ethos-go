package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
)

// DeleteHabit command encapsulates habit deletion
type DeleteHabit struct {
	HabitID string `validate:"uuid"`
	UserID  string `validate:"uuid"`
}

// DeleteHabitHandler processes habit deletion commands
type DeleteHabitHandler decorator.CommandHandler[DeleteHabit]

type deleteHabitHandler struct {
	repo      habit.Repository
	validator *validator.Validator
}

// NewDeleteHabitHandler creates a new handler with decorators
func NewDeleteHabitHandler(
	repo habit.Repository,
	validator *validator.Validator,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) DeleteHabitHandler {
	if repo == nil {
		panic("nil habit repository")
	}

	return decorator.ApplyCommandDecorators(
		deleteHabitHandler{
			repo:      repo,
			validator: validator,
		},
		log,
		metricsClient,
	)
}

func (h deleteHabitHandler) Handle(ctx context.Context, cmd DeleteHabit) error {
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

	// Delete the habit (with authorization check in repository)
	return h.repo.DeleteHabit(ctx, cmd.HabitID, cmd.UserID)
}
