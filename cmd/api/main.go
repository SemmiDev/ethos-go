package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/config"
	authtask "github.com/semmidev/ethos-go/internal/auth/adapters/task"
	authports "github.com/semmidev/ethos-go/internal/auth/ports"
	authsvc "github.com/semmidev/ethos-go/internal/auth/service"
	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/observability"
	habittask "github.com/semmidev/ethos-go/internal/habits/adapters/task"
	habitports "github.com/semmidev/ethos-go/internal/habits/ports"
	habitsvc "github.com/semmidev/ethos-go/internal/habits/service"
	notificationports "github.com/semmidev/ethos-go/internal/notifications/ports"
	notificationsvc "github.com/semmidev/ethos-go/internal/notifications/service"
)

// Build-time variables injected via ldflags
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

// Application holds all initialized components
type Application struct {
	DB           *sqlx.DB
	AsynqClient  *asynq.Client
	OTELProvider *observability.Provider

	// Module servers
	AuthServer          *authports.AuthOpenAPIServer
	HabitsServer        *habitports.OpenAPIServer
	NotificationsServer *notificationports.NotificationOpenAPIServer

	// Auth middleware
	AuthMiddleware func(http.Handler) http.Handler
}

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// Setup logger
	appLogger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	appLogger.Info(ctx, "starting app",
		logger.Field{Key: "env", Value: cfg.AppEnv},
		logger.Field{Key: "version", Value: version},
	)

	// Bootstrap application
	app, err := bootstrap(ctx, cfg, appLogger)
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}
	defer app.close(ctx)

	// Setup router
	router := NewRouter(RouterConfig{
		Config:              cfg,
		AuthServer:          app.AuthServer,
		HabitsServer:        app.HabitsServer,
		NotificationsServer: app.NotificationsServer,
		AuthMiddleware:      app.AuthMiddleware,
		OTELProvider:        app.OTELProvider,
	})

	// Create and start server
	server := NewServer(cfg, router, appLogger)

	go func() {
		if err := server.Start(ctx); err != nil && err != http.ErrServerClosed {
			appLogger.Error(ctx, err, "server failed")
			log.Fatalf("server startup failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		appLogger.Error(shutdownCtx, err, "server forced to shutdown")
	}

	appLogger.Info(ctx, "server stopped")
}

// bootstrap initializes all application dependencies and modules
func bootstrap(ctx context.Context, cfg *config.Config, appLogger logger.Logger) (*Application, error) {
	app := &Application{}

	// Initialize OpenTelemetry
	otelProvider, err := observability.New(ctx, observability.Config{
		ServiceName:    cfg.AppName,
		ServiceVersion: version,
		Environment:    cfg.AppEnv,
		OTLPEndpoint:   cfg.OTLPEndpoint,
		EnableTracing:  cfg.OTLPEnableTracing,
		EnableMetrics:  cfg.OTLPEnableMetrics,
		SampleRate:     cfg.OTLPSampleRate,
	})
	if err != nil {
		appLogger.Error(ctx, err, "failed to initialize OpenTelemetry")
		return nil, err
	}
	app.OTELProvider = otelProvider
	appLogger.Info(ctx, "OpenTelemetry initialized",
		logger.Field{Key: "tracing", Value: cfg.OTLPEnableTracing},
		logger.Field{Key: "metrics", Value: cfg.OTLPEnableMetrics},
	)

	// Initialize OpenTelemetry Metrics
	if _, err := observability.InitMetrics(ctx); err != nil {
		appLogger.Error(ctx, err, "failed to initialize metrics")
		return nil, err
	}

	// Initialize database
	db, err := database.NewSQLXConnection(cfg)
	if err != nil {
		appLogger.Error(ctx, err, "failed to initialize database")
		return nil, err
	}
	app.DB = db
	appLogger.Info(ctx, "database connection established")

	// Initialize Asynq client
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisDSN(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	app.AsynqClient = asynq.NewClient(redisOpt)
	appLogger.Info(ctx, "asynq client initialized")

	// Initialize metrics client
	metricsClient := &decorator.NoOpMetricsClient{}

	// Initialize task dispatcher for habits
	habitDispatcher := habittask.NewAsynqTaskDispatcher(app.AsynqClient, appLogger)

	// Initialize task dispatcher for auth
	taskDispatcher := authtask.NewAsynqTaskDispatcher(cfg, app.AsynqClient)

	// Initialize Auth module
	authApp := authsvc.NewApplication(ctx, cfg, app.DB, taskDispatcher, appLogger, metricsClient)
	app.AuthServer = authports.NewAuthOpenAPIServer(
		authApp.Commands.Register,
		authApp.Commands.Login,
		authApp.Commands.Logout,
		authApp.Commands.LogoutAll,
		authApp.Queries.ListSessions,
		authApp.Queries.GetProfile,
		authApp.Commands.UpdateProfile,
		authApp.Commands.ChangePassword,
		authApp.Commands.VerifyEmail,
		authApp.Commands.ResendVerification,
		authApp.Commands.ForgotPassword,
		authApp.Commands.ResetPassword,
		authApp.Commands.LoginGoogle,
		authApp.Queries.GetGoogleAuthURL,
		authApp.Commands.RevokeSessions,
		authApp.Commands.DeleteAccount,
		authApp.Queries.ExportUserData,
	)
	app.AuthMiddleware = authApp.AuthMiddleware

	// Initialize Habits module
	habitsApp := habitsvc.NewApplication(ctx, app.DB, habitDispatcher, appLogger, metricsClient)
	app.HabitsServer = habitports.NewOpenAPIServer(habitsApp)

	// Initialize Notifications module
	notificationsApp := notificationsvc.NewApplication(app.DB, appLogger, metricsClient, cfg)
	app.NotificationsServer = notificationports.NewNotificationOpenAPIServer(notificationsApp, cfg.VapidPublicKey)

	return app, nil
}

// close releases all application resources
func (app *Application) close(ctx context.Context) {
	if app.DB != nil {
		app.DB.Close()
	}
	if app.AsynqClient != nil {
		app.AsynqClient.Close()
	}
	if app.OTELProvider != nil {
		app.OTELProvider.Shutdown(ctx)
	}
}
