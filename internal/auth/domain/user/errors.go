package user

import "errors"

// Domain errors
var (
	ErrNotFound      = errors.New("user not found")
	ErrAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail  = errors.New("invalid email format")
)
