#!/bin/bash
# Simple script to test LSP server manually

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY="$SCRIPT_DIR/../lsp-server"

if [ ! -f "$BINARY" ]; then
    echo "Error: LSP server binary not found at $BINARY"
    echo "Please build it first with: make build"
    exit 1
fi

echo "Testing LSP Server..."
echo "====================="
echo ""

# Test 1: Initialize request
echo "Test 1: Sending initialize request..."
INIT_RESPONSE=$(echo '{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "processId": null,
    "clientInfo": {
      "name": "test-client",
      "version": "1.0.0"
    },
    "rootUri": null,
    "capabilities": {}
  }
}' | "$BINARY" 2>/dev/null)

if echo "$INIT_RESPONSE" | grep -q "serverInfo"; then
    echo "✓ Initialize request successful"
else
    echo "✗ Initialize request failed"
    echo "Response: $INIT_RESPONSE"
fi

echo ""
echo "Testing complete!"
echo ""
echo "For interactive testing, run the server and send JSON-RPC messages:"
echo "  $BINARY"
echo ""
echo "Or configure it in your editor using the configurations in README.md"
