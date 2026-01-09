package habitlog

import (
	"errors"
	"time"
)

// HabitLog represents a record of habit completion on a specific date
type HabitLog struct {
	logID     string
	habitID   string
	userID    string
	logDate   time.Time
	count     int
	note      *string // Nullable field - nil represents NULL in database
	createdAt time.Time
	updatedAt time.Time
}

// Domain errors - pure domain errors without infrastructure dependencies
var (
	ErrEmptyLogID   = errors.New("empty log id")
	ErrEmptyHabitID = errors.New("empty habit id")
	ErrEmptyUserID  = errors.New("empty user id")
	ErrInvalidCount = errors.New("count must be positive")
	ErrInvalidDate  = errors.New("invalid log date")
	ErrNotFound     = errors.New("habit log not found")
	ErrUnauthorized = errors.New("user cannot access this log")
)

// NewHabitLog creates a new habit log entry with validation
func NewHabitLog(
	logID, habitID, userID string,
	logDate time.Time,
	count int,
	note *string,
) (*HabitLog, error) {
	if logID == "" {
		return nil, ErrEmptyLogID
	}
	if habitID == "" {
		return nil, ErrEmptyHabitID
	}
	if userID == "" {
		return nil, ErrEmptyUserID
	}
	if count < 1 {
		return nil, ErrInvalidCount
	}
	if logDate.IsZero() {
		return nil, ErrInvalidDate
	}

	now := time.Now()
	return &HabitLog{
		logID:     logID,
		habitID:   habitID,
		userID:    userID,
		logDate:   logDate,
		count:     count,
		note:      note,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// UnmarshalHabitLogFromDatabase reconstructs a HabitLog from database
func UnmarshalHabitLogFromDatabase(
	logID, habitID, userID string,
	logDate time.Time,
	count int,
	note *string,
	createdAt, updatedAt time.Time,
) (*HabitLog, error) {
	if logID == "" {
		return nil, ErrEmptyLogID
	}
	if habitID == "" {
		return nil, ErrEmptyHabitID
	}
	if userID == "" {
		return nil, ErrEmptyUserID
	}

	return &HabitLog{
		logID:     logID,
		habitID:   habitID,
		userID:    userID,
		logDate:   logDate,
		count:     count,
		note:      note,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// Getters (read-only access)
func (l *HabitLog) LogID() string        { return l.logID }
func (l *HabitLog) HabitID() string      { return l.habitID }
func (l *HabitLog) UserID() string       { return l.userID }
func (l *HabitLog) LogDate() time.Time   { return l.logDate }
func (l *HabitLog) Count() int           { return l.count }
func (l *HabitLog) Note() *string        { return l.note }
func (l *HabitLog) CreatedAt() time.Time { return l.createdAt }
func (l *HabitLog) UpdatedAt() time.Time { return l.updatedAt }

// UpdateCount modifies the count for this log entry
func (l *HabitLog) UpdateCount(newCount int) error {
	if newCount < 1 {
		return ErrInvalidCount
	}
	l.count = newCount
	l.updatedAt = time.Now()
	return nil
}

// UpdateLogDate modifies the date for this log entry
func (l *HabitLog) UpdateLogDate(newDate time.Time) error {
	if newDate.IsZero() {
		return ErrInvalidDate
	}
	l.logDate = newDate
	l.updatedAt = time.Now()
	return nil
}

// UpdateNote modifies the note for this log entry
func (l *HabitLog) UpdateNote(newNote *string) {
	l.note = newNote
	l.updatedAt = time.Now()
}

// CanBeViewedBy checks if the user has permission to view this log
func (l *HabitLog) CanBeViewedBy(userID string) error {
	if l.userID != userID {
		return ErrUnauthorized
	}
	return nil
}

// CanBeModifiedBy checks if the user has permission to modify this log
func (l *HabitLog) CanBeModifiedBy(userID string) error {
	if l.userID != userID {
		return ErrUnauthorized
	}
	return nil
}
