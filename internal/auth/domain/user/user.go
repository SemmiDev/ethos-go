package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system
type User struct {
	UserID                 uuid.UUID  `db:"user_id"`
	Email                  string     `db:"email"`
	Name                   string     `db:"name"`
	HashedPassword         *string    `db:"hashed_password"`
	AuthProvider           string     `db:"auth_provider"`
	AuthProviderID         *string    `db:"auth_provider_id"`
	Timezone               string     `db:"timezone"`
	IsActive               bool       `db:"is_active"`
	IsVerified             bool       `db:"is_verified"`
	VerifyToken            *string    `db:"verify_token"`
	VerifyExpiresAt        *time.Time `db:"verify_expires_at"`
	PasswordResetToken     *string    `db:"password_reset_token"`
	PasswordResetExpiresAt *time.Time `db:"password_reset_expires_at"`
	CreatedAt              time.Time  `db:"created_at"`
	UpdatedAt              time.Time  `db:"updated_at"`
}

// NewUser creates a new user (factory constructor)
func NewUser(userID uuid.UUID, email, name, hashedPassword string) *User {
	now := time.Now()
	// Create pointer for string
	pwd := hashedPassword
	return &User{
		UserID:         userID,
		Email:          email,
		Name:           name,
		HashedPassword: &pwd,
		AuthProvider:   "email",
		AuthProviderID: nil,
		Timezone:       "Asia/Jakarta", // Default timezone
		IsActive:       true,
		IsVerified:     false,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// NewGoogleUser creates a new user from Google Auth
func NewGoogleUser(userID uuid.UUID, email, name, googleID string) *User {
	now := time.Now()
	providerID := googleID
	return &User{
		UserID:         userID,
		Email:          email,
		Name:           name,
		HashedPassword: nil,
		AuthProvider:   "google",
		AuthProviderID: &providerID,
		Timezone:       "Asia/Jakarta", // Default
		IsActive:       true,
		IsVerified:     true, // Google users are verified implicitly
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}
