# Boilerplate LSP Server in Go

A minimal but functional Language Server Protocol (LSP) implementation in Go that demonstrates the core concepts and features of LSP. This server can be used as a starting point for building custom language servers.

## Description

This is a boilerplate LSP server implementation that provides the foundational structure for creating language-specific tooling. It handles JSON-RPC 2.0 communication over stdio and implements essential LSP lifecycle methods and features.

The server is built using the official Go LSP libraries:
- `go.lsp.dev/protocol` - LSP protocol definitions and types
- `go.lsp.dev/jsonrpc2` - JSON-RPC 2.0 implementation

## Features Implemented

### LSP Lifecycle Methods
- ✅ **Initialize** - Establishes connection and returns server capabilities
- ✅ **Initialized** - Confirms successful initialization
- ✅ **Shutdown** - Gracefully prepares server for shutdown
- ✅ **Exit** - Terminates the server process

### Text Document Synchronization
- ✅ **textDocument/didOpen** - Tracks when documents are opened in the editor
- ✅ **textDocument/didChange** - Tracks document content changes (full sync mode)
- ✅ **textDocument/didClose** - Cleans up document state when closed

### Language Features
- ✅ **textDocument/hover** - Provides hover information for symbols
- ✅ **textDocument/completion** - Returns completion suggestions

### Additional Features
- ✅ Thread-safe document state management
- ✅ Comprehensive error handling and logging
- ✅ JSON-RPC 2.0 communication over stdio

## Architecture Overview

The server is organized into several components:

```
lsp-server/
├── main.go         # Entry point, stdio setup, connection management
├── server.go       # Core server struct, request routing, lifecycle handlers
├── handlers.go     # LSP feature handlers (hover, completion, etc.)
├── document.go     # Thread-safe document state management
├── go.mod          # Go module definition and dependencies
└── README.md       # This file
```

### Component Responsibilities

- **main.go**: Sets up stdio communication, creates the server, and manages the connection lifecycle
- **server.go**: Routes incoming LSP requests to appropriate handlers, implements lifecycle methods
- **handlers.go**: Implements LSP feature handlers like hover and completion
- **document.go**: Provides thread-safe storage and retrieval of document contents

## Installation

### Prerequisites
- Go 1.21 or higher

### Steps

1. Clone the repository:
```bash
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server
```

2. Download dependencies:
```bash
go mod download
```

3. Build the server:
```bash
go build -o lsp-server
```

## Usage

### Running the Server Standalone

You can test the server by running it directly (it expects JSON-RPC messages on stdin):

```bash
./lsp-server
```

### Integrating with Editors

#### Visual Studio Code

Create a VS Code extension or use a generic LSP client. Add to your `settings.json`:

```json
{
  "lsp-server.enable": true,
  "lsp-server.serverPath": "/path/to/lsp-server"
}
```

Or create a minimal extension with this configuration in `package.json`:

```json
{
  "contributes": {
    "configuration": {
      "type": "object",
      "title": "LSP Server",
      "properties": {
        "lspServer.trace.server": {
          "type": "string",
          "enum": ["off", "messages", "verbose"],
          "default": "off",
          "description": "Traces the communication between VS Code and the language server."
        }
      }
    }
  }
}
```

And in your extension's `extension.ts`:

```typescript
import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
  const serverExecutable = '/path/to/lsp-server';
  
  const serverOptions: ServerOptions = {
    command: serverExecutable,
    args: [],
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: 'file', language: 'plaintext' }],
  };

  client = new LanguageClient(
    'lspServer',
    'LSP Server',
    serverOptions,
    clientOptions
  );

  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
```

#### Neovim

Using Neovim's built-in LSP client:

```lua
-- In your init.lua or a separate config file
vim.api.nvim_create_autocmd("FileType", {
  pattern = "plaintext",  -- Change to your target filetype
  callback = function()
    vim.lsp.start({
      name = "lsp-server",
      cmd = {"/path/to/lsp-server"},
      root_dir = vim.fs.dirname(vim.fs.find({".git"}, { upward = true })[1]),
    })
  end,
})
```

Or using `nvim-lspconfig`:

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

-- Define the custom LSP server
if not configs.lsp_server then
  configs.lsp_server = {
    default_config = {
      cmd = {'/path/to/lsp-server'},
      filetypes = {'plaintext'},  -- Change to your target filetype
      root_dir = lspconfig.util.root_pattern('.git'),
      settings = {},
    },
  }
end

-- Setup the server
lspconfig.lsp_server.setup{}
```

#### Emacs (with lsp-mode)

Add to your Emacs configuration:

```elisp
(with-eval-after-load 'lsp-mode
  (add-to-list 'lsp-language-id-configuration '(plaintext-mode . "plaintext"))
  
  (lsp-register-client
   (make-lsp-client :new-connection (lsp-stdio-connection "/path/to/lsp-server")
                    :major-modes '(plaintext-mode)
                    :server-id 'lsp-server)))
```

## Build and Run Instructions

### Build

```bash
# Build for current platform
go build -o lsp-server

# Build for specific platforms
GOOS=linux GOARCH=amd64 go build -o lsp-server-linux-amd64
GOOS=darwin GOARCH=amd64 go build -o lsp-server-darwin-amd64
GOOS=windows GOARCH=amd64 go build -o lsp-server-windows-amd64.exe
```

### Run

```bash
# Run directly with go
go run .

# Or run the built binary
./lsp-server
```

### Development

For development, you can enable verbose logging by checking the log output. The server logs to stderr by default, which is captured by most LSP clients.

## Extending the Server

To add new LSP features:

1. **Add a new handler method** in `handlers.go`:
   ```go
   func (s *Server) handleNewFeature(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
       // Implementation
   }
   ```

2. **Register the handler** in `server.go`'s `Handle` method:
   ```go
   case protocol.MethodTextDocumentNewFeature:
       return s.handleNewFeature(ctx, reply, req)
   ```

3. **Update server capabilities** in the `handleInitialize` method to advertise the new feature.

## Testing

You can test the server using a simple test client or any LSP-compatible editor. For manual testing, you can send JSON-RPC messages via stdin:

```bash
# Example initialize request
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}' | ./lsp-server
```

## Contributing

This is a boilerplate implementation. Feel free to:
- Add more LSP features (go to definition, find references, rename, etc.)
- Implement language-specific functionality
- Improve error handling and edge cases
- Add comprehensive tests

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/)
- [go.lsp.dev Documentation](https://pkg.go.dev/go.lsp.dev)
- [Building LSP Servers](https://code.visualstudio.com/api/language-extensions/language-server-extension-guide)

## License

This boilerplate is provided as-is for educational and development purposes.
