package main

import (
	"net/http"

	genauth "github.com/semmidev/ethos-go/internal/generated/api/auth"
)

// corsMiddleware adds CORS headers to all responses
func corsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// scopeAwareAuthMiddleware returns a middleware that only applies authentication
// to routes that have bearerAuth security defined in the OpenAPI spec.
// The oapi-codegen generated code sets BearerAuthScopes in the context for protected routes.
func scopeAwareAuthMiddleware(authMiddleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the route requires authentication by looking for BearerAuthScopes in context
			// The generated OpenAPI wrapper sets this value for routes with bearerAuth security
			if scopes := r.Context().Value(genauth.BearerAuthScopes); scopes != nil {
				// This route requires authentication, apply the auth middleware
				authMiddleware(next).ServeHTTP(w, r)
				return
			}
			// No authentication required for this route
			next.ServeHTTP(w, r)
		})
	}
}
