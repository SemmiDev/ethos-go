package ports

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/semmidev/ethos-go/internal/common/apperror"
	"github.com/semmidev/ethos-go/internal/common/auth"
	"github.com/semmidev/ethos-go/internal/common/httputil"
	"github.com/semmidev/ethos-go/internal/common/model"
	notifications "github.com/semmidev/ethos-go/internal/generated/api/notifications"
	"github.com/semmidev/ethos-go/internal/notifications/app"
	"github.com/semmidev/ethos-go/internal/notifications/app/command"
	"github.com/semmidev/ethos-go/internal/notifications/app/query"
	"github.com/semmidev/ethos-go/internal/notifications/domain"
)

type NotificationOpenAPIServer struct {
	app            app.Application
	vapidPublicKey string
}

func NewNotificationOpenAPIServer(app app.Application, vapidPublicKey string) *NotificationOpenAPIServer {
	return &NotificationOpenAPIServer{
		app:            app,
		vapidPublicKey: vapidPublicKey,
	}
}

// Ensure implementation
var _ notifications.ServerInterface = (*NotificationOpenAPIServer)(nil)

// ListNotifications implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) ListNotifications(w http.ResponseWriter, r *http.Request, params notifications.ListNotificationsParams) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	filter := model.FilterFromRequest(r)
	if params.Page != nil {
		filter.CurrentPage = *params.Page
	}
	if params.PerPage != nil {
		filter.PerPage = *params.PerPage
	}
	if params.UnreadOnly != nil && *params.UnreadOnly {
		isActive := true
		filter.IsActive = &isActive // Map UnreadOnly to IsActive for now as per repository
	}

	result, err := s.app.Queries.ListNotifications.Handle(r.Context(), query.ListNotifications{
		UserID: user.UserID,
		Filter: filter,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	response := notifications.PaginatedNotificationListWrapper{
		Success: toBoolPtr(true),
		Message: toStringPtr("Notifications retrieved"),
		Data:    toNotificationsPtr(result.Notifications),
		Meta: &struct {
			Pagination *notifications.Pagination `json:"pagination,omitempty"`
		}{
			Pagination: toOpenAPIPagination(result.Pagination),
		},
	}

	render.JSON(w, r, response)
}

// GetUnreadCount implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	count, err := s.app.Queries.GetUnreadCount.Handle(r.Context(), query.GetUnreadCount{
		UserID: user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	response := notifications.UnreadCountResponse{
		Success: toBoolPtr(true),
		Message: toStringPtr("Unread count retrieved"),
		Data: &struct {
			Count *int `json:"count,omitempty"`
		}{
			Count: &count,
		},
	}

	render.JSON(w, r, response)
}

// MarkAsRead implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) MarkAsRead(w http.ResponseWriter, r *http.Request, notificationId openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	err = s.app.Commands.MarkAsRead.Handle(r.Context(), command.MarkAsRead{
		NotificationID: notificationId.String(),
		UserID:         user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Notification marked as read")
}

// MarkAllAsRead implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	err = s.app.Commands.MarkAllRead.Handle(r.Context(), command.MarkAllRead{
		UserID: user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "All notifications marked as read")
}

// CreateNotification implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) CreateNotification(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var body notifications.CreateNotificationJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid body"))
		return
	}

	// Default data if nil
	var data map[string]interface{}
	if body.Data != nil {
		data = *body.Data
	}

	err = s.app.Commands.CreateNotification.Handle(r.Context(), command.CreateNotification{
		UserID:  user.UserID,
		Type:    domain.NotificationType(body.Type),
		Title:   body.Title,
		Message: body.Message,
		Data:    data,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Created(w, r, nil, "Notification created")
}

// DeleteNotification implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) DeleteNotification(w http.ResponseWriter, r *http.Request, notificationId openapi_types.UUID) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	err = s.app.Commands.DeleteNotification.Handle(r.Context(), command.DeleteNotification{
		NotificationID: notificationId.String(),
		UserID:         user.UserID,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Notification deleted")
}

// GetVapidPublicKey implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) GetVapidPublicKey(w http.ResponseWriter, r *http.Request) {
	response := notifications.VapidPublicKeyResponse{
		Success: toBoolPtr(true),
		Message: toStringPtr("VAPID public key retrieved"),
		Data: &struct {
			VapidPublicKey *string `json:"vapid_public_key,omitempty"`
		}{
			VapidPublicKey: &s.vapidPublicKey,
		},
	}

	render.JSON(w, r, response)
}

// SubscribePush implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) SubscribePush(w http.ResponseWriter, r *http.Request) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var body notifications.SubscribePushJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	// Extract keys from the request
	p256dh := body.Keys.P256dh
	authKey := body.Keys.Auth

	err = s.app.Commands.SubscribePush.Handle(r.Context(), command.SubscribePush{
		UserID:    user.UserID,
		Endpoint:  body.Endpoint,
		P256dh:    p256dh,
		Auth:      authKey,
		UserAgent: r.UserAgent(),
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Push subscription saved")
}

// UnsubscribePush implements notifications.ServerInterface
func (s *NotificationOpenAPIServer) UnsubscribePush(w http.ResponseWriter, r *http.Request) {
	_, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httputil.Error(w, r, apperror.Unauthorized("unauthorized"))
		return
	}

	var body struct {
		Endpoint string `json:"endpoint"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.Error(w, r, apperror.ValidationFailed("invalid request body"))
		return
	}

	err = s.app.Commands.UnsubscribePush.Handle(r.Context(), command.UnsubscribePush{
		Endpoint: body.Endpoint,
	})
	if err != nil {
		httputil.Error(w, r, err)
		return
	}

	httputil.Success(w, r, nil, "Push subscription removed")
}

// Helper functions

func toBoolPtr(v bool) *bool {
	return &v
}

func toStringPtr(v string) *string {
	return &v
}

func toNotificationsPtr(notifs []domain.Notification) *[]notifications.Notification {
	res := make([]notifications.Notification, len(notifs))
	for i, n := range notifs {
		id, _ := uuid.Parse(n.ID) // Parsing UUID string to UUID type for OpenAPI
		created := n.CreatedAt

		var readAt *time.Time
		if n.ReadAt != nil {
			readAt = n.ReadAt
		}

		// Convert JSON data to map
		var dataMap map[string]interface{}
		if len(n.Data) > 0 {
			_ = json.Unmarshal(n.Data, &dataMap)
		}

		res[i] = notifications.Notification{
			Id:        &id,
			Type:      (*notifications.NotificationType)(&n.Type), // Using Type from domain
			Title:     &n.Title,
			Message:   &n.Message,
			IsRead:    &n.IsRead,
			CreatedAt: &created,
			ReadAt:    readAt,
			Data:      &dataMap,
		}
	}
	return &res
}

func toOpenAPIPagination(p *model.Paging) *notifications.Pagination {
	if p == nil {
		return nil
	}
	return &notifications.Pagination{
		HasPreviousPage:        &p.HasPreviousPage,
		HasNextPage:            &p.HasNextPage,
		CurrentPage:            &p.CurrentPage,
		PerPage:                &p.PerPage,
		TotalData:              &p.TotalData,
		TotalDataInCurrentPage: &p.TotalDataInCurrentPage,
		LastPage:               &p.LastPage,
		From:                   &p.From,
		To:                     &p.To,
	}
}
