#!/bin/bash
# Test script for the LSP server

# Start the server in the background
./lsp-server 2> test_server.log &
SERVER_PID=$!

# Give it a moment to start
sleep 1

# Function to send a message
send_message() {
    local message=$1
    local length=${#message}
    echo -e "Content-Length: $length\r\n\r\n$message"
}

# Initialize request
INIT_MSG='{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"processId":null,"rootUri":"file:///tmp","capabilities":{"textDocument":{"hover":{"contentFormat":["markdown","plaintext"]},"synchronization":{"didSave":true}}},"trace":"off","workspaceFolders":[{"uri":"file:///tmp","name":"tmp"}]}}'

# Initialized notification
INITIALIZED_MSG='{"jsonrpc":"2.0","method":"initialized","params":{}}'

# DidOpen notification
DIDOPEN_MSG='{"jsonrpc":"2.0","method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///tmp/test.txt","languageId":"text","version":1,"text":"Hello world\nTest document"}}}'

# Hover request
HOVER_MSG='{"jsonrpc":"2.0","id":2,"method":"textDocument/hover","params":{"textDocument":{"uri":"file:///tmp/test.txt"},"position":{"line":0,"character":1}}}'

# Shutdown request
SHUTDOWN_MSG='{"jsonrpc":"2.0","id":3,"method":"shutdown","params":{}}'

# Exit notification
EXIT_MSG='{"jsonrpc":"2.0","method":"exit","params":{}}'

# Send messages
{
    send_message "$INIT_MSG"
    sleep 0.5
    send_message "$INITIALIZED_MSG"
    sleep 0.5
    send_message "$DIDOPEN_MSG"
    sleep 0.5
    send_message "$HOVER_MSG"
    sleep 0.5
    send_message "$SHUTDOWN_MSG"
    sleep 0.5
    send_message "$EXIT_MSG"
} | timeout 10 ./lsp-server 2> test_server.log > test_output.log

# Check logs
echo "=== Server Log ==="
cat test_server.log

echo ""
echo "=== Server Responses ==="
cat test_output.log

# Cleanup
rm -f test_server.log test_output.log

echo ""
echo "Test completed!"
