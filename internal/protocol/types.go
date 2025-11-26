package protocol

// InitializeParams represents the parameters for the initialize request
type InitializeParams struct {
	ProcessID             *int                `json:"processId"`
	ClientInfo            *ClientInfo         `json:"clientInfo,omitempty"`
	Locale                string              `json:"locale,omitempty"`
	RootPath              *string             `json:"rootPath,omitempty"`
	RootURI               *string             `json:"rootUri"`
	InitializationOptions interface{}         `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities  `json:"capabilities"`
	Trace                 *string             `json:"trace,omitempty"`
	WorkspaceFolders      []WorkspaceFolder   `json:"workspaceFolders,omitempty"`
}

// ClientInfo represents information about the client
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// ClientCapabilities represents the capabilities provided by the client
type ClientCapabilities struct {
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Window       *WindowClientCapabilities       `json:"window,omitempty"`
	General      *GeneralClientCapabilities      `json:"general,omitempty"`
	Experimental interface{}                     `json:"experimental,omitempty"`
}

// WorkspaceClientCapabilities represents workspace-specific client capabilities
type WorkspaceClientCapabilities struct {
	ApplyEdit              bool                        `json:"applyEdit,omitempty"`
	WorkspaceEdit          *WorkspaceEditCapabilities  `json:"workspaceEdit,omitempty"`
	DidChangeConfiguration *DynamicRegistration        `json:"didChangeConfiguration,omitempty"`
	DidChangeWatchedFiles  *DynamicRegistration        `json:"didChangeWatchedFiles,omitempty"`
	Symbol                 *DynamicRegistration        `json:"symbol,omitempty"`
	ExecuteCommand         *DynamicRegistration        `json:"executeCommand,omitempty"`
	WorkspaceFolders       bool                        `json:"workspaceFolders,omitempty"`
	Configuration          bool                        `json:"configuration,omitempty"`
}

// WorkspaceEditCapabilities represents workspace edit capabilities
type WorkspaceEditCapabilities struct {
	DocumentChanges bool `json:"documentChanges,omitempty"`
}

// TextDocumentClientCapabilities represents text document-specific client capabilities
type TextDocumentClientCapabilities struct {
	Synchronization *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
	Completion      *DynamicRegistration                `json:"completion,omitempty"`
	Hover           *DynamicRegistration                `json:"hover,omitempty"`
	SignatureHelp   *DynamicRegistration                `json:"signatureHelp,omitempty"`
	References      *DynamicRegistration                `json:"references,omitempty"`
	DocumentHighlight *DynamicRegistration              `json:"documentHighlight,omitempty"`
	DocumentSymbol  *DynamicRegistration                `json:"documentSymbol,omitempty"`
	Formatting      *DynamicRegistration                `json:"formatting,omitempty"`
	RangeFormatting *DynamicRegistration                `json:"rangeFormatting,omitempty"`
	Definition      *DynamicRegistration                `json:"definition,omitempty"`
	CodeAction      *DynamicRegistration                `json:"codeAction,omitempty"`
}

// TextDocumentSyncClientCapabilities represents text document synchronization capabilities
type TextDocumentSyncClientCapabilities struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	WillSave            bool `json:"willSave,omitempty"`
	WillSaveWaitUntil   bool `json:"willSaveWaitUntil,omitempty"`
	DidSave             bool `json:"didSave,omitempty"`
}

// WindowClientCapabilities represents window-specific client capabilities
type WindowClientCapabilities struct {
	WorkDoneProgress bool `json:"workDoneProgress,omitempty"`
}

// GeneralClientCapabilities represents general client capabilities
type GeneralClientCapabilities struct {
	RegularExpressions *RegularExpressionsCapabilities `json:"regularExpressions,omitempty"`
	Markdown           *MarkdownCapabilities           `json:"markdown,omitempty"`
}

// RegularExpressionsCapabilities represents regular expression capabilities
type RegularExpressionsCapabilities struct {
	Engine  string `json:"engine"`
	Version string `json:"version,omitempty"`
}

// MarkdownCapabilities represents markdown capabilities
type MarkdownCapabilities struct {
	Parser  string   `json:"parser"`
	Version string   `json:"version,omitempty"`
	AllowedTags []string `json:"allowedTags,omitempty"`
}

// DynamicRegistration represents dynamic registration capability
type DynamicRegistration struct {
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// WorkspaceFolder represents a workspace folder
type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

// InitializeResult represents the result of the initialize request
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

// ServerCapabilities represents the capabilities provided by the server
type ServerCapabilities struct {
	TextDocumentSync   *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
	CompletionProvider *CompletionOptions       `json:"completionProvider,omitempty"`
	HoverProvider      bool                     `json:"hoverProvider,omitempty"`
	SignatureHelpProvider *SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`
	DefinitionProvider bool                     `json:"definitionProvider,omitempty"`
	ReferencesProvider bool                     `json:"referencesProvider,omitempty"`
	DocumentHighlightProvider bool              `json:"documentHighlightProvider,omitempty"`
	DocumentSymbolProvider bool                 `json:"documentSymbolProvider,omitempty"`
	WorkspaceSymbolProvider bool                `json:"workspaceSymbolProvider,omitempty"`
	CodeActionProvider bool                     `json:"codeActionProvider,omitempty"`
	DocumentFormattingProvider bool             `json:"documentFormattingProvider,omitempty"`
	DocumentRangeFormattingProvider bool        `json:"documentRangeFormattingProvider,omitempty"`
	RenameProvider     bool                     `json:"renameProvider,omitempty"`
}

// TextDocumentSyncOptions represents text document sync options
type TextDocumentSyncOptions struct {
	OpenClose         bool                 `json:"openClose,omitempty"`
	Change            TextDocumentSyncKind `json:"change,omitempty"`
	WillSave          bool                 `json:"willSave,omitempty"`
	WillSaveWaitUntil bool                 `json:"willSaveWaitUntil,omitempty"`
	Save              *SaveOptions         `json:"save,omitempty"`
}

// TextDocumentSyncKind defines how text documents are synced
type TextDocumentSyncKind int

const (
	// None means documents should not be synced
	None TextDocumentSyncKind = 0
	// Full means documents are synced by sending full content
	Full TextDocumentSyncKind = 1
	// Incremental means documents are synced by sending incremental changes
	Incremental TextDocumentSyncKind = 2
)

// SaveOptions represents save options
type SaveOptions struct {
	IncludeText bool `json:"includeText,omitempty"`
}

// CompletionOptions represents completion options
type CompletionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	ResolveProvider   bool     `json:"resolveProvider,omitempty"`
}

// SignatureHelpOptions represents signature help options
type SignatureHelpOptions struct {
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// ServerInfo represents information about the server
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// DidOpenTextDocumentParams represents the parameters for the textDocument/didOpen notification
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// TextDocumentItem represents a text document
type TextDocumentItem struct {
	URI        string `json:"uri"`
	LanguageID string `json:"languageId"`
	Version    int    `json:"version"`
	Text       string `json:"text"`
}

// DidChangeTextDocumentParams represents the parameters for the textDocument/didChange notification
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// VersionedTextDocumentIdentifier represents a versioned text document identifier
type VersionedTextDocumentIdentifier struct {
	URI     string `json:"uri"`
	Version int    `json:"version"`
}

// TextDocumentContentChangeEvent represents a change to a text document
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength *int   `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// Range represents a range in a text document
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Position represents a position in a text document
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// DidCloseTextDocumentParams represents the parameters for the textDocument/didClose notification
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// TextDocumentIdentifier represents a text document identifier
type TextDocumentIdentifier struct {
	URI string `json:"uri"`
}
