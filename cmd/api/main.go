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
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/semmidev/ethos-go/config"
	authtask "github.com/semmidev/ethos-go/internal/auth/adapters/task"
	authapp "github.com/semmidev/ethos-go/internal/auth/app"
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
	habitsapp "github.com/semmidev/ethos-go/internal/habits/app"
	habitports "github.com/semmidev/ethos-go/internal/habits/ports"
	habitsvc "github.com/semmidev/ethos-go/internal/habits/service"
	notificationsapp "github.com/semmidev/ethos-go/internal/notifications/app"
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

	// Initialize infrastructure
	otelProvider, db, asynqClient, err := initInfrastructure(ctx, cfg, appLogger)
	if err != nil {
		return err
	}
	defer otelProvider.Shutdown(ctx)
	defer db.Close()
	defer asynqClient.Close()

	// Initialize application modules
	authApp, habitsApp, notificationsApp := initModules(ctx, cfg, db, asynqClient, appLogger)

	// Create and start gRPC server
	grpcServer, grpcPort := createGRPCServer(authApp, habitsApp, notificationsApp)
	go runGRPCServer(ctx, grpcServer, grpcPort, appLogger)

	// Create gRPC-Gateway and HTTP server
	gwMux, err := createGatewayMux(ctx, grpcPort)
	if err != nil {
		return err
	}

	router := NewRouter(RouterConfig{
		Config:         cfg,
		GatewayMux:     gwMux,
		OTELProvider:   otelProvider,
		Logger:         appLogger,
		AuthMiddleware: authApp.AuthMiddleware,
	})

	httpServer := NewServer(cfg, router, appLogger)

	// Start HTTP server
	serverErrors := make(chan error, 1)
	go func() {
		if err := httpServer.Start(ctx); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case <-ctx.Done():
		appLogger.Info(ctx, "shutdown signal received")
	}

	// Graceful shutdown
	return gracefulShutdown(ctx, grpcServer, httpServer, appLogger)
}

// initInfrastructure initializes all infrastructure dependencies.
func initInfrastructure(
	ctx context.Context,
	cfg *config.Config,
	appLogger logger.Logger,
) (*observability.Provider, *sqlx.DB, *asynq.Client, error) {
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
		return nil, nil, nil, fmt.Errorf("failed to initialize OpenTelemetry: %w", err)
	}

	appLogger.Info(ctx, "OpenTelemetry initialized",
		logger.Field{Key: "tracing", Value: cfg.OTLPEnableTracing},
		logger.Field{Key: "metrics", Value: cfg.OTLPEnableMetrics},
	)

	if _, err := observability.InitMetrics(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// Initialize database
	db, err := database.NewSQLXConnection(cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	appLogger.Info(ctx, "database connection established")

	if err := database.RunMigrations(cfg.DSN(), migrations.FS, "."); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	appLogger.Info(ctx, "database migrations completed")

	// Initialize Asynq client
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisDSN(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	asynqClient := asynq.NewClient(redisOpt)
	appLogger.Info(ctx, "asynq client initialized")

	return otelProvider, db, asynqClient, nil
}

// initModules initializes all application modules.
func initModules(
	ctx context.Context,
	cfg *config.Config,
	db *sqlx.DB,
	asynqClient *asynq.Client,
	appLogger logger.Logger,
) (authapp.Application, habitsapp.Application, notificationsapp.Application) {
	metricsClient := metrics.NewPrometheusMetricsClient()

	// Initialize Outbox publisher
	outboxRepo := outbox.NewRepository(db)
	eventPublisher := outbox.NewPublisher(outboxRepo)

	// Initialize task dispatchers
	habitDispatcher := habittask.NewAsynqTaskDispatcher(asynqClient, appLogger)
	authTaskDispatcher := authtask.NewAsynqTaskDispatcher(cfg, asynqClient)

	// Initialize modules
	authApp := authsvc.NewApplication(ctx, cfg, db, authTaskDispatcher, eventPublisher, appLogger, metricsClient)
	habitsApp := habitsvc.NewApplication(ctx, db, habitDispatcher, eventPublisher, appLogger, metricsClient)
	notificationsApp := notificationsvc.NewApplication(db, appLogger, metricsClient, cfg)

	return authApp, habitsApp, notificationsApp
}

// createGRPCServer creates and configures the gRPC server.
func createGRPCServer(
	authApp authapp.Application,
	habitsApp habitsapp.Application,
	notificationsApp notificationsapp.Application,
) (*grpc.Server, string) {
	grpcPort := ":50051"

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

	habitsGRPCServer := habitports.NewHabitsGRPCServer(habitsApp)
	notificationsGRPCServer := notificationports.NewNotificationsGRPCServer(notificationsApp)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authports.UnaryAuthInterceptor(authApp.AuthService),
		),
	)

	authv1.RegisterAuthServiceServer(grpcServer, authGRPCServer)
	habitsv1.RegisterHabitsServiceServer(grpcServer, habitsGRPCServer)
	notificationsv1.RegisterNotificationsServiceServer(grpcServer, notificationsGRPCServer)
	reflection.Register(grpcServer)

	return grpcServer, grpcPort
}

// runGRPCServer starts the gRPC server.
func runGRPCServer(ctx context.Context, server *grpc.Server, port string, appLogger logger.Logger) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		appLogger.Error(ctx, err, "failed to listen on gRPC port")
		return
	}

	appLogger.Info(ctx, "starting gRPC server", logger.Field{Key: "port", Value: port})
	if err := server.Serve(listener); err != nil {
		appLogger.Error(ctx, err, "gRPC server error")
	}
}

// createGatewayMux creates the gRPC-Gateway mux.
func createGatewayMux(ctx context.Context, grpcPort string) (*runtime.ServeMux, error) {
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

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	grpcEndpoint := "localhost" + grpcPort

	if err := authv1.RegisterAuthServiceHandlerFromEndpoint(ctx, gwMux, grpcEndpoint, opts); err != nil {
		return nil, fmt.Errorf("failed to register auth gateway: %w", err)
	}
	if err := habitsv1.RegisterHabitsServiceHandlerFromEndpoint(ctx, gwMux, grpcEndpoint, opts); err != nil {
		return nil, fmt.Errorf("failed to register habits gateway: %w", err)
	}
	if err := notificationsv1.RegisterNotificationsServiceHandlerFromEndpoint(ctx, gwMux, grpcEndpoint, opts); err != nil {
		return nil, fmt.Errorf("failed to register notifications gateway: %w", err)
	}

	return gwMux, nil
}

// gracefulShutdown handles graceful shutdown of all servers.
func gracefulShutdown(ctx context.Context, grpcServer *grpc.Server, httpServer *Server, appLogger logger.Logger) error {
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	grpcServer.GracefulStop()
	appLogger.Info(ctx, "gRPC server stopped")

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	appLogger.Info(ctx, "server stopped gracefully")
	return nil
}

// customHeaderMatcher passes specific headers to gRPC metadata.
func customHeaderMatcher(key string) (string, bool) {
	switch key {
	case "Authorization", "X-Request-Id", "X-Session-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
