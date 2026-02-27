"""JSON-RPC 2.0 transport with Content-Length framing."""
from __future__ import annotations

import json
from typing import Any, BinaryIO

from lsp_server.jsonrpc.types import Request, Response


class Transport:
    """Reads and writes JSON-RPC messages over binary I/O streams
    using Content-Length framing."""

    def __init__(self, reader: BinaryIO, writer: BinaryIO) -> None:
        self._reader = reader
        self._writer = writer

    def read(self) -> Request:
        """Read a JSON-RPC request from the input stream."""
        content_length = -1

        # Read headers until empty line
        while True:
            line = self._reader.readline()
            if not line:
                raise EOFError("Unexpected end of input while reading headers")
            line_str = line.decode("ascii").strip()
            if line_str == "":
                break
            if line_str.lower().startswith("content-length:"):
                try:
                    content_length = int(line_str.split(":", 1)[1].strip())
                except ValueError:
                    raise ValueError(f"Invalid Content-Length value: {line_str}")

        if content_length < 0:
            raise ValueError("Missing Content-Length header")

        # Read exactly content_length bytes
        body = self._reader.read(content_length)
        if len(body) != content_length:
            raise ValueError(
                f"Expected {content_length} bytes, got {len(body)}"
            )

        data = json.loads(body)
        return Request(
            jsonrpc=data.get("jsonrpc", ""),
            id=data.get("id"),
            method=data.get("method", ""),
            params=data.get("params"),
        )

    def write(self, response: Response) -> None:
        """Write a JSON-RPC response to the output stream."""
        body = json.dumps(response.to_dict()).encode("utf-8")
        header = f"Content-Length: {len(body)}\r\n\r\n".encode("utf-8")
        self._writer.write(header)
        self._writer.write(body)
        self._writer.flush()

    def write_notification(self, method: str, params: Any = None) -> None:
        """Write a JSON-RPC notification to the output stream."""
        msg: dict[str, Any] = {
            "jsonrpc": "2.0",
            "method": method,
        }
        if params is not None:
            msg["params"] = params
        body = json.dumps(msg).encode("utf-8")
        header = f"Content-Length: {len(body)}\r\n\r\n".encode("utf-8")
        self._writer.write(header)
        self._writer.write(body)
        self._writer.flush()
