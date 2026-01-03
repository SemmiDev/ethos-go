package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain/push"
)

// SubscribePush command to subscribe to push notifications
type SubscribePush struct {
	UserID    string
	Endpoint  string
	P256dh    string
	Auth      string
	UserAgent string
}

// SubscribePushHandler handles push subscription
type SubscribePushHandler decorator.CommandHandler[SubscribePush]

type subscribePushHandler struct {
	repo push.Repository
}

// NewSubscribePushHandler creates a new handler
func NewSubscribePushHandler(
	repo push.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) SubscribePushHandler {
	return decorator.ApplyCommandDecorators[SubscribePush](
		subscribePushHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h subscribePushHandler) Handle(ctx context.Context, cmd SubscribePush) error {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return apperror.ValidationFailed("invalid user ID")
	}

	if cmd.Endpoint == "" {
		return apperror.ValidationFailed("endpoint is required")
	}
	if cmd.P256dh == "" {
		return apperror.ValidationFailed("p256dh key is required")
	}
	if cmd.Auth == "" {
		return apperror.ValidationFailed("auth key is required")
	}

	subscription := push.NewSubscription(userID, cmd.Endpoint, cmd.P256dh, cmd.Auth, cmd.UserAgent)
	if err := h.repo.Save(ctx, subscription); err != nil {
		return apperror.InternalError(err)
	}

	return nil
}

// UnsubscribePush command to unsubscribe from push notifications
type UnsubscribePush struct {
	Endpoint string
}

// UnsubscribePushHandler handles push unsubscription
type UnsubscribePushHandler decorator.CommandHandler[UnsubscribePush]

type unsubscribePushHandler struct {
	repo push.Repository
}

// NewUnsubscribePushHandler creates a new handler
func NewUnsubscribePushHandler(
	repo push.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) UnsubscribePushHandler {
	return decorator.ApplyCommandDecorators[UnsubscribePush](
		unsubscribePushHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h unsubscribePushHandler) Handle(ctx context.Context, cmd UnsubscribePush) error {
	if cmd.Endpoint == "" {
		return apperror.ValidationFailed("endpoint is required")
	}

	if err := h.repo.Delete(ctx, cmd.Endpoint); err != nil {
		return apperror.InternalError(err)
	}

	return nil
}
