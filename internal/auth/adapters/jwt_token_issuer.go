package adapters

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/config"
	"github.com/semmidev/ethos-go/internal/auth/domain/service"
)

// JWTTokenIssuer issues signed JWTs for access and refresh tokens.
//
// NOTE: This is an auth-module specific token implementation and is independent
// from pkg/token.
type JWTTokenIssuer struct {
	secretKey []byte
	issuer    string
}

func NewJWTTokenIssuer(cfg *config.Config) *JWTTokenIssuer {
	return &JWTTokenIssuer{
		secretKey: []byte(cfg.AuthJWTSecret),
		issuer:    cfg.AppName,
	}
}

type accessTokenClaims struct {
	jwt.RegisteredClaims
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
}

type refreshTokenClaims struct {
	jwt.RegisteredClaims
	SessionID string `json:"session_id"`
	Type      string `json:"type"`
}

func (j *JWTTokenIssuer) IssueAccessToken(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, expiresAt time.Time) (string, error) {
	_ = ctx
	now := time.Now()

	claims := &accessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		UserID:    userID.String(),
		SessionID: sessionID.String(),
		Type:      "access",
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(j.secretKey)
}

func (j *JWTTokenIssuer) IssueRefreshToken(ctx context.Context, sessionID uuid.UUID, expiresAt time.Time) (string, error) {
	_ = ctx
	now := time.Now()

	claims := &refreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   sessionID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		SessionID: sessionID.String(),
		Type:      "refresh",
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(j.secretKey)
}

// VerifyAccessToken validates an access token and returns its claims.
func (j *JWTTokenIssuer) VerifyAccessToken(ctx context.Context, tokenString string) (*service.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*accessTokenClaims); ok && token.Valid {
		// Validate issuer
		if claims.Issuer != j.issuer {
			return nil, jwt.ErrTokenInvalidIssuer
		}

		// Parse UUIDs
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return nil, err
		}

		sessionID, err := uuid.Parse(claims.SessionID)
		if err != nil {
			return nil, err
		}

		return &service.TokenClaims{
			UserID:    userID,
			SessionID: sessionID,
			IssuedAt:  claims.IssuedAt.Time.Unix(),
			ExpiresAt: claims.ExpiresAt.Time.Unix(),
		}, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

// VerifyRefreshToken validates a refresh token and returns its claims.
func (j *JWTTokenIssuer) VerifyRefreshToken(ctx context.Context, tokenString string) (*service.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*refreshTokenClaims); ok && token.Valid {
		// Validate issuer
		if claims.Issuer != j.issuer {
			return nil, jwt.ErrTokenInvalidIssuer
		}

		// Parse UUIDs
		sessionID, err := uuid.Parse(claims.SessionID)
		if err != nil {
			return nil, err
		}

		return &service.TokenClaims{
			SessionID: sessionID,
			IssuedAt:  claims.IssuedAt.Time.Unix(),
			ExpiresAt: claims.ExpiresAt.Time.Unix(),
		}, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
