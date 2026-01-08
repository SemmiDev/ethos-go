package command

import (
	"context"
	"time"

	authevents "github.com/semmidev/ethos-go/internal/auth/domain/events"
	"github.com/semmidev/ethos-go/internal/auth/domain/service"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
)

type ResetPasswordCommand struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordHandler decorator.CommandHandler[ResetPasswordCommand]

type resetPasswordHandler struct {
	userRepo       user.Repository
	passwordHasher service.PasswordHasher
	validator      *validator.Validator
	publisher      events.Publisher
}

func NewResetPasswordHandler(
	userRepo user.Repository,
	passwordHasher service.PasswordHasher,
	validator *validator.Validator,
	publisher events.Publisher, // Injected
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ResetPasswordHandler {
	return decorator.ApplyCommandDecorators[ResetPasswordCommand](
		resetPasswordHandler{
			userRepo:       userRepo,
			passwordHasher: passwordHasher,
			validator:      validator,
			publisher:      publisher,
		},
		log,
		metricsClient,
	)
}

func (h resetPasswordHandler) Handle(ctx context.Context, cmd ResetPasswordCommand) error {
	if err := h.validator.Validate(cmd); err != nil {
		return apperror.ValidationFailed(err.Error())
	}

	u, err := h.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return apperror.NotFound("User", cmd.Email)
	}

	if u.PasswordResetToken == nil || *u.PasswordResetToken != cmd.Code {
		return apperror.ValidationFailed("invalid reset token")
	}

	if u.PasswordResetExpiresAt != nil && u.PasswordResetExpiresAt.Before(time.Now()) {
		return apperror.ValidationFailed("reset token expired")
	}

	// Hash new password
	hashedPassword, err := h.passwordHasher.Hash(ctx, cmd.NewPassword)
	if err != nil {
		return apperror.InternalError(err)
	}

	u.HashedPassword = &hashedPassword
	u.PasswordResetToken = nil
	u.PasswordResetExpiresAt = nil

	if err := h.userRepo.Update(ctx, u); err != nil {
		return apperror.InternalError(err)
	}

	// Publish PasswordChanged event
	event := authevents.NewPasswordChanged(u.UserID.String(), u.Email)
	_ = h.publisher.Publish(ctx, event)

	return nil
}
