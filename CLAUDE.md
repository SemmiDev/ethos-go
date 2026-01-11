# CLAUDE.md - Ethos-Go Agent Guidelines

This document serves as the **primary source of truth** for AI agents working on the Ethos-Go project.

## 1. Project Overview

**Ethos-Go** is a professional Habit Tracking application built with Go backend and React frontend.

### Technology Stack

| Layer            | Technology                             |
| ---------------- | -------------------------------------- |
| **Backend**      | Go 1.25+, Chi Router                   |
| **Database**     | PostgreSQL 17                          |
| **Cache/Queue**  | Redis 8 + Asynq                        |
| **API**          | gRPC + gRPC-Gateway (Schema-First)     |
| **Frontend**     | React 19, Vite, Pure CSS               |
| **Architecture** | DDD, Clean Architecture, Single Binary |

## 2. Quick Start Commands

Run `make help` to see all available commands. Here are the most common:

| Category  | Command                        | Description                            |
| :-------- | :----------------------------- | :------------------------------------- |
| **Dev**   | `make dev`                     | Start full dev environment (Docker)    |
|           | `make stop`                    | Stop environment                       |
| **Code**  | `make generate-grpc`           | Generate gRPC & Gateway code           |
|           | `make fmt`                     | Format code (gofmt)                    |
|           | `make test`                    | Run all tests                          |
| **Buf**   | `make buf-lint`                | Lint Protobuf files                    |
|           | `make buf-generate`            | Generate Go code from Proto            |
| **DB**    | `make migrate-create name=foo` | Create new migration                   |
|           | `make migrate-up`              | Apply migrations                       |
| **Build** | `make build`                   | Build single binary (Backend+Frontend) |

## 3. Development Workflows

### Workflow A: Adding a New API Feature (gRPC)

1. **Define API**: Edit `api/proto/ethos/{module}/v1/{service}.proto`
2. **Lint**: Run `make buf-lint` to ensure style compliance.
3. **Generate**: Run `make generate-grpc` to regenerate Go stubs.
4. **Domain**: Create Entity & Repository Interface in `internal/{module}/domain/`
5. **Application**: Implement Command/Query Handler in `internal/{module}/app/`
6. **Infrastructure**: Implement Repository in `internal/{module}/adapters/`
7. **Ports**: Implement gRPC Server in `internal/{module}/ports/grpc_server.go`

### Workflow B: Database Migration

1. **Create**: `make migrate-create name=add_column_x`
2. **Edit SQL**:
   - Use `snake_case`
   - Primary keys: `uuid` (`gen_random_uuid()`)
   - Include `created_at` and `updated_at` (TIMESTAMPTZ)
3. **Apply**: `make migrate-up`

### Workflow C: Background Tasks

1. **Define Task**: `internal/{module}/adapters/task/` (Task Type & Payload)
2. **Implement Processor**: Logic to handle the task
3. **Register**: Add handler in `cmd/worker/main.go`

## 4. Directory Structure

```
internal/{module}/
├── domain/           # Pure Go. Business rules & interfaces.
│   ├── {entity}/     # Entity structs, value objects
│   └── repository.go # Persistence interface
├── app/              # Use Cases (CQRS)
│   ├── command/      # Write (State Change) - Returns ID/Error
│   └── query/        # Read (Data Retrieval) - Returns DTOs
├── adapters/         # Implementation Details
│   ├── {entity}_repository.go  # SQL/DB implementation
│   └── task/         # Asynq task definitions
└── ports/            # Entry Points
    └── grpc_server.go  # gRPC Handlers (Controller)
```

## 5. Coding Standards

### Error Handling

- Use `apperror` package for all domain errors.
- **NEVER** return raw database errors to the client.
- Map errors in the gRPC handler using `grpcutil.ToGRPCError`.

### Logging

- Use structural logging (`logger` package).
- Include `ctx` in all log calls for tracing.
- Format: `logger.Info(ctx, "event_name", logger.Field{...})`

### Context

- All blocking operations (DB, HTTP, Redis) **MUST** accept `context.Context`.

## 6. Testing Strategy

- **Unit Tests**: Domain logic and small components.
- **Integration Tests**: Repositories (with real DB).
- **E2E Tests**: Critical flows (API endpoints).

## 7. Agent Guidelines (Plan & Verify)

When working on complex tasks:

1. **Plan First**: Create a `plan.md` or `implementation_plan.md` for user approval if the task modifies architecture or critical paths.
2. **Atomic Changes**: Keep PRs/Commits small and focused.
3. **Verify**:
   - ALWAYS run `make test` after backend changes.
   - If modifying UI, ask user to verify visual changes or check build `npm run build`.
4. **Communication**: Be concise. Use Markdown. Don't over-explain obvious code.
