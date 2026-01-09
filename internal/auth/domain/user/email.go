package user

import (
	"errors"
	"regexp"
	"strings"
)

// Email is a value object representing a validated email address.
// It ensures emails are always lowercase and valid.
type Email struct {
	value string
}

// emailRegex provides basic email format validation.
// More comprehensive validation should happen at the email delivery layer.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ErrInvalidEmailFormat indicates the email format is invalid.
var ErrInvalidEmailFormat = errors.New("invalid email format")

// NewEmail creates a new Email value object with validation.
// Returns an error if the email format is invalid.
func NewEmail(raw string) (Email, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Email{}, ErrInvalidEmailFormat
	}

	if !emailRegex.MatchString(raw) {
		return Email{}, ErrInvalidEmailFormat
	}

	return Email{value: strings.ToLower(raw)}, nil
}

// MustNewEmail creates a new Email, panicking if invalid.
// Use only for testing or when the value is guaranteed valid.
func MustNewEmail(raw string) Email {
	e, err := NewEmail(raw)
	if err != nil {
		panic(err)
	}
	return e
}

// UnmarshalEmailFromDatabase reconstructs an Email from database storage.
// This trusts that the database value was validated on insert.
func UnmarshalEmailFromDatabase(value string) Email {
	return Email{value: value}
}

// String returns the email address as a string.
func (e Email) String() string {
	return e.value
}

// Equals checks if two emails are the same.
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

// IsEmpty returns true if the email is empty.
func (e Email) IsEmpty() bool {
	return e.value == ""
}
