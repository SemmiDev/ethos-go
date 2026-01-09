package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type DeleteNotification struct {
	NotificationID string
	UserID         string
}

type DeleteNotificationHandler decorator.CommandHandler[DeleteNotification]

type deleteNotificationHandler struct {
	repo domain.NotificationRepository
}

func NewDeleteNotificationHandler(
	repo domain.NotificationRepository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) DeleteNotificationHandler {
	return decorator.ApplyCommandDecorators(
		deleteNotificationHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h deleteNotificationHandler) Handle(ctx context.Context, cmd DeleteNotification) error {
	notif, err := h.repo.FindByID(ctx, cmd.NotificationID)
	if err != nil {
		return err
	}

	if notif.UserID != cmd.UserID {
		return apperror.Unauthorized("notification does not belong to user")
	}

	return h.repo.Delete(ctx, cmd.NotificationID)
}
