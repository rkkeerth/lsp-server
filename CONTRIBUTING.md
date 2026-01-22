# Contributing to LSP Server

Thank you for your interest in contributing! This document provides guidelines and information for contributing to this project.

## Project Structure

```
lsp-server/
├── main.go              # Entry point, stdio communication
├── server.go            # Core server, request routing
├── handlers.go          # LSP method implementations
├── document.go          # Document state management
├── go.mod               # Go module dependencies
├── Makefile             # Build automation
├── README.md            # Main documentation
├── QUICKSTART.md        # Quick start guide
└── examples/            # Example files and configurations
    ├── test.txt         # Example file for testing
    ├── test.sh          # Manual test script
    └── vscode-extension/ # VS Code extension example
```

## Development Setup

1. **Install Go 1.21+**
   ```bash
   go version  # Should be 1.21 or higher
   ```

2. **Clone and setup**
   ```bash
   cd lsp-server
   go mod download
   ```

3. **Build**
   ```bash
   make build
   ```

## Making Changes

### Adding a New LSP Method

1. **Add the handler in `server.go`:**
   ```go
   case "textDocument/newMethod":
       return s.handleTextDocumentNewMethod(ctx, reply, req)
   ```

2. **Implement the handler:**
   ```go
   func (s *Server) handleTextDocumentNewMethod(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
       var params protocol.NewMethodParams
       if err := json.Unmarshal(req.Params(), &params); err != nil {
           return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
       }
       
       // Implementation here
       
       return reply(ctx, result, nil)
   }
   ```

3. **Update capabilities in `handleInitialize()`** if needed:
   ```go
   NewMethodProvider: &protocol.NewMethodOptions{},
   ```

4. **Test the new method** with your editor

### Code Style Guidelines

- **Formatting**: Use `go fmt` or `gofmt`
- **Naming**: Follow Go naming conventions (camelCase for unexported, PascalCase for exported)
- **Errors**: Always check and handle errors appropriately
- **Logging**: Use `log.Printf` for debugging, include context
- **Comments**: Add godoc comments for exported functions

### Error Handling Pattern

```go
if err != nil {
    log.Printf("Error doing something: %v", err)
    return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InternalError, err.Error()))
}
```

### Thread Safety

When working with shared state:
- Use `sync.RWMutex` for read-heavy operations
- Lock for writes, RLock for reads
- Keep critical sections small
- Document thread-safety guarantees

Example from `document.go`:
```go
func (dm *DocumentManager) Get(uri string) (*Document, bool) {
    dm.mu.RLock()              // Read lock
    defer dm.mu.RUnlock()      // Ensure unlock
    
    doc, exists := dm.documents[uri]
    return doc, exists
}
```

## Testing

### Manual Testing

1. **Build the server:**
   ```bash
   make build
   ```

2. **Run with logging:**
   ```bash
   ./lsp-server 2> server.log
   ```

3. **Send test messages:**
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./lsp-server
   ```

4. **Or use the test script:**
   ```bash
   ./examples/test.sh
   ```

### Editor Testing

Test in a real editor to verify behavior:
1. Configure the server in VS Code/Vim/Emacs (see README.md)
2. Open the `examples/test.txt` file
3. Verify all features work as expected

### Adding Unit Tests

Create `*_test.go` files:

```go
package main

import "testing"

func TestDocumentManager(t *testing.T) {
    dm := NewDocumentManager()
    
    // Test opening
    dm.Open("file:///test.txt", "content", 1)
    
    doc, exists := dm.Get("file:///test.txt")
    if !exists {
        t.Error("Document should exist")
    }
    if doc.Content != "content" {
        t.Errorf("Expected 'content', got '%s'", doc.Content)
    }
}
```

Run tests:
```bash
go test -v ./...
```

## Common Improvements

### Priority Enhancements

1. **Incremental sync**: Change from full text sync to incremental
2. **Diagnostics**: Add error/warning reporting
3. **More LSP methods**: Formatting, references, rename, etc.
4. **Language parsing**: Replace pattern matching with real parsing
5. **Workspace support**: Handle multiple files and projects

### Performance Optimization

- Add caching for parsed syntax trees
- Implement proper indexing for symbol lookup
- Use goroutines for parallel processing
- Profile with `pprof` to find bottlenecks

### Quality Improvements

- Add comprehensive unit tests (target >80% coverage)
- Add integration tests with mock clients
- Improve error messages
- Add structured logging (replace log.Printf)
- Add metrics/telemetry

## Pull Request Process

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes** following the guidelines above

3. **Test thoroughly:**
   - Build succeeds
   - Manual testing works
   - No regressions in existing features

4. **Commit with clear messages:**
   ```bash
   git commit -m "Add support for textDocument/formatting
   
   - Implement formatting handler in server.go
   - Add format function in handlers.go
   - Update server capabilities
   - Add tests"
   ```

5. **Push and create a pull request**

## Code Review Checklist

- [ ] Code follows Go conventions and style
- [ ] All errors are handled appropriately
- [ ] Thread-safety is maintained for concurrent access
- [ ] Logging provides useful debugging information
- [ ] Changes are documented in code comments
- [ ] README.md is updated if needed
- [ ] Manual testing shows the feature works
- [ ] No breaking changes to existing functionality

## Architecture Decisions

When making significant changes, consider:

1. **Backward compatibility**: Can existing clients still work?
2. **Performance**: Will this impact latency or memory usage?
3. **Maintainability**: Is the code clear and well-organized?
4. **Standards compliance**: Does it follow LSP specification?

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [Go LSP Libraries](https://pkg.go.dev/go.lsp.dev)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Questions?

If you have questions about contributing:
1. Check the README.md and existing code
2. Look at the LSP specification
3. Review closed pull requests for similar changes
4. Open an issue for discussion

Happy coding! 🚀
