package push

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Subscription represents a Web Push subscription
type Subscription struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Endpoint  string
	P256dh    string // Public key for encryption
	Auth      string // Authentication secret
	UserAgent string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewSubscription creates a new subscription
func NewSubscription(userID uuid.UUID, endpoint, p256dh, auth, userAgent string) *Subscription {
	now := time.Now()
	return &Subscription{
		ID:        uuid.New(),
		UserID:    userID,
		Endpoint:  endpoint,
		P256dh:    p256dh,
		Auth:      auth,
		UserAgent: userAgent,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Repository defines the interface for push subscription persistence
type Repository interface {
	// Save saves or updates a push subscription
	Save(ctx context.Context, subscription *Subscription) error

	// FindByUserID finds all subscriptions for a user
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*Subscription, error)

	// FindByEndpoint finds a subscription by endpoint
	FindByEndpoint(ctx context.Context, endpoint string) (*Subscription, error)

	// Delete deletes a subscription by endpoint
	Delete(ctx context.Context, endpoint string) error

	// DeleteByUserID deletes all subscriptions for a user
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error

	// DeleteExpired deletes subscriptions older than the given duration
	DeleteExpired(ctx context.Context, olderThan time.Duration) (int64, error)
}
