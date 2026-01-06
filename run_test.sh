#!/bin/bash
# Simple script to test the LSP server

set -e

echo "======================================"
echo "Basic LSP Server Test Script"
echo "======================================"
echo ""

# Check Python version
echo "Checking Python version..."
python3 --version

echo ""
echo "Testing server compilation..."
python3 -m py_compile server.py
python3 -m py_compile test_client.py
echo "✓ Server code compiles successfully"

echo ""
echo "======================================"
echo "Running Test 1: Basic text file"
echo "======================================"
python3 test_client.py examples/test.txt

echo ""
echo "======================================"
echo "Running Test 2: Python code file"
echo "======================================"
python3 test_client.py examples/code_sample.py

echo ""
echo "======================================"
echo "All tests completed successfully!"
echo "======================================"
echo ""
echo "To view server logs:"
echo "  tail -f /tmp/lsp-server.log"
echo ""
echo "To use with your editor, see README.md"
