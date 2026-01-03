package app

import (
	"net/http"

	"github.com/semmidev/ethos-go/internal/auth/app/command"
	"github.com/semmidev/ethos-go/internal/auth/app/query"
)

// Application is the main application service facade for the auth module
type Application struct {
	Commands       Commands
	Queries        Queries
	AuthMiddleware func(http.Handler) http.Handler
}

// Commands groups all command handlers (write operations)
type Commands struct {
	Register           command.RegisterHandler
	Login              command.LoginHandler
	Logout             command.LogoutHandler
	LogoutAll          command.LogoutAllHandler
	RefreshToken       command.RefreshTokenHandler
	UpdateProfile      command.UpdateProfileHandler
	ChangePassword     command.ChangePasswordHandler
	VerifyEmail        command.VerifyEmailHandler
	ResendVerification command.ResendVerificationHandler
	ForgotPassword     command.ForgotPasswordHandler
	ResetPassword      command.ResetPasswordHandler
	LoginGoogle        command.LoginGoogleHandler
	RevokeSessions     command.RevokeAllOtherSessionsHandler
	DeleteAccount      command.DeleteAccountHandler
}

// Queries groups all query handlers (read operations)
type Queries struct {
	GetSession       query.GetSessionHandler
	ListSessions     query.ListSessionsHandler
	GetProfile       query.GetProfileHandler
	GetGoogleAuthURL query.GetGoogleAuthURLHandler
	ExportUserData   query.ExportUserDataHandler
}
