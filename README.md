# Task Management Go Backend

A robust, enterprise-ready Go backend built using clean architecture for project task management. It features workspace organization, secure authentication, multi-user task management, commenting, real-time updates via WebSockets, and asynchronous background worker processing for notification dispatching.

---

## Features

- **Authentication**: JWT-based secure user authentication (Register, Login, Password validation, Token parsing).
- **Workspace Management**: Multi-tenant workspace separation. Every user gets a default workspace on registration, with full support to read, update, or delete workspaces.
- **Task Management**: Create, view, update, and delete tasks within workspaces, transition task statuses, and assign tasks to workspace members.
- **Comments System**: Collaborative thread under tasks to write, fetch, and delete task comments.
- **Real-time WebSockets**: Push real-time event updates to connected clients using a centralized WebSocket Hub.
- **Asynchronous Background Workers**: Redis Pub/Sub combined with worker goroutines to process notification queues (task assignment, status changes) cleanly off the main thread.
- **Robust Error Handling**: Standardized central middleware for mapping database and application errors to precise client HTTP status codes (400, 401, 403, 404, 409, 500) and structured error payloads.

---

## Tech Stack

- **Core**: Go (Golang) v1.26.3
- **Web Framework**: Gin Gonic v1.12.0
- **Database Driver**: PostgreSQL (via pgx/v5 & database/sql stdlib wrapper)
- **Caching & Message Broker**: Redis v8 (using go-redis/v9)
- **Authentication**: JWT (JSON Web Tokens via golang-jwt/v4) & bcrypt hashing
- **Real-Time Pub/Sub**: WebSocket (gorilla/websocket) & Redis Channels
- **Database Migrations**: golang-migrate
- **Containerization**: Docker & Docker Compose

---

## Project Structure

```txt
cmd/
  ├── api/                 # API server entrypoint (runs Gin engine on :8080)
  └── worker/              # Background notification worker process
internal/
  ├── app/                 # Dependency injection and application bootstrapping
  ├── config/              # PostgreSQL & Redis connection configs and .env loading
  ├── middleware/          # JWT Auth, Logger, Request ID and recovery middlewares
  ├── modules/             # Core business modules containing handler, service, repo, and DTOs
  │     ├── auth/
  │     ├── comment/
  │     ├── notification/
  │     └── task/
  │     └── workspace/
  ├── realtime/            # WebSocket hub, subscribers, and ws request handlers
  └── router/              # Gin engine routing structure (public vs protected groups)
migrations/                # Up/Down SQL database migration & seed files
pkg/
  ├── apperror/            # Standard AppError struct, Postgres mapper, & Central HTTP error handler
  ├── cache/               # Redis connection wrapper for general caching
  ├── response/            # Structured API response utilities
  └── utils/               # Token utilities and context helpers (e.g. UUID parses, extracting claims)
```

---

## Prerequisites

Ensure you have the following installed on your machine before setting up the project:
- **Go**: Version `1.26.3` or higher ([Download](https://go.dev/dl/))
- **Docker & Docker Compose**: For launching PostgreSQL and Redis instantly ([Download](https://www.docker.com/products/docker-desktop/))
- **golang-migrate CLI**: (Optional, for running database migrations manually) ([Installation Guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))

---

## Environment Variables

The application requires a `.env` file in the root directory to run. Create a file named `.env` in the project root and populate it with the appropriate values:

```env
# Database Settings
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123456
DB_NAME=task_management
DB_SSLMODE=disable

# Authentication Secrets
SECRET_KEY=ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbl

# Redis Settings
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

> [!WARNING]
> Do not commit the `.env` file to your version control. Always keep environment secrets safe.

---

## Setup

Follow these steps to configure and boot up your local development environment:

### 1. Clone the repository and navigate to the directory
```bash
git clone <your-repository-url>
cd task-manament
```

### 2. Download Go module dependencies
```bash
go mod download
```

### 3. Start Database & Redis Services
Using Docker Compose is the recommended way to get the database and cache layers running:
```bash
docker compose up -d
```
This boots up a PostgreSQL 16 database and a Redis 8 container with port mappings `5432` and `6379` respectively.

### 4. Run Database Migrations & Seeds
Using the standard `golang-migrate` tool, run the database migrations and data seeds:
```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" up
```
This creates all the necessary tables (users, workspaces, tasks, notifications, comments) and seeds mock data for testing.

---

## Run Application

This project splits workload into two separate run paths: the HTTP API server and the background worker for notification queues.

### 1. Run the HTTP API Server
Open a terminal and run the API server. By default, it will start on port `8080`:
```bash
go run ./cmd/api
```

### 2. Run the Background Notification Worker
Open a second terminal window and run the background worker to consume events from the Redis queue and insert notifications to PostgreSQL:
```bash
go run ./cmd/worker
```

---

## Run Tests

### Run all tests in the workspace:
```bash
go test ./...
```

### Run tests and check for race conditions:
```bash
go test -race ./...
```

### Run tests with coverage:
```bash
go test -cover ./...
```

---

## API Testing

You can interact with the API endpoints using **curl**, **Postman**, or **Thunder Client**. Below are examples of primary actions:

### 1. Register User (Public)
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecurePassword123!",
    "full_name": "Test User"
  }'
```

### 2. Login User (Public)
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "SecurePassword123!"
  }'
```
*Note: Save the `token` from the response to authorize protected endpoints.*

### 3. Get Workspace (Protected)
```bash
curl -X GET http://localhost:8080/api/workspaces/ \
  -H "Authorization: Bearer <your_jwt_token>"
```

### 4. Create Task (Protected)
```bash
curl -X POST http://localhost:8080/api/tasks/ \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "workspace_id": "<workspace-uuid>",
    "title": "Build API Error Test Suite",
    "description": "Create automated tests for error status mapping",
    "status": "TODO",
    "priority": "HIGH",
    "due_date": "2026-06-30T12:00:00Z"
  }'
```

---

## Common Commands

| Command | Description |
|---|---|
| `go mod download` | Downloads all external dependencies listed in `go.mod` |
| `docker compose up -d` | Boots Postgres and Redis containers in detached mode |
| `docker compose down` | Stops and removes active container instances |
| `go run ./cmd/api` | Starts the HTTP Gin REST API server |
| `go run ./cmd/worker` | Starts the background notification subscription queue worker |
| `go test ./...` | Runs all test suites in the codebase |
| `go test -race ./...` | Runs test suites with Go race detector enabled |

---

## Migration Commands

Below are standard operations when using `golang-migrate`:

### Create a new migration file:
```bash
migrate create -ext sql -dir migrations -seq create_new_table
```

### Apply all up migrations:
```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" up
```

### Rollback the last migration:
```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" down 1
```

### Force migration state to a specific version (useful if dirty state occurs):
```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" force <version_number>
```

---

## Troubleshooting

### 1. Database Connection Refused
- **Symptom**: `Connet database faild: dial tcp 127.0.0.1:5432: connect: connection refused`
- **Fix**: Check if PostgreSQL container is running using `docker ps`. If not, run `docker compose up -d`. Check if port `5432` is occupied by another local PostgreSQL installation.

### 2. Redis Connection Refused
- **Symptom**: `Connet redis faild: dial tcp 127.0.0.1:6379: connect: connection refused`
- **Fix**: Verify Redis container is active. Check that your `.env` contains the correct `REDIS_ADDR=localhost:6379` and matches the exposed Docker compose ports.

### 3. Migration Dirty Database Error
- **Symptom**: `Dirty database version X. Fix and force version.`
- **Fix**: Open database terminal or client tool, resolve any SQL syntax errors manually, then run:
  `migrate -path migrations -database "postgres://..." force <last_successful_version>`

### 4. WebSocket Fails to Handshake
- **Symptom**: Clients cannot connect to `/api/ws/`
- **Fix**: The WebSocket endpoint `/api/ws/` is **protected** by authentication middleware. Ensure you are passing the JWT token correctly in the request (e.g. via token query parameter or standard HTTP Authorization headers if using a compliant client).

---

## Notes

- **Secrets Management**: Never add or push raw credentials into the repository. The `.env` file should remain excluded in `.gitignore`.
- **Startup Order**: Always run `docker compose up -d` and ensure your database migrations are fully configured before starting the main API process or worker.
- **Worker Execution**: The main API server delegates notifications through a Redis channel queue. Running the worker process (`cmd/worker`) is essential for those notifications to resolve and be stored/emitted.