package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// Handler processes a specific type of event
type Handler interface {
	// Handle processes the event
	Handle(ctx context.Context, data []byte) error
	// EventType returns the event type this handler processes
	EventType() string
}

// Consumer subscribes to and processes domain events
type Consumer struct {
	nc       *nats.Conn
	js       jetstream.JetStream
	stream   jetstream.Stream
	handlers map[string]Handler
	logger   logger.Logger
	ctx      context.Context
	cancel   context.CancelFunc
}

// ConsumerConfig holds consumer configuration
type ConsumerConfig struct {
	NATSConfig
	ConsumerName  string
	QueueGroup    string
	MaxDeliver    int
	AckWait       time.Duration
	FilterSubject string
}

// NewConsumer creates a new event consumer
func NewConsumer(ctx context.Context, cfg ConsumerConfig, log logger.Logger) (*Consumer, error) {
	// Connect to NATS
	opts := []nats.Option{
		nats.Name("ethos-event-consumer-" + cfg.ConsumerName),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.ReconnectWait(cfg.ReconnectWait),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("create JetStream context: %w", err)
	}

	streamName := cfg.StreamName
	if streamName == "" {
		streamName = StreamName
	}

	stream, err := js.Stream(ctx, streamName)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("get stream: %w", err)
	}

	consumerCtx, cancel := context.WithCancel(ctx)

	return &Consumer{
		nc:       nc,
		js:       js,
		stream:   stream,
		handlers: make(map[string]Handler),
		logger:   log,
		ctx:      consumerCtx,
		cancel:   cancel,
	}, nil
}

// RegisterHandler registers a handler for a specific event type
func (c *Consumer) RegisterHandler(h Handler) {
	c.handlers[h.EventType()] = h
	c.logger.Info(c.ctx, "registered event handler",
		logger.Field{Key: "event_type", Value: h.EventType()},
	)
}

// Start begins consuming events
func (c *Consumer) Start(ctx context.Context, consumerName, queueGroup string) error {
	// Create or get durable consumer
	consumer, err := c.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:          consumerName,
		Durable:       consumerName,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxDeliver:    3,
		FilterSubject: SubjectPrefix + ".>",
	})
	if err != nil {
		return fmt.Errorf("create consumer: %w", err)
	}

	c.logger.Info(ctx, "starting event consumer",
		logger.Field{Key: "consumer", Value: consumerName},
	)

	// Consume messages
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		c.handleMessage(ctx, msg)
	})
	if err != nil {
		return fmt.Errorf("start consuming: %w", err)
	}

	// Wait for context cancellation
	go func() {
		<-c.ctx.Done()
		cc.Stop()
	}()

	return nil
}

// handleMessage processes a single message
func (c *Consumer) handleMessage(ctx context.Context, msg jetstream.Msg) {
	// Extract event type from subject
	// Subject format: ethos.{module}.{entity}.{action}
	subject := msg.Subject()
	eventType := subject[len(SubjectPrefix)+1:] // Remove "ethos." prefix

	handler, ok := c.handlers[eventType]
	if !ok {
		c.logger.Debug(ctx, "no handler for event type",
			logger.Field{Key: "event_type", Value: eventType},
			logger.Field{Key: "subject", Value: subject},
		)
		// Ack anyway to avoid redelivery
		msg.Ack()
		return
	}

	// Process the event
	if err := handler.Handle(ctx, msg.Data()); err != nil {
		c.logger.Error(ctx, err, "failed to handle event",
			logger.Field{Key: "event_type", Value: eventType},
		)
		// Nak for redelivery
		msg.Nak()
		return
	}

	// Acknowledge successful processing
	msg.Ack()
	c.logger.Debug(ctx, "event processed",
		logger.Field{Key: "event_type", Value: eventType},
	)
}

// Close stops the consumer and closes the connection
func (c *Consumer) Close() error {
	c.cancel()
	c.nc.Close()
	return nil
}

// ParseEvent is a helper to unmarshal event data
func ParseEvent[T any](data []byte) (*T, error) {
	var event T
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("unmarshal event: %w", err)
	}
	return &event, nil
}
