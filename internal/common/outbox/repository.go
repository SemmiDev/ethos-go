package outbox

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/semmidev/ethos-go/internal/common/events"
)

// OutboxEntry represents an event stored in the outbox table
type OutboxEntry struct {
	ID            uuid.UUID       `db:"id"`
	EventType     string          `db:"event_type"`
	AggregateType string          `db:"aggregate_type"`
	AggregateID   string          `db:"aggregate_id"`
	Payload       json.RawMessage `db:"payload"`
	Metadata      json.RawMessage `db:"metadata"`
	CreatedAt     time.Time       `db:"created_at"`
	PublishedAt   *time.Time      `db:"published_at"`
	Published     bool            `db:"published"`
	RetryCount    int             `db:"retry_count"`
	LastError     *string         `db:"last_error"`
}

// Repository handles outbox persistence
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new outbox repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// Insert adds an event to the outbox (should be called within a transaction)
func (r *Repository) Insert(ctx context.Context, event events.Event, aggregateType string) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO outbox (id, event_type, aggregate_type, aggregate_id, payload)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = r.db.ExecContext(ctx, query,
		uuid.New(),
		event.EventType(),
		aggregateType,
		event.AggregateID(),
		payload,
	)
	return err
}

// GetUnpublished retrieves unpublished events for processing
func (r *Repository) GetUnpublished(ctx context.Context, limit int) ([]OutboxEntry, error) {
	query := `
		SELECT id, event_type, aggregate_type, aggregate_id, payload, metadata,
		       created_at, published_at, published, retry_count, last_error
		FROM outbox
		WHERE published = FALSE
		ORDER BY created_at ASC
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`
	var entries []OutboxEntry
	err := r.db.SelectContext(ctx, &entries, query, limit)
	return entries, err
}

// MarkPublished marks an entry as successfully published
func (r *Repository) MarkPublished(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE outbox
		SET published = TRUE, published_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// MarkFailed records a publish failure
func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID, errMsg string) error {
	query := `
		UPDATE outbox
		SET retry_count = retry_count + 1, last_error = $2
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id, errMsg)
	return err
}

// CleanupOld removes published entries older than the given duration
func (r *Repository) CleanupOld(ctx context.Context, olderThan time.Duration) (int64, error) {
	query := `
		DELETE FROM outbox
		WHERE published = TRUE AND published_at < $1
	`
	result, err := r.db.ExecContext(ctx, query, time.Now().Add(-olderThan))
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
