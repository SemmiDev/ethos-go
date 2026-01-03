package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
)

// UpdateHabit command encapsulates habit update input
type UpdateHabit struct {
	HabitID            string
	UserID             string
	Name               *string `json:"name" validate:"omitempty,min=3,max=100"`
	Description        *string `json:"description"` // Nullable
	Frequency          *string `json:"frequency" validate:"omitempty,oneof=daily weekly monthly custom"`
	RecurrenceDays     *int16  `json:"recurrence_days"`
	RecurrenceInterval *int    `json:"recurrence_interval"`
	TargetCount        *int    `json:"target_count" validate:"omitempty,min=1"`
	ReminderTime       *string `json:"reminder_time"` // Nullable - e.g. "08:00"
}

// UpdateHabitHandler processes habit update commands
type UpdateHabitHandler decorator.CommandHandler[UpdateHabit]

type updateHabitHandler struct {
	repo      habit.Repository
	validator *validator.Validator
}

// NewUpdateHabitHandler creates a new handler with decorators
func NewUpdateHabitHandler(
	repo habit.Repository,
	validator *validator.Validator,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) UpdateHabitHandler {
	if repo == nil {
		panic("nil habit repository")
	}

	return decorator.ApplyCommandDecorators[UpdateHabit](
		updateHabitHandler{
			repo:      repo,
			validator: validator,
		},
		log,
		metricsClient,
	)
}

func (h updateHabitHandler) Handle(ctx context.Context, cmd UpdateHabit) error {
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

	// Use repository UpdateFn pattern for transactional update
	return h.repo.UpdateHabit(
		ctx,
		cmd.HabitID,
		cmd.UserID,
		func(ctx context.Context, h *habit.Habit) (*habit.Habit, error) {
			// Apply updates if provided
			if cmd.Name != nil || cmd.Description != nil || cmd.Frequency != nil || cmd.RecurrenceDays != nil || cmd.RecurrenceInterval != nil || cmd.TargetCount != nil || cmd.ReminderTime != nil {
				// Resolve Frequency
				var freq habit.Frequency
				var err error
				if cmd.Frequency != nil {
					freq, err = habit.NewFrequency(*cmd.Frequency)
					if err != nil {
						return nil, err
					}
				} else {
					freq = h.Frequency()
				}

				// Resolve Recurrence
				currentRecurrence := h.Recurrence()
				days := currentRecurrence.Days()
				interval := currentRecurrence.Interval()

				if cmd.RecurrenceDays != nil {
					days = *cmd.RecurrenceDays
				}
				if cmd.RecurrenceInterval != nil && *cmd.RecurrenceInterval > 0 {
					interval = *cmd.RecurrenceInterval
				}

				recurrence, err := habit.NewRecurrence(days, interval)
				if err != nil {
					return nil, err
				}

				name := h.Name()
				if cmd.Name != nil {
					name = *cmd.Name
				}

				description := h.Description()
				if cmd.Description != nil {
					description = cmd.Description
				}

				targetCount := h.TargetCount()
				if cmd.TargetCount != nil {
					targetCount = *cmd.TargetCount
				}

				reminderTime := h.ReminderTime()
				if cmd.ReminderTime != nil {
					reminderTime = cmd.ReminderTime
				}

				if err := h.Update(name, description, freq, recurrence, targetCount, reminderTime); err != nil {
					return nil, err
				}
			}

			return h, nil
		},
	)
}
