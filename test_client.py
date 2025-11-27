#!/usr/bin/env python3
"""
Simple LSP client to test the basic LSP server.

This client demonstrates:
- Initializing a connection to the LSP server
- Opening a document
- Receiving diagnostics
- Shutting down the server

This implementation uses only Python standard library (no external dependencies).
"""

import json
import subprocess
import sys
import threading
import time
from pathlib import Path
from typing import Optional, Dict, Any, List


class LSPTestClient:
    """Simple LSP client for testing."""
    
    def __init__(self):
        self.process: Optional[subprocess.Popen] = None
        self.message_id = 0
        self.diagnostics: List[Dict[str, Any]] = []
        self.response_queue: Dict[int, Any] = {}
        
    def get_next_id(self) -> int:
        """Get next message ID."""
        self.message_id += 1
        return self.message_id
    
    def start_server(self, server_command: List[str]):
        """Start the LSP server process."""
        print("Starting LSP server...")
        self.process = subprocess.Popen(
            server_command,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=False,
            bufsize=0
        )
        
        # Start thread to read responses
        self.reader_thread = threading.Thread(target=self._read_messages, daemon=True)
        self.reader_thread.start()
        
        time.sleep(0.5)  # Give server time to start
        
    def _read_messages(self):
        """Read messages from server in background thread."""
        while True:
            try:
                message = self._read_message()
                if message is None:
                    break
                
                if "id" in message:
                    # This is a response
                    self.response_queue[message["id"]] = message
                elif message.get("method") == "textDocument/publishDiagnostics":
                    # This is a diagnostics notification
                    self._handle_diagnostics(message["params"])
                    
            except Exception as e:
                print(f"Error reading message: {e}")
                break
    
    def _read_message(self) -> Optional[Dict[str, Any]]:
        """Read a single message from the server."""
        if not self.process or not self.process.stdout:
            return None
        
        # Read headers
        headers = {}
        while True:
            line = self.process.stdout.readline()
            if not line:
                return None
            line = line.decode('utf-8').strip()
            if not line:
                break
            if ": " in line:
                key, value = line.split(": ", 1)
                headers[key] = value
        
        # Read content
        content_length = int(headers.get("Content-Length", 0))
        if content_length == 0:
            return None
        
        content = self.process.stdout.read(content_length).decode('utf-8')
        return json.loads(content)
    
    def _send_message(self, message: Dict[str, Any]):
        """Send a message to the server."""
        if not self.process or not self.process.stdin:
            raise RuntimeError("Server not started")
        
        content = json.dumps(message)
        content_bytes = content.encode('utf-8')
        content_length = len(content_bytes)
        
        header = f"Content-Length: {content_length}\r\n\r\n".encode('utf-8')
        self.process.stdin.write(header + content_bytes)
        self.process.stdin.flush()
    
    def _send_request(self, method: str, params: Any) -> int:
        """Send a request and return the message ID."""
        msg_id = self.get_next_id()
        message = {
            "jsonrpc": "2.0",
            "id": msg_id,
            "method": method,
            "params": params
        }
        self._send_message(message)
        return msg_id
    
    def _send_notification(self, method: str, params: Any):
        """Send a notification (no response expected)."""
        message = {
            "jsonrpc": "2.0",
            "method": method,
            "params": params
        }
        self._send_message(message)
    
    def _wait_for_response(self, msg_id: int, timeout: float = 5.0) -> Optional[Dict[str, Any]]:
        """Wait for a response with the given message ID."""
        start_time = time.time()
        while time.time() - start_time < timeout:
            if msg_id in self.response_queue:
                return self.response_queue.pop(msg_id)
            time.sleep(0.1)
        return None
    
    def _handle_diagnostics(self, params: Dict[str, Any]):
        """Handle diagnostics notification."""
        uri = params.get("uri", "")
        diagnostics = params.get("diagnostics", [])
        
        print(f"\n{'='*60}")
        print(f"Diagnostics received for: {uri}")
        print(f"{'='*60}")
        
        if not diagnostics:
            print("No issues found!")
        else:
            print(f"Found {len(diagnostics)} issue(s):\n")
            for diag in diagnostics:
                severity_map = {
                    1: "ERROR",
                    2: "WARNING",
                    3: "INFO",
                    4: "HINT"
                }
                severity = severity_map.get(diag.get("severity", 3), "UNKNOWN")
                line = diag["range"]["start"]["line"] + 1  # Convert to 1-indexed
                col = diag["range"]["start"]["character"] + 1
                message = diag.get("message", "")
                print(f"  [{severity}] Line {line}, Col {col}: {message}")
        
        print(f"{'='*60}\n")
        self.diagnostics = diagnostics
    
    def initialize(self):
        """Initialize the LSP connection."""
        print("Initializing LSP client...")
        
        params = {
            "processId": None,
            "rootUri": None,
            "capabilities": {}
        }
        
        msg_id = self._send_request("initialize", params)
        response = self._wait_for_response(msg_id)
        
        if response and "result" in response:
            print(f"Server initialized successfully")
            capabilities = response["result"].get("capabilities", {})
            print(f"Server capabilities: {json.dumps(capabilities, indent=2)}")
            
            # Send initialized notification
            self._send_notification("initialized", {})
            print("Client initialized notification sent")
        else:
            print("Failed to initialize server")
    
    def open_document(self, file_path: str):
        """Open a document for analysis."""
        print(f"\nOpening document: {file_path}")
        
        # Read the file content
        with open(file_path, 'r') as f:
            content = f.read()
        
        # Convert to file URI
        file_uri = Path(file_path).resolve().as_uri()
        
        # Send didOpen notification
        params = {
            "textDocument": {
                "uri": file_uri,
                "languageId": "plaintext",
                "version": 1,
                "text": content
            }
        }
        
        self._send_notification("textDocument/didOpen", params)
        
        # Wait for diagnostics to be processed
        time.sleep(1.0)
    
    def shutdown(self):
        """Shutdown the LSP server."""
        print("\nShutting down server...")
        
        # Send shutdown request
        msg_id = self._send_request("shutdown", None)
        response = self._wait_for_response(msg_id)
        
        # Send exit notification
        self._send_notification("exit", None)
        
        # Wait for process to exit
        if self.process:
            try:
                self.process.wait(timeout=2.0)
            except subprocess.TimeoutExpired:
                self.process.kill()
        
        print("Server shutdown complete")


def main(file_path: Optional[str] = None):
    """Run the test client."""
    if file_path is None:
        file_path = "examples/test.txt"
    
    if not Path(file_path).exists():
        print(f"Error: File not found: {file_path}")
        sys.exit(1)
    
    print("="*60)
    print("Basic LSP Server Test Client")
    print("="*60)
    
    client = LSPTestClient()
    
    try:
        # Start the server
        client.start_server([sys.executable, "server.py"])
        
        # Initialize
        client.initialize()
        
        # Open and analyze document
        client.open_document(file_path)
        
        # Show summary
        print(f"Test complete! Analyzed: {file_path}")
        print(f"Total diagnostics found: {len(client.diagnostics)}")
        
    except Exception as e:
        print(f"Error during test: {e}")
        import traceback
        traceback.print_exc()
    finally:
        # Shutdown
        try:
            client.shutdown()
        except:
            pass


if __name__ == "__main__":
    file_path = sys.argv[1] if len(sys.argv) > 1 else None
    main(file_path)
