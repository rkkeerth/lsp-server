.PHONY: all build test clean run install fmt vet lint

# Binary name
BINARY_NAME=lsp-server
BINARY_PATH=./$(BINARY_NAME)

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

all: test build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_PATH) -v

test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_PATH)

run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH)

install:
	@echo "Installing $(BINARY_NAME)..."
	$(GOCMD) install

fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

vet:
	@echo "Running go vet..."
	$(GOVET) ./...

lint: fmt vet

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Development helpers
dev-build: fmt vet build

.DEFAULT_GOAL := build
