# Task Management System

Simple Go backend API for a task management technical test. It provides JWT authentication and user-owned task CRUD with PostgreSQL persistence.

## Tech Stack

- Go 1.25
- Gin
- GORM
- PostgreSQL 16
- JWT
- bcrypt
- golang-migrate
- Docker and Docker Compose

## Folder Structure

```text
cmd/http                 Application entry point
app/config               Environment configuration
app/driver               HTTP server and route wiring
app/middleware           Gin middleware
modules/accounts/v1      Register and login module
modules/tasks/v1         Task CRUD module
pkg/apiresponse          Standard API response helpers
pkg/db                   PostgreSQL connection
pkg/helpers              JWT, password, and UUID helpers
pkg/xvalidator           Request validation helper
migrations               golang-migrate SQL files
```

## Prerequisites

- Go 1.25 via gvm
- PostgreSQL 16
- golang-migrate CLI
- Docker and Docker Compose, optional

## Setup With gvm

```bash
gvm use go1.25
go mod init api-task-management-system
go mod tidy
```

This repository is already initialized with module `api-task-management-system`, so run `go mod init` only when starting from an empty directory.

## Install golang-migrate CLI on Ubuntu

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Make sure your Go binary path is available:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Local Run

```bash
cp .env.example .env
docker compose up -d postgres
make migrate-up
make run
```

For local non-Docker runs, keep `DB_HOST=localhost` in `.env`.

## Docker Run

```bash
cp .env.example .env
```

When running the backend inside Docker Compose, set this in `.env`:

```env
DB_HOST=postgres
```

Then run:

```bash
make docker-up
```

If you run migrations from your host machine, keep `DB_MIGRATE_URL=postgres://postgres:postgres@localhost:5432/task_management_system?sslmode=disable`.

## API Endpoints

- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/tasks`
- `GET /api/tasks?status=pending`
- `GET /api/tasks?status=in-progress`
- `GET /api/tasks?status=done`
- `POST /api/tasks`
- `PUT /api/tasks/:id`
- `DELETE /api/tasks/:id`

Task routes require:

```http
Authorization: Bearer <token>
```

Logout is handled client-side by removing the token.

## Response Format

Success:

```json
{
  "meta": null,
  "message": "Success",
  "status": 1,
  "data": {}
}
```

Created:

```json
{
  "meta": null,
  "message": "Created successfully",
  "status": 1,
  "data": {}
}
```

Error:

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

## Sample curl

Register:

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Doe","email":"jane@example.com","password":"secret123"}'
```

Login:

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"jane@example.com","password":"secret123"}'
```

Create task:

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Write test API","description":"Finish the technical test","status":"pending","deadline":"2026-05-01"}'
```

List tasks:

```bash
curl http://localhost:8080/api/tasks \
  -H "Authorization: Bearer <token>"
```

Filter tasks by status:

```bash
curl "http://localhost:8080/api/tasks?status=in-progress" \
  -H "Authorization: Bearer <token>"
```

Update task:

```bash
curl -X PUT http://localhost:8080/api/tasks/<task_uuid> \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"Write test API","description":"Update docs and code","status":"done","deadline":"2026-05-02"}'
```

Delete task:

```bash
curl -X DELETE http://localhost:8080/api/tasks/<task_uuid> \
  -H "Authorization: Bearer <token>"
```

## Notes

- Passwords are hashed with bcrypt.
- JWT contains the internal numeric `user_id`.
- Task `user_id` is never accepted from request bodies.
- Users can only read, update, and delete their own tasks.
- Task UUIDs are used in public API paths.
