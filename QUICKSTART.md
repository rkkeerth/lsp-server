# LSP Server - Quick Reference

## Build & Run

```bash
# Build
go build -o lsp-server

# Run (expects JSON-RPC on stdin)
./lsp-server

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o lsp-server-linux
GOOS=darwin GOARCH=amd64 go build -o lsp-server-darwin
GOOS=windows GOARCH=amd64 go build -o lsp-server.exe
```

## Project Files

| File | Purpose |
|------|---------|
| `main.go` | Entry point, stdio setup, connection management |
| `server.go` | Request routing, lifecycle handlers |
| `handlers.go` | LSP feature implementations |
| `document.go` | Thread-safe document state management |
| `go.mod` | Module and dependency definitions |

## Implemented LSP Methods

### Lifecycle
- `initialize` → Returns server capabilities
- `initialized` → Initialization complete notification
- `shutdown` → Prepare for graceful shutdown
- `exit` → Terminate server process

### Text Synchronization
- `textDocument/didOpen` → Document opened
- `textDocument/didChange` → Document modified
- `textDocument/didClose` → Document closed

### Language Features
- `textDocument/hover` → Show hover information
- `textDocument/completion` → Provide completions

## Architecture

```
Editor (VS Code/Neovim/etc)
    ↕ JSON-RPC 2.0 over stdio
main.go (Connection Setup)
    ↓
server.go (Request Router)
    ↓
    ├→ handlers.go (Feature Implementations)
    └→ document.go (State Management)
```

## Adding New Features

1. **Add handler** in `handlers.go`:
```go
func (s *Server) handleNewFeature(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
    var params protocol.NewFeatureParams
    if err := json.Unmarshal(req.Params(), &params); err != nil {
        return reply(ctx, nil, err)
    }
    // Implementation
    return reply(ctx, result, nil)
}
```

2. **Register handler** in `server.go` Handle method:
```go
case protocol.MethodNewFeature:
    return s.handleNewFeature(ctx, reply, req)
```

3. **Update capabilities** in `handleInitialize`:
```go
NewFeatureProvider: &protocol.NewFeatureOptions{...},
```

## Testing

```bash
# Manual test with JSON-RPC message
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | ./lsp-server

# Integration test with editor - see EDITOR_SETUP.md
```

## Debugging

Server logs to stderr. Check:
- **VS Code**: Output panel → "LSP Server"
- **Neovim**: `:LspLog` or `~/.local/state/nvim/lsp.log`
- **Emacs**: `*lsp-log*` buffer (when `lsp-log-io` is enabled)

## Common Tasks

### Enable verbose logging in client

**Neovim:**
```lua
vim.lsp.set_log_level("debug")
```

**VS Code:**
```json
"lspServer.trace.server": "verbose"
```

### Check server is running

```bash
# Should show process if integrated with editor
ps aux | grep lsp-server
```

### Restart server

**Neovim:**
```vim
:LspRestart
```

**VS Code:**
Command Palette → "Reload Window"

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [go.lsp.dev Documentation](https://pkg.go.dev/go.lsp.dev)
- [README.md](./README.md) - Full documentation
- [EDITOR_SETUP.md](./EDITOR_SETUP.md) - Editor configurations
