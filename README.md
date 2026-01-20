# LSP Server Boilerplate (Go)

A boilerplate implementation of the [Language Server Protocol (LSP)](https://microsoft.github.io/language-server-protocol/) in Go. This project provides a foundation for building custom language servers for any programming language or domain-specific language.

## Features

- Full LSP message handling (JSON-RPC 2.0)
- Support for both stdio and TCP communication modes
- Modular handler architecture for easy extension
- Core LSP capabilities implemented:
  - Text document synchronization (open, change, save, close)
  - Completion with trigger characters
  - Hover information
  - Go to definition
  - Find references
  - Document symbols
  - Code formatting
  - Code actions

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── handler/
│   │   └── text_document.go  # Text document handlers
│   ├── lsp/
│   │   ├── handler.go        # LSP request/notification dispatcher
│   │   └── server.go         # Core LSP server implementation
│   └── protocol/
│       ├── messages.go       # JSON-RPC message types
│       └── types.go          # LSP type definitions
├── go.mod
├── Makefile
└── README.md
```

## Requirements

- Go 1.21 or later

## Quick Start

### Build

```bash
# Build the server
make build

# Or build directly with go
go build -o bin/lsp-server ./cmd/server
```

### Run

**Stdio mode (default):**
```bash
./bin/lsp-server
```

**TCP mode (for debugging):**
```bash
./bin/lsp-server -tcp -addr 127.0.0.1:7998
```

### Command Line Options

```
Usage: lsp-server [options]

Options:
  -version    Show version information
  -help       Show help message
  -tcp        Run server in TCP mode
  -addr       TCP address to listen on (default: 127.0.0.1:7998)
  -log        Log file path (default: stderr)
```

## Editor Integration

### VS Code

Create a VS Code extension or add to your settings:

```json
{
  "languageServerExample.serverPath": "/path/to/lsp-server"
}
```

### Neovim (with nvim-lspconfig)

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

configs.my_lsp = {
  default_config = {
    cmd = { '/path/to/lsp-server' },
    filetypes = { 'your-language' },
    root_dir = lspconfig.util.root_pattern('.git', 'go.mod'),
  },
}

lspconfig.my_lsp.setup{}
```

### Emacs (with lsp-mode)

```elisp
(lsp-register-client
  (make-lsp-client
    :new-connection (lsp-stdio-connection "/path/to/lsp-server")
    :major-modes '(your-language-mode)
    :server-id 'my-lsp))
```

## Extending the Server

### Adding New Capabilities

1. **Define protocol types** in `internal/protocol/types.go`

2. **Add handler methods** in `internal/handler/text_document.go`:
   ```go
   func (h *TextDocumentHandler) HandleMyFeature(ctx context.Context, params json.RawMessage) (interface{}, error) {
       // Parse params
       var p MyFeatureParams
       if err := json.Unmarshal(params, &p); err != nil {
           return nil, err
       }
       
       // Implement your logic
       return result, nil
   }
   ```

3. **Register the handler** in `internal/lsp/handler.go`:
   ```go
   case "textDocument/myFeature":
       return h.textDocuments.HandleMyFeature(ctx, params)
   ```

4. **Advertise the capability** in `handleInitialize`:
   ```go
   Capabilities: protocol.ServerCapabilities{
       // ...
       MyFeatureProvider: true,
   }
   ```

### Adding Document Analysis

The `TextDocumentHandler` maintains a map of open documents. You can access document content for analysis:

```go
doc, ok := h.GetDocument(uri)
if ok {
    content := doc.Content
    // Analyze content...
}
```

## Development

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make coverage
```

### Linting

```bash
# Format code
make fmt

# Run vet
make vet

# Run golangci-lint
make lint
```

### Building for Multiple Platforms

```bash
# Build for Linux, macOS, and Windows
make build-all
```

## Protocol Reference

This implementation follows the [LSP 3.17 specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/).

### Supported Methods

| Method | Type | Status |
|--------|------|--------|
| `initialize` | Request | ✅ Implemented |
| `initialized` | Notification | ✅ Implemented |
| `shutdown` | Request | ✅ Implemented |
| `exit` | Notification | ✅ Implemented |
| `textDocument/didOpen` | Notification | ✅ Implemented |
| `textDocument/didChange` | Notification | ✅ Implemented |
| `textDocument/didSave` | Notification | ✅ Implemented |
| `textDocument/didClose` | Notification | ✅ Implemented |
| `textDocument/completion` | Request | ✅ Placeholder |
| `textDocument/hover` | Request | ✅ Placeholder |
| `textDocument/definition` | Request | ✅ Placeholder |
| `textDocument/references` | Request | ✅ Placeholder |
| `textDocument/documentSymbol` | Request | ✅ Placeholder |
| `textDocument/formatting` | Request | ✅ Placeholder |
| `textDocument/codeAction` | Request | ✅ Placeholder |

## License

MIT License - feel free to use this as a starting point for your own language server.
