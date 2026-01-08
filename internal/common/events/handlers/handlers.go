package handlers

import (
	"context"
	"encoding/json"

	"github.com/semmidev/ethos-go/internal/common/events"
	"github.com/semmidev/ethos-go/internal/common/logger"
)

// UserRegisteredHandler handles UserRegistered events
type UserRegisteredHandler struct {
	logger logger.Logger
}

func NewUserRegisteredHandler(log logger.Logger) *UserRegisteredHandler {
	return &UserRegisteredHandler{logger: log}
}

func (h *UserRegisteredHandler) EventType() string {
	return "auth.user.registered"
}

func (h *UserRegisteredHandler) Handle(ctx context.Context, data []byte) error {
	event, err := events.ParseEvent[UserRegisteredEvent](data)
	if err != nil {
		return err
	}

	h.logger.Info(ctx, "handling user registered event",
		logger.Field{Key: "user_id", Value: event.UserID},
		logger.Field{Key: "email", Value: event.Email},
	)

	// Example: Create default notification preferences, welcome email, etc.
	// This is where you'd add business logic that reacts to new user registration

	return nil
}

// UserRegisteredEvent represents the event data
type UserRegisteredEvent struct {
	EventID      string `json:"event_id"`
	EventType    string `json:"event_type"`
	OccurredAt   string `json:"occurred_at"`
	AggregateID  string `json:"aggregate_id"`
	UserID       string `json:"user_id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	AuthProvider string `json:"auth_provider"`
}

// HabitCreatedHandler handles HabitCreated events
type HabitCreatedHandler struct {
	logger logger.Logger
}

func NewHabitCreatedHandler(log logger.Logger) *HabitCreatedHandler {
	return &HabitCreatedHandler{logger: log}
}

func (h *HabitCreatedHandler) EventType() string {
	return "habits.habit.created"
}

func (h *HabitCreatedHandler) Handle(ctx context.Context, data []byte) error {
	var event HabitCreatedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	h.logger.Info(ctx, "handling habit created event",
		logger.Field{Key: "habit_id", Value: event.HabitID},
		logger.Field{Key: "user_id", Value: event.UserID},
		logger.Field{Key: "name", Value: event.Name},
	)

	// Example: Create default reminders, analytics tracking, etc.

	return nil
}

// HabitCreatedEvent represents the event data
type HabitCreatedEvent struct {
	EventID     string `json:"event_id"`
	EventType   string `json:"event_type"`
	OccurredAt  string `json:"occurred_at"`
	AggregateID string `json:"aggregate_id"`
	HabitID     string `json:"habit_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Frequency   string `json:"frequency"`
	TargetCount int    `json:"target_count"`
}

// HabitCompletedHandler handles HabitCompleted events
type HabitCompletedHandler struct {
	logger logger.Logger
}

func NewHabitCompletedHandler(log logger.Logger) *HabitCompletedHandler {
	return &HabitCompletedHandler{logger: log}
}

func (h *HabitCompletedHandler) EventType() string {
	return "habits.habit.completed"
}

func (h *HabitCompletedHandler) Handle(ctx context.Context, data []byte) error {
	var event HabitCompletedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	h.logger.Info(ctx, "handling habit completed event",
		logger.Field{Key: "habit_id", Value: event.HabitID},
		logger.Field{Key: "user_id", Value: event.UserID},
		logger.Field{Key: "count", Value: event.Count},
	)

	// Example: Check for streak milestones, send notifications, etc.

	return nil
}

// HabitCompletedEvent represents the event data
type HabitCompletedEvent struct {
	EventID     string `json:"event_id"`
	EventType   string `json:"event_type"`
	OccurredAt  string `json:"occurred_at"`
	AggregateID string `json:"aggregate_id"`
	HabitID     string `json:"habit_id"`
	UserID      string `json:"user_id"`
	LogID       string `json:"log_id"`
	LogDate     string `json:"log_date"`
	Count       int    `json:"count"`
	TotalToday  int    `json:"total_today"`
}
