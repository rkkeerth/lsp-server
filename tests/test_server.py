"""Tests for LSP server request handling."""
import io
import json

import pytest

from lsp_server.jsonrpc.types import METHOD_NOT_FOUND, Request
from lsp_server.server import Server


def make_server():
    """Create a Server with dummy BytesIO reader/writer."""
    reader = io.BytesIO()
    writer = io.BytesIO()
    return Server(reader, writer)


def frame_message(obj):
    """Build a Content-Length framed JSON-RPC message."""
    body = json.dumps(obj).encode("utf-8")
    header = f"Content-Length: {len(body)}\r\n\r\n".encode("utf-8")
    return header + body


class TestHandleInitialize:
    def test_returns_correct_capabilities(self):
        server = make_server()
        req = Request(
            jsonrpc="2.0",
            id=1,
            method="initialize",
            params={"processId": None, "rootUri": None, "capabilities": {}},
        )
        resp = server.handle(req)
        assert resp is not None
        result = resp.result
        assert result is not None

        sync = result["capabilities"]["textDocumentSync"]
        assert sync["openClose"] is True
        assert sync["change"] == 1
        assert sync["save"]["includeText"] is False

        info = result["serverInfo"]
        assert info["name"] == "lsp-server"
        assert info["version"] == "0.1.0"


class TestHandleInitialized:
    def test_returns_none_and_sets_flag(self):
        server = make_server()
        req = Request(jsonrpc="2.0", id=None, method="initialized", params={})
        resp = server.handle(req)
        assert resp is None
        assert server.initialized is True


class TestHandleShutdown:
    def test_returns_null_result_and_sets_flag(self):
        server = make_server()
        req = Request(jsonrpc="2.0", id=2, method="shutdown")
        resp = server.handle(req)
        assert resp is not None
        assert resp.result is None
        assert server.shutdown is True


class TestHandleExit:
    def test_exit_after_shutdown_exits_zero(self):
        server = make_server()
        server.shutdown = True
        req = Request(jsonrpc="2.0", id=None, method="exit")
        with pytest.raises(SystemExit) as exc_info:
            server.handle(req)
        assert exc_info.value.code == 0

    def test_exit_without_shutdown_exits_one(self):
        server = make_server()
        req = Request(jsonrpc="2.0", id=None, method="exit")
        with pytest.raises(SystemExit) as exc_info:
            server.handle(req)
        assert exc_info.value.code == 1


class TestHandleDidOpen:
    def test_stores_document(self):
        server = make_server()
        req = Request(
            jsonrpc="2.0",
            id=None,
            method="textDocument/didOpen",
            params={
                "textDocument": {
                    "uri": "file:///test.py",
                    "languageId": "python",
                    "version": 1,
                    "text": "print('hello')",
                }
            },
        )
        resp = server.handle(req)
        assert resp is None
        doc = server.documents.get("file:///test.py")
        assert doc is not None
        assert doc.content == "print('hello')"
        assert doc.language_id == "python"
        assert doc.version == 1


class TestHandleDidChange:
    def test_updates_document_content(self):
        server = make_server()
        # First open a document
        open_req = Request(
            jsonrpc="2.0",
            id=None,
            method="textDocument/didOpen",
            params={
                "textDocument": {
                    "uri": "file:///test.py",
                    "languageId": "python",
                    "version": 1,
                    "text": "print('hello')",
                }
            },
        )
        server.handle(open_req)

        # Then change it
        change_req = Request(
            jsonrpc="2.0",
            id=None,
            method="textDocument/didChange",
            params={
                "textDocument": {"uri": "file:///test.py", "version": 2},
                "contentChanges": [{"text": "print('world')"}],
            },
        )
        resp = server.handle(change_req)
        assert resp is None
        doc = server.documents.get("file:///test.py")
        assert doc is not None
        assert doc.content == "print('world')"
        assert doc.version == 2


class TestHandleDidClose:
    def test_removes_document(self):
        server = make_server()
        # First open a document
        open_req = Request(
            jsonrpc="2.0",
            id=None,
            method="textDocument/didOpen",
            params={
                "textDocument": {
                    "uri": "file:///test.py",
                    "languageId": "python",
                    "version": 1,
                    "text": "print('hello')",
                }
            },
        )
        server.handle(open_req)

        # Then close it
        close_req = Request(
            jsonrpc="2.0",
            id=None,
            method="textDocument/didClose",
            params={"textDocument": {"uri": "file:///test.py"}},
        )
        resp = server.handle(close_req)
        assert resp is None
        assert server.documents.get("file:///test.py") is None


class TestHandleUnknown:
    def test_unknown_request_returns_method_not_found(self):
        server = make_server()
        req = Request(jsonrpc="2.0", id=5, method="textDocument/hover", params={})
        resp = server.handle(req)
        assert resp is not None
        assert resp.error is not None
        assert resp.error.code == METHOD_NOT_FOUND

    def test_unknown_notification_returns_none(self):
        server = make_server()
        req = Request(jsonrpc="2.0", id=None, method="$/cancelRequest", params={})
        resp = server.handle(req)
        assert resp is None


class TestServerRun:
    def test_processes_initialize_and_shutdown(self):
        init_msg = frame_message(
            {
                "jsonrpc": "2.0",
                "id": 1,
                "method": "initialize",
                "params": {
                    "processId": None,
                    "rootUri": None,
                    "capabilities": {},
                },
            }
        )
        shutdown_msg = frame_message(
            {"jsonrpc": "2.0", "id": 2, "method": "shutdown"}
        )
        reader = io.BytesIO(init_msg + shutdown_msg)
        writer = io.BytesIO()
        server = Server(reader, writer)
        server.run()

        output = writer.getvalue()
        assert len(output) > 0

        # Parse the framed responses from the output
        responses = []
        stream = io.BytesIO(output)
        while True:
            line = stream.readline()
            if not line:
                break
            line_str = line.decode("ascii").strip()
            if line_str.lower().startswith("content-length:"):
                content_length = int(line_str.split(":", 1)[1].strip())
                # Read the empty line separator
                stream.readline()
                body = stream.read(content_length)
                responses.append(json.loads(body))

        assert len(responses) == 2

        # First response is initialize
        assert responses[0]["id"] == 1
        assert "capabilities" in responses[0]["result"]
        assert responses[0]["result"]["serverInfo"]["name"] == "lsp-server"

        # Second response is shutdown
        assert responses[1]["id"] == 2
        assert responses[1]["result"] is None
