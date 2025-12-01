#!/usr/bin/env python3
"""
Simple LSP client example to test the server
This script demonstrates basic LSP protocol communication
"""

import json
import subprocess
import sys
from typing import Dict, Any, Optional


class LSPClient:
    def __init__(self, server_path: str):
        self.server = subprocess.Popen(
            [server_path],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            bufsize=0
        )
        self.request_id = 0

    def send_request(self, method: str, params: Optional[Dict[str, Any]] = None) -> int:
        """Send a request and return the request ID"""
        self.request_id += 1
        message = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": method,
            "params": params or {}
        }
        self._send_message(message)
        return self.request_id

    def send_notification(self, method: str, params: Optional[Dict[str, Any]] = None):
        """Send a notification (no response expected)"""
        message = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params or {}
        }
        self._send_message(message)

    def _send_message(self, message: Dict[str, Any]):
        """Send a JSON-RPC message with proper headers"""
        content = json.dumps(message)
        content_bytes = content.encode('utf-8')
        header = f"Content-Length: {len(content_bytes)}\r\n\r\n"
        full_message = header.encode('utf-8') + content_bytes
        
        print(f"→ Sending: {message['method']}")
        self.server.stdin.write(full_message)
        self.server.stdin.flush()

    def read_response(self) -> Optional[Dict[str, Any]]:
        """Read a response from the server"""
        # Read headers
        headers = {}
        while True:
            line = self.server.stdout.readline().decode('utf-8').strip()
            if not line:
                break
            if ':' in line:
                key, value = line.split(':', 1)
                headers[key.strip()] = value.strip()

        # Read content
        content_length = int(headers.get('Content-Length', 0))
        if content_length == 0:
            return None

        content = self.server.stdout.read(content_length).decode('utf-8')
        response = json.loads(content)
        
        print(f"← Received: {json.dumps(response, indent=2)}")
        return response

    def close(self):
        """Close the client connection"""
        self.server.stdin.close()
        self.server.stdout.close()
        self.server.stderr.close()
        self.server.wait()


def main():
    if len(sys.argv) < 2:
        print("Usage: python3 simple_client.py <path-to-lsp-server>")
        sys.exit(1)

    server_path = sys.argv[1]
    client = LSPClient(server_path)

    try:
        # Initialize
        print("\n=== Initializing server ===")
        client.send_request("initialize", {
            "processId": None,
            "rootUri": "file:///tmp/test-workspace",
            "capabilities": {
                "textDocument": {
                    "hover": {
                        "contentFormat": ["markdown", "plaintext"]
                    },
                    "synchronization": {
                        "didSave": True
                    }
                }
            },
            "trace": "off",
            "workspaceFolders": [
                {
                    "uri": "file:///tmp/test-workspace",
                    "name": "test-workspace"
                }
            ]
        })
        init_response = client.read_response()

        # Send initialized notification
        print("\n=== Server initialized ===")
        client.send_notification("initialized", {})

        # Open a document
        print("\n=== Opening document ===")
        client.send_notification("textDocument/didOpen", {
            "textDocument": {
                "uri": "file:///tmp/test.go",
                "languageId": "go",
                "version": 1,
                "text": "package main\n\nfunc HelloWorld() string {\n\treturn \"Hello, World!\"\n}\n"
            }
        })

        # Request hover information
        print("\n=== Requesting hover info ===")
        client.send_request("textDocument/hover", {
            "textDocument": {
                "uri": "file:///tmp/test.go"
            },
            "position": {
                "line": 2,
                "character": 5
            }
        })
        hover_response = client.read_response()

        # Change document
        print("\n=== Changing document ===")
        client.send_notification("textDocument/didChange", {
            "textDocument": {
                "uri": "file:///tmp/test.go",
                "version": 2
            },
            "contentChanges": [
                {
                    "text": "package main\n\nfunc HelloWorld() string {\n\treturn \"Hello, LSP!\"\n}\n"
                }
            ]
        })

        # Close document
        print("\n=== Closing document ===")
        client.send_notification("textDocument/didClose", {
            "textDocument": {
                "uri": "file:///tmp/test.go"
            }
        })

        # Shutdown
        print("\n=== Shutting down server ===")
        client.send_request("shutdown", {})
        shutdown_response = client.read_response()

        # Exit
        print("\n=== Exiting ===")
        client.send_notification("exit", {})

        print("\n✓ Test completed successfully!")

    except Exception as e:
        print(f"\n✗ Error: {e}", file=sys.stderr)
        sys.exit(1)
    finally:
        client.close()


if __name__ == "__main__":
    main()
