package events

import (
	"time"

	"github.com/semmidev/ethos-go/internal/common/random"
)

// Event is the base interface for all domain events
type Event interface {
	// EventID returns a unique identifier for this event instance
	EventID() string
	// EventType returns the type of event (e.g., "user.registered")
	EventType() string
	// OccurredAt returns when the event occurred
	OccurredAt() time.Time
	// AggregateID returns the ID of the aggregate that produced this event
	AggregateID() string
	// AggregateType returns the type of the aggregate (e.g., "user", "habit")
	AggregateType() string
}

// BaseEvent provides common implementation for all events
type BaseEvent struct {
	ID          string    `json:"event_id"`
	Type        string    `json:"event_type"`
	Occurred    time.Time `json:"occurred_at"`
	AggregateId string    `json:"aggregate_id"`
	AggType     string    `json:"aggregate_type"`
}

// NewBaseEvent creates a new base event with auto-generated ID and current timestamp
func NewBaseEvent(eventType, aggregateType, aggregateID string) BaseEvent {
	return BaseEvent{
		ID:          random.NewUUID().String(),
		Type:        eventType,
		Occurred:    time.Now().UTC(),
		AggregateId: aggregateID,
		AggType:     aggregateType,
	}
}

func (e BaseEvent) EventID() string       { return e.ID }
func (e BaseEvent) EventType() string     { return e.Type }
func (e BaseEvent) OccurredAt() time.Time { return e.Occurred }
func (e BaseEvent) AggregateID() string   { return e.AggregateId }
func (e BaseEvent) AggregateType() string { return e.AggType }

// EventMetadata contains optional metadata for events
type EventMetadata struct {
	CorrelationID string            `json:"correlation_id,omitempty"`
	CausationID   string            `json:"causation_id,omitempty"`
	UserID        string            `json:"user_id,omitempty"`
	Source        string            `json:"source,omitempty"`
	Version       int               `json:"version,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
}
