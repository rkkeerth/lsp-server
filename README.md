# Basic LSP Server

[![Python Version](https://img.shields.io/badge/python-3.8+-blue.svg)](https://www.python.org/downloads/)
[![LSP Version](https://img.shields.io/badge/LSP-3.17-green.svg)](https://microsoft.github.io/language-server-protocol/)
[![No Dependencies](https://img.shields.io/badge/dependencies-none-brightgreen.svg)](requirements.txt)

A lightweight, educational Language Server Protocol (LSP) implementation in Python. Built entirely with Python's standard library to demonstrate LSP fundamentals without external dependencies.

## 🎯 Project Purpose

This project serves as:

- **Learning Resource**: Understand LSP protocol internals by reading clean, well-documented code
- **Reference Implementation**: See how JSON-RPC 2.0 message handling works in practice
- **Starter Template**: Fork and extend to create custom language servers for your own needs
- **Proof of Concept**: Demonstrate that LSP servers don't require heavy frameworks

Perfect for developers who want to understand how language servers work under the hood or need a minimal foundation to build upon.

## ✨ Features

### Core LSP Capabilities

- **🔄 Lifecycle Management**: Full initialize/shutdown protocol compliance
  - Capability negotiation with clients
  - Graceful server shutdown
  - Proper state management

- **📄 Text Document Synchronization**: Real-time document tracking
  - `textDocument/didOpen` - Track newly opened documents
  - `textDocument/didChange` - Sync content changes (full document mode)
  - `textDocument/didClose` - Clean up closed documents

- **🔍 Diagnostic Analysis**: Intelligent code quality checks
  - `TODO` markers → Informational hints
  - `FIXME` markers → Warning alerts
  - Lines exceeding 120 characters → Style warnings
  - Duplicate consecutive lines → Code smell detection

- **📡 JSON-RPC 2.0 Protocol**: Standards-compliant message handling
  - Request/response correlation with message IDs
  - Notification support (fire-and-forget)
  - Proper error handling and reporting
  - Content-Length based message framing

## 📋 Requirements

**Minimum Requirements:**
- Python 3.8 or higher
- Standard library only (no pip packages needed!)
- Any OS: Linux, macOS, Windows

**For Editor Integration:**
- An LSP-compatible editor (VS Code, Neovim, Emacs, Sublime Text, Vim, etc.)
- Optional: Terminal for testing with the included test client

## 🚀 Quick Start

### Installation

**Option 1: Clone the Repository**
```bash
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server
```

**Option 2: Direct Download**
```bash
# Download just the server file
curl -O https://raw.githubusercontent.com/rkkeerth/lsp-server/main/server.py
chmod +x server.py
```

No pip install needed! The server is ready to run immediately:

```bash
python3 server.py
```

### 5-Minute Test Drive

Verify everything works in seconds:

```bash
# Run the included test client with example file
python3 test_client.py examples/test.txt
```

Expected output:
```
============================================================
Diagnostics received for: file:///path/to/test.txt
============================================================
Found 7 issue(s):

  [INFO] Line 5, Col 1: TODO found: Consider addressing this item
  [WARNING] Line 8, Col 1: FIXME found: This requires immediate attention
  [WARNING] Line 10, Col 121: Line too long (145 > 120 characters)
  ...
============================================================
```

✅ Success! Your LSP server is working correctly.

## 📖 Usage

### Running the Server Standalone

The LSP server communicates via standard input/output (stdin/stdout):

```bash
python3 server.py
```

The server will:
1. Start listening on stdin for JSON-RPC messages
2. Log activity to `/tmp/lsp-server.log`
3. Send responses to stdout
4. Run until receiving a shutdown request

### Testing with the Included Client

A full-featured test client demonstrates all server capabilities:

```bash
# Test with the example file
python3 test_client.py examples/test.txt

# Test with your own file
python3 test_client.py /path/to/your/file.txt

# Use the automated test script
./run_test.sh
```

The test client will:
1. ✅ Start the LSP server as a subprocess
2. ✅ Initialize the client-server connection
3. ✅ Open the specified document
4. ✅ Display diagnostics with severity levels and line numbers
5. ✅ Properly shut down the server

### Editor Integration

#### Visual Studio Code

**Method 1: Using a Generic LSP Extension**

Install an extension like ["Generic LSP Client"](https://marketplace.visualstudio.com/items?itemName=GregorBiswanger.genericlsp) and configure:

```json
{
  "genericlsp.languageServers": [
    {
      "command": "python3",
      "args": ["/absolute/path/to/lsp-server/server.py"],
      "filetypes": ["txt", "text"],
      "name": "basic-lsp-server"
    }
  ]
}
```

**Method 2: Create Custom Extension**

1. Use the [VS Code Extension Generator](https://code.visualstudio.com/api/get-started/your-first-extension)
2. Add LSP client dependency
3. Configure server command in extension code

#### Neovim (Native LSP)

Add to your `init.lua`:

```lua
-- Configure basic LSP server
vim.api.nvim_create_autocmd("FileType", {
  pattern = {"text", "txt"},
  callback = function()
    vim.lsp.start({
      name = 'basic-lsp-server',
      cmd = {'python3', '/absolute/path/to/lsp-server/server.py'},
      root_dir = vim.fs.dirname(vim.fs.find({'*.txt'}, { upward = true })[1]),
      settings = {},
    })
  end,
})
```

Or in `init.vim`:

```vim
lua << EOF
vim.lsp.start({
  name = 'basic-lsp-server',
  cmd = {'python3', '/absolute/path/to/lsp-server/server.py'},
  root_dir = vim.fn.getcwd(),
})
EOF
```

#### Emacs (lsp-mode)

Add to your Emacs configuration:

```elisp
(require 'lsp-mode)

(lsp-register-client
 (make-lsp-client 
   :new-connection (lsp-stdio-connection 
                     '("python3" "/absolute/path/to/lsp-server/server.py"))
   :major-modes '(text-mode)
   :server-id 'basic-lsp-server
   :priority -1))

(add-hook 'text-mode-hook #'lsp)
```

#### Sublime Text (LSP Package)

1. Install the [LSP package](https://packagecontrol.io/packages/LSP)
2. Add to LSP settings (`Preferences > Package Settings > LSP > Settings`):

```json
{
  "clients": {
    "basic-lsp-server": {
      "enabled": true,
      "command": ["python3", "/absolute/path/to/lsp-server/server.py"],
      "selector": "text.plain"
    }
  }
}
```

#### Vim (with vim-lsp)

```vim
if executable('python3')
  au User lsp_setup call lsp#register_server({
    \ 'name': 'basic-lsp-server',
    \ 'cmd': {server_info->['python3', '/absolute/path/to/lsp-server/server.py']},
    \ 'allowlist': ['text'],
    \ })
endif
```

## 🏗️ Architecture & Technical Details

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     LSP Client (Editor)                     │
│  (VS Code, Neovim, Emacs, Sublime Text, etc.)              │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ JSON-RPC 2.0 over stdin/stdout
                         │
┌────────────────────────▼────────────────────────────────────┐
│                    LSPServer (server.py)                    │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Message Handler                                    │    │
│  │  • Request routing (initialize, shutdown)          │    │
│  │  • Notification routing (didOpen, didChange, etc.) │    │
│  │  • Response generation                             │    │
│  └────────────────────────────────────────────────────┘    │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Document Manager                                   │    │
│  │  • In-memory document storage                      │    │
│  │  • Content synchronization                         │    │
│  └────────────────────────────────────────────────────┘    │
│  ┌────────────────────────────────────────────────────┐    │
│  │  Analysis Engine                                    │    │
│  │  • Pattern matching (TODO, FIXME)                  │    │
│  │  • Style checking (line length)                    │    │
│  │  • Code smell detection (duplicates)               │    │
│  └────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

### JSON-RPC 2.0 Communication Protocol

Messages are exchanged using the Language Server Protocol format:

**Message Structure:**
```
Content-Length: 123\r\n
\r\n
{JSON-RPC message}
```

**Request Example:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": null,
    "rootUri": "file:///path/to/workspace",
    "capabilities": {}
  }
}
```

**Response Example:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "capabilities": {
      "textDocumentSync": 1,
      "diagnosticProvider": {}
    }
  }
}
```

**Notification Example:**
```json
{
  "jsonrpc": "2.0",
  "method": "textDocument/didOpen",
  "params": {
    "textDocument": {
      "uri": "file:///path/to/file.txt",
      "languageId": "plaintext",
      "version": 1,
      "text": "file contents..."
    }
  }
}
```

### Message Flow Sequence

```
Client                                    Server
  │                                         │
  ├─────── initialize (request) ──────────>│
  │                                         │ (negotiate capabilities)
  │<────── initialize response ─────────────┤
  │                                         │
  ├─────── initialized (notification) ─────>│
  │                                         │
  ├─────── textDocument/didOpen ───────────>│
  │                                         │ (analyze document)
  │<────── publishDiagnostics ──────────────┤
  │                                         │
  ├─────── textDocument/didChange ─────────>│
  │                                         │ (re-analyze document)
  │<────── publishDiagnostics ──────────────┤
  │                                         │
  ├─────── shutdown (request) ─────────────>│
  │                                         │
  │<────── shutdown response ────────────────┤
  │                                         │
  ├─────── exit (notification) ─────────────>│
  │                                         │ (terminate)
```

### Core Components

#### 1. **LSPServer Class** (`server.py`)

Main server implementation with the following responsibilities:

- **Message I/O**: Read/write JSON-RPC messages via stdin/stdout
- **Request Handling**: Process requests that expect responses
- **Notification Handling**: Process fire-and-forget notifications
- **Document Lifecycle**: Track open documents and their content
- **Diagnostic Publishing**: Push analysis results to clients

**Key Methods:**
```python
def handle_initialize(params)       # Capability negotiation
def handle_shutdown(params)         # Prepare for termination
def handle_text_document_did_open(params)   # Document opened
def handle_text_document_did_change(params) # Document edited
def handle_text_document_did_close(params)  # Document closed
def analyze_document(text)          # Core analysis logic
def publish_diagnostics(uri, diags) # Send results to client
```

#### 2. **Analysis Engine**

The `analyze_document()` method scans text for issues using:

- **Regular expressions**: Pattern matching for TODO/FIXME
- **String operations**: Length checking, line comparison
- **Position tracking**: LSP-compatible line/character positions

**Diagnostic Structure:**
```python
{
  "range": {
    "start": {"line": 5, "character": 0},  # 0-indexed
    "end": {"line": 5, "character": 4}
  },
  "message": "TODO found: Consider addressing this item",
  "severity": 3,  # 1=Error, 2=Warning, 3=Info, 4=Hint
  "source": "basic-lsp-server"
}
```

#### 3. **Document State Management**

Documents are stored in memory with URI-to-content mapping:

```python
self.documents: Dict[str, str] = {
  "file:///path/to/file1.txt": "file1 contents...",
  "file:///path/to/file2.txt": "file2 contents...",
}
```

**Synchronization modes:**
- **Full (mode=1)**: Client sends entire document on each change
- **Incremental (mode=2)**: Client sends only changed portions (not implemented)

This server uses **full synchronization** for simplicity.

## 📊 LSP Specification Compliance

This implementation follows the [Language Server Protocol Specification v3.17](https://microsoft.github.io/language-server-protocol/) and implements the following subset:

### Lifecycle Messages

| Message | Type | Status | Description |
|---------|------|--------|-------------|
| `initialize` | Request | ✅ Implemented | Negotiate capabilities between client and server |
| `initialized` | Notification | ✅ Implemented | Confirm initialization complete |
| `shutdown` | Request | ✅ Implemented | Prepare for graceful shutdown |
| `exit` | Notification | ✅ Implemented | Terminate server process |

### Document Synchronization

| Message | Type | Status | Description |
|---------|------|--------|-------------|
| `textDocument/didOpen` | Notification | ✅ Implemented | Document opened in editor |
| `textDocument/didChange` | Notification | ✅ Implemented | Document content changed |
| `textDocument/didClose` | Notification | ✅ Implemented | Document closed in editor |
| `textDocument/didSave` | Notification | ⚠️ Received but no-op | Document saved (no special handling) |

### Diagnostics

| Message | Type | Status | Description |
|---------|------|--------|-------------|
| `textDocument/publishDiagnostics` | Notification | ✅ Implemented | Server pushes diagnostics to client |

### Server Capabilities

Advertised in the `initialize` response:

```json
{
  "capabilities": {
    "textDocumentSync": {
      "openClose": true,
      "change": 1,
      "save": {
        "includeText": false
      }
    },
    "diagnosticProvider": {
      "interFileDependencies": false,
      "workspaceDiagnostics": false
    }
  }
}
```

**Text Document Sync Mode:**
- `0` = None
- `1` = Full (entire document sent on change) ← **We use this**
- `2` = Incremental (only changes sent)

### Not Yet Implemented

These are common LSP features not included in this basic implementation:

| Feature | Message |
|---------|---------|
| Code Completion | `textDocument/completion` |
| Hover Information | `textDocument/hover` |
| Go to Definition | `textDocument/definition` |
| Find References | `textDocument/references` |
| Document Symbols | `textDocument/documentSymbol` |
| Code Formatting | `textDocument/formatting` |
| Code Actions | `textDocument/codeAction` |
| Rename | `textDocument/rename` |
| Workspace Symbols | `workspace/symbol` |

## 🛠️ Development Guide

### Project Structure

```
lsp-server/
├── server.py           # Main LSP server implementation
├── test_client.py      # Test client for validation
├── requirements.txt    # Empty (no dependencies!)
├── pyproject.toml      # Python project metadata
├── setup.py            # Setup script
├── run_test.sh         # Automated test runner
├── README.md           # This file
├── QUICKSTART.md       # Quick start guide
├── examples/
│   ├── test.txt       # Sample file with various issues
│   └── code_sample.py # Python code example
└── .gitignore

Total: ~600 lines of Python code
```

### Adding Custom Diagnostics

Extend the analysis engine by modifying `analyze_document()` in `server.py`:

**Example 1: Detect Hardcoded Credentials**
```python
# Add to analyze_document() method
match = re.search(r'password\s*=\s*["\'][\w]+["\']', line, re.IGNORECASE)
if match:
    diagnostics.append({
        "range": {
            "start": {"line": line_num, "character": match.start()},
            "end": {"line": line_num, "character": match.end()}
        },
        "message": "Hardcoded password detected - use environment variables",
        "severity": 1,  # Error
        "source": "basic-lsp-server"
    })
```

**Example 2: Detect Trailing Whitespace**
```python
if line.endswith(' ') or line.endswith('\t'):
    diagnostics.append({
        "range": {
            "start": {"line": line_num, "character": len(line.rstrip())},
            "end": {"line": line_num, "character": len(line)}
        },
        "message": "Trailing whitespace",
        "severity": 3,  # Info
        "source": "basic-lsp-server"
    })
```

**Example 3: Detect Missing Documentation**
```python
# Check for function definitions without docstrings
if re.match(r'^\s*def\s+\w+\(', line):
    # Check if next line is a docstring
    if line_num + 1 < len(lines):
        next_line = lines[line_num + 1].strip()
        if not next_line.startswith('"""') and not next_line.startswith("'''"):
            diagnostics.append({
                "range": {
                    "start": {"line": line_num, "character": 0},
                    "end": {"line": line_num, "character": len(line)}
                },
                "message": "Function missing docstring",
                "severity": 3,  # Info
                "source": "basic-lsp-server"
            })
```

### Diagnostic Severity Levels

| Severity | Value | Color (typical) | Use Case |
|----------|-------|-----------------|----------|
| Error | `1` | 🔴 Red | Syntax errors, critical issues |
| Warning | `2` | 🟡 Yellow | Style violations, code smells |
| Information | `3` | 🔵 Blue | Suggestions, TODOs |
| Hint | `4` | ⚪ Gray | Minor improvements, optimizations |

### Adding New LSP Features

**Step 1: Add Handler Method**
```python
def handle_text_document_hover(self, params: Dict[str, Any]) -> Dict[str, Any]:
    """Handle textDocument/hover request."""
    uri = params["textDocument"]["uri"]
    position = params["position"]
    line_num = position["line"]
    
    # Get document content
    if uri not in self.documents:
        return {"contents": ""}
    
    text = self.documents[uri]
    lines = text.split('\n')
    
    if line_num < len(lines):
        line = lines[line_num]
        return {
            "contents": {
                "kind": "markdown",
                "value": f"**Line {line_num + 1}**: `{line.strip()}`"
            }
        }
    
    return {"contents": ""}
```

**Step 2: Register Handler**
```python
def handle_request(self, request_id: Any, method: str, params: Dict[str, Any]) -> None:
    handlers = {
        "initialize": self.handle_initialize,
        "shutdown": self.handle_shutdown,
        "textDocument/hover": self.handle_text_document_hover,  # ← Add here
    }
    # ... rest of method
```

**Step 3: Update Server Capabilities**
```python
def handle_initialize(self, params: Dict[str, Any]) -> Dict[str, Any]:
    return {
        "capabilities": {
            "textDocumentSync": {...},
            "diagnosticProvider": {...},
            "hoverProvider": True,  # ← Add capability
        },
        # ...
    }
```

### Testing Your Changes

**1. Unit Test with Test Client**
```bash
# Create a test file
echo "TODO: Test hover functionality" > test_hover.txt

# Run test client
python3 test_client.py test_hover.txt

# Check logs for debugging
tail -f /tmp/lsp-server.log
```

**2. Integration Test with Editor**
```bash
# Start Neovim with custom config
nvim -u custom_lsp_config.lua test_file.txt

# Or VS Code
code --disable-extensions --add "path/to/lsp-extension"
```

**3. Manual Protocol Testing**
```bash
# Send raw JSON-RPC messages
echo -e 'Content-Length: 52\r\n\r\n{"jsonrpc":"2.0","id":1,"method":"initialize"}' | python3 server.py
```

### Logging and Debugging

**Log Location:** `/tmp/lsp-server.log`

**Monitor logs in real-time:**
```bash
tail -f /tmp/lsp-server.log
```

**Log levels:**
- `INFO`: Normal operation messages
- `WARNING`: Unexpected but handled situations
- `ERROR`: Exceptions and failures

**Add custom logging:**
```python
import logging
logger = logging.getLogger(__name__)

logger.info(f"Processing document: {uri}")
logger.warning(f"Unknown method: {method}")
logger.error(f"Failed to parse: {e}", exc_info=True)
```

### Performance Considerations

**Current Implementation:**
- ✅ Single-threaded: Simple, predictable
- ✅ Synchronous I/O: Easy to debug
- ❌ Blocks on slow operations: Not suitable for large files
- ❌ Full document sync: Inefficient for large files

**Production Improvements:**

1. **Async I/O**
```python
import asyncio

async def read_message(self):
    # Non-blocking message reading
    pass

async def handle_request(self, ...):
    # Concurrent request handling
    pass
```

2. **Incremental Sync**
```python
def handle_text_document_did_change(self, params):
    # Apply only changed portions
    for change in content_changes:
        start = change["range"]["start"]
        end = change["range"]["end"]
        # Update only affected range
```

3. **Background Analysis**
```python
import threading

def analyze_in_background(self, uri, text):
    thread = threading.Thread(
        target=self._analyze_and_publish,
        args=(uri, text)
    )
    thread.start()
```

4. **Caching**
```python
from functools import lru_cache

@lru_cache(maxsize=100)
def analyze_document(self, text: str):
    # Cache results for unchanged documents
    pass
```

## 🐛 Troubleshooting

### Common Issues and Solutions

#### Server Not Responding

**Symptoms:**
- Editor shows "LSP server disconnected"
- No diagnostics appear
- Hanging on initialization

**Solutions:**
```bash
# 1. Check if server is running
ps aux | grep server.py

# 2. Check logs for errors
tail -f /tmp/lsp-server.log

# 3. Verify Python version
python3 --version  # Must be 3.8+

# 4. Test server manually
echo 'Content-Length: 52\r\n\r\n{"jsonrpc":"2.0","id":1,"method":"initialize"}' | python3 server.py

# 5. Restart editor and try again
```

#### Diagnostics Not Appearing

**Symptoms:**
- Server connects but no issues shown
- File opens but analysis doesn't run

**Solutions:**
```bash
# 1. Verify file was opened correctly
grep "Document opened" /tmp/lsp-server.log

# 2. Check file URI format
# Correct: file:///absolute/path/to/file.txt
# Wrong: /relative/path/file.txt

# 3. Ensure file has detectable issues
echo "TODO: test issue" > test.txt
python3 test_client.py test.txt

# 4. Check editor LSP client configuration
# Some editors require explicit diagnostics request
```

#### Connection Refused / Initialization Failed

**Symptoms:**
- "Failed to start language server"
- "Connection closed by server"

**Solutions:**
```python
# 1. Check for port conflicts (if using TCP)
# Not applicable for this stdio-based server

# 2. Verify server path in editor config
# Must be absolute path: /absolute/path/to/server.py

# 3. Check file permissions
chmod +x server.py

# 4. Test with minimal client
python3 test_client.py examples/test.txt
```

#### Performance Issues

**Symptoms:**
- Slow typing response
- Editor freezing
- High CPU usage

**Solutions:**
```python
# 1. Check file size
ls -lh your-file.txt
# Files > 1MB may cause slowness

# 2. Reduce analysis frequency
# Modify handle_text_document_did_change to debounce:

import time

class LSPServer:
    def __init__(self):
        self.last_analysis_time = {}
        self.debounce_delay = 0.5  # seconds
    
    def handle_text_document_did_change(self, params):
        uri = params["textDocument"]["uri"]
        current_time = time.time()
        
        # Skip if analyzed recently
        if uri in self.last_analysis_time:
            if current_time - self.last_analysis_time[uri] < self.debounce_delay:
                return
        
        # ... rest of method
        self.last_analysis_time[uri] = current_time

# 3. Disable server for large files temporarily
```

#### Invalid JSON-RPC Messages

**Symptoms:**
- "Error parsing message"
- Server crashes on startup

**Solutions:**
```bash
# 1. Check message format
# Must have Content-Length header and proper JSON

# 2. Verify editor sends correct format
# Enable verbose logging in editor

# 3. Test with known-good client
python3 test_client.py examples/test.txt

# 4. Check for encoding issues
file -i server.py  # Should be utf-8
```

#### Editor-Specific Issues

**VS Code:**
```json
// Check output panel: View → Output → Select "LSP server"
// Enable trace logging:
{
  "lsp.trace.server": "verbose"
}
```

**Neovim:**
```lua
-- Check LSP logs
:lua vim.cmd('e ' .. vim.lsp.get_log_path())

-- Restart LSP client
:LspRestart

-- Check client status
:LspInfo
```

**Emacs:**
```elisp
;; Check *lsp-log* buffer
(switch-to-buffer "*lsp-log*")

;; Restart LSP
M-x lsp-workspace-restart

;; Enable debug mode
(setq lsp-log-io t)
```

### Debug Checklist

Before reporting issues, verify:

- [ ] Python 3.8+ installed: `python3 --version`
- [ ] Server file is executable: `chmod +x server.py`
- [ ] Server runs manually: `python3 server.py`
- [ ] Test client works: `python3 test_client.py examples/test.txt`
- [ ] Log file exists: `ls -l /tmp/lsp-server.log`
- [ ] No errors in logs: `tail -20 /tmp/lsp-server.log`
- [ ] Editor LSP client installed and configured
- [ ] Absolute paths used in editor configuration
- [ ] File URI format is correct

### Getting Help

If issues persist:

1. **Collect diagnostic information:**
```bash
# System info
python3 --version
uname -a

# Server logs
tail -50 /tmp/lsp-server.log

# Test client output
python3 test_client.py examples/test.txt 2>&1 | tee debug.log
```

2. **Create minimal reproduction:**
```bash
# Smallest file that reproduces the issue
echo "TODO: test" > minimal.txt
python3 test_client.py minimal.txt
```

3. **Check GitHub issues**: Search for similar problems
4. **Provide context**: Editor, OS, Python version, exact error messages

## 🤝 Contributing

Contributions are welcome! This project aims to remain simple and educational.

### Contribution Guidelines

**Good Contributions:**
- Bug fixes with test cases
- Documentation improvements
- New diagnostic rules with examples
- Editor integration guides
- Performance optimizations

**Please Avoid:**
- Adding external dependencies
- Overcomplicating the codebase
- Implementing advanced features that obscure learning
- Breaking changes to the API

### Development Setup

```bash
# Clone repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Create feature branch
git checkout -b feature/your-feature-name

# Make changes
vim server.py

# Test changes
python3 test_client.py examples/test.txt

# Check logs
tail -f /tmp/lsp-server.log

# Commit with clear message
git commit -m "Add: Detection for hardcoded credentials"

# Push and create pull request
git push origin feature/your-feature-name
```

### Testing Your Contributions

```bash
# 1. Run with test client
python3 test_client.py examples/test.txt

# 2. Test with multiple files
for file in examples/*.txt; do
    python3 test_client.py "$file"
done

# 3. Test with real editor
nvim test.txt  # With LSP configured

# 4. Check for regressions
diff <(python3 test_client.py examples/test.txt 2>&1) expected_output.txt
```

## 📄 License

This is a demonstration implementation for educational purposes. Feel free to use, modify, and distribute.

**MIT License** - See repository for full license text.

## 🙏 Acknowledgments

- **Microsoft**: For creating and maintaining the LSP specification
- **LSP Community**: For excellent documentation and examples
- **Contributors**: Everyone who has improved this project

## 📚 Additional Resources

### Files in This Repository

- **`server.py`**: Main LSP server implementation (~300 lines)
- **`test_client.py`**: Test client for validation (~200 lines)
- **`examples/test.txt`**: Sample file demonstrating all diagnostic rules
- **`examples/code_sample.py`**: Python code example
- **`QUICKSTART.md`**: 2-minute getting started guide
- **`run_test.sh`**: Automated test script

### Quick Reference Card

```bash
# Start server manually
python3 server.py

# Test with client
python3 test_client.py <file>

# Monitor logs
tail -f /tmp/lsp-server.log

# Quick test
./run_test.sh

# Integration test
nvim test.txt  # Or your editor of choice
```

### Server Capabilities Summary

| Capability | Supported | Notes |
|------------|-----------|-------|
| Initialize/Shutdown | ✅ | Full lifecycle management |
| Document Sync | ✅ | Full sync mode only |
| Diagnostics | ✅ | 4 types of checks |
| Hover | ❌ | Not implemented |
| Completion | ❌ | Not implemented |
| Go to Definition | ❌ | Not implemented |
| Find References | ❌ | Not implemented |
| Formatting | ❌ | Not implemented |
| Code Actions | ❌ | Not implemented |

---

**Questions? Issues? Contributions?**

Open an issue or pull request on [GitHub](https://github.com/rkkeerth/lsp-server)

Happy coding! 🚀
