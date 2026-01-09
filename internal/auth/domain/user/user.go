package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system
// Fields are private to enforce encapsulation - use getters for read access
type User struct {
	userID                 uuid.UUID
	email                  string
	name                   string
	hashedPassword         *string
	authProvider           string
	authProviderID         *string
	timezone               string
	isActive               bool
	isVerified             bool
	verifyToken            *string
	verifyExpiresAt        *time.Time
	passwordResetToken     *string
	passwordResetExpiresAt *time.Time
	createdAt              time.Time
	updatedAt              time.Time
}

// Getters for User fields

func (u *User) UserID() uuid.UUID                  { return u.userID }
func (u *User) Email() string                      { return u.email }
func (u *User) Name() string                       { return u.name }
func (u *User) HashedPassword() *string            { return u.hashedPassword }
func (u *User) AuthProvider() string               { return u.authProvider }
func (u *User) AuthProviderID() *string            { return u.authProviderID }
func (u *User) Timezone() string                   { return u.timezone }
func (u *User) IsActive() bool                     { return u.isActive }
func (u *User) IsVerified() bool                   { return u.isVerified }
func (u *User) VerifyToken() *string               { return u.verifyToken }
func (u *User) VerifyExpiresAt() *time.Time        { return u.verifyExpiresAt }
func (u *User) PasswordResetToken() *string        { return u.passwordResetToken }
func (u *User) PasswordResetExpiresAt() *time.Time { return u.passwordResetExpiresAt }
func (u *User) CreatedAt() time.Time               { return u.createdAt }
func (u *User) UpdatedAt() time.Time               { return u.updatedAt }

// Setters for mutable fields (business operations)

func (u *User) SetEmail(email string) {
	u.email = email
	u.updatedAt = time.Now()
}

func (u *User) SetName(name string) {
	u.name = name
	u.updatedAt = time.Now()
}

func (u *User) SetHashedPassword(hashedPassword string) {
	u.hashedPassword = &hashedPassword
	u.updatedAt = time.Now()
}

func (u *User) SetTimezone(timezone string) {
	u.timezone = timezone
	u.updatedAt = time.Now()
}

func (u *User) SetVerifyToken(token *string, expiresAt *time.Time) {
	u.verifyToken = token
	u.verifyExpiresAt = expiresAt
	u.updatedAt = time.Now()
}

func (u *User) SetPasswordResetToken(token *string, expiresAt *time.Time) {
	u.passwordResetToken = token
	u.passwordResetExpiresAt = expiresAt
	u.updatedAt = time.Now()
}

func (u *User) MarkVerified() {
	u.isVerified = true
	u.verifyToken = nil
	u.verifyExpiresAt = nil
	u.updatedAt = time.Now()
}

func (u *User) Deactivate() {
	u.isActive = false
	u.updatedAt = time.Now()
}

func (u *User) Activate() {
	u.isActive = true
	u.updatedAt = time.Now()
}

func (u *User) SetAuthProvider(provider string, providerID *string) {
	u.authProvider = provider
	u.authProviderID = providerID
	u.updatedAt = time.Now()
}

// NewUser creates a new user (factory constructor)
func NewUser(userID uuid.UUID, email, name, hashedPassword string) *User {
	now := time.Now()
	pwd := hashedPassword
	return &User{
		userID:         userID,
		email:          email,
		name:           name,
		hashedPassword: &pwd,
		authProvider:   "email",
		authProviderID: nil,
		timezone:       "Asia/Jakarta", // Default timezone
		isActive:       true,
		isVerified:     false,
		createdAt:      now,
		updatedAt:      now,
	}
}

// NewGoogleUser creates a new user from Google Auth
func NewGoogleUser(userID uuid.UUID, email, name, googleID string) *User {
	now := time.Now()
	providerID := googleID
	return &User{
		userID:         userID,
		email:          email,
		name:           name,
		hashedPassword: nil,
		authProvider:   "google",
		authProviderID: &providerID,
		timezone:       "Asia/Jakarta", // Default
		isActive:       true,
		isVerified:     true, // Google users are verified implicitly
		createdAt:      now,
		updatedAt:      now,
	}
}

// UnmarshalUserFromDatabase reconstructs a User from database fields
// This is used by the adapter layer to convert from database model to domain entity
func UnmarshalUserFromDatabase(
	userID uuid.UUID,
	email, name string,
	hashedPassword *string,
	authProvider string,
	authProviderID *string,
	timezone string,
	isActive, isVerified bool,
	verifyToken *string,
	verifyExpiresAt *time.Time,
	passwordResetToken *string,
	passwordResetExpiresAt *time.Time,
	createdAt, updatedAt time.Time,
) *User {
	return &User{
		userID:                 userID,
		email:                  email,
		name:                   name,
		hashedPassword:         hashedPassword,
		authProvider:           authProvider,
		authProviderID:         authProviderID,
		timezone:               timezone,
		isActive:               isActive,
		isVerified:             isVerified,
		verifyToken:            verifyToken,
		verifyExpiresAt:        verifyExpiresAt,
		passwordResetToken:     passwordResetToken,
		passwordResetExpiresAt: passwordResetExpiresAt,
		createdAt:              createdAt,
		updatedAt:              updatedAt,
	}
}
