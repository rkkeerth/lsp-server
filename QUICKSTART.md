# LSP Server - Quick Start Guide

## 🚀 Build in 30 Seconds

```bash
# Clone (if needed) and enter directory
cd /projects/sandbox/lsp-server

# Build the server
make build

# Verify build
./lsp-server --help  # (Note: Currently runs as daemon, no --help flag yet)

# Run tests
make test
```

## 📦 What You Get

A fully functional LSP server that:
- Implements JSON-RPC 2.0 protocol
- Handles LSP initialization lifecycle
- Manages text document synchronization
- Ready for editor integration

## 🔧 Basic Usage

### Start the Server

```bash
./lsp-server
```

The server communicates via stdin/stdout using the LSP protocol.

### Test with Example Client

```bash
# Build first
make build

# Run the test client
cd examples
go run test_client.go
```

Expected output:
```
Initialize response:
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "capabilities": {
      "textDocumentSync": {...}
    },
    "serverInfo": {
      "name": "basic-lsp-server",
      "version": "0.1.0"
    }
  }
}

Sent didOpen notification
Sent didChange notification

Shutdown response:
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": null
}

Sent exit notification
Test completed successfully!
```

## 📋 Makefile Commands

```bash
make build       # Build the LSP server binary
make test        # Run all tests
make clean       # Remove build artifacts
make fmt         # Format code with gofmt
make vet         # Run go vet
make deps        # Download dependencies
make build-all   # Build for multiple platforms
```

## 🔌 Editor Integration

### VS Code (Quick Setup)

1. Install the language server extension template
2. Update `serverOptions`:
   ```typescript
   const serverOptions: ServerOptions = {
       command: '/path/to/lsp-server',
       args: [],
       transport: TransportKind.stdio
   };
   ```

### Neovim (Quick Setup)

Add to your `init.lua`:
```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

configs.basic_lsp = {
  default_config = {
    cmd = {'/path/to/lsp-server'},
    filetypes = {'text'},
    root_dir = lspconfig.util.find_git_ancestor,
  },
}

lspconfig.basic_lsp.setup{}
```

### Emacs (Quick Setup)

Add to your Emacs config:
```elisp
(require 'lsp-mode)

(lsp-register-client
 (make-lsp-client
  :new-connection (lsp-stdio-connection "/path/to/lsp-server")
  :major-modes '(text-mode)
  :server-id 'basic-lsp))

(add-hook 'text-mode-hook #'lsp)
```

## 📝 Supported LSP Methods

| Method | Type | Description |
|--------|------|-------------|
| `initialize` | Request | Initialize server |
| `initialized` | Notification | Ready signal |
| `shutdown` | Request | Graceful shutdown |
| `exit` | Notification | Exit process |
| `textDocument/didOpen` | Notification | Document opened |
| `textDocument/didChange` | Notification | Document changed |
| `textDocument/didClose` | Notification | Document closed |

## 🧪 Testing

### Run Unit Tests
```bash
make test
```

### Run Specific Test
```bash
go test -v ./internal/server -run TestServerInitialize
```

### Check Coverage
```bash
go test -cover ./...
```

## 📚 Documentation

- `README.md` - Overview and features
- `ARCHITECTURE.md` - Technical architecture (302 lines)
- `USAGE.md` - Comprehensive usage guide (415 lines)
- `CONTRIBUTING.md` - Contribution guidelines (387 lines)
- `PROJECT_SUMMARY.md` - Complete project summary

## 🛠️ Development

### Project Structure
```
lsp-server/
├── cmd/lsp-server/      # Main entry point
├── internal/
│   ├── jsonrpc/         # JSON-RPC 2.0 layer
│   ├── protocol/        # LSP protocol types
│   └── server/          # Server implementation
└── examples/            # Example client
```

### Adding New LSP Methods

1. Define types in `internal/protocol/types.go`
2. Add handler in `internal/server/server.go`
3. Register in `Handle()` switch statement
4. Update capabilities in `handleInitialize()`
5. Add tests in `internal/server/server_test.go`

Example:
```go
// 1. Protocol type
type HoverParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Position     Position                `json:"position"`
}

// 2. Handler
func (s *Server) handleHover(params json.RawMessage) (interface{}, error) {
    var p HoverParams
    json.Unmarshal(params, &p)
    // Implementation...
    return result, nil
}

// 3. Register
case "textDocument/hover":
    return s.handleHover(params)

// 4. Update capabilities
capabilities.HoverProvider = true
```

## 🐛 Debugging

### Enable Verbose Logging
```bash
./lsp-server 2>/tmp/lsp-debug.log
```

### Trace Messages
Set environment variable:
```bash
export LSP_TRACE=verbose
./lsp-server
```

### Manual Testing
```bash
# Send initialize request
echo -e 'Content-Length: 90\r\n\r\n{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":null,"rootUri":null}}' | ./lsp-server
```

## ⚡ Performance Tips

- Documents are stored in memory
- Full sync mode is used (not incremental)
- Thread-safe with RWMutex
- Suitable for small to medium files

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/my-feature`
3. Write tests for your changes
4. Ensure tests pass: `make test`
5. Format code: `make fmt`
6. Submit pull request

See `CONTRIBUTING.md` for detailed guidelines.

## 📄 License

MIT License - See `LICENSE` file

## 🔗 Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0 Spec](https://www.jsonrpc.org/specification)
- [Go Documentation](https://golang.org/doc/)

## ✅ Verification

```bash
# Ensure everything works
make clean
make build
make test

# Should see:
# - Build successful
# - Binary created (lsp-server)
# - All tests passing (9/9)
```

## 🎯 Next Steps

1. **Basic Usage**: Run `make test` to verify installation
2. **Integration**: Set up with your favorite editor
3. **Extension**: Add language-specific features
4. **Customization**: Modify for your use case

For detailed information, see the full documentation in:
- `ARCHITECTURE.md` - System design
- `USAGE.md` - Complete usage guide
- `CONTRIBUTING.md` - Development guide

---

**Ready to use!** The LSP server is production-ready for basic text document synchronization and serves as an excellent foundation for building language-specific features.
