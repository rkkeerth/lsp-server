# Quick Start Guide

Get up and running with the LSP Server in 5 minutes!

## Prerequisites

- Go 1.21 or higher installed
- Git installed
- A code editor (VS Code, Neovim, or Emacs)

## Installation

### 1. Clone and Build

```bash
# Clone the repository
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server

# Download dependencies and build
make deps
make build
```

Or use the build script:

```bash
./build.sh
```

### 2. Verify Installation

```bash
# Check the binary was created
ls -l lsp-server

# Test it runs (it will wait for JSON-RPC input)
./lsp-server
# Press Ctrl+C to exit
```

## Testing the Server

### Quick Test with Example File

1. **Start the server with logging:**
   ```bash
   ./lsp-server 2> lsp.log &
   LSP_PID=$!
   ```

2. **In another terminal, monitor the logs:**
   ```bash
   tail -f lsp.log
   ```

3. **Stop the server when done:**
   ```bash
   kill $LSP_PID
   ```

### Test with VS Code

1. **Install Go extension** (if not already installed):
   - Open VS Code
   - Install "Go" extension by the Go Team at Google

2. **Open the example file:**
   ```bash
   code examples/sample.go
   ```

3. **Try LSP features:**
   - **Hover**: Hover over `Calculator`, `NewCalculator`, or method names
   - **Go to Definition**: Click on `NewCalculator()` and press F12
   - **Find References**: Right-click on `Calculator` → Find All References
   - **Completion**: Type `calc.` and see suggestions
   - **Outline**: Press Ctrl+Shift+O (Cmd+Shift+O on Mac)
   - **Diagnostics**: Notice the TODO and FIXME warnings

### Test with Neovim

1. **Configure nvim-lspconfig** (add to your init.lua):
   ```lua
   local lspconfig = require('lspconfig')
   local configs = require('lspconfig.configs')

   if not configs.lsp_server then
     configs.lsp_server = {
       default_config = {
         cmd = {'/path/to/lsp-server/lsp-server'},
         filetypes = {'go'},
         root_dir = lspconfig.util.root_pattern('.git', 'go.mod'),
       },
     }
   end

   lspconfig.lsp_server.setup{}
   ```

2. **Open example file:**
   ```bash
   nvim examples/sample.go
   ```

3. **Try LSP features:**
   - **Hover**: Press `K` on a symbol
   - **Go to Definition**: Press `gd` on a symbol
   - **Find References**: Press `gr` on a symbol
   - **Completion**: Press `<C-x><C-o>` (or your completion trigger)

## Understanding What It Does

### Features Demonstrated in sample.go

| Feature | Example | What to See |
|---------|---------|-------------|
| **Hover** | Hover over `Calculator` type | Type information displayed |
| **Definition** | F12 on `NewCalculator()` call | Jump to function definition |
| **References** | Find refs on `Calculator` | All usages highlighted |
| **Symbols** | Document outline | List of functions, types, etc. |
| **Completion** | Type `calc.` | Method suggestions appear |
| **Diagnostics** | TODO/FIXME comments | Warnings in Problems panel |

### How It Works

```
┌─────────────┐           ┌──────────────┐
│             │  stdin/   │              │
│  VS Code/   │  stdout   │  LSP Server  │
│  Neovim     │◄─────────►│  (this app)  │
│             │  JSON-RPC │              │
└─────────────┘           └──────────────┘
```

The LSP server:
1. Reads JSON-RPC messages from stdin
2. Processes requests (hover, definition, etc.)
3. Sends responses back via stdout
4. Logs debug info to stderr

## Next Steps

### Learn More

- **[README.md](README.md)**: Full feature documentation
- **[ARCHITECTURE.md](ARCHITECTURE.md)**: Deep dive into the design
- **[CONTRIBUTING.md](CONTRIBUTING.md)**: How to contribute

### Extend the Server

1. **Add your own language support:**
   - Modify `handlers/handlers.go` to recognize your language patterns
   - Update symbol extraction regex patterns
   - Add language-specific diagnostics

2. **Add new LSP features:**
   - See [LSP Specification](https://microsoft.github.io/language-server-protocol/)
   - Add capability in `server/server.go`
   - Implement handler in `handlers/handlers.go`

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Check for race conditions
go test -race ./...
```

## Troubleshooting

### Server Won't Start

**Problem**: `./lsp-server` command not found

**Solution**: 
```bash
# Make sure you built it
make build

# Check it exists
ls -l lsp-server

# Make it executable if needed
chmod +x lsp-server
```

### Features Not Working

**Problem**: Hover/completion doesn't work in editor

**Solution**:
1. Check server is configured correctly in editor settings
2. Look at server logs: `tail -f lsp.log`
3. Verify document is opened (check for "Document opened" in logs)
4. Try restarting the editor

### No Output

**Problem**: Running `./lsp-server` shows nothing

**Solution**: This is normal! The server waits for JSON-RPC input. It's meant to be used by an editor, not directly. Use `2> lsp.log` to see logs.

## Common Use Cases

### 1. Development Setup

```bash
# Terminal 1: Run tests on file change
while true; do
  inotifywait -e modify ./**/*.go
  make test
done

# Terminal 2: Edit code
vim server/server.go

# Terminal 3: Monitor logs
tail -f lsp.log
```

### 2. Manual Protocol Testing

Create a test file `test_initialize.json`:
```json
Content-Length: 247

{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":12345,"rootUri":"file:///tmp","capabilities":{}}}
```

Send it to the server:
```bash
cat test_initialize.json | ./lsp-server
```

### 3. Integration with Build System

```bash
# Run server as part of CI/CD
make build
make test
./lsp-server --version
```

## Resources

- **LSP Specification**: https://microsoft.github.io/language-server-protocol/
- **JSON-RPC 2.0**: https://www.jsonrpc.org/specification
- **Go LSP Libraries**: https://pkg.go.dev/go.lsp.dev

## Getting Help

- **Issues**: Check [existing issues](https://github.com/rkkeerth/lsp-server/issues)
- **Discussions**: Start a [discussion](https://github.com/rkkeerth/lsp-server/discussions)
- **Documentation**: Read the [full README](README.md)

---

**You're all set!** 🚀

Try opening `examples/sample.go` in your editor and exploring the LSP features. Happy coding!
