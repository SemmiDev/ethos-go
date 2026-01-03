package main

import (
	"net/http"
	"runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/semmidev/ethos-go/config"
	authports "github.com/semmidev/ethos-go/internal/auth/ports"
	"github.com/semmidev/ethos-go/internal/common/docs"
	"github.com/semmidev/ethos-go/internal/common/httputil"
	"github.com/semmidev/ethos-go/internal/common/observability"
	genauth "github.com/semmidev/ethos-go/internal/generated/api/auth"
	genhabits "github.com/semmidev/ethos-go/internal/generated/api/habits"
	gennotifications "github.com/semmidev/ethos-go/internal/generated/api/notifications"
	habitports "github.com/semmidev/ethos-go/internal/habits/ports"
	notificationports "github.com/semmidev/ethos-go/internal/notifications/ports"
	"github.com/semmidev/ethos-go/internal/web"
)

// RouterConfig contains all dependencies needed for router setup
type RouterConfig struct {
	Config              *config.Config
	AuthServer          *authports.AuthOpenAPIServer
	HabitsServer        *habitports.OpenAPIServer
	NotificationsServer *notificationports.NotificationOpenAPIServer
	AuthMiddleware      func(http.Handler) http.Handler
	OTELProvider        *observability.Provider
}

// NewRouter creates and configures the main chi router with all routes and middleware
func NewRouter(rc RouterConfig) chi.Router {
	r := chi.NewRouter()

	// Apply global middleware
	applyGlobalMiddleware(r, rc.Config)

	// Mount utility endpoints
	mountUtilityEndpoints(r, rc.Config, rc.OTELProvider)

	// Mount API documentation
	mountAPIDocs(r)

	// Mount API routes
	mountAPIRoutes(r, rc)

	// Mount SPA frontend (must be last)
	mountSPAHandler(r)

	return r
}

// applyGlobalMiddleware adds all global middleware to the router
func applyGlobalMiddleware(r chi.Router, cfg *config.Config) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsMiddleware())
	r.Use(observability.HTTPMiddleware(cfg.AppName))
}

// mountUtilityEndpoints adds health, version, metrics, and ping endpoints
func mountUtilityEndpoints(r chi.Router, cfg *config.Config, otelProvider *observability.Provider) {
	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httputil.Success(w, r, map[string]string{
			"status":  "healthy",
			"service": cfg.AppName,
		}, "Health check passed")
	})

	// Version info
	r.Get("/version", func(w http.ResponseWriter, r *http.Request) {
		httputil.Success(w, r, map[string]interface{}{
			"version":    version,
			"commit":     commit,
			"build_time": buildTime,
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
			"env":        cfg.AppEnv,
		}, "Version information")
	})

	// Prometheus metrics
	if otelProvider.PrometheusExporter != nil {
		r.Handle("/metrics", promhttp.Handler())
	}

	// Ping
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"pong"}`))
	})
}

// mountAPIDocs sets up API documentation endpoints
func mountAPIDocs(r chi.Router) {
	docsServer := docs.New(
		docs.Spec{Name: "Auth", Path: "/auth", GetSwagger: genauth.GetSwagger},
		docs.Spec{Name: "Habits", Path: "/habits", GetSwagger: genhabits.GetSwagger},
		docs.Spec{Name: "Notifications", Path: "/notifications", GetSwagger: gennotifications.GetSwagger},
	)
	docsServer.Mount(r)
}

// mountAPIRoutes mounts all API endpoints under /api prefix
func mountAPIRoutes(r chi.Router, rc RouterConfig) {
	r.Route("/api", func(api chi.Router) {
		// Auth routes with scope-aware authentication
		genauth.HandlerWithOptions(rc.AuthServer, genauth.ChiServerOptions{
			BaseURL:     "",
			BaseRouter:  api,
			Middlewares: []genauth.MiddlewareFunc{scopeAwareAuthMiddleware(rc.AuthMiddleware)},
		})

		// Protected routes (habits and notifications require auth)
		api.Group(func(protected chi.Router) {
			protected.Use(rc.AuthMiddleware)

			// Habits routes
			genhabits.HandlerFromMuxWithBaseURL(rc.HabitsServer, protected, "")

			// Notifications routes
			gennotifications.HandlerFromMuxWithBaseURL(rc.NotificationsServer, protected, "")
		})
	})
}

// mountSPAHandler serves the embedded frontend for SPA routing
func mountSPAHandler(r chi.Router) {
	if !web.IsFrontendBundled() {
		return
	}

	spaHandler, err := web.NewSPAHandler()
	if err != nil {
		return
	}

	// Static assets
	r.Handle("/assets/*", spaHandler)
	r.Handle("/manifest.json", spaHandler)
	r.Handle("/sw.js", spaHandler)
	r.Handle("/robots.txt", spaHandler)
	r.Handle("/favicon.ico", spaHandler)
	r.Handle("/icons/*", spaHandler)

	// Catch-all for SPA routing (must be last)
	r.NotFound(spaHandler.ServeHTTP)
}
