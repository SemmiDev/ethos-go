package outbox

import (
	"context"

	"github.com/semmidev/ethos-go/internal/common/events"
)

// Publisher publishes events to the database outbox table
type Publisher struct {
	repo *Repository
}

// NewPublisher creates a new outbox publisher
func NewPublisher(repo *Repository) *Publisher {
	return &Publisher{repo: repo}
}

// Publish persists the event to the outbox table
func (p *Publisher) Publish(ctx context.Context, event events.Event) error {
	return p.repo.Insert(ctx, event, event.AggregateType())
}

// PublishAll persists multiple events to the outbox table
func (p *Publisher) PublishAll(ctx context.Context, evts []events.Event) error {
	for _, event := range evts {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Close is a no-op for the outbox publisher
func (p *Publisher) Close() error {
	return nil
}

// Ensure Publisher implements events.Publisher
var _ events.Publisher = (*Publisher)(nil)
