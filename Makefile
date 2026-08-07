SHELL := /bin/bash

BACKEND_DIR := $(CURDIR)/src/backend
FRONTEND_DIR := $(CURDIR)/src/frontend
INTEGRATION_DIR := $(CURDIR)/tests/integration

SERVICE ?=
LOGS_TAIL ?=

PROD_COMPOSE := -f docker-compose.yml -f docker-compose.prod.yml

.PHONY: help init lint test test-integration api-spec api-generate build up dev dev-down down logs clean migrate \
        migrate-down migrate-status seed psql prod-up prod-down prod-logs prod-cert

help:
	@echo " "
	@echo "Targets:"
	@echo "  init            - Create .env and install backend and frontend dependencies"
	@echo "  lint            - Lint Go and frontend, check file size limit and locale sync"
	@echo "  test            - Run backend and frontend tests with coverage"
	@echo "  test-integration- Run end-to-end tests against a running stack (needs 'make up')"
	@echo "  api-spec        - Regenerate docs/openapi.yaml from the swaggo annotations in Go"
	@echo "  api-generate    - Regenerate the frontend API client from docs/openapi.yaml"
	@echo " "
	@echo "  build           - Build docker images (use SERVICE=... for a single service)"
	@echo "  up              - Start the stack and wait until healthy (use SERVICE=...)"
	@echo "  dev             - Start the stack with hot reload for backend and frontend"
	@echo "  dev-down        - Stop the hot reload stack"
	@echo "  down            - Stop and remove containers (use SERVICE=...)"
	@echo "  logs            - Follow logs (use SERVICE=... and/or LOGS_TAIL=...)"
	@echo "  clean           - Stop the stack and delete all volumes with data"
	@echo " "
	@echo "  migrate         - Apply pending database migrations"
	@echo "  migrate-down    - Roll back the last migration"
	@echo "  migrate-status  - Show which migrations are applied"
	@echo "  seed            - Reapply the demo data seed (00010)"
	@echo "  psql            - Open a psql shell in the database container"
	@echo " "
	@echo "  prod-up         - Start the production stack behind nginx (80/443)"
	@echo "  prod-down       - Stop the production stack"
	@echo "  prod-logs       - Follow production logs (use SERVICE=... and/or LOGS_TAIL=...)"
	@echo "  prod-cert       - Issue a Let's Encrypt certificate (DOMAIN=... LETSENCRYPT_EMAIL=...)"
	@echo " "

init:
	@test -f .env || cp .env.example .env
	cd "$(BACKEND_DIR)" && go mod download
	cd "$(FRONTEND_DIR)" && npm install --silent --no-audit --no-fund
	@echo " "
	@echo "Environment ready. .env holds working local defaults;"
	@echo "replace the secrets before exposing the stack outside localhost."
	@echo " "

lint:
	cd "$(BACKEND_DIR)" && golangci-lint run --config ../../.golangci.yaml ./...
	cd "$(INTEGRATION_DIR)" && go vet -tags=integration ./...
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

# Integration tests skip themselves unless the stack is up and INTEGRATION=1 is set.
test-integration:
	cd "$(INTEGRATION_DIR)" && INTEGRATION=1 go test -tags=integration -count=1 -v ./...
	@echo " "
	@echo "Integration tests completed!"

api-spec:
	@bash ./scripts/generate-openapi.sh
	@echo " "
	@echo "docs/openapi.yaml regenerated from the Go annotations"
	@echo "Run 'make api-generate' to refresh the frontend client too."
	@echo " "

api-generate:
	cd "$(FRONTEND_DIR)" && npm run api:generate
	@echo " "
	@echo "API client regenerated from docs/openapi.yaml"

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

dev:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build --wait --remove-orphans $(SERVICE)
	@echo " "
	@echo "  Hot reload is on: backend rebuilds on .go changes, frontend uses Vite HMR"
	@echo "  frontend -> http://localhost:3000"
	@echo "  backend  -> http://localhost:8080/api/v1"
	@echo " "

dev-down:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down --remove-orphans
	@echo " "
	@echo "Dev stack stopped!"

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

migrate-down:
	docker compose run --rm --entrypoint /app/migrate migrate down
	@echo " "
	@echo "Last migration rolled back!"

migrate-status:
	docker compose run --rm --entrypoint /app/migrate migrate status

seed:
	docker compose run --rm --entrypoint /app/migrate migrate down
	docker compose run --rm migrate up
	@echo " "
	@echo "Demo data reseeded. Every demo account uses the password: demo1234"
	@echo " "

psql:
	docker compose exec postgres psql -U "$${POSTGRES_USER:-avito}" -d "$${POSTGRES_DB:-avito}"

prod-up:
	docker compose $(PROD_COMPOSE) up -d --build --remove-orphans
	@echo " "
	@echo "  Production stack is up behind nginx on ports 80 and 443"
	@echo "  Next steps for HTTPS are in deploy/README.md"
	@echo " "

prod-down:
	docker compose $(PROD_COMPOSE) down --remove-orphans
	@echo " "
	@echo "Production stack stopped!"

prod-logs:
	@if [ -n "$(SERVICE)" ]; then \
		docker compose $(PROD_COMPOSE) logs -f --tail=$${LOGS_TAIL:-100} $(SERVICE); \
	else \
		docker compose $(PROD_COMPOSE) logs -f --tail=$${LOGS_TAIL:-100}; \
	fi

prod-cert:
	@test -n "$(DOMAIN)" || { echo "DOMAIN is required: make prod-cert DOMAIN=example.com LETSENCRYPT_EMAIL=you@example.com"; exit 1; }
	@test -n "$(LETSENCRYPT_EMAIL)" || { echo "LETSENCRYPT_EMAIL is required"; exit 1; }
	DOMAIN=$(DOMAIN) LETSENCRYPT_EMAIL=$(LETSENCRYPT_EMAIL) \
		docker compose $(PROD_COMPOSE) --profile certbot run --rm certbot
	@echo " "
	@echo "Certificate issued. Enable the HTTPS server block in deploy/nginx/conf.d/app.conf"
	@echo " "
