# Development Guide

This guide provides information for developers who want to extend or modify the LSP server.

## Project Structure

```
lsp-server/
├── main.go                      # Server entry point
├── go.mod                       # Go module definition
├── go.sum                       # Dependency checksums
├── Makefile                     # Build automation
├── protocol/
│   └── messages.go              # LSP protocol types
├── server/
│   ├── server.go                # Main server implementation
│   ├── document.go              # Document management
│   └── server_test.go           # Unit tests
└── examples/
    ├── simple_client.py         # Python test client
    └── vscode-extension/        # VS Code extension example
```

## Code Organization

### main.go

The entry point that:
- Creates the logger (writes to stderr)
- Initializes the server
- Sets up JSON-RPC 2.0 communication over stdin/stdout
- Waits for the connection to close

### protocol/messages.go

Contains all LSP protocol message types:
- Request/response parameters
- Capabilities structures
- Document-related types
- Position and range types

### server/server.go

The main server logic:
- Request dispatcher (`Handle` method)
- Lifecycle handlers (initialize, shutdown, exit)
- Document synchronization handlers (didOpen, didChange, didClose)
- Language feature handlers (hover)

### server/document.go

Document management:
- Thread-safe document store
- Document content tracking
- Helper methods (e.g., `GetWordAtPosition`)

## Adding New Features

### 1. Adding a New LSP Method

Example: Adding code completion support

#### Step 1: Define protocol types

In `protocol/messages.go`:

```go
// CompletionParams represents the parameters for textDocument/completion
type CompletionParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Position     Position               `json:"position"`
    Context      *CompletionContext     `json:"context,omitempty"`
}

type CompletionContext struct {
    TriggerKind      CompletionTriggerKind `json:"triggerKind"`
    TriggerCharacter *string               `json:"triggerCharacter,omitempty"`
}

type CompletionTriggerKind int

const (
    Invoked                    CompletionTriggerKind = 1
    TriggerCharacter           CompletionTriggerKind = 2
    TriggerForIncompleteCompletion CompletionTriggerKind = 3
)

// CompletionList represents a list of completion items
type CompletionList struct {
    IsIncomplete bool             `json:"isIncomplete"`
    Items        []CompletionItem `json:"items"`
}

type CompletionItem struct {
    Label         string              `json:"label"`
    Kind          CompletionItemKind  `json:"kind,omitempty"`
    Detail        string              `json:"detail,omitempty"`
    Documentation *MarkupContent      `json:"documentation,omitempty"`
    InsertText    string              `json:"insertText,omitempty"`
}

type CompletionItemKind int

const (
    TextCompletion     CompletionItemKind = 1
    MethodCompletion   CompletionItemKind = 2
    FunctionCompletion CompletionItemKind = 3
    // ... more kinds
)
```

#### Step 2: Update server capabilities

In `server/server.go`, update `handleInitialize`:

```go
result := protocol.InitializeResult{
    Capabilities: protocol.ServerCapabilities{
        TextDocumentSync: protocol.TextDocumentSyncOptions{
            OpenClose: true,
            Change:    protocol.Full,
        },
        HoverProvider: true,
        CompletionProvider: &protocol.CompletionOptions{
            TriggerCharacters: []string{".", ":"},
            ResolveProvider:   false,
        },
    },
    // ...
}
```

#### Step 3: Add handler method

In `server/server.go`:

```go
func (s *Server) handleTextDocumentCompletion(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
    var params protocol.CompletionParams
    if err := json.Unmarshal(*req.Params, &params); err != nil {
        s.logger.Printf("Error unmarshaling completion params: %v", err)
        return nil, err
    }

    doc, exists := s.documents.Get(params.TextDocument.URI)
    if !exists {
        return nil, nil
    }

    // Get the current line text
    if params.Position.Line >= len(doc.Lines) {
        return nil, nil
    }
    
    // Generate completion items based on context
    items := []protocol.CompletionItem{
        {
            Label:      "exampleFunction",
            Kind:       protocol.FunctionCompletion,
            Detail:     "func() string",
            InsertText: "exampleFunction()",
        },
        // Add more items...
    }

    return protocol.CompletionList{
        IsIncomplete: false,
        Items:        items,
    }, nil
}
```

#### Step 4: Register in dispatcher

In the `Handle` method's switch statement:

```go
case "textDocument/completion":
    return s.handleTextDocumentCompletion(ctx, req)
```

#### Step 5: Add tests

In `server/server_test.go`:

```go
func TestCompletion(t *testing.T) {
    logger := log.New(os.Stderr, "[TEST] ", log.Lshortfile)
    srv := NewServer(logger)

    // Open a document
    srv.documents.Open("file:///tmp/test.go", "go", 1, "package main\n\nfunc ")

    params := protocol.CompletionParams{
        TextDocument: protocol.TextDocumentIdentifier{
            URI: "file:///tmp/test.go",
        },
        Position: protocol.Position{
            Line:      2,
            Character: 5,
        },
    }

    paramsJSON, _ := json.Marshal(params)
    rawParams := json.RawMessage(paramsJSON)
    req := &jsonrpc2.Request{
        Method: "textDocument/completion",
        Params: &rawParams,
    }

    result, err := srv.handleTextDocumentCompletion(context.Background(), req)
    if err != nil {
        t.Fatalf("handleTextDocumentCompletion failed: %v", err)
    }

    completionList, ok := result.(protocol.CompletionList)
    if !ok {
        t.Fatalf("Expected CompletionList, got %T", result)
    }

    if len(completionList.Items) == 0 {
        t.Error("Expected completion items")
    }
}
```

### 2. Adding Diagnostics (Push Notifications)

Diagnostics are pushed from server to client. Example implementation:

```go
// In protocol/messages.go
type PublishDiagnosticsParams struct {
    URI         string       `json:"uri"`
    Version     *int         `json:"version,omitempty"`
    Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
    Range              Range              `json:"range"`
    Severity           DiagnosticSeverity `json:"severity,omitempty"`
    Code               interface{}        `json:"code,omitempty"`
    Source             string             `json:"source,omitempty"`
    Message            string             `json:"message"`
    RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

type DiagnosticSeverity int

const (
    SeverityError       DiagnosticSeverity = 1
    SeverityWarning     DiagnosticSeverity = 2
    SeverityInformation DiagnosticSeverity = 3
    SeverityHint        DiagnosticSeverity = 4
)
```

To send diagnostics, you need access to the connection:

```go
// In server.go, modify the Server struct
type Server struct {
    documents *DocumentStore
    logger    *log.Logger
    conn      *jsonrpc2.Conn  // Add this
}

// Update handlers to publish diagnostics
func (s *Server) handleTextDocumentDidOpen(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
    // ... existing code ...
    
    // Analyze and send diagnostics
    s.publishDiagnostics(ctx, doc.URI)
    
    return nil, nil
}

func (s *Server) publishDiagnostics(ctx context.Context, uri string) {
    doc, exists := s.documents.Get(uri)
    if !exists {
        return
    }

    diagnostics := []protocol.Diagnostic{}
    
    // Example: Check for lines longer than 80 characters
    for i, line := range doc.Lines {
        if len(line) > 80 {
            diagnostics = append(diagnostics, protocol.Diagnostic{
                Range: protocol.Range{
                    Start: protocol.Position{Line: i, Character: 80},
                    End:   protocol.Position{Line: i, Character: len(line)},
                },
                Severity: protocol.SeverityWarning,
                Source:   "basic-lsp",
                Message:  "Line exceeds 80 characters",
            })
        }
    }

    params := protocol.PublishDiagnosticsParams{
        URI:         uri,
        Diagnostics: diagnostics,
    }

    s.conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}
```

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test -v -run TestServerInitialize ./server/
```

### Manual Testing

Use the Python client:

```bash
make example
```

Or test manually with JSON-RPC messages:

```bash
./lsp-server 2> server.log
```

Then paste JSON-RPC messages with proper headers.

### Debugging

Enable verbose logging in the server:

```go
// In main.go
logger := log.New(os.Stderr, "[LSP] ", log.Ldate|log.Ltime|log.Lshortfile)
```

View logs:

```bash
./lsp-server 2> debug.log
# In another terminal
tail -f debug.log
```

## Best Practices

### Error Handling

Always handle errors properly and log them:

```go
if err := json.Unmarshal(*req.Params, &params); err != nil {
    s.logger.Printf("Error unmarshaling params: %v", err)
    return nil, &jsonrpc2.Error{
        Code:    jsonrpc2.CodeParseError,
        Message: err.Error(),
    }
}
```

### Thread Safety

Use the document store's mutex properly:

```go
// Document store operations are already thread-safe
doc, exists := s.documents.Get(uri)
```

For new concurrent operations, use sync.RWMutex:

```go
type MyStore struct {
    mu   sync.RWMutex
    data map[string]interface{}
}

func (s *MyStore) Get(key string) (interface{}, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    val, ok := s.data[key]
    return val, ok
}
```

### Logging

Log important events with context:

```go
s.logger.Printf("Processing %s for %s at %d:%d", 
    req.Method, params.TextDocument.URI, 
    params.Position.Line, params.Position.Character)
```

### Performance

For large documents, consider:
- Incremental sync instead of full sync
- Caching analysis results
- Debouncing diagnostic updates
- Background processing for expensive operations

## References

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0](https://www.jsonrpc.org/specification)
- [Go LSP Libraries](https://github.com/topics/language-server-protocol?l=go)
- [sourcegraph/jsonrpc2](https://github.com/sourcegraph/jsonrpc2)
