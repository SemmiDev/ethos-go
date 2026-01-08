package query

import (
	"context"
	"time"
)

// ExportDataRepository abstracts data fetching for GDPR export
// This keeps the query handler clean and testable
type ExportDataRepository interface {
	GetUserHabits(ctx context.Context, userID string) ([]ExportedHabit, error)
	GetUserHabitLogs(ctx context.Context, userID string) ([]ExportedHabitLog, error)
	GetUserNotifications(ctx context.Context, userID string) ([]ExportedNotif, error)
}

// ExportedHabit represents a habit for GDPR export
type ExportedHabit struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Frequency    string    `json:"frequency"`
	TargetCount  int       `json:"target_count"`
	IsActive     bool      `json:"is_active"`
	ReminderTime *string   `json:"reminder_time"`
	CreatedAt    time.Time `json:"created_at"`
}

// ExportedHabitLog represents a habit log for GDPR export
type ExportedHabitLog struct {
	ID        string    `json:"id"`
	HabitID   string    `json:"habit_id"`
	LogDate   string    `json:"log_date"`
	Count     int       `json:"count"`
	Note      *string   `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

// ExportedNotif represents a notification for GDPR export
type ExportedNotif struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Data      []byte    `json:"data"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
