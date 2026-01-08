# LSP Server

A complete Language Server Protocol (LSP) server implementation in GoLang that follows the LSP specification and provides intelligent language features for code editors and IDEs.

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Development Setup](#development-setup)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)
- [Contact & Support](#contact--support)

## Overview

The LSP Server is a comprehensive implementation of the Language Server Protocol, designed to provide intelligent language features to any editor or IDE that supports LSP. Built with GoLang, this server emphasizes performance, reliability, and extensibility.

### What is LSP?

The Language Server Protocol (LSP) is an open standard developed by Microsoft that defines the protocol used between editors/IDEs and language servers. It enables features like:

- Code completion
- Go to definition
- Find references
- Hover information
- Document symbols
- And many more...

This project implements a fully functional LSP server that can be integrated with popular editors like VS Code, Neovim, Emacs, and others.

## Key Features

### Lifecycle Management
- **Initialize**: Negotiates capabilities between client and server
- **Initialized**: Confirms successful initialization
- **Shutdown**: Gracefully shuts down the server
- **Exit**: Terminates the server process

### Text Document Synchronization
- **didOpen**: Tracks when documents are opened in the editor
- **didChange**: Updates document content on changes (full synchronization)
- **didClose**: Removes documents from tracking when closed
- **didSave**: Handles document save events

### Language Features
- **Hover**: Displays information about symbols when hovering over them
- **Go to Definition**: Navigates to symbol definitions
- **Find References**: Locates all usages of a symbol across documents
- **Document Symbols**: Provides document outline and navigation
- **Workspace Symbols**: Searches for symbols across the entire workspace
- **Code Completion**: Offers intelligent code suggestions including keywords and symbols
- **Diagnostics**: Reports errors, warnings, and hints (TODO/FIXME detection)

### Architecture Highlights
- **Thread-Safe**: Safe concurrent request handling with mutex protection
- **Structured Logging**: Comprehensive logging using zap
- **JSON-RPC 2.0**: Standard protocol compliance
- **Extensible**: Easy to add new features and capabilities

## Installation

### Prerequisites

- **Go**: Version 1.21 or higher
- **Git**: For cloning the repository

### Building from Source

```bash
# Clone the repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Download dependencies
go mod download

# Build the server
go build -o lsp-server .

# Verify the build
./lsp-server --version
```

### Install with Go

```bash
go install github.com/rkkeerth/lsp-server@latest
```

## Usage

### Running the Server

The LSP server communicates via stdin/stdout, which is the standard for LSP implementations:

```bash
./lsp-server
```

### Integration with Editors

#### VS Code

Create or update your VS Code extension's `package.json`:

```json
{
  "contributes": {
    "configuration": {
      "type": "object",
      "title": "LSP Server Configuration",
      "properties": {
        "lspServer.executablePath": {
          "type": "string",
          "default": "lsp-server",
          "description": "Path to the LSP server executable"
        }
      }
    }
  }
}
```

Or use in a client extension:

```typescript
import { LanguageClient, LanguageClientOptions, ServerOptions } from 'vscode-languageclient/node';

const serverOptions: ServerOptions = {
  command: 'lsp-server',
  args: []
};

const clientOptions: LanguageClientOptions = {
  documentSelector: [{ scheme: 'file', language: 'yourlanguage' }]
};

const client = new LanguageClient('lspServer', 'LSP Server', serverOptions, clientOptions);
client.start();
```

#### Neovim

Add to your Neovim configuration (init.lua):

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

-- Define the custom LSP server
if not configs.lsp_server then
  configs.lsp_server = {
    default_config = {
      cmd = {'/path/to/lsp-server'},
      filetypes = {'go', 'python', 'javascript'}, -- Adjust as needed
      root_dir = lspconfig.util.root_pattern('.git', 'go.mod'),
      settings = {},
    },
  }
end

-- Start the LSP server
lspconfig.lsp_server.setup{
  on_attach = function(client, bufnr)
    -- Your keybindings here
    local bufopts = { noremap=true, silent=true, buffer=bufnr }
    vim.keymap.set('n', 'gd', vim.lsp.buf.definition, bufopts)
    vim.keymap.set('n', 'K', vim.lsp.buf.hover, bufopts)
    vim.keymap.set('n', 'gr', vim.lsp.buf.references, bufopts)
  end,
}
```

#### Emacs (lsp-mode)

```elisp
(with-eval-after-load 'lsp-mode
  (add-to-list 'lsp-language-id-configuration '(your-mode . "language-id"))
  (lsp-register-client
   (make-lsp-client
    :new-connection (lsp-stdio-connection "/path/to/lsp-server")
    :major-modes '(your-mode)
    :server-id 'lsp-server)))
```

#### Sublime Text

Add to your LSP settings:

```json
{
  "clients": {
    "lsp-server": {
      "enabled": true,
      "command": ["/path/to/lsp-server"],
      "selector": "source.yourlanguage"
    }
  }
}
```

## Configuration

### Environment Variables

The LSP server can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `LSP_LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
| `LSP_LOG_FILE` | Path to log file | `stderr` |
| `LSP_MAX_WORKERS` | Maximum concurrent workers | `10` |

### Runtime Configuration

The server can be configured through initialization parameters sent by the client during the `initialize` request:

```json
{
  "initializationOptions": {
    "logLevel": "debug",
    "diagnosticsEnabled": true,
    "completionEnabled": true
  }
}
```

## Development Setup

### For Contributors

1. **Fork and Clone**

```bash
git clone https://github.com/YOUR_USERNAME/lsp-server.git
cd lsp-server
```

2. **Install Dependencies**

```bash
go mod download
```

3. **Run Tests**

```bash
go test ./...
```

4. **Run with Debug Logging**

```bash
./lsp-server 2> lsp-server.log
```

### Project Structure

```
lsp-server/
├── main.go                 # Entry point and server initialization
├── server/
│   └── server.go          # Core LSP server and request routing
├── document/
│   └── manager.go         # Document management and text operations
├── handlers/
│   ├── handlers.go        # LSP feature implementations
│   └── symbol_index.go    # Symbol indexing for fast lookups
├── go.mod                 # Go module dependencies
└── README.md              # This file
```

### Key Components

- **Server**: Manages JSON-RPC communication, handles client requests, and routes them to appropriate handlers
- **Document Manager**: Maintains the state of all open documents with thread-safe access
- **Handlers**: Implements LSP features like hover, completion, definition, and references
- **Symbol Index**: Maintains a fast lookup index for workspace symbols

### Adding New Features

To add a new LSP feature:

1. Add the capability to `server/server.go` in `setupCapabilities()`
2. Add a new handler method in `handlers/handlers.go`
3. Add the route in `server/server.go` in the `handle()` method
4. Add tests for the new feature
5. Update this README with the new feature

Example:

```go
// In handlers/handlers.go
func (h *Handler) NewFeature(uri protocol.DocumentURI, params YourParams) *YourResult {
    // Implementation here
}

// In server/server.go
case "textDocument/newFeature":
    // Handle the request
```

## API Documentation

### Protocol Implementation

This server implements the Language Server Protocol using:

- **JSON-RPC 2.0**: Standard protocol for client-server communication
- **stdin/stdout**: Communication channel as per LSP specification
- **Header-based messaging**: Each message includes `Content-Length` header

### Message Format

```
Content-Length: <bytes>\r\n
\r\n
{JSON-RPC message}
```

### Example Messages

#### Initialize Request

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": 12345,
    "rootUri": "file:///path/to/project",
    "capabilities": {}
  }
}
```

#### Initialize Response

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
      }
    },
    "serverInfo": {
      "name": "lsp-server",
      "version": "1.0.0"
    }
  }
}
```

#### Hover Request

```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "textDocument/hover",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.go"
    },
    "position": {
      "line": 2,
      "character": 5
    }
  }
}
```

### Supported Methods

| Method | Description | Status |
|--------|-------------|--------|
| `initialize` | Initialize the server | ✅ Implemented |
| `initialized` | Initialization complete | ✅ Implemented |
| `shutdown` | Shutdown the server | ✅ Implemented |
| `exit` | Exit the server | ✅ Implemented |
| `textDocument/didOpen` | Document opened | ✅ Implemented |
| `textDocument/didChange` | Document changed | ✅ Implemented |
| `textDocument/didSave` | Document saved | ✅ Implemented |
| `textDocument/didClose` | Document closed | ✅ Implemented |
| `textDocument/hover` | Hover information | ✅ Implemented |
| `textDocument/definition` | Go to definition | ✅ Implemented |
| `textDocument/references` | Find references | ✅ Implemented |
| `textDocument/documentSymbol` | Document symbols | ✅ Implemented |
| `workspace/symbol` | Workspace symbols | ✅ Implemented |
| `textDocument/completion` | Code completion | ✅ Implemented |
| `textDocument/publishDiagnostics` | Diagnostics | ✅ Implemented |

## Contributing

We welcome contributions from the community! Here's how you can help:

### Ways to Contribute

- **Bug Reports**: Report bugs through GitHub Issues
- **Feature Requests**: Suggest new features or improvements
- **Code Contributions**: Submit pull requests with bug fixes or new features
- **Documentation**: Improve documentation and examples
- **Testing**: Write tests to improve code coverage

### Contribution Guidelines

1. **Fork** the repository
2. **Create a branch** for your feature (`git checkout -b feature/amazing-feature`)
3. **Write tests** for your changes
4. **Ensure tests pass** (`go test ./...`)
5. **Format your code** (`go fmt ./...`)
6. **Commit your changes** (`git commit -m 'Add some amazing feature'`)
7. **Push to the branch** (`git push origin feature/amazing-feature`)
8. **Open a Pull Request**

### Code Standards

- Follow Go conventions and best practices
- Write clear, concise commit messages
- Include tests for new features
- Maintain thread safety in concurrent code
- Document exported functions and types
- Keep error handling comprehensive

### Pull Request Process

1. Update the README.md with details of changes if applicable
2. Update the documentation with any new features or changes
3. The PR will be merged once you have the sign-off of the maintainers

## License

This project is licensed under the MIT License - see below for details:

```
MIT License

Copyright (c) 2025 rkkeerth

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## Contact & Support

### Maintainer

- **GitHub**: [@rkkeerth](https://github.com/rkkeerth)
- **Repository**: [lsp-server](https://github.com/rkkeerth/lsp-server)

### Getting Help

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/rkkeerth/lsp-server/issues)
- **Discussions**: Join discussions in [GitHub Discussions](https://github.com/rkkeerth/lsp-server/discussions)
- **Pull Requests**: Contribute code via [Pull Requests](https://github.com/rkkeerth/lsp-server/pulls)

### Resources

- [Language Server Protocol Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [LSP Go Protocol Types](https://pkg.go.dev/go.lsp.dev/protocol)
- [Go Documentation](https://golang.org/doc/)

### Acknowledgments

Special thanks to:
- The Language Server Protocol team at Microsoft
- The Go LSP library maintainers
- All contributors to this project

---

**Built with ❤️ by the community**
