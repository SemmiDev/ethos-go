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

type ForgotPasswordCommand struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordHandler decorator.CommandHandler[ForgotPasswordCommand]

type forgotPasswordHandler struct {
	userRepo   user.Repository
	validator  *validator.Validator
	dispatcher gateway.TaskDispatcher
}

func NewForgotPasswordHandler(
	userRepo user.Repository,
	validator *validator.Validator,
	dispatcher gateway.TaskDispatcher,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ForgotPasswordHandler {
	return decorator.ApplyCommandDecorators(
		forgotPasswordHandler{
			userRepo:   userRepo,
			validator:  validator,
			dispatcher: dispatcher,
		},
		log,
		metricsClient,
	)
}

func (h forgotPasswordHandler) Handle(ctx context.Context, cmd ForgotPasswordCommand) error {
	if err := h.validator.Validate(cmd); err != nil {
		return apperror.ValidationFailed(err.Error())
	}

	u, err := h.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		// Return success even if not found to prevent user enumeration
		return nil
	}

	// Generate code
	code, err := random.GenerateNumericOTP(6)
	if err != nil {
		return apperror.InternalError(err)
	}
	expiresAt := time.Now().Add(15 * time.Minute)

	// Use domain setter
	u.SetPasswordResetToken(&code, &expiresAt)

	if err := h.userRepo.Update(ctx, u); err != nil {
		return apperror.InternalError(err)
	}

	// Enqueue task
	payload := &gateway.PayloadSendForgotPasswordEmail{
		UserID:                     u.UserID(),
		Name:                       u.Name(),
		Email:                      u.Email(),
		VerificationCode:           code,
		VerificationCodeExpiration: 15,
	}

	if err := h.dispatcher.DispatchSendForgotPasswordEmail(ctx, payload); err != nil {
		return fmt.Errorf("failed to enqueue email: %w", err)
	}

	return nil
}
