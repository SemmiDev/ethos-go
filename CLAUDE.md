# CLAUDE.md - Ethos-Go Agent Guidelines

This document serves as the **primary source of truth** for AI agents working on the Ethos-Go project.

## 1. Project Overview

**Ethos-Go** is a professional Habit Tracking application built with Go backend and React frontend.

### Technology Stack

| Layer            | Technology                     |
| ---------------- | ------------------------------ |
| **Backend**      | Go 1.21+, Chi Router, sqlx/pgx |
| **Database**     | PostgreSQL                     |
| **Cache/Queue**  | Redis + Asynq                  |
| **API**          | OpenAPI 3.0 (Schema-First)     |
| **Frontend**     | React 19, Vite, Pure CSS       |
| **Architecture** | DDD, Clean Architecture, CQRS  |

## 2. Quick Start Commands

Run `make help` to see all available commands. Here are the most common:

| Task                      | Command                                   |
| ------------------------- | ----------------------------------------- |
| **Start dev environment** | `make dev`                                |
| **Stop environment**      | `make stop`                               |
| **View logs**             | `make logs`                               |
| **Generate API code**     | `make generate`                           |
| **Create migration**      | `make migrate-create name=add_foo_column` |
| **Run migrations**        | `make migrate-up`                         |
| **Rebuild backend**       | `make rebuild-app`                        |
| **Rebuild frontend**      | `make rebuild-frontend`                   |
| **Run tests**             | `make test`                               |
| **Format code**           | `make fmt`                                |
| **Build binary**          | `make build`                              |

## 3. Development Workflows

### Workflow A: Adding a New API Feature

**Step 1: Define the API Contract (OpenAPI)**

```bash
# Edit the OpenAPI spec
vim api/openapi/{module}.yml    # auth.yml, habits.yml, or notifications.yml

# Generate Go server code
make generate
```

The generated code goes to `internal/generated/api/{module}/`.

**Step 2: Implement Domain Layer** (`internal/{module}/domain/`)

```go
// Define entity in domain/{entity}/{entity}.go
type Habit struct {
    ID        string
    Name      string
    Frequency Frequency
    // ...
}

// Define repository interface in domain/repository.go
type HabitRepository interface {
    Create(ctx context.Context, habit *habit.Habit) error
    FindByID(ctx context.Context, id string) (*habit.Habit, error)
}
```

**Step 3: Implement Application Layer** (`internal/{module}/app/`)

```go
// Command: internal/{module}/app/command/create_habit.go
type CreateHabitCommand struct {
    UserID    string
    Name      string
    Frequency string
}

type CreateHabitHandler struct {
    repo   domain.HabitRepository
    logger logger.Logger
}

func (h *CreateHabitHandler) Handle(ctx context.Context, cmd CreateHabitCommand) (string, error) {
    // Business logic here
}
```

**Step 4: Implement Repository** (`internal/{module}/adapters/`)

```go
// internal/{module}/adapters/{entity}_repository.go
type PostgresHabitRepository struct {
    db *sqlx.DB
}

func (r *PostgresHabitRepository) Create(ctx context.Context, h *habit.Habit) error {
    query := `INSERT INTO habits (id, user_id, name, frequency) VALUES ($1, $2, $3, $4)`
    _, err := r.db.ExecContext(ctx, query, h.ID, h.UserID, h.Name, h.Frequency)
    return err
}
```

**Step 5: Implement HTTP Handler** (`internal/{module}/ports/`)

Implement the generated OpenAPI interface in `openapi_server.go`.

---

### Workflow B: Database Migration

**Step 1: Create Migration Files**

```bash
# Creates migrations/000XXX_add_reminder_time.up.sql and .down.sql
make migrate-create name=add_reminder_time
```

**Step 2: Write SQL**

```sql
-- migrations/000XXX_add_reminder_time.up.sql
ALTER TABLE habits ADD COLUMN reminder_time TIME;

-- migrations/000XXX_add_reminder_time.down.sql
ALTER TABLE habits DROP COLUMN reminder_time;
```

**Migration Standards:**

- Use `uuid` for primary keys (`gen_random_uuid()`)
- Always include `created_at TIMESTAMPTZ DEFAULT NOW()`
- Always include `updated_at TIMESTAMPTZ DEFAULT NOW()`
- Use `snake_case` for all identifiers
- Use `TIMESTAMPTZ` not `TIMESTAMP`

**Step 3: Apply Migration**

```bash
# In Docker environment
make migrate-up

# Or with explicit DATABASE_URL
DATABASE_URL="postgres://user:pass@localhost:5432/ethosgo?sslmode=disable" make migrate-up
```

**Rollback if needed:**

```bash
make migrate-down      # Rollback last migration
make migrate-force version=5  # Force to specific version (when dirty)
```

---

### Workflow C: Adding Background Tasks

**Step 1: Define Task Type** (`internal/{module}/adapters/task/`)

```go
const TaskSendReminder = "notification:send_reminder"

type SendReminderPayload struct {
    UserID  string `json:"user_id"`
    HabitID string `json:"habit_id"`
}

func NewSendReminderTask(userID, habitID string) (*asynq.Task, error) {
    payload, _ := json.Marshal(SendReminderPayload{UserID: userID, HabitID: habitID})
    return asynq.NewTask(TaskSendReminder, payload), nil
}
```

**Step 2: Implement Processor**

```go
func (p *TaskProcessor) ProcessSendReminder(ctx context.Context, t *asynq.Task) error {
    var payload SendReminderPayload
    if err := json.Unmarshal(t.Payload(), &payload); err != nil {
        return err
    }
    // Process the task
    return nil
}
```

**Step 3: Register in Worker** (`cmd/worker/main.go`)

```go
mux.HandleFunc(task.TaskSendReminder, processor.ProcessSendReminder)
```

## 4. Directory Structure

```
internal/{module}/
├── domain/           # Pure Go. Business rules & interfaces.
│   ├── {entity}/     # Entity structs, value objects
│   ├── gateway/      # Interfaces for external comms
│   └── repository.go # Data persistence interface
├── app/              # Application layer
│   ├── command/      # Write operations (change state)
│   └── query/        # Read operations (return data)
├── adapters/         # Infrastructure implementations
│   ├── {entity}_repository.go  # Database implementations
│   └── task/         # Background task definitions
└── ports/            # Entry points
    └── openapi_server.go  # HTTP handlers
```

## 5. Code Patterns

### Error Handling

```go
import "github.com/semmidev/ethos-go/internal/common/apperror"

// Create domain errors
return apperror.NewNotFound("habit", id)
return apperror.NewValidation("name is required")
return apperror.NewUnauthorized("invalid credentials")
```

### HTTP Responses

```go
import "github.com/semmidev/ethos-go/internal/common/httputil"

httputil.Success(w, r, data, "Habit created successfully")
httputil.Error(w, r, err)
httputil.ValidationError(w, r, validationErrors)
```

### Logging

```go
logger.Info(ctx, "habit created",
    logger.Field{Key: "habit_id", Value: id},
    logger.Field{Key: "user_id", Value: userID},
)
logger.Error(ctx, err, "failed to create habit")
```

### Context

All blocking operations MUST accept `context.Context`:

```go
func (r *Repo) FindByID(ctx context.Context, id string) (*Entity, error)
```

## 6. Testing

```bash
make test           # Run all tests
make test-short     # Run short tests only
make test-coverage  # Generate coverage report
```

## 7. Docker & Deployment

```bash
make dev            # Start full development stack
make compose-build  # Rebuild and start
make compose-logs   # View all logs
make logs-app       # View backend logs only
make rebuild-app    # Rebuild backend only

# Build production binary
make build          # Includes embedded frontend
make build-backend  # Backend only
```

## 8. Observability

Services available in dev environment:

- **Grafana**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **Jaeger**: http://localhost:16686

API Documentation:

- **Auth API**: http://localhost:8080/auth/doc
- **Habits API**: http://localhost:8080/habits/doc
- **Notifications API**: http://localhost:8080/notifications/doc
