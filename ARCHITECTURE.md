# Architecture Documentation

This document provides an in-depth explanation of the LSP server architecture, design decisions, and implementation details.

## Overview

This is a basic but complete implementation of a Language Server Protocol (LSP) server in Go. It demonstrates fundamental LSP concepts including:

- JSON-RPC 2.0 message protocol
- LSP lifecycle management
- Text document synchronization
- Concurrent request handling
- Thread-safe state management

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         LSP Client                          │
│                  (Editor/IDE/Test Client)                   │
└────────────────────┬────────────────────────────────────────┘
                     │
                     │ stdin/stdout
                     │ JSON-RPC 2.0
                     │
┌────────────────────▼────────────────────────────────────────┐
│                       LSP Server                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Message Reader/Writer                   │  │
│  │  - Read from stdin with buffering                   │  │
│  │  - Parse Content-Length headers                     │  │
│  │  - Write to stdout with proper formatting           │  │
│  └──────────────────┬───────────────────────────────────┘  │
│                     │                                       │
│  ┌──────────────────▼───────────────────────────────────┐  │
│  │              Message Handler                         │  │
│  │  - Parse JSON-RPC messages                          │  │
│  │  - Route to appropriate handlers                    │  │
│  │  - Process requests concurrently                    │  │
│  └──────┬────────────────────────────────────────┬──────┘  │
│         │                                        │          │
│  ┌──────▼──────────┐                    ┌───────▼───────┐  │
│  │    Lifecycle     │                    │  Text Document │  │
│  │    Handlers      │                    │    Handlers    │  │
│  │  - initialize    │                    │  - didOpen     │  │
│  │  - initialized   │                    │  - didChange   │  │
│  │  - shutdown      │                    │  - didClose    │  │
│  │  - exit          │                    │                │  │
│  └─────────┬────────┘                    └────────┬───────┘  │
│            │                                      │          │
│            └──────────┬───────────────────────────┘          │
│                       │                                      │
│  ┌────────────────────▼──────────────────────────────────┐  │
│  │              Server State                             │  │
│  │  - Current state (uninitialized/initialized/etc)     │  │
│  │  - Open documents map (URI -> Document)              │  │
│  │  - Thread-safe with RWMutex                          │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## Component Details

### 1. Message Protocol Layer

**File**: `internal/lsp/server.go` (readMessage, writeMessage)

**Responsibilities**:
- Parse LSP message format (Content-Length header + JSON body)
- Handle message serialization/deserialization
- Buffer management for stdin/stdout

**Design Decisions**:
- Uses `bufio.Reader` for efficient reading
- Implements header parsing for Content-Length
- Handles both `\r\n` and `\n` line endings for compatibility

**Message Format**:
```
Content-Length: 123\r\n
\r\n
{"jsonrpc":"2.0","id":1,"method":"initialize",...}
```

### 2. JSON-RPC Layer

**File**: `internal/lsp/protocol.go`

**Responsibilities**:
- Define JSON-RPC 2.0 message structures
- Provide type-safe representation of LSP protocol
- Support requests, responses, and notifications

**Key Types**:
- `Request`: Represents a request (has ID)
- `Response`: Represents a response (has ID and Result or Error)
- `Notification`: Represents a notification (no ID)
- `Error`: Represents JSON-RPC error

**Design Decisions**:
- Uses `interface{}` for ID to support both string and number IDs
- Separates requests and notifications at the type level
- Includes all core LSP protocol types

### 3. Server Core

**File**: `internal/lsp/server.go`

**Responsibilities**:
- Manage server lifecycle
- Route messages to handlers
- Maintain server state
- Coordinate concurrent operations

**Key Structures**:
```go
type Server struct {
    reader    *bufio.Reader          // Input stream
    writer    io.Writer              // Output stream
    mu        sync.RWMutex           // Protects state and documents
    state     ServerState            // Current server state
    documents map[string]*Document   // Open documents
}
```

**State Machine**:
```
Uninitialized → Initializing → Initialized → ShuttingDown → Shutdown
     │              │              │              │            │
     └──────────────┴──────────────┴──────────────┴────────────┘
                        (error states)
```

**Design Decisions**:
- Uses goroutines for concurrent message handling
- Read-write locks for efficient concurrent access
- Separate state management methods for thread safety

### 4. Lifecycle Management

**File**: `internal/lsp/lifecycle.go`

**Responsibilities**:
- Handle initialization sequence
- Negotiate capabilities
- Manage shutdown sequence

**Initialization Flow**:
```
Client                          Server
  │                               │
  ├─────── initialize ───────────>│
  │                               ├── Validate state
  │                               ├── Parse capabilities
  │                               ├── Set state: Initializing
  │<──── InitializeResult ────────┤
  │                               ├── Set state: Initialized
  ├─────── initialized ──────────>│
  │                               ├── Log confirmation
  │                               │
```

**Shutdown Flow**:
```
Client                          Server
  │                               │
  ├─────── shutdown ─────────────>│
  │                               ├── Set state: ShuttingDown
  │<──── null result ─────────────┤
  │                               │
  ├─────── exit ─────────────────>│
  │                               ├── Set state: Shutdown
  │                               ├── os.Exit(0)
```

**Design Decisions**:
- Strict state validation
- Proper error responses for invalid states
- Exit codes: 0 for clean shutdown, 1 otherwise

### 5. Text Document Synchronization

**File**: `internal/lsp/textdocument.go`

**Responsibilities**:
- Track open documents
- Handle document changes
- Maintain document state

**Document Lifecycle**:
```
didOpen → didChange (multiple) → didClose
   │           │                      │
   └───────────┴──────────────────────┘
          (document in memory)
```

**Synchronization Mode**:
- Uses **Full Document Sync** (TextDocumentSyncKind.Full)
- Client sends entire document on each change
- Simpler than incremental sync
- Sufficient for most use cases

**Design Decisions**:
- Store complete document content in memory
- Track version numbers for consistency
- Thread-safe document access via server mutex

## Concurrency Model

### Message Processing

Each incoming message is processed in its own goroutine:

```go
for {
    msg := readMessage()
    go handleMessage(msg)  // Concurrent processing
}
```

**Benefits**:
- Non-blocking message processing
- Better responsiveness
- Handles multiple simultaneous requests

**Considerations**:
- Requires thread-safe state management
- Must handle concurrent document access
- Potential for race conditions without proper locking

### Thread Safety

**Read-Write Locks**:
```go
type Server struct {
    mu        sync.RWMutex
    documents map[string]*Document
}

// Read operations (multiple concurrent readers)
func (s *Server) GetDocument(uri string) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    // Read access
}

// Write operations (exclusive access)
func (s *Server) SetDocument(uri string, doc *Document) {
    s.mu.Lock()
    defer s.mu.Unlock()
    // Write access
}
```

**Benefits**:
- Multiple concurrent readers
- Exclusive writer access
- Prevents data races

## Data Flow

### Request Flow

```
1. Client sends request via stdin
2. Server reads message (readMessage)
3. Parse JSON into Request struct
4. Route to handler based on method
5. Handler processes (concurrent goroutine)
6. Handler generates result
7. Server sends response (sendResponse)
8. Response written to stdout
```

### Notification Flow

```
1. Client sends notification via stdin
2. Server reads message (readMessage)
3. Parse JSON into Request struct (ID is nil)
4. Route to handler based on method
5. Handler processes (concurrent goroutine)
6. Handler updates state
7. No response sent (notifications don't expect responses)
```

## Error Handling

### Error Codes

Standard JSON-RPC error codes:
- `-32700`: Parse error
- `-32600`: Invalid request
- `-32601`: Method not found
- `-32602`: Invalid params
- `-32603`: Internal error

### Error Response Example

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32601,
    "message": "Method not found: textDocument/completion"
  }
}
```

### Error Handling Strategy

1. **Parse Errors**: Log and continue (don't crash)
2. **Invalid Requests**: Send error response
3. **Unknown Methods**: Send MethodNotFound error
4. **Handler Errors**: Log and send InternalError

## Capability Negotiation

### Server Capabilities

Declared in initialize response:

```go
ServerCapabilities{
    TextDocumentSync: TextDocumentSyncOptions{
        OpenClose: true,           // Support didOpen/didClose
        Change:    Full,           // Full document sync
        Save:      &SaveOptions{   // Save notifications
            IncludeText: false,
        },
    },
    CompletionProvider: &CompletionOptions{
        TriggerCharacters: []string{".", ":"},
    },
}
```

### Client Capabilities

Received in initialize request:
- Used to determine what features client supports
- Server can adjust behavior based on client capabilities
- Currently logged but not used (basic implementation)

## Performance Considerations

### Memory Management

- **Documents**: Stored in memory (full content)
- **Trade-off**: Speed vs. memory usage
- **Suitable for**: Projects with reasonable file sizes
- **Not suitable for**: Very large files (>10MB)

### Message Buffering

- Uses `bufio.Reader` for efficient reading
- Buffers reduce syscall overhead
- Important for high-frequency notifications

### Concurrency

- **Pro**: Handles multiple requests concurrently
- **Con**: Higher memory usage per goroutine
- **Trade-off**: Responsiveness vs. resource usage

## Extensibility Points

### Adding New Capabilities

1. **Protocol Types**: Add to `protocol.go`
2. **Capability Declaration**: Update `lifecycle.go`
3. **Handler Implementation**: Create new handler
4. **Route Registration**: Add to `server.go`

### Example: Adding Completion

```go
// 1. Protocol types (already defined)
type CompletionParams struct { ... }
type CompletionList struct { ... }

// 2. Declare capability
ServerCapabilities{
    CompletionProvider: &CompletionOptions{
        TriggerCharacters: []string{".", ":"},
    },
}

// 3. Implement handler
func (s *Server) handleTextDocumentCompletion(req Request) {
    // Implementation
}

// 4. Register route
case "textDocument/completion":
    s.handleTextDocumentCompletion(request)
```

## Security Considerations

### Input Validation

- All JSON parsing is error-checked
- Content-Length validated before reading
- Unknown methods handled gracefully

### Resource Limits

**Current Implementation**:
- No limits on document size
- No limits on number of open documents
- No request rate limiting

**Production Considerations**:
- Add maximum document size limits
- Limit number of concurrent documents
- Implement request rate limiting
- Add timeout for long-running operations

### Denial of Service

**Potential Vulnerabilities**:
- Large Content-Length values
- Rapid message flooding
- Memory exhaustion from large documents

**Mitigations** (not implemented, but recommended):
- Content-Length upper bound
- Request rate limiting
- Memory usage monitoring
- Graceful degradation

## Testing Strategy

### Unit Tests

- Test individual components in isolation
- Mock dependencies (reader/writer)
- Test both success and error cases
- Use table-driven tests

### Integration Tests

- Test complete message flows
- Use test client to simulate real client
- Verify state transitions
- Test concurrent operations

### Manual Testing

- VS Code extension for real-world testing
- Test client for automated scenarios
- Netcat for low-level protocol testing

## Future Enhancements

### Short Term

- [ ] Add more LSP capabilities (hover, completion, goto definition)
- [ ] Implement incremental document sync
- [ ] Add comprehensive error recovery
- [ ] Improve logging with levels

### Long Term

- [ ] Add document analysis features
- [ ] Implement workspace symbol search
- [ ] Add code actions and quick fixes
- [ ] Support multiple workspace folders
- [ ] Add configuration support

## References

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0](https://www.jsonrpc.org/specification)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)

## Conclusion

This architecture provides a solid foundation for an LSP server implementation. It demonstrates:

- Proper protocol implementation
- Thread-safe concurrent operations
- Clear separation of concerns
- Extensibility for future features

The design prioritizes simplicity and clarity, making it an excellent starting point for learning LSP or building more complex language servers.
