#!/usr/bin/env python3
"""
Basic Language Server Protocol (LSP) Server Implementation

This server demonstrates core LSP capabilities:
- Initialize and shutdown lifecycle management
- Text document synchronization
- Basic diagnostics (detects TODO, FIXME, and simple patterns)

This implementation uses only Python standard library (no external dependencies).
"""

import json
import logging
import re
import sys
from typing import Dict, List, Any, Optional

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    filename="/tmp/lsp-server.log"
)
logger = logging.getLogger(__name__)


class LSPServer:
    """Basic Language Server Protocol implementation."""
    
    def __init__(self):
        self.documents: Dict[str, str] = {}
        self.running = True
        self.initialized = False
        
    def analyze_document(self, text: str) -> List[Dict[str, Any]]:
        """
        Analyze text document and return diagnostics.
        
        This function detects:
        - TODO comments (informational)
        - FIXME comments (warning)
        - Lines longer than 120 characters (warning)
        - Duplicate consecutive lines (information)
        """
        diagnostics = []
        lines = text.split('\n')
        
        for line_num, line in enumerate(lines):
            # Check for TODO markers
            match = re.search(r'\bTODO\b', line, re.IGNORECASE)
            if match:
                diagnostics.append({
                    "range": {
                        "start": {"line": line_num, "character": match.start()},
                        "end": {"line": line_num, "character": match.end()}
                    },
                    "message": "TODO found: Consider addressing this item",
                    "severity": 3,  # Information
                    "source": "basic-lsp-server"
                })
            
            # Check for FIXME markers
            match = re.search(r'\bFIXME\b', line, re.IGNORECASE)
            if match:
                diagnostics.append({
                    "range": {
                        "start": {"line": line_num, "character": match.start()},
                        "end": {"line": line_num, "character": match.end()}
                    },
                    "message": "FIXME found: This requires immediate attention",
                    "severity": 2,  # Warning
                    "source": "basic-lsp-server"
                })
            
            # Check for lines that are too long
            if len(line) > 120:
                diagnostics.append({
                    "range": {
                        "start": {"line": line_num, "character": 120},
                        "end": {"line": line_num, "character": len(line)}
                    },
                    "message": f"Line too long ({len(line)} > 120 characters)",
                    "severity": 2,  # Warning
                    "source": "basic-lsp-server"
                })
            
            # Check for duplicate consecutive lines
            if line_num > 0 and line.strip() and line == lines[line_num - 1]:
                diagnostics.append({
                    "range": {
                        "start": {"line": line_num, "character": 0},
                        "end": {"line": line_num, "character": len(line)}
                    },
                    "message": "Duplicate line detected",
                    "severity": 3,  # Information
                    "source": "basic-lsp-server"
                })
        
        logger.info(f"Analysis complete: found {len(diagnostics)} diagnostics")
        return diagnostics
    
    def handle_initialize(self, params: Dict[str, Any]) -> Dict[str, Any]:
        """Handle initialize request."""
        logger.info("Handling initialize request")
        self.initialized = True
        return {
            "capabilities": {
                "textDocumentSync": {
                    "openClose": True,
                    "change": 1,  # Full document sync
                    "save": {"includeText": False}
                },
                "diagnosticProvider": {
                    "interFileDependencies": False,
                    "workspaceDiagnostics": False
                }
            },
            "serverInfo": {
                "name": "basic-lsp-server",
                "version": "0.1.0"
            }
        }
    
    def handle_shutdown(self, params: Dict[str, Any]) -> None:
        """Handle shutdown request."""
        logger.info("Handling shutdown request")
        self.running = False
        return None
    
    def handle_text_document_did_open(self, params: Dict[str, Any]) -> None:
        """Handle textDocument/didOpen notification."""
        text_document = params.get("textDocument", {})
        uri = text_document.get("uri")
        text = text_document.get("text", "")
        
        logger.info(f"Document opened: {uri}")
        self.documents[uri] = text
        
        # Analyze and publish diagnostics
        diagnostics = self.analyze_document(text)
        self.publish_diagnostics(uri, diagnostics)
    
    def handle_text_document_did_change(self, params: Dict[str, Any]) -> None:
        """Handle textDocument/didChange notification."""
        text_document = params.get("textDocument", {})
        uri = text_document.get("uri")
        content_changes = params.get("contentChanges", [])
        
        logger.info(f"Document changed: {uri}")
        
        # Full document sync
        if content_changes:
            text = content_changes[0].get("text", "")
            self.documents[uri] = text
            
            # Re-analyze and publish diagnostics
            diagnostics = self.analyze_document(text)
            self.publish_diagnostics(uri, diagnostics)
    
    def handle_text_document_did_close(self, params: Dict[str, Any]) -> None:
        """Handle textDocument/didClose notification."""
        text_document = params.get("textDocument", {})
        uri = text_document.get("uri")
        
        logger.info(f"Document closed: {uri}")
        
        # Remove from cache and clear diagnostics
        if uri in self.documents:
            del self.documents[uri]
        self.publish_diagnostics(uri, [])
    
    def publish_diagnostics(self, uri: str, diagnostics: List[Dict[str, Any]]) -> None:
        """Publish diagnostics to the client."""
        notification = {
            "jsonrpc": "2.0",
            "method": "textDocument/publishDiagnostics",
            "params": {
                "uri": uri,
                "diagnostics": diagnostics
            }
        }
        self.send_message(notification)
    
    def send_message(self, message: Dict[str, Any]) -> None:
        """Send a message to the client via stdout."""
        content = json.dumps(message)
        content_length = len(content.encode('utf-8'))
        
        response = f"Content-Length: {content_length}\r\n\r\n{content}"
        sys.stdout.write(response)
        sys.stdout.flush()
        
        logger.info(f"Sent message: {message.get('method', message.get('id', 'unknown'))}")
    
    def handle_request(self, request_id: Any, method: str, params: Dict[str, Any]) -> None:
        """Handle a request message."""
        logger.info(f"Handling request: {method}")
        
        handlers = {
            "initialize": self.handle_initialize,
            "shutdown": self.handle_shutdown,
        }
        
        handler = handlers.get(method)
        if handler:
            try:
                result = handler(params)
                response = {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "result": result
                }
                self.send_message(response)
            except Exception as e:
                logger.error(f"Error handling request {method}: {e}", exc_info=True)
                error_response = {
                    "jsonrpc": "2.0",
                    "id": request_id,
                    "error": {
                        "code": -32603,
                        "message": str(e)
                    }
                }
                self.send_message(error_response)
        else:
            logger.warning(f"Unknown request method: {method}")
    
    def handle_notification(self, method: str, params: Dict[str, Any]) -> None:
        """Handle a notification message."""
        logger.info(f"Handling notification: {method}")
        
        handlers = {
            "initialized": lambda p: None,  # No-op
            "textDocument/didOpen": self.handle_text_document_did_open,
            "textDocument/didChange": self.handle_text_document_did_change,
            "textDocument/didClose": self.handle_text_document_did_close,
            "exit": lambda p: sys.exit(0),
        }
        
        handler = handlers.get(method)
        if handler:
            try:
                handler(params)
            except Exception as e:
                logger.error(f"Error handling notification {method}: {e}", exc_info=True)
        else:
            logger.warning(f"Unknown notification method: {method}")
    
    def read_message(self) -> Optional[Dict[str, Any]]:
        """Read a message from stdin."""
        # Read headers
        headers = {}
        while True:
            line = sys.stdin.readline()
            if not line:
                return None
            line = line.strip()
            if not line:
                break
            key, value = line.split(": ", 1)
            headers[key] = value
        
        # Read content
        content_length = int(headers.get("Content-Length", 0))
        if content_length == 0:
            return None
        
        content = sys.stdin.read(content_length)
        message = json.loads(content)
        
        logger.info(f"Received message: {message.get('method', message.get('id', 'unknown'))}")
        return message
    
    def run(self) -> None:
        """Main server loop."""
        logger.info("Starting Basic LSP Server...")
        
        try:
            while self.running:
                message = self.read_message()
                if not message:
                    break
                
                method = message.get("method")
                params = message.get("params", {})
                request_id = message.get("id")
                
                if request_id is not None:
                    # This is a request
                    self.handle_request(request_id, method, params)
                else:
                    # This is a notification
                    self.handle_notification(method, params)
        
        except KeyboardInterrupt:
            logger.info("Server interrupted")
        except Exception as e:
            logger.error(f"Server error: {e}", exc_info=True)
        finally:
            logger.info("Server stopped")


def main():
    """Start the LSP server."""
    server = LSPServer()
    server.run()


if __name__ == "__main__":
    main()
