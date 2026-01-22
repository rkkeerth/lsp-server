#!/bin/bash
# Simple test script for the LSP server
# This demonstrates the JSON-RPC messages that editors send to the server

# Note: This is for demonstration purposes. In practice, use an LSP client library.

echo "Testing LSP Server"
echo "=================="
echo ""
echo "To test the server manually, you can send JSON-RPC messages via stdin."
echo "Example initialize request:"
echo ""
cat << 'EOF'
Content-Length: 147

{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":null,"rootUri":"file:///tmp/test","capabilities":{},"workspaceFolders":null}}
EOF

echo ""
echo ""
echo "Build the server first with:"
echo "  go build -o lsp-server"
echo ""
echo "Then you can test it with an LSP client or by piping JSON-RPC messages."
echo ""
echo "For automated testing, consider using:"
echo "  - An LSP test framework"
echo "  - Editor-specific LSP test utilities"
echo "  - Custom JSON-RPC test client"
