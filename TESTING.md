# Testing Guide

This document describes how to test the LSP server implementation.

## Quick Start

```bash
# Build the server
make build

# Run the test client
cd examples
go run test_client.go ../lsp-server
```

## Testing Methods

### 1. Automated Test Client

The `examples/test_client.go` program provides automated testing of core LSP functionality.

**Usage:**
```bash
# Build the server first
make build

# Run the test client
cd examples
go run test_client.go ../lsp-server

# Or specify a custom server path
go run test_client.go /path/to/your/lsp-server
```

**What it tests:**
- Initialize handshake
- Initialized notification
- textDocument/didOpen
- textDocument/didChange
- textDocument/didClose
- Shutdown sequence
- Exit notification

**Expected output:**
```
=== Test 1: Initialize ===
Initialize response:
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "capabilities": { ... }
  }
}

=== Test 2: Initialized Notification ===
Sent initialized notification
...
=== All tests completed successfully! ===
```

### 2. Manual Testing with Netcat

You can manually test LSP messages using netcat:

```bash
# Start the server
./lsp-server 2>server.log

# In another terminal, prepare a test message
cat > init.json << 'EOF'
{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":null,"rootUri":"file:///tmp","capabilities":{}}}
EOF

# Calculate content length and send
CONTENT=$(cat init.json)
LENGTH=${#CONTENT}
echo -e "Content-Length: $LENGTH\r\n\r\n$CONTENT" | nc localhost 9999
```

### 3. VS Code Extension Testing

**Setup:**

1. Navigate to the example extension:
```bash
cd examples/vscode-extension
```

2. Install dependencies:
```bash
npm install
```

3. Compile the extension:
```bash
npm run compile
```

4. Build the LSP server:
```bash
cd ../..
make build
```

**Running:**

1. Open `examples/vscode-extension` in VS Code
2. Press F5 to launch Extension Development Host
3. In the new window, create a `.txt` file
4. Type some text and observe the server logs

**Viewing Logs:**

- Server stderr: Check the Output panel → "Basic LSP Server"
- LSP trace: Set `"basicLspServer.trace.server": "verbose"` in settings

### 4. Unit Testing

Run Go unit tests:

```bash
make test
```

Or run with coverage:

```bash
go test -cover ./...
```

## Test Scenarios

### Scenario 1: Basic Document Lifecycle

1. Initialize server
2. Open document
3. Modify document
4. Close document
5. Shutdown server

**Expected behavior:**
- Server accepts all messages
- Document state is tracked correctly
- No errors in logs

### Scenario 2: Multiple Documents

1. Initialize server
2. Open document A
3. Open document B
4. Modify document A
5. Modify document B
6. Close document A
7. Modify document B (should still work)
8. Close document B

**Expected behavior:**
- Server tracks multiple documents independently
- Closing one document doesn't affect others

### Scenario 3: Invalid Messages

1. Initialize server
2. Send malformed JSON
3. Send request with unknown method
4. Send notification with invalid params

**Expected behavior:**
- Server responds with appropriate error codes
- Server continues operating after errors
- No crashes

### Scenario 4: Concurrent Requests

1. Initialize server
2. Send multiple didChange notifications rapidly
3. Verify all changes are processed

**Expected behavior:**
- Server handles concurrent messages
- No race conditions or deadlocks
- Document state is consistent

## Debugging

### Enable Verbose Logging

The server logs to stderr. Redirect to a file for analysis:

```bash
./lsp-server 2>debug.log
```

### Common Issues

**Issue: "Content-Length header missing"**
- Ensure messages include proper headers
- Format: `Content-Length: {bytes}\r\n\r\n{json}`

**Issue: "Method not found"**
- Check method name matches LSP specification
- Verify server supports the method

**Issue: "Server already initialized"**
- Don't send initialize twice
- Check initialization sequence

### Message Debugging

To see all messages, check server logs:

```bash
tail -f server.log | grep "Received"
```

## Performance Testing

### Message Throughput

Test how many messages the server can handle:

```bash
# Send 1000 didChange notifications
for i in {1..1000}; do
  echo "Sending message $i"
  # Send didChange message
done
```

### Memory Usage

Monitor server memory:

```bash
# Get server PID
PID=$(pgrep lsp-server)

# Monitor memory
watch -n 1 "ps -p $PID -o rss,vsz,cmd"
```

## CI/CD Integration

### GitHub Actions Example

```yaml
name: Test LSP Server

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      
      - name: Build
        run: make build
      
      - name: Test
        run: make test
      
      - name: Integration Test
        run: |
          cd examples
          timeout 30s go run test_client.go ../lsp-server
```

## Troubleshooting

### Server Won't Start

1. Check binary exists: `ls -la lsp-server`
2. Check permissions: `chmod +x lsp-server`
3. Try running directly: `./lsp-server`

### No Response from Server

1. Check server is running: `ps aux | grep lsp-server`
2. Verify stdin/stdout aren't blocked
3. Check message format is correct

### Test Client Fails

1. Build server first: `make build`
2. Check server path is correct
3. Look at server logs in stderr

## Additional Resources

- [LSP Specification](https://microsoft.github.io/language-server-protocol/)
- [JSON-RPC 2.0](https://www.jsonrpc.org/specification)
- [Testing LSP Servers](https://github.com/Microsoft/language-server-protocol/blob/main/testingLanguageServers.md)
