package adapters

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// BcryptPasswordHasher implements PasswordHasher using bcrypt.
//
// We return (false, nil) for mismatched passwords so callers can treat it as
// an authentication failure, while other bcrypt errors bubble up.
type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher() *BcryptPasswordHasher {
	return &BcryptPasswordHasher{cost: bcrypt.DefaultCost}
}

func (h *BcryptPasswordHasher) Hash(ctx context.Context, password string) (string, error) {
	_ = ctx
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *BcryptPasswordHasher) Compare(ctx context.Context, hashedPassword, plainPassword string) (bool, error) {
	_ = ctx
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
