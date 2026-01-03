package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
	"github.com/semmidev/ethos-go/internal/notifications/domain/push"
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
	repo        domain.NotificationRepository
	pushService push.Service
	log         logger.Logger
}

func NewCreateNotificationHandler(
	repo domain.NotificationRepository,
	pushService push.Service,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) CreateNotificationHandler {
	return decorator.ApplyCommandDecorators[CreateNotification](
		createNotificationHandler{
			repo:        repo,
			pushService: pushService,
			log:         log,
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

	// Send Push Notification (best effort)
	// We don't want to fail the request if push fails
	go func() {
		userID, err := uuid.Parse(cmd.UserID)
		if err != nil {
			return
		}

		payload := push.Payload{
			Title: cmd.Title,
			Body:  cmd.Message,
			Data:  make(map[string]string),
			Icon:  "/icon-192x192.png", // Default icon
			Badge: "/badge-72x72.png",  // Default badge
		}

		// Convert arbitrary data to string map for push payload
		if cmd.Data != nil {
			for k, v := range cmd.Data {
				payload.Data[k] = fmt.Sprint(v)
			}
		}

		// Use background context as the request context might be cancelled
		if err := h.pushService.SendNotification(context.Background(), userID, payload); err != nil {
			h.log.Error(context.Background(), err, "failed to send push notification",
				logger.Field{Key: "user_id", Value: cmd.UserID},
			)
		}
	}()

	return nil
}
