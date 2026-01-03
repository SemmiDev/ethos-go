# CLAUDE.md - Ethos-Go Agent Guidelines

This document serves as the **primary source of truth** for AI agents working on the Ethos-Go project. It details the architectural decisions, coding standards, and workflows required to maintain the project's integrity.

## 1. Project Overview

**Ethos-Go** is a professional Habit Tracking application built with a focus on scalability, maintainability, and a premium user experience.

### Technology Stack

-   **Backend Core**: Go 1.21+, standard library focus.
-   **Web Framework**: Chi Router (lightweight, idiomatic).
-   **Database**: PostgreSQL (Driver: `pgx`, Builder: `sqlx`).
-   **Async Processing**: Redis with `Asynq` for robust background tasks.
-   **API Definition**: OpenAPI 3.0 (Schema-First approach).
-   **Frontend**: React 19, Vite, TailwindCSS (Utility-first), Custom "Banking" Design System.
-   **Architecture**: Domain-Driven Design (DDD), Clean Architecture, CQRS (Command Query Responsibility Segregation).

## 2. Architecture & Directory Structure

The project follows a strict **Dependency Rule**: Dependencies point **INWARDS**. The inner layers know nothing about the outer layers.

```text
internal/{module}/
├── domain/                  # [INNERMOST] Pure Go. Business Rules.
│   ├── {entity}/            # E.g., habit/habit.go (Structs, Value Objects)
│   ├── gateway/             # Interfaces for external comms (Email, Queue, etc.)
│   └── repository.go        # Interface for data persistence.
├── app/                     # [APPLICATION] Orchestration.
│   ├── command/             # Write side (Changes state).
│   ├── query/               # Read side (Returns views).
│   └── dto/                 # Data Transfer Objects (optional).
├── ports/                   # [INTERFACE] Entry points.
│   ├── http/                # REST API Handlers.
│   └── openapi_types.gen.go # Generated OpenAPI structs.
└── infrastructure/          # [OUTERMOST] Tools & Frameworks.
    ├── persistence/         # Database implementations (SQL).
    └── adapter/             # Implementations of Gateways (Redis, SMTP).
```

### Key Concepts

#### 1. Domain Gateways (`domain/gateway`)

The **Gateway Pattern** is used to invert dependencies. If the Domain needs to send an email or dispatch a background job, it defines an **Interface** in `domain/gateway`.

-   **Why?** The Domain should not depend on Redis or SMTP libraries.
-   **How?** The `infrastructure` layer implements this interface (e.g., `AsynqTaskDispatcher`), and it is injected into the Application layer at runtime.

#### 2. CQRS (Command Query Responsibility Segregation)

We separate **Write** and **Read** operations.

-   **Commands**: Modifies state. Returns `error` or `(ID, error)`. found in `app/command`.
-   **Queries**: Reads state. Returns `(DTO, error)`. found in `app/query`.
-   **Handlers**: Each Command/Query has a dedicated Handler structs (e.g., `CreateHabitHandler`).

## 3. Implementation Workflows

### Workflow: Adding a New API Feature

1.  **Contract First (OpenAPI)**

    -   Modify `api/openapi/{module}.yml`.
    -   Define the path, method, and request/response bodies.
    -   Run `make gen-api` -> Generates server interface in `internal/generated/`.

2.  **Domain Layer (The "What")**

    -   Define the Entity in `internal/{module}/domain/{entity}.go`.
    -   Define the Repository Interface in `internal/{module}/domain/repository.go`.
    -   _Rule_: Use rich domain models (methods on structs) where possible.

3.  **Application Layer (The "How")**

    -   Create a command/query struct (e.g., `CreateHabitCommand`).
    -   Create a handler struct (`CreateHabitHandler`).
    -   Inject dependencies (Repositories, Gateways) into the handler.
    -   Implement the `Handle(ctx, cmd)` method.

4.  **Infrastructure Layer (The "Mechanism")**

    -   Implement the Repository Interface in `internal/{module}/infrastructure/persistence/postgres/`.
    -   Use `sqlx` for queries. Always write explicit SQL (no ORMs like Gorm).

5.  **Interface Layer (The "Exposure")**
    -   Implement the generated OpenAPI interface in `internal/{module}/ports/http/`.
    -   Bind the HTTP request to the Command/Query.
    -   Invoke the Application Handler.
    -   Map the result to a standardized JSON response.

### Workflow: Database Migration

1.  **Generate Files**: Use standard naming `YYYYMMDD_description`.
    -   `migrations/00000X_create_table_foo.up.sql`
    -   `migrations/00000X_create_table_foo.down.sql`
2.  **Standards**:
    -   Use `uuid` for primary keys (`gen_random_uuid()`).
    -   Always include `created_at` and `updated_at` with `timestamptz`.
    -   Use `snake_case` for all database identifiers.
3.  **Execute**: Run `make migrate-up`.

## 4. Frontend Guidelines

### Banking Design System

We utilize a hybrid styling approach: **Tailwind Utility Classes** + **Custom Semantic Components**.

-   **Design Philosophy**: "Banking Grade". Robust, solid, high-contrast, professional.
-   **CSS Architecture**:
    -   `src/index.css`: Contains the "Banking Design System" layers.
    -   Use classes like `.btn-banking-primary`, `.card-banking`, `.input-banking`.
    -   Avoid using raw Tailwind (e.g., `bg-blue-500 rounded-md p-4`) for reusable UI elements.

### React Patterns

-   **State**: Use `Zustand` for global state (User session, Theme). Local state for forms.
-   **API Client**: Use the configured `axios` instance which handles auth headers automatically.
-   **Clean Code**:
    -   Isolate complex logic into custom hooks (`useHabitStats`).
    -   Keep components pure and presentational where possible.

## 5. Coding Standards & Rules

### General

-   **Errors**: Use the `internal/common/apperror` package.
    -   `apperror.New(code, msg)` -> Mapped to HTTP Status codes automatically.
-   **Context**: All blocking operations (DB, API calls) MUST take `context.Context`.
-   **Logging**: Use structured logging. Never `fmt.Println` in production code.

### "Gateway" Specifics

-   **Naming**: `TaskDispatcher` (Interface), `RedisTaskDispatcher` (Implementation).
-   **Location**: Interfaces in `domain/gateway`. Implementations in `infrastructure/adapter`.
-   **Usage**: The Application layer depends on the Interface. The Main (`cmd/api`) injects the Implementation.

## 6. Cheatsheet

| Task                 | Command                             |
| :------------------- | :---------------------------------- |
| **Generate Go Code** | `make gen-api` (after editing .yml) |
| **Run Migrations**   | `make migrate-up`                   |
| **Run Server**       | `go run cmd/api/main.go`            |
| **Run Frontend**     | `cd frontend && npm run dev`        |
| **Run Tests**        | `go test ./...`                     |
| **Lint Frontend**    | `npm run lint`                      |
