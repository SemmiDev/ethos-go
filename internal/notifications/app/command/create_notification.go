package command

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type CreateNotification struct {
	UserID  string
	Type    domain.NotificationType
	Title   string
	Message string
	Data    map[string]interface{}
}

type CreateNotificationHandler decorator.CommandHandler[CreateNotification]

type createNotificationHandler struct {
	repo domain.NotificationRepository
	log  logger.Logger
}

func NewCreateNotificationHandler(
	repo domain.NotificationRepository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) CreateNotificationHandler {
	return decorator.ApplyCommandDecorators(
		createNotificationHandler{
			repo: repo,
			log:  log,
		},
		log,
		metricsClient,
	)
}

func (h createNotificationHandler) Handle(ctx context.Context, cmd CreateNotification) error {
	notif, err := domain.NewNotification(cmd.UserID, cmd.Type, cmd.Title, cmd.Message, cmd.Data)
	if err != nil {
		return err
	}
	if err := h.repo.Create(ctx, notif); err != nil {
		return err
	}

	return nil
}
