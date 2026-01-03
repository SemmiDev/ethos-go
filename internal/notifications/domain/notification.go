package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	TypeStreakMilestone NotificationType = "streak_milestone"
	TypeHabitReminder   NotificationType = "habit_reminder"
	TypeAchievement     NotificationType = "achievement"
	TypeSystem          NotificationType = "system"
	TypeWelcome         NotificationType = "welcome"
)

type Notification struct {
	ID        string           `db:"notification_id" json:"id"`
	UserID    string           `db:"user_id" json:"user_id"`
	Type      NotificationType `db:"type" json:"type"`
	Title     string           `db:"title" json:"title"`
	Message   string           `db:"message" json:"message"`
	Data      json.RawMessage  `db:"data" json:"data"`
	IsRead    bool             `db:"is_read" json:"is_read"`
	CreatedAt time.Time        `db:"created_at" json:"created_at"`
	ReadAt    *time.Time       `db:"read_at" json:"read_at"`
}

func NewNotification(userID string, notifType NotificationType, title, message string, data map[string]interface{}) (*Notification, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return &Notification{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Data:      jsonData,
		IsRead:    false,
		CreatedAt: time.Now(),
	}, nil
}

func (n *Notification) MarkAsRead() {
	now := time.Now()
	n.IsRead = true
	n.ReadAt = &now
}
