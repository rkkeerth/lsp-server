# Basic LSP Server

A simple Language Server Protocol (LSP) server implementation demonstrating core LSP capabilities. This implementation uses **only Python standard library** with no external dependencies.

## Features

- **Initialize and Shutdown Lifecycle**: Proper handling of client initialization requests with server capabilities and graceful shutdown
- **Text Document Synchronization**: Maintains synchronized state of text documents as they are opened, modified, and closed
- **Basic Diagnostics**: Analyzes text documents and reports:
  - TODO markers (informational)
  - FIXME markers (warnings)
  - Lines longer than 120 characters (warnings)
  - Duplicate consecutive lines (informational)
- **JSON-RPC 2.0 Message Handling**: Processes messages over standard input/output streams with proper request/response correlation

## Requirements

- Python 3.8 or higher
- **No external dependencies required!** This implementation uses only Python standard library.

## Installation

1. Clone or navigate to the repository:
```bash
cd /projects/sandbox/lsp-server
```

2. No installation needed! The server is ready to run:
```bash
python server.py
```

## Usage

### Running the Server

The LSP server communicates via standard input/output (stdin/stdout). To run it directly:

```bash
python server.py
```

### Testing with the Test Client

A simple test client is provided to demonstrate the server's functionality:

```bash
python test_client.py examples/test.txt
```

This will:
1. Start the LSP server as a subprocess
2. Initialize the client-server connection
3. Open the test document
4. Display any diagnostics found
5. Properly shut down the server

### Integration with Editors

#### VS Code

To use with VS Code, you can create an extension or use the server with existing LSP client extensions. Add the following to your VS Code settings:

```json
{
  "lsp-sample.trace.server": "verbose",
  "lsp-sample.serverPath": "/path/to/lsp-server/server.py"
}
```

#### Neovim

With Neovim's built-in LSP client:

```lua
vim.lsp.start({
  name = 'basic-lsp-server',
  cmd = {'python3', '/path/to/lsp-server/server.py'},
  root_dir = vim.fs.dirname(vim.fs.find({'*.txt'}, { upward = true })[1]),
})
```

#### Emacs (lsp-mode)

```elisp
(lsp-register-client
 (make-lsp-client :new-connection (lsp-stdio-connection '("python3" "/path/to/lsp-server/server.py"))
                  :major-modes '(text-mode)
                  :server-id 'basic-lsp-server))
```

## Example Output

When analyzing a document, the server will provide diagnostics like:

```
============================================================
Diagnostics received for: file:///path/to/test.txt
============================================================
Found 7 issue(s):

  [INFO] Line 5, Col 1: TODO found: Consider addressing this item
  [WARNING] Line 8, Col 1: FIXME found: This requires immediate attention
  [WARNING] Line 10, Col 121: Line too long (145 > 120 characters)
  [INFO] Line 13, Col 1: Duplicate line detected
  [INFO] Line 16, Col 1: TODO found: Consider addressing this item
  [WARNING] Line 17, Col 1: FIXME found: This requires immediate attention
============================================================
```

## Architecture

The server implements the LSP protocol from scratch using only Python's standard library:

### JSON-RPC 2.0 Communication

- Messages are exchanged via stdin/stdout
- Each message has a `Content-Length` header followed by JSON content
- Supports both requests (expect responses) and notifications (no response)

### Key Components

- **LSPServer class**: Main server implementation with:
  - Document lifecycle handlers (didOpen, didChange, didClose)
  - Initialize and shutdown handlers
  - Diagnostic analysis engine
  - JSON-RPC message routing
  
- **analyze_document()**: Core analysis function that scans text for issues using regular expressions

### Message Flow

```
Client → initialize (request) → Server
Client ← initialize result ← Server
Client → initialized (notification) → Server
Client → textDocument/didOpen (notification) → Server
Client ← textDocument/publishDiagnostics (notification) ← Server
Client → shutdown (request) → Server
Client ← shutdown result ← Server
Client → exit (notification) → Server
```

## Logging

The server logs its activity to `/tmp/lsp-server.log` for debugging purposes. You can monitor it with:

```bash
tail -f /tmp/lsp-server.log
```

## Development

### Adding New Diagnostics

To add new diagnostic rules, modify the `analyze_document()` method in `server.py`:

```python
# Example: Detect hardcoded passwords
match = re.search(r'password\s*=\s*["\']', line, re.IGNORECASE)
if match:
    diagnostics.append({
        "range": {
            "start": {"line": line_num, "character": match.start()},
            "end": {"line": line_num, "character": match.end()}
        },
        "message": "Hardcoded password detected",
        "severity": 1,  # Error
        "source": "basic-lsp-server"
    })
```

### Diagnostic Severity Levels

- `1`: Error (red)
- `2`: Warning (yellow)
- `3`: Information (blue)
- `4`: Hint (gray)

### Testing

1. Create a test file in the `examples/` directory
2. Run the test client: `python test_client.py examples/your-test-file.txt`
3. Verify diagnostics are reported correctly
4. Check the log file at `/tmp/lsp-server.log` for detailed server activity

## LSP Specification Compliance

This implementation follows the [Language Server Protocol Specification](https://microsoft.github.io/language-server-protocol/) and implements:

### Lifecycle Messages
- **initialize**: Returns server capabilities
- **initialized**: Acknowledgment notification (no-op)
- **shutdown**: Prepares for server shutdown
- **exit**: Terminates the server process

### Document Synchronization
- **textDocument/didOpen**: Document opened in editor
- **textDocument/didChange**: Document content changed (full sync mode)
- **textDocument/didClose**: Document closed in editor

### Diagnostics
- **textDocument/publishDiagnostics**: Server pushes diagnostics to client

### Capabilities
- Text Document Sync: Full document synchronization
- Diagnostic Provider: On-demand diagnostics

## Implementation Details

### Why No External Dependencies?

This implementation deliberately avoids external libraries to:
1. Demonstrate understanding of the LSP protocol itself
2. Minimize installation complexity
3. Serve as an educational example
4. Enable deployment in restricted environments

### Thread Safety

The server is single-threaded and processes messages sequentially. For production use, consider:
- Asynchronous message handling
- Background document analysis
- Thread-safe document storage

## Limitations

This is a basic implementation for demonstration purposes. Production LSP servers typically include:
- Code completion (textDocument/completion)
- Go to definition (textDocument/definition)
- Find references (textDocument/references)
- Code formatting (textDocument/formatting)
- Hover information (textDocument/hover)
- Incremental document synchronization
- More sophisticated diagnostics
- Configuration options via workspace/didChangeConfiguration
- Multi-workspace support
- Symbol search
- Code actions

## Extending the Server

To add new LSP features:

1. Add a handler method to the `LSPServer` class
2. Register it in the appropriate handler dictionary
3. Update server capabilities in `handle_initialize()`
4. Test with the test client or a real editor

Example - Adding hover support:

```python
def handle_text_document_hover(self, params: Dict[str, Any]) -> Dict[str, Any]:
    """Handle textDocument/hover request."""
    uri = params["textDocument"]["uri"]
    position = params["position"]
    
    # Implement hover logic here
    return {
        "contents": {
            "kind": "plaintext",
            "value": "Hover information"
        }
    }

# Add to request handlers
handlers = {
    "initialize": self.handle_initialize,
    "shutdown": self.handle_shutdown,
    "textDocument/hover": self.handle_text_document_hover,  # New
}

# Update capabilities in handle_initialize
"capabilities": {
    "textDocumentSync": {...},
    "hoverProvider": True,  # New
}
```

## Troubleshooting

### Server not responding
- Check `/tmp/lsp-server.log` for errors
- Ensure Python 3.8+ is installed
- Verify the server process is running

### Diagnostics not appearing
- Check that the file was opened with textDocument/didOpen
- Verify the file URI format is correct
- Look for errors in the log file

### Connection issues
- Ensure stdin/stdout are not being used by other processes
- Check that the editor is configured correctly
- Test with the provided test client first

## License

This is a demonstration implementation for educational purposes.
