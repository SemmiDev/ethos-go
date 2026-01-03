package user

import (
	"context"

	"github.com/google/uuid"
)

// UserReader provides read-only access to user data.
// The auth module needs to look up users during login but shouldn't modify them.
// That's the responsibility of the user module.
type UserReader interface {
	// FindByEmail looks up a user by their email address.
	// Returns an error if no user is found with that email.
	FindByEmail(ctx context.Context, email string) (*User, error)

	// FindByID looks up a user by their unique identifier.
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
}
