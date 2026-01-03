package domain

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/model"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	FindByID(ctx context.Context, id string) (*Notification, error)
	List(ctx context.Context, userID string, filter model.Filter) ([]Notification, *model.Paging, error)
	Update(ctx context.Context, notification *Notification) error
	Delete(ctx context.Context, id string) error
	MarkAllAsRead(ctx context.Context, userID string) error
	GetUnreadCount(ctx context.Context, userID string) (int, error)
}
