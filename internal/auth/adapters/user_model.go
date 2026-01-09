package adapters

import (
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
)

// UserModel is the database representation of a User
// This DTO has db tags for sqlx scanning, keeping infrastructure concerns out of domain
type UserModel struct {
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

// ToUser converts the database model to a domain entity
func (m *UserModel) ToUser() *user.User {
	return user.UnmarshalUserFromDatabase(
		m.UserID,
		m.Email,
		m.Name,
		m.HashedPassword,
		m.AuthProvider,
		m.AuthProviderID,
		m.Timezone,
		m.IsActive,
		m.IsVerified,
		m.VerifyToken,
		m.VerifyExpiresAt,
		m.PasswordResetToken,
		m.PasswordResetExpiresAt,
		m.CreatedAt,
		m.UpdatedAt,
	)
}

// UserModelFromUser converts a domain entity to a database model
func UserModelFromUser(u *user.User) *UserModel {
	return &UserModel{
		UserID:                 u.UserID(),
		Email:                  u.Email(),
		Name:                   u.Name(),
		HashedPassword:         u.HashedPassword(),
		AuthProvider:           u.AuthProvider(),
		AuthProviderID:         u.AuthProviderID(),
		Timezone:               u.Timezone(),
		IsActive:               u.IsActive(),
		IsVerified:             u.IsVerified(),
		VerifyToken:            u.VerifyToken(),
		VerifyExpiresAt:        u.VerifyExpiresAt(),
		PasswordResetToken:     u.PasswordResetToken(),
		PasswordResetExpiresAt: u.PasswordResetExpiresAt(),
		CreatedAt:              u.CreatedAt(),
		UpdatedAt:              u.UpdatedAt(),
	}
}
