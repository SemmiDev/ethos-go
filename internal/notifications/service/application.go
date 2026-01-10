package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/config"
	"github.com/semmidev/ethos-go/internal/common/decorator"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/adapters"
	"github.com/semmidev/ethos-go/internal/notifications/app"
	"github.com/semmidev/ethos-go/internal/notifications/app/command"
	"github.com/semmidev/ethos-go/internal/notifications/app/query"
)

func NewApplication(
	db *sqlx.DB,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
	_ *config.Config, // config parameter kept for API compatibility but no longer used for VAPID
) app.Application {
	repo := adapters.NewNotificationPostgresRepository(db)

	return app.Application{
		Commands: app.Commands{
			CreateNotification: command.NewCreateNotificationHandler(
				repo,
				log,
				metricsClient,
			),
			MarkAsRead: command.NewMarkAsReadHandler(
				repo,
				log,
				metricsClient,
			),
			MarkAllRead: command.NewMarkAllReadHandler(
				repo,
				log,
				metricsClient,
			),
			DeleteNotification: command.NewDeleteNotificationHandler(
				repo,
				log,
				metricsClient,
			),
		},
		Queries: app.Queries{
			ListNotifications: query.NewListNotificationsHandler(
				repo,
				log,
				metricsClient,
			),
			GetUnreadCount: query.NewGetUnreadCountHandler(
				repo,
				log,
				metricsClient,
			),
		},
	}
}
