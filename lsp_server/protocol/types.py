"""LSP protocol types as dataclasses."""
from __future__ import annotations

from dataclasses import dataclass, field
from enum import IntEnum
from typing import Any, Optional


class TextDocumentSyncKind(IntEnum):
    NONE = 0
    FULL = 1
    INCREMENTAL = 2


class DiagnosticSeverity(IntEnum):
    ERROR = 1
    WARNING = 2
    INFORMATION = 3
    HINT = 4


@dataclass
class Position:
    line: int = 0
    character: int = 0

    def to_dict(self) -> dict:
        return {"line": self.line, "character": self.character}

    @classmethod
    def from_dict(cls, data: dict) -> Position:
        if not data:
            return cls()
        return cls(
            line=data.get("line", 0),
            character=data.get("character", 0),
        )


@dataclass
class Range:
    start: Position = field(default_factory=Position)
    end: Position = field(default_factory=Position)

    def to_dict(self) -> dict:
        return {"start": self.start.to_dict(), "end": self.end.to_dict()}

    @classmethod
    def from_dict(cls, data: dict) -> Range:
        if not data:
            return cls()
        return cls(
            start=Position.from_dict(data.get("start", {})),
            end=Position.from_dict(data.get("end", {})),
        )


@dataclass
class Location:
    uri: str = ""
    range: Range = field(default_factory=Range)

    def to_dict(self) -> dict:
        return {"uri": self.uri, "range": self.range.to_dict()}

    @classmethod
    def from_dict(cls, data: dict) -> Location:
        if not data:
            return cls()
        return cls(
            uri=data.get("uri", ""),
            range=Range.from_dict(data.get("range", {})),
        )


@dataclass
class SaveOptions:
    include_text: bool = False

    def to_dict(self) -> dict:
        return {"includeText": self.include_text}

    @classmethod
    def from_dict(cls, data: dict) -> SaveOptions:
        if not data:
            return cls()
        return cls(include_text=data.get("includeText", False))


@dataclass
class TextDocumentSyncOptions:
    open_close: bool = False
    change: TextDocumentSyncKind = TextDocumentSyncKind.NONE
    save: Optional[SaveOptions] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {
            "openClose": self.open_close,
            "change": self.change.value,
        }
        if self.save is not None:
            d["save"] = self.save.to_dict()
        return d

    @classmethod
    def from_dict(cls, data: dict) -> TextDocumentSyncOptions:
        if not data:
            return cls()
        save_data = data.get("save")
        return cls(
            open_close=data.get("openClose", False),
            change=TextDocumentSyncKind(data.get("change", 0)),
            save=SaveOptions.from_dict(save_data) if save_data is not None else None,
        )


@dataclass
class ServerInfo:
    name: str = ""
    version: Optional[str] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {"name": self.name}
        if self.version is not None:
            d["version"] = self.version
        return d

    @classmethod
    def from_dict(cls, data: dict) -> ServerInfo:
        if not data:
            return cls()
        return cls(
            name=data.get("name", ""),
            version=data.get("version"),
        )


@dataclass
class ServerCapabilities:
    text_document_sync: Optional[TextDocumentSyncOptions] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {}
        if self.text_document_sync is not None:
            d["textDocumentSync"] = self.text_document_sync.to_dict()
        return d

    @classmethod
    def from_dict(cls, data: dict) -> ServerCapabilities:
        if not data:
            return cls()
        tds = data.get("textDocumentSync")
        return cls(
            text_document_sync=(
                TextDocumentSyncOptions.from_dict(tds) if tds is not None else None
            ),
        )


@dataclass
class InitializeResult:
    capabilities: ServerCapabilities = field(default_factory=ServerCapabilities)
    server_info: Optional[ServerInfo] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {"capabilities": self.capabilities.to_dict()}
        if self.server_info is not None:
            d["serverInfo"] = self.server_info.to_dict()
        return d

    @classmethod
    def from_dict(cls, data: dict) -> InitializeResult:
        if not data:
            return cls()
        si = data.get("serverInfo")
        return cls(
            capabilities=ServerCapabilities.from_dict(data.get("capabilities", {})),
            server_info=ServerInfo.from_dict(si) if si is not None else None,
        )


@dataclass
class WorkspaceFolder:
    uri: str = ""
    name: str = ""

    def to_dict(self) -> dict:
        return {"uri": self.uri, "name": self.name}

    @classmethod
    def from_dict(cls, data: dict) -> WorkspaceFolder:
        if not data:
            return cls()
        return cls(uri=data.get("uri", ""), name=data.get("name", ""))


@dataclass
class TextDocumentSyncClientCapabilities:
    dynamic_registration: bool = False
    will_save: bool = False
    will_save_wait_until: bool = False
    did_save: bool = False

    def to_dict(self) -> dict:
        return {
            "dynamicRegistration": self.dynamic_registration,
            "willSave": self.will_save,
            "willSaveWaitUntil": self.will_save_wait_until,
            "didSave": self.did_save,
        }

    @classmethod
    def from_dict(cls, data: dict) -> TextDocumentSyncClientCapabilities:
        if not data:
            return cls()
        return cls(
            dynamic_registration=data.get("dynamicRegistration", False),
            will_save=data.get("willSave", False),
            will_save_wait_until=data.get("willSaveWaitUntil", False),
            did_save=data.get("didSave", False),
        )


@dataclass
class TextDocumentClientCapabilities:
    synchronization: Optional[TextDocumentSyncClientCapabilities] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {}
        if self.synchronization is not None:
            d["synchronization"] = self.synchronization.to_dict()
        return d

    @classmethod
    def from_dict(cls, data: dict) -> TextDocumentClientCapabilities:
        if not data:
            return cls()
        sync = data.get("synchronization")
        return cls(
            synchronization=(
                TextDocumentSyncClientCapabilities.from_dict(sync)
                if sync is not None
                else None
            ),
        )


@dataclass
class WorkspaceClientCapabilities:
    workspace_folders: bool = False

    def to_dict(self) -> dict:
        return {"workspaceFolders": self.workspace_folders}

    @classmethod
    def from_dict(cls, data: dict) -> WorkspaceClientCapabilities:
        if not data:
            return cls()
        return cls(workspace_folders=data.get("workspaceFolders", False))


@dataclass
class ClientCapabilities:
    text_document: Optional[TextDocumentClientCapabilities] = None
    workspace: Optional[WorkspaceClientCapabilities] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {}
        if self.text_document is not None:
            d["textDocument"] = self.text_document.to_dict()
        if self.workspace is not None:
            d["workspace"] = self.workspace.to_dict()
        return d

    @classmethod
    def from_dict(cls, data: dict) -> ClientCapabilities:
        if not data:
            return cls()
        td = data.get("textDocument")
        ws = data.get("workspace")
        return cls(
            text_document=(
                TextDocumentClientCapabilities.from_dict(td) if td is not None else None
            ),
            workspace=(
                WorkspaceClientCapabilities.from_dict(ws) if ws is not None else None
            ),
        )


@dataclass
class InitializeParams:
    process_id: Optional[int] = None
    root_uri: Optional[str] = None
    capabilities: ClientCapabilities = field(default_factory=ClientCapabilities)
    initialization_options: Any = None
    trace: Optional[str] = None
    workspace_folders: Optional[list] = None

    def to_dict(self) -> dict:
        d: dict[str, Any] = {
            "capabilities": self.capabilities.to_dict(),
        }
        if self.process_id is not None:
            d["processId"] = self.process_id
        if self.root_uri is not None:
            d["rootUri"] = self.root_uri
        if self.initialization_options is not None:
            d["initializationOptions"] = self.initialization_options
        if self.trace is not None:
            d["trace"] = self.trace
        if self.workspace_folders is not None:
            d["workspaceFolders"] = [
                wf.to_dict() if hasattr(wf, "to_dict") else wf
                for wf in self.workspace_folders
            ]
        return d

    @classmethod
    def from_dict(cls, data: dict) -> InitializeParams:
        if not data:
            return cls()
        wf_data = data.get("workspaceFolders")
        workspace_folders = None
        if wf_data is not None:
            workspace_folders = [WorkspaceFolder.from_dict(wf) for wf in wf_data]
        return cls(
            process_id=data.get("processId"),
            root_uri=data.get("rootUri"),
            capabilities=ClientCapabilities.from_dict(data.get("capabilities", {})),
            initialization_options=data.get("initializationOptions"),
            trace=data.get("trace"),
            workspace_folders=workspace_folders,
        )


@dataclass
class TextDocumentIdentifier:
    uri: str = ""

    def to_dict(self) -> dict:
        return {"uri": self.uri}

    @classmethod
    def from_dict(cls, data: dict) -> TextDocumentIdentifier:
        if not data:
            return cls()
        return cls(uri=data.get("uri", ""))


@dataclass
class VersionedTextDocumentIdentifier:
    uri: str = ""
    version: int = 0

    def to_dict(self) -> dict:
        return {"uri": self.uri, "version": self.version}

    @classmethod
    def from_dict(cls, data: dict) -> VersionedTextDocumentIdentifier:
        if not data:
            return cls()
        return cls(
            uri=data.get("uri", ""),
            version=data.get("version", 0),
        )


@dataclass
class TextDocumentItem:
    uri: str = ""
    language_id: str = ""
    version: int = 0
    text: str = ""

    def to_dict(self) -> dict:
        return {
            "uri": self.uri,
            "languageId": self.language_id,
            "version": self.version,
            "text": self.text,
        }

    @classmethod
    def from_dict(cls, data: dict) -> TextDocumentItem:
        if not data:
            return cls()
        return cls(
            uri=data.get("uri", ""),
            language_id=data.get("languageId", ""),
            version=data.get("version", 0),
            text=data.get("text", ""),
        )


@dataclass
class TextDocumentContentChangeEvent:
    range: Optional[Range] = None
    range_length: int = 0
    text: str = ""

    def to_dict(self) -> dict:
        d: dict[str, Any] = {"text": self.text}
        if self.range is not None:
            d["range"] = self.range.to_dict()
        if self.range_length != 0:
            d["rangeLength"] = self.range_length
        return d

    @classmethod
    def from_dict(cls, data: dict) -> TextDocumentContentChangeEvent:
        if not data:
            return cls()
        r = data.get("range")
        return cls(
            range=Range.from_dict(r) if r is not None else None,
            range_length=data.get("rangeLength", 0),
            text=data.get("text", ""),
        )


@dataclass
class Diagnostic:
    range: Range = field(default_factory=Range)
    severity: Optional[DiagnosticSeverity] = None
    code: Any = None
    source: Optional[str] = None
    message: str = ""

    def to_dict(self) -> dict:
        d: dict[str, Any] = {
            "range": self.range.to_dict(),
            "message": self.message,
        }
        if self.severity is not None:
            d["severity"] = self.severity.value
        if self.code is not None:
            d["code"] = self.code
        if self.source is not None:
            d["source"] = self.source
        return d

    @classmethod
    def from_dict(cls, data: dict) -> Diagnostic:
        if not data:
            return cls()
        sev = data.get("severity")
        return cls(
            range=Range.from_dict(data.get("range", {})),
            severity=DiagnosticSeverity(sev) if sev is not None else None,
            code=data.get("code"),
            source=data.get("source"),
            message=data.get("message", ""),
        )


@dataclass
class PublishDiagnosticsParams:
    uri: str = ""
    version: Optional[int] = None
    diagnostics: list = field(default_factory=list)

    def to_dict(self) -> dict:
        d: dict[str, Any] = {
            "uri": self.uri,
            "diagnostics": [
                diag.to_dict() if hasattr(diag, "to_dict") else diag
                for diag in self.diagnostics
            ],
        }
        if self.version is not None:
            d["version"] = self.version
        return d

    @classmethod
    def from_dict(cls, data: dict) -> PublishDiagnosticsParams:
        if not data:
            return cls()
        return cls(
            uri=data.get("uri", ""),
            version=data.get("version"),
            diagnostics=[
                Diagnostic.from_dict(d) for d in data.get("diagnostics", [])
            ],
        )


@dataclass
class DidOpenTextDocumentParams:
    text_document: TextDocumentItem = field(default_factory=TextDocumentItem)

    def to_dict(self) -> dict:
        return {"textDocument": self.text_document.to_dict()}

    @classmethod
    def from_dict(cls, data: dict) -> DidOpenTextDocumentParams:
        if not data:
            return cls()
        return cls(
            text_document=TextDocumentItem.from_dict(data.get("textDocument", {})),
        )


@dataclass
class DidChangeTextDocumentParams:
    text_document: VersionedTextDocumentIdentifier = field(
        default_factory=VersionedTextDocumentIdentifier
    )
    content_changes: list = field(default_factory=list)

    def to_dict(self) -> dict:
        return {
            "textDocument": self.text_document.to_dict(),
            "contentChanges": [
                cc.to_dict() if hasattr(cc, "to_dict") else cc
                for cc in self.content_changes
            ],
        }

    @classmethod
    def from_dict(cls, data: dict) -> DidChangeTextDocumentParams:
        if not data:
            return cls()
        return cls(
            text_document=VersionedTextDocumentIdentifier.from_dict(
                data.get("textDocument", {})
            ),
            content_changes=[
                TextDocumentContentChangeEvent.from_dict(cc)
                for cc in data.get("contentChanges", [])
            ],
        )


@dataclass
class DidCloseTextDocumentParams:
    text_document: TextDocumentIdentifier = field(
        default_factory=TextDocumentIdentifier
    )

    def to_dict(self) -> dict:
        return {"textDocument": self.text_document.to_dict()}

    @classmethod
    def from_dict(cls, data: dict) -> DidCloseTextDocumentParams:
        if not data:
            return cls()
        return cls(
            text_document=TextDocumentIdentifier.from_dict(
                data.get("textDocument", {})
            ),
        )
