# Makefile

.PHONY: help build run stop clean migrate-up migrate-down test

help:
	@echo "Available commands:"
	@echo "  make build        - Build Docker images"
	@echo "  make run          - Run all services"
	@echo "  make stop         - Stop all services"
	@echo "  make clean        - Stop and remove containers, volumes"
	@echo "  make migrate-up   - Run migrations"
	@echo "  make migrate-down - Rollback migrations"
	@echo "  make logs         - View logs"
	@echo "  make test         - Run tests"
	@echo "  make dev          - Run in development mode"

build:
	docker compose build

run:
	docker compose up -d

stop:
	docker compose down

clean:
	docker compose down -v

migrate-up:
	docker compose exec app /app/migrate -path /app/migrations -database "$$DATABASE_URL" up

migrate-down:
	docker compose exec app /app/migrate -path /app/migrations -database "$$DATABASE_URL" down

logs:
	docker compose logs -f

test:
	docker compose run --rm app go test -v ./...

dev:
	docker compose -f docker compose.dev.yml up

shell:
	docker compose exec app /bin/sh

psql:
	docker compose exec postgres psql -U postgres -d mydb