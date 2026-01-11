package adapters

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/domain/service"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	authctx "github.com/semmidev/ethos-go/internal/auth/infrastructure/context"
	"github.com/semmidev/ethos-go/internal/auth/infrastructure/token"
)

// TokenVerifier is the interface for verifying access tokens
type TokenVerifier interface {
	VerifyAccessToken(ctx context.Context, tokenString string) (*service.TokenClaims, error)
}

// UserFinder is an interface for finding users by ID
type UserFinder interface {
	FindByID(ctx context.Context, userID uuid.UUID) (*user.User, error)
}

// AuthService implements the AuthServiceInterface for gRPC authentication
type AuthService struct {
	tokenVerifier TokenVerifier
	userRepo      UserFinder
}

// NewAuthService creates a new AuthService
func NewAuthService(tokenVerifier TokenVerifier, userRepo UserFinder) *AuthService {
	return &AuthService{
		tokenVerifier: tokenVerifier,
		userRepo:      userRepo,
	}
}

// ValidateToken validates a JWT token and returns its payload
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*token.Payload, error) {
	claims, err := s.tokenVerifier.VerifyAccessToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	return &token.Payload{
		UserID:    claims.UserID,
		SessionID: claims.SessionID,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}, nil
}

// GetUserByID retrieves a user by ID and returns auth context user
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (authctx.User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return authctx.User{}, err
	}

	u, err := s.userRepo.FindByID(ctx, uid)
	if err != nil {
		return authctx.User{}, err
	}

	return authctx.User{
		UserID: userID,
		Email:  u.Email(),
	}, nil
}
