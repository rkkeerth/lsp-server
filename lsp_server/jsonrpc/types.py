"""JSON-RPC 2.0 types and error codes."""
from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Optional


# JSON-RPC 2.0 error codes
PARSE_ERROR = -32700
INVALID_REQUEST = -32600
METHOD_NOT_FOUND = -32601
INVALID_PARAMS = -32602
INTERNAL_ERROR = -32603
SERVER_NOT_INITIALIZED = -32002
REQUEST_CANCELLED = -32800
CONTENT_MODIFIED = -32801


@dataclass
class ResponseError:
    code: int = 0
    message: str = ""
    data: Any = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {
            "code": self.code,
            "message": self.message,
        }
        if self.data is not None:
            d["data"] = self.data
        return d


@dataclass
class Request:
    jsonrpc: str = ""
    id: Any = None
    method: str = ""
    params: Any = None

    def is_notification(self) -> bool:
        return self.id is None


@dataclass
class Response:
    jsonrpc: str = "2.0"
    id: Any = None
    result: Any = None
    error: Optional[ResponseError] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {
            "jsonrpc": self.jsonrpc,
            "id": self.id,
        }
        if self.error is not None:
            d["error"] = self.error.to_dict()
        else:
            d["result"] = self.result
        if self.result is not None and self.error is not None:
            d["result"] = self.result
        return d


@dataclass
class Notification:
    jsonrpc: str = "2.0"
    method: str = ""
    params: Any = None


def new_response(id: Any, result: Any) -> Response:
    return Response(jsonrpc="2.0", id=id, result=result)


def new_error_response(id: Any, code: int, message: str) -> Response:
    return Response(
        jsonrpc="2.0",
        id=id,
        error=ResponseError(code=code, message=message),
    )
