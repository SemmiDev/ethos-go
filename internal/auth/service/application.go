package service

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/config"
	"github.com/semmidev/ethos-go/internal/auth/adapters"
	"github.com/semmidev/ethos-go/internal/auth/adapters/google"
	"github.com/semmidev/ethos-go/internal/auth/app"
	"github.com/semmidev/ethos-go/internal/auth/app/command"
	"github.com/semmidev/ethos-go/internal/auth/app/query"
	"github.com/semmidev/ethos-go/internal/auth/domain/gateway"
	"github.com/semmidev/ethos-go/internal/auth/domain/session"
	"github.com/semmidev/ethos-go/internal/auth/ports"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
)

// NewApplication creates and wires all dependencies for the auth module
func NewApplication(
	_ context.Context,
	cfg *config.Config,
	db *sqlx.DB,
	dispatcher gateway.TaskDispatcher,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) app.Application {
	// Create adapters (infrastructure)
	userRepo := adapters.NewPostgresUserRepository(db)
	sessionRepo := adapters.NewPostgresSessionRepository(db)
	passwordHasher := adapters.NewBcryptPasswordHasher()
	tokenIssuer := adapters.NewJWTTokenIssuer(cfg)
	validate := validator.New("en")
	googleService := google.NewService(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleCallbackURL,
	)

	// Create domain services
	authService := session.NewAuthenticationService(
		time.Duration(cfg.AuthAccessTokenExpiry)*time.Minute,
		time.Duration(cfg.AuthRefreshTokenExpiry)*time.Minute,
	)

	// Create command and query handlers
	return app.Application{
		AuthMiddleware: ports.AuthMiddleware(tokenIssuer, userRepo),
		Commands: app.Commands{
			Register: command.NewRegisterHandler(
				userRepo,
				passwordHasher,
				validate,
				dispatcher,
				log,
				metricsClient,
			),
			Login: command.NewLoginHandler(
				sessionRepo,
				userRepo,
				passwordHasher,
				tokenIssuer,
				authService,
				validate,
				log,
				metricsClient,
			),
			Logout: command.NewLogoutHandler(
				sessionRepo,
				log,
				metricsClient,
			),
			LogoutAll: command.NewLogoutAllHandler(
				sessionRepo,
				log,
				metricsClient,
			),
			RefreshToken: command.NewRefreshTokenHandler(
				sessionRepo,
				tokenIssuer,
				authService,
				log,
				metricsClient,
			),
			UpdateProfile: command.NewUpdateProfileHandler(
				userRepo,
				log,
				metricsClient,
			),
			ChangePassword: command.NewChangePasswordHandler(
				userRepo,
				log,
				metricsClient,
			),
			VerifyEmail: command.NewVerifyEmailHandler(
				userRepo,
				validate,
				log,
				metricsClient,
			),
			ResendVerification: command.NewResendVerificationHandler(
				userRepo,
				validate,
				dispatcher,
				log,
				metricsClient,
			),
			ForgotPassword: command.NewForgotPasswordHandler(
				userRepo,
				validate,
				dispatcher,
				log,
				metricsClient,
			),
			ResetPassword: command.NewResetPasswordHandler(
				userRepo,
				passwordHasher,
				validate,
				log,
				metricsClient,
			),
			LoginGoogle: command.NewLoginGoogleHandler(
				googleService,
				userRepo,
				sessionRepo,
				tokenIssuer,
				authService,
				log,
				metricsClient,
			),
			RevokeSessions: command.NewRevokeAllOtherSessionsHandler(
				sessionRepo,
				log,
				metricsClient,
			),
			DeleteAccount: command.NewDeleteAccountHandler(
				userRepo,
				sessionRepo,
				log,
				metricsClient,
			),
		},
		Queries: app.Queries{
			GetSession: query.NewGetSessionHandler(
				sessionRepo,
				log,
				metricsClient,
			),
			ListSessions: query.NewListSessionsHandler(
				sessionRepo,
				log,
				metricsClient,
			),
			GetProfile: query.NewGetProfileHandler(
				userRepo,
				log,
				metricsClient,
			),
			GetGoogleAuthURL: query.NewGetGoogleAuthURLHandler(
				googleService,
				log,
				metricsClient,
			),
			ExportUserData: query.NewExportUserDataHandler(
				userRepo,
				db,
				log,
				metricsClient,
			),
		},
	}
}
