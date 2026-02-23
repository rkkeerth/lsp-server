# LSP Server Boilerplate

A boilerplate Language Server Protocol (LSP) implementation in Go. This provides a working foundation that can be extended to build custom language servers.

## Features

- JSON-RPC 2.0 communication over stdio
- Server lifecycle management (initialize, shutdown, exit)
- Text document synchronization (open, close, change)
- Scaffolding for common LSP features:
  - Completion
  - Hover
  - Diagnostics

## Prerequisites

- Go 1.21 or later

## Building

```bash
# Build all packages
go build ./...

# Build the server binary
go build -o lsp-server ./cmd/server
```

## Running

```bash
# Run the server (communicates over stdio)
./lsp-server

# Show version
./lsp-server --version
```

## Testing

```bash
go test ./...
```

## Project Structure

```
├── cmd/
│   └── server/
│       └── main.go           # Entry point
├── internal/
│   ├── document/
│   │   └── document.go       # Text document management
│   ├── handler/
│   │   ├── completion.go     # Completion handler
│   │   ├── diagnostics.go    # Diagnostics helper
│   │   ├── hover.go          # Hover handler
│   │   ├── lifecycle.go      # Initialize/shutdown handlers
│   │   └── textdocument.go   # Document sync handlers
│   ├── protocol/
│   │   ├── jsonrpc.go        # JSON-RPC 2.0 implementation
│   │   ├── methods.go        # LSP method constants
│   │   └── types.go          # LSP protocol types
│   └── server/
│       └── server.go         # Server core
├── go.mod
└── README.md
```

## Extending the Server

### Adding a New Feature

1. **Add the method constant** in `internal/protocol/methods.go`
2. **Add parameter/result types** in `internal/protocol/types.go`
3. **Create a handler** in `internal/handler/`
4. **Register the handler** in `internal/server/server.go`
5. **Advertise the capability** in `HandleInitialize`

### Example: Adding Go-to-Definition

```go
// 1. Add method constant (methods.go)
const MethodTextDocumentDefinition = "textDocument/definition"

// 2. Add types (types.go)
type DefinitionParams struct {
    TextDocumentPositionParams
}

// 3. Create handler (handler/definition.go)
type Definition struct {
    store *document.Store
}

func (h *Definition) Handle(params *protocol.DefinitionParams) *protocol.Location {
    // Implementation here
    return &protocol.Location{
        URI:   params.TextDocument.URI,
        Range: protocol.Range{...},
    }
}

// 4. Register in server.go handleRequest()
case protocol.MethodTextDocumentDefinition:
    // ...

// 5. Add capability in lifecycle.go
DefinitionProvider: true,
```

## Editor Integration

### VS Code

Create a VS Code extension that spawns this server. In your extension's `activate` function:

```typescript
const serverOptions: ServerOptions = {
    run: { command: '/path/to/lsp-server' },
    debug: { command: '/path/to/lsp-server' }
};

const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: 'file', language: 'yourlang' }],
};

const client = new LanguageClient(
    'yourLangServer',
    'Your Language Server',
    serverOptions,
    clientOptions
);

client.start();
```

### Neovim

Using nvim-lspconfig:

```lua
local configs = require('lspconfig.configs')

configs.yourlsp = {
    default_config = {
        cmd = { '/path/to/lsp-server' },
        filetypes = { 'yourlang' },
        root_dir = function(fname)
            return vim.fn.getcwd()
        end,
    },
}

require('lspconfig').yourlsp.setup{}
```

## License

MIT
