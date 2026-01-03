package testutil

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer wraps a testcontainers PostgreSQL instance
type PostgresContainer struct {
	*postgres.PostgresContainer
	ConnectionString string
	DB               *sqlx.DB
}

// NewPostgresContainer creates and starts a new PostgreSQL container for testing
func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	container, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("ethos_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PostgresContainer{
		PostgresContainer: container,
		ConnectionString:  connStr,
		DB:                db,
	}, nil
}

// RunMigrations applies the database schema for testing
func (c *PostgresContainer) RunMigrations(ctx context.Context) error {
	schema := `
		-- Users
		CREATE TABLE IF NOT EXISTS "users" (
			"user_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
			"name" varchar(255) NOT NULL,
			"email" varchar(255) UNIQUE NOT NULL,
			"avatar" varchar(500),
			"is_active" BOOLEAN NOT NULL DEFAULT TRUE,
			"hashed_password" varchar(255),
			"password_changed_at" timestamptz,
			"created_at" timestamptz NOT NULL DEFAULT (now()),
			"updated_at" timestamptz NOT NULL DEFAULT (now()),
			"is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
			"verify_token" VARCHAR(255),
			"verify_expires_at" TIMESTAMPTZ,
			"password_reset_token" VARCHAR(255),
			"password_reset_expires_at" TIMESTAMPTZ,
			"timezone" VARCHAR(50) DEFAULT 'Asia/Jakarta',
			"auth_provider" varchar(50) DEFAULT 'email',
			"auth_provider_id" varchar(255)
		);

		-- Sessions
		CREATE TABLE IF NOT EXISTS "sessions" (
			"session_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
			"user_id" uuid NOT NULL REFERENCES "users"("user_id") ON DELETE CASCADE,
			"refresh_token" varchar(500) NOT NULL,
			"user_agent" varchar(500) NOT NULL,
			"client_ip" varchar(50) NOT NULL,
			"is_blocked" boolean NOT NULL DEFAULT false,
			"expires_at" timestamptz NOT NULL,
			"created_at" timestamptz NOT NULL DEFAULT (now()),
			"updated_at" timestamptz NOT NULL DEFAULT (now())
		);

		-- Habits
		CREATE TABLE IF NOT EXISTS "habits" (
			"habit_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
			"user_id" uuid NOT NULL REFERENCES "users"("user_id") ON DELETE CASCADE,
			"name" varchar(255) NOT NULL,
			"description" text,
			"frequency" varchar(20) NOT NULL DEFAULT 'daily',
			"target_count" integer DEFAULT 1,
			"is_active" boolean DEFAULT true,
			"reminder_time" VARCHAR(5),
			"created_at" timestamptz NOT NULL DEFAULT (now()),
			"updated_at" timestamptz NOT NULL DEFAULT (now()),
			"recurrence_days" SMALLINT DEFAULT 127,
			"recurrence_interval" INT DEFAULT 1,
			CONSTRAINT valid_frequency CHECK (frequency IN ('daily', 'weekly', 'monthly'))
		);

		-- Habit Logs
		CREATE TABLE IF NOT EXISTS "habit_logs" (
			"log_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
			"habit_id" uuid NOT NULL REFERENCES "habits"("habit_id") ON DELETE CASCADE,
			"user_id" uuid NOT NULL REFERENCES "users"("user_id") ON DELETE CASCADE,
			"log_date" date NOT NULL,
			"count" integer DEFAULT 1,
			"note" text,
			"created_at" timestamptz NOT NULL DEFAULT (now()),
			"updated_at" timestamptz NOT NULL DEFAULT (now())
		);

		-- Habit Stats
		CREATE TABLE IF NOT EXISTS habit_stats (
			habit_id UUID PRIMARY KEY REFERENCES habits(habit_id) ON DELETE CASCADE,
			current_streak INT NOT NULL DEFAULT 0,
			longest_streak INT NOT NULL DEFAULT 0,
			total_completions INT NOT NULL DEFAULT 0,
			last_completed_at DATE,
			consistency_score DECIMAL(5,2) DEFAULT 0.0,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		-- Habit Vacations
		CREATE TABLE IF NOT EXISTS habit_vacations (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			habit_id UUID NOT NULL REFERENCES habits(habit_id) ON DELETE CASCADE,
			start_date DATE NOT NULL,
			end_date DATE,
			reason TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			CONSTRAINT valid_vacation_dates CHECK (end_date IS NULL OR end_date >= start_date)
		);
	`

	_, err := c.DB.ExecContext(ctx, schema)
	return err
}

// Cleanup cleans up the container and database connection
func (c *PostgresContainer) Cleanup(ctx context.Context) error {
	if c.DB != nil {
		c.DB.Close()
	}
	return c.Terminate(ctx)
}

// TruncateTables clears all data from tables (useful between tests)
func (c *PostgresContainer) TruncateTables(ctx context.Context) error {
	tables := []string{"habit_vacations", "habit_stats", "habit_logs", "habits", "sessions", "users"}
	for _, table := range tables {
		_, err := c.DB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return err
		}
	}
	return nil
}
