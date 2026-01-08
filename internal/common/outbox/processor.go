package outbox

import (
	"context"
	"time"

	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// Processor polls the outbox and publishes events
type Processor struct {
	repo      *Repository
	publisher events.Publisher
	logger    logger.Logger
	interval  time.Duration
	batchSize int
}

// NewProcessor creates a new outbox processor
func NewProcessor(
	repo *Repository,
	publisher events.Publisher,
	log logger.Logger,
	interval time.Duration,
	batchSize int,
) *Processor {
	if interval == 0 {
		interval = 5 * time.Second
	}
	if batchSize == 0 {
		batchSize = 100
	}
	return &Processor{
		repo:      repo,
		publisher: publisher,
		logger:    log,
		interval:  interval,
		batchSize: batchSize,
	}
}

// Start begins the outbox polling loop
func (p *Processor) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	p.logger.Info(ctx, "outbox processor started",
		logger.Field{Key: "interval", Value: p.interval.String()},
		logger.Field{Key: "batch_size", Value: p.batchSize},
	)

	for {
		select {
		case <-ctx.Done():
			p.logger.Info(ctx, "outbox processor stopped")
			return
		case <-ticker.C:
			p.process(ctx)
		}
	}
}

func (p *Processor) process(ctx context.Context) {
	entries, err := p.repo.GetUnpublished(ctx, p.batchSize)
	if err != nil {
		p.logger.Error(ctx, err, "failed to get unpublished outbox entries")
		return
	}

	if len(entries) == 0 {
		return
	}

	p.logger.Debug(ctx, "processing outbox entries",
		logger.Field{Key: "count", Value: len(entries)},
	)

	for _, entry := range entries {
		// Create a simple event wrapper for publishing
		evt := &outboxEvent{
			id:            entry.ID.String(),
			eventType:     entry.EventType,
			aggregateID:   entry.AggregateID,
			aggregateType: entry.AggregateType,
			createdAt:     entry.CreatedAt,
			payload:       entry.Payload,
		}

		if err := p.publisher.Publish(ctx, evt); err != nil {
			p.logger.Error(ctx, err, "failed to publish outbox event",
				logger.Field{Key: "event_id", Value: entry.ID.String()},
				logger.Field{Key: "event_type", Value: entry.EventType},
			)
			_ = p.repo.MarkFailed(ctx, entry.ID, err.Error())
			continue
		}

		if err := p.repo.MarkPublished(ctx, entry.ID); err != nil {
			p.logger.Error(ctx, err, "failed to mark event as published",
				logger.Field{Key: "event_id", Value: entry.ID.String()},
			)
		}
	}
}

// outboxEvent wraps an outbox entry for publishing
type outboxEvent struct {
	id            string
	eventType     string
	aggregateID   string
	aggregateType string
	createdAt     time.Time
	payload       []byte
}

func (e *outboxEvent) EventID() string       { return e.id }
func (e *outboxEvent) EventType() string     { return e.eventType }
func (e *outboxEvent) OccurredAt() time.Time { return e.createdAt }
func (e *outboxEvent) AggregateID() string   { return e.aggregateID }
func (e *outboxEvent) AggregateType() string { return e.aggregateType }

// MarshalJSON returns the stored payload
func (e *outboxEvent) MarshalJSON() ([]byte, error) {
	return e.payload, nil
}
