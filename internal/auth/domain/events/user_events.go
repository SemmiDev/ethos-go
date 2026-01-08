package events

import (
	"time"

	commonevents "github.com/semmidev/ethos-go/internal/common/events"
)

// Event subjects
const (
	UserRegisteredType     = "auth.user.registered"
	UserVerifiedType       = "auth.user.verified"
	PasswordChangedType    = "auth.user.password_changed"
	UserLoggedInType       = "auth.user.logged_in"
	PasswordResetRequested = "auth.user.password_reset_requested"
)

// UserRegistered is emitted when a new user registers
type UserRegistered struct {
	commonevents.BaseEvent
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AuthProvider string `json:"auth_provider"`
}

// NewUserRegistered creates a new UserRegistered event
func NewUserRegistered(userID, email, name, authProvider string) UserRegistered {
	return UserRegistered{
		BaseEvent:    commonevents.NewBaseEvent(UserRegisteredType, "user", userID),
		UserID:       userID,
		Email:        email,
		Name:         name,
		AuthProvider: authProvider,
	}
}

// UserVerified is emitted when a user verifies their email
type UserVerified struct {
	commonevents.BaseEvent
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	VerifiedAt time.Time `json:"verified_at"`
}

// NewUserVerified creates a new UserVerified event
func NewUserVerified(userID, email string) UserVerified {
	return UserVerified{
		BaseEvent:  commonevents.NewBaseEvent(UserVerifiedType, "user", userID),
		UserID:     userID,
		Email:      email,
		VerifiedAt: time.Now().UTC(),
	}
}

// PasswordChanged is emitted when a user changes their password
type PasswordChanged struct {
	commonevents.BaseEvent
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	ChangedAt time.Time `json:"changed_at"`
}

// NewPasswordChanged creates a new PasswordChanged event
func NewPasswordChanged(userID, email string) PasswordChanged {
	return PasswordChanged{
		BaseEvent: commonevents.NewBaseEvent(PasswordChangedType, "user", userID),
		UserID:    userID,
		Email:     email,
		ChangedAt: time.Now().UTC(),
	}
}

// UserLoggedIn is emitted when a user logs in
type UserLoggedIn struct {
	commonevents.BaseEvent
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
}

// NewUserLoggedIn creates a new UserLoggedIn event
func NewUserLoggedIn(userID, email, userAgent, clientIP string) UserLoggedIn {
	return UserLoggedIn{
		BaseEvent: commonevents.NewBaseEvent(UserLoggedInType, "user", userID),
		UserID:    userID,
		Email:     email,
		UserAgent: userAgent,
		ClientIP:  clientIP,
	}
}
