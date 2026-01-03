package query

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// ExportUserDataQuery request to export all user data
type ExportUserDataQuery struct {
	UserID string
}

// ExportedData contains all user data bundled for GDPR export
type ExportedData struct {
	ExportedAt    time.Time          `json:"exported_at"`
	User          ExportedUser       `json:"user"`
	Habits        []ExportedHabit    `json:"habits"`
	HabitLogs     []ExportedHabitLog `json:"habit_logs"`
	Notifications []ExportedNotif    `json:"notifications"`
}

type ExportedUser struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	Timezone     string    `json:"timezone"`
	AuthProvider string    `json:"auth_provider"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
}

type ExportedHabit struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  *string   `json:"description"`
	Frequency    string    `json:"frequency"`
	TargetCount  int       `json:"target_count"`
	IsActive     bool      `json:"is_active"`
	ReminderTime *string   `json:"reminder_time"`
	CreatedAt    time.Time `json:"created_at"`
}

type ExportedHabitLog struct {
	ID        string    `json:"id"`
	HabitID   string    `json:"habit_id"`
	LogDate   string    `json:"log_date"`
	Count     int       `json:"count"`
	Note      *string   `json:"note"`
	CreatedAt time.Time `json:"created_at"`
}

type ExportedNotif struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Title     string          `json:"title"`
	Message   string          `json:"message"`
	Data      json.RawMessage `json:"data"`
	IsRead    bool            `json:"is_read"`
	CreatedAt time.Time       `json:"created_at"`
}

// ExportUserDataHandler handles data export queries
type ExportUserDataHandler decorator.QueryHandler[ExportUserDataQuery, ExportedData]

type exportUserDataHandler struct {
	userRepo user.Repository
	db       *sqlx.DB
}

// NewExportUserDataHandler creates a new handler
func NewExportUserDataHandler(
	userRepo user.Repository,
	db *sqlx.DB,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ExportUserDataHandler {
	return decorator.ApplyQueryDecorators[ExportUserDataQuery, ExportedData](
		exportUserDataHandler{
			userRepo: userRepo,
			db:       db,
		},
		log,
		metricsClient,
	)
}

func (h exportUserDataHandler) Handle(ctx context.Context, q ExportUserDataQuery) (ExportedData, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return ExportedData{}, apperror.ValidationFailed("invalid user ID")
	}

	// Fetch user
	u, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ExportedData{}, apperror.NotFound("user", q.UserID)
	}

	exportedUser := ExportedUser{
		ID:           u.UserID.String(),
		Email:        u.Email,
		Name:         u.Name,
		Timezone:     u.Timezone,
		AuthProvider: u.AuthProvider,
		IsVerified:   u.IsVerified,
		CreatedAt:    u.CreatedAt,
	}

	// Fetch habits
	var habits []ExportedHabit
	habitQuery := `SELECT habit_id, name, description, frequency, target_count, is_active, reminder_time, created_at
	               FROM habits WHERE user_id = $1 ORDER BY created_at`
	rows, err := h.db.QueryxContext(ctx, habitQuery, q.UserID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var hab struct {
				HabitID      string    `db:"habit_id"`
				Name         string    `db:"name"`
				Description  *string   `db:"description"`
				Frequency    string    `db:"frequency"`
				TargetCount  int       `db:"target_count"`
				IsActive     bool      `db:"is_active"`
				ReminderTime *string   `db:"reminder_time"`
				CreatedAt    time.Time `db:"created_at"`
			}
			if rows.StructScan(&hab) == nil {
				habits = append(habits, ExportedHabit{
					ID:           hab.HabitID,
					Name:         hab.Name,
					Description:  hab.Description,
					Frequency:    hab.Frequency,
					TargetCount:  hab.TargetCount,
					IsActive:     hab.IsActive,
					ReminderTime: hab.ReminderTime,
					CreatedAt:    hab.CreatedAt,
				})
			}
		}
	}

	// Fetch habit logs
	var logs []ExportedHabitLog
	logQuery := `SELECT log_id, habit_id, log_date, count, note, created_at
	             FROM habit_logs WHERE user_id = $1 ORDER BY log_date DESC`
	logRows, err := h.db.QueryxContext(ctx, logQuery, q.UserID)
	if err == nil {
		defer logRows.Close()
		for logRows.Next() {
			var log struct {
				LogID     string    `db:"log_id"`
				HabitID   string    `db:"habit_id"`
				LogDate   time.Time `db:"log_date"`
				Count     int       `db:"count"`
				Note      *string   `db:"note"`
				CreatedAt time.Time `db:"created_at"`
			}
			if logRows.StructScan(&log) == nil {
				logs = append(logs, ExportedHabitLog{
					ID:        log.LogID,
					HabitID:   log.HabitID,
					LogDate:   log.LogDate.Format("2006-01-02"),
					Count:     log.Count,
					Note:      log.Note,
					CreatedAt: log.CreatedAt,
				})
			}
		}
	}

	// Fetch notifications
	var notifications []ExportedNotif
	notifQuery := `SELECT notification_id, type, title, message, data, is_read, created_at
	               FROM notifications WHERE user_id = $1 ORDER BY created_at DESC`
	notifRows, err := h.db.QueryxContext(ctx, notifQuery, q.UserID)
	if err == nil {
		defer notifRows.Close()
		for notifRows.Next() {
			var notif struct {
				NotificationID string          `db:"notification_id"`
				Type           string          `db:"type"`
				Title          string          `db:"title"`
				Message        string          `db:"message"`
				Data           json.RawMessage `db:"data"`
				IsRead         bool            `db:"is_read"`
				CreatedAt      time.Time       `db:"created_at"`
			}
			if notifRows.StructScan(&notif) == nil {
				notifications = append(notifications, ExportedNotif{
					ID:        notif.NotificationID,
					Type:      notif.Type,
					Title:     notif.Title,
					Message:   notif.Message,
					Data:      notif.Data,
					IsRead:    notif.IsRead,
					CreatedAt: notif.CreatedAt,
				})
			}
		}
	}

	return ExportedData{
		ExportedAt:    time.Now(),
		User:          exportedUser,
		Habits:        habits,
		HabitLogs:     logs,
		Notifications: notifications,
	}, nil
}
