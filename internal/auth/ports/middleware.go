package ports

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/semmidev/ethos-go/internal/auth/domain/service"
	"github.com/semmidev/ethos-go/internal/auth/domain/user"
	authctx "github.com/semmidev/ethos-go/internal/auth/infrastructure/context"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/httputil"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// contextKey is a custom type for context keys to avoid collisions.
// Using a custom type instead of string prevents other packages from
// accidentally overwriting our context values.
type contextKey string

const (
	userIDKey    contextKey = "user_id"
	sessionIDKey contextKey = "session_id"
	emailKey     contextKey = "email"
)

// AuthMiddleware creates middleware that validates access tokens and populates
// the request context with user information. This is the bridge between the
// HTTP layer and our authentication domain.
//
// How it works:
// 1. Extract the Authorization header from the request
// 2. Verify the token is valid and not expired
// 3. Extract claims (user ID, session ID) from the token
// 4. Add these claims to the request context
// 5. Call the next handler with the enriched context
//
// Any handlers that run after this middleware can trust that the request
// is authenticated and can safely extract user info from the context.
func AuthMiddleware(tokenVerifier service.TokenVerifier, userReader user.UserReader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract the token from the Authorization header
			// Standard format is: "Bearer <token>"
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondUnauthorized(w, r, "missing authorization header")
				return
			}

			// Split "Bearer <token>" into ["Bearer", "<token>"]
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondUnauthorized(w, r, "invalid authorization header format")
				return
			}

			token := parts[1]

			// Verify the token and extract its claims
			// This checks the signature, expiration, and issuer
			claims, err := tokenVerifier.VerifyAccessToken(r.Context(), token)
			if err != nil {
				// Token is invalid, expired, or malformed
				log.Printf("Token verification failed: %v", err) // Simple log for debug
				httputil.Error(w, r, apperror.Unauthorized("invalid or expired token"))
				return
			}

			// Optionally, fetch full user details to ensure account is still active
			// This adds a database call but prevents deleted/blocked users from accessing APIs
			user, err := userReader.FindByID(r.Context(), claims.UserID)
			if err != nil {
				respondUnauthorized(w, r, "user not found")
				return
			}

			if !user.IsActive {
				respondUnauthorized(w, r, "user account is not active")
				return
			}

			// Add authentication info to the request context
			// Downstream handlers can now safely access this data
			ctx := r.Context()
			ctx = context.WithValue(ctx, userIDKey, claims.UserID.String())
			ctx = context.WithValue(ctx, sessionIDKey, claims.SessionID.String())
			ctx = context.WithValue(ctx, emailKey, user.Email)

			// Also set the common auth context for other modules (like habits) that rely on it
			ctx = authctx.ContextWithUser(ctx, authctx.User{
				UserID: claims.UserID.String(),
				Email:  user.Email,
			})

			// Enrich wide event with user context for Canonical Log Lines
			// This adds user info to the single comprehensive log per request
			logger.AddUserContext(ctx, claims.UserID.String(), user.Email)

			// Call the next handler with the enriched context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// OptionalAuthMiddleware is similar to AuthMiddleware but doesn't reject
// requests with missing or invalid tokens. Instead, it populates the context
// if a valid token is present, but allows the request through either way.
//
// This is useful for endpoints that behave differently for authenticated vs
// anonymous users (like a home page that shows a "Login" button for anonymous
// users and a "Dashboard" button for authenticated users).
func OptionalAuthMiddleware(tokenVerifier service.TokenVerifier, userReader user.UserReader) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// No token provided - that's fine, just continue
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				// Invalid format - continue without auth
				next.ServeHTTP(w, r)
				return
			}

			token := parts[1]

			claims, err := tokenVerifier.VerifyAccessToken(r.Context(), token)
			if err != nil {
				// Invalid token - continue without auth
				next.ServeHTTP(w, r)
				return
			}

			// Token is valid - enrich the context
			user, err := userReader.FindByID(r.Context(), claims.UserID)
			if err != nil || !user.IsActive {
				// User not found or inactive - continue without auth
				next.ServeHTTP(w, r)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, userIDKey, claims.UserID.String())
			ctx = context.WithValue(ctx, sessionIDKey, claims.SessionID.String())
			ctx = context.WithValue(ctx, emailKey, user.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRoles is middleware that checks if the authenticated user has one of
// the required roles. This implements role-based access control (RBAC).
//
// NOTE: This is a placeholder showing the pattern. In a real app,
// you'd fetch user roles from the database and check them here.
func RequireRoles(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// In a real implementation, you'd:
			// 1. Get user ID from context
			// 2. Fetch user roles from database
			// 3. Check if they have any of the required roles
			// 4. If not, return 403 Forbidden

			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}

// GetUserIDFromContext extracts the authenticated user ID from the request context.
// This should only be called in handlers that run after AuthMiddleware.
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey).(string)
	return userID, ok
}

// GetSessionIDFromContext extracts the session ID from the request context.
func GetSessionIDFromContext(ctx context.Context) (string, bool) {
	sessionID, ok := ctx.Value(sessionIDKey).(string)
	return sessionID, ok
}

// GetEmailFromContext extracts the user's email from the request context.
func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(emailKey).(string)
	return email, ok
}

// respondUnauthorized sends a 401 Unauthorized response with a message.
func respondUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	httputil.Error(w, r, apperror.Unauthorized(message))
}
