# load the environment variables
include .env

DB_URI = "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=disable"

# targets
.PHONY: db-up db-down build run migrate-up migrate-down help db-reset

db-up:
	docker compose up -d

db-down:
	docker compose down

build:
	go build -o chronotube cmd/main.go

run: build
	./chronotube

migrate-up:
	cd db/schema && goose postgres $(DB_URI) up && cd ../..

migrate-down:
	cd db/schema && goose postgres $(DB_URI) down && cd ../..

db-reset:
	cd db/schema && goose postgres $(DB_URI) down-to 0 && cd ../..

.DEFAULT_GOAL := help

help:
	@echo "Available targets:"
	@echo "  db-up        - Starts up the database container"
	@echo "  db-down      - Stops and removes the database container"
	@echo "  build        - Builds the Chronotube application"
	@echo "  run          - Runs the Chronotube application"
	@echo "  migrate-up   - Applies database migrations (up)"
	@echo "  migrate-down - Rolls back database migrations (down)"
	@echo "  db-reset     - Resets the database to initial state"
