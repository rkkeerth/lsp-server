package protocol

// TextDocumentSyncKind defines how text documents are synced
type TextDocumentSyncKind int

const (
	TextDocumentSyncKindNone        TextDocumentSyncKind = 0
	TextDocumentSyncKindFull        TextDocumentSyncKind = 1
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

// Position in a text document
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range in a text document
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location represents a location inside a resource
type Location struct {
	URI   string `json:"uri"`
	Range Range  `json:"range"`
}

// TextDocumentIdentifier identifies a text document
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}

// VersionedTextDocumentIdentifier identifies a specific version
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// TextDocumentItem is an item to transfer a text document from client to server
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// TextDocumentPositionParams is a parameter literal used in requests
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// DiagnosticSeverity for diagnostic messages
type DiagnosticSeverity int

const (
	DiagnosticSeverityError       DiagnosticSeverity = 1
	DiagnosticSeverityWarning     DiagnosticSeverity = 2
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	DiagnosticSeverityHint        DiagnosticSeverity = 4
)

// Diagnostic represents a diagnostic like a compiler error
type Diagnostic struct {
	Range    Range               `json:"range"`
	Severity *DiagnosticSeverity `json:"severity,omitempty"`
	Code     *string             `json:"code,omitempty"`
	Source   *string             `json:"source,omitempty"`
	Message  string              `json:"message"`
}

// PublishDiagnosticsParams are the params for textDocument/publishDiagnostics
type PublishDiagnosticsParams struct {
	URI         string       `json:"uri"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// MarkupKind describes the content type
type MarkupKind string

const (
	MarkupKindPlainText MarkupKind = "plaintext"
	MarkupKindMarkdown  MarkupKind = "markdown"
)

// MarkupContent represents a string value with markup kind
type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}

// Hover is the result of a hover request
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}

// HoverParams are the params for textDocument/hover
type HoverParams struct {
	TextDocumentPositionParams
}

// CompletionItemKind describes the kind of completion item
type CompletionItemKind int

const (
	CompletionItemKindText        CompletionItemKind = 1
	CompletionItemKindMethod      CompletionItemKind = 2
	CompletionItemKindFunction    CompletionItemKind = 3
	CompletionItemKindConstructor CompletionItemKind = 4
	CompletionItemKindField       CompletionItemKind = 5
	CompletionItemKindVariable    CompletionItemKind = 6
	CompletionItemKindClass       CompletionItemKind = 7
	CompletionItemKindInterface   CompletionItemKind = 8
	CompletionItemKindModule      CompletionItemKind = 9
	CompletionItemKindProperty    CompletionItemKind = 10
	CompletionItemKindKeyword     CompletionItemKind = 14
	CompletionItemKindSnippet     CompletionItemKind = 15
)

// CompletionItem represents a completion suggestion
type CompletionItem struct {
	Label         string              `json:"label"`
	Kind          *CompletionItemKind `json:"kind,omitempty"`
	Detail        *string             `json:"detail,omitempty"`
	Documentation *string             `json:"documentation,omitempty"`
	InsertText    *string             `json:"insertText,omitempty"`
}

// CompletionList represents a collection of completion items
type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

// CompletionParams are the params for textDocument/completion
type CompletionParams struct {
	TextDocumentPositionParams
}

// CompletionOptions are server capabilities for completion
type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
}

// TextDocumentSyncOptions specify how text documents are synced
type TextDocumentSyncOptions struct {
	OpenClose bool                 `json:"openClose,omitempty"`
	Change    TextDocumentSyncKind `json:"change,omitempty"`
}

// ServerCapabilities define what the server can do
type ServerCapabilities struct {
	TextDocumentSync   *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
	CompletionProvider *CompletionOptions       `json:"completionProvider,omitempty"`
	HoverProvider      bool                     `json:"hoverProvider,omitempty"`
}

// ClientCapabilities define what the client supports
type ClientCapabilities struct {
	// Add specific client capabilities as needed
}

// InitializeParams are sent from client to server during initialize
type InitializeParams struct {
	ProcessID    *int               `json:"processId"`
	RootURI      *string            `json:"rootUri"`
	Capabilities ClientCapabilities `json:"capabilities"`
}

// InitializeResult is returned from server to client after initialize
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
}

// TextDocumentContentChangeEvent describes a change to a text document
type TextDocumentContentChangeEvent struct {
	Range *Range `json:"range,omitempty"`
	Text  string `json:"text"`
}

// DidOpenTextDocumentParams are sent when a document is opened
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidCloseTextDocumentParams are sent when a document is closed
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DidChangeTextDocumentParams are sent when a document changes
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}
