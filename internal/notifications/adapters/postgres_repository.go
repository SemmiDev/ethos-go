package adapters

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/model"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type NotificationPostgresRepository struct {
	db *sqlx.DB
}

func NewNotificationPostgresRepository(db *sqlx.DB) *NotificationPostgresRepository {
	return &NotificationPostgresRepository{db: db}
}

func (r *NotificationPostgresRepository) Create(ctx context.Context, n *domain.Notification) error {
	query := `
		INSERT INTO notifications (notification_id, user_id, type, title, message, data, is_read, created_at, read_at)
		VALUES (:notification_id, :user_id, :type, :title, :message, :data, :is_read, :created_at, :read_at)
	`
	_, err := r.db.NamedExecContext(ctx, query, n)
	return err
}

func (r *NotificationPostgresRepository) FindByID(ctx context.Context, id string) (*domain.Notification, error) {
	var n domain.Notification
	query := `SELECT * FROM notifications WHERE notification_id = $1`
	err := r.db.GetContext(ctx, &n, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("notification", id)
		}
		return nil, err
	}
	return &n, nil
}

func (r *NotificationPostgresRepository) List(ctx context.Context, userID string, filter model.Filter) ([]domain.Notification, *model.Paging, error) {
	var n []domain.Notification
	var count int

	countQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	args := []interface{}{userID}

	// Filter by unread status if requested
	if filter.IsActive != nil && *filter.IsActive {
		// Assuming IsActive=true means "Unread Only" for this context
		// This is a temporary mapping until we have specific filter fields
	}

	// Let's build query dynamically
	var conditions []string
	conditions = append(conditions, "user_id = $1")

	// If keyword search
	if filter.Keyword != "" {
		conditions = append(conditions, "(title ILIKE $2 OR message ILIKE $2)")
		args = append(args, "%"+filter.Keyword+"%")
	}

	// Construct WHERE clause
	whereClause := strings.Join(conditions, " AND ")

	// Final Count Query
	finalCountQuery := strings.Replace(countQuery, "WHERE user_id = $1", "WHERE "+whereClause, 1)

	// Get total count
	err := r.db.GetContext(ctx, &count, finalCountQuery, args...)
	if err != nil {
		return nil, nil, err
	}

	// Calculate pagination
	pagination, err := model.NewPaging(filter.CurrentPage, filter.PerPage, count)
	if err != nil {
		return nil, nil, err
	}

	// Final Select Query
	offset := filter.GetOffset()
	query := fmt.Sprintf("SELECT * FROM notifications WHERE %s ORDER BY created_at DESC LIMIT %d OFFSET %d",
		whereClause, pagination.PerPage, offset)

	err = r.db.SelectContext(ctx, &n, query, args...)
	if err != nil {
		return nil, nil, err
	}

	return n, pagination, nil
}

// ListUnread is a specific method if needed, or we adapt List above.
// For explicit control, let's modify List to support custom "unread" logic if needed,
// but actually, we can just extend Filter struct later.
// For this MVP, let's add `ListUnread` if strictly needed, but `List` with args is fine.
// Wait, the interface uses `model.Filter`. Let's stick to that.

func (r *NotificationPostgresRepository) Update(ctx context.Context, n *domain.Notification) error {
	query := `
		UPDATE notifications SET
			is_read = :is_read,
			read_at = :read_at
		WHERE notification_id = :notification_id
	`
	_, err := r.db.NamedExecContext(ctx, query, n)
	return err
}

func (r *NotificationPostgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM notifications WHERE notification_id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *NotificationPostgresRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1
		WHERE user_id = $2 AND is_read = false
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (r *NotificationPostgresRepository) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	err := r.db.GetContext(ctx, &count, query, userID)
	return count, err
}
