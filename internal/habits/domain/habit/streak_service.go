package habit

import (
	"sort"
	"time"

	"github.com/semmidev/ethos-go/internal/habits/domain/habitlog"
)

// StreakService calculates streak statistics for a habit
type StreakService struct{}

// NewStreakService creates a new StreakService
func NewStreakService() *StreakService {
	return &StreakService{}
}

// CalculateStreak computes the current and longest streak for a habit based on logs and vacations
func (s *StreakService) CalculateStreak(
	habit *Habit,
	logs []*habitlog.HabitLog,
	vacations []*HabitVacation,
	today time.Time,
) *HabitStats {
	stats := NewHabitStats(habit.HabitID())

	if len(logs) == 0 {
		return stats
	}

	// Sort logs by date descending (most recent first)
	sortedLogs := make([]*habitlog.HabitLog, len(logs))
	copy(sortedLogs, logs)
	sort.Slice(sortedLogs, func(i, j int) bool {
		return sortedLogs[i].LogDate().After(sortedLogs[j].LogDate())
	})

	// Create a set of completion dates
	completionDates := make(map[string]bool)
	for _, log := range sortedLogs {
		dateKey := log.LogDate().Format("2006-01-02")
		completionDates[dateKey] = true
	}

	// Create a set of vacation dates
	isVacationDate := func(date time.Time) bool {
		for _, v := range vacations {
			if v.IsActiveOn(date) {
				return true
			}
		}
		return false
	}

	// Should complete on this date?
	shouldComplete := func(date time.Time) bool {
		return habit.Recurrence().ShouldCompleteOn(date, habit.Frequency(), habit.CreatedAt())
	}

	// Calculate current streak
	currentStreak := 0
	longestStreak := 0
	tempStreak := 0
	totalCompletions := len(completionDates)

	// Start from today and go backwards
	checkDate := today
	firstCheck := true

	for {
		dateKey := checkDate.Format("2006-01-02")

		// If it's a vacation day, skip it (don't break streak)
		if isVacationDate(checkDate) {
			checkDate = checkDate.AddDate(0, 0, -1)
			continue
		}

		// If this day requires completion
		if shouldComplete(checkDate) {
			if completionDates[dateKey] {
				// Completed - increase streak
				if firstCheck {
					currentStreak++
				}
				tempStreak++
			} else {
				// Not completed - streak broken
				if firstCheck {
					// Check if today is the first miss
					firstCheck = false
				}
				if tempStreak > longestStreak {
					longestStreak = tempStreak
				}
				tempStreak = 0
				break
			}
		}

		firstCheck = false

		// Don't go before habit creation
		if checkDate.Before(habit.CreatedAt()) {
			break
		}

		// Move to previous day
		checkDate = checkDate.AddDate(0, 0, -1)
	}

	// Final longest streak check
	if tempStreak > longestStreak {
		longestStreak = tempStreak
	}

	// Calculate consistency score (last 30 days)
	consistencyScore := s.CalculateConsistency(habit, completionDates, vacations, today, 30)

	// Find last completed date
	var lastCompletedAt *time.Time
	if len(sortedLogs) > 0 {
		lastDate := sortedLogs[0].LogDate()
		lastCompletedAt = &lastDate
	}

	// Update stats
	stats = UnmarshalStatsFromDatabase(
		habit.HabitID(),
		currentStreak,
		longestStreak,
		totalCompletions,
		lastCompletedAt,
		consistencyScore,
		time.Now(),
	)

	return stats
}

// CalculateConsistency computes the consistency percentage over a given number of days
func (s *StreakService) CalculateConsistency(
	habit *Habit,
	completionDates map[string]bool,
	vacations []*HabitVacation,
	today time.Time,
	days int,
) float64 {
	isVacationDate := func(date time.Time) bool {
		for _, v := range vacations {
			if v.IsActiveOn(date) {
				return true
			}
		}
		return false
	}

	shouldComplete := func(date time.Time) bool {
		return habit.Recurrence().ShouldCompleteOn(date, habit.Frequency(), habit.CreatedAt())
	}

	expectedDays := 0
	completedDays := 0

	for i := 0; i < days; i++ {
		checkDate := today.AddDate(0, 0, -i)

		// Skip dates before habit creation
		if checkDate.Before(habit.CreatedAt()) {
			continue
		}

		// Skip vacation days
		if isVacationDate(checkDate) {
			continue
		}

		// If this day required completion
		if shouldComplete(checkDate) {
			expectedDays++
			dateKey := checkDate.Format("2006-01-02")
			if completionDates[dateKey] {
				completedDays++
			}
		}
	}

	if expectedDays == 0 {
		return 100.0 // No expected days = 100% consistency
	}

	return float64(completedDays) / float64(expectedDays) * 100.0
}
