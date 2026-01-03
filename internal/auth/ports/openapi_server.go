package ports

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/auth/app/command"
	"github.com/semmidev/ethos-go/internal/auth/app/query"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	commonAuth "github.com/semmidev/ethos-go/internal/common/auth"
	"github.com/semmidev/ethos-go/internal/common/httputil"
	"github.com/semmidev/ethos-go/internal/common/model"
	"github.com/semmidev/ethos-go/internal/generated/api/auth"
)

type AuthOpenAPIServer struct {
	registerHandler           command.RegisterHandler
	loginHandler              command.LoginHandler
	logoutHandler             command.LogoutHandler
	logoutAllHandler          command.LogoutAllHandler
	listSessionsHandler       query.ListSessionsHandler
	getProfileHandler         query.GetProfileHandler
	updateProfileHandler      command.UpdateProfileHandler
	changePasswordHandler     command.ChangePasswordHandler
	verifyEmailHandler        command.VerifyEmailHandler
	resendVerificationHandler command.ResendVerificationHandler
	forgotPasswordHandler     command.ForgotPasswordHandler
	resetPasswordHandler      command.ResetPasswordHandler
	loginGoogleHandler        command.LoginGoogleHandler
	getGoogleAuthURLHandler   query.GetGoogleAuthURLHandler
	revokeSessionsHandler     command.RevokeAllOtherSessionsHandler
	deleteAccountHandler      command.DeleteAccountHandler
	exportDataHandler         query.ExportUserDataHandler
}

// Ensure AuthOpenAPIServer implements auth.ServerInterface
var _ auth.ServerInterface = (*AuthOpenAPIServer)(nil)

func NewAuthOpenAPIServer(
	registerHandler command.RegisterHandler,
	loginHandler command.LoginHandler,
	logoutHandler command.LogoutHandler,
	logoutAllHandler command.LogoutAllHandler,
	listSessionsHandler query.ListSessionsHandler,
	getProfileHandler query.GetProfileHandler,
	updateProfileHandler command.UpdateProfileHandler,
	changePasswordHandler command.ChangePasswordHandler,
	verifyEmailHandler command.VerifyEmailHandler,
	resendVerificationHandler command.ResendVerificationHandler,
	forgotPasswordHandler command.ForgotPasswordHandler,
	resetPasswordHandler command.ResetPasswordHandler,
	loginGoogleHandler command.LoginGoogleHandler,
	getGoogleAuthURLHandler query.GetGoogleAuthURLHandler,
	revokeSessionsHandler command.RevokeAllOtherSessionsHandler,
	deleteAccountHandler command.DeleteAccountHandler,
	exportDataHandler query.ExportUserDataHandler,
) *AuthOpenAPIServer {
	return &AuthOpenAPIServer{
		registerHandler:           registerHandler,
		loginHandler:              loginHandler,
		logoutHandler:             logoutHandler,
		logoutAllHandler:          logoutAllHandler,
		listSessionsHandler:       listSessionsHandler,
		getProfileHandler:         getProfileHandler,
		updateProfileHandler:      updateProfileHandler,
		changePasswordHandler:     changePasswordHandler,
		verifyEmailHandler:        verifyEmailHandler,
		resendVerificationHandler: resendVerificationHandler,
		forgotPasswordHandler:     forgotPasswordHandler,
		resetPasswordHandler:      resetPasswordHandler,
		loginGoogleHandler:        loginGoogleHandler,
		getGoogleAuthURLHandler:   getGoogleAuthURLHandler,
		revokeSessionsHandler:     revokeSessionsHandler,
		deleteAccountHandler:      deleteAccountHandler,
		exportDataHandler:         exportDataHandler,
	}
}

// List all user sessions
// (GET /auth/sessions)
func (s *AuthOpenAPIServer) ListSessions(w http.ResponseWriter, r *http.Request, params auth.ListSessionsParams) {
	// Get user from context
	user, err := commonAuth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	// Get current session ID
	currentSessionID, ok := GetSessionIDFromContext(r.Context())
	if !ok {
		httputil.Error(w, r, apperror.Unauthorized("session not found"))
		return
	}

	// Parse filter from query parameters
	filter := model.FilterFromRequest(r)

	// Handle specific session filters
	includeBlocked := false
	if params.IncludeBlocked != nil {
		includeBlocked = *params.IncludeBlocked
	}
	includeExpired := false
	if params.IncludeExpired != nil {
		includeExpired = *params.IncludeExpired
	}

	// Execute query
	q := query.ListSessionsQuery{
		UserID:           user.UserID,
		CurrentSessionID: currentSessionID,
		IncludeBlocked:   includeBlocked,
		IncludeExpired:   includeExpired,
		Filter:           filter,
	}

	result, err := s.listSessionsHandler.Handle(r.Context(), q)
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	// Map DTOs to response
	sessionsList := make([]auth.Session, 0, len(result.Sessions))
	for _, sess := range result.Sessions {
		sessionID, _ := uuid.Parse(sess.SessionID)

		sessionsList = append(sessionsList, auth.Session{
			SessionId: &sessionID,
			UserAgent: &sess.UserAgent,
			ClientIp:  &sess.ClientIP,
			IsBlocked: &sess.IsBlocked,
			ExpiresAt: &sess.ExpiresAt,
			CreatedAt: &sess.CreatedAt,
			IsActive:  &sess.IsActive,
			IsCurrent: &sess.IsCurrent,
		})
	}

	// Return paginated response
	httputil.SuccessPaginated(w, r, sessionsList, result.Pagination, "Sessions retrieved successfully")
}

// Register a new user
// (POST /auth/register)
func (s *AuthOpenAPIServer) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.RegisterCommand{
		Name:     req.Name,
		Email:    string(req.Email),
		Password: req.Password,
	}

	result, err := s.registerHandler.Handle(r.Context(), cmd)
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Created(w, r, map[string]interface{}{
		"user_id": result.UserID.String(),
		"email":   result.Email,
		"name":    result.Name,
	}, "User registered successfully")
}

// Authenticate user
// (POST /auth/login)
func (s *AuthOpenAPIServer) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	userAgent := r.UserAgent()
	if userAgent == "" {
		userAgent = "Unknown"
	}

	clientIP := r.RemoteAddr

	cmd := command.LoginCommand{
		Email:     string(req.Email),
		Password:  req.Password,
		UserAgent: userAgent,
		ClientIP:  clientIP,
	}

	result, err := s.loginHandler.Handle(r.Context(), cmd)
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"session_id":    result.SessionID,
		"user_id":       result.UserID,
		"expires_at":    result.ExpiresAt,
	}, "Login successful")
}

// Logout current session
// (POST /auth/logout)
func (s *AuthOpenAPIServer) Logout(w http.ResponseWriter, r *http.Request) {
	var req auth.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.LogoutCommand{
		SessionID: req.SessionId.String(),
	}

	if err := s.logoutHandler.Handle(r.Context(), cmd); err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "logged out successfully")
}

// Logout from all devices
// (POST /auth/logout-all)
func (s *AuthOpenAPIServer) LogoutAll(w http.ResponseWriter, r *http.Request) {
	var req auth.LogoutAllRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.LogoutAllCommand{
		UserID: req.UserId.String(),
	}

	if err := s.logoutAllHandler.Handle(r.Context(), cmd); err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "logged out from all devices successfully")
}

// Get current user profile
// (GET /auth/profile)
func (s *AuthOpenAPIServer) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := commonAuth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	result, err := s.getProfileHandler.Handle(r.Context(), query.GetProfileQuery{
		UserID: user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, map[string]interface{}{
		"user_id":    result.UserID,
		"name":       result.Name,
		"email":      result.Email,
		"timezone":   result.Timezone,
		"created_at": result.CreatedAt,
	}, "Profile retrieved successfully")
}

// Update user profile
// (PUT /auth/profile)
func (s *AuthOpenAPIServer) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	user, err := commonAuth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var req auth.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	var email *string
	if req.Email != nil {
		e := string(*req.Email)
		email = &e
	}

	result, err := s.updateProfileHandler.Handle(r.Context(), command.UpdateProfileCommand{
		UserID:   user.UserID,
		Name:     req.Name,
		Email:    email,
		Timezone: req.Timezone,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, map[string]interface{}{
		"user_id":    result.UserID,
		"name":       result.Name,
		"email":      result.Email,
		"timezone":   result.Timezone,
		"created_at": result.CreatedAt,
	}, "Profile updated successfully")
}

// Change user password
// (POST /auth/change-password)
func (s *AuthOpenAPIServer) ChangePassword(w http.ResponseWriter, r *http.Request) {
	user, err := commonAuth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var req auth.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	err = s.changePasswordHandler.Handle(r.Context(), command.ChangePasswordCommand{
		UserID:          user.UserID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Password changed successfully")
}

// Verify email address
// (POST /auth/verify-email)
func (s *AuthOpenAPIServer) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req auth.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.VerifyEmailCommand{
		Email: string(req.Email),
		Code:  req.Code,
	}

	if err := s.verifyEmailHandler.Handle(r.Context(), cmd); err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Email verified successfully")
}

// Resend verification email
// (POST /auth/resend-verification)
func (s *AuthOpenAPIServer) ResendVerification(w http.ResponseWriter, r *http.Request) {
	var req auth.ResendVerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.ResendVerificationCommand{
		Email: string(req.Email),
	}

	if err := s.resendVerificationHandler.Handle(r.Context(), cmd); err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Verification email sent")
}

// Request password reset
// (POST /auth/forgot-password)
func (s *AuthOpenAPIServer) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.ForgotPasswordCommand{
		Email: string(req.Email),
	}

	if err := s.forgotPasswordHandler.Handle(r.Context(), cmd); err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Password reset email sent")
}

// Reset password
// (POST /auth/reset-password)
func (s *AuthOpenAPIServer) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req auth.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.ResetPasswordCommand{
		Email:       string(req.Email),
		Code:        req.Code,
		NewPassword: req.NewPassword,
	}

	if err := s.resetPasswordHandler.Handle(r.Context(), cmd); err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Password reset successfully")
}

// googleLogin handles google login request
// (GET /auth/google/login)
func (s *AuthOpenAPIServer) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := "state-token" // TODO: Implement CSRF protection properly

	url, err := s.getGoogleAuthURLHandler.Handle(r.Context(), query.GetGoogleAuthURLQuery{State: state})
	if err != nil {
		httputil.Error(w, r, apperror.InternalError(err))
		return
	}

	httputil.Success(w, r, map[string]string{
		"url": url,
	}, "Google Login URL generated")
}

// googleCallback handles google oauth callback
// (POST /auth/google/callback)
func (s *AuthOpenAPIServer) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Code string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	userAgent := r.UserAgent()
	if userAgent == "" {
		userAgent = "Unknown"
	}
	clientIP := r.RemoteAddr

	cmd := command.LoginGoogleCommand{
		Code:      body.Code,
		UserAgent: userAgent,
		ClientIP:  clientIP,
	}

	result, err := s.loginGoogleHandler.Handle(r.Context(), cmd)
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, map[string]interface{}{
		"access_token":  result.AccessToken,
		"refresh_token": result.RefreshToken,
		"session_id":    result.SessionID,
		"user_id":       result.UserID,
		"expires_at":    result.ExpiresAt,
	}, "Login successful")
}

// Revoke all other sessions
// (DELETE /auth/sessions/other)
func (s *AuthOpenAPIServer) RevokeOtherSessions(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, err := commonAuth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	// Get current session ID
	currentSessionID, ok := GetSessionIDFromContext(r.Context())
	if !ok {
		httputil.Error(w, r, apperror.Unauthorized("session not found"))
		return
	}

	cmd := command.RevokeAllOtherSessionsCommand{
		UserID:           user.UserID,
		CurrentSessionID: currentSessionID,
	}

	result, err := s.revokeSessionsHandler.Handle(r.Context(), cmd)
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, map[string]interface{}{
		"revoked_count": result.RevokedCount,
	}, "Other sessions revoked successfully")
}

// Export user data (GDPR)
// (GET /auth/export)
func (s *AuthOpenAPIServer) ExportUserData(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, err := commonAuth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	result, err := s.exportDataHandler.Handle(r.Context(), query.ExportUserDataQuery{
		UserID: user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, result, "User data exported successfully")
}

// Delete user account
// (DELETE /auth/account)
func (s *AuthOpenAPIServer) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, err := commonAuth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var req struct {
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	cmd := command.DeleteAccountCommand{
		UserID:          user.UserID,
		Password:        req.Password,
		ConfirmDeletion: true,
	}

	if err := s.deleteAccountHandler.Handle(r.Context(), cmd); err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Account deleted successfully")
}
