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
	"github.com/semmidev/ethos-go/internal/notifications/service/webpush"
)

func NewApplication(
	db *sqlx.DB,
	log logger.Logger,
	metricsClient decorator.MetricsClient,
	cfg *config.Config,
) app.Application {
	repo := adapters.NewNotificationPostgresRepository(db)
	pushRepo := adapters.NewPushSubscriptionRepository(db)

	pushService := webpush.NewService(
		pushRepo,
		cfg.VapidPublicKey,
		cfg.VapidPrivateKey,
		cfg.VapidSubject,
		log,
	)

	return app.Application{
		Commands: app.Commands{
			CreateNotification: command.NewCreateNotificationHandler(
				repo,
				pushService,
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
			SubscribePush: command.NewSubscribePushHandler(
				pushRepo,
				log,
				metricsClient,
			),
			UnsubscribePush: command.NewUnsubscribePushHandler(
				pushRepo,
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
