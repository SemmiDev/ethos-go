package ports

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	authctx "github.com/semmidev/ethos-go/internal/auth/infrastructure/context"
	"github.com/semmidev/ethos-go/internal/common/grpcutil"
	"github.com/semmidev/ethos-go/internal/common/model"
	commonv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/common/v1"
	notificationsv1 "github.com/semmidev/ethos-go/internal/generated/grpc/ethos/notifications/v1"
	"github.com/semmidev/ethos-go/internal/notifications/app"
	"github.com/semmidev/ethos-go/internal/notifications/app/command"
	"github.com/semmidev/ethos-go/internal/notifications/app/query"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

// NotificationsGRPCServer implements the gRPC NotificationsService interface.
type NotificationsGRPCServer struct {
	notificationsv1.UnimplementedNotificationsServiceServer
	app app.Application
}

// NewNotificationsGRPCServer creates a new NotificationsGRPCServer.
func NewNotificationsGRPCServer(application app.Application) *NotificationsGRPCServer {
	return &NotificationsGRPCServer{app: application}
}

// CreateNotification creates a new notification.
func (s *NotificationsGRPCServer) CreateNotification(ctx context.Context, req *notificationsv1.CreateNotificationRequest) (*notificationsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	var data map[string]interface{}
	if req.Data != nil {
		data = req.Data.AsMap()
	}

	// Convert string type to domain.NotificationType
	notifType := domain.NotificationType(req.Type)

	cmd := command.CreateNotification{
		UserID:  user.UserID,
		Type:    notifType,
		Title:   req.Title,
		Message: req.Message,
		Data:    data,
	}

	if err := s.app.Commands.CreateNotification.Handle(ctx, cmd); err != nil {
		return nil, toNotificationsGRPCError(err)
	}

	return &notificationsv1.SuccessResponse{
		Success: true,
		Message: "Notification created successfully",
	}, nil
}

// ListNotifications returns notifications for the authenticated user.
func (s *NotificationsGRPCServer) ListNotifications(ctx context.Context, req *notificationsv1.ListNotificationsRequest) (*notificationsv1.ListNotificationsResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	filter := model.NewFilter()
	if req.Page > 0 {
		filter.CurrentPage = int(req.Page)
	}
	if req.PerPage > 0 {
		filter.PerPage = int(req.PerPage)
	}

	// Set unread filter if requested
	if req.UnreadOnly {
		isActive := true
		filter.IsActive = &isActive // Using IsActive to represent unread-only filter
	}

	result, err := s.app.Queries.ListNotifications.Handle(ctx, query.ListNotifications{
		UserID: user.UserID,
		Filter: filter,
	})
	if err != nil {
		return nil, toNotificationsGRPCError(err)
	}

	notifications := make([]*notificationsv1.Notification, 0, len(result.Notifications))
	for _, n := range result.Notifications {
		notifications = append(notifications, toProtoNotification(n))
	}

	return &notificationsv1.ListNotificationsResponse{
		Success: true,
		Message: "Notifications retrieved successfully",
		Data:    notifications,
		Meta: &commonv1.Meta{
			Pagination: &commonv1.PaginationResponse{
				HasPreviousPage:        result.Pagination.HasPreviousPage,
				HasNextPage:            result.Pagination.HasNextPage,
				CurrentPage:            int32(result.Pagination.CurrentPage),
				PerPage:                int32(result.Pagination.PerPage),
				TotalData:              int32(result.Pagination.TotalData),
				TotalDataInCurrentPage: int32(result.Pagination.TotalDataInCurrentPage),
				LastPage:               int32(result.Pagination.LastPage),
				From:                   int32(result.Pagination.From),
				To:                     int32(result.Pagination.To),
			},
		},
	}, nil
}

// GetUnreadCount returns the count of unread notifications.
func (s *NotificationsGRPCServer) GetUnreadCount(ctx context.Context, req *notificationsv1.GetUnreadCountRequest) (*notificationsv1.UnreadCountResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	count, err := s.app.Queries.GetUnreadCount.Handle(ctx, query.GetUnreadCount{
		UserID: user.UserID,
	})
	if err != nil {
		return nil, toNotificationsGRPCError(err)
	}

	return &notificationsv1.UnreadCountResponse{
		Success: true,
		Message: "Unread count retrieved successfully",
		Data: &notificationsv1.UnreadCountData{
			Count: int32(count),
		},
	}, nil
}

// MarkAsRead marks a notification as read.
func (s *NotificationsGRPCServer) MarkAsRead(ctx context.Context, req *notificationsv1.MarkAsReadRequest) (*notificationsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.MarkAsRead{
		NotificationID: req.NotificationId,
		UserID:         user.UserID,
	}

	if err := s.app.Commands.MarkAsRead.Handle(ctx, cmd); err != nil {
		return nil, toNotificationsGRPCError(err)
	}

	return &notificationsv1.SuccessResponse{
		Success: true,
		Message: "Notification marked as read",
	}, nil
}

// MarkAllAsRead marks all notifications as read.
func (s *NotificationsGRPCServer) MarkAllAsRead(ctx context.Context, req *notificationsv1.MarkAllAsReadRequest) (*notificationsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.MarkAllRead{
		UserID: user.UserID,
	}

	if err := s.app.Commands.MarkAllRead.Handle(ctx, cmd); err != nil {
		return nil, toNotificationsGRPCError(err)
	}

	return &notificationsv1.SuccessResponse{
		Success: true,
		Message: "All notifications marked as read",
	}, nil
}

// DeleteNotification deletes a notification.
func (s *NotificationsGRPCServer) DeleteNotification(ctx context.Context, req *notificationsv1.DeleteNotificationRequest) (*notificationsv1.SuccessResponse, error) {
	user, err := authctx.UserFromCtx(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "unauthorized")
	}

	cmd := command.DeleteNotification{
		NotificationID: req.NotificationId,
		UserID:         user.UserID,
	}

	if err := s.app.Commands.DeleteNotification.Handle(ctx, cmd); err != nil {
		return nil, toNotificationsGRPCError(err)
	}

	return &notificationsv1.SuccessResponse{
		Success: true,
		Message: "Notification deleted successfully",
	}, nil
}

// toProtoNotification converts a domain.Notification to a protobuf Notification.
func toProtoNotification(n domain.Notification) *notificationsv1.Notification {
	notifType := notificationsv1.NotificationType_NOTIFICATION_TYPE_SYSTEM
	switch n.Type {
	case domain.TypeStreakMilestone:
		notifType = notificationsv1.NotificationType_NOTIFICATION_TYPE_STREAK_MILESTONE
	case domain.TypeHabitReminder:
		notifType = notificationsv1.NotificationType_NOTIFICATION_TYPE_HABIT_REMINDER
	case domain.TypeAchievement:
		notifType = notificationsv1.NotificationType_NOTIFICATION_TYPE_ACHIEVEMENT
	case domain.TypeWelcome:
		notifType = notificationsv1.NotificationType_NOTIFICATION_TYPE_WELCOME
	}

	notif := &notificationsv1.Notification{
		Id:        n.ID,
		Type:      notifType,
		Title:     n.Title,
		Message:   n.Message,
		IsRead:    n.IsRead,
		CreatedAt: timestamppb.New(n.CreatedAt),
	}

	// Convert JSON data to protobuf Struct
	if len(n.Data) > 0 {
		var dataMap map[string]interface{}
		if err := json.Unmarshal(n.Data, &dataMap); err == nil {
			if data, err := structpb.NewStruct(dataMap); err == nil {
				notif.Data = data
			}
		}
	}

	if n.ReadAt != nil {
		notif.ReadAt = timestamppb.New(*n.ReadAt)
	}

	return notif
}

// toNotificationsGRPCError converts application errors to gRPC status errors.
func toNotificationsGRPCError(err error) error {
	return grpcutil.ToGRPCError(err)
}
