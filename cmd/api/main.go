package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/semmidev/ethos-go/config"
	authtask "github.com/semmidev/ethos-go/internal/auth/adapters/task"
	authports "github.com/semmidev/ethos-go/internal/auth/ports"
	authsvc "github.com/semmidev/ethos-go/internal/auth/service"
	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/grpcutil"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/metrics"
	"github.com/semmidev/ethos-go/internal/common/observability"
	"github.com/semmidev/ethos-go/internal/common/outbox"
	authv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/auth/v1"
	habitsv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/habits/v1"
	notificationsv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/notifications/v1"
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
		logger.Field{Key: "commit", Value: commit},
		logger.Field{Key: "build_time", Value: buildTime},
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
	authGRPCServer := authports.NewAuthGRPCServer(
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
	habitsGRPCServer := habitports.NewHabitsGRPCServer(habitsApp)

	// Initialize Notifications module
	notificationsApp := notificationsvc.NewApplication(db, appLogger, metricsClient, cfg)
	notificationsGRPCServer := notificationports.NewNotificationsGRPCServer(notificationsApp)

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authports.UnaryAuthInterceptor(authApp.AuthService),
		),
	)

	// Register gRPC services
	authv1.RegisterAuthServiceServer(grpcServer, authGRPCServer)
	habitsv1.RegisterHabitsServiceServer(grpcServer, habitsGRPCServer)
	notificationsv1.RegisterNotificationsServiceServer(grpcServer, notificationsGRPCServer)
	reflection.Register(grpcServer) // Enable gRPC reflection for debugging

	// Start gRPC server
	grpcPort := ":50051"
	grpcListener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port: %w", err)
	}

	go func() {
		appLogger.Info(ctx, "starting gRPC server", logger.Field{Key: "port", Value: grpcPort})
		if err := grpcServer.Serve(grpcListener); err != nil {
			appLogger.Error(ctx, err, "gRPC server error")
		}
	}()

	// Create gRPC-Gateway mux
	gwMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customHeaderMatcher),
		runtime.WithErrorHandler(grpcutil.CustomHTTPError),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)

	// Connect gRPC-Gateway to gRPC server
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcEndpoint := "localhost" + grpcPort

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, gwMux, grpcEndpoint, opts); err != nil {
		return fmt.Errorf("failed to register auth gateway: %w", err)
	}
	if err := habitsv1.RegisterHabitsServiceHandlerFromEndpoint(ctx, gwMux, grpcEndpoint, opts); err != nil {
		return fmt.Errorf("failed to register habits gateway: %w", err)
	}
	if err := notificationsv1.RegisterNotificationsServiceHandlerFromEndpoint(ctx, gwMux, grpcEndpoint, opts); err != nil {
		return fmt.Errorf("failed to register notifications gateway: %w", err)
	}

	// Create HTTP router with gRPC-Gateway
	router := NewRouter(RouterConfig{
		Config:         cfg,
		GatewayMux:     gwMux,
		OTELProvider:   otelProvider,
		Logger:         appLogger,
		AuthMiddleware: authApp.AuthMiddleware,
	})

	// Create HTTP server
	httpServer := NewServer(cfg, router, appLogger)

	// Start HTTP server in a goroutine
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

	// Stop gRPC server
	grpcServer.GracefulStop()
	appLogger.Info(ctx, "gRPC server stopped")

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	appLogger.Info(ctx, "server stopped gracefully")
	return nil
}

// customHeaderMatcher passes specific headers to gRPC metadata
func customHeaderMatcher(key string) (string, bool) {
	switch key {
	case "Authorization", "X-Request-Id", "X-Session-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
