package events

import (
	"time"

	commonevents "github.com/semmidev/ethos-go/internal/common/events"
)

// Event subjects
const (
	HabitCreatedType     = "habits.habit.created"
	HabitCompletedType   = "habits.habit.completed"
	HabitDeactivatedType = "habits.habit.deactivated"
	HabitActivatedType   = "habits.habit.activated"
	StreakMilestoneType  = "habits.streak.milestone"
)

// HabitCreated is emitted when a new habit is created
type HabitCreated struct {
	commonevents.BaseEvent
	HabitID     string `json:"habit_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Frequency   string `json:"frequency"`
	TargetCount int    `json:"target_count"`
}

// NewHabitCreated creates a new HabitCreated event
func NewHabitCreated(habitID, userID, name, frequency string, targetCount int) HabitCreated {
	return HabitCreated{
		BaseEvent:   commonevents.NewBaseEvent(HabitCreatedType, "habit", habitID),
		HabitID:     habitID,
		UserID:      userID,
		Name:        name,
		Frequency:   frequency,
		TargetCount: targetCount,
	}
}

// HabitCompleted is emitted when a habit is logged/completed
type HabitCompleted struct {
	commonevents.BaseEvent
	HabitID    string    `json:"habit_id"`
	UserID     string    `json:"user_id"`
	LogID      string    `json:"log_id"`
	LogDate    time.Time `json:"log_date"`
	Count      int       `json:"count"`
	TotalToday int       `json:"total_today"`
}

// NewHabitCompleted creates a new HabitCompleted event
func NewHabitCompleted(habitID, userID, logID string, logDate time.Time, count, totalToday int) HabitCompleted {
	return HabitCompleted{
		BaseEvent:  commonevents.NewBaseEvent(HabitCompletedType, "habit", habitID),
		HabitID:    habitID,
		UserID:     userID,
		LogID:      logID,
		LogDate:    logDate,
		Count:      count,
		TotalToday: totalToday,
	}
}

// StreakMilestone is emitted when a user reaches a streak milestone
type StreakMilestone struct {
	commonevents.BaseEvent
	HabitID       string `json:"habit_id"`
	UserID        string `json:"user_id"`
	HabitName     string `json:"habit_name"`
	CurrentStreak int    `json:"current_streak"`
	Milestone     int    `json:"milestone"` // 7, 30, 100, etc.
}

// NewStreakMilestone creates a new StreakMilestone event
func NewStreakMilestone(habitID, userID, habitName string, currentStreak, milestone int) StreakMilestone {
	return StreakMilestone{
		BaseEvent:     commonevents.NewBaseEvent(StreakMilestoneType, "habit", habitID),
		HabitID:       habitID,
		UserID:        userID,
		HabitName:     habitName,
		CurrentStreak: currentStreak,
		Milestone:     milestone,
	}
}

// HabitDeactivated is emitted when a habit is deactivated
type HabitDeactivated struct {
	commonevents.BaseEvent
	HabitID string `json:"habit_id"`
	UserID  string `json:"user_id"`
}

// NewHabitDeactivated creates a new HabitDeactivated event
func NewHabitDeactivated(habitID, userID string) HabitDeactivated {
	return HabitDeactivated{
		BaseEvent: commonevents.NewBaseEvent(HabitDeactivatedType, "habit", habitID),
		HabitID:   habitID,
		UserID:    userID,
	}
}

// HabitActivated is emitted when a habit is activated
type HabitActivated struct {
	commonevents.BaseEvent
	HabitID string `json:"habit_id"`
	UserID  string `json:"user_id"`
}

// NewHabitActivated creates a new HabitActivated event
func NewHabitActivated(habitID, userID string) HabitActivated {
	return HabitActivated{
		BaseEvent: commonevents.NewBaseEvent(HabitActivatedType, "habit", habitID),
		HabitID:   habitID,
		UserID:    userID,
	}
}
