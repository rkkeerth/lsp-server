# LSP Server in Go

A complete Language Server Protocol (LSP) server implementation written in GoLang. This server provides fundamental language intelligence features for text editors and IDEs through the standardized LSP protocol.

## Overview

This LSP server implements the Language Server Protocol, which enables editors to provide advanced language features such as code completion, hover information, and go-to-definition functionality. The server communicates with clients (editors) using JSON-RPC 2.0 over stdin/stdout.

### Features

- **Full LSP Lifecycle Support**
  - Server initialization and capability negotiation
  - Graceful shutdown handling
  
- **Document Synchronization**
  - `textDocument/didOpen`: Track newly opened documents
  - `textDocument/didChange`: Monitor real-time content changes
  - `textDocument/didSave`: Handle document save events
  - `textDocument/didClose`: Clean up closed documents

- **Language Intelligence**
  - `textDocument/completion`: Smart code completion with keyword suggestions
  - `textDocument/hover`: Rich hover information with markdown formatting
  - `textDocument/definition`: Go-to-definition functionality for symbols

- **Robust Architecture**
  - Thread-safe document state management
  - JSON-RPC 2.0 communication layer
  - Comprehensive error handling and logging
  - Concurrent request handling

## Architecture

```
┌─────────────────────────────────────────────┐
│           Editor/IDE (Client)               │
│  (VS Code, Vim, Emacs, etc.)               │
└─────────────────┬───────────────────────────┘
                  │
                  │ JSON-RPC 2.0 over stdin/stdout
                  │
┌─────────────────▼───────────────────────────┐
│            LSP Server (main.go)             │
│  ┌──────────────────────────────────────┐  │
│  │  JSON-RPC Connection Handler         │  │
│  │  (stdio stream, message routing)     │  │
│  └────────────┬─────────────────────────┘  │
│               │                             │
│  ┌────────────▼─────────────────────────┐  │
│  │      Server Core (server.go)         │  │
│  │  - Request dispatcher                │  │
│  │  - Capability management             │  │
│  │  - Lifecycle handlers                │  │
│  └────────────┬─────────────────────────┘  │
│               │                             │
│  ┌────────────▼─────────────────────────┐  │
│  │  Protocol Handlers (handlers.go)     │  │
│  │  - Completion logic                  │  │
│  │  - Hover information                 │  │
│  │  - Definition finder                 │  │
│  └────────────┬─────────────────────────┘  │
│               │                             │
│  ┌────────────▼─────────────────────────┐  │
│  │  Document Manager (document.go)      │  │
│  │  - Thread-safe document storage      │  │
│  │  - Version tracking                  │  │
│  │  - Content synchronization           │  │
│  └──────────────────────────────────────┘  │
└─────────────────────────────────────────────┘
```

### Components

1. **main.go**: Entry point that sets up stdin/stdout communication and starts the server
2. **server.go**: Core server logic with request handling and routing
3. **handlers.go**: Implementation of LSP protocol handlers for completion, hover, and definition
4. **document.go**: Thread-safe document state management
5. **go.mod**: Go module configuration with LSP protocol dependencies

## Building

### Prerequisites

- Go 1.21 or later
- Internet connection for downloading dependencies (first build only)

### Build Steps

```bash
# Clone or navigate to the repository
cd lsp-server

# Download dependencies
go mod download

# Build the server
go build -o lsp-server .

# Or install it to your GOPATH/bin
go install
```

The build produces an executable named `lsp-server` (or `lsp-server.exe` on Windows).

## Running

The LSP server communicates over stdin/stdout and is designed to be launched by an editor. For testing purposes, you can run it manually:

```bash
./lsp-server
```

The server will log activity to stderr, which can be viewed in your editor's LSP log output.

## Editor Configuration

### Visual Studio Code

Create or update `.vscode/settings.json` in your workspace:

```json
{
  "languageServerExample.trace.server": "verbose",
  "languageServerExample.serverPath": "/path/to/lsp-server"
}
```

For a complete VS Code extension, create a minimal extension in `.vscode/extensions/lsp-example/`:

**package.json**:
```json
{
  "name": "lsp-example",
  "version": "0.1.0",
  "engines": { "vscode": "^1.75.0" },
  "activationEvents": ["onLanguage:plaintext"],
  "main": "./extension.js",
  "contributes": {
    "configuration": {
      "type": "object",
      "title": "LSP Example",
      "properties": {
        "lspExample.serverPath": {
          "type": "string",
          "default": "lsp-server",
          "description": "Path to the LSP server executable"
        }
      }
    }
  }
}
```

**extension.js**:
```javascript
const vscode = require('vscode');
const { LanguageClient } = require('vscode-languageclient/node');

let client;

function activate(context) {
    const config = vscode.workspace.getConfiguration('lspExample');
    const serverPath = config.get('serverPath') || 'lsp-server';

    const serverOptions = {
        command: serverPath,
        args: []
    };

    const clientOptions = {
        documentSelector: [{ scheme: 'file', language: 'plaintext' }]
    };

    client = new LanguageClient(
        'lspExample',
        'LSP Example',
        serverOptions,
        clientOptions
    );

    client.start();
}

function deactivate() {
    if (client) {
        return client.stop();
    }
}

module.exports = { activate, deactivate };
```

### Vim/Neovim

Using **vim-lsp**:

```vim
if executable('lsp-server')
    au User lsp_setup call lsp#register_server({
        \ 'name': 'lsp-server',
        \ 'cmd': {server_info->['lsp-server']},
        \ 'allowlist': ['text', 'markdown'],
        \ })
endif

" Enable LSP features
function! s:on_lsp_buffer_enabled() abort
    setlocal omnifunc=lsp#complete
    setlocal signcolumn=yes
    nmap <buffer> gd <plug>(lsp-definition)
    nmap <buffer> K <plug>(lsp-hover)
endfunction

augroup lsp_install
    au!
    autocmd User lsp_buffer_enabled call s:on_lsp_buffer_enabled()
augroup END
```

Using **coc.nvim**:

Add to `:CocConfig`:

```json
{
  "languageserver": {
    "example": {
      "command": "lsp-server",
      "filetypes": ["text", "markdown"],
      "rootPatterns": [".git/"]
    }
  }
}
```

Using **nvim-lspconfig** (Neovim 0.5+):

```lua
local lspconfig = require('lspconfig')
local configs = require('lspconfig.configs')

if not configs.example_lsp then
  configs.example_lsp = {
    default_config = {
      cmd = {'lsp-server'},
      filetypes = {'text', 'markdown'},
      root_dir = lspconfig.util.root_pattern('.git'),
      settings = {},
    },
  }
end

lspconfig.example_lsp.setup{}
```

### Emacs

Using **lsp-mode**:

Add to your `init.el` or `.emacs`:

```elisp
(require 'lsp-mode)

(add-to-list 'lsp-language-id-configuration '(text-mode . "text"))

(lsp-register-client
 (make-lsp-client :new-connection (lsp-stdio-connection "lsp-server")
                  :major-modes '(text-mode)
                  :server-id 'example-lsp))

(add-hook 'text-mode-hook #'lsp)
```

Using **eglot**:

```elisp
(require 'eglot)

(add-to-list 'eglot-server-programs '(text-mode . ("lsp-server")))

(add-hook 'text-mode-hook 'eglot-ensure)
```

## Development

### Project Structure

```
lsp-server/
├── main.go           # Entry point and stdio communication setup
├── server.go         # Core server and request dispatcher
├── handlers.go       # LSP protocol handler implementations
├── document.go       # Document state management
├── go.mod            # Go module dependencies
└── README.md         # This file
```

### Adding New Features

To add a new LSP method:

1. Add a case in `server.go`'s `Handle()` method
2. Implement the handler function following the pattern:
   ```go
   func (s *Server) handleNewMethod(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error
   ```
3. Update server capabilities in `handleInitialize()` if needed
4. Implement the actual logic in `handlers.go` or a new file

### Testing

Test the server manually with an editor or use a test client:

```bash
# Run with verbose logging
./lsp-server 2> lsp-server.log

# In another terminal, send test requests
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}' | ./lsp-server
```

## Logging

The server logs all activity to stderr. To capture logs:

```bash
./lsp-server 2> server.log
```

When configured in an editor, logs typically appear in the editor's LSP log panel.

## Extending the Server

This implementation provides a foundation for building language-specific LSP servers. To extend it:

1. **Language-Specific Parsing**: Replace the simple pattern matching in `handlers.go` with a proper parser for your target language
2. **Symbol Table**: Implement a symbol table to track variable/function definitions across files
3. **Type System**: Add type inference and checking for statically-typed languages
4. **Diagnostics**: Implement `textDocument/publishDiagnostics` to report errors and warnings
5. **Code Actions**: Add `textDocument/codeAction` for refactoring and quick fixes
6. **Formatting**: Implement `textDocument/formatting` for code formatting
7. **References**: Add `textDocument/references` to find all symbol references
8. **Rename**: Implement `textDocument/rename` for symbol renaming

## Troubleshooting

**Server doesn't start in editor:**
- Check the server path is correct in editor configuration
- Verify the binary has execute permissions: `chmod +x lsp-server`
- Review editor's LSP logs for error messages

**No completions/hover/definitions:**
- Ensure the file type is correctly configured in editor settings
- Check that documents are being opened (review server logs)
- Verify the server capabilities are correctly reported during initialization

**Performance issues:**
- Consider implementing incremental document synchronization instead of full text sync
- Add caching for parsed syntax trees
- Optimize symbol lookup with proper indexing

## License

This is an example implementation for educational purposes.

## Contributing

To contribute improvements:

1. Add comprehensive error handling for edge cases
2. Implement additional LSP methods (formatting, code actions, etc.)
3. Add unit tests for document management and handlers
4. Improve language intelligence with real parsing
5. Add support for multi-file projects and workspaces

## References

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [go.lsp.dev Documentation](https://pkg.go.dev/go.lsp.dev)
