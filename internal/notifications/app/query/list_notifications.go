package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/model"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type ListNotifications struct {
	UserID string
	Filter model.Filter
}

type ListNotificationsResult struct {
	Notifications []domain.Notification `json:"notifications"`
	Pagination    *model.Paging         `json:"pagination"`
}

type ListNotificationsHandler decorator.QueryHandler[ListNotifications, *ListNotificationsResult]

type listNotificationsHandler struct {
	repo domain.NotificationRepository
}

func NewListNotificationsHandler(
	repo domain.NotificationRepository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ListNotificationsHandler {
	return decorator.ApplyQueryDecorators[ListNotifications, *ListNotificationsResult](
		listNotificationsHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h listNotificationsHandler) Handle(ctx context.Context, q ListNotifications) (*ListNotificationsResult, error) {
	notifs, paging, err := h.repo.List(ctx, q.UserID, q.Filter)
	if err != nil {
		return nil, err
	}

	return &ListNotificationsResult{
		Notifications: notifs,
		Pagination:    paging,
	}, nil
}
