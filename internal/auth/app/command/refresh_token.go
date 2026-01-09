package command

import (
	"context"
	"time"

	"github.com/semmidev/ethos-go/internal/auth/domain/service"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

type RefreshTokenCommand struct {
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenHandler decorator.CommandHandlerWithResult[RefreshTokenCommand, *RefreshTokenResult]

type refreshTokenHandler struct {
	sessionRepo session.Repository
	tokenIssuer service.TokenIssuer
	authService *session.AuthenticationService
}

func NewRefreshTokenHandler(
	sessionRepo session.Repository,
	tokenIssuer service.TokenIssuer,
	authService *session.AuthenticationService,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) RefreshTokenHandler {
	return decorator.ApplyCommandResultDecorators(
		refreshTokenHandler{
			sessionRepo: sessionRepo,
			tokenIssuer: tokenIssuer,
			authService: authService,
		},
		log,
		metricsClient,
	)
}

func (h refreshTokenHandler) Handle(ctx context.Context, cmd RefreshTokenCommand) (*RefreshTokenResult, error) {
	// Find session by refresh token
	sess, err := h.sessionRepo.FindByRefreshToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, apperror.NotFound("session", "")
	}

	// Validate session
	if !sess.IsValid() {
		return nil, apperror.SessionExpired(nil)
	}

	// Calculate new expiration times
	now := time.Now()
	accessTokenExpiry := now.Add(h.authService.AccessTokenTTL())
	// We might or might not want to extend the refresh token.
	// Common practice: Rotation. Issue a new one, invalidate the old one.
	// But let's check if we want to extend the session life or keep the original session end.
	// If the session model has a fixed absolute expiry, we should respect it.
	// If the session has a sliding window, we extend it.
	// Looking at session.go (not read but inferred), it likely has ExpiresAt.
	// Let's assume sliding window for now (refresh token extends session).
	refreshTokenExpiry := now.Add(h.authService.RefreshTokenTTL())

	// Issue new access token
	accessToken, err := h.tokenIssuer.IssueAccessToken(ctx, sess.UserID(), accessTokenExpiry)
	if err != nil {
		return nil, apperror.InternalError(err)
	}

	// Issue new refresh token
	newRefreshToken, err := h.tokenIssuer.IssueRefreshToken(ctx, sess.SessionID(), refreshTokenExpiry)
	if err != nil {
		return nil, apperror.InternalError(err)
	}

	// Update session with new refresh token and expiry
	sess.Refresh(newRefreshToken, refreshTokenExpiry)

	if err := h.sessionRepo.Update(ctx, sess); err != nil {
		return nil, apperror.DatabaseError("update session", err)
	}

	return &RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
