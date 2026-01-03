package push

import (
	"context"

	"github.com/google/uuid"
)

// Service defines the interface for sending push notifications
type Service interface {
	SendNotification(ctx context.Context, userID uuid.UUID, payload Payload) error
	SendToAll(ctx context.Context, userIDs []uuid.UUID, payload Payload) error
}

// Payload represents the data sent in a push notification
type Payload struct {
	Title   string            `json:"title"`
	Body    string            `json:"body"`
	Icon    string            `json:"icon,omitempty"`
	Badge   string            `json:"badge,omitempty"`
	Tag     string            `json:"tag,omitempty"`
	Data    map[string]string `json:"data,omitempty"`
	Actions []Action          `json:"actions,omitempty"`
}

// Action represents an action button
type Action struct {
	Action string `json:"action"`
	Title  string `json:"title"`
	Icon   string `json:"icon,omitempty"`
}
