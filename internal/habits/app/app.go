package app

import (
	"github.com/semmidev/ethos-go/internal/habits/app/command"
	"github.com/semmidev/ethos-go/internal/habits/app/query"
)

// Application is the main application service facade for the habits module
type Application struct {
	Commands Commands
	Queries  Queries
}

// Commands groups all command handlers (write operations)
type Commands struct {
	CreateHabit     command.CreateHabitHandler
	UpdateHabit     command.UpdateHabitHandler
	DeleteHabit     command.DeleteHabitHandler
	ActivateHabit   command.ActivateHabitHandler
	DeactivateHabit command.DeactivateHabitHandler
	LogHabit        command.LogHabitHandler
	UpdateHabitLog  command.UpdateHabitLogHandler
	DeleteHabitLog  command.DeleteHabitLogHandler
}

// Queries groups all query handlers (read operations)
type Queries struct {
	GetHabit           query.GetHabitHandler
	ListHabits         query.ListHabitsHandler
	GetHabitLogs       query.GetHabitLogsHandler
	GetHabitStats      query.GetHabitStatsHandler
	GetDashboard       query.GetDashboardHandler
	GetWeeklyAnalytics query.GetWeeklyAnalyticsHandler
	GetHabitsDue       query.GetHabitsDueHandler
}
