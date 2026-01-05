# Contributing to LSP Server

Thank you for your interest in contributing to the LSP Server project! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Issue Reporting](#issue-reporting)

## Code of Conduct

This project follows a standard code of conduct:

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Maintain professional communication

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/lsp-server.git
   cd lsp-server
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/rkkeerth/lsp-server.git
   ```
4. **Create a branch** for your changes:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- A text editor with Go support (VS Code, GoLand, Neovim, etc.)

### Installation

```bash
# Download dependencies
go mod download

# Build the project
make build

# Run tests
make test

# Run formatters and linters
make check
```

### Running the Server

```bash
# Run directly
./lsp-server

# Run with debug logs
./lsp-server 2> debug.log
```

## Project Structure

```
lsp-server/
├── main.go              # Entry point
├── server/              # Server implementation
│   └── server.go
├── document/            # Document management
│   ├── manager.go
│   └── manager_test.go
├── handlers/            # LSP feature handlers
│   ├── handlers.go
│   ├── symbol_index.go
│   ├── handlers_test.go
│   └── symbol_index_test.go
├── examples/            # Example files for testing
└── .vscode/            # VS Code configuration
```

## Coding Standards

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

### Key Principles

1. **Formatting**: Use `gofmt` (run `make fmt`)
2. **Naming**: 
   - Use camelCase for variables
   - Use PascalCase for exported identifiers
   - Use descriptive names (avoid single letters except in short scopes)
3. **Documentation**:
   - Document all exported functions, types, and constants
   - Use complete sentences starting with the identifier name
   - Example:
     ```go
     // NewManager creates a new document manager instance
     func NewManager() *Manager {
     ```

4. **Error Handling**:
   - Always check and handle errors
   - Provide context when wrapping errors
   - Log errors appropriately

5. **Concurrency**:
   - Protect shared state with mutexes
   - Prefer channels for communication
   - Document goroutine lifecycles

### Code Examples

#### Good Example

```go
// GetWordAt retrieves the word at a specific position in a document.
// It returns an empty string if the position is invalid.
func (d *Document) GetWordAt(pos protocol.Position) string {
    if int(pos.Line) >= len(d.Lines) {
        return ""
    }

    line := d.Lines[pos.Line]
    if int(pos.Character) >= len(line) {
        return ""
    }

    // Find word boundaries
    start := int(pos.Character)
    end := int(pos.Character)

    for start > 0 && isWordChar(rune(line[start-1])) {
        start--
    }

    for end < len(line) && isWordChar(rune(line[end])) {
        end++
    }

    return line[start:end]
}
```

#### Bad Example

```go
// get word
func (d *Document) GetWordAt(p protocol.Position) string {
    l := d.Lines[p.Line]  // No bounds checking
    // ... implementation
}
```

## Testing Guidelines

### Writing Tests

1. **Unit Tests**: Test individual functions and methods
   ```go
   func TestDocumentGetWordAt(t *testing.T) {
       doc := &Document{
           Lines: []string{"hello world"},
       }
       
       word := doc.GetWordAt(protocol.Position{Line: 0, Character: 0})
       if word != "hello" {
           t.Errorf("Expected 'hello', got '%s'", word)
       }
   }
   ```

2. **Table-Driven Tests**: Use for multiple test cases
   ```go
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
               result := ProcessInput(tt.input)
               if result != tt.expected {
                   t.Errorf("Expected %s, got %s", tt.expected, result)
               }
           })
       }
   }
   ```

3. **Test Coverage**: Aim for >80% coverage
   ```bash
   go test -cover ./...
   ```

4. **Race Detection**: Always check for race conditions
   ```bash
   go test -race ./...
   ```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detection
go test -race ./...

# Run specific package
go test ./document/

# Run specific test
go test -run TestManagerOpen ./document/
```

## Pull Request Process

### Before Submitting

1. **Update your branch** with upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run all checks**:
   ```bash
   make check    # Format and lint
   make test     # Run tests
   ```

3. **Update documentation** if needed:
   - Update README.md for user-facing changes
   - Update ARCHITECTURE.md for architectural changes
   - Add/update code comments

4. **Test manually** with an editor:
   - Build the server
   - Test with VS Code, Neovim, or another LSP client
   - Verify your changes work as expected

### Commit Guidelines

Use clear, descriptive commit messages:

```
Add workspace symbol search optimization

- Implement symbol indexing for faster lookups
- Add caching for frequently accessed symbols
- Update tests to cover new functionality

Fixes #123
```

**Format**:
- First line: Brief summary (50 chars or less)
- Blank line
- Detailed description
- Reference to issue if applicable

### Submitting the PR

1. **Push your branch**:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create Pull Request** on GitHub:
   - Provide a clear title and description
   - Reference related issues
   - Include screenshots/examples if applicable
   - Mark as draft if work in progress

3. **PR Description Template**:
   ```markdown
   ## Description
   Brief description of changes
   
   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update
   
   ## Testing
   - [ ] Unit tests added/updated
   - [ ] Manual testing completed
   - [ ] No race conditions
   
   ## Checklist
   - [ ] Code follows style guidelines
   - [ ] Documentation updated
   - [ ] Tests pass
   - [ ] No breaking changes (or documented)
   ```

### Review Process

- Maintainers will review your PR
- Address feedback and comments
- Update your PR based on feedback
- Once approved, a maintainer will merge

## Issue Reporting

### Bug Reports

Use this template:

```markdown
## Description
Clear description of the bug

## Steps to Reproduce
1. Open file X
2. Trigger action Y
3. Observe error Z

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g., macOS 13.0]
- Go Version: [e.g., 1.21.0]
- Editor: [e.g., VS Code 1.85.0]
- LSP Server Version: [e.g., 1.0.0]

## Logs
```
Paste relevant logs here
```

## Additional Context
Any other relevant information
```

### Feature Requests

Use this template:

```markdown
## Feature Description
Clear description of the proposed feature

## Use Case
Why is this feature needed? What problem does it solve?

## Proposed Implementation
Ideas for how this could be implemented

## Alternatives Considered
Other approaches you've thought about

## Additional Context
Any other relevant information
```

## Development Tips

### Debugging

1. **Structured Logging**: Use zap for debugging
   ```go
   h.logger.Debug("Processing request",
       zap.String("uri", string(uri)),
       zap.Int("line", int(pos.Line)),
   )
   ```

2. **LSP Inspector**: Use tools like:
   - VS Code: Enable LSP logging
   - Neovim: `:lua vim.lsp.set_log_level("debug")`

3. **JSON-RPC Inspection**: Monitor stdin/stdout with a proxy

### Common Pitfalls

1. **Writing to stdout**: Never use `fmt.Println()` - corrupts protocol
2. **Ignoring errors**: Always handle errors properly
3. **Race conditions**: Always test with `-race` flag
4. **Unbounded operations**: Consider performance with large files

### Useful Commands

```bash
# Format all code
make fmt

# Run linters
make vet

# Build for multiple platforms
make build-all

# Install locally
make install

# Clean build artifacts
make clean
```

## Getting Help

- **Documentation**: Read README.md and ARCHITECTURE.md
- **Examples**: Check the examples/ directory
- **Issues**: Search existing issues for similar problems
- **Discussions**: Start a discussion for questions

## Recognition

Contributors will be recognized in:
- GitHub contributors list
- Release notes for significant contributions
- Project documentation

Thank you for contributing! 🎉
