# LSP Server Usage Guide

This guide explains how to use the basic LSP server implementation.

## Building the Server

### Using Make

```bash
make build
```

### Using Go directly

```bash
go build -o lsp-server ./cmd/lsp-server
```

## Running the Server

The LSP server communicates via standard input/output using the LSP protocol:

```bash
./lsp-server
```

The server will:
1. Listen for JSON-RPC messages on stdin
2. Send responses on stdout
3. Use the LSP message format (Content-Length header + JSON body)

## Testing the Server

### Run Unit Tests

```bash
make test
# or
go test -v ./...
```

### Using the Example Client

Build the server first, then run the test client:

```bash
make build
cd examples
go run test_client.go
```

The test client will:
1. Start the LSP server
2. Send an initialize request
3. Send text document notifications (open, change)
4. Send shutdown request
5. Exit

## Integration with Editors/IDEs

### Visual Studio Code

Create a VS Code extension with the following configuration:

```typescript
import * as path from 'path';
import { workspace, ExtensionContext } from 'vscode';
import {
    LanguageClient,
    LanguageClientOptions,
    ServerOptions,
    TransportKind
} from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: ExtensionContext) {
    const serverExecutable = '/path/to/lsp-server';
    
    const serverOptions: ServerOptions = {
        command: serverExecutable,
        args: [],
        transport: TransportKind.stdio
    };

    const clientOptions: LanguageClientOptions = {
        documentSelector: [{ scheme: 'file', language: 'plaintext' }],
        synchronize: {
            fileEvents: workspace.createFileSystemWatcher('**/*')
        }
    };

    client = new LanguageClient(
        'basicLspServer',
        'Basic LSP Server',
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

### Neovim

Configure with `nvim-lspconfig`:

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

-- Define the LSP server configuration
if not configs.basic_lsp then
  configs.basic_lsp = {
    default_config = {
      cmd = {'/path/to/lsp-server'},
      filetypes = {'text', 'plaintext'},
      root_dir = function(fname)
        return lspconfig.util.find_git_ancestor(fname) or vim.fn.getcwd()
      end,
      settings = {},
    },
  }
end

-- Start the LSP server
lspconfig.basic_lsp.setup{}
```

### Emacs (with lsp-mode)

Add to your Emacs configuration:

```elisp
(require 'lsp-mode)

(lsp-register-client
 (make-lsp-client
  :new-connection (lsp-stdio-connection "/path/to/lsp-server")
  :major-modes '(text-mode)
  :server-id 'basic-lsp))

(add-hook 'text-mode-hook #'lsp)
```

## Message Format

All messages follow the LSP specification with a Content-Length header:

```
Content-Length: <length>\r\n
\r\n
<JSON content>
```

### Example Request

```
Content-Length: 123\r\n
\r\n
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":1234,"rootUri":"file:///workspace"}}
```

### Example Response

```
Content-Length: 98\r\n
\r\n
{"jsonrpc":"2.0","id":1,"result":{"capabilities":{"textDocumentSync":{"openClose":true}}}}
```

## Supported Methods

### Lifecycle

- `initialize` - Initialize the server and negotiate capabilities
- `initialized` - Notification that client is initialized
- `shutdown` - Prepare for shutdown
- `exit` - Exit the server process

### Text Synchronization

- `textDocument/didOpen` - Document opened in editor
- `textDocument/didChange` - Document content changed
- `textDocument/didClose` - Document closed in editor

## Message Examples

### Initialize Request

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": 12345,
    "rootUri": "file:///workspace",
    "capabilities": {
      "textDocument": {
        "synchronization": {
          "dynamicRegistration": true,
          "willSave": true,
          "didSave": true
        }
      }
    }
  }
}
```

### Initialize Response

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
          "includeText": true
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

### didOpen Notification

```json
{
  "jsonrpc": "2.0",
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///workspace/file.txt",
      "languageId": "plaintext",
      "version": 1,
      "text": "Hello, World!"
    }
  }
}
```

### didChange Notification

```json
{
  "jsonrpc": "2.0",
  "method": "textDocument/didChange",
  "params": {
    "textDocument": {
      "uri": "file:///workspace/file.txt",
      "version": 2
    },
    "contentChanges": [
      {
        "text": "Hello, LSP!"
      }
    ]
  }
}
```

## Debugging

### Enable Logging

Redirect stderr to a log file:

```bash
./lsp-server 2>/tmp/lsp-server.log
```

### Test with netcat

You can manually test the server using netcat:

```bash
# Start server
./lsp-server

# In another terminal, send a message
echo -e 'Content-Length: 90\r\n\r\n{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":1234,"rootUri":null}}' | ./lsp-server
```

### Trace Messages

To see all messages exchanged, use the test client with verbose output:

```bash
cd examples
go run test_client.go 2>&1 | tee /tmp/lsp-trace.log
```

## Extending the Server

### Adding New Capabilities

1. Define the protocol types in `internal/protocol/types.go`
2. Add handler method in `internal/server/server.go`
3. Register the method in the `Handle()` switch
4. Update server capabilities in `handleInitialize()`

### Example: Adding Hover Support

```go
// In internal/protocol/types.go
type HoverParams struct {
    TextDocument TextDocumentIdentifier `json:"textDocument"`
    Position     Position                `json:"position"`
}

type Hover struct {
    Contents MarkupContent `json:"contents"`
}

type MarkupContent struct {
    Kind  string `json:"kind"`
    Value string `json:"value"`
}

// In internal/server/server.go
func (s *Server) handleHover(params json.RawMessage) (interface{}, error) {
    var hoverParams protocol.HoverParams
    if err := json.Unmarshal(params, &hoverParams); err != nil {
        return nil, err
    }
    
    // Get document
    doc, exists := s.GetDocument(hoverParams.TextDocument.URI)
    if !exists {
        return nil, nil
    }
    
    // Return hover information
    return protocol.Hover{
        Contents: protocol.MarkupContent{
            Kind:  "plaintext",
            Value: "Hover information for the document",
        },
    }, nil
}

// In Handle() method
case "textDocument/hover":
    return s.handleHover(params)

// In handleInitialize()
capabilities.HoverProvider = true
```

## Performance Tips

1. **Document Size**: Be mindful of large documents in full sync mode
2. **Concurrency**: The server uses mutexes; avoid long-running operations in handlers
3. **Memory**: Documents are kept in memory; implement cleanup for unused documents
4. **Logging**: Avoid excessive logging in production

## Troubleshooting

### Server Won't Start

- Check that the binary has execute permissions: `chmod +x lsp-server`
- Verify Go version: `go version` (requires Go 1.21+)
- Check for port conflicts if using TCP transport

### Client Can't Connect

- Verify the server is running: `ps aux | grep lsp-server`
- Check that stdin/stdout are not being used by other processes
- Ensure the client is using the correct communication method (stdio)

### Messages Not Processed

- Verify message format includes Content-Length header
- Check JSON is valid: Use `jq` or online JSON validators
- Ensure `\r\n\r\n` separator is present after headers
- Check content length matches actual JSON length

### Test Failures

```bash
# Run specific test
go test -v ./internal/server -run TestServerInitialize

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0](https://www.jsonrpc.org/specification)
- [Go Documentation](https://golang.org/doc/)
- [LSP Implementations](https://langserver.org/)
