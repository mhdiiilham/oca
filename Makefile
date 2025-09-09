# Makefile for oca project

MODULE_NAME := oca

.PHONY: all test lint tidy

all: test

# Run unit tests with race detector and coverage
test:
	@echo "Running tests..."
	go test ./... -race -cover

# Run go vet and lint checks
lint:
	@echo "Running go vet..."
	go vet ./...

# Ensure go.mod and go.sum are tidy
tidy:
	@echo "Tidying modules..."
	go mod tidy
