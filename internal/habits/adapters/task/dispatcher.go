package task

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/internal/common/logger"
	domaintask "github.com/semmidev/ethos-go/internal/habits/domain/task"
)

// Task constants for habits module
const TaskHabitCreated = "habits:created"

// HabitCreatedPayload contains data for the habit created task
type HabitCreatedPayload struct {
	HabitID string `json:"habit_id"`
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
}

// AsynqTaskDispatcher dispatches habit-related tasks to the queue
type AsynqTaskDispatcher struct {
	client *asynq.Client
	logger logger.Logger
}

// Ensure AsynqTaskDispatcher implements domaintask.TaskDispatcher
var _ domaintask.TaskDispatcher = (*AsynqTaskDispatcher)(nil)

func NewAsynqTaskDispatcher(client *asynq.Client, logger logger.Logger) *AsynqTaskDispatcher {
	return &AsynqTaskDispatcher{client: client, logger: logger}
}

func (d *AsynqTaskDispatcher) DispatchHabitCreated(ctx context.Context, habitID, userID, name string) error {
	payload, err := json.Marshal(HabitCreatedPayload{HabitID: habitID, UserID: userID, Name: name})
	if err != nil {
		return err
	}

	task := asynq.NewTask(TaskHabitCreated, payload)
	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		d.logger.Error(ctx, err, "failed to enqueue habit created task")
		return err
	}

	d.logger.Info(ctx, "dispatched habit created task", logger.Field{Key: "task_id", Value: info.ID})
	return nil
}
