# Makefile — convenience shortcuts for common developer tasks
# Run `make help` to see all available commands.

.PHONY: help run build test sqlc migrate-up migrate-down docker-up docker-down lint

# ── Defaults ───────────────────────────────────────────────────────────────────
BINARY      := ./bin/server
MAIN        := ./cmd/server/main.go
MIGRATE_URL ?= $(shell grep DATABASE_URL .env | cut -d '=' -f2-)
MIGRATIONS  := ./db/migrations

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ── Development ────────────────────────────────────────────────────────────────
run: ## Run the server directly with go run (hot-path for development)
	go run $(MAIN)

build: ## Compile the binary to ./bin/server
	@mkdir -p bin
	go build -o $(BINARY) $(MAIN)

test: ## Run all unit tests with verbose output
	go test ./... -v -race -count=1

# ── Code Generation ────────────────────────────────────────────────────────────
sqlc: ## Re-generate the SQLC database layer from SQL queries
	sqlc generate

# ── Database Migrations ────────────────────────────────────────────────────────
migrate-up: ## Apply all pending migrations
	migrate -path $(MIGRATIONS) -database "$(MIGRATE_URL)" -verbose up

migrate-down: ## Roll back the last applied migration
	migrate -path $(MIGRATIONS) -database "$(MIGRATE_URL)" -verbose down 1

migrate-drop: ## Drop everything — DESTRUCTIVE, development only!
	migrate -path $(MIGRATIONS) -database "$(MIGRATE_URL)" drop -f

# ── Docker ─────────────────────────────────────────────────────────────────────
docker-up: ## Start the full stack (Postgres + API) with Docker Compose
	docker compose up --build -d

docker-down: ## Stop and remove containers
	docker compose down

docker-logs: ## Tail live logs from all containers
	docker compose logs -f

# ── Lint ───────────────────────────────────────────────────────────────────────
lint: ## Run golangci-lint (install it first: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# ── Tidy ───────────────────────────────────────────────────────────────────────
tidy: ## Download and tidy Go module dependencies
	go mod tidy
