package command

import (
	"context"
	"time"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	"github.com/semmidev/ethos-go/internal/habits/domain/habitlog"
)

// UpdateHabitLog command encapsulates habit log updates
type UpdateHabitLog struct {
	LogID   string
	UserID  string
	Count   *int       `json:"count" validate:"omitempty,min=1"`
	Note    *string    `json:"note"`
	LogDate *time.Time `json:"log_date"`
}

// UpdateHabitLogHandler processes habit log update commands
type UpdateHabitLogHandler decorator.CommandHandler[UpdateHabitLog]

type updateHabitLogHandler struct {
	repo      habitlog.Repository
	validator *validator.Validator
}

// NewUpdateHabitLogHandler creates a new handler with decorators
func NewUpdateHabitLogHandler(
	repo habitlog.Repository,
	validator *validator.Validator,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) UpdateHabitLogHandler {
	if repo == nil {
		panic("nil habit log repository")
	}

	return decorator.ApplyCommandDecorators(
		updateHabitLogHandler{
			repo:      repo,
			validator: validator,
		},
		log,
		metricsClient,
	)
}

func (h updateHabitLogHandler) Handle(ctx context.Context, cmd UpdateHabitLog) error {
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

	return h.repo.UpdateHabitLog(
		ctx,
		cmd.LogID,
		cmd.UserID,
		func(ctx context.Context, log *habitlog.HabitLog) (*habitlog.HabitLog, error) {
			if cmd.Count != nil {
				if err := log.UpdateCount(*cmd.Count); err != nil {
					return nil, err
				}
			}
			if cmd.Note != nil {
				log.UpdateNote(cmd.Note)
			}
			if cmd.LogDate != nil {
				if err := log.UpdateLogDate(*cmd.LogDate); err != nil {
					return nil, err
				}
			}
			return log, nil
		},
	)
}
