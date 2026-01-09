package command

import (
	"context"
	"time"

	authevents "github.com/semmidev/ethos-go/internal/auth/domain/events"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
)

type VerifyEmailCommand struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}

type VerifyEmailHandler decorator.CommandHandler[VerifyEmailCommand]

type verifyEmailHandler struct {
	userRepo  user.Repository
	validator *validator.Validator
	publisher events.Publisher
}

func NewVerifyEmailHandler(
	userRepo user.Repository,
	validator *validator.Validator,
	publisher events.Publisher, // Injected publisher
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) VerifyEmailHandler {
	return decorator.ApplyCommandDecorators(
		verifyEmailHandler{
			userRepo:  userRepo,
			validator: validator,
			publisher: publisher,
		},
		log,
		metricsClient,
	)
}

func (h verifyEmailHandler) Handle(ctx context.Context, cmd VerifyEmailCommand) error {
	if err := h.validator.Validate(cmd); err != nil {
		return apperror.ValidationFailed(err.Error())
	}

	u, err := h.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return apperror.NotFound("User", cmd.Email)
	}

	if u.IsVerified() {
		return nil
	}

	if u.VerifyToken() == nil || *u.VerifyToken() != cmd.Code {
		return apperror.ValidationFailed("invalid verification code")
	}

	if u.VerifyExpiresAt() != nil && u.VerifyExpiresAt().Before(time.Now()) {
		return apperror.ValidationFailed("verification code expired")
	}

	// Use domain method to mark as verified
	u.MarkVerified()

	if err := h.userRepo.Update(ctx, u); err != nil {
		return apperror.InternalError(err)
	}

	// Publish UserVerified event
	event := authevents.NewUserVerified(u.UserID().String(), u.Email())
	_ = h.publisher.Publish(ctx, event)

	return nil
}
