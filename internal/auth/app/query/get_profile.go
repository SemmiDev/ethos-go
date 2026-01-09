package query

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// GetProfileQuery gets user profile
type GetProfileQuery struct {
	UserID string
}

// ProfileResult contains user profile data
type ProfileResult struct {
	UserID    string
	Name      string
	Email     string
	Timezone  string
	CreatedAt time.Time
}

// GetProfileHandler handles profile queries
type GetProfileHandler decorator.QueryHandler[GetProfileQuery, ProfileResult]

type getProfileHandler struct {
	repo user.UserReader
}

// NewGetProfileHandler creates a new handler with decorators
func NewGetProfileHandler(
	repo user.UserReader,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetProfileHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyQueryDecorators[GetProfileQuery, ProfileResult](
		getProfileHandler{repo: repo},
		log,
		metricsClient,
	)
}

func (h getProfileHandler) Handle(ctx context.Context, q GetProfileQuery) (ProfileResult, error) {
	userID, err := uuid.Parse(q.UserID)
	if err != nil {
		return ProfileResult{}, apperror.ValidationFailed("invalid user ID")
	}

	existingUser, err := h.repo.FindByID(ctx, userID)
	if err != nil {
		return ProfileResult{}, apperror.NotFound("user", q.UserID)
	}

	// Use getter methods instead of direct field access
	return ProfileResult{
		UserID:    existingUser.UserID().String(),
		Name:      existingUser.Name(),
		Email:     existingUser.Email(),
		Timezone:  existingUser.Timezone(),
		CreatedAt: existingUser.CreatedAt(),
	}, nil
}
