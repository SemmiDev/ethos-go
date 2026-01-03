package query

import "time"

// Habit represents a read model for habit queries (optimized for UI)
type Habit struct {
	HabitID      string    `json:"habit_id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description,omitempty"` // Nullable field
	Frequency    string    `json:"frequency"`
	TargetCount  int       `json:"target_count"`
	ReminderTime *string   `json:"reminder_time,omitempty"` // Nullable field
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// HabitLog represents a read model for habit log queries
type HabitLog struct {
	LogID     string    `json:"log_id"`
	HabitID   string    `json:"habit_id"`
	UserID    string    `json:"user_id"`
	LogDate   time.Time `json:"log_date"`
	Count     int       `json:"count"`
	Note      *string   `json:"note,omitempty"` // Nullable field
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HabitStats represents statistics for a habit
type HabitStats struct {
	HabitID          string     `json:"habit_id"`
	HabitName        string     `json:"habit_name"`
	CurrentStreak    int        `json:"current_streak"`
	LongestStreak    int        `json:"longest_streak"`
	TotalCompletions int        `json:"total_completions"`
	CompletionRate   float64    `json:"completion_rate"` // Percentage
	ThisWeekCount    int        `json:"this_week_count"`
	ThisMonthCount   int        `json:"this_month_count"`
	LastLogDate      *time.Time `json:"last_log_date,omitempty"`
}

// DashboardSummary represents overall user statistics
type DashboardSummary struct {
	TotalActiveHabits     int          `json:"total_active_habits"`
	TotalCompletionsToday int          `json:"total_completions_today"`
	TotalCompletionsWeek  int          `json:"total_completions_week"`
	TotalCompletionsMonth int          `json:"total_completions_month"`
	BestStreak            int          `json:"best_streak"`
	CurrentStreak         int          `json:"current_streak"`
	LongestStreak         int          `json:"longest_streak"`
	WeeklyCompletion      int          `json:"weekly_completion"` // Percentage 0-100
	TotalLogs             int          `json:"total_logs"`
	HabitSummaries        []HabitStats `json:"habit_summaries"`
}

// WeeklyAnalytics represents weekly analytics data
type WeeklyAnalytics struct {
	Days              []DailyAnalytics `json:"days"`
	AverageCompletion int              `json:"average_completion"`
}

// DailyAnalytics represents analytics for a single day
type DailyAnalytics struct {
	DayName              string `json:"day_name"`
	Date                 string `json:"date"`
	LogsCount            int    `json:"logs_count"`
	CompletionPercentage int    `json:"completion_percentage"`
}

// ReminderHabit represents a habit that needs a reminder (due today, not completed)
type ReminderHabit struct {
	UserID       string  `db:"user_id"`
	HabitID      string  `db:"habit_id"`
	HabitName    string  `db:"name"`
	ReminderTime *string `db:"reminder_time"`
}
