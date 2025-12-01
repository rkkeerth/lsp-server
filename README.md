# Basic LSP Server

A simple Language Server Protocol (LSP) server implementation in Go that demonstrates core LSP functionality.

## Features

- **JSON-RPC 2.0 Communication**: Full bidirectional communication over stdin/stdout
- **LSP Lifecycle Management**: 
  - `initialize`: Server initialization with capability negotiation
  - `initialized`: Confirmation notification
  - `shutdown`: Graceful shutdown request
  - `exit`: Server exit notification
- **Text Document Synchronization**: 
  - `textDocument/didOpen`: Track newly opened documents
  - `textDocument/didChange`: Monitor document changes (full sync mode)
  - `textDocument/didClose`: Clean up closed documents
- **Language Features**:
  - `textDocument/hover`: Provides hover information for identifiers
- **Error Handling**: Proper error responses and logging
- **Thread-Safe Document Management**: Concurrent-safe document store

## Prerequisites

- Go 1.21 or later
- An LSP-compatible editor or IDE (VS Code, Neovim, Emacs, etc.)

## Building the Server

```bash
# Download dependencies
go mod download

# Build the server
go build -o lsp-server .
```

This will create an executable named `lsp-server` in the current directory.

## Running the Server

The LSP server communicates via stdin/stdout, so it's typically launched by an editor/IDE. However, you can test it manually:

```bash
# Run the server (it will wait for input on stdin)
./lsp-server
```

The server logs diagnostic information to stderr, so you'll see startup messages like:
```
[LSP] 2025/12/01 19:06:27 Starting LSP server...
[LSP] 2025/12/01 19:06:27 LSP server is ready and listening on stdin/stdout
```

## Connecting to Editors

### VS Code

1. Create a VS Code extension or use an existing one like `vscode-languageclient`

2. Add configuration to your `settings.json`:

```json
{
  "languageServerExample.server": {
    "command": "/path/to/lsp-server",
    "args": []
  }
}
```

3. Or create a minimal extension with a `package.json`:

```json
{
  "name": "basic-lsp-client",
  "version": "0.1.0",
  "engines": {
    "vscode": "^1.75.0"
  },
  "activationEvents": ["onLanguage:plaintext"],
  "main": "./extension.js",
  "dependencies": {
    "vscode-languageclient": "^8.0.0"
  }
}
```

4. And an `extension.js`:

```javascript
const { LanguageClient } = require('vscode-languageclient/node');

function activate(context) {
  const serverOptions = {
    command: '/path/to/lsp-server',
    args: []
  };

  const clientOptions = {
    documentSelector: [{ scheme: 'file', language: 'plaintext' }]
  };

  const client = new LanguageClient(
    'basicLspServer',
    'Basic LSP Server',
    serverOptions,
    clientOptions
  );

  client.start();
}

exports.activate = activate;
```

### Neovim

Using the built-in LSP client in Neovim (0.5+):

```lua
-- Add to your init.lua or init.vim (as lua block)
vim.api.nvim_create_autocmd("FileType", {
  pattern = "text",
  callback = function()
    vim.lsp.start({
      name = "basic-lsp-server",
      cmd = {"/path/to/lsp-server"},
      root_dir = vim.fs.dirname(vim.fs.find({"go.mod", ".git"}, { upward = true })[1]),
    })
  end,
})
```

### Emacs

Using `lsp-mode`:

```elisp
(with-eval-after-load 'lsp-mode
  (add-to-list 'lsp-language-id-configuration '(text-mode . "text"))
  (lsp-register-client
   (make-lsp-client
    :new-connection (lsp-stdio-connection "/path/to/lsp-server")
    :major-modes '(text-mode)
    :server-id 'basic-lsp-server)))
```

## Testing the Server Manually

You can test the server manually using JSON-RPC messages:

1. Start the server:
```bash
./lsp-server 2> server.log
```

2. Send an initialize request (paste this followed by Enter):
```json
Content-Length: 324

{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":null,"rootUri":"file:///tmp","capabilities":{"textDocument":{"hover":{"contentFormat":["markdown","plaintext"]},"synchronization":{"didSave":true}}},"trace":"off","workspaceFolders":[{"uri":"file:///tmp","name":"tmp"}]}}
```

The server will respond with its capabilities.

3. Send an initialized notification:
```json
Content-Length: 51

{"jsonrpc":"2.0","method":"initialized","params":{}}
```

4. Open a document:
```json
Content-Length: 168

{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///tmp/test.txt","languageId":"text","version":1,"text":"Hello world\nTest document"}}}
```

5. Request hover information:
```json
Content-Length: 137

{"jsonrpc":"2.0","id":2,"method":"textDocument/hover","params":{"textDocument":{"uri":"file:///tmp/test.txt"},"position":{"line":0,"character":1}}}
```

## Project Structure

```
lsp-server/
├── main.go              # Entry point, sets up JSON-RPC connection
├── go.mod               # Go module dependencies
├── protocol/
│   └── messages.go      # LSP protocol message types and structures
└── server/
    ├── server.go        # Main server logic and request handlers
    └── document.go      # Document store and management
```

## Architecture

### Communication Layer

The server uses the `jsonrpc2` library to handle JSON-RPC 2.0 communication over stdin/stdout. Messages follow the LSP specification with `Content-Length` headers.

### Document Management

Documents are stored in a thread-safe `DocumentStore` that tracks:
- Document URI
- Language ID
- Version number
- Full content and line-by-line representation

### Request Handling

The server implements a request dispatcher that routes LSP methods to appropriate handlers. Each handler:
1. Unmarshals the request parameters
2. Performs the requested operation
3. Returns a response or error

### Hover Feature

The hover implementation:
1. Locates the word at the cursor position
2. Identifies word boundaries using alphanumeric characters and underscores
3. Returns formatted markdown with the word and document language

## Logging

The server logs all activity to stderr, including:
- Server startup/shutdown
- Incoming requests and their methods
- Document operations (open, change, close)
- Hover requests with position and word information

View logs by redirecting stderr:
```bash
./lsp-server 2> server.log
```

## Extending the Server

### Adding New Language Features

1. Define the request/response types in `protocol/messages.go`
2. Add a case in the `Handle` method switch statement in `server/server.go`
3. Implement the handler method
4. Update the server capabilities in `handleInitialize`

Example for adding diagnostics:

```go
// In protocol/messages.go
type PublishDiagnosticsParams struct {
    URI         string       `json:"uri"`
    Diagnostics []Diagnostic `json:"diagnostics"`
}

type Diagnostic struct {
    Range    Range  `json:"range"`
    Severity int    `json:"severity"`
    Message  string `json:"message"`
}

// In server/server.go
func (s *Server) publishDiagnostics(ctx context.Context, conn *jsonrpc2.Conn, uri string) {
    // Analyze document and create diagnostics
    params := protocol.PublishDiagnosticsParams{
        URI:         uri,
        Diagnostics: []protocol.Diagnostic{},
    }
    conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}
```

### Supporting Incremental Sync

To support incremental document synchronization:

1. Change the sync mode in capabilities:
```go
Change: protocol.Incremental,
```

2. Update the `handleTextDocumentDidChange` method to apply partial changes based on the `Range` field in change events.

## Troubleshooting

### Server not responding

- Check that the server is running: `ps aux | grep lsp-server`
- Verify logs in stderr output
- Ensure the editor is configured with the correct server path

### Hover not working

- Confirm the document is opened with `textDocument/didOpen`
- Check that the cursor is on an identifier (alphanumeric or underscore)
- Review server logs for the hover request

### Connection issues

- Verify stdin/stdout are not being used by other processes
- Check that the editor is sending properly formatted JSON-RPC messages
- Ensure `Content-Length` headers are correct

## Contributing

Contributions are welcome! Some ideas for enhancements:

- Add code completion support
- Implement diagnostics/linting
- Add go-to-definition functionality
- Support incremental text synchronization
- Add workspace symbols support
- Implement document formatting

## References

- [LSP Specification](https://microsoft.github.io/language-server-protocol/specifications/specification-current/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [Go LSP Libraries](https://github.com/sourcegraph/jsonrpc2)
