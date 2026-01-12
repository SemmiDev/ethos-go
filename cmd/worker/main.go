package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/config"
	authadapter "github.com/semmidev/ethos-go/internal/auth/adapters"
	authtask "github.com/semmidev/ethos-go/internal/auth/adapters/task"
	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/email"
	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/events/handlers"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/metrics"
	"github.com/semmidev/ethos-go/internal/common/outbox"
	habittask "github.com/semmidev/ethos-go/internal/habits/adapters/task"
	habitsvc "github.com/semmidev/ethos-go/internal/habits/service"
	notifadapter "github.com/semmidev/ethos-go/internal/notifications/adapters"
	notiftask "github.com/semmidev/ethos-go/internal/notifications/adapters/task"
	notificationsvc "github.com/semmidev/ethos-go/internal/notifications/service"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, _, _ io.Writer) error {
	// Set up signal handling
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Load Configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup Logger
	appLogger, err := logger.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	appLogger.Info(ctx, "starting worker",
		logger.Field{Key: "env", Value: cfg.AppEnv},
	)

	// Initialize Database Connection
	db, err := database.NewSQLXConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	appLogger.Info(ctx, "database connection established")

	// Initialize Dependencies
	metricsClient := metrics.NewPrometheusMetricsClient()
	sessionRepo := authadapter.NewSessionPostgresRepository(db)
	userRepo := authadapter.NewUserPostgresRepository(db)

	// Create UserProvider adapter - this allows other modules to access user data via interface
	userProvider := authadapter.NewUserProviderAdapter(userRepo)

	// Create notification repository for cross-module communication
	notifRepo := notifadapter.NewNotificationPostgresRepository(db)

	// Initialize NATS
	var eventPublisher events.Publisher
	var eventConsumer *events.Consumer

	if cfg.NATSUrl != "" {
		// NATS Publisher
		natsPublisher, err := events.NewNATSPublisher(ctx, events.NATSConfig{
			URL:           cfg.NATSUrl,
			StreamName:    cfg.NATSStreamName,
			MaxReconnects: cfg.NATSMaxReconnects,
			ReconnectWait: 2 * time.Second,
		}, appLogger)
		if err != nil {
			appLogger.Error(ctx, err, "failed to initialize NATS publisher")
			// We continue, but outbox processor won't be able to publish
			eventPublisher = events.NewNoOpPublisher()
		} else {
			eventPublisher = natsPublisher
			defer natsPublisher.Close()
			appLogger.Info(ctx, "NATS publisher initialized")
		}

		// NATS Consumer
		natsConsumer, err := events.NewConsumer(ctx, events.ConsumerConfig{
			NATSConfig: events.NATSConfig{
				URL:           cfg.NATSUrl,
				StreamName:    cfg.NATSStreamName,
				MaxReconnects: cfg.NATSMaxReconnects,
				ReconnectWait: 2 * time.Second,
			},
			ConsumerName: cfg.NATSConsumerName,
			QueueGroup:   cfg.NATSConsumerName + "-group", // Load balance among workers
		}, appLogger)
		if err != nil {
			appLogger.Error(ctx, err, "failed to initialize NATS consumer")
		} else {
			eventConsumer = natsConsumer
			defer eventConsumer.Close()
			appLogger.Info(ctx, "NATS consumer initialized")

			// Register Event Handlers with cross-module dependencies
			// UserRegisteredHandler: uses UserProvider (Auth) + NotificationRepository (Notifications)
			eventConsumer.RegisterHandler(handlers.NewUserRegisteredHandler(appLogger, userProvider, notifRepo))
			eventConsumer.RegisterHandler(handlers.NewHabitCreatedHandler(appLogger))
			eventConsumer.RegisterHandler(handlers.NewHabitCompletedHandler(appLogger))

			// Start Consumer
			if err := eventConsumer.Start(ctx, cfg.NATSConsumerName, cfg.NATSConsumerName+"-group"); err != nil {
				appLogger.Error(ctx, err, "failed to start NATS consumer")
			}
		}
	} else {
		eventPublisher = events.NewNoOpPublisher()
		appLogger.Warn(ctx, "NATS not configured, skipping event integration")
	}

	// Initialize Outbox Processor
	outboxRepo := outbox.NewRepository(db)
	outboxProcessor := outbox.NewProcessor(
		outboxRepo,
		eventPublisher,
		appLogger,
		1*time.Second, // Poll every second
		50,            // Batch size
	)
	go outboxProcessor.Start(ctx) // Start in background

	// Initialize Asynq Client
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisDSN(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	// Initialize task dispatcher for habits
	habitDispatcher := habittask.NewAsynqTaskDispatcher(asynqClient, appLogger)
	habitsApp := habitsvc.NewApplication(ctx, db, habitDispatcher, eventPublisher, appLogger, metricsClient)

	// Notifications App
	notificationsApp := notificationsvc.NewApplication(db, appLogger, metricsClient, cfg)

	// Setup Asynq Server (The Worker)
	srv := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"default": 1,
			},
			Logger: NewAsynqLogger(appLogger),
		},
	)

	// Register Task Processors
	mux := asynq.NewServeMux()

	// Session Cleanup Processor
	sessionCleanupProcessor := authtask.NewSessionCleanupProcessor(sessionRepo, appLogger)
	mux.Handle(authtask.TaskSessionCleanup, sessionCleanupProcessor)

	// Notification Task Processor
	notifProcessor := notiftask.NewTaskProcessor(notificationsApp, habitsApp, appLogger)
	mux.HandleFunc(notiftask.TaskProcessReminders, notifProcessor.ProcessTask)
	mux.HandleFunc(habittask.TaskHabitCreated, notifProcessor.ProcessHabitCreatedTask)

	// Email Task Processor
	smtpClient, err := email.NewSMTPClient(cfg, appLogger)
	if err != nil {
		return fmt.Errorf("failed to initialize smtp client: %w", err)
	}

	authTaskProcessor := authtask.NewTaskProcessor(appLogger, smtpClient)
	mux.HandleFunc(authtask.TaskSendVerifyEmail, authTaskProcessor.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(authtask.TaskSendForgotPasswordEmail, authTaskProcessor.ProcessTaskSendForgotPasswordEmail)

	// Setup Scheduler
	scheduler := asynq.NewScheduler(
		redisOpt,
		&asynq.SchedulerOpts{
			Logger: NewAsynqLogger(appLogger),
		},
	)

	// Register scheduled tasks
	if _, err := scheduler.Register("@every 15m", authtask.NewSessionCleanupTask()); err != nil {
		return fmt.Errorf("failed to register cleanup schedule: %w", err)
	}

	if _, err := scheduler.Register("* * * * *", notiftask.NewProcessRemindersTask()); err != nil {
		return fmt.Errorf("failed to register notification schedule: %w", err)
	}

	appLogger.Info(ctx, "starting worker and scheduler")

	// Run Scheduler in a goroutine
	schedulerErrors := make(chan error, 1)
	go func() {
		if err := scheduler.Run(); err != nil {
			schedulerErrors <- err
		}
	}()

	// Run Server in a goroutine
	serverErrors := make(chan error, 1)
	go func() {
		if err := srv.Run(mux); err != nil {
			serverErrors <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-schedulerErrors:
		return fmt.Errorf("scheduler failed: %w", err)
	case err := <-serverErrors:
		return fmt.Errorf("worker server failed: %w", err)
	case <-ctx.Done():
		appLogger.Info(ctx, "shutdown signal received")
	}

	// Graceful shutdown
	srv.Shutdown()
	scheduler.Shutdown()

	appLogger.Info(ctx, "worker stopped gracefully")
	return nil
}

// NewAsynqLogger adapts our structured logger to asynq logger interface
func NewAsynqLogger(l logger.Logger) asynq.Logger {
	return &asynqLoggerAdapter{l}
}

type asynqLoggerAdapter struct {
	logger logger.Logger
}

func (l *asynqLoggerAdapter) Debug(args ...interface{}) {
	l.logger.Debug(context.Background(), "asynq", logger.Field{Key: "msg", Value: args})
}

func (l *asynqLoggerAdapter) Info(args ...interface{}) {
	l.logger.Info(context.Background(), "asynq", logger.Field{Key: "msg", Value: args})
}

func (l *asynqLoggerAdapter) Warn(args ...interface{}) {
	l.logger.Warn(context.Background(), "asynq", logger.Field{Key: "msg", Value: args})
}

func (l *asynqLoggerAdapter) Error(args ...interface{}) {
	l.logger.Error(context.Background(), nil, "asynq", logger.Field{Key: "msg", Value: args})
}

func (l *asynqLoggerAdapter) Fatal(args ...interface{}) {
	l.logger.Error(context.Background(), nil, "asynq fatal", logger.Field{Key: "msg", Value: args})
	os.Exit(1)
}
