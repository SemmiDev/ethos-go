package command

import (
	"context"
	"time"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
	"github.com/semmidev/ethos-go/internal/habits/domain/habitlog"
)

// LogHabit command encapsulates logging a habit completion
type LogHabit struct {
	LogID   string
	HabitID string
	UserID  string
	LogDate time.Time `json:"log_date" validate:"required"`
	Count   int       `json:"count" validate:"required,min=1"`
	Note    *string   `json:"note"`
}

// LogHabitHandler processes habit logging commands
type LogHabitHandler decorator.CommandHandler[LogHabit]

type logHabitHandler struct {
	habitRepo habit.Repository
	logRepo   habitlog.Repository
	validator *validator.Validator
	streakSvc *habit.StreakService
}

// NewLogHabitHandler creates a new handler with decorators
func NewLogHabitHandler(
	habitRepo habit.Repository,
	logRepo habitlog.Repository,
	validator *validator.Validator,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) LogHabitHandler {
	if habitRepo == nil {
		panic("nil habit repository")
	}
	if logRepo == nil {
		panic("nil habit log repository")
	}

	return decorator.ApplyCommandDecorators[LogHabit](
		logHabitHandler{
			habitRepo: habitRepo,
			logRepo:   logRepo,
			validator: validator,
			streakSvc: habit.NewStreakService(),
		},
		log,
		metricsClient,
	)
}

func (h logHabitHandler) Handle(ctx context.Context, cmd LogHabit) error {
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

	// Verify habit exists and belongs to user
	_, err := h.habitRepo.GetHabit(ctx, cmd.HabitID, cmd.UserID)
	if err != nil {
		return err
	}

	// Create new log entry
	newLog, err := habitlog.NewHabitLog(
		cmd.LogID,
		cmd.HabitID,
		cmd.UserID,
		cmd.LogDate,
		cmd.Count,
		cmd.Note,
	)
	if err != nil {
		return err
	}

	if err := h.logRepo.AddHabitLog(ctx, newLog); err != nil {
		return err
	}

	// Recalculate streaks and stats
	// 1. Get habit aggregate (already fetched)
	habitAgg, err := h.habitRepo.GetHabit(ctx, cmd.HabitID, cmd.UserID)
	if err != nil {
		return err // Should not happen as we checked above
	}

	// 2. Fetch all logs for this habit (needed for accurate streak calc)
	logs, err := h.logRepo.ListHabitLogs(ctx, cmd.HabitID, cmd.UserID)
	if err != nil {
		return err
	}

	// 3. Fetch active vacations
	vacations, err := h.habitRepo.ListVacations(ctx, cmd.HabitID)
	if err != nil {
		return err
	}

	// 4. Calculate stats
	stats := h.streakSvc.CalculateStreak(habitAgg, logs, vacations, time.Now())

	// 5. Persist stats
	if err := h.habitRepo.UpsertStats(ctx, stats); err != nil {
		return err
	}

	return nil
}
