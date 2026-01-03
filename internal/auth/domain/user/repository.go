package user

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines user persistence operations
type Repository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, userID uuid.UUID) error
}
