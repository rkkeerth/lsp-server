// Package protocol defines LSP protocol types.
package protocol

// InitializeParams represents the parameters for the initialize request.
type InitializeParams struct {
	ProcessID             *int               `json:"processId"`
	RootURI               *string            `json:"rootUri"`
	Capabilities          ClientCapabilities `json:"capabilities"`
	InitializationOptions interface{}        `json:"initializationOptions,omitempty"`
	Trace                 string             `json:"trace,omitempty"`
	WorkspaceFolders      []WorkspaceFolder  `json:"workspaceFolders,omitempty"`
}

// ClientCapabilities represents the capabilities provided by the client.
type ClientCapabilities struct {
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
}

// TextDocumentClientCapabilities represents text document client capabilities.
type TextDocumentClientCapabilities struct {
	Synchronization *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
}

// TextDocumentSyncClientCapabilities represents sync client capabilities.
type TextDocumentSyncClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	WillSave            bool `json:"willSave,omitempty"`
	WillSaveWaitUntil   bool `json:"willSaveWaitUntil,omitempty"`
	DidSave             bool `json:"didSave,omitempty"`
}

// WorkspaceClientCapabilities represents workspace client capabilities.
type WorkspaceClientCapabilities struct {
	WorkspaceFolders bool `json:"workspaceFolders,omitempty"`
}

// WorkspaceFolder represents a workspace folder.
type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// InitializeResult represents the result of the initialize request.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

// ServerCapabilities represents the capabilities provided by the server.
type ServerCapabilities struct {
	TextDocumentSync *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
}

// TextDocumentSyncOptions represents text document sync options.
type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose,omitempty"`
	Change    TextDocumentSyncKind `json:"change,omitempty"`
	Save      *SaveOptions         `json:"save,omitempty"`
}

// TextDocumentSyncKind defines how the document is synced.
type TextDocumentSyncKind int

const (
	// SyncNone means documents should not be synced at all.
	SyncNone TextDocumentSyncKind = 0
	// SyncFull means documents are synced by sending the full content.
	SyncFull TextDocumentSyncKind = 1
	// SyncIncremental means documents are synced by sending incremental updates.
	SyncIncremental TextDocumentSyncKind = 2
)

// SaveOptions represents save options.
type SaveOptions struct {
	IncludeText bool `json:"includeText,omitempty"`
}

// ServerInfo represents information about the server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// TextDocumentIdentifier identifies a text document.
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// VersionedTextDocumentIdentifier identifies a specific version of a text document.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// TextDocumentItem represents a text document.
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// TextDocumentContentChangeEvent represents a change to a text document.
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength int    `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// Position represents a position in a text document.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range represents a range in a text document.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location represents a location inside a resource.
type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// Diagnostic represents a diagnostic (error, warning, etc.).
type Diagnostic struct {
	Range    Range            `json:"range"`
	Severity DiagnosticSeverity `json:"severity,omitempty"`
	Code     interface{}      `json:"code,omitempty"`
	Source   string           `json:"source,omitempty"`
	Message  string           `json:"message"`
}

// DiagnosticSeverity represents the severity of a diagnostic.
type DiagnosticSeverity int

const (
	// SeverityError reports an error.
	SeverityError DiagnosticSeverity = 1
	// SeverityWarning reports a warning.
	SeverityWarning DiagnosticSeverity = 2
	// SeverityInformation reports an information.
	SeverityInformation DiagnosticSeverity = 3
	// SeverityHint reports a hint.
	SeverityHint DiagnosticSeverity = 4
)

// PublishDiagnosticsParams represents the parameters for the publishDiagnostics notification.
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Version     *int         `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// DidOpenTextDocumentParams represents the parameters for the textDocument/didOpen notification.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidChangeTextDocumentParams represents the parameters for the textDocument/didChange notification.
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// DidCloseTextDocumentParams represents the parameters for the textDocument/didClose notification.
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}
