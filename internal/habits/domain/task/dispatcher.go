package task

import "context"

// TaskDispatcher interface for dispatching habit-related background tasks
type TaskDispatcher interface {
	DispatchHabitCreated(ctx context.Context, habitID, userID, name string) error
}
