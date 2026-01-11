package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Token interface {
	// GenerateToken membuat token baru untuk user ID tertentu.
	GenerateToken(userID uuid.UUID, exp time.Duration) (*GenerateTokenResponse, error)
	// VerifyBearerToken memverifikasi token dari header otorisasi "Bearer".
	VerifyBearerToken(authHeader string) (*Payload, error)
	// VerifyToken memverifikasi token string dan mengembalikan payload-nya.
	VerifyToken(token string) (*Payload, error)
}

type GenerateTokenResponse struct {
	Value     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Payload struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

var (
	ErrInvalidTokenFormat = errors.New("token format tidak valid")
	ErrInvalidToken       = errors.New("token tidak valid")
	ErrExpiredToken       = errors.New("token sudah kedaluwarsa")
)
