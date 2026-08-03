SHELL := /bin/bash

BACKEND_DIR := $(CURDIR)/src/backend
FRONTEND_DIR := $(CURDIR)/src/frontend

SERVICE ?=
LOGS_TAIL ?=

.PHONY: help init lint test build up down logs clean migrate

help:
	@echo " "
	@echo "Targets:"
	@echo "  init            - Create .env and install backend and frontend dependencies"
	@echo "  lint            - Lint Go and frontend, check file size limit and locale sync"
	@echo "  test            - Run backend and frontend tests with coverage"
	@echo " "
	@echo "  build           - Build docker images (use SERVICE=... for a single service)"
	@echo "  up              - Start the stack and wait until healthy (use SERVICE=...)"
	@echo "  down            - Stop and remove containers (use SERVICE=...)"
	@echo "  logs            - Follow logs (use SERVICE=... and/or LOGS_TAIL=...)"
	@echo "  clean           - Stop the stack and delete all volumes with data"
	@echo " "
	@echo "  migrate         - Apply pending database migrations"
	@echo " "

init:
	@test -f .env || cp .env.example .env
	cd "$(BACKEND_DIR)" && go mod download
	cd "$(FRONTEND_DIR)" && npm ci --silent
	@echo " "
	@echo "Environment ready. Fill in POSTGRES_PASSWORD and JWT_SECRET in .env"
	@echo " "

lint:
	cd "$(BACKEND_DIR)" && golangci-lint run --config ../../.golangci.yaml ./...
	cd "$(FRONTEND_DIR)" && npm run lint -- --max-warnings=0
	cd "$(FRONTEND_DIR)" && npm run typecheck
	@bash ./scripts/check-file-length.sh
	@node ./scripts/check-locales.mjs
	@echo " "
	@echo "Linting completed!"

test:
	cd "$(BACKEND_DIR)" && go test -race -coverprofile=coverage.out -covermode=atomic ./...
	cd "$(BACKEND_DIR)" && go tool cover -func=coverage.out | tail -1
	cd "$(FRONTEND_DIR)" && npm run test:coverage
	@echo " "
	@echo "Tests completed!"

build:
	@if [ -n "$(SERVICE)" ]; then \
		docker compose build $(SERVICE); \
	else \
		docker compose build; \
	fi
	@echo " "
	@echo "Build completed!"

up:
	@if [ -n "$(SERVICE)" ]; then \
		docker compose up -d --build --wait --remove-orphans $(SERVICE); \
	else \
		docker compose up -d --build --wait --remove-orphans; \
	fi
	@echo " "
	@echo "  frontend -> http://localhost:3000"
	@echo "  backend  -> http://localhost:8080/api/v1"
	@echo "  metrics  -> http://localhost:8080/metrics"
	@echo " "

down:
	@if [ -n "$(SERVICE)" ]; then \
		docker compose stop $(SERVICE) || true; \
		docker compose rm -f $(SERVICE) || true; \
	else \
		docker compose down --remove-orphans; \
	fi
	@echo " "
	@echo "Containers stopped and removed!"

logs:
	@if [ -n "$(SERVICE)" ] && [ -n "$(LOGS_TAIL)" ]; then \
		docker compose logs -f --tail=$(LOGS_TAIL) $(SERVICE); \
	elif [ -n "$(SERVICE)" ]; then \
		docker compose logs -f --tail=100 $(SERVICE); \
	elif [ -n "$(LOGS_TAIL)" ]; then \
		docker compose logs -f --tail=$(LOGS_TAIL); \
	else \
		docker compose logs -f --tail=100; \
	fi

clean:
	docker compose down -v --remove-orphans
	@echo " "
	@echo "Containers and volumes removed!"

migrate:
	docker compose run --rm migrate up
	@echo " "
	@echo "Migrations applied!"
