package ports

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/semmidev/ethos-go/internal/auth/app/command"
	"github.com/semmidev/ethos-go/internal/auth/app/query"
	authctx "github.com/semmidev/ethos-go/internal/auth/infrastructure/context"
	"github.com/semmidev/ethos-go/internal/common/grpcutil"
	"github.com/semmidev/ethos-go/internal/common/model"
	authv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/auth/v1"
	commonv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/common/v1"
)

// AuthGRPCServer implements the gRPC AuthService interface.
type AuthGRPCServer struct {
	authv1.UnimplementedAuthServiceServer
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

// NewAuthGRPCServer creates a new AuthGRPCServer.
func NewAuthGRPCServer(
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
) *AuthGRPCServer {
	return &AuthGRPCServer{
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

// Register creates a new user account.
func (s *AuthGRPCServer) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	cmd := command.RegisterCommand{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	result, err := s.registerHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		Data: &authv1.RegisterData{
			UserId: result.UserID.String(),
			Email:  result.Email,
			Name:   result.Name,
		},
	}, nil
}

// Login authenticates a user and returns tokens.
func (s *AuthGRPCServer) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	cmd := command.LoginCommand{
		Email:     req.Email,
		Password:  req.Password,
		UserAgent: "gRPC-Client",
		ClientIP:  "unknown",
	}

	result, err := s.loginHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.LoginResponse{
		Success: true,
		Data: &authv1.LoginData{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			SessionId:    result.SessionID,
			UserId:       result.UserID,
			ExpiresAt:    result.ExpiresAt,
		},
	}, nil
}

// GoogleLogin returns the Google OAuth login URL.
func (s *AuthGRPCServer) GoogleLogin(ctx context.Context, req *authv1.GoogleLoginRequest) (*authv1.GoogleLoginResponse, error) {
	state := "state-token"

	url, err := s.getGoogleAuthURLHandler.Handle(ctx, query.GetGoogleAuthURLQuery{State: state})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.GoogleLoginResponse{
		Success: true,
		Data: &authv1.GoogleLoginData{
			Url: url,
		},
	}, nil
}

// GoogleCallback handles the Google OAuth callback.
func (s *AuthGRPCServer) GoogleCallback(ctx context.Context, req *authv1.GoogleCallbackRequest) (*authv1.LoginResponse, error) {
	cmd := command.LoginGoogleCommand{
		Code:      req.Code,
		UserAgent: "gRPC-Client",
		ClientIP:  "unknown",
	}

	result, err := s.loginGoogleHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.LoginResponse{
		Success: true,
		Data: &authv1.LoginData{
			AccessToken:  result.AccessToken,
			RefreshToken: result.RefreshToken,
			SessionId:    result.SessionID,
			UserId:       result.UserID,
			ExpiresAt:    result.ExpiresAt,
		},
	}, nil
}

// Logout terminates the specified session.
func (s *AuthGRPCServer) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	cmd := command.LogoutCommand{
		SessionID: req.SessionId,
	}

	if err := s.logoutHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.LogoutResponse{
		Success: true,
		Message: "Logged out successfully",
	}, nil
}

// LogoutAll terminates all sessions for a user.
func (s *AuthGRPCServer) LogoutAll(ctx context.Context, req *authv1.LogoutAllRequest) (*authv1.LogoutResponse, error) {
	cmd := command.LogoutAllCommand{
		UserID: req.UserId,
	}

	if err := s.logoutAllHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.LogoutResponse{
		Success: true,
		Message: "Logged out from all devices successfully",
	}, nil
}

// ListSessions returns all sessions for the authenticated user.
func (s *AuthGRPCServer) ListSessions(ctx context.Context, req *authv1.ListSessionsRequest) (*authv1.ListSessionsResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	currentSessionID := user.SessionID

	filter := model.NewFilter()
	if req.Page > 0 {
		filter.CurrentPage = int(req.Page)
	}
	if req.PerPage > 0 {
		filter.PerPage = int(req.PerPage)
	}

	q := query.ListSessionsQuery{
		UserID:           user.UserID,
		CurrentSessionID: currentSessionID,
		IncludeBlocked:   req.IncludeBlocked,
		IncludeExpired:   req.IncludeExpired,
		Filter:           filter,
	}

	result, err := s.listSessionsHandler.Handle(ctx, q)
	if err != nil {
		return nil, toGRPCError(err)
	}

	sessions := make([]*authv1.Session, 0, len(result.Sessions))
	for _, sess := range result.Sessions {
		sessions = append(sessions, &authv1.Session{
			SessionId: sess.SessionID,
			UserAgent: sess.UserAgent,
			ClientIp:  sess.ClientIP,
			IsBlocked: sess.IsBlocked,
			ExpiresAt: timestamppb.New(sess.ExpiresAt),
			CreatedAt: timestamppb.New(sess.CreatedAt),
			IsActive:  sess.IsActive,
			IsCurrent: sess.IsCurrent,
		})
	}

	return &authv1.ListSessionsResponse{
		Success: true,
		Message: "Sessions retrieved successfully",
		Data:    sessions,
		Meta: &commonv1.Meta{
			Pagination: &commonv1.PaginationResponse{
				HasPreviousPage:        result.Pagination.HasPreviousPage,
				HasNextPage:            result.Pagination.HasNextPage,
				CurrentPage:            int32(result.Pagination.CurrentPage),
				PerPage:                int32(result.Pagination.PerPage),
				TotalData:              int32(result.Pagination.TotalData),
				TotalDataInCurrentPage: int32(result.Pagination.TotalDataInCurrentPage),
				LastPage:               int32(result.Pagination.LastPage),
				From:                   int32(result.Pagination.From),
				To:                     int32(result.Pagination.To),
			},
		},
	}, nil
}

// RevokeOtherSessions revokes all sessions except the current one.
func (s *AuthGRPCServer) RevokeOtherSessions(ctx context.Context, req *authv1.RevokeOtherSessionsRequest) (*authv1.RevokeOtherSessionsResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	currentSessionID := user.SessionID

	cmd := command.RevokeAllOtherSessionsCommand{
		UserID:           user.UserID,
		CurrentSessionID: currentSessionID,
	}

	result, err := s.revokeSessionsHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.RevokeOtherSessionsResponse{
		Success:      true,
		Message:      "Other sessions revoked successfully",
		RevokedCount: int32(result.RevokedCount),
	}, nil
}

// GetProfile retrieves the current user's profile.
func (s *AuthGRPCServer) GetProfile(ctx context.Context, req *authv1.GetProfileRequest) (*authv1.ProfileResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	result, err := s.getProfileHandler.Handle(ctx, query.GetProfileQuery{
		UserID: user.UserID,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.ProfileResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data: &authv1.ProfileData{
			UserId:    result.UserID,
			Name:      result.Name,
			Email:     result.Email,
			Timezone:  result.Timezone,
			CreatedAt: timestamppb.New(result.CreatedAt),
		},
	}, nil
}

// UpdateProfile updates the current user's profile.
func (s *AuthGRPCServer) UpdateProfile(ctx context.Context, req *authv1.UpdateProfileRequest) (*authv1.ProfileResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.UpdateProfileCommand{
		UserID:   user.UserID,
		Name:     req.Name,
		Email:    req.Email,
		Timezone: req.Timezone,
	}

	result, err := s.updateProfileHandler.Handle(ctx, cmd)
	if err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.ProfileResponse{
		Success: true,
		Message: "Profile updated successfully",
		Data: &authv1.ProfileData{
			UserId:    result.UserID,
			Name:      result.Name,
			Email:     result.Email,
			Timezone:  result.Timezone,
			CreatedAt: timestamppb.New(result.CreatedAt),
		},
	}, nil
}

// ChangePassword changes the user's password.
func (s *AuthGRPCServer) ChangePassword(ctx context.Context, req *authv1.ChangePasswordRequest) (*authv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.ChangePasswordCommand{
		UserID:          user.UserID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	if err := s.changePasswordHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.SuccessResponse{
		Success: true,
		Message: "Password changed successfully",
	}, nil
}

// VerifyEmail verifies the user's email address.
func (s *AuthGRPCServer) VerifyEmail(ctx context.Context, req *authv1.VerifyEmailRequest) (*authv1.SuccessResponse, error) {
	cmd := command.VerifyEmailCommand{
		Email: req.Email,
		Code:  req.Code,
	}

	if err := s.verifyEmailHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.SuccessResponse{
		Success: true,
		Message: "Email verified successfully",
	}, nil
}

// ResendVerification resends the verification email.
func (s *AuthGRPCServer) ResendVerification(ctx context.Context, req *authv1.ResendVerificationRequest) (*authv1.SuccessResponse, error) {
	cmd := command.ResendVerificationCommand{
		Email: req.Email,
	}

	if err := s.resendVerificationHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.SuccessResponse{
		Success: true,
		Message: "Verification email sent",
	}, nil
}

// ForgotPassword initiates the password reset flow.
func (s *AuthGRPCServer) ForgotPassword(ctx context.Context, req *authv1.ForgotPasswordRequest) (*authv1.SuccessResponse, error) {
	cmd := command.ForgotPasswordCommand{
		Email: req.Email,
	}

	if err := s.forgotPasswordHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.SuccessResponse{
		Success: true,
		Message: "Password reset email sent",
	}, nil
}

// ResetPassword completes the password reset flow.
func (s *AuthGRPCServer) ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*authv1.SuccessResponse, error) {
	cmd := command.ResetPasswordCommand{
		Email:       req.Email,
		Code:        req.Code,
		NewPassword: req.NewPassword,
	}

	if err := s.resetPasswordHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.SuccessResponse{
		Success: true,
		Message: "Password reset successfully",
	}, nil
}

// ExportUserData exports all user data (GDPR compliance).
func (s *AuthGRPCServer) ExportUserData(ctx context.Context, req *authv1.ExportUserDataRequest) (*authv1.ExportUserDataResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	result, err := s.exportDataHandler.Handle(ctx, query.ExportUserDataQuery{
		UserID: user.UserID,
	})
	if err != nil {
		return nil, toGRPCError(err)
	}

	data, err := json.Marshal(result)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to serialize user data")
	}

	return &authv1.ExportUserDataResponse{
		Success: true,
		Data:    data,
	}, nil
}

// DeleteAccount permanently deletes the user account.
func (s *AuthGRPCServer) DeleteAccount(ctx context.Context, req *authv1.DeleteAccountRequest) (*authv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.DeleteAccountCommand{
		UserID:          user.UserID,
		Password:        req.Password,
		ConfirmDeletion: true,
	}

	if err := s.deleteAccountHandler.Handle(ctx, cmd); err != nil {
		return nil, toGRPCError(err)
	}

	return &authv1.SuccessResponse{
		Success: true,
		Message: "Account deleted successfully",
	}, nil
}

// toGRPCError converts application errors to gRPC status errors.
func toGRPCError(err error) error {
	return grpcutil.ToGRPCError(err)
}
