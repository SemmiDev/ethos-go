# Ethos - Modern Habit Tracker

<div align="center">
  <img src="logo.jpg" alt="Ethos Logo" width="120" height="120" />

  <h3>High-Performance Go Backend + Embedded React Frontend</h3>

  <p>
    <a href="https://go.dev"><img src="https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go" alt="Go"></a>
    <a href="https://react.dev"><img src="https://img.shields.io/badge/React-19-61DAFB?style=for-the-badge&logo=react" alt="React"></a>
    <a href="https://postgres.org"><img src="https://img.shields.io/badge/PostgreSQL-17-4169E1?style=for-the-badge&logo=postgresql" alt="PostgreSQL"></a>
    <a href="https://docker.com"><img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker" alt="Docker"></a>
  </p>

  <p>_Build better habits, one day at a time._</p>

  <p>
    <a href="#key-features">Features</a> •
    <a href="#quick-start">Quick Start</a> •
    <a href="#architecture">Architecture</a> •
    <a href="docs/screenshots/">Screenshots</a> •
    <a href="CLAUDE.md">Docs</a>
  </p>
</div>

---

## Overview

**Ethos** is a production-grade habit tracking application built to demonstrate **advanced Go patterns**, **Clean Architecture**, and **Single Binary Architecture**.

It simplifies deployment by embedding the React frontend directly into the Go binary. **One file is all you need.**

## Screenshots

<div align="center">

|                      Landing Page                      |                        Dashboard                         |
| :----------------------------------------------------: | :------------------------------------------------------: |
| <img src="docs/screenshots/landing.png" width="400" /> | <img src="docs/screenshots/dashboard.png" width="400" /> |

|                        Habits                         |                        Analytics                         |
| :---------------------------------------------------: | :------------------------------------------------------: |
| <img src="docs/screenshots/habits.png" width="400" /> | <img src="docs/screenshots/analytics.png" width="400" /> |

</div>

## Key Features

- **Habit Tracking**: Create, track, and visualize daily habits.
- **Why It Matters**: Build consistency with streaks and analytics.
- **Single Binary**: Zero-config deployment (Go + React embedded).
- **Clean Architecture**: Domain-Driven Design (DDD) & CQRS.
- **gRPC API**: High-performance gRPC services with JSON Gateway.
- **Observability**: LGTM Stack (Loki, Grafana, Tempo, Mimir).
- **Security**: JWT Auth, Rate Limiting, Secure Headers.

## Quick Start

### 1. The Easy Way (Docker)

```bash
git clone https://github.com/semmidev/ethos-go.git
cd ethos-go
cp .env.example .env

# Starts Backend, Frontend, DB, Redis, and Observability stack
make dev
```

Visit **http://localhost:8080**.

### 2. The Hacker Way (Single Binary)

```bash
# Build everything into one binary
make build

# Run it
./build/ethos-go
```

## Architecture

Ethos adheres to **Clean/Hexagonal Architecture** with **CQRS**:

```mermaid
graph TB
    subgraph "Ports (Driving Adapters)"
        GRPC["gRPC Server"]
        HTTP["HTTP/REST Gateway"]
    end

    subgraph "Application Layer"
        CMD["Commands<br/>(Write Operations)"]
        QRY["Queries<br/>(Read Operations)"]
        DEC["Decorators<br/>(Logging, Metrics)"]
    end

    subgraph "Domain Layer"
        ENT["Entities<br/>(User, Habit, Session)"]
        SVC["Domain Services"]
        PORT["Port Interfaces<br/>(Repository, Reader, Writer)"]
    end

    subgraph "Adapters (Driven)"
        REPO["PostgreSQL<br/>Repositories"]
        JWT["JWT Token<br/>Issuer"]
        HASH["BCrypt<br/>Hasher"]
        OAUTH["Google OAuth"]
        CACHE["Redis Cache"]
    end

    GRPC --> CMD
    GRPC --> QRY
    HTTP --> CMD
    HTTP --> QRY

    CMD --> DEC
    QRY --> DEC
    DEC --> ENT
    DEC --> SVC
    CMD --> PORT
    QRY --> PORT

    PORT -.->|implements| REPO
    PORT -.->|implements| JWT
    PORT -.->|implements| HASH
    PORT -.->|implements| OAUTH
    PORT -.->|implements| CACHE
```

### Project Structure

```bash
ethos-go/
├── cmd/                    # Application entry points
│   ├── server/             # Main API server
│   └── worker/             # Background job worker
├── internal/
│   ├── auth/               # Authentication module
│   ├── habits/             # Habit tracking module
│   ├── notifications/      # Notification module
│   ├── common/             # Shared utilities
│   └── generated/          # Generated gRPC code
├── api/                    # Protocol Buffer definitions
├── migrations/             # Database migrations
└── frontend/               # Embedded React application
```

### Layer Responsibilities

| Layer           | Directory   | Responsibility                                                             |
| --------------- | ----------- | -------------------------------------------------------------------------- |
| **Domain**      | `domain/`   | Business entities, rules, and port interfaces. Zero external dependencies. |
| **Application** | `app/`      | Use cases (Commands/Queries), orchestrates domain logic.                   |
| **Adapters**    | `adapters/` | Implements domain ports (PostgreSQL, JWT, BCrypt).                         |
| **Ports**       | `ports/`    | Entry points (gRPC servers, HTTP handlers).                                |
| **Service**     | `service/`  | Dependency injection and wiring.                                           |

### Design Decisions

#### 1. Domain Entity Encapsulation

Entities use **private fields** with **getters/setters** to enforce invariants:

```go
// ✅ Good: Encapsulated entity
type User struct {
    userID  uuid.UUID  // Private
    email   string
}

func (u *User) Email() string { return u.email }
func (u *User) SetEmail(email string) { u.email = email }
```

**Why?** Prevents invalid state mutations and enforces business rules at the domain level.

#### 2. Interface Segregation (ISP)

Repository interfaces are split by read/write concerns:

```go
type UserReader interface {
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByID(ctx context.Context, userID uuid.UUID) (*User, error)
}

type UserWriter interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
}

type Repository interface {
    UserReader
    UserWriter
}
```

**Why?** Query handlers only need `UserReader`, avoiding unnecessary dependencies.

#### 3. CQRS Pattern

Commands (writes) and Queries (reads) are separated into distinct handlers:

```go
// Commands: Mutate state
type Commands struct {
    Register  RegisterHandler
    Login     LoginHandler
    // ...
}

// Queries: Read-only operations
type Queries struct {
    GetProfile   GetProfileHandler
    ListSessions ListSessionsHandler
    // ...
}
```

**Why?** Enables independent scaling, clearer intent, and simpler testing.

#### 4. Decorator Pattern for Cross-Cutting Concerns

Logging, metrics, and tracing are wrapped around handlers:

```go
func NewLoginHandler(...) LoginHandler {
    return decorator.ApplyCommandResultDecorators(
        loginHandler{...},
        logger,
        metricsClient,
    )
}
```

**Why?** Keeps business logic clean; cross-cutting concerns are composable.

#### 5. Dependency Inversion

Domain defines interfaces; adapters implement them:

```go
// Domain defines the contract
type TokenIssuer interface {
    IssueAccessToken(ctx, userID, sessionID, expiresAt) (string, error)
}

// Adapter implements it
type JWTTokenIssuer struct { ... }
func (j *JWTTokenIssuer) IssueAccessToken(...) (string, error) { ... }
```

**Why?** Domain remains testable and framework-agnostic.

#### 6. Thin Controllers (Ports)

gRPC/HTTP handlers only translate requests and delegate to use cases:

```go
func (s *AuthGRPCServer) Login(ctx, req) (*LoginResponse, error) {
    cmd := command.LoginCommand{Email: req.Email, Password: req.Password}
    result, err := s.loginHandler.Handle(ctx, cmd)
    if err != nil {
        return nil, toGRPCError(err)
    }
    return &LoginResponse{AccessToken: result.AccessToken}, nil
}
```

**Why?** Keeps transport concerns separate from business logic.

See **[CLAUDE.md](CLAUDE.md)** for detailed development guidelines.

## Tech Stack

| Component | Tech                     |
| :-------- | :----------------------- |
| **Lang**  | Go 1.25+                 |
| **Web**   | React 19, Vite, Tailwind |
| **API**   | gRPC + Buf               |
| **Data**  | PostgreSQL 17, Redis 8   |
| **Ops**   | Docker, K8s, Grafana     |

## Contributing

Pull requests are welcome. For major changes, please open an issue first.

## License

[MIT](LICENSE)
