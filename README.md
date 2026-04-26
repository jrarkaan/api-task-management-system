# Task Management System API

A simple RESTful backend API for a personal Task Management System. This project supports user registration, login with JWT authentication, and user-owned task CRUD operations.

The backend is built with Go, Gin, GORM, PostgreSQL, golang-migrate, Zap logger, and Swagger/OpenAPI documentation.

## Quick Start for Reviewer

> No `.env` setup required — it is generated automatically.

```bash
make docker-start
```

This single command will:

1. Auto-generate `.env` with Docker-friendly defaults (if it does not exist)
2. Start PostgreSQL and wait for it to be healthy
3. Run all database migrations
4. Start the backend

The API will be available at:

```text
http://localhost:8080
```

Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

Health check:

```bash
curl http://localhost:8080/health
```

To stop and remove containers:

```bash
make docker-down
```

To wipe the database and start fresh:

```bash
make docker-reset
```

## Tech Stack

- Go 1.25
- Gin
- GORM
- PostgreSQL 16
- JWT authentication
- bcrypt password hashing
- golang-migrate
- Zap structured logger
- Swagger/OpenAPI via swaggo
- Docker and Docker Compose

## Features

- User registration
- User login with JWT token
- Protected task routes using `Authorization: Bearer <token>`
- Create, list, update, and delete tasks
- Task ownership validation per authenticated user
- Task filtering by status
- Task pagination with `page` and `limit`
- Partial task update
- Standardized API response format
- SQL migrations using golang-migrate
- Swagger UI documentation
- CORS enabled for development with:
  - Allowed origins: *
  - Allowed methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
  - Authorization header supported

## Folder Structure

```text
cmd/http                 Application entry point
app/config               Environment configuration
app/driver               HTTP server and route wiring
app/middleware           Gin middleware for auth, request logging, and recovery
modules/accounts/v1      Register and login module
modules/tasks/v1         Task CRUD module
pkg/apiresponse          Standard API response helpers and Swagger DTOs
pkg/db                   Transaction manager
pkg/db/pg                PostgreSQL connection
pkg/helpers              JWT, password, and UUID helpers
pkg/logger               Zap logger wrapper
pkg/pagination           Pagination helper
pkg/xvalidator           Request validation helper
migrations               golang-migrate SQL files
docs                     Generated Swagger/OpenAPI files
```

## Prerequisites

- Go 1.25, recommended via `gvm`
- Docker and Docker Compose
- golang-migrate CLI
- swag CLI, only needed when regenerating Swagger docs

## Environment Variables

Create `.env` from the example file:

```bash
cp .env.example .env
```

Default local development values:

```env
APP_NAME=task-management-system
APP_ENV=local
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=task_management_system
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Jakarta

DB_MIGRATE_URL=postgres://postgres:postgres@localhost:5432/task_management_system?sslmode=disable

JWT_SECRET=change-me
JWT_EXPIRES_HOURS=24
```

### Local app vs Docker app DB host

Use `DB_HOST=localhost` when running the Go app directly on your machine:

```bash
make run
```

Use `DB_HOST=postgres` when running the backend service inside Docker Compose:

```env
DB_HOST=postgres
DB_MIGRATE_URL=postgres://postgres:postgres@postgres:5432/task_management_system?sslmode=disable
```

If migrations are run from your host machine, keep `DB_MIGRATE_URL` pointed to `localhost`.

## Install Tools

### Use Go 1.25 with gvm

```bash
gvm use go1.25
```

The repository is already initialized as:

```bash
module api-task-management-system
```

Run module tidy after pulling or changing dependencies:

```bash
go mod tidy
```

### Install golang-migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Make sure your Go binary path is available:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Install swag CLI

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## Run Locally

Start PostgreSQL only:

```bash
docker compose up -d postgres
```

Run migrations:

```bash
make migrate-up
```

Run the API:

```bash
make run
```

The API will be available at:

```text
http://localhost:8080
```

Health check:

```bash
curl http://localhost:8080/health
```

## Run With Docker Compose

The easiest way to run the full stack (PostgreSQL + migrations + backend) in one command:

```bash
make docker-start
```

`.env` is auto-generated if it does not exist. The generated file uses `DB_HOST=postgres` and a `DB_MIGRATE_URL` pointing to the `postgres` service, which are required for Docker Compose networking.

View logs:

```bash
make docker-logs
```

Stop containers (keep volumes):

```bash
make docker-down
```

Stop containers and delete database data:

```bash
make docker-reset
```

Restart everything:

```bash
make docker-restart
```

## Database Migration

Run all pending migrations:

```bash
make migrate-up
```

Rollback one migration step:

```bash
make migrate-down
```

Force migration version if needed:

```bash
make migrate-force VERSION=1
```

Migration files are located in:

```text
migrations/
```

## Swagger API Documentation

Generate Swagger docs:

```bash
make swagger
```

Run the app:

```bash
make run
```

Open Swagger UI:

```text
http://localhost:8080/swagger/index.html
```

For protected endpoints, click **Authorize** and enter:

```text
Bearer <your_jwt_token>
```

## API Endpoints

### Health

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/health` | No | Health check |

### Auth

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/accounts/v1/auth/register` | No | Register a new user |
| POST | `/accounts/v1/auth/login` | No | Login and get JWT token |

### Tasks

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/tasks` | Yes | List authenticated user's tasks |
| GET | `/tasks?status=pending` | Yes | Filter tasks by status |
| GET | `/tasks?page=1&limit=10` | Yes | List tasks with pagination |
| GET | `/tasks?status=done&page=1&limit=10` | Yes | Filter and paginate tasks |
| POST | `/tasks` | Yes | Create a task |
| PUT | `/tasks/{id}` | Yes | Partially update a task by UUID |
| DELETE | `/tasks/{id}` | Yes | Delete a task by UUID |

Task routes require:

```http
Authorization: Bearer <token>
```

Logout is handled on the client side by removing the stored token.

## Response Format

Success response:

```json
{
  "meta": null,
  "message": "Success",
  "status": 1,
  "data": {}
}
```

Created response:

```json
{
  "meta": null,
  "message": "Created successfully",
  "status": 1,
  "data": {}
}
```

Error response:

```json
{
  "message": "",
  "status": 0,
  "error": {
    "code": 400,
    "message": "validation error message",
    "status": true
  }
}
```

Paginated list response example:

```json
{
  "meta": {
    "page": 1,
    "limit": 10,
    "total_rows": 25,
    "total_pages": 3,
    "has_next": true,
    "has_prev": false
  },
  "message": "Success",
  "status": 1,
  "data": []
}
```

## Sample curl Commands

### Register

```bash
curl -X POST http://localhost:8080/accounts/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"secret123"}'
```

### Login

```bash
curl -X POST http://localhost:8080/accounts/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"secret123"}'
```

Copy the token from the login response and use it in protected task requests.

### Create Task

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Write test API","description":"Finish the technical test","status":"pending","deadline":"2026-05-01"}'
```

### List Tasks

```bash
curl http://localhost:8080/tasks \
  -H "Authorization: Bearer <token>"
```

### Filter Tasks by Status

```bash
curl "http://localhost:8080/tasks?status=in-progress" \
  -H "Authorization: Bearer <token>"
```

### List Tasks With Pagination

```bash
curl "http://localhost:8080/tasks?page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

### Filter Tasks With Pagination

```bash
curl "http://localhost:8080/tasks?status=done&page=1&limit=10" \
  -H "Authorization: Bearer <token>"
```

### Update Only Status

```bash
curl -X PUT http://localhost:8080/tasks/<task_uuid> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"status":"done"}'
```

### Update Only Title

```bash
curl -X PUT http://localhost:8080/tasks/<task_uuid> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"New title"}'
```

### Update Title and Deadline

```bash
curl -X PUT http://localhost:8080/tasks/<task_uuid> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Submit technical test","deadline":"2026-04-30"}'
```

### Invalid Empty Update Body

```bash
curl -X PUT http://localhost:8080/tasks/<task_uuid> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{}'
```

Expected response:

```json
{
  "message": "",
  "status": 0,
  "error": {
    "code": 400,
    "message": "at least one field must be provided",
    "status": true
  }
}
```

### Delete Task

```bash
curl -X DELETE http://localhost:8080/tasks/<task_uuid> \
  -H "Authorization: Bearer <token>"
```

## Makefile Commands

### Docker (Reviewer Workflow)

| Command | Description |
|---|---|
| `make docker-start` | Auto-generate `.env`, run migrations, start backend |
| `make docker-down` | Stop and remove containers (keep volumes) |
| `make docker-reset` | Stop and remove containers **and** volumes (wipes DB) |
| `make docker-stop` | Gracefully stop containers without removing them |
| `make docker-restart` | `docker-down` then `docker-start` |
| `make docker-logs` | Follow Docker Compose logs |
| `make docker-up` | Start all services (build backend image) |

### Local Development

| Command | Description |
|---|---|
| `make env` | Auto-generate `.env` if it does not exist |
| `make run` | Run API locally with nodemon hot-reload |
| `make build` | Build application binary |
| `make tidy` | Run `go mod tidy` |
| `make swagger` | Generate Swagger docs |
| `make swagger-fmt` | Format Swagger annotations |
| `make migrate-up` | Apply database migrations |
| `make migrate-down` | Rollback database migrations |
| `make migrate-force VERSION=1` | Force migration version |

## Notes

- Task IDs in update and delete endpoints use task UUID.
- Task `status` must be one of `pending`, `in-progress`, or `done`.
- Task `deadline` format must be `YYYY-MM-DD`.
- Update task supports partial update, but at least one field must be provided.
- Passwords are hashed using bcrypt.
- JWT token expiration is configured through `JWT_EXPIRES_HOURS`.
