# LSP Server

A boilerplate Language Server Protocol (LSP) server written in Python using only the standard library.

## Features

- JSON-RPC 2.0 transport over stdin/stdout
- Full text document synchronization
- Core LSP lifecycle methods

## Requirements

Python 3.11+ (standard library only, no dependencies)

## Usage

The server communicates over stdin/stdout using the LSP protocol. It is designed to be launched by an editor or IDE that supports the Language Server Protocol.

```bash
python3 main.py
```

### Example: VS Code Configuration

Create a file at `.vscode/settings.json` in your project:

```json
{
  "lsp-server.command": ["python3", "/path/to/main.py"]
}
```

Or configure it in an extension's `package.json`:

```json
{
  "contributes": {
    "configuration": {
      "properties": {
        "lsp-server.command": {
          "type": "array",
          "default": ["python3", "main.py"],
          "description": "Command to start the LSP server"
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
    cmd = { 'python3', '/path/to/main.py' },
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
├── main.py              # Entry point
├── lsp_server/
│   ├── __init__.py
│   ├── jsonrpc/
│   │   ├── __init__.py
│   │   ├── types.py     # JSON-RPC 2.0 types
│   │   └── transport.py # stdin/stdout transport
│   ├── protocol/
│   │   ├── __init__.py
│   │   ├── types.py     # LSP protocol types
│   │   └── methods.py   # LSP method constants
│   ├── document/
│   │   ├── __init__.py
│   │   └── manager.py   # Document state management
│   └── server/
│       ├── __init__.py
│       └── server.py    # Server implementation
└── tests/               # Test suite
```

## License

MIT
