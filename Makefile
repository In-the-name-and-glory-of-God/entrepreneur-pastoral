.PHONY: help build run air docker-build docker-up docker-down docker-logs test lint
.DEFAULT_GOAL := help

# Variables
APP_NAME=entrepreneur-pastoral
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Available commands:"
	@echo "  build          Build the Go binary"
	@echo "  run            Run the application"
	@echo "  air            Run the application with live reloading"
	@echo "  docker-build   Build the Docker image"
	@echo "  docker-up      Start the services using docker-compose"
	@echo "  docker-down    Stop the services using docker-compose"
	@echo "  docker-logs    View the logs of the services"
	@echo "  test           Run the tests"
	@echo "  lint           Run the linter"

build:
	@echo "Building the application..."
	@go build -o build/$(APP_NAME) ./cmd/server/main.go

run:
	@echo "Running the application..."
	@go run ./cmd/server/main.go

docker-build:
	@echo "Building the Docker image..."
	@docker-compose build

docker-up:
	@echo "Starting the services..."
	@docker-compose up -d

docker-down:
	@echo "Stopping the services..."
	@docker-compose down

docker-logs:
	@echo "Viewing the logs..."
	@docker-compose logs -f

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run
