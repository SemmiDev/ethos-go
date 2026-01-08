package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

const (
	// StreamName is the JetStream stream for domain events
	StreamName = "ETHOS_EVENTS"
	// SubjectPrefix is the prefix for all event subjects
	SubjectPrefix = "ethos"
)

// NATSPublisher publishes events to NATS JetStream
type NATSPublisher struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	logger logger.Logger
}

// NATSConfig holds NATS connection configuration
type NATSConfig struct {
	URL           string
	StreamName    string
	MaxReconnects int
	ReconnectWait time.Duration
}

// NewNATSPublisher creates a new NATS JetStream publisher
func NewNATSPublisher(ctx context.Context, cfg NATSConfig, log logger.Logger) (*NATSPublisher, error) {
	// Connect to NATS
	opts := []nats.Option{
		nats.Name("ethos-event-publisher"),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			if err != nil {
				log.Error(ctx, err, "NATS disconnected")
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Info(ctx, "NATS reconnected", logger.Field{Key: "url", Value: nc.ConnectedUrl()})
		}),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("create JetStream context: %w", err)
	}

	// Ensure stream exists
	streamName := cfg.StreamName
	if streamName == "" {
		streamName = StreamName
	}

	_, err = js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        streamName,
		Description: "Ethos domain events stream",
		Subjects:    []string{SubjectPrefix + ".>"},
		Storage:     jetstream.FileStorage,
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      7 * 24 * time.Hour, // Keep events for 7 days
		MaxMsgs:     -1,
		MaxBytes:    -1,
		Duplicates:  5 * time.Minute,
	})
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("create/update stream: %w", err)
	}

	log.Info(ctx, "NATS JetStream publisher initialized",
		logger.Field{Key: "url", Value: cfg.URL},
		logger.Field{Key: "stream", Value: streamName},
	)

	return &NATSPublisher{
		nc:     nc,
		js:     js,
		logger: log,
	}, nil
}

// Publish publishes a single event to NATS JetStream
func (p *NATSPublisher) Publish(ctx context.Context, event Event) error {
	subject := p.buildSubject(event.EventType())

	// Serialize event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	// Publish with deduplication ID
	_, err = p.js.Publish(ctx, subject, data,
		jetstream.WithMsgID(event.EventID()),
	)
	if err != nil {
		p.logger.Error(ctx, err, "failed to publish event",
			logger.Field{Key: "event_type", Value: event.EventType()},
			logger.Field{Key: "event_id", Value: event.EventID()},
		)
		return fmt.Errorf("publish event: %w", err)
	}

	p.logger.Debug(ctx, "event published",
		logger.Field{Key: "event_type", Value: event.EventType()},
		logger.Field{Key: "event_id", Value: event.EventID()},
		logger.Field{Key: "subject", Value: subject},
	)

	return nil
}

// PublishAll publishes multiple events
func (p *NATSPublisher) PublishAll(ctx context.Context, events []Event) error {
	for _, event := range events {
		if err := p.Publish(ctx, event); err != nil {
			return err
		}
	}
	return nil
}

// Close closes the NATS connection
func (p *NATSPublisher) Close() error {
	p.nc.Close()
	return nil
}

// buildSubject converts event type to NATS subject
// e.g., "user.registered" -> "ethos.user.registered"
func (p *NATSPublisher) buildSubject(eventType string) string {
	return SubjectPrefix + "." + eventType
}

// Ensure NATSPublisher implements Publisher
var _ Publisher = (*NATSPublisher)(nil)

// ErrPublishFailed is returned when event publishing fails
var ErrPublishFailed = errors.New("failed to publish event")
