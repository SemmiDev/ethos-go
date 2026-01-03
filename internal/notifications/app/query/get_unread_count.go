package query

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type GetUnreadCount struct {
	UserID string
}

type GetUnreadCountHandler decorator.QueryHandler[GetUnreadCount, int]

type getUnreadCountHandler struct {
	repo domain.NotificationRepository
}

func NewGetUnreadCountHandler(
	repo domain.NotificationRepository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetUnreadCountHandler {
	return decorator.ApplyQueryDecorators[GetUnreadCount, int](
		getUnreadCountHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h getUnreadCountHandler) Handle(ctx context.Context, q GetUnreadCount) (int, error) {
	return h.repo.GetUnreadCount(ctx, q.UserID)
}
