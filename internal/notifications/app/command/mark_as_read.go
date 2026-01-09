package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type MarkAsRead struct {
	NotificationID string
	UserID         string
}

type MarkAsReadHandler decorator.CommandHandler[MarkAsRead]

type markAsReadHandler struct {
	repo domain.NotificationRepository
}

func NewMarkAsReadHandler(
	repo domain.NotificationRepository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) MarkAsReadHandler {
	return decorator.ApplyCommandDecorators(
		markAsReadHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h markAsReadHandler) Handle(ctx context.Context, cmd MarkAsRead) error {
	notif, err := h.repo.FindByID(ctx, cmd.NotificationID)
	if err != nil {
		return err
	}

	if notif.UserID != cmd.UserID {
		return apperror.Unauthorized("notification does not belong to user")
	}

	notif.MarkAsRead()

	return h.repo.Update(ctx, notif)
}
