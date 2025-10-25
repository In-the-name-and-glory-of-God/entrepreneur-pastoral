.PHONY: help build docker-up docker-down docker-logs go-mod-verify go-vet lint test build pipeline
.DEFAULT_GOAL := help

# Variables
APP_NAME=entrepreneur-pastoral
GO_FILES=$(shell find . -name '*.go' -not -path "./vendor/*")

help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Available commands:"
	@echo "  build          Build the Go binary"
	@echo "  docker-up      Start the services using docker-compose"
	@echo "  docker-down    Stop the services using docker-compose"
	@echo "  docker-logs    View the logs of the services"
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
	
