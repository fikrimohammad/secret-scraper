APP_NAME = secret-scraper
BUILD_DIR = bin
MAIN = ./cmd/main.go

.PHONY: all build run dev clean test lint tidy vendor

all: build

## Build the binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN)
	@echo "Binary: $(BUILD_DIR)/$(APP_NAME)"

## Run without building
run:
	go run $(MAIN)

## Build and run
dev: build
	./$(BUILD_DIR)/$(APP_NAME)

## Run tests
test:
	go test ./... -v -race -count=1

## Run tests with coverage
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

## Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

## Tidy go modules
tidy:
	go mod tidy

## Update vendor directory
vendor:
	go mod vendor

## Remove build artifacts
clean:
	@rm -rf $(BUILD_DIR) coverage.out coverage.html
	@echo "Cleaned."

## Show help
help:
	@echo "Available targets:"
	@grep -E '^##' Makefile | sed 's/## /  /'