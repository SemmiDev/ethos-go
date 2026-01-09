package habit

import (
	"context"
	"time"
)

// HabitReader provides read-only access to habit data.
// Use this interface when you only need to query habits.
type HabitReader interface {
	// GetHabit retrieves a habit by ID for a specific user.
	GetHabit(ctx context.Context, habitID, userID string) (*Habit, error)

	// ListHabitsByUser returns all habits for a user.
	ListHabitsByUser(ctx context.Context, userID string) ([]*Habit, error)
}

// HabitWriter provides write operations for habit data.
// Use this interface when you need to create or modify habits.
type HabitWriter interface {
	// AddHabit creates a new habit.
	AddHabit(ctx context.Context, habit *Habit) error

	// UpdateHabit modifies an existing habit using a callback function.
	UpdateHabit(
		ctx context.Context,
		habitID, userID string,
		updateFn func(ctx context.Context, h *Habit) (*Habit, error),
	) error

	// DeleteHabit removes a habit.
	DeleteHabit(ctx context.Context, habitID, userID string) error
}

// StatsRepository provides operations for habit statistics.
type StatsRepository interface {
	// GetStats retrieves habit statistics.
	GetStats(ctx context.Context, habitID string) (*HabitStats, error)

	// UpsertStats creates or updates habit statistics.
	UpsertStats(ctx context.Context, stats *HabitStats) error
}

// VacationRepository provides operations for habit vacations.
type VacationRepository interface {
	// AddVacation creates a new habit vacation.
	AddVacation(ctx context.Context, vacation *HabitVacation) error

	// GetActiveVacation retrieves the currently active vacation for a habit.
	GetActiveVacation(ctx context.Context, habitID string) (*HabitVacation, error)

	// EndVacation ends a vacation at the specified date.
	EndVacation(ctx context.Context, vacationID string, endDate time.Time) error

	// ListVacations returns all vacations for a habit.
	ListVacations(ctx context.Context, habitID string) ([]*HabitVacation, error)
}

// Repository combines all habit repository interfaces.
// This is the full interface that adapters implement.
// Consumers should depend on the smallest interface they need.
type Repository interface {
	HabitReader
	HabitWriter
	StatsRepository
	VacationRepository
}
