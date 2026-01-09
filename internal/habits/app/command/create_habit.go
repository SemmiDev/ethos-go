package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	habitevents "github.com/semmidev/ethos-go/internal/habits/domain/events"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
	domaintask "github.com/semmidev/ethos-go/internal/habits/domain/task"
)

// CreateHabit command encapsulates habit creation input
type CreateHabit struct {
	HabitID            string
	UserID             string
	Name               string  `json:"name" validate:"required,min=3,max=100"`
	Description        *string `json:"description"`
	Frequency          string  `json:"frequency" validate:"required,oneof=daily weekly monthly custom"`
	RecurrenceDays     *int16  `json:"recurrence_days"`     // Bitmask: Sun=1, Mon=2, etc. nil = all days
	RecurrenceInterval *int    `json:"recurrence_interval"` // Every N periods. nil = 1
	TargetCount        int     `json:"target_count" validate:"required,min=1"`
	ReminderTime       *string `json:"reminder_time"`
}

// CreateHabitHandler processes habit creation commands
type CreateHabitHandler decorator.CommandHandler[CreateHabit]

type createHabitHandler struct {
	repo       habit.Repository
	validator  *validator.Validator
	dispatcher domaintask.TaskDispatcher
	publisher  events.Publisher
}

// NewCreateHabitHandler creates a new handler with decorators
func NewCreateHabitHandler(
	repo habit.Repository,
	validator *validator.Validator,
	dispatcher domaintask.TaskDispatcher,
	publisher events.Publisher, // Injected publisher
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) CreateHabitHandler {
	if repo == nil {
		panic("nil habit repository")
	}

	return decorator.ApplyCommandDecorators(
		createHabitHandler{
			repo:       repo,
			validator:  validator,
			dispatcher: dispatcher,
			publisher:  publisher,
		},
		log,
		metricsClient,
	)
}

func (h createHabitHandler) Handle(ctx context.Context, cmd CreateHabit) error {
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

	// Create frequency value object
	frequency, err := habit.NewFrequency(cmd.Frequency)
	if err != nil {
		return err
	}

	// Create recurrence value object (use defaults if not provided)
	recurrenceDays := habit.AllDays
	if cmd.RecurrenceDays != nil {
		recurrenceDays = *cmd.RecurrenceDays
	}
	recurrenceInterval := 1
	if cmd.RecurrenceInterval != nil && *cmd.RecurrenceInterval > 0 {
		recurrenceInterval = *cmd.RecurrenceInterval
	}
	recurrence, err := habit.NewRecurrence(recurrenceDays, recurrenceInterval)
	if err != nil {
		recurrence = habit.DefaultRecurrence()
	}

	// Create new habit aggregate
	newHabit, err := habit.NewHabit(
		cmd.HabitID,
		cmd.UserID,
		cmd.Name,
		cmd.Description,
		frequency,
		recurrence,
		cmd.TargetCount,
		cmd.ReminderTime,
	)
	if err != nil {
		return err
	}

	// Persist the habit
	if err := h.repo.AddHabit(ctx, newHabit); err != nil {
		return err
	}

	// Notify about habit creation
	// We ignore error here to not fail the transaction if notification fails,
	// but in production, we should log it properly (h.notifier implementation typically logs).
	_ = h.dispatcher.DispatchHabitCreated(ctx, cmd.HabitID, cmd.UserID, cmd.Name)

	// Publish HabitCreated event
	event := habitevents.NewHabitCreated(
		newHabit.HabitID(),
		newHabit.UserID(),
		newHabit.Name(),
		newHabit.Frequency().String(),
		newHabit.TargetCount(),
	)
	_ = h.publisher.Publish(ctx, event)

	return nil
}
