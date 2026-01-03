package command

import (
	"context"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// RevokeAllOtherSessionsCommand revokes all sessions except the current one
type RevokeAllOtherSessionsCommand struct {
	UserID           string
	CurrentSessionID string
}

// RevokeAllOtherSessionsResult contains the count of revoked sessions
type RevokeAllOtherSessionsResult struct {
	RevokedCount int
}

// RevokeAllOtherSessionsHandler handles session revocation
type RevokeAllOtherSessionsHandler decorator.CommandHandlerWithResult[RevokeAllOtherSessionsCommand, RevokeAllOtherSessionsResult]

type revokeAllOtherSessionsHandler struct {
	sessionRepo session.Repository
}

// NewRevokeAllOtherSessionsHandler creates a new handler
func NewRevokeAllOtherSessionsHandler(
	sessionRepo session.Repository,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) RevokeAllOtherSessionsHandler {
	if sessionRepo == nil {
		panic("nil session repo")
	}

	return decorator.ApplyCommandResultDecorators[RevokeAllOtherSessionsCommand, RevokeAllOtherSessionsResult](
		revokeAllOtherSessionsHandler{sessionRepo: sessionRepo},
		log,
		metricsClient,
	)
}

func (h revokeAllOtherSessionsHandler) Handle(ctx context.Context, cmd RevokeAllOtherSessionsCommand) (RevokeAllOtherSessionsResult, error) {
	userID, err := uuid.Parse(cmd.UserID)
	if err != nil {
		return RevokeAllOtherSessionsResult{}, apperror.ValidationFailed("invalid user ID")
	}

	currentSessionID, err := uuid.Parse(cmd.CurrentSessionID)
	if err != nil {
		return RevokeAllOtherSessionsResult{}, apperror.ValidationFailed("invalid session ID")
	}

	// Get all sessions for the user
	sessions, err := h.sessionRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return RevokeAllOtherSessionsResult{}, apperror.InternalError(err)
	}

	revokedCount := 0
	for _, sess := range sessions {
		// Skip the current session
		if sess.SessionID == currentSessionID {
			continue
		}

		// Block/revoke this session
		sess.IsBlocked = true
		if err := h.sessionRepo.Update(ctx, sess); err != nil {
			// Log but continue with other sessions
			continue
		}
		revokedCount++
	}

	return RevokeAllOtherSessionsResult{RevokedCount: revokedCount}, nil
}
