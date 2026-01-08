package command

import (
	"context"
	"time"

	"github.com/google/uuid"
	authevents "github.com/semmidev/ethos-go/internal/auth/domain/events"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"golang.org/x/crypto/bcrypt"
)

// UpdateProfileCommand for updating user profile
type UpdateProfileCommand struct {
	UserID   string
	Name     *string
	Email    *string
	Timezone *string
}

// UpdateProfileResult contains the updated profile data
type UpdateProfileResult struct {
	UserID    string
	Name      string
	Email     string
	Timezone  string
	CreatedAt time.Time
}

// UpdateProfileHandler handles profile updates
type UpdateProfileHandler decorator.CommandHandlerWithResult[UpdateProfileCommand, UpdateProfileResult]

type updateProfileHandler struct {
	repo user.Repository
}

// NewUpdateProfileHandler creates a new handler with decorators
func NewUpdateProfileHandler(
	repo user.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) UpdateProfileHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandResultDecorators[UpdateProfileCommand, UpdateProfileResult](
		updateProfileHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h updateProfileHandler) Handle(ctx context.Context, cmd UpdateProfileCommand) (UpdateProfileResult, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return UpdateProfileResult{}, apperror.ValidationFailed("invalid user ID")
	}

	existingUser, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return UpdateProfileResult{}, apperror.NotFound("user", cmd.UserID)
	}

	// Update fields if provided
	if cmd.Name != nil && *cmd.Name != "" {
		existingUser.Name = *cmd.Name
	}
	if cmd.Email != nil && *cmd.Email != "" {
		existingUser.Email = *cmd.Email
	}
	if cmd.Timezone != nil && *cmd.Timezone != "" {
		existingUser.Timezone = *cmd.Timezone
	}
	existingUser.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, existingUser); err != nil {
		return UpdateProfileResult{}, apperror.InternalError(err)
	}

	return UpdateProfileResult{
		UserID:    existingUser.UserID.String(),
		Name:      existingUser.Name,
		Email:     existingUser.Email,
		Timezone:  existingUser.Timezone,
		CreatedAt: existingUser.CreatedAt,
	}, nil
}

// ChangePasswordCommand for changing user password
type ChangePasswordCommand struct {
	UserID          string
	CurrentPassword string
	NewPassword     string
}

// ChangePasswordHandler handles password changes
type ChangePasswordHandler decorator.CommandHandler[ChangePasswordCommand]

type changePasswordHandler struct {
	repo      user.Repository
	publisher events.Publisher
}

// NewChangePasswordHandler creates a new handler with decorators
func NewChangePasswordHandler(
	repo user.Repository,
	publisher events.Publisher, // Injected
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) ChangePasswordHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators[ChangePasswordCommand](
		changePasswordHandler{
			repo:      repo,
			publisher: publisher,
		},
		log,
		metricsClient,
	)
}

func (h changePasswordHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return apperror.ValidationFailed("invalid user ID")
	}

	existingUser, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return apperror.NotFound("user", cmd.UserID)
	}

	// Verify current password
	if existingUser.HashedPassword == nil {
		return apperror.Unauthorized("user has no password set")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*existingUser.HashedPassword), []byte(cmd.CurrentPassword)); err != nil {
		return apperror.Unauthorized("current password is incorrect")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperror.InternalError(err)
	}

	hashedStr := string(hashedPassword)
	existingUser.HashedPassword = &hashedStr
	existingUser.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, existingUser); err != nil {
		return apperror.InternalError(err)
	}

	// Publish PasswordChanged event
	event := authevents.NewPasswordChanged(existingUser.UserID.String(), existingUser.Email)
	_ = h.publisher.Publish(ctx, event)

	return nil
}
