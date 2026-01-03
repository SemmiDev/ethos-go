package webpush

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/uuid"
	"github.com/semmidev/ethos-go/internal/common/logger"
	"github.com/semmidev/ethos-go/internal/notifications/domain/push"
)

// Service handles Web Push notifications
type Service struct {
	repo         push.Repository
	vapidPublic  string
	vapidPrivate string
	vapidSubject string // mailto: or https:// URL
	logger       logger.Logger
	httpClient   *http.Client
}

// NewService creates a new Web Push service
func NewService(
	repo push.Repository,
	vapidPublicKey string,
	vapidPrivateKey string,
	vapidSubject string,
	log logger.Logger,
) *Service {
	return &Service{
		repo:         repo,
		vapidPublic:  vapidPublicKey,
		vapidPrivate: vapidPrivateKey,
		vapidSubject: vapidSubject,
		logger:       log,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetVAPIDPublicKey returns the VAPID public key for client-side subscription
func (s *Service) GetVAPIDPublicKey() string {
	return s.vapidPublic
}

// Subscribe saves a new push subscription for a user
func (s *Service) Subscribe(ctx context.Context, userID uuid.UUID, endpoint, p256dh, auth, userAgent string) error {
	subscription := push.NewSubscription(userID, endpoint, p256dh, auth, userAgent)
	return s.repo.Save(ctx, subscription)
}

// Unsubscribe removes a push subscription
func (s *Service) Unsubscribe(ctx context.Context, endpoint string) error {
	return s.repo.Delete(ctx, endpoint)
}

// UnsubscribeUser removes all push subscriptions for a user
func (s *Service) UnsubscribeUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteByUserID(ctx, userID)
}

// SendNotification sends a push notification to a user
func (s *Service) SendNotification(ctx context.Context, userID uuid.UUID, payload push.Payload) error {
	subscriptions, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get subscriptions: %w", err)
	}

	if len(subscriptions) == 0 {
		s.logger.Info(ctx, "no push subscriptions for user", logger.Field{Key: "user_id", Value: userID.String()})
		return nil
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send to all subscriptions concurrently
	var wg sync.WaitGroup
	errChan := make(chan error, len(subscriptions))

	for _, sub := range subscriptions {
		wg.Add(1)
		go func(sub *push.Subscription) {
			defer wg.Done()

			err := s.sendToSubscription(ctx, sub, payloadBytes)
			if err != nil {
				errChan <- err
				// If the subscription is invalid (gone/expired), delete it
				if s.isSubscriptionGone(err) {
					s.logger.Info(ctx, "removing expired subscription",
						logger.Field{Key: "endpoint", Value: sub.Endpoint[:50] + "..."},
					)
					_ = s.repo.Delete(ctx, sub.Endpoint)
				}
			}
		}(sub)
	}

	wg.Wait()
	close(errChan)

	// Collect errors
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		s.logger.Error(ctx, errors[0], "some push notifications failed",
			logger.Field{Key: "total", Value: len(subscriptions)},
			logger.Field{Key: "failed", Value: len(errors)},
		)
	}

	successCount := len(subscriptions) - len(errors)
	if successCount > 0 {
		s.logger.Info(ctx, "push notifications sent",
			logger.Field{Key: "user_id", Value: userID.String()},
			logger.Field{Key: "count", Value: successCount},
		)
	}

	return nil
}

// SendToAll sends a push notification to multiple users
func (s *Service) SendToAll(ctx context.Context, userIDs []uuid.UUID, payload push.Payload) error {
	var wg sync.WaitGroup
	for _, userID := range userIDs {
		wg.Add(1)
		go func(uid uuid.UUID) {
			defer wg.Done()
			if err := s.SendNotification(ctx, uid, payload); err != nil {
				s.logger.Error(ctx, err, "failed to send push to user",
					logger.Field{Key: "user_id", Value: uid.String()},
				)
			}
		}(userID)
	}
	wg.Wait()
	return nil
}

// sendToSubscription sends a push notification to a single subscription
func (s *Service) sendToSubscription(ctx context.Context, sub *push.Subscription, payload []byte) error {
	subscription := &webpush.Subscription{
		Endpoint: sub.Endpoint,
		Keys: webpush.Keys{
			P256dh: sub.P256dh,
			Auth:   sub.Auth,
		},
	}

	options := &webpush.Options{
		Subscriber:      s.vapidSubject,
		VAPIDPublicKey:  s.vapidPublic,
		VAPIDPrivateKey: s.vapidPrivate,
		TTL:             86400, // 24 hours
		Urgency:         webpush.UrgencyNormal,
		HTTPClient:      s.httpClient,
	}

	resp, err := webpush.SendNotification(payload, subscription, options)
	if err != nil {
		return fmt.Errorf("failed to send push: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("push service returned status %d", resp.StatusCode)
	}

	return nil
}

// isSubscriptionGone checks if the error indicates the subscription is no longer valid
func (s *Service) isSubscriptionGone(err error) bool {
	// Check for HTTP 404 (Not Found) or 410 (Gone)
	// These indicate the subscription is no longer valid
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "410") || contains(errStr, "404") || contains(errStr, "gone")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
