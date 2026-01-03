package command

import (
	"context"
	"fmt"
	"time"

	"github.com/semmidev/ethos-go/internal/auth/domain/gateway"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/random"
	"github.com/semmidev/ethos-go/internal/common/validator"
)

type ResendVerificationCommand struct {
	Email string `json:"email" validate:"required,email"`
}

type ResendVerificationHandler decorator.CommandHandler[ResendVerificationCommand]

type resendVerificationHandler struct {
	userRepo   user.Repository
	validator  *validator.Validator
	dispatcher gateway.TaskDispatcher
}

func NewResendVerificationHandler(
	userRepo user.Repository,
	validator *validator.Validator,
	dispatcher gateway.TaskDispatcher,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ResendVerificationHandler {
	return decorator.ApplyCommandDecorators[ResendVerificationCommand](
		resendVerificationHandler{
			userRepo:   userRepo,
			validator:  validator,
			dispatcher: dispatcher,
		},
		log,
		metricsClient,
	)
}

func (h resendVerificationHandler) Handle(ctx context.Context, cmd ResendVerificationCommand) error {
	if err := h.validator.Validate(cmd); err != nil {
		return apperror.ValidationFailed(err.Error())
	}

	u, err := h.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return apperror.NotFound("User", cmd.Email)
	}

	if u.IsVerified {
		return apperror.ValidationFailed("user already verified")
	}

	// Generate code
	code, err := random.GenerateNumericOTP(6)
	if err != nil {
		return apperror.InternalError(err)
	}
	expiresAt := time.Now().Add(15 * time.Minute)

	u.VerifyToken = &code
	u.VerifyExpiresAt = &expiresAt

	if err := h.userRepo.Update(ctx, u); err != nil {
		return apperror.InternalError(err)
	}

	// Enqueue task
	payload := &gateway.PayloadSendVerifyEmail{
		UserID:                     u.UserID,
		Name:                       u.Name,
		Email:                      u.Email,
		VerificationCode:           code,
		VerificationCodeExpiration: 15,
	}

	if err := h.dispatcher.DispatchSendVerifyEmail(ctx, payload); err != nil {
		// Just log error, don't fail transaction as DB update succeeded
		// In strictly atomic env, we'd use outbox pattern
		// For now, return error so client might retry
		return fmt.Errorf("failed to enqueue email: %w", err)
	}

	return nil
}
