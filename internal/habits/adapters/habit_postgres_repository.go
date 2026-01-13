package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/model"
	"github.com/semmidev/ethos-go/internal/habits/app/query"
	"github.com/semmidev/ethos-go/internal/habits/domain/habit"
)

type habitModel struct {
	HabitID            string         `db:"habit_id"`
	UserID             string         `db:"user_id"`
	Name               string         `db:"name"`
	Description        sql.NullString `db:"description"`
	Frequency          string         `db:"frequency"`
	RecurrenceDays     int16          `db:"recurrence_days"`
	RecurrenceInterval int            `db:"recurrence_interval"`
	TargetCount        int            `db:"target_count"`
	ReminderTime       sql.NullString `db:"reminder_time"`
	IsActive           bool           `db:"is_active"`
	CreatedAt          time.Time      `db:"created_at"`
	UpdatedAt          time.Time      `db:"updated_at"`
}

type statsModel struct {
	HabitID          string     `db:"habit_id"`
	CurrentStreak    int        `db:"current_streak"`
	LongestStreak    int        `db:"longest_streak"`
	TotalCompletions int        `db:"total_completions"`
	LastCompletedAt  *time.Time `db:"last_completed_at"`
	ConsistencyScore float64    `db:"consistency_score"`
	UpdatedAt        time.Time  `db:"updated_at"`
}

type vacationModel struct {
	ID        string     `db:"id"`
	HabitID   string     `db:"habit_id"`
	StartDate time.Time  `db:"start_date"`
	EndDate   *time.Time `db:"end_date"`
	Reason    *string    `db:"reason"`
	CreatedAt time.Time  `db:"created_at"`
}

type HabitPostgresRepository struct {
	db database.DBTX
}

func NewHabitPostgresRepository(db database.DBTX) *HabitPostgresRepository {
	return &HabitPostgresRepository{db: db}
}

func (r *HabitPostgresRepository) AddHabit(ctx context.Context, h *habit.Habit) error {
	query := `
        INSERT INTO habits (habit_id, user_id, name, description, frequency, target_count, reminder_time, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `
	// Convert *string to sql.NullString for database insert
	var description sql.NullString
	if h.Description() != nil {
		description = sql.NullString{String: *h.Description(), Valid: true}
	}

	var reminderTime sql.NullString
	if h.ReminderTime() != nil {
		reminderTime = sql.NullString{String: *h.ReminderTime(), Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		h.HabitID(),
		h.UserID(),
		h.Name(),
		description,
		h.Frequency().String(),
		h.TargetCount(),
		reminderTime,
		h.IsActive(),
		h.CreatedAt(),
		h.UpdatedAt(),
	)
	return err
}

func (r *HabitPostgresRepository) GetHabit(ctx context.Context, habitID, userID string) (*habit.Habit, error) {
	var model habitModel
	query := `SELECT * FROM habits WHERE habit_id = $1`
	err := r.db.GetContext(ctx, &model, query, habitID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, habit.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	h, err := r.unmarshalHabit(model)
	if err != nil {
		return nil, err
	}

	if err := h.CanBeViewedBy(userID); err != nil {
		return nil, err
	}

	return h, nil
}

func (r *HabitPostgresRepository) UpdateHabit(
	ctx context.Context,
	habitID, userID string,
	updateFn func(ctx context.Context, h *habit.Habit) (*habit.Habit, error),
) error {
	return database.RunInTx(ctx, r.db, func(tx database.DBTX) error {
		var model habitModel
		query := `SELECT * FROM habits WHERE habit_id = $1 FOR UPDATE`
		err := tx.GetContext(ctx, &model, query, habitID)
		if errors.Is(err, sql.ErrNoRows) {
			return habit.ErrNotFound
		}
		if err != nil {
			return err
		}

		h, err := r.unmarshalHabit(model)
		if err != nil {
			return err
		}

		if err := h.CanBeViewedBy(userID); err != nil {
			return err
		}

		updatedHabit, err := updateFn(ctx, h)
		if err != nil {
			return err
		}

		// Convert *string to sql.NullString for database update
		var description sql.NullString
		if updatedHabit.Description() != nil {
			description = sql.NullString{String: *updatedHabit.Description(), Valid: true}
		}

		var reminderTime sql.NullString
		if updatedHabit.ReminderTime() != nil {
			reminderTime = sql.NullString{String: *updatedHabit.ReminderTime(), Valid: true}
		}

		updateQuery := `
        UPDATE habits
        SET name = $1, description = $2, frequency = $3, target_count = $4, reminder_time = $5, is_active = $6, updated_at = $7
        WHERE habit_id = $8
    `
		_, err = tx.ExecContext(ctx, updateQuery,
			updatedHabit.Name(),
			description,
			updatedHabit.Frequency().String(),
			updatedHabit.TargetCount(),
			reminderTime,
			updatedHabit.IsActive(),
			updatedHabit.UpdatedAt(),
			habitID,
		)
		return err
	})
}

func (r *HabitPostgresRepository) DeleteHabit(ctx context.Context, habitID, userID string) error {
	h, err := r.GetHabit(ctx, habitID, userID)
	if err != nil {
		return err
	}

	if err := h.CanBeViewedBy(userID); err != nil {
		return err
	}

	query := `DELETE FROM habits WHERE habit_id = $1`
	_, err = r.db.ExecContext(ctx, query, habitID)
	return err
}

func (r *HabitPostgresRepository) ListHabitsByUser(ctx context.Context, userID string) ([]*habit.Habit, error) {
	var models []habitModel
	query := `SELECT * FROM habits WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	habits := make([]*habit.Habit, len(models))
	for i, m := range models {
		h, err := r.unmarshalHabit(m)
		if err != nil {
			return nil, err
		}
		habits[i] = h
	}

	return habits, nil
}

// Habit Stats

func (r *HabitPostgresRepository) GetStats(ctx context.Context, habitID string) (*habit.HabitStats, error) {
	var model statsModel
	query := `SELECT * FROM habit_stats WHERE habit_id = $1`
	err := r.db.GetContext(ctx, &model, query, habitID)
	if errors.Is(err, sql.ErrNoRows) {
		return habit.NewHabitStats(habitID), nil
	}
	if err != nil {
		return nil, err
	}

	return habit.UnmarshalStatsFromDatabase(
		model.HabitID,
		model.CurrentStreak,
		model.LongestStreak,
		model.TotalCompletions,
		model.LastCompletedAt,
		model.ConsistencyScore,
		model.UpdatedAt,
	), nil
}

func (r *HabitPostgresRepository) UpsertStats(ctx context.Context, stats *habit.HabitStats) error {
	query := `
		INSERT INTO habit_stats (habit_id, current_streak, longest_streak, total_completions, last_completed_at, consistency_score, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (habit_id) DO UPDATE SET
			current_streak = EXCLUDED.current_streak,
			longest_streak = EXCLUDED.longest_streak,
			total_completions = EXCLUDED.total_completions,
			last_completed_at = EXCLUDED.last_completed_at,
			consistency_score = EXCLUDED.consistency_score,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.ExecContext(ctx, query,
		stats.HabitID(),
		stats.CurrentStreak(),
		stats.LongestStreak(),
		stats.TotalCompletions(),
		stats.LastCompletedAt(),
		stats.ConsistencyScore(),
		stats.UpdatedAt(),
	)
	return err
}

// Habit Vacations

func (r *HabitPostgresRepository) AddVacation(ctx context.Context, vacation *habit.HabitVacation) error {
	query := `
		INSERT INTO habit_vacations (id, habit_id, start_date, end_date, reason, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		vacation.ID(),
		vacation.HabitID(),
		vacation.StartDate(),
		vacation.EndDate(),
		vacation.Reason(),
		vacation.CreatedAt(),
	)
	return err
}

func (r *HabitPostgresRepository) GetActiveVacation(ctx context.Context, habitID string) (*habit.HabitVacation, error) {
	var model vacationModel
	query := `
		SELECT * FROM habit_vacations
		WHERE habit_id = $1 AND (end_date IS NULL OR end_date >= CURRENT_DATE)
		ORDER BY start_date DESC LIMIT 1
	`
	err := r.db.GetContext(ctx, &model, query, habitID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil // No active vacation
	}
	if err != nil {
		return nil, err
	}

	return habit.UnmarshalVacationFromDatabase(
		model.ID,
		model.HabitID,
		model.StartDate,
		model.EndDate,
		model.Reason,
		model.CreatedAt,
	), nil
}

func (r *HabitPostgresRepository) EndVacation(ctx context.Context, vacationID string, endDate time.Time) error {
	query := `UPDATE habit_vacations SET end_date = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, endDate, vacationID)
	return err
}

func (r *HabitPostgresRepository) ListVacations(ctx context.Context, habitID string) ([]*habit.HabitVacation, error) {
	var models []vacationModel
	query := `SELECT * FROM habit_vacations WHERE habit_id = $1 ORDER BY start_date DESC`
	err := r.db.SelectContext(ctx, &models, query, habitID)
	if err != nil {
		return nil, err
	}

	vacations := make([]*habit.HabitVacation, len(models))
	for i, m := range models {
		vacations[i] = habit.UnmarshalVacationFromDatabase(
			m.ID,
			m.HabitID,
			m.StartDate,
			m.EndDate,
			m.Reason,
			m.CreatedAt,
		)
	}
	return vacations, nil
}

// Query read model implementations

func (r *HabitPostgresRepository) GetHabitQuery(ctx context.Context, habitID, userID string) (*query.Habit, error) {
	var model habitModel
	q := `SELECT * FROM habits WHERE habit_id = $1`
	err := r.db.GetContext(ctx, &model, q, habitID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, habit.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	// Authorization check
	if model.UserID != userID {
		return nil, habit.ErrUnauthorized
	}

	return &query.Habit{
		HabitID:      model.HabitID,
		UserID:       model.UserID,
		Name:         model.Name,
		Description:  nullStringToPtr(model.Description),
		Frequency:    model.Frequency,
		TargetCount:  model.TargetCount,
		ReminderTime: nullStringToPtr(model.ReminderTime),
		IsActive:     model.IsActive,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}, nil
}

func (r *HabitPostgresRepository) ListHabits(ctx context.Context, userID string, filter model.Filter) ([]query.Habit, int, error) {
	// Build WHERE conditions
	conditions := []string{"user_id = $1"}
	args := []interface{}{userID}
	argIndex := 2

	// Status filter
	if filter.ActiveOnly() {
		conditions = append(conditions, "is_active = true")
	} else if filter.InactiveOnly() {
		conditions = append(conditions, "is_active = false")
	}

	// Keyword search (search in name and description)
	if filter.HasKeyword() {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	// Date range filters
	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.StartDate)
		argIndex++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.EndDate)
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM habits WHERE %s", whereClause)
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderBy := "created_at"
	if filter.HasSort() {
		// Validate sort column to prevent SQL injection
		allowedColumns := map[string]bool{
			"name": true, "created_at": true, "updated_at": true, "is_active": true,
		}
		if allowedColumns[filter.SortBy] {
			orderBy = filter.SortBy
		}
	}
	orderDirection := "ASC"
	if filter.IsDesc() {
		orderDirection = "DESC"
	}

	// Build the main query with pagination
	var q string
	if filter.IsUnlimitedPage() {
		q = fmt.Sprintf(
			"SELECT * FROM habits WHERE %s ORDER BY %s %s",
			whereClause, orderBy, orderDirection,
		)
	} else {
		q = fmt.Sprintf(
			"SELECT * FROM habits WHERE %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
			whereClause, orderBy, orderDirection, argIndex, argIndex+1,
		)
		args = append(args, filter.GetLimit(), filter.GetOffset())
	}

	var models []habitModel
	if err := r.db.SelectContext(ctx, &models, q, args...); err != nil {
		return nil, 0, err
	}

	habits := make([]query.Habit, len(models))
	for i, m := range models {
		habits[i] = query.Habit{
			HabitID:      m.HabitID,
			UserID:       m.UserID,
			Name:         m.Name,
			Description:  nullStringToPtr(m.Description),
			Frequency:    m.Frequency,
			TargetCount:  m.TargetCount,
			ReminderTime: nullStringToPtr(m.ReminderTime),
			IsActive:     m.IsActive,
			CreatedAt:    m.CreatedAt,
			UpdatedAt:    m.UpdatedAt,
		}
	}
	return habits, totalCount, nil
}

func (r *HabitPostgresRepository) unmarshalHabit(model habitModel) (*habit.Habit, error) {
	return habit.UnmarshalHabitFromDatabase(
		model.HabitID,
		model.UserID,
		model.Name,
		nullStringToPtr(model.Description),
		model.Frequency,
		model.RecurrenceDays,
		model.RecurrenceInterval,
		model.TargetCount,
		nullStringToPtr(model.ReminderTime),
		model.IsActive,
		model.CreatedAt,
		model.UpdatedAt,
	)
}

// nullStringToPtr converts sql.NullString to *string
// Returns nil if NullString is not valid, otherwise returns pointer to the string value
func nullStringToPtr(ns sql.NullString) *string {
	if !ns.Valid {
		return nil
	}
	return &ns.String
}
