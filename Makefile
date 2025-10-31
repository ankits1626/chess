.PHONY: up down logs backend-logs restart clean migrate-up migrate-down migrate-reset db-shell test help

# Start all services
up:
	docker compose up --build -d

# Stop all services
down:
	docker compose down

# View logs for all services
logs:
	docker compose logs -f

# View logs for backend service only
backend-logs:
	docker compose logs -f backend

# Restart all services
restart:
	docker compose restart

# Stop all services and remove volumes
clean:
	docker compose down -v

# Run database migrations
migrate-up:
	@echo "Running database migrations..."
	docker exec -i chess-coach-postgres psql -U postgres -d chess_coach < backend/internal/db/migrations/001_initial_schema.sql
	@echo "Migrations completed!"

# Drop all tables (destructive!)
migrate-down:
	@echo "WARNING: This will drop all tables!"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		docker exec -i chess-coach-postgres psql -U postgres -d chess_coach -c "DROP TABLE IF EXISTS moves CASCADE;"; \
		docker exec -i chess-coach-postgres psql -U postgres -d chess_coach -c "DROP TABLE IF EXISTS games CASCADE;"; \
		docker exec -i chess-coach-postgres psql -U postgres -d chess_coach -c "DROP TABLE IF EXISTS users CASCADE;"; \
		docker exec -i chess-coach-postgres psql -U postgres -d chess_coach -c "DROP EXTENSION IF EXISTS \"uuid-ossp\";"; \
		echo "All tables dropped!"; \
	fi

# Reset database (drop + recreate)
migrate-reset: migrate-down migrate-up

# Open PostgreSQL shell
db-shell:
	docker exec -it chess-coach-postgres psql -U postgres -d chess_coach

# Run tests
test:
	cd backend && go test ./... -v

# Show help
help:
	@echo "Available commands:"
	@echo "  make up            - Start all services"
	@echo "  make down          - Stop all services"
	@echo "  make logs          - View logs for all services"
	@echo "  make backend-logs  - View backend logs only"
	@echo "  make restart       - Restart all services"
	@echo "  make clean         - Stop services and remove volumes"
	@echo "  make migrate-up    - Run database migrations"
	@echo "  make migrate-down  - Drop all database tables"
	@echo "  make migrate-reset - Reset database (drop + migrate)"
	@echo "  make db-shell      - Open PostgreSQL shell"
	@echo "  make test          - Run all Go tests"
	@echo "  make help          - Show this help message"
