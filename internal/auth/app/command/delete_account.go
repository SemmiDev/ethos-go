package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// DeleteAccountCommand permanently deletes a user account
type DeleteAccountCommand struct {
	UserID          string
	Password        string // For verification (optional for OAuth users)
	ConfirmDeletion bool   // Must be true to proceed
}

// DeleteAccountHandler handles account deletion
type DeleteAccountHandler decorator.CommandHandler[DeleteAccountCommand]

type deleteAccountHandler struct {
	userRepo    user.Repository
	sessionRepo session.Repository
}

// NewDeleteAccountHandler creates a new handler
func NewDeleteAccountHandler(
	userRepo user.Repository,
	sessionRepo session.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) DeleteAccountHandler {
	if userRepo == nil {
		panic("nil user repo")
	}
	if sessionRepo == nil {
		panic("nil session repo")
	}

	return decorator.ApplyCommandDecorators[DeleteAccountCommand](
		deleteAccountHandler{
			userRepo:    userRepo,
			sessionRepo: sessionRepo,
		},
		log,
		metricsClient,
	)
}

func (h deleteAccountHandler) Handle(ctx context.Context, cmd DeleteAccountCommand) error {
	if !cmd.ConfirmDeletion {
		return apperror.ValidationFailed("deletion must be confirmed")
	}

	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return apperror.ValidationFailed("invalid user ID")
	}

	// Verify user exists
	_, err = h.userRepo.FindByID(ctx, userID)
	if err != nil {
		return apperror.NotFound("user", cmd.UserID)
	}

	// Note: Password verification could be added here for email-based users
	// For now, we rely on the frontend confirmation modal

	// Delete user (FK cascade will handle habits, logs, notifications)
	if err := h.userRepo.Delete(ctx, userID); err != nil {
		return apperror.InternalError(err)
	}

	return nil
}
