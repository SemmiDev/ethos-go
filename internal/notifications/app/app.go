package app

import (
	"github.com/semmidev/ethos-go/internal/notifications/app/command"
	"github.com/semmidev/ethos-go/internal/notifications/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateNotification command.CreateNotificationHandler
	MarkAsRead         command.MarkAsReadHandler
	MarkAllRead        command.MarkAllReadHandler
	DeleteNotification command.DeleteNotificationHandler
}

type Queries struct {
	ListNotifications query.ListNotificationsHandler
	GetUnreadCount    query.GetUnreadCountHandler
}
