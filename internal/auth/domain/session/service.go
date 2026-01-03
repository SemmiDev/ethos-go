package session

import "time"

// AuthenticationService handles domain logic for authentication
type AuthenticationService struct {
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAuthenticationService(accessTTL, refreshTTL time.Duration) *AuthenticationService {
	return &AuthenticationService{
		accessTokenTTL:  accessTTL,
		refreshTokenTTL: refreshTTL,
	}
}

func (s *AuthenticationService) AccessTokenTTL() time.Duration {
	return s.accessTokenTTL
}

func (s *AuthenticationService) RefreshTokenTTL() time.Duration {
	return s.refreshTokenTTL
}
