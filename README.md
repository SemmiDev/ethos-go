# Ethos - Modern Habit Tracker

### High-Performance Go Backend + Embedded React Frontend

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)
![React](https://img.shields.io/badge/React-18-61DAFB?style=flat&logo=react)
![Clean Architecture](https://img.shields.io/badge/Architecture-Clean%2FDDD-orange)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)
![Status](https://img.shields.io/badge/Status-Production%20Ready-success)

_Build better habits, one day at a time._

[Features](#key-features) •
[Architecture](#architecture) •
[Getting Started](#getting-started) •
[Deployment](#deployment) •
[Documentation](#documentation)

---

## Overview

**Ethos** is a production-grade habit tracking application built to demonstrate **advanced Go patterns**, **Domain-Driven Design (DDD)**, and **Single Binary Architecture**.

It combines a robust Go backend with a modern React frontend into a **single executable**, making deployment as simple as copying one file. No external web servers or confusing CORS configurations required.

## Screenshots

<div align="center">

|               Landing Page               |                  Dashboard                   |
| :--------------------------------------: | :------------------------------------------: |
| ![Landing](docs/screenshots/landing.png) | ![Dashboard](docs/screenshots/dashboard.png) |

|                 Habits                 |                  Analytics                   |
| :------------------------------------: | :------------------------------------------: |
| ![Habits](docs/screenshots/habits.png) | ![Analytics](docs/screenshots/analytics.png) |

</div>

## Key Features

### Core Functionality

- **Habit Tracking**: Create, track, and visualize daily habits.
- **Smart Analytics**: Insightful charts and streaks functionality.
- **Notifications**: Email and push notifications (SMTP).
- **Gamification**: Earn badges and track progress.
- **Command Palette**: `Cmd+K` interface for power users.
- **PWA Support**: Installable on mobile and desktop.

### Technical Highlights

- **Single Binary**: Frontend embedded into Go binary (Zero-config deployment).
- **Clean Architecture**: Strict separation of concerns (Domain, Application, Infrastructure).
- **CQRS**: Command Query Responsibility Segregation for scalability.
- **Observability**: Full LGTM Stack (Loki, Grafana, Tempo, Mimir) + OpenTelemetry.
- **Concurrency**: Worker pool pattern for background tasks (emails, data processing).
- **Security**: JWT Authentication, Rate Limiting, and Secure Headers.

## Tech Stack

| Component          | Technology                                      |
| ------------------ | ----------------------------------------------- |
| **Language**       | Go 1.25+                                        |
| **Frontend**       | React, Vite, Framer Motion, TailwindCSS         |
| **Database**       | PostgreSQL 17                                   |
| **Caching**        | Redis 8                                         |
| **Observability**  | OpenTelemetry, Grafana, Loki, Tempo, Prometheus |
| **Infrastructure** | Docker, Kubernetes, Makefile                    |
| **Routing**        | Chi Router                                      |
| **Documentation**  | Swagger/OpenAPI                                 |

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 20+ (for frontend development)
- Docker & Docker Compose

### Quick Start (Single Command)

Most easy way to start everything (Backend + Frontend + DB + Monitoring):

```bash
# Clone the repo
git clone https://github.com/semmidev/ethos-go.git
cd ethos-go

# Create env file
cp .env.example .env

# Start everything with Docker Compose
docker-compose -f compose.dev.yml up -d --build
```

Visit **http://localhost:8080** to see the app!

### Development Mode (Hot Reload)

If you want to work on the code:

**1. Backend**

```bash
make run
# Runs on localhost:8080
```

**2. Frontend** (in a new terminal)

```bash
cd frontend
npm install
npm run dev
# Runs on localhost:3001 (proxies API requests to 8080)
```

## Architecture

Ethos strictly follows **Clean Architecture** principles as detailed in [CLAUDE.md](CLAUDE.md):

### Folder Structure

```bash
ethos-go/
├── cmd/api/            # Main entry point
├── internal/
│   ├── {module}/       # Feature Modules (Auth, Habits, etc.)
│   │   ├── domain/     # Logic & Gateways
│   │   ├── app/        # CQRS (Commands/Queries)
│   │   ├── infra/      # DB & Adapters
│   │   └── ports/      # HTTP Handlers
│   ├── common/         # Shared kernels
│   └── generated/      # Generated Code (OpenAPI, SQLC)
├── frontend/           # React Application
├── k8s/                # Kubernetes Manifests
└── deployments/        # Dockerfiles
```

For detailed coding guidelines and workflows, see **[CLAUDE.md](CLAUDE.md)**.

## Single Binary Build

We use Go's `embed` package to bundle the React frontend into the binary.

```bash
# 1. Build Frontend + Backend into one binary
make build

# 2. Run it
./build/ethos-go

# 3. Access App
# open http://localhost:8080
```

This single 25MB binary contains **everything** needed to run the application (except external DB/Redis).

## Observability

Ethos comes with a pre-configured **Grafana Stack**.
Access Grafana at **http://localhost:3000** (user: `admin`, pass: `admin`).

- **Dashboards**: Pre-built dashboards for Go Runtime, HTTP Metrics, and Business Logic.
- **Traces**: End-to-end request tracing with Tempo.
- **Logs**: Structured logging with Loki.

## Documentation

Detailed documentation is available in the `docs/` folder:

- [Agent Guidelines (CLAUDE.md)](CLAUDE.md) - **Primary source of truth for development.**
- [Backend Guide](GUIDE.md) - Deep dive into patterns.
- [API Documentation](docs/openapi/) - OpenAPI/Swagger specs.
- [Production Readiness](docs/PRODUCTION_READINESS.md) - Checklist for deployment.

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
  <sub>Built with ❤️ by the Sammi Aldhi Yanto</sub>
</div>
