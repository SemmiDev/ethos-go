package query

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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

// ExportUserDataHandler handles data export queries
type ExportUserDataHandler decorator.QueryHandler[ExportUserDataQuery, ExportedData]

type exportUserDataHandler struct {
	userRepo   user.Repository
	exportRepo ExportDataRepository
}

// NewExportUserDataHandler creates a new handler
func NewExportUserDataHandler(
	userRepo user.Repository,
	exportRepo ExportDataRepository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ExportUserDataHandler {
	return decorator.ApplyQueryDecorators[ExportUserDataQuery, ExportedData](
		exportUserDataHandler{
			userRepo:   userRepo,
			exportRepo: exportRepo,
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

	// Use getter methods for domain entity
	exportedUser := ExportedUser{
		ID:           u.UserID().String(),
		Email:        u.Email(),
		Name:         u.Name(),
		Timezone:     u.Timezone(),
		AuthProvider: u.AuthProvider(),
		IsVerified:   u.IsVerified(),
		CreatedAt:    u.CreatedAt(),
	}

	// Fetch habits via repository
	habits, err := h.exportRepo.GetUserHabits(ctx, q.UserID)
	if err != nil {
		habits = []ExportedHabit{} // graceful fallback
	}

	// Fetch habit logs via repository
	logs, err := h.exportRepo.GetUserHabitLogs(ctx, q.UserID)
	if err != nil {
		logs = []ExportedHabitLog{} // graceful fallback
	}

	// Fetch notifications via repository
	notifs, err := h.exportRepo.GetUserNotifications(ctx, q.UserID)
	if err != nil {
		notifs = []ExportedNotif{} // graceful fallback
	}

	// Convert ExportedNotif Data field to json.RawMessage for response
	notifications := make([]ExportedNotif, len(notifs))
	for i, n := range notifs {
		notifications[i] = ExportedNotif{
			ID:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Message:   n.Message,
			Data:      n.Data,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt,
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

// Ensure types are exported for JSON marshaling
var _ json.Marshaler = (*ExportedData)(nil)

func (e ExportedData) MarshalJSON() ([]byte, error) {
	type Alias ExportedData
	return json.Marshal(&struct{ Alias }{Alias: Alias(e)})
}
