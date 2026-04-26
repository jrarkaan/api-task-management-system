ifneq (,$(wildcard .env))
include .env
export
endif

DB_MIGRATE_URL ?= postgres://postgres:postgres@localhost:5432/task_management_system?sslmode=disable
VERSION ?= 1

.PHONY: run build tidy migrate-up migrate-down migrate-force docker-up docker-down docker-logs

run:
	go run ./cmd/http

build:
	go build ./cmd/http

tidy:
	go mod tidy

migrate-up:
	migrate -path migrations -database "$(DB_MIGRATE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DB_MIGRATE_URL)" down

migrate-force:
	migrate -path migrations -database "$(DB_MIGRATE_URL)" force $(VERSION)

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f
