package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/config"
	authtask "github.com/semmidev/ethos-go/internal/auth/adapters/task"
	authports "github.com/semmidev/ethos-go/internal/auth/ports"
	authsvc "github.com/semmidev/ethos-go/internal/auth/service"
	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/metrics"
	"github.com/semmidev/ethos-go/internal/common/observability"
	"github.com/semmidev/ethos-go/internal/common/outbox"
	habittask "github.com/semmidev/ethos-go/internal/habits/adapters/task"
	habitports "github.com/semmidev/ethos-go/internal/habits/ports"
	habitsvc "github.com/semmidev/ethos-go/internal/habits/service"
	notificationports "github.com/semmidev/ethos-go/internal/notifications/ports"
	notificationsvc "github.com/semmidev/ethos-go/internal/notifications/service"
	"github.com/semmidev/ethos-go/migrations"
)

// Build-time variables injected via ldflags
var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

// main is deliberately kept simple: it only calls run().
func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run is the real entry point for the application.
func run(ctx context.Context, _, _ io.Writer) error {
	// Set up signal handling - context is cancelled on SIGINT or SIGTERM
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup logger
	appLogger, err := logger.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	appLogger.Info(ctx, "starting app",
		logger.Field{Key: "env", Value: cfg.AppEnv},
		logger.Field{Key: "version", Value: version},
	)

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
		return fmt.Errorf("failed to initialize OpenTelemetry: %w", err)
	}
	defer otelProvider.Shutdown(ctx)

	appLogger.Info(ctx, "OpenTelemetry initialized",
		logger.Field{Key: "tracing", Value: cfg.OTLPEnableTracing},
		logger.Field{Key: "metrics", Value: cfg.OTLPEnableMetrics},
	)

	// Initialize OpenTelemetry Metrics
	if _, err := observability.InitMetrics(ctx); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Initialize database
	db, err := database.NewSQLXConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()
	appLogger.Info(ctx, "database connection established")

	// Run database migrations
	if err := database.RunMigrations(cfg.DSN(), migrations.FS, "."); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	appLogger.Info(ctx, "database migrations completed")

	// Initialize Asynq client
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisDSN(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()
	appLogger.Info(ctx, "asynq client initialized")

	// Initialize metrics client
	metricsClient := metrics.NewPrometheusMetricsClient()

	// Initialize Outbox (for reliable event publishing)
	outboxRepo := outbox.NewRepository(db)
	outboxPublisher := outbox.NewPublisher(outboxRepo)
	appLogger.Info(ctx, "Outbox publisher initialized")

	// Use OutboxPublisher for services
	eventPublisher := outboxPublisher

	// Initialize task dispatchers
	habitDispatcher := habittask.NewAsynqTaskDispatcher(asynqClient, appLogger)
	authTaskDispatcher := authtask.NewAsynqTaskDispatcher(cfg, asynqClient)

	// Initialize Auth module
	authApp := authsvc.NewApplication(ctx, cfg, db, authTaskDispatcher, eventPublisher, appLogger, metricsClient)
	authServer := authports.NewAuthOpenAPIServer(
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

	// Initialize Habits module
	habitsApp := habitsvc.NewApplication(ctx, db, habitDispatcher, eventPublisher, appLogger, metricsClient)
	habitsServer := habitports.NewOpenAPIServer(habitsApp)

	// Initialize Notifications module
	notificationsApp := notificationsvc.NewApplication(db, appLogger, metricsClient, cfg)
	notificationsServer := notificationports.NewNotificationOpenAPIServer(notificationsApp)

	// Setup router
	router := NewRouter(RouterConfig{
		Config:              cfg,
		AuthServer:          authServer,
		HabitsServer:        habitsServer,
		NotificationsServer: notificationsServer,
		AuthMiddleware:      authApp.AuthMiddleware,
		OTELProvider:        otelProvider,
		Logger:              appLogger,
	})

	// Create HTTP server
	httpServer := NewServer(cfg, router, appLogger)

	// Start server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		if err := httpServer.Start(ctx); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Wait for either context cancellation or server error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		appLogger.Info(ctx, "shutdown signal received")
	}

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	appLogger.Info(ctx, "server stopped gracefully")
	return nil
}
