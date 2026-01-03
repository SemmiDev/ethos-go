package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type MarkAllRead struct {
	UserID string
}

type MarkAllReadHandler decorator.CommandHandler[MarkAllRead]

type markAllReadHandler struct {
	repo domain.NotificationRepository
}

func NewMarkAllReadHandler(
	repo domain.NotificationRepository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) MarkAllReadHandler {
	return decorator.ApplyCommandDecorators[MarkAllRead](
		markAllReadHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h markAllReadHandler) Handle(ctx context.Context, cmd MarkAllRead) error {
	return h.repo.MarkAllAsRead(ctx, cmd.UserID)
}
