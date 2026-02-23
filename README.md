# Go LSP Server

A Language Server Protocol (LSP) implementation in Go, providing a foundation for building language servers for editors and IDEs.

## Features

- JSON-RPC 2.0 protocol over stdin/stdout
- LSP lifecycle management (initialize, shutdown, exit)
- Document synchronization (full sync mode)
  - textDocument/didOpen
  - textDocument/didChange
  - textDocument/didClose
- Extensible architecture for adding language features

## Building

```bash
go build ./...
```

## Testing

```bash
go test ./... -v
```

## Usage

The server communicates via stdin/stdout using the LSP protocol:

```bash
./lsp-server
```

### Example Integration

The server expects JSON-RPC messages with Content-Length headers:

```
Content-Length: 123\r\n
\r\n
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{...}}
```

### Initialize Request

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": 1234,
    "rootUri": "file:///path/to/project",
    "capabilities": {}
  }
}
```

### Server Capabilities

The server responds with its capabilities:

```json
{
  "capabilities": {
    "textDocumentSync": {
      "openClose": true,
      "change": 1,
      "save": {"includeText": true}
    }
  },
  "serverInfo": {
    "name": "go-lsp-server",
    "version": "0.1.0"
  }
}
```

## Extending

To add new language features:

1. Define types in `lsp/` package
2. Add method constants to `lsp/methods.go`
3. Implement handlers in `server/handlers.go`
4. Register handlers in `server/server.go`

## Architecture

```
main.go          - Entry point
rpc/             - JSON-RPC 2.0 protocol layer
  rpc.go         - Message encoding/decoding
lsp/             - LSP type definitions
  types.go       - Core types (Position, Range, etc.)
  initialize.go  - Initialize request/response types
  textdocument.go- Document sync types
  methods.go     - Method name constants
server/          - Server implementation
  server.go      - Main loop and dispatch
  handlers.go    - Method handlers
  state.go       - Document state management
```

## License

MIT
