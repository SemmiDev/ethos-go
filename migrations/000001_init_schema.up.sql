-- ============================================================================
-- HABIT TRACKING APPLICATION - COMPLETE SCHEMA (CONSOLIDATED)
-- ============================================================================

-- ============================================================================
-- USERS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS "users" (
  "user_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "avatar" varchar(500),
  "is_active" BOOLEAN NOT NULL DEFAULT TRUE,
  "hashed_password" varchar(255),  -- Nullable for OAuth users
  "password_changed_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  -- User verification fields
  "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,
  "verify_token" VARCHAR(255),
  "verify_expires_at" TIMESTAMPTZ,
  "password_reset_token" VARCHAR(255),
  "password_reset_expires_at" TIMESTAMPTZ,
  -- User timezone
  "timezone" VARCHAR(50) DEFAULT 'Asia/Jakarta',
  -- OAuth fields
  "auth_provider" varchar(50) DEFAULT 'email',
  "auth_provider_id" varchar(255)
);

CREATE INDEX IF NOT EXISTS idx_users_verify_token ON "users"("verify_token");
CREATE INDEX IF NOT EXISTS idx_users_reset_token ON "users"("password_reset_token");

COMMENT ON COLUMN users.timezone IS 'User timezone in IANA format (e.g., Asia/Jakarta, America/New_York, UTC)';

-- ============================================================================
-- SESSIONS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS "sessions" (
  "session_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "refresh_token" varchar(500) NOT NULL,
  "user_agent" varchar(500) NOT NULL,
  "client_ip" varchar(50) NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT fk_sessions_user FOREIGN KEY ("user_id") REFERENCES "users"("user_id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON "sessions"("user_id");
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON "sessions"("expires_at") WHERE NOT is_blocked;

-- ============================================================================
-- HABITS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS "habits" (
  "habit_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "name" varchar(255) NOT NULL,
  "description" text,
  "frequency" varchar(20) NOT NULL DEFAULT 'daily',
  "target_count" integer DEFAULT 1,
  "is_active" boolean DEFAULT true,
  "reminder_time" VARCHAR(5),
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  -- Advanced Recurrence
  "recurrence_days" SMALLINT DEFAULT 127,
  "recurrence_interval" INT DEFAULT 1,
  CONSTRAINT fk_habits_user FOREIGN KEY ("user_id") REFERENCES "users"("user_id") ON DELETE CASCADE,
  CONSTRAINT valid_frequency CHECK (frequency IN ('daily', 'weekly', 'monthly'))
);

CREATE INDEX IF NOT EXISTS idx_habits_user_active ON "habits"("user_id", "is_active");

COMMENT ON TABLE habits IS 'Daftar kebiasaan yang ingin dijalankan oleh user';
COMMENT ON COLUMN habits.frequency IS 'Frekuensi habit: daily, weekly, monthly';
COMMENT ON COLUMN habits.recurrence_days IS 'Bitmask for days: Sun=1, Mon=2, Tue=4, Wed=8, Thu=16, Fri=32, Sat=64. 127=all';

-- ============================================================================
-- HABIT LOGS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS "habit_logs" (
  "log_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "habit_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "log_date" date NOT NULL,
  "count" integer DEFAULT 1,
  "note" text,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  CONSTRAINT fk_habit_logs_habit FOREIGN KEY ("habit_id") REFERENCES "habits"("habit_id") ON DELETE CASCADE,
  CONSTRAINT fk_habit_logs_user FOREIGN KEY ("user_id") REFERENCES "users"("user_id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_habit_logs_habit_date ON "habit_logs"("habit_id", "log_date" DESC);
CREATE INDEX IF NOT EXISTS idx_habit_logs_user_date ON "habit_logs"("user_id", "log_date" DESC);

-- ============================================================================
-- HABIT STATS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS habit_stats (
    habit_id UUID PRIMARY KEY REFERENCES habits(habit_id) ON DELETE CASCADE,
    current_streak INT NOT NULL DEFAULT 0,
    longest_streak INT NOT NULL DEFAULT 0,
    total_completions INT NOT NULL DEFAULT 0,
    last_completed_at DATE,
    consistency_score DECIMAL(5,2) DEFAULT 0.0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================================================
-- HABIT VACATIONS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS habit_vacations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    habit_id UUID NOT NULL REFERENCES habits(habit_id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE,
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT valid_vacation_dates CHECK (end_date IS NULL OR end_date >= start_date)
);

CREATE INDEX IF NOT EXISTS idx_habit_vacations_habit_id ON habit_vacations(habit_id);
CREATE INDEX IF NOT EXISTS idx_habit_vacations_active ON habit_vacations(habit_id, start_date, end_date);

-- ============================================================================
-- NOTIFICATIONS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS "notifications" (
  "notification_id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "type" varchar(50) NOT NULL,
  "title" varchar(255) NOT NULL,
  "message" text NOT NULL,
  "data" jsonb,
  "is_read" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "read_at" timestamptz,
  CONSTRAINT fk_notifications_user FOREIGN KEY ("user_id") REFERENCES "users"("user_id") ON DELETE CASCADE,
  CONSTRAINT valid_notification_type CHECK (type IN ('streak_milestone', 'habit_reminder', 'achievement', 'system', 'welcome'))
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON "notifications"("user_id", "created_at" DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON "notifications"("user_id", "is_read") WHERE NOT is_read;

-- ============================================================================
-- PUSH SUBSCRIPTIONS TABLE
-- ============================================================================
CREATE TABLE IF NOT EXISTS push_subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    endpoint TEXT NOT NULL UNIQUE,
    p256dh TEXT NOT NULL,
    auth TEXT NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_push_subscriptions_user_id ON push_subscriptions(user_id);
CREATE INDEX idx_push_subscriptions_created_at ON push_subscriptions(created_at);
