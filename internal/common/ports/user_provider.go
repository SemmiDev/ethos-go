// Package ports defines the interfaces that allow modules to communicate
// without direct dependencies. These interfaces follow the Dependency Inversion Principle.
package ports

import (
	"context"
)

// UserInfo contains minimal user information needed by other modules.
// This is an anti-corruption layer - it decouples modules from the full User entity.
type UserInfo struct {
	UserID   string
	Email    string
	Name     string
	Timezone string
}

// UserProvider is an interface that allows other modules to query user data
// without depending on the Auth module's internal implementation.
//
// Example usage:
//   - Notifications module needs user email to send notifications
//   - Habits module needs user timezone for scheduling
//
// The Auth module provides an implementation, but consumers only depend on this interface.
type UserProvider interface {
	// GetUserByID retrieves basic user information by ID.
	// Returns ErrUserNotFound if the user doesn't exist.
	GetUserByID(ctx context.Context, userID string) (*UserInfo, error)

	// GetUserByEmail retrieves basic user information by email.
	// Returns ErrUserNotFound if the user doesn't exist.
	GetUserByEmail(ctx context.Context, email string) (*UserInfo, error)
}
