package ports

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/semmidev/ethos-go/internal/auth/app"
	authctx "github.com/semmidev/ethos-go/internal/auth/infrastructure/context"
)

// publicMethods lists gRPC methods that don't require authentication
var publicMethods = map[string]bool{
	"/ethos.auth.v1.AuthService/Register":           true,
	"/ethos.auth.v1.AuthService/Login":              true,
	"/ethos.auth.v1.AuthService/GoogleLogin":        true,
	"/ethos.auth.v1.AuthService/GoogleCallback":     true,
	"/ethos.auth.v1.AuthService/VerifyEmail":        true,
	"/ethos.auth.v1.AuthService/ResendVerification": true,
	"/ethos.auth.v1.AuthService/ForgotPassword":     true,
	"/ethos.auth.v1.AuthService/ResetPassword":      true,
}

// UnaryAuthInterceptor creates a gRPC unary interceptor for authentication
func UnaryAuthInterceptor(authSvc app.AuthServiceInterface) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Skip authentication for public methods
		if publicMethods[info.FullMethod] {
			return handler(ctx, req)
		}

		// Extract token from metadata
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 {
			// Also check grpcgateway-authorization header (set by gRPC-Gateway)
			authHeader = md.Get("grpcgateway-authorization")
		}
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		if tokenString == authHeader[0] {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
		}

		// Validate token
		payload, err := authSvc.ValidateToken(ctx, tokenString)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		// Get user from claims
		user, err := authSvc.GetUserByID(ctx, payload.UserID.String())
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "user not found")
		}

		// Add session ID from payload
		user.SessionID = payload.SessionID.String()

		// Add user to context
		ctx = authctx.ContextWithUser(ctx, user)

		return handler(ctx, req)
	}
}
