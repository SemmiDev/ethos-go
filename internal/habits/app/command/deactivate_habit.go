package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
)

// DeactivateHabit command encapsulates habit deactivation
type DeactivateHabit struct {
	HabitID string `validate:"uuid"`
	UserID  string `validate:"uuid"`
}

// DeactivateHabitHandler processes habit deactivation commands
type DeactivateHabitHandler decorator.CommandHandler[DeactivateHabit]

type deactivateHabitHandler struct {
	repo      habit.Repository
	validator *validator.Validator
}

// NewDeactivateHabitHandler creates a new handler with decorators
func NewDeactivateHabitHandler(
	repo habit.Repository,
	validator *validator.Validator,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) DeactivateHabitHandler {
	if repo == nil {
		panic("nil habit repository")
	}

	return decorator.ApplyCommandDecorators[DeactivateHabit](
		deactivateHabitHandler{
			repo:      repo,
			validator: validator,
		},
		log,
		metricsClient,
	)
}

func (h deactivateHabitHandler) Handle(ctx context.Context, cmd DeactivateHabit) error {
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

	// Use repository UpdateFn pattern
	return h.repo.UpdateHabit(
		ctx,
		cmd.HabitID,
		cmd.UserID,
		func(ctx context.Context, habit *habit.Habit) (*habit.Habit, error) {
			// Apply domain behavior
			if err := habit.Deactivate(); err != nil {
				return nil, err
			}
			return habit, nil
		},
	)
}
