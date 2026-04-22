.PHONY: run db-up db-down up down logs backend-logs frontend-logs

-include .env
export

DATABASE_URL ?= postgres://packwise:packwise@localhost:5432/pack_wise?sslmode=disable

db-up:
	docker compose up -d db

db-down:
	docker compose stop db

run:
	cd backend && PORT=$(BACKEND_PORT) DATABASE_URL=$(DATABASE_URL) CORS_ALLOW_ORIGIN=$(CORS_ALLOW_ORIGIN) go run .

up:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f backend

backend-logs:
	docker compose logs -f backend

frontend-logs:
	docker compose logs -f frontend
