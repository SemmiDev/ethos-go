package habit

import (
	"errors"
	"time"
)

// HabitVacation represents a period where a habit is paused
type HabitVacation struct {
	id        string
	habitID   string
	startDate time.Time
	endDate   *time.Time // nil = ongoing
	reason    *string
	createdAt time.Time
}

// Vacation domain errors - pure domain errors without infrastructure dependencies
var (
	ErrVacationEmptyID      = errors.New("vacation id cannot be empty")
	ErrVacationEmptyHabitID = errors.New("habit id cannot be empty")
	ErrVacationInvalidDates = errors.New("end date must be after start date")
	ErrVacationNotFound     = errors.New("vacation not found")
	ErrVacationAlreadyEnded = errors.New("vacation has already ended")
	ErrVacationOverlap      = errors.New("vacation dates overlap with an existing vacation")
)

// NewHabitVacation creates a new vacation period for a habit
func NewHabitVacation(id, habitID string, startDate time.Time, reason *string) (*HabitVacation, error) {
	if id == "" {
		return nil, ErrVacationEmptyID
	}
	if habitID == "" {
		return nil, ErrVacationEmptyHabitID
	}

	return &HabitVacation{
		id:        id,
		habitID:   habitID,
		startDate: startDate,
		endDate:   nil, // Open-ended until explicitly ended
		reason:    reason,
		createdAt: time.Now(),
	}, nil
}

// UnmarshalVacationFromDatabase reconstructs a HabitVacation from database
func UnmarshalVacationFromDatabase(
	id, habitID string,
	startDate time.Time,
	endDate *time.Time,
	reason *string,
	createdAt time.Time,
) *HabitVacation {
	return &HabitVacation{
		id:        id,
		habitID:   habitID,
		startDate: startDate,
		endDate:   endDate,
		reason:    reason,
		createdAt: createdAt,
	}
}

// Getters
func (v HabitVacation) ID() string           { return v.id }
func (v HabitVacation) HabitID() string      { return v.habitID }
func (v HabitVacation) StartDate() time.Time { return v.startDate }
func (v HabitVacation) EndDate() *time.Time  { return v.endDate }
func (v HabitVacation) Reason() *string      { return v.reason }
func (v HabitVacation) CreatedAt() time.Time { return v.createdAt }

// IsOngoing returns true if the vacation has not ended
func (v HabitVacation) IsOngoing() bool {
	return v.endDate == nil
}

// IsActiveOn returns true if the given date falls within the vacation period
func (v HabitVacation) IsActiveOn(date time.Time) bool {
	// Normalize to date only (ignore time component)
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	startOnly := time.Date(v.startDate.Year(), v.startDate.Month(), v.startDate.Day(), 0, 0, 0, 0, v.startDate.Location())

	if dateOnly.Before(startOnly) {
		return false
	}

	if v.endDate == nil {
		return true // Ongoing vacation
	}

	endOnly := time.Date(v.endDate.Year(), v.endDate.Month(), v.endDate.Day(), 0, 0, 0, 0, v.endDate.Location())
	return !dateOnly.After(endOnly)
}

// End ends the vacation on the given date
func (v *HabitVacation) End(endDate time.Time) error {
	if v.endDate != nil {
		return ErrVacationAlreadyEnded
	}

	// Normalize dates for comparison
	endOnly := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())
	startOnly := time.Date(v.startDate.Year(), v.startDate.Month(), v.startDate.Day(), 0, 0, 0, 0, v.startDate.Location())

	if endOnly.Before(startOnly) {
		return ErrVacationInvalidDates
	}

	v.endDate = &endDate
	return nil
}
