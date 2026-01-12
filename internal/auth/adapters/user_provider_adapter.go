package adapters

import (
	"context"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/ports"
)

// UserProviderAdapter implements ports.UserProvider using the Auth module's UserReader.
// This adapter allows other modules to access user data via a shared interface
// without directly depending on the Auth module's internal types.
type UserProviderAdapter struct {
	userReader user.UserReader
}

// NewUserProviderAdapter creates a new UserProviderAdapter.
// This is created in the Auth module's service bootstrap and can be passed
// to other modules (like Notifications) that need user information.
func NewUserProviderAdapter(userReader user.UserReader) *UserProviderAdapter {
	return &UserProviderAdapter{userReader: userReader}
}

// GetUserByID retrieves user information by ID.
// Implements ports.UserProvider interface.
func (a *UserProviderAdapter) GetUserByID(ctx context.Context, userID string) (*ports.UserInfo, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	u, err := a.userReader.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &ports.UserInfo{
		UserID:   u.UserID().String(),
		Email:    u.Email(),
		Name:     u.Name(),
		Timezone: u.Timezone(),
	}, nil
}

// GetUserByEmail retrieves user information by email.
// Implements ports.UserProvider interface.
func (a *UserProviderAdapter) GetUserByEmail(ctx context.Context, email string) (*ports.UserInfo, error) {
	u, err := a.userReader.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &ports.UserInfo{
		UserID:   u.UserID().String(),
		Email:    u.Email(),
		Name:     u.Name(),
		Timezone: u.Timezone(),
	}, nil
}

// Compile-time check that UserProviderAdapter implements ports.UserProvider
var _ ports.UserProvider = (*UserProviderAdapter)(nil)
