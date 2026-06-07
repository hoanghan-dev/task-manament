# Task Management API

Task Management API is a backend project built with Go, Gin, PostgreSQL, Redis, Docker, JWT Authentication, WebSocket, background worker, and integration tests.

The system supports authentication, workspace management, task management, task assignment, task status update, comments, notification jobs, realtime WebSocket events, health check API, Redis caching, and permission validation.

## Tech Stack

* Go
* Gin Web Framework
* PostgreSQL
* Redis
* Docker
* Docker Compose
* JWT Authentication
* WebSocket
* golang-migrate
* Testify
* net/http/httptest

## Features

* User register and login
* JWT authentication middleware
* Workspace management
* Task CRUD
* Assign task
* Update task status
* Comment on task
* Delete comment
* Notification worker
* Redis queue
* Redis cache
* WebSocket realtime notification
* Health check API
* Standard API response
* Standard application error handling
* Unit tests
* Integration tests
* Unauthorized tests
* Permission tests
* Dockerfile build
* Docker Compose full system run

## Project Structure

```bash
.
├── cmd
│   ├── api
│   │   └── main.go
│   └── worker
│       └── main_worker.go
│
├── internal
│   ├── app
│   │   └── app.go
│   │
│   ├── config
│   │   ├── postgres_connect_db.go
│   │   └── redis_connect_db.go
│   │
│   ├── middleware
│   │   ├── auth.go
│   │   ├── logging.go
│   │   └── request_id.go
│   │
│   ├── modules
│   │   ├── auth
│   │   ├── workspace
│   │   ├── task
│   │   ├── comment
│   │   ├── notification
│   │   └── health
│   │
│   ├── realtime
│   │   ├── client_ws.go
│   │   ├── event_ws.go
│   │   ├── hander_ws.go
│   │   ├── hub_ws.go
│   │   └── subscriber_ws.go
│   │
│   └── router
│       └── router.go
│
├── migrations
│   ├── 000001_create_users_table.up.sql
│   ├── 000002_create_workspace_table.up.sql
│   ├── 000003_create_tasks.up.sql
│   ├── 000004_data_seed.up.sql
│   ├── 000005_create_notifications_table.up.sql
│   └── 000006_create_comments_table.up.sql
│
├── pkg
│   ├── apperror
│   ├── cache
│   ├── response
│   └── utils
│
├── test
│   └── integration
│       ├── .env.test
│       ├── api_flow_test.go
│       ├── permission_test.go
│       ├── unauthorized_test.go
│       ├── test_helpers_test.go
│       └── test_setup_test.go
│
├── logs
│   └── request.log
│
├── docker-compose.yml
├── Dockerfile
├── env.example
├── go.mod
├── go.sum
└── README.md
```

## Prerequisites

Before running this project, make sure you have installed:

* Go
* Docker Desktop
* Docker Compose
* golang-migrate
* Git

Check versions:

```bash
go version
docker --version
docker compose version
migrate -version
```

## Environment Configuration

The project uses environment variables for database, Redis, JWT, and server configuration.

Create a `.env` file in the root folder.

Copy from `env.example`:

```bash
cp env.example .env
```

On Windows PowerShell:

```powershell
Copy-Item env.example .env
```

Example `.env`:

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123456
DB_NAME=task_management
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379

JWT_SECRET=your_secret_key
```

If you run the API inside Docker Compose, the database and Redis host should usually be service names:

```env
DB_HOST=postgres
REDIS_HOST=redis
```

If you run the API locally on your machine, use:

```env
DB_HOST=localhost
REDIS_HOST=localhost
```

Do not commit the real `.env` file to GitHub.

## Build Docker Image Manually

Build Docker image from `Dockerfile`:

```bash
docker build -t task-management-api .
```

Check built images:

```bash
docker images
```

Run API container manually:

```bash
docker run --rm -p 8080:8080 --env-file .env task-management-api
```

Note: If you run the API manually with `docker run`, PostgreSQL and Redis must already be running.

If PostgreSQL and Redis are running on your host machine, your `.env` may need:

```env
DB_HOST=host.docker.internal
REDIS_HOST=host.docker.internal
```

If PostgreSQL and Redis are running inside Docker Compose, it is recommended to run the full system using Docker Compose instead of manual `docker run`.

## Build with Docker Compose

Build all services defined in `docker-compose.yml`:

```bash
docker compose build
```

Build only API service:

```bash
docker compose build api
```

Build only worker service:

```bash
docker compose build worker
```

Rebuild and start all services:

```bash
docker compose up -d --build
```

This command is usually the most convenient command because it builds images and starts containers at the same time.

## Run Project with Docker Compose

Start PostgreSQL, Redis, API, and worker:

```bash
docker compose up -d
```

Start and rebuild if source code or Dockerfile changed:

```bash
docker compose up -d --build
```

Check running containers:

```bash
docker compose ps
```

View API logs:

```bash
docker compose logs -f api
```

View worker logs:

```bash
docker compose logs -f worker
```

View PostgreSQL logs:

```bash
docker compose logs -f postgres
```

View Redis logs:

```bash
docker compose logs -f redis
```

View all logs:

```bash
docker compose logs -f
```

Stop all containers:

```bash
docker compose down
```

Stop containers and remove volumes:

```bash
docker compose down -v
```

Use `docker compose down -v` when you want to reset the database completely.

## Run Database Migration

This project uses SQL migration files inside the `migrations` folder.

Migration files:

```bash
000001_create_users_table
000002_create_workspace_table
000003_create_tasks
000004_data_seed
000005_create_notifications_table
000006_create_comments_table
```

Run migration:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" up
```

Rollback all migrations:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" down
```

Rollback one migration:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" down 1
```

Check migration version:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" version
```

Force migration version if database is dirty:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" force <version>
```

Example:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" force 6
```

Then run migration again:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" up
```

## Run API Locally

Start PostgreSQL and Redis first.

Then run API:

```bash
go run cmd/api/main.go
```

The API will run at:

```bash
http://localhost:8080
```

## Run Worker Locally

The worker is separated from the API process.

Run worker:

```bash
go run cmd/worker/main_worker.go
```

The worker listens to Redis queue and handles notification jobs.

## API Base URL

```bash
http://localhost:8080/api
```

## API Routes

### Public Routes

These routes do not require authentication.

| Method | Endpoint             | Description       |
| ------ | -------------------- | ----------------- |
| POST   | `/api/auth/register` | Register new user |
| POST   | `/api/auth/login`    | Login user        |
| GET    | `/api/health/`       | Check API health  |

### Protected Routes

These routes require JWT token.

Send token in request header:

```http
Authorization: Bearer <access_token>
```

| Method | Endpoint                   | Description             |
| ------ | -------------------------- | ----------------------- |
| GET    | `/api/tasks/`              | Get all tasks           |
| GET    | `/api/tasks/:id`           | Get task by ID          |
| POST   | `/api/tasks/`              | Create task             |
| PUT    | `/api/tasks/:id`           | Update task             |
| DELETE | `/api/tasks/:id`           | Delete task             |
| PATCH  | `/api/tasks/assign`        | Assign task to user     |
| PATCH  | `/api/tasks/:id/status`    | Update task status      |
| GET    | `/api/workspaces/`         | Get workspace           |
| PUT    | `/api/workspaces/`         | Update workspace        |
| DELETE | `/api/workspaces/:wpId`    | Delete workspace        |
| POST   | `/api/tasks/:id/comments`  | Create comment for task |
| GET    | `/api/tasks/:id/comments`  | Get comments of task    |
| DELETE | `/api/comments/:commentId` | Delete comment          |
| GET    | `/api/ws/`                 | Connect WebSocket       |

## API Usage

### Health Check

```http
GET /api/health/
```

Curl:

```bash
curl http://localhost:8080/api/health/
```

Example response:

```json
{
  "success": true,
  "message": "Health check successfully",
  "data": {
    "status": "OK"
  }
}
```

The exact response may be different depending on the implementation in `internal/modules/health`.

## Authentication APIs

### Register

```http
POST /api/auth/register
```

Request body:

```json
{
  "email": "user@example.com",
  "password": "123456",
  "full_name": "Test User"
}
```

Curl:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "123456",
    "full_name": "Test User"
  }'
```

### Login

```http
POST /api/auth/login
```

Request body:

```json
{
  "email": "user@example.com",
  "password": "123456"
}
```

Curl:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "123456"
  }'
```

After login, copy the access token from response and use it for protected APIs.

Example protected header:

```http
Authorization: Bearer <access_token>
```

## Workspace APIs

### Get Workspace

```http
GET /api/workspaces/
```

Curl:

```bash
curl -X GET http://localhost:8080/api/workspaces/ \
  -H "Authorization: Bearer <access_token>"
```

### Update Workspace

```http
PUT /api/workspaces/
```

Request body:

```json
{
  "name": "Updated Workspace",
  "description": "Updated workspace description"
}
```

Curl:

```bash
curl -X PUT http://localhost:8080/api/workspaces/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "name": "Updated Workspace",
    "description": "Updated workspace description"
  }'
```

### Delete Workspace

```http
DELETE /api/workspaces/:wpId
```

Curl:

```bash
curl -X DELETE http://localhost:8080/api/workspaces/<workspace_id> \
  -H "Authorization: Bearer <access_token>"
```

## Task APIs

### Get All Tasks

```http
GET /api/tasks/
```

Curl:

```bash
curl -X GET http://localhost:8080/api/tasks/ \
  -H "Authorization: Bearer <access_token>"
```

### Get Task By ID

```http
GET /api/tasks/:id
```

Curl:

```bash
curl -X GET http://localhost:8080/api/tasks/<task_id> \
  -H "Authorization: Bearer <access_token>"
```

### Create Task

```http
POST /api/tasks/
```

Request body:

```json
{
  "title": "Learn Go",
  "description": "Practice Go backend",
  "workspace_id": "<workspace_id>"
}
```

Curl:

```bash
curl -X POST http://localhost:8080/api/tasks/ \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title": "Learn Go",
    "description": "Practice Go backend",
    "workspace_id": "<workspace_id>"
  }'
```

### Update Task

```http
PUT /api/tasks/:id
```

Request body:

```json
{
  "title": "Updated Task",
  "description": "Updated task description",
  "status": "IN_PROGRESS"
}
```

Curl:

```bash
curl -X PUT http://localhost:8080/api/tasks/<task_id> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "title": "Updated Task",
    "description": "Updated task description",
    "status": "IN_PROGRESS"
  }'
```

### Delete Task

```http
DELETE /api/tasks/:id
```

Curl:

```bash
curl -X DELETE http://localhost:8080/api/tasks/<task_id> \
  -H "Authorization: Bearer <access_token>"
```

### Assign Task

```http
PATCH /api/tasks/assign
```

Request body:

```json
{
  "task_id": "<task_id>",
  "workspace_id": "<workspace_id>",
  "assignee_id": "<user_id>"
}
```

Curl:

```bash
curl -X PATCH http://localhost:8080/api/tasks/assign \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "task_id": "<task_id>",
    "workspace_id": "<workspace_id>",
    "assignee_id": "<user_id>"
  }'
```

### Update Task Status

```http
PATCH /api/tasks/:id/status
```

Request body:

```json
{
  "status": "DONE"
}
```

Curl:

```bash
curl -X PATCH http://localhost:8080/api/tasks/<task_id>/status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "status": "DONE"
  }'
```

## Comment APIs

### Create Comment

```http
POST /api/tasks/:id/comments
```

`id` is the task ID.

Request body:

```json
{
  "content": "This is a comment"
}
```

Curl:

```bash
curl -X POST http://localhost:8080/api/tasks/<task_id>/comments \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "content": "This is a comment"
  }'
```

### Get Comments By Task

```http
GET /api/tasks/:id/comments
```

`id` is the task ID.

Curl:

```bash
curl -X GET http://localhost:8080/api/tasks/<task_id>/comments \
  -H "Authorization: Bearer <access_token>"
```

### Delete Comment

```http
DELETE /api/comments/:commentId
```

Curl:

```bash
curl -X DELETE http://localhost:8080/api/comments/<comment_id> \
  -H "Authorization: Bearer <access_token>"
```

## WebSocket

The project has realtime WebSocket support in:

```bash
internal/realtime
```

WebSocket route:

```http
GET /api/ws/
```

Connect to WebSocket endpoint:

```bash
ws://localhost:8080/api/ws/
```

Because WebSocket is protected by `AuthMiddleware`, send JWT token using the same authentication method supported by the middleware.

If the middleware reads token from the `Authorization` header, use:

```http
Authorization: Bearer <access_token>
```

When a notification event occurs, the server can send realtime event data to the connected user.

Example event:

```json
{
  "type": "task.assigned",
  "data": {
    "sender_id": "<sender_id>",
    "receiver_id": "<receiver_id>",
    "task_id": "<task_id>"
  }
}
```

The project also has update status event support:

```bash
internal/modules/task/events/update_status_event.go
internal/modules/notification/jobs/update_status_job.go
```

## Run Unit Tests

Run all tests:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

Run specific module test:

```bash
go test -v ./internal/modules/auth/services
```

```bash
go test -v ./internal/modules/workspace/services
```

```bash
go test -v ./internal/modules/task/services
```

```bash
go test -v ./internal/modules/comment/services
```

## Run Integration Tests

Integration tests are located in:

```bash
test/integration
```

Test files:

```bash
api_flow_test.go
permission_test.go
unauthorized_test.go
test_helpers_test.go
test_setup_test.go
```

The integration test environment file is:

```bash
test/integration/.env.test
```

Run all integration tests:

```bash
go test -v ./test/integration
```

Run API flow test:

```bash
go test -v ./test/integration -run TestAPIFlow
```

Run unauthorized test:

```bash
go test -v ./test/integration -run TestUnauthorized
```

Run permission test:

```bash
go test -v ./test/integration -run TestPermission
```

If the test uses a separate test database, make sure the test database exists and migration has been run for that database.

Example test database migration:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management_test?sslmode=disable" up
```

## Run Race Detector

Run race detector for all packages:

```bash
go test -race ./...
```

Run race detector for integration tests:

```bash
go test -race -v ./test/integration
```

## Logs

Request logs are stored in:

```bash
logs/request.log
```

Integration test logs are stored in:

```bash
test/integration/logs/request.log
```

## Error Handling

The project uses standard application errors in:

```bash
pkg/apperror
```

Files:

```bash
app_error.go
handler.go
postgres.go
```

The project also uses standard API response helpers in:

```bash
pkg/response/api_response.go
```

Common error types:

* Bad Request
* Unauthorized
* Forbidden
* Not Found
* Conflict
* Validation Error
* Internal Server Error

## Redis Usage

Redis is used for:

* Caching task data
* Queueing notification jobs
* Worker processing

Redis cache code is located in:

```bash
pkg/cache/redis_cache.go
```

Redis connection code is located in:

```bash
internal/config/redis_connect_db.go
```

## Worker

The worker entry point is:

```bash
cmd/worker/main_worker.go
```

Notification worker code is located in:

```bash
internal/modules/notification/worker/notification_worker.go
```

The worker handles jobs from Redis queue, such as:

* Assign task notification
* Update task status notification

## Docker Commands Summary

### Build Docker image manually

```bash
docker build -t task-management-api .
```

### Run Docker image manually

```bash
docker run --rm -p 8080:8080 --env-file .env task-management-api
```

### Build all Docker Compose services

```bash
docker compose build
```

### Build only API service

```bash
docker compose build api
```

### Build only worker service

```bash
docker compose build worker
```

### Start full system

```bash
docker compose up -d
```

### Rebuild and start full system

```bash
docker compose up -d --build
```

### Stop full system

```bash
docker compose down
```

### Stop full system and remove volumes

```bash
docker compose down -v
```

### View all logs

```bash
docker compose logs -f
```

### View API logs

```bash
docker compose logs -f api
```

### View worker logs

```bash
docker compose logs -f worker
```

### View PostgreSQL logs

```bash
docker compose logs -f postgres
```

### View Redis logs

```bash
docker compose logs -f redis
```

### Remove unused Docker images

```bash
docker image prune
```

### Remove unused Docker resources

```bash
docker system prune
```

## Common Commands

### Run API locally

```bash
go run cmd/api/main.go
```

### Run worker locally

```bash
go run cmd/worker/main_worker.go
```

### Run all tests

```bash
go test ./...
```

### Run all tests with verbose output

```bash
go test -v ./...
```

### Run integration tests

```bash
go test -v ./test/integration
```

### Run race detector

```bash
go test -race ./...
```

### Run migration

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" up
```

## Common Issues

### Error: relation "users" does not exist

The database tables have not been created yet.

Run migration:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" up
```

### Docker Compose warning: DB_NAME variable is not set

Make sure `.env` exists in the root folder and contains the required variables.

Check `.env`:

```bash
cat .env
```

On Windows PowerShell:

```powershell
Get-Content .env
```

### API cannot connect to PostgreSQL

Check if PostgreSQL container is running:

```bash
docker compose ps
```

Check database logs:

```bash
docker compose logs -f postgres
```

Also verify these variables:

```env
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
```

If the API runs inside Docker Compose, use:

```env
DB_HOST=postgres
```

If the API runs locally, use:

```env
DB_HOST=localhost
```

### API cannot connect to Redis

Check Redis container:

```bash
docker compose ps
```

Check Redis logs:

```bash
docker compose logs -f redis
```

If the API runs inside Docker Compose, use:

```env
REDIS_HOST=redis
```

If the API runs locally, use:

```env
REDIS_HOST=localhost
```

### Port is already in use

Stop old containers:

```bash
docker compose down
```

Or check which process is using the port.

On Windows:

```powershell
netstat -ano | findstr :8080
```

### Dirty database migration

Check version:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" version
```

Force version:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" force <version>
```

Then run migration again:

```bash
migrate -path migrations -database "postgres://postgres:123456@localhost:5432/task_management?sslmode=disable" up
```

### Docker image does not update after code changes

Rebuild image:

```bash
docker compose up -d --build
```

Or rebuild manually:

```bash
docker build -t task-management-api .
```

### Reset database completely

```bash
docker compose down -v
docker compose up -d --build
```

Then run migration again if migration is not automatically executed.

## Development Notes

* Do not commit `.env`.
* Use `env.example` to document required environment variables.
* Run migration before using the API.
* Run tests before creating a pull request.
* Use Docker Compose for easier setup.
* API and worker are two separate processes.
* Integration tests should use `test/integration/.env.test`.
* Permission tests must verify that one user cannot update, delete, or comment on another user's private task.
* Unauthorized tests must verify missing or invalid token returns `401 Unauthorized`.

## Suggested Git Commit

```bash
git add README.md
git commit -m "docs: improve README setup run test and API usage"
```
