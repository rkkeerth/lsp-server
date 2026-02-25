# LSP Server

A boilerplate Language Server Protocol (LSP) server written in Go using only the standard library.

## Features

- JSON-RPC 2.0 transport over stdin/stdout
- Full text document synchronization
- Core LSP lifecycle methods

## Building

```bash
go build -o lsp-server .
```

## Usage

The server communicates over stdin/stdout using the LSP protocol. It is designed to be launched by an editor or IDE that supports the Language Server Protocol.

### Example: VS Code Configuration

Create a file at `.vscode/settings.json` in your project:

```json
{
  "lsp-server.path": "/path/to/lsp-server"
}
```

Or configure it in an extension's `package.json`:

```json
{
  "contributes": {
    "configuration": {
      "properties": {
        "lsp-server.path": {
          "type": "string",
          "default": "lsp-server",
          "description": "Path to the LSP server executable"
        }
      }
    }
  }
}
```

### Example: Neovim Configuration

Using nvim-lspconfig:

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

configs.lsp_server = {
  default_config = {
    cmd = { '/path/to/lsp-server' },
    filetypes = { 'your-language' },
    root_dir = lspconfig.util.root_pattern('.git'),
  },
}

lspconfig.lsp_server.setup({})
```

## Supported LSP Methods

### Lifecycle

| Method | Type | Description |
|--------|------|-------------|
| `initialize` | Request | Initialize the server with client capabilities |
| `initialized` | Notification | Client acknowledgment after initialization |
| `shutdown` | Request | Prepare the server for exit |
| `exit` | Notification | Exit the server process |

### Document Synchronization

| Method | Type | Description |
|--------|------|-------------|
| `textDocument/didOpen` | Notification | Document opened in editor |
| `textDocument/didChange` | Notification | Document content changed (full sync) |
| `textDocument/didClose` | Notification | Document closed in editor |

## Server Capabilities

The server advertises the following capabilities:

```json
{
  "textDocumentSync": {
    "openClose": true,
    "change": 1,
    "save": {
      "includeText": false
    }
  }
}
```

- `change: 1` indicates full document sync (the entire document content is sent on each change)

## Architecture

```
lsp-server/
├── main.go           # Entry point
├── go.mod            # Go module definition
├── jsonrpc/
│   ├── types.go      # JSON-RPC 2.0 types
│   └── transport.go  # stdin/stdout transport
├── protocol/
│   ├── types.go      # LSP protocol types
│   └── methods.go    # LSP method constants
├── document/
│   └── manager.go    # Document state management
└── server/
    └── server.go     # Server implementation
```

## License

MIT
