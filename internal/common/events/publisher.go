package events

import "context"

// Publisher defines the interface for publishing domain events
type Publisher interface {
	// Publish publishes a single event
	Publish(ctx context.Context, event Event) error
	// PublishAll publishes multiple events atomically where possible
	PublishAll(ctx context.Context, events []Event) error
	// Close closes the publisher connection
	Close() error
}

// NoOpPublisher is a no-op implementation for testing
type NoOpPublisher struct{}

func NewNoOpPublisher() *NoOpPublisher {
	return &NoOpPublisher{}
}

func (p *NoOpPublisher) Publish(ctx context.Context, event Event) error {
	return nil
}

func (p *NoOpPublisher) PublishAll(ctx context.Context, events []Event) error {
	return nil
}

func (p *NoOpPublisher) Close() error {
	return nil
}
