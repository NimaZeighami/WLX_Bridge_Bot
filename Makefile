# Project variables
BINARY_NAME ?= bridgebot
BUILD_DIR ?= bin
CMD_DIR ?= ./cmd/bridgebot
PKG_LIST := $(shell go list ./...)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS = -X "main.Version=$(VERSION)"

# Tools
GOLANGCI_LINT ?= golangci-lint
GOTEST ?= go test
GOBUILD ?= go build
GOMOD ?= go mod
GOCLEAN ?= go clean

.PHONY: all build run test lint format vet tidy clean install-tools docker-build docker-run

all: build

## Build the Go binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

## Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

## Run all tests with coverage
test:
	@echo "Running tests..."
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...

## Run the Go linter
lint:
	@echo "Linting code..."
	$(GOLANGCI_LINT) run

## Format the code
format:
	@echo "Formatting code..."
	go fmt ./...

## Vet the code
vet:
	@echo "Running go vet..."
	go vet ./...

## Tidy dependencies
tidy:
	@echo "Tidying Go modules..."
	$(GOMOD) tidy

## Clean up generated files
clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR) coverage.out

## Install development tools
install-tools:
	@echo "Installing required tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## Cross-compile for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(CMD_DIR)

## Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):latest .

## Run Docker container
docker-run: docker-build
	@echo "Running Docker container..."
	docker run --rm -p 8080:8080 $(BINARY_NAME):latest
