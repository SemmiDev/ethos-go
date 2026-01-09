package user

import (
	"context"

	"github.com/google/uuid"
)

// UserWriter provides write operations for user data.
// Use this interface when you need to create or modify users.
type UserWriter interface {
	// Create stores a new user in the database.
	Create(ctx context.Context, user *User) error

	// Update modifies an existing user in the database.
	Update(ctx context.Context, user *User) error

	// Delete permanently removes a user from the database.
	Delete(ctx context.Context, userID uuid.UUID) error
}

// Repository combines all user repository interfaces.
// This is the full interface that adapters implement.
// Consumers should depend on the smallest interface they need.
type Repository interface {
	UserReader // FindByEmail, FindByID
	UserWriter // Create, Update, Delete
}
