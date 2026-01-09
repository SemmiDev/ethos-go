package adapters

import (
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
)

// SessionModel is the database representation of a Session.
// This separates infrastructure concerns (db tags) from domain logic.
type SessionModel struct {
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

// ToSession converts a database model to a domain Session entity.
func (m *SessionModel) ToSession() *session.Session {
	return session.UnmarshalSessionFromDatabase(
		m.SessionID,
		m.UserID,
		m.RefreshToken,
		m.UserAgent,
		m.ClientIP,
		m.IsBlocked,
		m.ExpiresAt,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// SessionModelFromSession converts a domain Session entity to a database model.
func SessionModelFromSession(s *session.Session) *SessionModel {
	return &SessionModel{
		SessionID:    s.SessionID(),
		UserID:       s.UserID(),
		RefreshToken: s.RefreshToken(),
		UserAgent:    s.UserAgent(),
		ClientIP:     s.ClientIP(),
		IsBlocked:    s.IsBlocked(),
		ExpiresAt:    s.ExpiresAt(),
		CreatedAt:    s.CreatedAt(),
		UpdatedAt:    s.UpdatedAt(),
	}
}
