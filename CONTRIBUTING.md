# Contributing to LSP Server

Thank you for your interest in contributing to the LSP Server project!

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Basic understanding of the Language Server Protocol
- Familiarity with JSON-RPC 2.0

### Setup Development Environment

1. Clone the repository:
```bash
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server
```

2. Install dependencies:
```bash
make deps
```

3. Build the project:
```bash
make build
```

4. Run tests:
```bash
make test
```

## Development Workflow

### Project Structure

```
lsp-server/
├── cmd/lsp-server/       # Main application entry point
├── internal/
│   ├── jsonrpc/          # JSON-RPC 2.0 implementation
│   ├── protocol/         # LSP protocol types
│   └── server/           # LSP server logic
├── examples/             # Example clients and usage
├── Makefile              # Build automation
└── go.mod                # Go module definition
```

### Making Changes

1. **Create a branch** for your feature or fix:
```bash
git checkout -b feature/your-feature-name
```

2. **Make your changes** following the code style guidelines

3. **Write tests** for your changes:
   - Add unit tests in `*_test.go` files
   - Ensure tests are comprehensive and cover edge cases

4. **Run tests and formatting**:
```bash
make fmt
make vet
make test
```

5. **Commit your changes** with clear commit messages:
```bash
git commit -m "Add feature: description of your changes"
```

## Code Style Guidelines

### General Principles

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (run `make fmt`)
- Run `go vet` to catch common mistakes (run `make vet`)
- Write clear, self-documenting code
- Add comments for exported functions and types

### Naming Conventions

- Use camelCase for local variables: `documentURI`
- Use PascalCase for exported names: `InitializeParams`
- Use descriptive names: prefer `textDocument` over `td`
- Constants in PascalCase: `JSONRPCVersion`

### Error Handling

```go
// Good: Wrap errors with context
if err != nil {
    return nil, fmt.Errorf("failed to parse document: %w", err)
}

// Bad: Return raw errors
if err != nil {
    return nil, err
}
```

### Testing

```go
// Good: Descriptive test names
func TestServerHandlesInitializeRequest(t *testing.T) {
    // Test implementation
}

// Good: Use table-driven tests for multiple cases
func TestMultipleCases(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"case1", "input1", "output1"},
        {"case2", "input2", "output2"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test logic
        })
    }
}
```

## Adding New LSP Features

### 1. Protocol Types

Add new types to `internal/protocol/types.go`:

```go
// HoverParams represents parameters for textDocument/hover
type HoverParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Position     Position                `json:"position"`
}

// Hover represents hover information
type Hover struct {
    Contents MarkupContent `json:"contents"`
    Range    *Range        `json:"range,omitempty"`
}
```

### 2. Server Handler

Add handler to `internal/server/server.go`:

```go
// handleHover handles the textDocument/hover request
func (s *Server) handleHover(params json.RawMessage) (interface{}, error) {
    var hoverParams protocol.HoverParams
    if err := json.Unmarshal(params, &hoverParams); err != nil {
        return nil, fmt.Errorf("invalid hover params: %w", err)
    }
    
    // Implementation
    
    return result, nil
}
```

### 3. Register Method

Update the `Handle()` method:

```go
func (s *Server) Handle(method string, params json.RawMessage) (interface{}, error) {
    switch method {
    // ... existing cases ...
    case "textDocument/hover":
        return s.handleHover(params)
    // ...
    }
}
```

### 4. Update Capabilities

In `handleInitialize()`:

```go
capabilities := protocol.ServerCapabilities{
    // ... existing capabilities ...
    HoverProvider: true,
}
```

### 5. Write Tests

Add tests in `internal/server/server_test.go`:

```go
func TestServerHover(t *testing.T) {
    server := NewServer()
    
    // Setup
    // ...
    
    // Test hover
    hoverParams := protocol.HoverParams{
        TextDocument: protocol.TextDocumentIdentifier{
            URI: "file:///test.txt",
        },
        Position: protocol.Position{Line: 0, Character: 0},
    }
    
    paramsJSON, err := json.Marshal(hoverParams)
    if err != nil {
        t.Fatalf("Failed to marshal params: %v", err)
    }
    
    result, err := server.handleHover(paramsJSON)
    if err != nil {
        t.Fatalf("handleHover failed: %v", err)
    }
    
    // Assertions
    // ...
}
```

## Testing Guidelines

### Running Tests

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/server/

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Good Tests

1. **Test one thing per test**: Each test should verify one behavior
2. **Use descriptive names**: Test names should describe what's being tested
3. **Arrange-Act-Assert**: Structure tests clearly
4. **Test edge cases**: Include boundary conditions and error cases
5. **Use test fixtures**: Create helper functions for common setup

### Example Test Structure

```go
func TestFeature(t *testing.T) {
    // Arrange: Set up test data and state
    server := NewServer()
    input := "test input"
    
    // Act: Execute the functionality
    result, err := server.ProcessInput(input)
    
    // Assert: Verify the results
    if err != nil {
        t.Fatalf("Unexpected error: %v", err)
    }
    
    if result != "expected output" {
        t.Errorf("Expected 'expected output', got '%s'", result)
    }
}
```

## Documentation

### Code Documentation

- Add godoc comments for all exported types and functions
- Use complete sentences
- Include examples where helpful

```go
// InitializeParams represents the parameters sent from the client
// to the server during the initialize request. It includes client
// capabilities and workspace information.
type InitializeParams struct {
    // ProcessID is the process ID of the client process
    ProcessID *int `json:"processId"`
    // ...
}
```

### Architecture Documentation

When adding significant features:
1. Update `ARCHITECTURE.md` with design decisions
2. Update `USAGE.md` with usage examples
3. Add examples to `examples/` directory if applicable

## Pull Request Process

1. **Update documentation** for any changed functionality
2. **Add tests** that verify your changes work correctly
3. **Ensure all tests pass**: `make test`
4. **Format your code**: `make fmt`
5. **Run go vet**: `make vet`
6. **Update README.md** if adding new features

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe the tests you added/ran

## Checklist
- [ ] Code follows project style guidelines
- [ ] Added/updated tests
- [ ] All tests pass
- [ ] Updated documentation
- [ ] No breaking changes (or documented if necessary)
```

## Code Review

All submissions require review. We use GitHub pull requests for this purpose.

### Review Criteria

- Code quality and readability
- Test coverage and quality
- Documentation completeness
- Adherence to LSP specification
- Performance considerations
- Security implications

## Reporting Issues

### Bug Reports

Include:
- Go version: `go version`
- Operating system
- Steps to reproduce
- Expected behavior
- Actual behavior
- Relevant logs or error messages

### Feature Requests

Include:
- Use case description
- Proposed solution
- Alternative solutions considered
- Impact on existing functionality

## Communication

- **Issues**: For bug reports and feature requests
- **Pull Requests**: For code contributions
- **Discussions**: For questions and general discussion

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue for any questions about contributing!

Thank you for contributing to LSP Server! 🎉
