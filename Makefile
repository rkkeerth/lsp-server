.PHONY: all build clean test run install deps help

# Binary name
BINARY_NAME=lsp-server
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME).exe

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Build the project
all: clean deps build

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	@echo "Building $(BINARY_NAME)..."
	@go build -o $(BINARY_NAME) -v

clean: ## Remove build artifacts
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -f $(BINARY_UNIX)
	@rm -f $(BINARY_WINDOWS)
	@rm -f *.log
	@go clean

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

run: build ## Build and run the server
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

install: ## Install the binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@go install

# Cross compilation
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build -o $(BINARY_UNIX) -v

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build -o $(BINARY_WINDOWS) -v

build-all: build-linux build-windows ## Build for all platforms
	@echo "Building for all platforms..."

# Development helpers
fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golint (requires golint to be installed)
	@echo "Running golint..."
	@golint ./...

check: fmt vet ## Run formatters and linters
	@echo "Code check complete"

# Docker support (if needed)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME):latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run -i $(BINARY_NAME):latest
