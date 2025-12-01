.PHONY: build test clean run install example

# Build the LSP server
build:
	go build -o lsp-server .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	rm -f lsp-server coverage.out coverage.html
	rm -f test_server.log test_output.log

# Run the server (useful for manual testing)
run: build
	./lsp-server

# Install dependencies
install:
	go mod download
	go mod tidy

# Run the example client
example: build
	python3 examples/simple_client.py ./lsp-server

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	go vet ./...

# Full check: format, lint, and test
check: fmt lint test

# Show help
help:
	@echo "Available targets:"
	@echo "  build          - Build the LSP server binary"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove build artifacts"
	@echo "  run            - Build and run the server"
	@echo "  install        - Download and tidy dependencies"
	@echo "  example        - Run the example Python client"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Run Go linter"
	@echo "  check          - Format, lint, and test"
	@echo "  help           - Show this help message"
