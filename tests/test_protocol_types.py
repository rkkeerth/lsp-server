"""Tests for LSP protocol types serialization and deserialization."""
from lsp_server.protocol.types import (
    Diagnostic,
    DiagnosticSeverity,
    DidChangeTextDocumentParams,
    DidCloseTextDocumentParams,
    DidOpenTextDocumentParams,
    InitializeParams,
    InitializeResult,
    Position,
    Range,
    ServerCapabilities,
    ServerInfo,
    TextDocumentSyncKind,
    TextDocumentSyncOptions,
)


class TestTextDocumentSyncKind:
    def test_none_is_zero(self):
        assert TextDocumentSyncKind.NONE == 0

    def test_full_is_one(self):
        assert TextDocumentSyncKind.FULL == 1

    def test_incremental_is_two(self):
        assert TextDocumentSyncKind.INCREMENTAL == 2


class TestDiagnosticSeverity:
    def test_error_is_one(self):
        assert DiagnosticSeverity.ERROR == 1

    def test_warning_is_two(self):
        assert DiagnosticSeverity.WARNING == 2

    def test_information_is_three(self):
        assert DiagnosticSeverity.INFORMATION == 3

    def test_hint_is_four(self):
        assert DiagnosticSeverity.HINT == 4


class TestPosition:
    def test_to_dict(self):
        pos = Position(line=5, character=10)
        assert pos.to_dict() == {"line": 5, "character": 10}

    def test_from_dict(self):
        pos = Position.from_dict({"line": 3, "character": 7})
        assert pos.line == 3
        assert pos.character == 7

    def test_round_trip(self):
        original = Position(line=12, character=4)
        restored = Position.from_dict(original.to_dict())
        assert restored.line == original.line
        assert restored.character == original.character


class TestRange:
    def test_to_dict(self):
        r = Range(
            start=Position(line=1, character=0),
            end=Position(line=1, character=10),
        )
        d = r.to_dict()
        assert d == {
            "start": {"line": 1, "character": 0},
            "end": {"line": 1, "character": 10},
        }

    def test_from_dict(self):
        r = Range.from_dict({
            "start": {"line": 2, "character": 5},
            "end": {"line": 3, "character": 0},
        })
        assert r.start.line == 2
        assert r.start.character == 5
        assert r.end.line == 3
        assert r.end.character == 0

    def test_round_trip(self):
        original = Range(
            start=Position(line=10, character=2),
            end=Position(line=10, character=20),
        )
        restored = Range.from_dict(original.to_dict())
        assert restored.start.line == original.start.line
        assert restored.start.character == original.start.character
        assert restored.end.line == original.end.line
        assert restored.end.character == original.end.character


class TestInitializeResult:
    def test_to_dict_with_capabilities(self):
        result = InitializeResult(
            capabilities=ServerCapabilities(
                text_document_sync=TextDocumentSyncOptions(
                    open_close=True,
                    change=TextDocumentSyncKind.FULL,
                )
            ),
            server_info=ServerInfo(name="test-server", version="1.0.0"),
        )
        d = result.to_dict()
        assert "capabilities" in d
        assert d["capabilities"]["textDocumentSync"]["openClose"] is True
        assert d["capabilities"]["textDocumentSync"]["change"] == 1
        assert d["serverInfo"]["name"] == "test-server"
        assert d["serverInfo"]["version"] == "1.0.0"

    def test_from_dict_round_trip(self):
        data = {
            "capabilities": {
                "textDocumentSync": {
                    "openClose": True,
                    "change": 2,
                }
            },
            "serverInfo": {
                "name": "my-server",
                "version": "0.1.0",
            },
        }
        result = InitializeResult.from_dict(data)
        assert result.capabilities.text_document_sync is not None
        assert result.capabilities.text_document_sync.open_close is True
        assert result.capabilities.text_document_sync.change == TextDocumentSyncKind.INCREMENTAL
        assert result.server_info is not None
        assert result.server_info.name == "my-server"
        assert result.server_info.version == "0.1.0"
        # Round-trip back to dict
        restored = result.to_dict()
        assert restored["capabilities"]["textDocumentSync"]["openClose"] is True
        assert restored["capabilities"]["textDocumentSync"]["change"] == 2
        assert restored["serverInfo"]["name"] == "my-server"


class TestInitializeParams:
    def test_from_dict_with_typical_client_json(self):
        data = {
            "processId": 1234,
            "rootUri": "file:///home/user/project",
            "capabilities": {
                "textDocument": {
                    "synchronization": {
                        "dynamicRegistration": True,
                        "willSave": True,
                        "didSave": True,
                    }
                },
                "workspace": {
                    "workspaceFolders": True,
                },
            },
            "trace": "verbose",
            "workspaceFolders": [
                {"uri": "file:///home/user/project", "name": "project"},
            ],
        }
        params = InitializeParams.from_dict(data)
        assert params.process_id == 1234
        assert params.root_uri == "file:///home/user/project"
        assert params.trace == "verbose"
        assert params.capabilities.text_document is not None
        assert params.capabilities.text_document.synchronization is not None
        assert params.capabilities.text_document.synchronization.dynamic_registration is True
        assert params.capabilities.workspace is not None
        assert params.capabilities.workspace.workspace_folders is True
        assert params.workspace_folders is not None
        assert len(params.workspace_folders) == 1
        assert params.workspace_folders[0].uri == "file:///home/user/project"
        assert params.workspace_folders[0].name == "project"


class TestDidOpenTextDocumentParams:
    def test_from_dict_and_to_dict_round_trip(self):
        data = {
            "textDocument": {
                "uri": "file:///test.py",
                "languageId": "python",
                "version": 1,
                "text": "print('hello')\n",
            }
        }
        params = DidOpenTextDocumentParams.from_dict(data)
        assert params.text_document.uri == "file:///test.py"
        assert params.text_document.language_id == "python"
        assert params.text_document.version == 1
        assert params.text_document.text == "print('hello')\n"
        # Round-trip
        restored = params.to_dict()
        assert restored["textDocument"]["uri"] == "file:///test.py"
        assert restored["textDocument"]["languageId"] == "python"
        assert restored["textDocument"]["version"] == 1
        assert restored["textDocument"]["text"] == "print('hello')\n"


class TestDidChangeTextDocumentParams:
    def test_from_dict_with_content_changes(self):
        data = {
            "textDocument": {
                "uri": "file:///test.py",
                "version": 2,
            },
            "contentChanges": [
                {"text": "print('world')\n"},
                {
                    "range": {
                        "start": {"line": 0, "character": 0},
                        "end": {"line": 0, "character": 5},
                    },
                    "text": "echo",
                },
            ],
        }
        params = DidChangeTextDocumentParams.from_dict(data)
        assert params.text_document.uri == "file:///test.py"
        assert params.text_document.version == 2
        assert len(params.content_changes) == 2
        assert params.content_changes[0].text == "print('world')\n"
        assert params.content_changes[0].range is None
        assert params.content_changes[1].text == "echo"
        assert params.content_changes[1].range is not None
        assert params.content_changes[1].range.start.line == 0
        assert params.content_changes[1].range.end.character == 5


class TestDidCloseTextDocumentParams:
    def test_from_dict_and_to_dict_round_trip(self):
        data = {"textDocument": {"uri": "file:///test.py"}}
        params = DidCloseTextDocumentParams.from_dict(data)
        assert params.text_document.uri == "file:///test.py"
        restored = params.to_dict()
        assert restored == data


class TestDiagnostic:
    def test_to_dict_includes_severity_value(self):
        diag = Diagnostic(
            range=Range(
                start=Position(line=0, character=0),
                end=Position(line=0, character=5),
            ),
            severity=DiagnosticSeverity.WARNING,
            message="unused variable",
            source="linter",
        )
        d = diag.to_dict()
        assert d["severity"] == 2
        assert d["message"] == "unused variable"
        assert d["source"] == "linter"
        assert d["range"]["start"]["line"] == 0

    def test_from_dict_restores_enum(self):
        data = {
            "range": {
                "start": {"line": 1, "character": 0},
                "end": {"line": 1, "character": 10},
            },
            "severity": 1,
            "message": "syntax error",
        }
        diag = Diagnostic.from_dict(data)
        assert diag.severity == DiagnosticSeverity.ERROR
        assert isinstance(diag.severity, DiagnosticSeverity)
        assert diag.message == "syntax error"


class TestTextDocumentSyncOptions:
    def test_to_dict_camel_case(self):
        opts = TextDocumentSyncOptions(
            open_close=True,
            change=TextDocumentSyncKind.FULL,
        )
        d = opts.to_dict()
        assert d["openClose"] is True
        assert d["change"] == 1
        assert "save" not in d

    def test_to_dict_with_save(self):
        from lsp_server.protocol.types import SaveOptions

        opts = TextDocumentSyncOptions(
            open_close=True,
            change=TextDocumentSyncKind.FULL,
            save=SaveOptions(include_text=True),
        )
        d = opts.to_dict()
        assert d["save"]["includeText"] is True


class TestServerInfo:
    def test_to_dict_omits_version_when_none(self):
        info = ServerInfo(name="test-server", version=None)
        d = info.to_dict()
        assert d == {"name": "test-server"}
        assert "version" not in d

    def test_to_dict_includes_version_when_set(self):
        info = ServerInfo(name="test-server", version="1.0.0")
        d = info.to_dict()
        assert d == {"name": "test-server", "version": "1.0.0"}
