.PHONY: all build clean install test help

# Binary name
BINARY_NAME=lsp-server

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

all: clean build

help:
	@echo "LSP Server - Makefile commands:"
	@echo "  make build     - Build the LSP server binary"
	@echo "  make clean     - Remove built binaries and clean cache"
	@echo "  make install   - Install the binary to GOPATH/bin"
	@echo "  make test      - Run tests"
	@echo "  make deps      - Download and verify dependencies"
	@echo "  make help      - Show this help message"

deps:
	$(GOMOD) download
	$(GOMOD) verify

build: deps
	$(GOBUILD) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

install: deps
	$(GOINSTALL)

test:
	$(GOTEST) -v ./...

# Cross-compilation targets
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-linux-amd64 -v

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-windows-amd64.exe -v

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)-darwin-amd64 -v

build-all: build-linux build-windows build-darwin
	@echo "Built binaries for all platforms"
