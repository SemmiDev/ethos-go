package habitlog

import (
	"context"
	"time"
)

// Repository defines the interface for habit log persistence
type Repository interface {
	// AddHabitLog creates a new habit log entry
	AddHabitLog(ctx context.Context, log *HabitLog) error

	// GetHabitLog retrieves a single log by ID with authorization
	GetHabitLog(ctx context.Context, logID, userID string) (*HabitLog, error)

	// UpdateHabitLog uses the updateFn pattern for transactional updates
	UpdateHabitLog(
		ctx context.Context,
		logID, userID string,
		updateFn func(ctx context.Context, log *HabitLog) (*HabitLog, error),
	) error

	// DeleteHabitLog removes a log entry with authorization
	DeleteHabitLog(ctx context.Context, logID, userID string) error

	// GetHabitLogByDate finds a log for a specific habit on a specific date
	GetHabitLogByDate(ctx context.Context, habitID string, date time.Time, userID string) (*HabitLog, error)

	// ListHabitLogs retrieves all logs for a habit (used for streak calculation)
	ListHabitLogs(ctx context.Context, habitID, userID string) ([]*HabitLog, error)
}
