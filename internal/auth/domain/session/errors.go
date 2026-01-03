package session

import "errors"

// Domain errors
var (
	ErrNotFound       = errors.New("session not found")
	ErrSessionExpired = errors.New("session expired")
	ErrSessionBlocked = errors.New("session blocked")
)
