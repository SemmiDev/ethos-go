package adapters

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/internal/auth/app/query"
)

// ExportDataPostgresRepository implements query.ExportDataRepository
type ExportDataPostgresRepository struct {
	db *sqlx.DB
}

// NewExportDataPostgresRepository creates a new export data repository
func NewExportDataPostgresRepository(db *sqlx.DB) *ExportDataPostgresRepository {
	return &ExportDataPostgresRepository{db: db}
}

// GetUserHabits fetches all habits for a user
func (r *ExportDataPostgresRepository) GetUserHabits(ctx context.Context, userID string) ([]query.ExportedHabit, error) {
	q := `SELECT habit_id, name, description, frequency, target_count, is_active, reminder_time, created_at
	      FROM habits WHERE user_id = $1 ORDER BY created_at`

	rows, err := r.db.QueryxContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var habits []query.ExportedHabit
	for rows.Next() {
		var h struct {
			HabitID      string    `db:"habit_id"`
			Name         string    `db:"name"`
			Description  *string   `db:"description"`
			Frequency    string    `db:"frequency"`
			TargetCount  int       `db:"target_count"`
			IsActive     bool      `db:"is_active"`
			ReminderTime *string   `db:"reminder_time"`
			CreatedAt    time.Time `db:"created_at"`
		}
		if err := rows.StructScan(&h); err != nil {
			continue
		}
		habits = append(habits, query.ExportedHabit{
			ID:           h.HabitID,
			Name:         h.Name,
			Description:  h.Description,
			Frequency:    h.Frequency,
			TargetCount:  h.TargetCount,
			IsActive:     h.IsActive,
			ReminderTime: h.ReminderTime,
			CreatedAt:    h.CreatedAt,
		})
	}
	return habits, nil
}

// GetUserHabitLogs fetches all habit logs for a user
func (r *ExportDataPostgresRepository) GetUserHabitLogs(ctx context.Context, userID string) ([]query.ExportedHabitLog, error) {
	q := `SELECT log_id, habit_id, log_date, count, note, created_at
	      FROM habit_logs WHERE user_id = $1 ORDER BY log_date DESC`

	rows, err := r.db.QueryxContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []query.ExportedHabitLog
	for rows.Next() {
		var l struct {
			LogID     string    `db:"log_id"`
			HabitID   string    `db:"habit_id"`
			LogDate   time.Time `db:"log_date"`
			Count     int       `db:"count"`
			Note      *string   `db:"note"`
			CreatedAt time.Time `db:"created_at"`
		}
		if err := rows.StructScan(&l); err != nil {
			continue
		}
		logs = append(logs, query.ExportedHabitLog{
			ID:        l.LogID,
			HabitID:   l.HabitID,
			LogDate:   l.LogDate.Format("2006-01-02"),
			Count:     l.Count,
			Note:      l.Note,
			CreatedAt: l.CreatedAt,
		})
	}
	return logs, nil
}

// GetUserNotifications fetches all notifications for a user
func (r *ExportDataPostgresRepository) GetUserNotifications(ctx context.Context, userID string) ([]query.ExportedNotif, error) {
	q := `SELECT notification_id, type, title, message, data, is_read, created_at
	      FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryxContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []query.ExportedNotif
	for rows.Next() {
		var n struct {
			NotificationID string          `db:"notification_id"`
			Type           string          `db:"type"`
			Title          string          `db:"title"`
			Message        string          `db:"message"`
			Data           json.RawMessage `db:"data"`
			IsRead         bool            `db:"is_read"`
			CreatedAt      time.Time       `db:"created_at"`
		}
		if err := rows.StructScan(&n); err != nil {
			continue
		}
		notifications = append(notifications, query.ExportedNotif{
			ID:        n.NotificationID,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			Data:      n.Data,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
		})
	}
	return notifications, nil
}
