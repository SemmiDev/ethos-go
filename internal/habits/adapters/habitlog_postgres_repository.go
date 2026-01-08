package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/model"
	"github.com/semmidev/ethos-go/internal/habits/app/query"
	"github.com/semmidev/ethos-go/internal/habits/domain/habitlog"
)

type habitLogModel struct {
	LogID     string         `db:"log_id"`
	HabitID   string         `db:"habit_id"`
	UserID    string         `db:"user_id"`
	LogDate   time.Time      `db:"log_date"`
	Count     int            `db:"count"`
	Note      sql.NullString `db:"note"` // Nullable field
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type HabitLogPostgresRepository struct {
	db database.DBTX
}

func NewHabitLogPostgresRepository(db database.DBTX) *HabitLogPostgresRepository {
	return &HabitLogPostgresRepository{db: db}
}

// Domain repository implementation

func (r *HabitLogPostgresRepository) AddHabitLog(ctx context.Context, log *habitlog.HabitLog) error {
	q := `
		INSERT INTO habit_logs (log_id, habit_id, user_id, log_date, count, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	// Convert *string to sql.NullString for database insert
	var note sql.NullString
	if log.Note() != nil {
		note = sql.NullString{String: *log.Note(), Valid: true}
	}

	_, err := r.db.ExecContext(ctx, q,
		log.LogID(),
		log.HabitID(),
		log.UserID(),
		log.LogDate(),
		log.Count(),
		note,
		log.CreatedAt(),
		log.UpdatedAt(),
	)
	return err
}

func (r *HabitLogPostgresRepository) GetHabitLog(ctx context.Context, logID, userID string) (*habitlog.HabitLog, error) {
	var model habitLogModel
	q := `SELECT * FROM habit_logs WHERE log_id = $1`
	err := r.db.GetContext(ctx, &model, q, logID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, habitlog.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	log, err := r.unmarshalHabitLog(model)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if err := log.CanBeViewedBy(userID); err != nil {
		return nil, err
	}

	return log, nil
}

func (r *HabitLogPostgresRepository) ListHabitLogs(ctx context.Context, habitID, userID string) ([]*habitlog.HabitLog, error) {
	var models []habitLogModel
	query := `SELECT * FROM habit_logs WHERE habit_id = $1 AND user_id = $2 ORDER BY log_date DESC`
	err := r.db.SelectContext(ctx, &models, query, habitID, userID)
	if err != nil {
		return nil, err
	}

	logs := make([]*habitlog.HabitLog, len(models))
	for i, m := range models {
		log, err := r.unmarshalHabitLog(m)
		if err != nil {
			return nil, err
		}
		logs[i] = log
	}
	return logs, nil
}

func (r *HabitLogPostgresRepository) UpdateHabitLog(
	ctx context.Context,
	logID, userID string,
	updateFn func(ctx context.Context, log *habitlog.HabitLog) (*habitlog.HabitLog, error),
) error {
	return database.RunInTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var model habitLogModel
		q := `SELECT * FROM habit_logs WHERE log_id = $1 FOR UPDATE`
		err := tx.GetContext(ctx, &model, q, logID)
		if errors.Is(err, sql.ErrNoRows) {
			return habitlog.ErrNotFound
		}
		if err != nil {
			return err
		}

		log, err := r.unmarshalHabitLog(model)
		if err != nil {
			return err
		}

		// Authorization check
		if err := log.CanBeModifiedBy(userID); err != nil {
			return err
		}

		// Apply update function
		updatedLog, err := updateFn(ctx, log)
		if err != nil {
			return err
		}

		// Convert *string to sql.NullString for database update
		var note sql.NullString
		if updatedLog.Note() != nil {
			note = sql.NullString{String: *updatedLog.Note(), Valid: true}
		}

		// Persist changes
		updateQuery := `
		UPDATE habit_logs
		SET count = $1, note = $2, log_date = $3, updated_at = $4
		WHERE log_id = $5
	`
		_, err = tx.ExecContext(ctx, updateQuery,
			updatedLog.Count(),
			note,
			updatedLog.LogDate(),
			updatedLog.UpdatedAt(),
			logID,
		)
		return err
	})
}

func (r *HabitLogPostgresRepository) DeleteHabitLog(ctx context.Context, logID, userID string) error {
	// Get log first to check authorization
	log, err := r.GetHabitLog(ctx, logID, userID)
	if err != nil {
		return err
	}

	if err := log.CanBeModifiedBy(userID); err != nil {
		return err
	}

	q := `DELETE FROM habit_logs WHERE log_id = $1`
	_, err = r.db.ExecContext(ctx, q, logID)
	return err
}

func (r *HabitLogPostgresRepository) GetHabitLogByDate(
	ctx context.Context,
	habitID string,
	date time.Time,
	userID string,
) (*habitlog.HabitLog, error) {
	var model habitLogModel
	q := `SELECT * FROM habit_logs WHERE habit_id = $1 AND log_date = $2`
	err := r.db.GetContext(ctx, &model, q, habitID, date)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, habitlog.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	log, err := r.unmarshalHabitLog(model)
	if err != nil {
		return nil, err
	}

	// Authorization check
	if err := log.CanBeViewedBy(userID); err != nil {
		return nil, err
	}

	return log, nil
}

// Query read model implementations

func (r *HabitLogPostgresRepository) GetHabitLogs(
	ctx context.Context,
	habitID, userID string,
	filter model.Filter,
) ([]query.HabitLog, int, error) {
	// Build WHERE conditions
	conditions := []string{"habit_id = $1", "user_id = $2"}
	args := []interface{}{habitID, userID}
	argIndex := 3

	// Date range filters (use filter's StartDate/EndDate)
	if filter.StartDate != nil {
		conditions = append(conditions, fmt.Sprintf("log_date >= $%d", argIndex))
		args = append(args, *filter.StartDate)
		argIndex++
	}
	if filter.EndDate != nil {
		conditions = append(conditions, fmt.Sprintf("log_date <= $%d", argIndex))
		args = append(args, *filter.EndDate)
		argIndex++
	}

	// Keyword search in note
	if filter.HasKeyword() {
		conditions = append(conditions, fmt.Sprintf("note ILIKE $%d", argIndex))
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total for pagination
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM habit_logs WHERE %s", whereClause)
	var totalCount int
	if err := r.db.GetContext(ctx, &totalCount, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Build ORDER BY clause
	orderBy := "log_date"
	if filter.HasSort() {
		allowedColumns := map[string]bool{
			"log_date": true, "created_at": true, "count": true,
		}
		if allowedColumns[filter.SortBy] {
			orderBy = filter.SortBy
		}
	}
	orderDirection := "DESC" // Default to descending for logs (most recent first)
	if !filter.IsDesc() && filter.HasSort() {
		orderDirection = "ASC"
	}

	// Build the main query with pagination
	var q string
	if filter.IsUnlimitedPage() {
		q = fmt.Sprintf(
			"SELECT * FROM habit_logs WHERE %s ORDER BY %s %s",
			whereClause, orderBy, orderDirection,
		)
	} else {
		q = fmt.Sprintf(
			"SELECT * FROM habit_logs WHERE %s ORDER BY %s %s LIMIT $%d OFFSET $%d",
			whereClause, orderBy, orderDirection, argIndex, argIndex+1,
		)
		args = append(args, filter.GetLimit(), filter.GetOffset())
	}

	var models []habitLogModel
	if err := r.db.SelectContext(ctx, &models, q, args...); err != nil {
		return nil, 0, err
	}

	logs := make([]query.HabitLog, len(models))
	for i, m := range models {
		logs[i] = query.HabitLog{
			LogID:     m.LogID,
			HabitID:   m.HabitID,
			UserID:    m.UserID,
			LogDate:   m.LogDate,
			Count:     m.Count,
			Note:      nullStringToPtr(m.Note),
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}

	return logs, totalCount, nil
}

// Helper methods

func (r *HabitLogPostgresRepository) unmarshalHabitLog(model habitLogModel) (*habitlog.HabitLog, error) {
	return habitlog.UnmarshalHabitLogFromDatabase(
		model.LogID,
		model.HabitID,
		model.UserID,
		model.LogDate,
		model.Count,
		nullStringToPtr(model.Note),
		model.CreatedAt,
		model.UpdatedAt,
	)
}
