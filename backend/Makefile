.PHONY: help install-tools sqlc-generate migrate-up migrate-down docker-up docker-down logs test

# Detect Docker Compose command (V1 vs V2)
DOCKER_COMPOSE := $(shell command -v docker-compose 2> /dev/null)
ifndef DOCKER_COMPOSE
	DOCKER_COMPOSE := docker compose
endif

help:
	@echo "Chess Coach Backend - Commands"
	@echo ""
	@echo "Development:"
	@echo "  make install-tools   - Install sqlc, migrate, etc."
	@echo "  make sqlc-generate   - Generate sqlc code"
	@echo "  make migrate-up      - Run migrations"
	@echo "  make migrate-down    - Rollback migrations"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up       - Start all containers"
	@echo "  make docker-down     - Stop all containers"
	@echo "  make docker-rebuild  - Rebuild and restart"
	@echo "  make logs            - View logs"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run tests"
	@echo "  make test-coverage   - Run tests with coverage"

install-tools:
	@echo "Installing Go tools..."
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed!"

sqlc-generate:
	sqlc generate

migrate-up:
	migrate -path db/migrations -database "postgres://chess_coach:chess_coach_dev@localhost:5432/chess_coach_dev?sslmode=disable" up

migrate-down:
	migrate -path db/migrations -database "postgres://chess_coach:chess_coach_dev@localhost:5432/chess_coach_dev?sslmode=disable" down

docker-up:
	$(DOCKER_COMPOSE) up -d

docker-down:
	$(DOCKER_COMPOSE) down

docker-rebuild:
	$(DOCKER_COMPOSE) up -d --build

logs:
	$(DOCKER_COMPOSE) logs -f api

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out