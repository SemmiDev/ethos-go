package session

import (
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/common/random"
)

// Session represents an authenticated user session in our system.
// Each session is tied to a specific device and contains security information
// like IP address and user agent to help detect suspicious activity.
// Fields are private to enforce encapsulation - use getters for read access.
type Session struct {
	sessionID    uuid.UUID
	userID       uuid.UUID
	refreshToken string
	userAgent    string
	clientIP     string
	isBlocked    bool
	expiresAt    time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

// Getters for Session fields

func (s *Session) SessionID() uuid.UUID { return s.sessionID }
func (s *Session) UserID() uuid.UUID    { return s.userID }
func (s *Session) RefreshToken() string { return s.refreshToken }
func (s *Session) UserAgent() string    { return s.userAgent }
func (s *Session) ClientIP() string     { return s.clientIP }
func (s *Session) IsBlocked() bool      { return s.isBlocked }
func (s *Session) ExpiresAt() time.Time { return s.expiresAt }
func (s *Session) CreatedAt() time.Time { return s.createdAt }
func (s *Session) UpdatedAt() time.Time { return s.updatedAt }

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
		sessionID:    sessionID,
		userID:       userID,
		refreshToken: refreshToken,
		userAgent:    userAgent,
		clientIP:     clientIP,
		isBlocked:    false,
		expiresAt:    expiresAt,
		createdAt:    now,
		updatedAt:    now,
	}
}

// UnmarshalSessionFromDatabase reconstructs a Session from database fields.
// This is used by the adapter layer to convert from database model to domain entity.
func UnmarshalSessionFromDatabase(
	sessionID uuid.UUID,
	userID uuid.UUID,
	refreshToken string,
	userAgent string,
	clientIP string,
	isBlocked bool,
	expiresAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *Session {
	return &Session{
		sessionID:    sessionID,
		userID:       userID,
		refreshToken: refreshToken,
		userAgent:    userAgent,
		clientIP:     clientIP,
		isBlocked:    isBlocked,
		expiresAt:    expiresAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// IsExpired checks if the session has passed its expiration time.
// This is important for security - expired sessions should not be usable.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

// IsValid checks if a session can be used for authentication.
// A session is valid only if it's not blocked and hasn't expired.
func (s *Session) IsValid() bool {
	return !s.isBlocked && !s.IsExpired()
}

// Block marks this session as blocked, preventing further use.
// This is useful when we detect suspicious activity or when a user
// explicitly logs out from a specific device.
func (s *Session) Block() {
	s.isBlocked = true
	s.updatedAt = time.Now()
}

// MatchesToken checks if the provided token matches this session's refresh token.
// We use this during token refresh to verify the client has the correct token.
func (s *Session) MatchesToken(token string) bool {
	return s.refreshToken == token
}

// Refresh updates the session with a new refresh token and expiration time.
// This is called when a client uses their refresh token to get a new access token.
func (s *Session) Refresh(newToken string, newExpiry time.Time) {
	s.refreshToken = newToken
	s.expiresAt = newExpiry
	s.updatedAt = time.Now()
}
