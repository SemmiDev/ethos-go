package user

import (
	"errors"
	"strings"
	"time"
)

// Timezone is a value object representing a validated timezone.
// It ensures timezones can be loaded by the time package.
type Timezone struct {
	value string
}

// ErrInvalidTimezone indicates the timezone is not valid/loadable.
var ErrInvalidTimezone = errors.New("invalid timezone")

// DefaultTimezone is the default timezone used when none is specified.
const DefaultTimezone = "Asia/Jakarta"

// NewTimezone creates a new Timezone value object with validation.
// Returns an error if the timezone cannot be loaded.
func NewTimezone(raw string) (Timezone, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return NewTimezone(DefaultTimezone)
	}

	// Validate the timezone can be loaded
	_, err := time.LoadLocation(raw)
	if err != nil {
		return Timezone{}, ErrInvalidTimezone
	}

	return Timezone{value: raw}, nil
}

// MustNewTimezone creates a new Timezone, panicking if invalid.
// Use only for testing or when the value is guaranteed valid.
func MustNewTimezone(raw string) Timezone {
	t, err := NewTimezone(raw)
	if err != nil {
		panic(err)
	}
	return t
}

// UnmarshalTimezoneFromDatabase reconstructs a Timezone from database storage.
// This trusts that the database value was validated on insert.
func UnmarshalTimezoneFromDatabase(value string) Timezone {
	return Timezone{value: value}
}

// String returns the timezone name.
func (t Timezone) String() string {
	return t.value
}

// Location returns the time.Location for this timezone.
// Returns UTC if the timezone cannot be loaded (should not happen if properly validated).
func (t Timezone) Location() *time.Location {
	loc, err := time.LoadLocation(t.value)
	if err != nil {
		return time.UTC
	}
	return loc
}

// Equals checks if two timezones are the same.
func (t Timezone) Equals(other Timezone) bool {
	return t.value == other.value
}

// IsDefault returns true if this is the default timezone.
func (t Timezone) IsDefault() bool {
	return t.value == DefaultTimezone
}
