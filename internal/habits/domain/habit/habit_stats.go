package habit

import (
	"time"
)

// HabitStats represents pre-calculated statistics for a habit
type HabitStats struct {
	habitID          string
	currentStreak    int
	longestStreak    int
	totalCompletions int
	lastCompletedAt  *time.Time
	consistencyScore float64 // 0.0 to 100.0
	updatedAt        time.Time
}

// NewHabitStats creates a new HabitStats instance
func NewHabitStats(habitID string) *HabitStats {
	return &HabitStats{
		habitID:          habitID,
		currentStreak:    0,
		longestStreak:    0,
		totalCompletions: 0,
		lastCompletedAt:  nil,
		consistencyScore: 0.0,
		updatedAt:        time.Now(),
	}
}

// UnmarshalStatsFromDatabase reconstructs HabitStats from database
func UnmarshalStatsFromDatabase(
	habitID string,
	currentStreak, longestStreak, totalCompletions int,
	lastCompletedAt *time.Time,
	consistencyScore float64,
	updatedAt time.Time,
) *HabitStats {
	return &HabitStats{
		habitID:          habitID,
		currentStreak:    currentStreak,
		longestStreak:    longestStreak,
		totalCompletions: totalCompletions,
		lastCompletedAt:  lastCompletedAt,
		consistencyScore: consistencyScore,
		updatedAt:        updatedAt,
	}
}

// Getters
func (s HabitStats) HabitID() string             { return s.habitID }
func (s HabitStats) CurrentStreak() int          { return s.currentStreak }
func (s HabitStats) LongestStreak() int          { return s.longestStreak }
func (s HabitStats) TotalCompletions() int       { return s.totalCompletions }
func (s HabitStats) LastCompletedAt() *time.Time { return s.lastCompletedAt }
func (s HabitStats) ConsistencyScore() float64   { return s.consistencyScore }
func (s HabitStats) UpdatedAt() time.Time        { return s.updatedAt }

// UpdateStreak updates the current streak and potentially the longest streak
func (s *HabitStats) UpdateStreak(newStreak int, completedAt time.Time) {
	s.currentStreak = newStreak
	if newStreak > s.longestStreak {
		s.longestStreak = newStreak
	}
	s.lastCompletedAt = &completedAt
	s.totalCompletions++
	s.updatedAt = time.Now()
}

// ResetStreak resets the current streak to zero (streak broken)
func (s *HabitStats) ResetStreak() {
	s.currentStreak = 0
	s.updatedAt = time.Now()
}

// UpdateConsistency updates the consistency score
func (s *HabitStats) UpdateConsistency(score float64) {
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}
	s.consistencyScore = score
	s.updatedAt = time.Now()
}
