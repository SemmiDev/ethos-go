package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TokenIssuer is an interface for creating authentication tokens.
// The app layer doesn't care whether we use JWT, Paseto, or some
// other token format - it just needs the ability to create tokens with certain claims.
type TokenIssuer interface {
	// IssueAccessToken creates a short-lived token for API access.
	// The token contains the user ID and expires after a configured duration.
	IssueAccessToken(ctx context.Context, userID uuid.UUID, expiresAt time.Time) (string, error)

	// IssueRefreshToken creates a long-lived token for obtaining new access tokens.
	// This token should be stored securely and used only to refresh access tokens.
	IssueRefreshToken(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (string, error)
}

// TokenClaims represents the validated information extracted from a token.
type TokenClaims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	IssuedAt  int64
	ExpiresAt int64
}

// TokenVerifier validates tokens and extracts their claims.
// This is separate from TokenIssuer following the Interface Segregation Principle -
// some components only need to verify tokens, not issue them.
type TokenVerifier interface {
	// VerifyAccessToken validates an access token and returns its claims.
	// Returns an error if the token is invalid, expired, or malformed.
	VerifyAccessToken(ctx context.Context, token string) (*TokenClaims, error)

	// VerifyRefreshToken validates a refresh token and returns its claims.
	VerifyRefreshToken(ctx context.Context, token string) (*TokenClaims, error)
}
