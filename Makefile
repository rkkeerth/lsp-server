.PHONY: build clean test run fmt vet

# Build the LSP server
build:
	go build -o lsp-server ./cmd/lsp-server

# Clean build artifacts
clean:
	rm -f lsp-server

# Run tests
test:
	go test -v ./...

# Run the server
run: build
	./lsp-server

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Install dependencies
deps:
	go mod download
	go mod tidy

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o lsp-server-linux-amd64 ./cmd/lsp-server
	GOOS=darwin GOARCH=amd64 go build -o lsp-server-darwin-amd64 ./cmd/lsp-server
	GOOS=windows GOARCH=amd64 go build -o lsp-server-windows-amd64.exe ./cmd/lsp-server
