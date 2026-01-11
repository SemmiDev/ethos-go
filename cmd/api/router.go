package main

import (
	"net/http"
	goruntime "runtime"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/semmidev/ethos-go/config"
	"github.com/semmidev/ethos-go/internal/common/httputil"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/observability"
	"github.com/semmidev/ethos-go/internal/web"
)

// RouterConfig contains all dependencies needed for router setup
type RouterConfig struct {
	Config         *config.Config
	GatewayMux     *runtime.ServeMux
	OTELProvider   *observability.Provider
	Logger         logger.Logger
	AuthMiddleware func(http.Handler) http.Handler
}

// NewRouter creates and configures the main chi router with all routes and middleware
func NewRouter(rc RouterConfig) chi.Router {
	r := chi.NewRouter()

	// Apply global middleware
	applyGlobalMiddleware(r, rc)

	// Mount utility endpoints
	mountUtilityEndpoints(r, rc.Config, rc.OTELProvider)

	// Mount gRPC-Gateway API routes
	mountGatewayRoutes(r, rc)

	// Mount SPA frontend (must be last)
	mountSPAHandler(r)

	return r
}

// applyGlobalMiddleware adds all global middleware to the router
func applyGlobalMiddleware(r chi.Router, rc RouterConfig) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsMiddleware())
	r.Use(observability.HTTPMiddleware(rc.Config.AppName))

	// Event middleware (Canonical Log Lines)
	if rc.Logger != nil {
		sampler := logger.NewSampler(logger.SamplerConfig{
			Enabled:        true,
			BaseRate:       rc.Config.EventSampleRate,
			P99ThresholdMs: rc.Config.EventP99ThresholdMs,
		})
		r.Use(logger.EventMiddleware(logger.EventMiddlewareConfig{
			ServiceName: rc.Config.AppName,
			Version:     version,
			Environment: rc.Config.AppEnv,
			Logger:      rc.Logger,
			Sampler:     sampler,
		}))
	} else {
		r.Use(middleware.Logger)
	}
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
			"go_version": goruntime.Version(),
			"os":         goruntime.GOOS,
			"arch":       goruntime.GOARCH,
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

// mountGatewayRoutes mounts the gRPC-Gateway handler for API routes
func mountGatewayRoutes(r chi.Router, rc RouterConfig) {
	// Mount gRPC-Gateway under /v1 (the paths defined in proto files)
	r.Mount("/v1", rc.GatewayMux)

	// Also mount under /api for backward compatibility during transition
	// This allows existing clients to continue using /api/* paths
	r.Mount("/api", http.StripPrefix("/api", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// Rewrite /api/auth/* to /v1/auth/*, /api/habits/* to /v1/habits/*, etc.
		req.URL.Path = "/v1" + req.URL.Path
		rc.GatewayMux.ServeHTTP(w, req)
	})))
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
