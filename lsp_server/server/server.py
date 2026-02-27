"""LSP server implementation."""
from __future__ import annotations

import logging
import sys
from typing import BinaryIO, Optional

from lsp_server.jsonrpc import (
    Transport,
    Request,
    Response,
    new_response,
    new_error_response,
    METHOD_NOT_FOUND,
    INVALID_PARAMS,
)
from lsp_server.protocol.methods import (
    METHOD_INITIALIZE,
    METHOD_INITIALIZED,
    METHOD_SHUTDOWN,
    METHOD_EXIT,
    METHOD_DID_OPEN,
    METHOD_DID_CHANGE,
    METHOD_DID_CLOSE,
)
from lsp_server.protocol.types import (
    InitializeParams,
    InitializeResult,
    ServerCapabilities,
    TextDocumentSyncOptions,
    TextDocumentSyncKind,
    SaveOptions,
    ServerInfo,
    DidOpenTextDocumentParams,
    DidChangeTextDocumentParams,
    DidCloseTextDocumentParams,
)
from lsp_server.document import Manager


class Server:
    """Language Server Protocol server."""

    def __init__(self, reader: BinaryIO, writer: BinaryIO) -> None:
        self.transport = Transport(reader, writer)
        self.documents = Manager()
        self.initialized = False
        self.shutdown = False
        self.logger = logging.getLogger("lsp")
        if not self.logger.handlers:
            handler = logging.StreamHandler(sys.stderr)
            handler.setFormatter(logging.Formatter("[lsp] %(message)s"))
            self.logger.addHandler(handler)
            self.logger.setLevel(logging.DEBUG)

    def run(self) -> None:
        """Main loop: read requests, dispatch, write responses."""
        while True:
            try:
                request = self.transport.read()
            except EOFError:
                return
            except Exception:
                self.logger.exception("Error reading request")
                raise

            response = self.handle(request)
            if response is not None:
                self.transport.write(response)

    def handle(self, request: Request) -> Optional[Response]:
        """Dispatch a request to the appropriate handler."""
        self.logger.info("Received method: %s", request.method)

        if request.method == METHOD_INITIALIZE:
            return self._handle_initialize(request)
        elif request.method == METHOD_INITIALIZED:
            self._handle_initialized(request)
            return None
        elif request.method == METHOD_SHUTDOWN:
            return self._handle_shutdown(request)
        elif request.method == METHOD_EXIT:
            self._handle_exit(request)
            return None
        elif request.method == METHOD_DID_OPEN:
            self._handle_did_open(request)
            return None
        elif request.method == METHOD_DID_CHANGE:
            self._handle_did_change(request)
            return None
        elif request.method == METHOD_DID_CLOSE:
            self._handle_did_close(request)
            return None
        else:
            if request.is_notification():
                self.logger.info("Ignoring unknown notification: %s", request.method)
                return None
            return new_error_response(
                request.id, METHOD_NOT_FOUND, f"Method not found: {request.method}"
            )

    def _handle_initialize(self, request: Request) -> Response:
        """Handle initialize request."""
        params = InitializeParams.from_dict(request.params or {})
        self.logger.info(
            "Initializing server, client process ID: %s", params.process_id
        )

        result = InitializeResult(
            capabilities=ServerCapabilities(
                text_document_sync=TextDocumentSyncOptions(
                    open_close=True,
                    change=TextDocumentSyncKind.FULL,
                    save=SaveOptions(include_text=False),
                ),
            ),
            server_info=ServerInfo(name="lsp-server", version="0.1.0"),
        )
        return new_response(request.id, result.to_dict())

    def _handle_initialized(self, request: Request) -> None:
        """Handle initialized notification."""
        self.initialized = True
        self.logger.info("Server initialized")

    def _handle_shutdown(self, request: Request) -> Response:
        """Handle shutdown request."""
        self.shutdown = True
        self.logger.info("Shutdown requested")
        return new_response(request.id, None)

    def _handle_exit(self, request: Request) -> None:
        """Handle exit notification."""
        if self.shutdown:
            sys.exit(0)
        else:
            sys.exit(1)

    def _handle_did_open(self, request: Request) -> None:
        """Handle textDocument/didOpen notification."""
        params = DidOpenTextDocumentParams.from_dict(request.params or {})
        td = params.text_document
        self.documents.open(td.uri, td.language_id, td.version, td.text)
        self.logger.info("Document opened: %s", td.uri)

    def _handle_did_change(self, request: Request) -> None:
        """Handle textDocument/didChange notification."""
        params = DidChangeTextDocumentParams.from_dict(request.params or {})
        if params.content_changes:
            content = params.content_changes[-1].text
            self.documents.change(
                params.text_document.uri, params.text_document.version, content
            )
            self.logger.info("Document changed: %s", params.text_document.uri)

    def _handle_did_close(self, request: Request) -> None:
        """Handle textDocument/didClose notification."""
        params = DidCloseTextDocumentParams.from_dict(request.params or {})
        self.documents.close(params.text_document.uri)
        self.logger.info("Document closed: %s", params.text_document.uri)
