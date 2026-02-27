"""Tests for JSON-RPC 2.0 types and error codes."""
from lsp_server.jsonrpc.types import (
    CONTENT_MODIFIED,
    INTERNAL_ERROR,
    INVALID_PARAMS,
    INVALID_REQUEST,
    METHOD_NOT_FOUND,
    PARSE_ERROR,
    REQUEST_CANCELLED,
    SERVER_NOT_INITIALIZED,
    Request,
    Response,
    ResponseError,
    new_error_response,
    new_response,
)


class TestRequest:
    def test_is_notification_true_when_id_none(self):
        req = Request(jsonrpc="2.0", id=None, method="textDocument/didOpen")
        assert req.is_notification() is True

    def test_is_notification_false_when_id_set(self):
        req = Request(jsonrpc="2.0", id=1, method="initialize")
        assert req.is_notification() is False


class TestNewResponse:
    def test_creates_correct_response(self):
        resp = new_response(42, {"capabilities": {}})
        assert resp.jsonrpc == "2.0"
        assert resp.id == 42
        assert resp.result == {"capabilities": {}}
        assert resp.error is None


class TestNewErrorResponse:
    def test_creates_correct_error_response(self):
        resp = new_error_response(7, METHOD_NOT_FOUND, "Method not found")
        assert resp.jsonrpc == "2.0"
        assert resp.id == 7
        assert resp.error is not None
        assert resp.error.code == METHOD_NOT_FOUND
        assert resp.error.message == "Method not found"
        assert resp.result is None


class TestResponseToDict:
    def test_successful_response_includes_result(self):
        resp = new_response(1, {"key": "value"})
        d = resp.to_dict()
        assert d["jsonrpc"] == "2.0"
        assert d["id"] == 1
        assert d["result"] == {"key": "value"}
        assert "error" not in d

    def test_error_response_includes_error_omits_result(self):
        resp = new_error_response(1, PARSE_ERROR, "Parse error")
        d = resp.to_dict()
        assert d["jsonrpc"] == "2.0"
        assert d["id"] == 1
        assert "error" in d
        assert d["error"]["code"] == PARSE_ERROR
        assert d["error"]["message"] == "Parse error"
        assert "result" not in d

    def test_none_result_includes_result_null(self):
        resp = new_response(1, None)
        d = resp.to_dict()
        assert "result" in d
        assert d["result"] is None


class TestResponseError:
    def test_to_dict_omits_data_when_none(self):
        err = ResponseError(code=INTERNAL_ERROR, message="Internal error", data=None)
        d = err.to_dict()
        assert d["code"] == INTERNAL_ERROR
        assert d["message"] == "Internal error"
        assert "data" not in d

    def test_to_dict_includes_data_when_set(self):
        err = ResponseError(code=INTERNAL_ERROR, message="Internal error", data={"detail": "stack trace"})
        d = err.to_dict()
        assert d["code"] == INTERNAL_ERROR
        assert d["message"] == "Internal error"
        assert d["data"] == {"detail": "stack trace"}


class TestErrorCodes:
    def test_parse_error(self):
        assert PARSE_ERROR == -32700

    def test_invalid_request(self):
        assert INVALID_REQUEST == -32600

    def test_method_not_found(self):
        assert METHOD_NOT_FOUND == -32601

    def test_invalid_params(self):
        assert INVALID_PARAMS == -32602

    def test_internal_error(self):
        assert INTERNAL_ERROR == -32603

    def test_server_not_initialized(self):
        assert SERVER_NOT_INITIALIZED == -32002

    def test_request_cancelled(self):
        assert REQUEST_CANCELLED == -32800

    def test_content_modified(self):
        assert CONTENT_MODIFIED == -32801
