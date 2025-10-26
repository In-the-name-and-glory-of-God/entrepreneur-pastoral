.PHONY: help build docker-up docker-down docker-logs pipeline go-mod-verify go-vet lint test
.DEFAULT_GOAL := help

# Variables
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")
MIGRATIONS_PATH=./pkg/database/migrations
DATABASE_ADDR="postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?$(DB_PARAMS)&search_path=public"

help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Available commands:"
	@echo "  build          Build the Go binary"
	@echo "  docker-up      Start the services using docker-compose"
	@echo "  docker-down    Stop the services using docker-compose"
	@echo "  docker-logs    View the logs of the services"
	@echo "	 migration		Create the migrations"
	@echo "  migration-up   Run the migrations"
	@echo "  migration-down Rollback the last migration or a specific number of migrations"
	@echo "  pipeline   	Runs go-mod-verify, go-vet, lint, test and build"
	@echo "  go-mod-verify  Verify dependencies"
	@echo "  go-vet        	Analyze source code"
	@echo "  lint           Run the linter"
	@echo "  test           Run the tests"

build:
	@echo "Building the application..."
	@go build -o build/$(APP_NAME) ./cmd/server/main.go

docker-up:
	@echo "Starting the services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping the services..."
	@docker-compose down

docker-logs:
	@echo "Viewing the logs..."
	@docker-compose logs -f

.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATIONS_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migration-up:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DATABASE_ADDR) up

.PHONY: migrate-down
migration-down:
	@migrate -path=$(MIGRATIONS_PATH) -database=$(DATABASE_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

pipeline: go-mod-verify go-vet lint test build
	@echo "Entrepreneur Pastoral pipeline execution done."

go-mod-verify:
	@echo "Verifying modules..."
	@go mod verify

go-vet:
	@echo "Analyzing source code..."
	@go vet ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

test:
	@echo "Running tests..."
	@go test -v ./...
	
