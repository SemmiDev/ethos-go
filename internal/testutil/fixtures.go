package testutil

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// TestUser represents a test user fixture
type TestUser struct {
	UserID         uuid.UUID
	Email          string
	Name           string
	HashedPassword string
}

// CreateTestUser creates a user directly in the database for testing
func (c *PostgresContainer) CreateTestUser(ctx context.Context, user TestUser) error {
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}

	query := `
		INSERT INTO users (user_id, email, name, hashed_password, is_verified, is_active)
		VALUES ($1, $2, $3, $4, true, true)
	`
	_, err := c.DB.ExecContext(ctx, query, user.UserID, user.Email, user.Name, user.HashedPassword)
	return err
}

// TestHabit represents a test habit fixture
type TestHabit struct {
	HabitID     string
	UserID      string
	Name        string
	Description *string
	Frequency   string
	TargetCount int
	IsActive    bool
}

// CreateTestHabit creates a habit directly in the database for testing
func (c *PostgresContainer) CreateTestHabit(ctx context.Context, habit TestHabit) error {
	if habit.HabitID == "" {
		habit.HabitID = uuid.New().String()
	}
	if habit.Frequency == "" {
		habit.Frequency = "daily"
	}
	if habit.TargetCount == 0 {
		habit.TargetCount = 1
	}

	now := time.Now()
	query := `
		INSERT INTO habits (habit_id, user_id, name, description, frequency, target_count, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := c.DB.ExecContext(ctx, query,
		habit.HabitID, habit.UserID, habit.Name, habit.Description,
		habit.Frequency, habit.TargetCount, habit.IsActive, now, now,
	)
	return err
}

// DefaultTestUser returns a default test user
func DefaultTestUser() TestUser {
	return TestUser{
		UserID:         uuid.New(),
		Email:          "test@example.com",
		Name:           "Test User",
		HashedPassword: "$2a$10$abcdefghijklmnopqrstuv", // dummy bcrypt hash
	}
}
