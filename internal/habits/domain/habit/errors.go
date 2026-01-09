package habit

import "errors"

// Domain errors - pure domain errors without infrastructure dependencies
// These errors are translated to apperror.AppError at the adapter/port layer
var (
	// Business logic errors
	ErrAlreadyActive   = errors.New("habit is already active")
	ErrAlreadyInactive = errors.New("habit is already inactive")

	// Validation errors
	ErrEmptyName          = errors.New("habit name cannot be empty")
	ErrInvalidTargetCount = errors.New("target count must be positive")
	ErrInvalidReminder    = errors.New("invalid reminder time format (HH:MM)")
	ErrEmptyHabitID       = errors.New("empty habit id")
	ErrEmptyUserID        = errors.New("empty user id")

	// Access errors
	ErrNotFound     = errors.New("habit not found")
	ErrUnauthorized = errors.New("user cannot access this habit")
)
