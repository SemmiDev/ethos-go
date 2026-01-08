package adapters

import (
	"context"
	"database/sql"
	"time"

	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/habits/app/query"
)

// StatsRepository handles statistics calculations
type StatsRepository struct {
	db database.DBTX
}

func NewStatsRepository(db database.DBTX) *StatsRepository {
	return &StatsRepository{db: db}
}

// GetHabitStats calculates statistics for a single habit
func (r *StatsRepository) GetHabitStats(ctx context.Context, habitID, userID string) (*query.HabitStats, error) {
	// Get habit info
	var habitName string
	err := r.db.GetContext(ctx, &habitName, `SELECT name FROM habits WHERE habit_id = $1 AND user_id = $2`, habitID, userID)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows
	}
	if err != nil {
		return nil, err
	}

	stats := &query.HabitStats{
		HabitID:   habitID,
		HabitName: habitName,
	}

	// Total completions
	err = r.db.GetContext(ctx, &stats.TotalCompletions,
		`SELECT COALESCE(SUM(count), 0) FROM habit_logs WHERE habit_id = $1`, habitID)
	if err != nil {
		return nil, err
	}

	// Last log date
	var lastDate sql.NullTime
	err = r.db.GetContext(ctx, &lastDate,
		`SELECT MAX(log_date) FROM habit_logs WHERE habit_id = $1`, habitID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if lastDate.Valid {
		stats.LastLogDate = &lastDate.Time
	}

	// Current streak and longest streak
	stats.CurrentStreak = r.calculateCurrentStreak(ctx, habitID)
	stats.LongestStreak = r.calculateLongestStreak(ctx, habitID)

	// This week count
	weekStart := startOfWeek(time.Now())
	err = r.db.GetContext(ctx, &stats.ThisWeekCount,
		`SELECT COALESCE(SUM(count), 0) FROM habit_logs WHERE habit_id = $1 AND log_date >= $2`,
		habitID, weekStart)
	if err != nil {
		return nil, err
	}

	// This month count
	monthStart := startOfMonth(time.Now())
	err = r.db.GetContext(ctx, &stats.ThisMonthCount,
		`SELECT COALESCE(SUM(count), 0) FROM habit_logs WHERE habit_id = $1 AND log_date >= $2`,
		habitID, monthStart)
	if err != nil {
		return nil, err
	}

	// Completion rate (last 30 days)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	var daysLogged int
	err = r.db.GetContext(ctx, &daysLogged,
		`SELECT COUNT(DISTINCT log_date) FROM habit_logs WHERE habit_id = $1 AND log_date >= $2`,
		habitID, thirtyDaysAgo)
	if err != nil {
		return nil, err
	}
	stats.CompletionRate = float64(daysLogged) / 30.0 * 100.0

	return stats, nil
}

// GetDashboard calculates dashboard summary for a user
func (r *StatsRepository) GetDashboard(ctx context.Context, userID string) (*query.DashboardSummary, error) {
	summary := &query.DashboardSummary{
		HabitSummaries: []query.HabitStats{},
	}

	// Total active habits
	err := r.db.GetContext(ctx, &summary.TotalActiveHabits,
		`SELECT COUNT(*) FROM habits WHERE user_id = $1 AND is_active = true`, userID)
	if err != nil {
		return nil, err
	}

	// Today's completions
	today := time.Now().Truncate(24 * time.Hour)
	err = r.db.GetContext(ctx, &summary.TotalCompletionsToday,
		`SELECT COALESCE(SUM(count), 0) FROM habit_logs WHERE user_id = $1 AND log_date = $2`,
		userID, today)
	if err != nil {
		return nil, err
	}

	// This week's completions
	weekStart := startOfWeek(time.Now())
	err = r.db.GetContext(ctx, &summary.TotalCompletionsWeek,
		`SELECT COALESCE(SUM(count), 0) FROM habit_logs WHERE user_id = $1 AND log_date >= $2`,
		userID, weekStart)
	if err != nil {
		return nil, err
	}

	// This month's completions
	monthStart := startOfMonth(time.Now())
	err = r.db.GetContext(ctx, &summary.TotalCompletionsMonth,
		`SELECT COALESCE(SUM(count), 0) FROM habit_logs WHERE user_id = $1 AND log_date >= $2`,
		userID, monthStart)
	if err != nil {
		return nil, err
	}

	// Total logs all time
	err = r.db.GetContext(ctx, &summary.TotalLogs,
		`SELECT COALESCE(COUNT(*), 0) FROM habit_logs WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	// Get habit summaries for all active habits
	var habitIDs []string
	err = r.db.SelectContext(ctx, &habitIDs,
		`SELECT habit_id FROM habits WHERE user_id = $1 AND is_active = true`, userID)
	if err != nil {
		return nil, err
	}

	bestStreak := 0
	maxCurrentStreak := 0
	maxLongestStreak := 0
	for _, habitID := range habitIDs {
		habitStats, err := r.GetHabitStats(ctx, habitID, userID)
		if err != nil {
			continue // Skip habits with errors
		}
		summary.HabitSummaries = append(summary.HabitSummaries, *habitStats)
		if habitStats.LongestStreak > bestStreak {
			bestStreak = habitStats.LongestStreak
		}
		if habitStats.CurrentStreak > maxCurrentStreak {
			maxCurrentStreak = habitStats.CurrentStreak
		}
		if habitStats.LongestStreak > maxLongestStreak {
			maxLongestStreak = habitStats.LongestStreak
		}
	}
	summary.BestStreak = bestStreak
	summary.CurrentStreak = maxCurrentStreak
	summary.LongestStreak = maxLongestStreak

	// Calculate weekly completion percentage
	// (days with at least one log in this week / 7) * 100
	var daysWithLogs int
	err = r.db.GetContext(ctx, &daysWithLogs,
		`SELECT COUNT(DISTINCT log_date) FROM habit_logs WHERE user_id = $1 AND log_date >= $2`,
		userID, weekStart)
	if err != nil {
		daysWithLogs = 0
	}
	summary.WeeklyCompletion = int(float64(daysWithLogs) / 7.0 * 100.0)

	return summary, nil
}

// Helper methods for streak calculation

func (r *StatsRepository) calculateCurrentStreak(ctx context.Context, habitID string) int {
	// Get all log dates in descending order
	var dates []time.Time
	err := r.db.SelectContext(ctx, &dates,
		`SELECT DISTINCT log_date FROM habit_logs WHERE habit_id = $1 ORDER BY log_date DESC LIMIT 365`,
		habitID)
	if err != nil || len(dates) == 0 {
		return 0
	}

	// Check if the most recent log is today or yesterday
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.AddDate(0, 0, -1)

	if !dates[0].Equal(today) && !dates[0].Equal(yesterday) {
		return 0 // Streak is broken
	}

	// Count consecutive days
	streak := 1
	for i := 1; i < len(dates); i++ {
		expectedDate := dates[i-1].AddDate(0, 0, -1)
		if dates[i].Equal(expectedDate) {
			streak++
		} else {
			break
		}
	}

	return streak
}

func (r *StatsRepository) calculateLongestStreak(ctx context.Context, habitID string) int {
	var dates []time.Time
	err := r.db.SelectContext(ctx, &dates,
		`SELECT DISTINCT log_date FROM habit_logs WHERE habit_id = $1 ORDER BY log_date ASC`,
		habitID)
	if err != nil || len(dates) == 0 {
		return 0
	}

	maxStreak := 1
	currentStreak := 1

	for i := 1; i < len(dates); i++ {
		expectedDate := dates[i-1].AddDate(0, 0, 1)
		if dates[i].Equal(expectedDate) {
			currentStreak++
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}
		} else {
			currentStreak = 1
		}
	}

	return maxStreak
}

// GetWeeklyAnalytics returns completion data for the last 7 days
func (r *StatsRepository) GetWeeklyAnalytics(ctx context.Context, userID string) (*query.WeeklyAnalytics, error) {
	analytics := &query.WeeklyAnalytics{
		Days: make([]query.DailyAnalytics, 0, 7),
	}

	// Get total active habits count for calculating percentages
	var activeHabitsCount int
	err := r.db.GetContext(ctx, &activeHabitsCount,
		`SELECT COUNT(*) FROM habits WHERE user_id = $1 AND is_active = true`, userID)
	if err != nil {
		return nil, err
	}

	// If no active habits, return zeros
	if activeHabitsCount == 0 {
		activeHabitsCount = 1 // Avoid division by zero
	}

	// Get logs for each of the last 7 days
	today := time.Now().Truncate(24 * time.Hour)
	totalCompletion := 0

	for i := 6; i >= 0; i-- {
		day := today.AddDate(0, 0, -i)
		dayName := day.Format("Mon")
		dateStr := day.Format("2006-01-02")

		var logsCount int
		err := r.db.GetContext(ctx, &logsCount,
			`SELECT COUNT(DISTINCT habit_id) FROM habit_logs WHERE user_id = $1 AND log_date = $2`,
			userID, day)
		if err != nil {
			logsCount = 0
		}

		// Calculate completion percentage (habits logged / total active habits)
		completionPercentage := int(float64(logsCount) / float64(activeHabitsCount) * 100.0)
		if completionPercentage > 100 {
			completionPercentage = 100
		}
		totalCompletion += completionPercentage

		analytics.Days = append(analytics.Days, query.DailyAnalytics{
			DayName:              dayName,
			Date:                 dateStr,
			LogsCount:            logsCount,
			CompletionPercentage: completionPercentage,
		})
	}

	analytics.AverageCompletion = totalCompletion / 7

	return analytics, nil
}

// GetHabitsDueForReminder returns habits that are active, daily, have no logs for today,
// and either have reminder_time matching the current time in user's timezone, or have NULL reminder_time at 8 PM user's local time.
func (r *StatsRepository) GetHabitsDueForReminder(ctx context.Context) ([]query.ReminderHabit, error) {
	var habits []query.ReminderHabit
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Use PostgreSQL timezone functions to compare reminder_time with current time in user's timezone
	// The key is: TO_CHAR(NOW() AT TIME ZONE u.timezone, 'HH24:MI') gives current time in user's local timezone
	sqlQuery := `
		SELECT h.user_id, h.habit_id, h.name, h.reminder_time
		FROM habits h
		JOIN users u ON h.user_id = u.user_id
		LEFT JOIN habit_logs l ON h.habit_id = l.habit_id AND l.log_date = $1
		WHERE h.is_active = true
		  AND h.frequency = 'daily'
		  AND l.habit_id IS NULL
		  AND (
		      -- Habit has custom reminder_time and it matches current time in user's timezone
		      (h.reminder_time IS NOT NULL AND h.reminder_time = TO_CHAR(NOW() AT TIME ZONE COALESCE(u.timezone, 'UTC'), 'HH24:MI'))
		      OR
		      -- Habit has no custom reminder_time and it's 8 PM in user's timezone (default)
		      (h.reminder_time IS NULL AND TO_CHAR(NOW() AT TIME ZONE COALESCE(u.timezone, 'UTC'), 'HH24:MI') = '20:00')
		  )
	`

	err := r.db.SelectContext(ctx, &habits, sqlQuery, today)
	return habits, err
}

// Time helper functions

func startOfWeek(t time.Time) time.Time {
	// Start of week (Monday)
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday
	}
	return t.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
}

func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
