package command

import (
	"context"
	"time"

	authevents "github.com/semmidev/ethos-go/internal/auth/domain/events"
	"github.com/semmidev/ethos-go/internal/auth/domain/service"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/random"
	"github.com/semmidev/ethos-go/internal/common/validator"
)

// LoginCommand contains all the data needed to authenticate a user
type LoginCommand struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	UserAgent string `json:"user_agent"`
	ClientIP  string `json:"client_ip"`
}

// LoginResult contains everything the client needs after successful authentication
type LoginResult struct {
	AccessToken  string
	RefreshToken string
	SessionID    string
	UserID       string
	ExpiresAt    int64
}

// LoginHandler processes login commands
type LoginHandler decorator.CommandHandlerWithResult[LoginCommand, *LoginResult]

type loginHandler struct {
	sessionRepo    session.Repository
	userRepo       user.Repository
	passwordHasher service.PasswordHasher
	tokenIssuer    service.TokenIssuer
	authService    *session.AuthenticationService
	validator      *validator.Validator
	publisher      events.Publisher
}

func NewLoginHandler(
	sessionRepo session.Repository,
	userRepo user.Repository,
	passwordHasher service.PasswordHasher,
	tokenIssuer service.TokenIssuer,
	authService *session.AuthenticationService,
	validator *validator.Validator,
	publisher events.Publisher, // Injected publisher
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) LoginHandler {
	return decorator.ApplyCommandResultDecorators[LoginCommand, *LoginResult](
		loginHandler{
			sessionRepo:    sessionRepo,
			userRepo:       userRepo,
			passwordHasher: passwordHasher,
			tokenIssuer:    tokenIssuer,
			authService:    authService,
			validator:      validator,
			publisher:      publisher,
		},
		log,
		metricsClient,
	)
}

func (h loginHandler) Handle(ctx context.Context, cmd LoginCommand) (*LoginResult, error) {
	// Validate input
	if err := h.validator.Validate(cmd); err != nil {
		if validationErrors, ok := validator.GetValidationErrors(err); ok {
			details := make(map[string]interface{})
			for _, ve := range validationErrors {
				details[ve.Field] = ve.Message
			}
			return nil, apperror.ValidationFailedWithDetails("validation failed", details)
		}
		return nil, apperror.ValidationFailed(err.Error())
	}

	// Look up the user
	foundUser, err := h.userRepo.FindByEmail(ctx, cmd.Email)
	if err != nil {
		// Just generic error to avoid user enumeration
		return nil, apperror.InvalidCredentials(nil)
	}

	// Verify password
	if foundUser.HashedPassword == nil {
		return nil, apperror.InvalidCredentials(nil)
	}
	passwordMatches, err := h.passwordHasher.Compare(ctx, *foundUser.HashedPassword, cmd.Password)
	if err != nil {
		return nil, apperror.InternalError(err)
	}

	if !passwordMatches {
		return nil, apperror.InvalidCredentials(nil)
	}

	if !foundUser.IsVerified {
		return nil, apperror.Unauthorized("Please verify your email address")
	}

	// Calculate token expiration times
	now := time.Now()
	accessTokenExpiry := now.Add(h.authService.AccessTokenTTL())
	refreshTokenExpiry := now.Add(h.authService.RefreshTokenTTL())

	// Issue access token
	accessToken, err := h.tokenIssuer.IssueAccessToken(ctx, foundUser.UserID, accessTokenExpiry)
	if err != nil {
		return nil, apperror.InternalError(err)
	}

	// Issue refresh token
	sessionID := random.NewUUID()
	refreshToken, err := h.tokenIssuer.IssueRefreshToken(ctx, sessionID, refreshTokenExpiry)
	if err != nil {
		return nil, apperror.InternalError(err)
	}

	// Create session
	newSession := session.NewSession(
		sessionID,
		foundUser.UserID,
		refreshToken,
		cmd.UserAgent,
		cmd.ClientIP,
		refreshTokenExpiry,
	)

	// Persist the session
	if err := h.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, apperror.DatabaseError("create session", err)
	}

	// Publish UserLoggedIn event
	event := authevents.NewUserLoggedIn(
		foundUser.UserID.String(),
		foundUser.Email,
		cmd.UserAgent,
		cmd.ClientIP,
	)
	_ = h.publisher.Publish(ctx, event)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		SessionID:    sessionID.String(),
		UserID:       foundUser.UserID.String(),
		ExpiresAt:    accessTokenExpiry.Unix(),
	}, nil
}
