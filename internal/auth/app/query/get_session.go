package query

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// SessionDTO represents session data in a format suitable for clients.
type SessionDTO struct {
	SessionID string    `json:"session_id"`
	UserAgent string    `json:"user_agent"`
	ClientIP  string    `json:"client_ip"`
	IsBlocked bool      `json:"is_blocked"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`  // Computed: not blocked and not expired
	IsCurrent bool      `json:"is_current"` // Is this the session making the request?
}

// GetSessionQuery requests information about a specific session.
type GetSessionQuery struct {
	SessionID        string
	UserID           string // For authorization - users can only see their own sessions
	CurrentSessionID string // To mark which session is current
}

// GetSessionHandler retrieves session information.
type GetSessionHandler decorator.QueryHandler[GetSessionQuery, *SessionDTO]

type getSessionHandler struct {
	sessionRepo session.Repository
}

// NewGetSessionHandler creates a handler with its dependencies.
func NewGetSessionHandler(
	sessionRepo session.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) GetSessionHandler {
	return decorator.ApplyQueryDecorators[GetSessionQuery, *SessionDTO](
		getSessionHandler{sessionRepo: sessionRepo},
		log,
		metricsClient,
	)
}

// Handle executes the query to get session details.
// Translates domain errors to AppErrors.
func (h getSessionHandler) Handle(ctx context.Context, query GetSessionQuery) (*SessionDTO, error) {
	// Parse the session ID
	sessionID, err := uuid.Parse(query.SessionID)
	if err != nil {
		return nil, apperror.InvalidInput("session_id", "invalid UUID format").
			WithError(err)
	}

	// Parse the user ID for authorization
	userID, err := uuid.Parse(query.UserID)
	if err != nil {
		return nil, apperror.InvalidInput("user_id", "invalid UUID format").
			WithError(err)
	}

	// Fetch the session from the repository
	sess, err := h.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, h.translateError(err, "find session")
	}

	// Verify the session belongs to the requesting user
	if sess.UserID() != userID {
		return nil, apperror.OperationNotAllowed(
			"view session",
			"session does not belong to this user",
		).WithDetails("session_id", query.SessionID).
			WithDetails("user_id", query.UserID)
	}

	// Convert domain model to DTO
	isCurrent := query.SessionID == query.CurrentSessionID
	return toSessionDTO(sess, isCurrent), nil
}

// translateError converts domain errors to AppErrors
func (h getSessionHandler) translateError(err error, operation string) *apperror.AppError {
	switch {
	case errors.Is(err, session.ErrNotFound):
		return apperror.NotFound("session", "")
	}

	if appErr := apperror.GetAppError(err); appErr != nil {
		return appErr
	}

	return apperror.DatabaseError(operation, err)
}
