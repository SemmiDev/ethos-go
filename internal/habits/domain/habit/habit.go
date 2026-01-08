package habit

import (
	"errors"
	"time"

	"github.com/semmidev/ethos-go/internal/common/apperror"
)

type Habit struct {
	habitID      string
	userID       string
	name         string
	description  *string // Nullable field - nil represents NULL in database
	frequency    Frequency
	recurrence   Recurrence // Advanced recurrence (days, interval)
	targetCount  int
	reminderTime *string // Nullable field - e.g. "08:00"
	isActive     bool
	createdAt    time.Time
	updatedAt    time.Time
}

var (
	ErrAlreadyActive      = errors.New("habit is already active")
	ErrAlreadyInactive    = errors.New("habit is already inactive")
	ErrEmptyName          = apperror.ValidationFailed("habit name cannot be empty")
	ErrInvalidTargetCount = apperror.ValidationFailed("target count must be positive")
	ErrInvalidReminder    = apperror.ValidationFailed("invalid reminder time format (HH:MM)")
	ErrNotFound           = apperror.NotFound("habit", "")
	ErrUnauthorized       = apperror.Unauthorized("user cannot access this habit")
)

func NewHabit(
	habitID, userID, name string,
	description *string,
	frequency Frequency,
	recurrence Recurrence,
	targetCount int,
	reminderTime *string,
) (*Habit, error) {
	if habitID == "" {
		return nil, errors.New("empty habit id")
	}
	if userID == "" {
		return nil, errors.New("empty user id")
	}
	if name == "" {
		return nil, ErrEmptyName
	}
	if targetCount < 1 {
		return nil, ErrInvalidTargetCount
	}
	if err := frequency.Validate(); err != nil {
		return nil, err
	}
	if reminderTime != nil {
		if _, err := time.Parse("15:04", *reminderTime); err != nil {
			return nil, ErrInvalidReminder
		}
	}

	now := time.Now()
	return &Habit{
		habitID:      habitID,
		userID:       userID,
		name:         name,
		description:  description,
		frequency:    frequency,
		recurrence:   recurrence,
		targetCount:  targetCount,
		reminderTime: reminderTime,
		isActive:     true,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func UnmarshalHabitFromDatabase(
	habitID, userID, name string,
	description *string,
	frequencyStr string,
	recurrenceDays int16,
	recurrenceInterval int,
	targetCount int,
	reminderTime *string,
	isActive bool,
	createdAt, updatedAt time.Time,
) (*Habit, error) {
	frequency, err := NewFrequency(frequencyStr)
	if err != nil {
		return nil, err
	}

	recurrence, err := NewRecurrence(recurrenceDays, recurrenceInterval)
	if err != nil {
		// Fallback to default recurrence if invalid
		recurrence = DefaultRecurrence()
	}

	h := &Habit{
		habitID:      habitID,
		userID:       userID,
		name:         name,
		description:  description,
		frequency:    frequency,
		recurrence:   recurrence,
		targetCount:  targetCount,
		reminderTime: reminderTime,
		isActive:     isActive,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}

	return h, nil
}

func (h Habit) HabitID() string        { return h.habitID }
func (h Habit) UserID() string         { return h.userID }
func (h Habit) Name() string           { return h.name }
func (h Habit) Description() *string   { return h.description }
func (h Habit) Frequency() Frequency   { return h.frequency }
func (h Habit) Recurrence() Recurrence { return h.recurrence }
func (h Habit) TargetCount() int       { return h.targetCount }
func (h Habit) ReminderTime() *string  { return h.reminderTime }
func (h Habit) IsActive() bool         { return h.isActive }
func (h Habit) CreatedAt() time.Time   { return h.createdAt }
func (h Habit) UpdatedAt() time.Time   { return h.updatedAt }

func (h Habit) CanBeViewedBy(userID string) error {
	if h.userID != userID {
		return ErrUnauthorized
	}
	return nil
}
