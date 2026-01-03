package main

import (
	"context"
	"log"
	"os"

	"github.com/hibiken/asynq"
	"github.com/semmidev/ethos-go/config"
	authadapter "github.com/semmidev/ethos-go/internal/auth/adapters"
	authtask "github.com/semmidev/ethos-go/internal/auth/adapters/task"
	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/email"
	"github.com/semmidev/ethos-go/internal/common/logger"
	habittask "github.com/semmidev/ethos-go/internal/habits/adapters/task"
	habitsvc "github.com/semmidev/ethos-go/internal/habits/service"
	notiftask "github.com/semmidev/ethos-go/internal/notifications/adapters/task"
	notificationsvc "github.com/semmidev/ethos-go/internal/notifications/service"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	// 2. Setup Logger
	appLogger, err := logger.New(cfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	ctx := context.Background()
	appLogger.Info(ctx, "starting worker",
		logger.Field{Key: "env", Value: cfg.AppEnv},
	)

	// 3. Initialize Database Connection
	db, err := database.NewSQLXConnection(cfg)
	if err != nil {
		appLogger.Error(ctx, err, "failed to connect to database")
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		appLogger.Error(ctx, err, "failed to ping database")
		os.Exit(1)
	}
	appLogger.Info(ctx, "database connection established")

	// 4. Initialize Dependency Apps
	metricsClient := &decorator.NoOpMetricsClient{}

	// Auth Dependencies
	sessionRepo := authadapter.NewPostgresSessionRepository(db)

	// Initialize Asynq Client for dispatcher usage
	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisDSN(),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	// Initialize task dispatcher for habits
	habitDispatcher := habittask.NewAsynqTaskDispatcher(asynqClient, appLogger)

	habitsApp := habitsvc.NewApplication(ctx, db, habitDispatcher, appLogger, metricsClient)

	// Notifications App
	notificationsApp := notificationsvc.NewApplication(db, appLogger, metricsClient, cfg)

	// 5. Setup Asynq Server (The Worker)
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

	// 6. Register Task Processors
	mux := asynq.NewServeMux()

	// Register Session Cleanup Processor
	sessionCleanupProcessor := authtask.NewSessionCleanupProcessor(sessionRepo, appLogger)
	mux.Handle(authtask.TaskSessionCleanup, sessionCleanupProcessor)

	// Register Notification Task Processor
	notifProcessor := notiftask.NewTaskProcessor(notificationsApp, habitsApp, appLogger)

	// Task: Process Daily Reminders
	mux.HandleFunc(notiftask.TaskProcessReminders, notifProcessor.ProcessTask)

	// Task: Habit Created (Immediate Notification)
	mux.HandleFunc(habittask.TaskHabitCreated, notifProcessor.ProcessHabitCreatedTask)

	// Register Email Task Processor
	smtpClient, err := email.NewSMTPClient(cfg, appLogger)
	if err != nil {
		appLogger.Error(ctx, err, "failed to initialize smtp client")
		os.Exit(1)
	}

	authTaskProcessor := authtask.NewTaskProcessor(appLogger, smtpClient)
	mux.HandleFunc(authtask.TaskSendVerifyEmail, authTaskProcessor.ProcessTaskSendVerifyEmail)
	mux.HandleFunc(authtask.TaskSendForgotPasswordEmail, authTaskProcessor.ProcessTaskSendForgotPasswordEmail)

	// 7. Setup Scheduler
	scheduler := asynq.NewScheduler(
		redisOpt,
		&asynq.SchedulerOpts{
			Logger: NewAsynqLogger(appLogger),
		},
	)

	// Schedule cleanup every 15 minutes
	if _, err := scheduler.Register("@every 15m", authtask.NewSessionCleanupTask()); err != nil {
		appLogger.Error(ctx, err, "failed to register cleanup schedule")
		os.Exit(1)
	}

	// Schedule Notification Reminders (every minute to support custom reminder times)
	if _, err := scheduler.Register("* * * * *", notiftask.NewProcessRemindersTask()); err != nil {
		appLogger.Error(ctx, err, "failed to register notification schedule")
		os.Exit(1)
	}

	// 8. Run Everything
	appLogger.Info(ctx, "starting worker and scheduler")

	// Run Scheduler in a goroutine
	go func() {
		if err := scheduler.Run(); err != nil {
			appLogger.Error(ctx, err, "scheduler failed")
			os.Exit(1)
		}
	}()

	// Run Server (Blocking)
	if err := srv.Run(mux); err != nil {
		appLogger.Error(ctx, err, "worker server failed")
		os.Exit(1)
	}
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
