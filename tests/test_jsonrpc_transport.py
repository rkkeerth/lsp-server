"""Tests for JSON-RPC 2.0 transport with Content-Length framing."""
import io
import json

import pytest

from lsp_server.jsonrpc.transport import Transport
from lsp_server.jsonrpc.types import new_response


def _frame(obj: dict) -> bytes:
    """Encode a dict as a Content-Length framed message."""
    body = json.dumps(obj).encode("utf-8")
    header = f"Content-Length: {len(body)}\r\n\r\n".encode("utf-8")
    return header + body


class TestTransportRead:
    def test_valid_message_returns_correct_request(self):
        msg = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {"capabilities": {}},
        }
        reader = io.BytesIO(_frame(msg))
        writer = io.BytesIO()
        transport = Transport(reader, writer)

        req = transport.read()
        assert req.jsonrpc == "2.0"
        assert req.id == 1
        assert req.method == "initialize"
        assert req.params == {"capabilities": {}}

    def test_missing_content_length_raises_value_error(self):
        # Send a message without Content-Length header (just blank line then body)
        raw = b"SomeHeader: value\r\n\r\n{}"
        reader = io.BytesIO(raw)
        writer = io.BytesIO()
        transport = Transport(reader, writer)

        with pytest.raises(ValueError, match="Missing Content-Length"):
            transport.read()

    def test_empty_input_raises_eof_error(self):
        reader = io.BytesIO(b"")
        writer = io.BytesIO()
        transport = Transport(reader, writer)

        with pytest.raises(EOFError):
            transport.read()

    def test_read_multiple_sequential_messages(self):
        msg1 = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": None,
        }
        msg2 = {
            "jsonrpc": "2.0",
            "id": 2,
            "method": "shutdown",
        }
        raw = _frame(msg1) + _frame(msg2)
        reader = io.BytesIO(raw)
        writer = io.BytesIO()
        transport = Transport(reader, writer)

        req1 = transport.read()
        assert req1.id == 1
        assert req1.method == "initialize"

        req2 = transport.read()
        assert req2.id == 2
        assert req2.method == "shutdown"


class TestTransportWrite:
    def test_writes_correct_content_length_and_json_body(self):
        reader = io.BytesIO()
        writer = io.BytesIO()
        transport = Transport(reader, writer)

        resp = new_response(1, {"result": "ok"})
        transport.write(resp)

        output = writer.getvalue()
        # Parse the header and body
        header_end = output.index(b"\r\n\r\n")
        header = output[:header_end].decode("utf-8")
        body = output[header_end + 4:]

        assert header.startswith("Content-Length: ")
        content_length = int(header.split(": ", 1)[1])
        assert content_length == len(body)

        parsed = json.loads(body)
        assert parsed["jsonrpc"] == "2.0"
        assert parsed["id"] == 1
        assert parsed["result"] == {"result": "ok"}


class TestTransportWriteNotification:
    def test_write_notification_without_params(self):
        reader = io.BytesIO()
        writer = io.BytesIO()
        transport = Transport(reader, writer)

        transport.write_notification("window/logMessage")

        output = writer.getvalue()
        header_end = output.index(b"\r\n\r\n")
        body = output[header_end + 4:]
        parsed = json.loads(body)

        assert parsed["jsonrpc"] == "2.0"
        assert parsed["method"] == "window/logMessage"
        assert "params" not in parsed

    def test_write_notification_with_params(self):
        reader = io.BytesIO()
        writer = io.BytesIO()
        transport = Transport(reader, writer)

        transport.write_notification("window/logMessage", {"type": 3, "message": "hello"})

        output = writer.getvalue()
        header_end = output.index(b"\r\n\r\n")
        body = output[header_end + 4:]
        parsed = json.loads(body)

        assert parsed["jsonrpc"] == "2.0"
        assert parsed["method"] == "window/logMessage"
        assert parsed["params"] == {"type": 3, "message": "hello"}
