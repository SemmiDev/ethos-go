package habit

import (
	"context"
	"time"
)

type Repository interface {
	// Habit CRUD
	AddHabit(ctx context.Context, habit *Habit) error
	GetHabit(ctx context.Context, habitID, userID string) (*Habit, error)
	UpdateHabit(
		ctx context.Context,
		habitID, userID string,
		updateFn func(ctx context.Context, h *Habit) (*Habit, error),
	) error
	DeleteHabit(ctx context.Context, habitID, userID string) error
	ListHabitsByUser(ctx context.Context, userID string) ([]*Habit, error)

	// Habit Stats
	GetStats(ctx context.Context, habitID string) (*HabitStats, error)
	UpsertStats(ctx context.Context, stats *HabitStats) error

	// Habit Vacations
	AddVacation(ctx context.Context, vacation *HabitVacation) error
	GetActiveVacation(ctx context.Context, habitID string) (*HabitVacation, error)
	EndVacation(ctx context.Context, vacationID string, endDate time.Time) error
	ListVacations(ctx context.Context, habitID string) ([]*HabitVacation, error)
}
