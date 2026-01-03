package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

type LogoutCommand struct {
	SessionID string
}

type LogoutHandler decorator.CommandHandler[LogoutCommand]

type logoutHandler struct {
	sessionRepo session.Repository
}

func NewLogoutHandler(
	sessionRepo session.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) LogoutHandler {
	return decorator.ApplyCommandDecorators[LogoutCommand](
		logoutHandler{sessionRepo: sessionRepo},
		log,
		metricsClient,
	)
}

func (h logoutHandler) Handle(ctx context.Context, cmd LogoutCommand) error {
	sessionID, err := uuid.Parse(cmd.SessionID)
	if err != nil {
		return apperror.ValidationFailed("invalid session ID")
	}

	return h.sessionRepo.Delete(ctx, sessionID)
}

type LogoutAllCommand struct {
	UserID string
}

type LogoutAllHandler decorator.CommandHandler[LogoutAllCommand]

type logoutAllHandler struct {
	sessionRepo session.Repository
}

func NewLogoutAllHandler(
	sessionRepo session.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) LogoutAllHandler {
	return decorator.ApplyCommandDecorators[LogoutAllCommand](
		logoutAllHandler{sessionRepo: sessionRepo},
		log,
		metricsClient,
	)
}

func (h logoutAllHandler) Handle(ctx context.Context, cmd LogoutAllCommand) error {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return apperror.ValidationFailed("invalid user ID")
	}

	return h.sessionRepo.DeleteAllByUserID(ctx, userID)
}
