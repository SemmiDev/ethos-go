package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/common/random"
)

// Session represents an authenticated user session in our system.
// Each session is tied to a specific device and contains security information
// like IP address and user agent to help detect suspicious activity.
type Session struct {
	SessionID    uuid.UUID `db:"session_id"`
	UserID       uuid.UUID `db:"user_id"`
	RefreshToken string    `db:"refresh_token"`
	UserAgent    string    `db:"user_agent"`
	ClientIP     string    `db:"client_ip"`
	IsBlocked    bool      `db:"is_blocked"`
	ExpiresAt    time.Time `db:"expires_at"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// NewSession creates a new session for a user. This is the only way to construct
// a valid session, ensuring all required fields are set properly.
func NewSession(
	sessionID uuid.UUID,
	userID uuid.UUID,
	refreshToken string,
	userAgent string,
	clientIP string,
	expiresAt time.Time,
) *Session {
	if sessionID == uuid.Nil {
		sessionID = random.NewUUID()
	}

	now := time.Now()

	return &Session{
		SessionID:    sessionID,
		UserID:       userID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		IsBlocked:    false,
		ExpiresAt:    expiresAt,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsExpired checks if the session has passed its expiration time.
// This is important for security - expired sessions should not be usable.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if a session can be used for authentication.
// A session is valid only if it's not blocked and hasn't expired.
func (s *Session) IsValid() bool {
	return !s.IsBlocked && !s.IsExpired()
}

// Block marks this session as blocked, preventing further use.
// This is useful when we detect suspicious activity or when a user
// explicitly logs out from a specific device.
func (s *Session) Block() {
	s.IsBlocked = true
	s.UpdatedAt = time.Now()
}

// MatchesToken checks if the provided token matches this session's refresh token.
// We use this during token refresh to verify the client has the correct token.
func (s *Session) MatchesToken(token string) bool {
	return s.RefreshToken == token
}

// Refresh updates the session with a new refresh token and expiration time.
// This is called when a client uses their refresh token to get a new access token.
func (s *Session) Refresh(newToken string, newExpiry time.Time) {
	s.RefreshToken = newToken
	s.ExpiresAt = newExpiry
	s.UpdatedAt = time.Now()
}
