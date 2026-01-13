package service

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/database"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/common/validator"
	"github.com/semmidev/ethos-go/internal/habits/adapters"
	"github.com/semmidev/ethos-go/internal/habits/app"
	"github.com/semmidev/ethos-go/internal/habits/app/command"
	"github.com/semmidev/ethos-go/internal/habits/app/query"
	domaintask "github.com/semmidev/ethos-go/internal/habits/domain/task"
)

// NewApplication creates and wires all dependencies for the habits module
func NewApplication(
	ctx context.Context,
	db database.DBTX,
	dispatcher domaintask.TaskDispatcher,
	eventPublisher events.Publisher, // Added eventPublisher
	log logger.Logger,
	metricsClient decorator.MetricsClient,
) app.Application {
	// Create repository instances
	habitRepo := adapters.NewHabitPostgresRepository(db)
	habitLogRepo := adapters.NewHabitLogPostgresRepository(db)
	statsRepo := adapters.NewStatsRepository(db)
	validate := validator.New("en")

	// Create command handlers with decorators
	return app.Application{
		Commands: app.Commands{
			CreateHabit: command.NewCreateHabitHandler(
				habitRepo,
				validate,
				dispatcher,
				eventPublisher,
				log,
				metricsClient,
			),
			UpdateHabit: command.NewUpdateHabitHandler(
				habitRepo,
				validate,
				log,
				metricsClient,
			),
			DeleteHabit: command.NewDeleteHabitHandler(
				habitRepo,
				validate,
				log,
				metricsClient,
			),
			ActivateHabit: command.NewActivateHabitHandler(
				habitRepo,
				validate,
				eventPublisher,
				log,
				metricsClient,
			),
			DeactivateHabit: command.NewDeactivateHabitHandler(
				habitRepo,
				validate,
				eventPublisher,
				log,
				metricsClient,
			),
			LogHabit: command.NewLogHabitHandler(
				habitRepo,
				habitLogRepo,
				validate,
				eventPublisher,
				log,
				metricsClient,
			),
			UpdateHabitLog: command.NewUpdateHabitLogHandler(
				habitLogRepo,
				validate,
				log,
				metricsClient,
			),
			DeleteHabitLog: command.NewDeleteHabitLogHandler(
				habitLogRepo,
				validate,
				log,
				metricsClient,
			),
		},
		Queries: app.Queries{
			GetHabit: query.NewGetHabitHandler(
				habitRepo,
				log,
				metricsClient,
			),
			ListHabits: query.NewListHabitsHandler(
				habitRepo,
				log,
				metricsClient,
			),
			GetHabitLogs: query.NewGetHabitLogsHandler(
				habitLogRepo,
				log,
				metricsClient,
			),
			GetHabitStats: query.NewGetHabitStatsHandler(
				statsRepo,
				log,
				metricsClient,
			),
			GetDashboard: query.NewGetDashboardHandler(
				statsRepo,
				log,
				metricsClient,
			),
			GetWeeklyAnalytics: query.NewGetWeeklyAnalyticsHandler(
				statsRepo,
				log,
				metricsClient,
			),
			GetHabitsDue: query.NewGetHabitsDueHandler(
				statsRepo,
				log,
				metricsClient,
			),
		},
	}
}
