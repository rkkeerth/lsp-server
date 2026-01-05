#!/bin/bash

# Build script for LSP server

set -e

echo "Building LSP server..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | awk '{print $3}')
echo "Using $GO_VERSION"

# Download dependencies
echo "Downloading dependencies..."
go mod download

# Run tests if requested
if [ "$1" = "--test" ]; then
    echo "Running tests..."
    go test -v ./...
fi

# Build the binary
echo "Building binary..."
go build -o lsp-server .

echo "Build complete! Binary: ./lsp-server"
echo ""
echo "To run the server:"
echo "  ./lsp-server"
echo ""
echo "To see debug logs:"
echo "  ./lsp-server 2> lsp-server.log"
