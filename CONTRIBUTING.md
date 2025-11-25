# Contributing to LSP Server

Thank you for your interest in contributing to this LSP server implementation! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)

### Setting Up Development Environment

1. Fork the repository
2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/lsp-server.git
cd lsp-server
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

### 1. Create a Branch

Create a feature branch for your work:

```bash
git checkout -b feature/your-feature-name
```

Use descriptive branch names:
- `feature/add-hover-support`
- `fix/initialization-bug`
- `docs/improve-readme`

### 2. Make Changes

Follow these guidelines:

- Write clear, concise commit messages
- Keep commits focused and atomic
- Add tests for new features
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Format code
make fmt

# Run linters
make vet

# Test with the example client
cd examples
go run test_client.go ../lsp-server
```

### 4. Submit a Pull Request

1. Push your branch to your fork
2. Open a pull request against the main repository
3. Provide a clear description of your changes
4. Reference any related issues

## Code Style

### Go Code Standards

- Follow standard Go conventions and idioms
- Use `gofmt` for formatting (run `make fmt`)
- Run `go vet` to catch common issues (run `make vet`)
- Write clear, self-documenting code
- Add comments for exported functions and types

### Example:

```go
// HandleRequest processes an incoming LSP request and returns a response.
// It validates the request format and routes to the appropriate handler.
func (s *Server) HandleRequest(req Request) (*Response, error) {
    // Implementation
}
```

## Project Structure

Understanding the project structure helps with contributions:

```
lsp-server/
├── main.go                    # Entry point
├── internal/lsp/             # LSP implementation (not exported)
│   ├── protocol.go           # LSP protocol types
│   ├── server.go             # Core server logic
│   ├── lifecycle.go          # Initialization/shutdown
│   ├── textdocument.go       # Document synchronization
│   └── *_test.go             # Unit tests
├── examples/                 # Example clients and tests
└── docs/                     # Additional documentation
```

## Adding New Features

### Adding a New LSP Capability

1. **Define Protocol Types** (`internal/lsp/protocol.go`):
```go
// HoverParams defines parameters for textDocument/hover
type HoverParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Position     Position                `json:"position"`
}

// Hover represents hover information
type Hover struct {
    Contents string `json:"contents"`
    Range    *Range `json:"range,omitempty"`
}
```

2. **Update Server Capabilities** (`internal/lsp/lifecycle.go`):
```go
result := InitializeResult{
    Capabilities: ServerCapabilities{
        // ... existing capabilities
        HoverProvider: true, // Add this
    },
}
```

3. **Implement Handler** (new file or existing file):
```go
// handleTextDocumentHover handles textDocument/hover requests
func (s *Server) handleTextDocumentHover(request Request) {
    // Implementation
}
```

4. **Register Handler** (`internal/lsp/server.go`):
```go
case "textDocument/hover":
    s.handleTextDocumentHover(request)
```

5. **Add Tests**:
```go
func TestHandleHover(t *testing.T) {
    // Test implementation
}
```

6. **Update Documentation**:
- Add method to README.md supported methods table
- Update capability list
- Add example usage if needed

## Testing Guidelines

### Unit Tests

- Place tests in `*_test.go` files
- Test both success and error cases
- Use table-driven tests for multiple scenarios
- Mock external dependencies

Example:
```go
func TestDocumentManagement(t *testing.T) {
    tests := []struct {
        name     string
        uri      string
        content  string
        expected bool
    }{
        {"valid document", "file:///test.txt", "content", true},
        {"empty uri", "", "content", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

- Use the test client in `examples/test_client.go`
- Test complete message flows
- Verify state changes
- Check error handling

## Documentation

### Code Documentation

- Comment all exported types and functions
- Use godoc conventions
- Include examples where helpful

### User Documentation

- Update README.md for user-facing changes
- Add to TESTING.md for new test procedures
- Update examples for new features

## Pull Request Process

1. **Before Submitting**:
   - Ensure all tests pass
   - Run code formatters and linters
   - Update documentation
   - Add tests for new features

2. **PR Description**:
   - Clearly describe what changes you've made
   - Explain why the change is needed
   - Link to related issues
   - Include testing steps if applicable

3. **Review Process**:
   - Address reviewer feedback promptly
   - Keep discussions focused and professional
   - Be open to suggestions and alternatives

4. **After Approval**:
   - Squash commits if requested
   - Ensure branch is up to date with main
   - Maintainers will merge when ready

## Issue Guidelines

### Reporting Bugs

Include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Go version and OS
- Relevant logs or error messages

Example:
```markdown
**Description**: Server crashes on didChange notification

**Steps to Reproduce**:
1. Initialize server
2. Open document
3. Send didChange with empty contentChanges

**Expected**: Server handles gracefully
**Actual**: Server crashes with panic

**Environment**:
- Go version: 1.21
- OS: Ubuntu 22.04
- LSP Server version: 0.1.0

**Logs**:
```
panic: index out of range
...
```
```

### Suggesting Features

Include:
- Clear description of the feature
- Use cases and benefits
- Possible implementation approach
- Related LSP specification sections

## Communication

- Be respectful and professional
- Ask questions if something is unclear
- Help others when you can
- Stay on topic in discussions

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

If you have questions about contributing:
- Open an issue for discussion
- Check existing issues and PRs
- Review the LSP specification

## Recognition

Contributors will be recognized in:
- Git commit history
- Release notes for significant contributions
- Future CONTRIBUTORS.md file

Thank you for contributing to the LSP Server project!
