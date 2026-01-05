# LSP Server Architecture

This document describes the architecture and design decisions of the LSP server implementation.

## Overview

The LSP server is built following a modular, layered architecture that separates concerns and promotes maintainability. It follows the official [Language Server Protocol specification](https://microsoft.github.io/language-server-protocol/) and implements JSON-RPC 2.0 for communication.

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────┐
│                      LSP Client                         │
│              (VS Code, Neovim, Emacs, etc.)             │
└─────────────────────┬───────────────────────────────────┘
                      │ stdin/stdout
                      │ JSON-RPC 2.0
                      │
┌─────────────────────▼───────────────────────────────────┐
│                   main.go                               │
│              (Entry Point & Logger Setup)               │
└─────────────────────┬───────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────┐
│                  server/server.go                       │
│          ┌─────────────────────────────────┐            │
│          │  Request Router & Dispatcher    │            │
│          │  - JSON-RPC message handling    │            │
│          │  - Lifecycle management         │            │
│          │  - Capability negotiation       │            │
│          └─────────────┬───────────────────┘            │
└────────────────────────┼────────────────────────────────┘
                         │
         ┌───────────────┼───────────────┐
         │               │               │
┌────────▼──────┐ ┌─────▼────────┐ ┌───▼──────────┐
│  document/    │ │  handlers/   │ │  handlers/   │
│  manager.go   │ │  handlers.go │ │  symbol_     │
│               │ │              │ │  index.go    │
│ - Document    │ │ - Hover      │ │              │
│   storage     │ │ - Definition │ │ - Symbol     │
│ - Content     │ │ - References │ │   indexing   │
│   tracking    │ │ - Symbols    │ │ - Fast       │
│ - Thread-safe │ │ - Completion │ │   lookups    │
│   access      │ │ - Diagnostics│ │              │
└───────────────┘ └──────────────┘ └──────────────┘
```

## Core Components

### 1. Main Entry Point (`main.go`)

**Responsibility**: Application initialization and lifecycle

- Initializes structured logging with Zap
- Creates the server instance
- Sets up stdin/stdout communication
- Handles graceful shutdown

**Key Design Decisions**:
- Uses structured logging to stderr (stdout reserved for LSP protocol)
- Minimal logic - delegates everything to the server package

### 2. Server Package (`server/`)

**Responsibility**: Protocol handling and request routing

#### Components:

##### Server (`server.go`)
- Manages JSON-RPC 2.0 connection
- Routes incoming requests to appropriate handlers
- Manages server state (initialized, shutdown)
- Thread-safe state management with mutexes
- Publishes notifications to clients

**Key Methods**:
```go
Run(ctx, stdin, stdout)      // Starts the server
handle(ctx, reply, request)  // Routes requests
setupCapabilities()          // Defines server capabilities
```

**Thread Safety**:
- Uses `sync.RWMutex` for state protection
- Concurrent request handling
- Safe shutdown mechanism

### 3. Document Management (`document/`)

**Responsibility**: Document lifecycle and content management

#### Components:

##### Manager (`manager.go`)
- Maintains open documents in memory
- Provides thread-safe CRUD operations
- Tracks document versions
- Splits content into lines for analysis

**Data Structure**:
```go
type Document struct {
    URI     DocumentURI
    Content string
    Version int32
    Lines   []string  // Pre-split for performance
}
```

**Thread Safety**:
- `sync.RWMutex` for concurrent read/write
- Multiple readers, single writer pattern
- Efficient for read-heavy workloads

**Key Operations**:
- `Open()`: Add new document
- `Update()`: Update existing document
- `Close()`: Remove document
- `Get()`: Retrieve document (thread-safe)
- `GetAll()`: List all documents

### 4. Handler Package (`handlers/`)

**Responsibility**: Implement LSP features

#### Components:

##### Handler (`handlers.go`)
Implements all LSP feature methods:

**Feature Categories**:

1. **Navigation Features**:
   - `Hover()`: Symbol information on hover
   - `Definition()`: Go to definition
   - `References()`: Find all references

2. **Symbol Features**:
   - `DocumentSymbols()`: Extract symbols from document
   - `WorkspaceSymbols()`: Search across workspace

3. **Code Intelligence**:
   - `Completion()`: Code completion suggestions
   - `GetDiagnostics()`: Syntax and semantic analysis

**Implementation Details**:

- **Pattern Matching**: Uses regex for symbol detection
  ```go
  funcRegex := regexp.MustCompile(`func\s+(\w+)\s*\(`)
  typeRegex := regexp.MustCompile(`type\s+(\w+)\s+(struct|interface)`)
  ```

- **Context-Aware Analysis**: Examines surrounding code for better results

- **Performance**: Pre-compiled regex patterns, efficient string operations

##### Symbol Index (`symbol_index.go`)
Fast symbol lookup and indexing:

```go
type SymbolIndex struct {
    symbols map[string][]SymbolInfo
    mu      sync.RWMutex
}
```

**Operations**:
- `Add()`: Index a symbol
- `Get()`: Retrieve symbols by name
- `Remove()`: Remove symbols by URI
- `Clear()`: Reset index

**Thread Safety**: Protected with `sync.RWMutex`

## Communication Protocol

### JSON-RPC 2.0 Message Format

Messages are transmitted over stdin/stdout with HTTP-style headers:

```
Content-Length: 123\r\n
\r\n
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "textDocument/hover",
  "params": {...}
}
```

### Message Types

1. **Request**: Requires a response
   ```json
   {"jsonrpc":"2.0", "id":1, "method":"...", "params":{}}
   ```

2. **Response**: Reply to a request
   ```json
   {"jsonrpc":"2.0", "id":1, "result":{}}
   ```

3. **Notification**: No response expected
   ```json
   {"jsonrpc":"2.0", "method":"...", "params":{}}
   ```

### Lifecycle Flow

```
Client                          Server
  │                               │
  ├─ initialize ─────────────────►│
  │                               ├─ Process capabilities
  │◄────────── InitializeResult ─┤
  │                               │
  ├─ initialized ────────────────►│
  │                               ├─ Server ready
  │                               │
  ├─ textDocument/didOpen ───────►│
  │                               ├─ Track document
  │◄────── publishDiagnostics ────┤
  │                               │
  ├─ textDocument/hover ─────────►│
  │◄────────── Hover result ──────┤
  │                               │
  ├─ shutdown ───────────────────►│
  │◄────────── null ──────────────┤
  │                               │
  ├─ exit ───────────────────────►│
  │                               ├─ Server exits
```

## Concurrency Model

### Design Principles

1. **Non-blocking I/O**: JSON-RPC connection runs in goroutine
2. **Thread-safe data structures**: All shared state protected
3. **Read-heavy optimization**: `RWMutex` allows concurrent reads
4. **No race conditions**: Verified with `go test -race`

### Synchronization Points

- Document manager: `sync.RWMutex`
- Symbol index: `sync.RWMutex`
- Server state: `sync.RWMutex`

### Request Handling

Each request is handled concurrently:
```go
conn.Go(ctx, handler)  // Non-blocking
```

## Error Handling Strategy

### Levels

1. **Protocol Errors**: Return JSON-RPC error responses
   ```go
   return reply(ctx, nil, jsonrpc2.ErrMethodNotFound)
   ```

2. **Application Errors**: Log and return nil results
   ```go
   logger.Error("Failed to process", zap.Error(err))
   return reply(ctx, nil, nil)
   ```

3. **Fatal Errors**: Log and exit
   ```go
   logger.Fatal("Cannot start server", zap.Error(err))
   ```

### Logging

Uses structured logging (Zap):
- **Debug**: Detailed request/response info
- **Info**: Important state changes
- **Warn**: Non-critical issues
- **Error**: Failures that don't stop server
- **Fatal**: Critical failures requiring shutdown

Logs go to stderr to avoid interfering with stdin/stdout protocol.

## Performance Considerations

### Optimizations

1. **Pre-split lines**: Documents store lines as `[]string`
   - Avoids repeated splitting during analysis
   - O(1) line access

2. **Compiled regex patterns**: Reused across invocations
   - One-time compilation cost
   - Fast matching

3. **Efficient string operations**: 
   - Minimal allocations
   - In-place character checks

4. **Symbol indexing**: O(1) lookup by name
   - Pre-indexed during document analysis
   - Fast workspace-wide searches

### Memory Management

- Documents are kept in memory (typical for LSP)
- Old documents released on close
- No memory leaks (verified with profiling)

### Scalability

- **Concurrent requests**: Goroutine per request
- **Multiple documents**: HashMap for O(1) access
- **Large files**: Line-based processing (streaming potential)

## Extension Points

### Adding New Features

1. **Add capability** in `setupCapabilities()`:
   ```go
   capabilities.SignatureHelpProvider = &SignatureHelpOptions{...}
   ```

2. **Implement handler** in `handlers/handlers.go`:
   ```go
   func (h *Handler) SignatureHelp(...) *SignatureHelp {
       // Implementation
   }
   ```

3. **Route request** in `handle()`:
   ```go
   case "textDocument/signatureHelp":
       return s.handleSignatureHelp(ctx, reply, req)
   ```

### Custom Language Support

The server is language-agnostic by design. To support specific languages:

1. **Update symbol extraction** in `extractDocumentSymbols()`
2. **Add language-specific patterns** to regex collection
3. **Implement semantic analysis** in `GetDiagnostics()`

## Testing Strategy

### Unit Tests

- Document manager operations
- Symbol index operations
- Handler logic (with mocks)

### Integration Tests

- Full lifecycle (initialize → work → shutdown)
- Multi-document scenarios
- Concurrent request handling

### Manual Testing

- Use example files in `examples/`
- Test with real editors
- Monitor logs for issues

## Dependencies

### External Libraries

1. **go.lsp.dev/protocol** (v0.12.0)
   - LSP type definitions
   - Protocol constants
   - Version compatibility

2. **go.lsp.dev/jsonrpc2** (v0.10.0)
   - JSON-RPC 2.0 implementation
   - Stream handling
   - Connection management

3. **go.uber.org/zap** (v1.26.0)
   - Structured logging
   - High performance
   - Production-ready

### Standard Library Usage

- `context`: Cancellation and timeouts
- `sync`: Thread synchronization
- `encoding/json`: JSON parsing
- `io`: Stream handling
- `regexp`: Pattern matching
- `strings`: Text processing

## Security Considerations

### Input Validation

- JSON unmarshaling with error checking
- Bounds checking on positions
- Safe string operations (no buffer overflows)

### Resource Limits

Current implementation has no explicit limits. Production deployments should consider:
- Maximum document size
- Maximum number of documents
- Request timeout
- Memory limits

### No Privilege Escalation

- Reads stdin/stdout only
- No file system access (client provides content)
- No network access
- No shell execution

## Future Enhancements

### Planned Features

1. **Incremental Sync**: Reduce bandwidth for large files
2. **Semantic Tokens**: Syntax highlighting
3. **Code Actions**: Quick fixes and refactorings
4. **Signature Help**: Parameter hints
5. **Call Hierarchy**: Function call trees
6. **Type Hierarchy**: Type relationship navigation

### Performance Improvements

1. **Incremental Parsing**: Parse only changed regions
2. **Background Indexing**: Index symbols asynchronously
3. **Caching**: Cache expensive operations
4. **Streaming**: Process large files in chunks

### Architecture Evolution

1. **Plugin System**: Load language-specific handlers
2. **Configuration**: Runtime configuration support
3. **Metrics**: Performance and usage metrics
4. **Remote Protocol**: Support network communication

## Conclusion

This LSP server provides a solid foundation for language intelligence features. The modular architecture allows easy extension, the concurrency model ensures responsiveness, and adherence to the LSP specification ensures compatibility with all major editors.

The implementation prioritizes:
- **Correctness**: Full protocol compliance
- **Performance**: Efficient algorithms and data structures
- **Maintainability**: Clear separation of concerns
- **Extensibility**: Easy to add new features
- **Safety**: Thread-safe, no data races

For questions or contributions, see the main README.md.
