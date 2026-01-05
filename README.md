# LSP Server

A complete Language Server Protocol (LSP) server implementation in GoLang that follows the LSP specification and provides intelligent language features for code editors and IDEs.

## Features

This LSP server implements the following core features:

### Lifecycle Management
- **Initialize**: Negotiates capabilities between client and server
- **Initialized**: Confirms successful initialization
- **Shutdown**: Gracefully shuts down the server
- **Exit**: Terminates the server process

### Text Document Synchronization
- **didOpen**: Tracks when documents are opened in the editor
- **didChange**: Updates document content on changes (full synchronization)
- **didClose**: Removes documents from tracking when closed
- **didSave**: Handles document save events

### Language Features
- **Hover**: Displays information about symbols when hovering over them
- **Go to Definition**: Navigates to symbol definitions
- **Find References**: Locates all usages of a symbol across documents
- **Document Symbols**: Provides document outline and navigation
- **Workspace Symbols**: Searches for symbols across the entire workspace
- **Code Completion**: Offers intelligent code suggestions including keywords and symbols
- **Diagnostics**: Reports errors, warnings, and hints (TODO/FIXME detection)

## Architecture

The project is structured with clear separation of concerns:

```
lsp-server/
├── main.go                 # Entry point and server initialization
├── server/
│   └── server.go          # Core LSP server and request routing
├── document/
│   └── manager.go         # Document management and text operations
├── handlers/
│   ├── handlers.go        # LSP feature implementations
│   └── symbol_index.go    # Symbol indexing for fast lookups
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

### Key Components

- **Server**: Manages JSON-RPC communication, handles client requests, and routes them to appropriate handlers
- **Document Manager**: Maintains the state of all open documents with thread-safe access
- **Handlers**: Implements LSP features like hover, completion, definition, and references
- **Symbol Index**: Maintains a fast lookup index for workspace symbols

## Protocol Implementation

This server implements the Language Server Protocol using:

- **JSON-RPC 2.0**: Standard protocol for client-server communication
- **stdin/stdout**: Communication channel as per LSP specification
- **Header-based messaging**: Each message includes `Content-Length` header

### Message Format

```
Content-Length: <bytes>\r\n
\r\n
{JSON-RPC message}
```

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Download dependencies
go mod download

# Build the server
go build -o lsp-server .
```

## Usage

### Running the Server

The LSP server communicates via stdin/stdout, which is the standard for LSP implementations:

```bash
./lsp-server
```

### Integration with Editors

#### VS Code

Create a `.vscode/launch.json` configuration:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "type": "node",
      "request": "launch",
      "name": "Launch LSP Server",
      "program": "${workspaceFolder}/path/to/lsp-server"
    }
  ]
}
```

Or integrate via a VS Code extension with the following configuration in `package.json`:

```json
{
  "contributes": {
    "configuration": {
      "type": "object",
      "title": "LSP Server Configuration",
      "properties": {
        "lspServer.executablePath": {
          "type": "string",
          "default": "lsp-server",
          "description": "Path to the LSP server executable"
        }
      }
    }
  }
}
```

#### Neovim

Add to your Neovim configuration:

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

-- Define the custom LSP server
if not configs.lsp_server then
  configs.lsp_server = {
    default_config = {
      cmd = {'/path/to/lsp-server'},
      filetypes = {'<your-language>'},
      root_dir = lspconfig.util.root_pattern('.git', 'go.mod'),
      settings = {},
    },
  }
end

-- Start the LSP server
lspconfig.lsp_server.setup{}
```

#### Emacs (lsp-mode)

```elisp
(with-eval-after-load 'lsp-mode
  (add-to-list 'lsp-language-id-configuration '(<mode> . "language-id"))
  (lsp-register-client
   (make-lsp-client
    :new-connection (lsp-stdio-connection "/path/to/lsp-server")
    :major-modes '(<your-mode>)
    :server-id 'lsp-server)))
```

## Examples

### Testing with LSP Client

You can test the server using any LSP client. Here's an example using a simple JSON-RPC communication:

#### Initialize Request

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": 12345,
    "rootUri": "file:///path/to/project",
    "capabilities": {}
  }
}
```

#### Initialize Response

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "capabilities": {
      "textDocumentSync": {
        "openClose": true,
        "change": 1
      },
      "hoverProvider": true,
      "definitionProvider": true,
      "referencesProvider": true,
      "documentSymbolProvider": true,
      "workspaceSymbolProvider": true,
      "completionProvider": {
        "triggerCharacters": [".", ":", ">"]
      }
    },
    "serverInfo": {
      "name": "lsp-server",
      "version": "1.0.0"
    }
  }
}
```

### Document Operations

#### Open Document

```json
{
  "jsonrpc": "2.0",
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.go",
      "languageId": "go",
      "version": 1,
      "text": "package main\n\nfunc main() {\n\tprintln(\"Hello\")\n}"
    }
  }
}
```

#### Hover Request

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "textDocument/hover",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.go"
    },
    "position": {
      "line": 2,
      "character": 5
    }
  }
}
```

#### Completion Request

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "textDocument/completion",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.go"
    },
    "position": {
      "line": 3,
      "character": 1
    }
  }
}
```

## Development

### Running Tests

```bash
go test ./...
```

### Running with Debug Logging

The server uses structured logging with zap. Logs are written to stderr to avoid interfering with stdin/stdout communication:

```bash
./lsp-server 2> lsp-server.log
```

### Adding New Features

To add a new LSP feature:

1. Add the capability to `server/server.go` in `setupCapabilities()`
2. Add a new handler method in `handlers/handlers.go`
3. Add the route in `server/server.go` in the `handle()` method
4. Update this README with the new feature

## Thread Safety

The server is designed to handle concurrent requests safely:

- Document manager uses `sync.RWMutex` for safe concurrent access
- Symbol index is protected with mutex locks
- Server state (initialized, shutdown) is protected with mutex

## Error Handling

The server implements comprehensive error handling:

- All JSON unmarshaling errors are caught and logged
- Invalid requests return appropriate JSON-RPC error responses
- Document not found scenarios are handled gracefully
- Connection errors are logged with context

## Dependencies

- `go.lsp.dev/protocol` - LSP protocol types and structures
- `go.lsp.dev/jsonrpc2` - JSON-RPC 2.0 implementation
- `go.uber.org/zap` - Structured logging

## Contributing

Contributions are welcome! Please ensure:

1. Code follows Go conventions and best practices
2. All features include appropriate error handling
3. Changes maintain thread safety
4. Documentation is updated accordingly

## Troubleshooting

### Server Not Starting

- Ensure Go 1.21+ is installed: `go version`
- Check dependencies are downloaded: `go mod download`
- Verify the executable has proper permissions

### Connection Issues

- Ensure your editor is configured to use stdin/stdout
- Check server logs for initialization errors
- Verify the server path is correct in editor configuration

### Feature Not Working

- Check if the capability is enabled in server response
- Review server logs for errors (stderr)
- Ensure the document is properly opened with `textDocument/didOpen`

## References

- [Language Server Protocol Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [LSP Go Protocol Types](https://pkg.go.dev/go.lsp.dev/protocol)

## Version History

### 1.0.0 (Current)
- Initial release
- Full LSP lifecycle support
- Text document synchronization
- Hover, definition, and references support
- Document and workspace symbols
- Code completion
- Basic diagnostics

## Future Enhancements

- Incremental document synchronization
- Advanced syntax analysis
- Code formatting support
- Refactoring operations
- Signature help
- Code actions and quick fixes
- Semantic highlighting
- Call hierarchy
- Type hierarchy

---

**Author**: rkkeerth  
**Repository**: https://github.com/rkkeerth/lsp-server
