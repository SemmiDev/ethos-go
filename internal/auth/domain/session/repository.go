package session

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines what operations our domain needs for session persistence.
// Notice this is an interface in the domain layer - the domain defines what it needs,
// and the infrastructure layer provides the implementation. This is the Dependency
// Inversion Principle in action.
type Repository interface {
	// Create stores a new session in the database.
	// Returns an error if the session already exists or if there's a database problem.
	Create(ctx context.Context, session *Session) error

	// FindByID retrieves a session by its unique identifier.
	// Returns ErrSessionNotFound if no session exists with that ID.
	FindByID(ctx context.Context, sessionID uuid.UUID) (*Session, error)

	// FindByRefreshToken looks up a session using its refresh token.
	// This is used during token refresh to locate the session.
	// Returns ErrSessionNotFound if no session has that token.
	FindByRefreshToken(ctx context.Context, refreshToken string) (*Session, error)

	// FindAllByUserID returns all sessions for a specific user.
	// This allows users to see all their active login sessions across devices.
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]*Session, error)

	// Update modifies an existing session in the database.
	// This is used when refreshing tokens or blocking sessions.
	Update(ctx context.Context, session *Session) error

	// Delete permanently removes a session from the database.
	// This is used during logout to invalidate the session.
	Delete(ctx context.Context, sessionID uuid.UUID) error

	// DeleteAllByUserID removes all sessions for a user.
	// This is useful for "logout everywhere" functionality.
	DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired removes all sessions that have passed their expiration time.
	// This cleanup operation should be run periodically by a background worker.
	DeleteExpired(ctx context.Context) (int64, error)
}
