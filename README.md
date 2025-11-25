# LSP Server - Basic Language Server Protocol Implementation in Go

A minimal but functional Language Server Protocol (LSP) server implementation in Go, demonstrating the core concepts of LSP including initialization, lifecycle management, and text document synchronization.

## Features

- ✅ **JSON-RPC 2.0 Communication**: Full implementation of JSON-RPC message handling
- ✅ **LSP Lifecycle Management**: Proper initialization, shutdown, and exit handling
- ✅ **Text Document Synchronization**: Support for didOpen, didChange, and didClose notifications
- ✅ **Concurrent Message Processing**: Handles multiple requests concurrently
- ✅ **Thread-Safe Document Store**: Safe concurrent access to document state
- ✅ **Full Document Sync**: Complete document synchronization (Full sync mode)

## Project Structure

```
lsp-server/
├── main.go                    # Entry point
├── internal/lsp/
│   ├── protocol.go           # LSP protocol types and structures
│   ├── server.go             # Core server implementation
│   ├── lifecycle.go          # Initialization and lifecycle handlers
│   └── textdocument.go       # Text document synchronization handlers
├── go.mod                    # Go module definition
└── README.md                 # This file
```

## Installation

### Prerequisites

- Go 1.21 or higher

### Building from Source

```bash
# Clone the repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Build the binary
go build -o lsp-server .

# Or install globally
go install .
```

## Usage

### Running the Server

The LSP server communicates via stdin/stdout following the LSP specification:

```bash
./lsp-server
```

The server will log diagnostic information to stderr while using stdout for LSP communication.

### Testing with a Client

You can test the server with any LSP-compatible editor or client. Here's an example using Visual Studio Code:

1. **Configure VS Code** by creating a `.vscode/settings.json`:

```json
{
  "languageServerExample.trace.server": "verbose"
}
```

2. **Manual Testing with netcat** (for debugging):

```bash
# Start the server
./lsp-server

# Send initialize request (you'll need to format this properly with headers)
```

### Example LSP Message Flow

1. **Initialize Request**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": 1234,
    "rootUri": "file:///path/to/workspace",
    "capabilities": {}
  }
}
```

2. **Initialize Response**:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "capabilities": {
      "textDocumentSync": {
        "openClose": true,
        "change": 1,
        "save": {
          "includeText": false
        }
      }
    },
    "serverInfo": {
      "name": "basic-lsp-server",
      "version": "0.1.0"
    }
  }
}
```

3. **Initialized Notification**:
```json
{
  "jsonrpc": "2.0",
  "method": "initialized",
  "params": {}
}
```

4. **Text Document Did Open Notification**:
```json
{
  "jsonrpc": "2.0",
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.txt",
      "languageId": "plaintext",
      "version": 1,
      "text": "Hello, World!"
    }
  }
}
```

## Supported LSP Methods

### Lifecycle Methods

| Method | Type | Description |
|--------|------|-------------|
| `initialize` | Request | Initialize the server with client capabilities |
| `initialized` | Notification | Confirm initialization complete |
| `shutdown` | Request | Request server shutdown |
| `exit` | Notification | Exit the server process |

### Text Document Synchronization Methods

| Method | Type | Description |
|--------|------|-------------|
| `textDocument/didOpen` | Notification | Document opened in the client |
| `textDocument/didChange` | Notification | Document content changed |
| `textDocument/didClose` | Notification | Document closed in the client |

## Architecture

### Message Flow

```
Client -> stdin -> Server.readMessage() -> Server.handleMessage() -> Handler
                                                                         |
Client <- stdout <- Server.writeMessage() <- Server.sendResponse() <----+
```

### Key Components

1. **Server**: Main server structure managing state and connections
2. **Protocol Types**: Complete LSP type definitions for requests/responses
3. **Lifecycle Handlers**: Initialize, shutdown, and exit handling
4. **Text Document Handlers**: Document synchronization logic
5. **Message Router**: Dispatches messages to appropriate handlers

### Concurrency Model

- Each incoming message is processed in a separate goroutine
- Document store uses read-write locks for thread-safe access
- Server state is protected by mutex locks

## Development

### Running Tests

```bash
go test ./...
```

### Code Organization

- `internal/lsp/`: LSP implementation (not exported)
  - `protocol.go`: Type definitions matching LSP specification
  - `server.go`: Core server logic and message handling
  - `lifecycle.go`: Initialize, shutdown, exit handlers
  - `textdocument.go`: Document synchronization handlers

### Adding New Capabilities

To add new LSP capabilities:

1. Define types in `protocol.go`
2. Add server capability in `lifecycle.go` (initialize handler)
3. Implement handler in appropriate file
4. Register handler in `server.go` (handleMessage method)

Example: Adding hover support

```go
// 1. Already defined in protocol.go: HoverParams, Hover types

// 2. Update initialize handler in lifecycle.go:
result := InitializeResult{
    Capabilities: ServerCapabilities{
        // ... existing capabilities
        HoverProvider: true,
    },
}

// 3. Create handler in new file or existing file:
func (s *Server) handleTextDocumentHover(request Request) {
    // Implementation
}

// 4. Register in server.go:
case "textDocument/hover":
    s.handleTextDocumentHover(request)
```

## Technical Details

### JSON-RPC 2.0 Message Format

Messages are sent using the following format:

```
Content-Length: {length}\r\n
\r\n
{json-content}
```

Where `{length}` is the byte length of the JSON content.

### Document Synchronization

This server uses **Full Document Sync** mode (TextDocumentSyncKind.Full), meaning:
- On `didChange`, the entire document content is sent
- No incremental updates are needed
- Simpler implementation, suitable for most use cases

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

### Development Guidelines

1. Follow Go best practices and idioms
2. Maintain thread safety for concurrent operations
3. Add appropriate logging for debugging
4. Follow LSP specification closely
5. Keep code simple and readable

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/specifications/specification-current/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [Go Documentation](https://golang.org/doc/)

## License

MIT License - See LICENSE file for details

## Author

rkkeerth

## Version

0.1.0 - Initial implementation with basic LSP features
