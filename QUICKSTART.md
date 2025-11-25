# Quick Start Guide

Get up and running with the LSP server in 5 minutes!

## Prerequisites

- Go 1.21 or higher installed
- Basic understanding of command line

## Installation

### 1. Clone the Repository

```bash
git clone https://github.com/rkkeerth/lsp-server.git
cd lsp-server
```

### 2. Build the Server

```bash
# Using Make (recommended)
make build

# Or using Go directly
go build -o lsp-server .
```

You should now have an `lsp-server` executable in your directory.

## Quick Test

### Test with the Included Test Client

The easiest way to verify everything works:

```bash
# Run the test client
cd examples
go run test_client.go ../lsp-server
```

**Expected output:**
```
=== Test 1: Initialize ===
Initialize response:
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "capabilities": {
      "textDocumentSync": {
        "openClose": true,
        "change": 1
      }
    }
  }
}
...
=== All tests completed successfully! ===
```

If you see this, congratulations! Your LSP server is working correctly.

## What Just Happened?

The test client just:
1. ✅ Started your LSP server
2. ✅ Sent an initialize request
3. ✅ Opened a test document
4. ✅ Modified the document
5. ✅ Closed the document
6. ✅ Shut down the server cleanly

## Next Steps

### Use with VS Code

See the detailed guide in `examples/vscode-extension/README.md`

**Quick version:**
1. Navigate to the extension directory
2. Install dependencies: `npm install`
3. Compile: `npm run compile`
4. Press F5 in VS Code to launch
5. Open a `.txt` file to activate

### Explore the Code

```bash
# View the main entry point
cat main.go

# Check the protocol definitions
cat internal/lsp/protocol.go

# See how messages are handled
cat internal/lsp/server.go
```

### Run Tests

```bash
# Run all unit tests
make test

# Or with Go directly
go test ./...
```

### Read the Documentation

- **README.md** - Complete feature overview
- **ARCHITECTURE.md** - Deep dive into the design
- **TESTING.md** - Comprehensive testing guide
- **CONTRIBUTING.md** - How to contribute

## Common Issues

### "go: command not found"

**Solution:** Install Go from https://golang.org/dl/

### "Permission denied" when running lsp-server

**Solution:** Make the binary executable:
```bash
chmod +x lsp-server
```

### Test client can't find server

**Solution:** Make sure you're running from the examples directory:
```bash
cd examples
go run test_client.go ../lsp-server
```

## Development Mode

For active development:

```bash
# Format, vet, and build in one command
make dev-build

# Or run individual steps
make fmt    # Format code
make vet    # Run static analysis
make build  # Build binary
```

## What's Included?

```
lsp-server/
├── main.go                   # Entry point - Start here!
├── internal/lsp/            # Core implementation
│   ├── protocol.go          # LSP types
│   ├── server.go            # Main server logic
│   ├── lifecycle.go         # Initialize/shutdown
│   ├── textdocument.go      # Document sync
│   └── server_test.go       # Unit tests
├── examples/
│   ├── test_client.go       # Automated test client
│   └── vscode-extension/    # VS Code integration
├── Makefile                 # Build automation
├── README.md                # Full documentation
├── ARCHITECTURE.md          # Design details
├── TESTING.md               # Testing guide
└── CONTRIBUTING.md          # Contribution guide
```

## Key Features Implemented

✅ **JSON-RPC 2.0** - Full protocol implementation  
✅ **Lifecycle Management** - Initialize, shutdown, exit  
✅ **Document Sync** - didOpen, didChange, didClose  
✅ **Concurrent Processing** - Handle multiple requests  
✅ **Thread-Safe** - Safe concurrent document access  
✅ **Well Tested** - Unit tests and integration tests  

## Supported LSP Methods

| Method | Type | Status |
|--------|------|--------|
| `initialize` | Request | ✅ |
| `initialized` | Notification | ✅ |
| `shutdown` | Request | ✅ |
| `exit` | Notification | ✅ |
| `textDocument/didOpen` | Notification | ✅ |
| `textDocument/didChange` | Notification | ✅ |
| `textDocument/didClose` | Notification | ✅ |

## Learning Path

1. **Start here** - This quick start guide ✓
2. **Run the test** - Verify it works ✓
3. **Read main.go** - Understand the entry point
4. **Explore server.go** - See message handling
5. **Check protocol.go** - Learn LSP types
6. **Read ARCHITECTURE.md** - Understand the design
7. **Try TESTING.md** - Learn testing approaches
8. **Extend it!** - Add new capabilities

## Need Help?

- 📖 Read **README.md** for detailed documentation
- 🏗️ Check **ARCHITECTURE.md** for design details
- 🧪 See **TESTING.md** for testing approaches
- 🤝 Review **CONTRIBUTING.md** for contribution guidelines
- 🐛 Open an issue on GitHub for bugs
- 💡 Open a discussion for questions

## Useful Commands Cheat Sheet

```bash
# Build
make build                    # Build the binary
go build -o lsp-server .     # Build with Go

# Test
make test                     # Run unit tests
go test -v ./...             # Run tests with verbose output
cd examples && go run test_client.go ../lsp-server  # Integration test

# Development
make fmt                      # Format code
make vet                      # Run static analysis
make dev-build               # Format, vet, and build

# Clean
make clean                    # Remove build artifacts

# Run
./lsp-server                 # Start the server (manual mode)
./lsp-server 2>server.log    # Start and log to file
```

## Sample Session

Here's what a typical LSP session looks like:

```
1. Client → Server: initialize
2. Server → Client: initialize result (capabilities)
3. Client → Server: initialized
4. Client → Server: textDocument/didOpen
5. Client → Server: textDocument/didChange (multiple times)
6. Client → Server: textDocument/didClose
7. Client → Server: shutdown
8. Server → Client: shutdown result
9. Client → Server: exit
10. Server exits
```

## What Makes This Implementation Special?

1. **Educational** - Clear, well-documented code
2. **Complete** - Full LSP lifecycle implementation
3. **Tested** - Unit tests and integration tests
4. **Concurrent** - Handles multiple requests efficiently
5. **Extensible** - Easy to add new capabilities
6. **Production-Ready** - Thread-safe and robust

## Next Features to Explore

Want to extend the server? Try adding:
- **Hover support** - Show information on hover
- **Completion** - Provide code completion
- **Go to definition** - Navigate to symbol definitions
- **Diagnostics** - Report errors and warnings
- **Incremental sync** - More efficient document updates

See ARCHITECTURE.md for implementation guidance!

## Success!

You now have a working LSP server! 🎉

Explore the code, run tests, and start building your own language features.

Happy coding! 🚀
