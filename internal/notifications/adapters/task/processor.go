package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/internal/common/logger"
	habittask "github.com/semmidev/ethos-go/internal/habits/adapters/task"
	habitsapp "github.com/semmidev/ethos-go/internal/habits/app"
	habitsquery "github.com/semmidev/ethos-go/internal/habits/app/query"
	notifapp "github.com/semmidev/ethos-go/internal/notifications/app"
	"github.com/semmidev/ethos-go/internal/notifications/app/command"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

const (
	TaskProcessReminders = "notifications:process_reminders"
)

// TaskProcessor handles processing of notification-related background tasks
type TaskProcessor struct {
	notifApp  notifapp.Application
	habitsApp habitsapp.Application
	logger    logger.Logger
}

func NewTaskProcessor(
	notifApp notifapp.Application,
	habitsApp habitsapp.Application,
	logger logger.Logger,
) *TaskProcessor {
	return &TaskProcessor{
		notifApp:  notifApp,
		habitsApp: habitsApp,
		logger:    logger,
	}
}

// NewProcessRemindersTask creates a task to process reminders
func NewProcessRemindersTask() *asynq.Task {
	return asynq.NewTask(TaskProcessReminders, nil)
}

// ProcessTask implements asynq.Handler for reminders
func (p *TaskProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	currentTime := time.Now().Format("15:04")

	p.logger.Info(ctx, "processing habit reminders task",
		logger.Field{Key: "current_time", Value: currentTime},
	)

	// Get Habits Due - the query already filters correctly:
	// - At 8 PM: returns habits with matching time OR NULL reminder_time
	// - At other times: returns only habits with matching reminder_time
	habits, err := p.habitsApp.Queries.GetHabitsDue.Handle(ctx, habitsquery.GetHabitsDue{})
	if err != nil {
		p.logger.Error(ctx, err, "failed to get habits due")
		return err
	}

	count := 0
	for _, habit := range habits {
		title := "Habit Reminder"
		message := fmt.Sprintf("Don't forget to complete '%s' today!", habit.HabitName)

		err := p.notifApp.Commands.CreateNotification.Handle(ctx, command.CreateNotification{
			UserID:  habit.UserID,
			Type:    domain.TypeHabitReminder,
			Title:   title,
			Message: message,
			Data: map[string]interface{}{
				"habit_id": habit.HabitID,
			},
		})

		if err != nil {
			p.logger.Error(ctx, err, "failed to create notification", logger.Field{Key: "user_id", Value: habit.UserID})
			continue
		}
		count++
	}

	p.logger.Info(ctx, "processed reminders", logger.Field{Key: "count", Value: count})
	return nil
}

// ProcessHabitCreatedTask handles immediate notification creation when a habit is created
func (p *TaskProcessor) ProcessHabitCreatedTask(ctx context.Context, t *asynq.Task) error {
	p.logger.Info(ctx, "processing habit created task")

	var payload habittask.HabitCreatedPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to parse task payload: %w", err)
	}

	title := "New Habit Started!"
	message := fmt.Sprintf("You've started tracking '%s'. We believe in you!", payload.Name)

	err := p.notifApp.Commands.CreateNotification.Handle(ctx, command.CreateNotification{
		UserID:  payload.UserID,
		Type:    domain.TypeWelcome,
		Title:   title,
		Message: message,
		Data: map[string]interface{}{
			"habit_id": payload.HabitID,
		},
	})
	if err != nil {
		p.logger.Error(ctx, err, "failed to create welcome notification")
		return err
	}

	p.logger.Info(ctx, "sent welcome notification", logger.Field{Key: "user_id", Value: payload.UserID})
	return nil
}
