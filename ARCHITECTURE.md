# LSP Server Architecture

This document describes the architecture and design of the basic LSP server implementation.

## Overview

The Language Server Protocol (LSP) server is designed with a modular architecture that separates concerns:

1. **JSON-RPC Layer** - Handles low-level communication protocol
2. **Protocol Layer** - Defines LSP-specific data structures
3. **Server Layer** - Implements LSP business logic
4. **Main Entry Point** - Orchestrates the components

## Architecture Diagram

```
┌─────────────────────────────────────┐
│         Client (Editor/IDE)         │
└────────────┬────────────────────────┘
             │ stdin/stdout
             │ JSON-RPC 2.0
             ▼
┌─────────────────────────────────────┐
│       cmd/lsp-server/main.go        │
│     (Application Entry Point)       │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│      internal/jsonrpc/rpc.go        │
│   (JSON-RPC 2.0 Communication)      │
│  - Message reading/writing          │
│  - Request/Response handling        │
│  - Header parsing                   │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│    internal/server/server.go        │
│     (LSP Server Logic)              │
│  - Initialize/Shutdown              │
│  - Document management              │
│  - Method dispatching               │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│    internal/protocol/types.go       │
│     (LSP Protocol Types)            │
│  - InitializeParams                 │
│  - ServerCapabilities               │
│  - TextDocument types               │
└─────────────────────────────────────┘
```

## Component Details

### 1. JSON-RPC Layer (`internal/jsonrpc`)

#### Message Types

The JSON-RPC layer implements three message types:

- **Request**: Has an ID and expects a response
- **Response**: Reply to a request, contains result or error
- **Notification**: No ID, no response expected

#### Communication Protocol

The server uses the LSP message format:

```
Content-Length: <length>\r\n
\r\n
<JSON content>
```

Key files:
- `message.go` - Message type definitions and constructors
- `rpc.go` - Core RPC communication logic

#### Key Functions

- `ReadMessage()` - Reads LSP-formatted messages from input stream
- `WriteMessage()` - Writes LSP-formatted messages to output stream
- `ProcessMessage()` - Routes incoming messages to appropriate handlers
- `Run()` - Main message processing loop

### 2. Protocol Layer (`internal/protocol`)

Defines all LSP protocol structures according to the specification:

#### Initialize/Shutdown

- `InitializeParams` - Client initialization parameters
- `InitializeResult` - Server capabilities response
- `ClientCapabilities` - What the client supports
- `ServerCapabilities` - What the server supports

#### Text Document Synchronization

- `DidOpenTextDocumentParams` - Document opened notification
- `DidChangeTextDocumentParams` - Document changed notification
- `DidCloseTextDocumentParams` - Document closed notification
- `TextDocumentItem` - Complete document representation
- `TextDocumentContentChangeEvent` - Document change event

#### Sync Modes

The server supports:
- **Full sync** (`TextDocumentSyncKind.Full`) - Client sends entire document on each change
- **Incremental sync** (`TextDocumentSyncKind.Incremental`) - Client sends only changes

Currently implemented: **Full sync**

### 3. Server Layer (`internal/server`)

#### Server Structure

```go
type Server struct {
    initialized bool
    shutdown    bool
    documents   map[string]*Document
    mu          sync.RWMutex
}
```

#### Document Management

Documents are stored in-memory with:
- URI (unique identifier)
- Language ID
- Version number
- Full text content

#### Method Handlers

The server implements handlers for:

1. **initialize** - Returns server capabilities
2. **initialized** - Marks server as ready
3. **shutdown** - Prepares for graceful shutdown
4. **exit** - Exits the process
5. **textDocument/didOpen** - Opens and stores a document
6. **textDocument/didChange** - Updates document content
7. **textDocument/didClose** - Removes document from memory

#### Thread Safety

All document operations are protected by a `sync.RWMutex` to ensure thread-safe access.

### 4. Main Entry Point (`cmd/lsp-server`)

Responsibilities:
1. Create LSP server instance
2. Create JSON-RPC handler with stdin/stdout
3. Start the message processing loop
4. Handle errors and exit codes

## Communication Flow

### 1. Initialize Sequence

```
Client                          Server
  │                               │
  ├──── initialize request ──────►│
  │                               │ (process capabilities)
  │◄──── initialize response ─────┤
  │                               │
  ├──── initialized notification ►│
  │                               │ (server ready)
```

### 2. Document Lifecycle

```
Client                          Server
  │                               │
  ├──── didOpen ─────────────────►│ (store document)
  │                               │
  ├──── didChange ───────────────►│ (update document)
  ├──── didChange ───────────────►│ (update document)
  │                               │
  ├──── didClose ────────────────►│ (remove document)
```

### 3. Shutdown Sequence

```
Client                          Server
  │                               │
  ├──── shutdown request ────────►│
  │                               │ (cleanup)
  │◄──── shutdown response ───────┤
  │                               │
  ├──── exit notification ───────►│ (exit process)
```

## Extensibility

### Adding New LSP Methods

To add support for new LSP methods:

1. **Define protocol types** in `internal/protocol/types.go`
2. **Add handler** in `internal/server/server.go`
3. **Register method** in the `Handle()` switch statement
4. **Update capabilities** in the initialize response

Example:

```go
// 1. Define types
type HoverParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Position     Position                `json:"position"`
}

type Hover struct {
    Contents string `json:"contents"`
    Range    *Range `json:"range,omitempty"`
}

// 2. Add handler
func (s *Server) handleHover(params json.RawMessage) (interface{}, error) {
    var hoverParams HoverParams
    if err := json.Unmarshal(params, &hoverParams); err != nil {
        return nil, err
    }
    
    // Implement hover logic
    return Hover{Contents: "Hover information"}, nil
}

// 3. Register in Handle()
case "textDocument/hover":
    return s.handleHover(params)

// 4. Update capabilities
capabilities.HoverProvider = true
```

### Custom Features

Add custom server-specific features by:

1. Extending `ServerCapabilities` with custom fields
2. Implementing handlers for custom methods
3. Documenting custom behavior in client configuration

## Testing

### Unit Tests

Each component has unit tests:
- `internal/jsonrpc/message_test.go` - Message construction tests
- `internal/server/server_test.go` - Server logic tests

Run tests:
```bash
go test ./...
```

### Integration Testing

Use the example client in `examples/test_client.go`:

```bash
# Build server
make build

# Run test client
cd examples
go run test_client.go
```

## Performance Considerations

1. **Concurrency**: The server uses mutex locks for document access
2. **Memory**: Documents are stored entirely in memory
3. **Synchronization**: Full sync mode is simpler but less efficient for large documents

## Future Enhancements

Potential improvements:

1. **Incremental sync** - More efficient document updates
2. **Language features** - Add hover, completion, goto definition
3. **Diagnostics** - Report errors and warnings
4. **Configuration** - Support workspace configuration
5. **Multi-workspace** - Handle multiple workspace folders
6. **Logging** - Structured logging with levels
7. **Metrics** - Performance monitoring
8. **Caching** - Cache parsed document structures

## References

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [Go LSP Libraries](https://github.com/sourcegraph/go-lsp)
