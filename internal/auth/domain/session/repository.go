package session

import (
	"context"

	"github.com/google/uuid"
)

// SessionReader provides read-only access to session data.
// Use this interface when you only need to query sessions.
type SessionReader interface {
	// FindByID retrieves a session by its unique identifier.
	FindByID(ctx context.Context, sessionID uuid.UUID) (*Session, error)

	// FindByRefreshToken looks up a session using its refresh token.
	FindByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)

	// FindAllByUserID returns all sessions for a specific user.
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Session, error)
}

// SessionWriter provides write operations for session data.
// Use this interface when you need to create or modify sessions.
type SessionWriter interface {
	// Create stores a new session in the database.
	Create(ctx context.Context, session *Session) error

	// Update modifies an existing session in the database.
	Update(ctx context.Context, session *Session) error

	// Delete permanently removes a session from the database.
	Delete(ctx context.Context, sessionID uuid.UUID) error
}

// SessionMaintainer provides maintenance operations for sessions.
// Use this interface for cleanup tasks and bulk operations.
type SessionMaintainer interface {
	// DeleteAllByUserID removes all sessions for a user.
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired removes all sessions that have passed their expiration time.
	DeleteExpired(ctx context.Context) (int64, error)
}

// Repository combines all session repository interfaces.
// This is the full interface that adapters implement.
// Consumers should depend on the smallest interface they need.
type Repository interface {
	SessionReader
	SessionWriter
	SessionMaintainer
}
