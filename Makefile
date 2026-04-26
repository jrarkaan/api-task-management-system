ifneq (,$(wildcard .env))
include .env
export
endif

DB_MIGRATE_URL ?= postgres://postgres:postgres@localhost:5432/task_management_system?sslmode=disable
PROD_COMPOSE=docker compose --env-file .env.production -f docker-compose.cloudflared.yml
VERSION ?= 1

.PHONY: run build tidy swagger swagger-fmt \
	migrate-up migrate-down migrate-force \
	env \
	docker-up docker-start docker-stop docker-down docker-reset docker-restart docker-logs \
	prod-up prod-down prod-restart prod-logs prod-migrate prod-build

# ─── Local Development ───────────────────────────────────────────────────────

run:
	nodemon --exec go run cmd/http/main.go --signal SIGTERM

swagger:
	swag init -g cmd/http/main.go -o docs

swagger-fmt:
	swag fmt

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

# ─── Environment ─────────────────────────────────────────────────────────────

env:
	@if [ ! -f .env ]; then \
		echo "Generating .env with Docker-friendly defaults..."; \
		printf 'APP_NAME=task-management-system\n' > .env; \
		printf 'APP_ENV=local\n' >> .env; \
		printf 'APP_PORT=8080\n' >> .env; \
		printf '\n' >> .env; \
		printf 'DB_HOST=postgres\n' >> .env; \
		printf 'DB_PORT=5432\n' >> .env; \
		printf 'DB_USER=postgres\n' >> .env; \
		printf 'DB_PASSWORD=postgres\n' >> .env; \
		printf 'DB_NAME=task_management_system\n' >> .env; \
		printf 'DB_SSLMODE=disable\n' >> .env; \
		printf 'DB_TIMEZONE=Asia/Jakarta\n' >> .env; \
		printf '\n' >> .env; \
		printf 'DB_MIGRATE_URL=postgres://postgres:postgres@postgres:5432/task_management_system?sslmode=disable\n' >> .env; \
		printf '\n' >> .env; \
		printf 'JWT_SECRET=change-me\n' >> .env; \
		printf 'JWT_EXPIRES_HOURS=24\n' >> .env; \
		echo ".env generated."; \
	else \
		echo ".env already exists, skipping generation."; \
	fi

# ─── Docker Compose ───────────────────────────────────────────────────────────

docker-start: env
	docker compose up -d postgres
	docker compose run --rm migrate
	docker compose up -d backend
	@echo ""
	@echo "✅ Backend running at     http://localhost:8080"
	@echo "📄 Swagger UI running at  http://localhost:8080/swagger/index.html"

docker-stop:
	docker compose stop

docker-down:
	docker compose down

docker-reset:
	docker compose down -v

docker-restart: docker-down docker-start

docker-up:
	docker compose up -d --build

docker-logs:
	docker compose logs -f

# ─── Production (Cloudflare Tunnel) ──────────────────────────────────────────

prod-up:
	$(PROD_COMPOSE) up -d --build

prod-down:
	$(PROD_COMPOSE) down

prod-restart:
	$(PROD_COMPOSE) down
	$(PROD_COMPOSE) up -d --build

prod-logs:
	$(PROD_COMPOSE) logs -f

prod-migrate:
	$(PROD_COMPOSE) run --rm migrate

prod-build:
	$(PROD_COMPOSE) build --no-cache backend
