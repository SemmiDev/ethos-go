package adapters

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/internal/notifications/domain/push"
)

type pushSubscriptionModel struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	Endpoint  string    `db:"endpoint"`
	P256dh    string    `db:"p256dh"`
	Auth      string    `db:"auth"`
	UserAgent *string   `db:"user_agent"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// PushSubscriptionRepository implements push.Repository
type PushSubscriptionRepository struct {
	db *sqlx.DB
}

// NewPushSubscriptionRepository creates a new push subscription repository
func NewPushSubscriptionRepository(db *sqlx.DB) *PushSubscriptionRepository {
	return &PushSubscriptionRepository{db: db}
}

// Save saves or updates a push subscription (upsert)
func (r *PushSubscriptionRepository) Save(ctx context.Context, subscription *push.Subscription) error {
	query := `
		INSERT INTO push_subscriptions (id, user_id, endpoint, p256dh, auth, user_agent, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (endpoint) DO UPDATE SET
			user_id = EXCLUDED.user_id,
			p256dh = EXCLUDED.p256dh,
			auth = EXCLUDED.auth,
			user_agent = EXCLUDED.user_agent,
			updated_at = EXCLUDED.updated_at
	`

	var userAgent *string
	if subscription.UserAgent != "" {
		userAgent = &subscription.UserAgent
	}

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID,
		subscription.UserID,
		subscription.Endpoint,
		subscription.P256dh,
		subscription.Auth,
		userAgent,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)

	return err
}

// FindByUserID finds all subscriptions for a user
func (r *PushSubscriptionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*push.Subscription, error) {
	query := `SELECT id, user_id, endpoint, p256dh, auth, user_agent, created_at, updated_at
	          FROM push_subscriptions WHERE user_id = $1`

	var models []pushSubscriptionModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	subscriptions := make([]*push.Subscription, len(models))
	for i, m := range models {
		subscriptions[i] = unmarshalPushSubscription(m)
	}

	return subscriptions, nil
}

// FindByEndpoint finds a subscription by endpoint
func (r *PushSubscriptionRepository) FindByEndpoint(ctx context.Context, endpoint string) (*push.Subscription, error) {
	query := `SELECT id, user_id, endpoint, p256dh, auth, user_agent, created_at, updated_at
	          FROM push_subscriptions WHERE endpoint = $1`

	var model pushSubscriptionModel
	err := r.db.GetContext(ctx, &model, query, endpoint)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return unmarshalPushSubscription(model), nil
}

// Delete deletes a subscription by endpoint
func (r *PushSubscriptionRepository) Delete(ctx context.Context, endpoint string) error {
	query := `DELETE FROM push_subscriptions WHERE endpoint = $1`
	_, err := r.db.ExecContext(ctx, query, endpoint)
	return err
}

// DeleteByUserID deletes all subscriptions for a user
func (r *PushSubscriptionRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `DELETE FROM push_subscriptions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

// DeleteExpired deletes subscriptions older than the given duration
func (r *PushSubscriptionRepository) DeleteExpired(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `DELETE FROM push_subscriptions WHERE updated_at < $1`
	cutoff := time.Now().Add(-olderThan)

	result, err := r.db.ExecContext(ctx, query, cutoff)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func unmarshalPushSubscription(m pushSubscriptionModel) *push.Subscription {
	userAgent := ""
	if m.UserAgent != nil {
		userAgent = *m.UserAgent
	}

	return &push.Subscription{
		ID:        m.ID,
		UserID:    m.UserID,
		Endpoint:  m.Endpoint,
		P256dh:    m.P256dh,
		Auth:      m.Auth,
		UserAgent: userAgent,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
