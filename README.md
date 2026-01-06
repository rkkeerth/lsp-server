# LSP Server

A comprehensive Language Server Protocol (LSP) server implementation in Go that provides intelligent language features for code editors and IDEs. This server implements the LSP specification to enable features like code completion, hover information, go-to-definition, find references, and more.

## What is an LSP Server?

The Language Server Protocol (LSP) is an open, JSON-RPC-based protocol created by Microsoft that standardizes the communication between development tools and language servers. It allows editors and IDEs to provide intelligent language features without having to implement language-specific logic for each tool.

Key benefits of LSP:
- **Write Once, Use Everywhere**: A single language server can be used across multiple editors
- **Rich Language Features**: Provides code completion, diagnostics, refactoring, and navigation
- **Separation of Concerns**: Language intelligence is decoupled from the editor
- **Protocol Standardization**: Well-defined communication protocol between client and server

## Features

This LSP server implementation provides a complete set of language intelligence features:

### 🔄 Lifecycle Management
- **Initialize**: Negotiates capabilities between client and server
- **Initialized**: Confirms successful initialization
- **Shutdown**: Gracefully shuts down the server
- **Exit**: Terminates the server process

### 📄 Text Document Synchronization
- **didOpen**: Tracks when documents are opened in the editor
- **didChange**: Updates document content as you type (full document synchronization)
- **didClose**: Removes documents from tracking when closed
- **didSave**: Handles document save events

### 🚀 Language Intelligence Features
- **Hover**: Displays detailed information about symbols under the cursor
- **Go to Definition**: Navigates directly to symbol definitions
- **Find References**: Locates all usages of a symbol across the workspace
- **Document Symbols**: Provides document outline for quick navigation
- **Workspace Symbols**: Searches for symbols across the entire workspace
- **Code Completion**: Offers intelligent code suggestions including keywords and symbols
- **Diagnostics**: Reports errors, warnings, and hints (including TODO/FIXME detection)

### 🔒 Robustness
- Thread-safe document management with concurrent access support
- Comprehensive error handling and recovery
- Structured logging for debugging and monitoring
- JSON-RPC 2.0 compliant message handling

## Installation

### Prerequisites

- **Go 1.21 or higher**: [Download Go](https://golang.org/dl/)
- **Git**: For cloning the repository

### Building from Source

```bash
# Clone the repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Download dependencies
go mod download

# Build the server binary
go build -o lsp-server .

# (Optional) Install globally
go install .
```

### Verify Installation

```bash
# Check that the binary was created
./lsp-server --version

# Or if installed globally
lsp-server --version
```

## Usage

### Running the Server

The LSP server communicates via stdin/stdout, which is the standard for LSP implementations:

```bash
# Start the server (it will wait for LSP messages on stdin)
./lsp-server

# Redirect logs to a file for debugging
./lsp-server 2> lsp-server.log
```

The server uses:
- **stdin**: Receives LSP requests from the client
- **stdout**: Sends LSP responses to the client
- **stderr**: Outputs diagnostic and debug logs

### Integration with Code Editors

#### Visual Studio Code

Create a custom VS Code extension or configure an existing LSP client extension:

**Option 1: Using a generic LSP client extension**

1. Install the "Generic LSP Client" extension
2. Configure in `.vscode/settings.json`:

```json
{
  "genericLSPClient.languageServerConfigs": {
    "mylanguage": {
      "command": "/path/to/lsp-server",
      "languageIds": ["mylanguage"]
    }
  }
}
```

**Option 2: Create a custom extension**

In your extension's `extension.ts`:

```typescript
import * as vscode from 'vscode';
import { LanguageClient, LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node';

let client: LanguageClient;

export function activate(context: vscode.ExtensionContext) {
  const serverOptions: ServerOptions = {
    command: '/path/to/lsp-server',
    args: []
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: 'file', language: 'mylanguage' }],
  };

  client = new LanguageClient('lspServer', 'LSP Server', serverOptions, clientOptions);
  client.start();
}

export function deactivate(): Thenable<void> | undefined {
  if (!client) {
    return undefined;
  }
  return client.stop();
}
```

#### Neovim

Add to your Neovim configuration (Lua):

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

-- Define the custom LSP server
if not configs.lsp_server then
  configs.lsp_server = {
    default_config = {
      cmd = {'/path/to/lsp-server'},
      filetypes = {'mylanguage'},
      root_dir = lspconfig.util.root_pattern('.git', 'go.mod'),
      settings = {},
    },
  }
end

-- Start the LSP server
lspconfig.lsp_server.setup{
  on_attach = function(client, bufnr)
    -- Enable completion
    vim.api.nvim_buf_set_option(bufnr, 'omnifunc', 'v:lua.vim.lsp.omnifunc')
    
    -- Key mappings
    local opts = { noremap=true, silent=true, buffer=bufnr }
    vim.keymap.set('n', 'gd', vim.lsp.buf.definition, opts)
    vim.keymap.set('n', 'K', vim.lsp.buf.hover, opts)
    vim.keymap.set('n', 'gr', vim.lsp.buf.references, opts)
  end
}
```

#### Emacs (lsp-mode)

Add to your Emacs configuration:

```elisp
(with-eval-after-load 'lsp-mode
  (add-to-list 'lsp-language-id-configuration '(mylanguage-mode . "mylanguage"))
  (lsp-register-client
   (make-lsp-client
    :new-connection (lsp-stdio-connection "/path/to/lsp-server")
    :major-modes '(mylanguage-mode)
    :server-id 'lsp-server)))

(add-hook 'mylanguage-mode-hook #'lsp)
```

#### Sublime Text

Install the LSP package and add to `LSP.sublime-settings`:

```json
{
  "clients": {
    "lsp-server": {
      "enabled": true,
      "command": ["/path/to/lsp-server"],
      "selector": "source.mylanguage"
    }
  }
}
```

## Getting Started

### Quick Start Guide

1. **Build the server**: `go build -o lsp-server .`
2. **Configure your editor**: Use the integration examples above
3. **Open a file**: The server will initialize when you open a supported file
4. **Use features**: Try hovering over code, using auto-completion, or navigating to definitions

### Example LSP Communication

Here's how the LSP protocol works with example messages:

#### 1. Initialize Request (Client → Server)

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": 12345,
    "rootUri": "file:///home/user/project",
    "capabilities": {
      "textDocument": {
        "hover": {
          "contentFormat": ["markdown", "plaintext"]
        },
        "completion": {
          "completionItem": {
            "snippetSupport": true
          }
        }
      }
    }
  }
}
```

#### 2. Initialize Response (Server → Client)

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
      },
      "diagnosticProvider": true
    },
    "serverInfo": {
      "name": "lsp-server",
      "version": "1.0.0"
    }
  }
}
```

#### 3. Document Open Notification (Client → Server)

```json
{
  "jsonrpc": "2.0",
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///home/user/project/main.go",
      "languageId": "go",
      "version": 1,
      "text": "package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"
    }
  }
}
```

#### 4. Hover Request (Client → Server)

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "textDocument/hover",
  "params": {
    "textDocument": {
      "uri": "file:///home/user/project/main.go"
    },
    "position": {
      "line": 2,
      "character": 5
    }
  }
}
```

#### 5. Hover Response (Server → Client)

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "contents": {
      "kind": "markdown",
      "value": "**main**\n\nThe main function is the entry point of the program."
    }
  }
}
```

## Configuration

### Server Configuration Options

The server can be configured via initialization parameters. Common options include:

- **Root URI**: The workspace root directory
- **Initialization Options**: Custom server-specific settings
- **Client Capabilities**: Features supported by the client

### Environment Variables

- `LSP_LOG_LEVEL`: Set logging level (debug, info, warn, error)
- `LSP_LOG_FILE`: Path to log file (default: stderr)

Example:
```bash
export LSP_LOG_LEVEL=debug
export LSP_LOG_FILE=/tmp/lsp-server.log
./lsp-server
```

## Architecture

### Project Structure

```
lsp-server/
├── main.go                     # Entry point and server initialization
├── server/
│   └── server.go              # Core LSP server and request routing
├── document/
│   ├── manager.go             # Document state management
│   └── manager_test.go        # Document manager tests
├── handlers/
│   ├── handlers.go            # LSP feature implementations
│   ├── symbol_index.go        # Symbol indexing for fast lookups
│   └── *_test.go              # Handler tests
├── examples/
│   ├── README.md              # Examples documentation
│   └── sample.go              # Sample code for testing
├── go.mod                     # Go module dependencies
├── Makefile                   # Build and test automation
└── README.md                  # This file
```

### Key Components

#### Server
The core server manages JSON-RPC 2.0 communication between the client and server. It:
- Reads messages from stdin with proper header parsing
- Routes requests to appropriate handlers
- Sends responses back via stdout
- Manages server lifecycle and state

#### Document Manager
Maintains the state of all open documents with thread-safe access:
- Stores document content and version information
- Provides atomic read/write operations using `sync.RWMutex`
- Tracks document metadata (URI, language ID, version)

#### Handlers
Implements LSP features:
- **Lifecycle handlers**: Initialize, shutdown, exit
- **Text sync handlers**: didOpen, didChange, didClose, didSave
- **Feature handlers**: Hover, definition, references, completion, symbols

#### Symbol Index
Fast symbol lookup across the workspace:
- Indexes symbols from all open documents
- Provides efficient search capabilities
- Updates incrementally as documents change

### Communication Flow

```
┌─────────────┐                                    ┌─────────────┐
│   Editor    │                                    │ LSP Server  │
│  (Client)   │                                    │             │
└──────┬──────┘                                    └──────┬──────┘
       │                                                  │
       │  1. Initialize Request                          │
       │ ─────────────────────────────────────────────> │
       │                                                  │
       │  2. Initialize Response (Capabilities)          │
       │ <───────────────────────────────────────────── │
       │                                                  │
       │  3. Initialized Notification                    │
       │ ─────────────────────────────────────────────> │
       │                                                  │
       │  4. textDocument/didOpen                        │
       │ ─────────────────────────────────────────────> │
       │                                                  │
       │  5. textDocument/hover                          │
       │ ─────────────────────────────────────────────> │
       │                                                  │
       │  6. Hover Response                              │
       │ <───────────────────────────────────────────── │
       │                                                  │
       │  7. textDocument/didChange                      │
       │ ─────────────────────────────────────────────> │
       │                                                  │
       │  8. textDocument/publishDiagnostics            │
       │ <───────────────────────────────────────────── │
       │                                                  │
```

### Message Format

LSP uses a header-based message format over stdin/stdout:

```
Content-Length: <number-of-bytes>\r\n
\r\n
<json-rpc-content>
```

Example:
```
Content-Length: 124\r\n
\r\n
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":1234,"rootUri":"file:///path","capabilities":{}}}
```

## Contributing

We welcome contributions! Here's how you can help:

### Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/lsp-server.git
   cd lsp-server
   ```
3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
4. **Make your changes** and commit them:
   ```bash
   git commit -am "Add new feature"
   ```
5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
6. **Open a Pull Request** on GitHub

### Development Guidelines

#### Code Style
- Follow Go best practices and conventions
- Use `gofmt` and `golint` for code formatting
- Write clear, self-documenting code with meaningful variable names
- Add comments for complex logic or non-obvious behavior

#### Testing
- Write unit tests for new features
- Ensure all tests pass: `go test ./...`
- Aim for high test coverage
- Include both positive and negative test cases

#### Thread Safety
- Use proper synchronization primitives (`sync.Mutex`, `sync.RWMutex`)
- Avoid data races (test with `go test -race`)
- Document thread-safety guarantees

#### Error Handling
- Handle all errors explicitly
- Provide meaningful error messages
- Log errors with appropriate context
- Return proper JSON-RPC error responses

#### Documentation
- Update README.md for new features
- Add godoc comments for exported functions
- Include examples where appropriate
- Update CHANGELOG.md

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Run specific test
go test -v -run TestFunctionName ./package
```

### Building and Testing Locally

```bash
# Install dependencies
go mod download

# Build
go build -o lsp-server .

# Run
./lsp-server

# Format code
gofmt -s -w .

# Run linter
golangci-lint run
```

### Pull Request Checklist

Before submitting a PR, ensure:
- [ ] Code follows Go conventions
- [ ] All tests pass
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] Commit messages are clear and descriptive
- [ ] No unnecessary dependencies added
- [ ] Code is properly formatted

## License

This project is licensed under the MIT License. This means you are free to:
- Use the software commercially
- Modify the source code
- Distribute the software
- Use it privately

See the [LICENSE](LICENSE) file for the full license text.

## Support and Contact

### Getting Help

- **Issues**: Report bugs or request features on [GitHub Issues](https://github.com/rkkeerth/lsp-server/issues)
- **Discussions**: Ask questions in [GitHub Discussions](https://github.com/rkkeerth/lsp-server/discussions)
- **Documentation**: Check the [LSP Specification](https://microsoft.github.io/language-server-protocol/)

### Reporting Bugs

When reporting bugs, please include:
- Go version: `go version`
- Operating system and version
- Steps to reproduce the issue
- Expected vs actual behavior
- Relevant log output (from stderr)
- Minimal code example if applicable

### Feature Requests

We welcome feature requests! Please:
- Check if the feature already exists or is planned
- Describe the use case and benefits
- Provide examples of how it would work
- Consider if it aligns with the LSP specification

## Resources

### LSP Documentation
- [Language Server Protocol Specification](https://microsoft.github.io/language-server-protocol/)
- [LSP Overview](https://microsoft.github.io/language-server-protocol/overviews/lsp/overview/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)

### Go LSP Libraries
- [go.lsp.dev/protocol](https://pkg.go.dev/go.lsp.dev/protocol) - LSP protocol types
- [go.lsp.dev/jsonrpc2](https://pkg.go.dev/go.lsp.dev/jsonrpc2) - JSON-RPC implementation

### Editor Integration Guides
- [VS Code LSP Extension Guide](https://code.visualstudio.com/api/language-extensions/language-server-extension-guide)
- [Neovim LSP Documentation](https://neovim.io/doc/user/lsp.html)
- [Emacs lsp-mode](https://emacs-lsp.github.io/lsp-mode/)

## Troubleshooting

### Common Issues

#### Server Not Starting
**Problem**: Server doesn't start or exits immediately

**Solutions**:
- Verify Go installation: `go version` (requires 1.21+)
- Check build was successful: `ls -l lsp-server`
- Ensure executable permissions: `chmod +x lsp-server`
- Review stderr logs for initialization errors

#### Connection Issues
**Problem**: Editor cannot connect to server

**Solutions**:
- Verify server path in editor configuration
- Check that server is using stdin/stdout (not TCP)
- Review editor's LSP client logs
- Ensure no other process is interfering with stdin/stdout

#### Features Not Working
**Problem**: Hover, completion, or other features don't work

**Solutions**:
- Check server capabilities in initialize response
- Ensure document was opened with `textDocument/didOpen`
- Review server logs for errors
- Verify editor supports the LSP feature
- Check that the cursor position is valid

#### Performance Issues
**Problem**: Server is slow or uses too much memory

**Solutions**:
- Check for large files causing indexing delays
- Review log level (debug logging is verbose)
- Monitor goroutine count and memory usage
- Consider incremental document sync for large files

#### Diagnostic Output

```bash
# Enable debug logging
export LSP_LOG_LEVEL=debug
./lsp-server 2> debug.log

# Check server output
tail -f debug.log

# Test server manually
echo 'Content-Length: 124\r\n\r\n{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":1234}}' | ./lsp-server
```

## Roadmap

### Current Version (1.0.0)
- ✅ Full LSP lifecycle support
- ✅ Text document synchronization
- ✅ Hover information
- ✅ Go to definition
- ✅ Find references
- ✅ Document and workspace symbols
- ✅ Code completion
- ✅ Basic diagnostics

### Planned Features

#### Version 1.1
- [ ] Incremental document synchronization
- [ ] Improved diagnostics with quickfixes
- [ ] Code formatting support
- [ ] Enhanced symbol resolution

#### Version 1.2
- [ ] Signature help
- [ ] Rename refactoring
- [ ] Code actions and quick fixes
- [ ] Document links

#### Version 2.0
- [ ] Semantic highlighting
- [ ] Call hierarchy
- [ ] Type hierarchy
- [ ] Inline values
- [ ] Advanced refactoring operations

### Contributing to the Roadmap

Have ideas for new features? We'd love to hear them! Please:
1. Check existing issues and discussions
2. Open a feature request with detailed use cases
3. Consider contributing the implementation
4. Participate in design discussions

## Acknowledgments

This project is built with excellent open-source libraries:
- [go.lsp.dev](https://go.lsp.dev/) - LSP protocol types and JSON-RPC implementation
- [Zap](https://github.com/uber-go/zap) - Blazing fast structured logging
- The Go team for an excellent programming language and tooling

Special thanks to:
- Microsoft for creating and maintaining the LSP specification
- The LSP community for comprehensive documentation and examples
- All contributors and users of this project

---

**Author**: [rkkeerth](https://github.com/rkkeerth)  
**Repository**: [github.com/rkkeerth/lsp-server](https://github.com/rkkeerth/lsp-server)  
**Version**: 1.0.0  
**License**: MIT
