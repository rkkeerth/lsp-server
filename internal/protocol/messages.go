// Package protocol defines the LSP (Language Server Protocol) message types.
package protocol

import "encoding/json"

// Message is the base interface for all JSON-RPC messages.
type Message interface {
	isMessage()
}

// RequestMessage represents a JSON-RPC request message.
type RequestMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func (RequestMessage) isMessage() {}

// ResponseMessage represents a JSON-RPC response message.
type ResponseMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

func (ResponseMessage) isMessage() {}

// NotificationMessage represents a JSON-RPC notification message.
type NotificationMessage struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func (NotificationMessage) isMessage() {}

// ResponseError represents a JSON-RPC error.
type ResponseError struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorCode represents a JSON-RPC error code.
type ErrorCode int

const (
	// ParseError is returned when invalid JSON was received by the server.
	ParseError ErrorCode = -32700
	// InvalidRequest is returned when the JSON sent is not a valid Request object.
	InvalidRequest ErrorCode = -32600
	// MethodNotFound is returned when the method does not exist / is not available.
	MethodNotFound ErrorCode = -32601
	// InvalidParams is returned when invalid method parameter(s) are used.
	InvalidParams ErrorCode = -32602
	// InternalError is returned when there is an internal JSON-RPC error.
	InternalError ErrorCode = -32603

	// LSP specific error codes
	// ServerNotInitialized is returned when a request is sent before the server is initialized.
	ServerNotInitialized ErrorCode = -32002
	// UnknownErrorCode is used for errors that don't have a specific code.
	UnknownErrorCode ErrorCode = -32001
	// RequestCancelled is returned when a request is cancelled.
	RequestCancelled ErrorCode = -32800
	// ContentModified is returned when the content that the request was made on has been modified.
	ContentModified ErrorCode = -32801
)

// InitializeParams represents the parameters sent in an initialize request.
type InitializeParams struct {
	// The process ID of the parent process that started the server.
	ProcessID *int32 `json:"processId"`
	// The rootUri of the workspace.
	RootURI DocumentURI `json:"rootUri,omitempty"`
	// The rootPath of the workspace (deprecated, use rootUri).
	RootPath string `json:"rootPath,omitempty"`
	// User provided initialization options.
	InitializationOptions interface{} `json:"initializationOptions,omitempty"`
	// The capabilities provided by the client (editor or tool).
	Capabilities ClientCapabilities `json:"capabilities"`
	// The workspace folders configured in the client when the server starts.
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders,omitempty"`
}

// ClientCapabilities defines capabilities for dynamic registration, workspace and text document features.
type ClientCapabilities struct {
	// Workspace specific client capabilities.
	Workspace *WorkspaceClientCapabilities `json:"workspace,omitempty"`
	// Text document specific client capabilities.
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	// Experimental client capabilities.
	Experimental interface{} `json:"experimental,omitempty"`
}

// WorkspaceClientCapabilities defines capabilities the editor/tool provides on the workspace.
type WorkspaceClientCapabilities struct {
	// The client supports applying batch edits to the workspace.
	ApplyEdit bool `json:"applyEdit,omitempty"`
	// Capabilities specific to `WorkspaceEdit`s.
	WorkspaceEdit *WorkspaceEditClientCapabilities `json:"workspaceEdit,omitempty"`
	// Capabilities specific to the `workspace/didChangeConfiguration` notification.
	DidChangeConfiguration *DidChangeConfigurationClientCapabilities `json:"didChangeConfiguration,omitempty"`
	// Capabilities specific to the `workspace/didChangeWatchedFiles` notification.
	DidChangeWatchedFiles *DidChangeWatchedFilesClientCapabilities `json:"didChangeWatchedFiles,omitempty"`
	// Capabilities specific to the `workspace/symbol` request.
	Symbol *WorkspaceSymbolClientCapabilities `json:"symbol,omitempty"`
	// Capabilities specific to the `workspace/executeCommand` request.
	ExecuteCommand *ExecuteCommandClientCapabilities `json:"executeCommand,omitempty"`
	// The client has support for workspace folders.
	WorkspaceFolders bool `json:"workspaceFolders,omitempty"`
	// The client supports `workspace/configuration` requests.
	Configuration bool `json:"configuration,omitempty"`
}

// WorkspaceEditClientCapabilities defines capabilities specific to WorkspaceEdit.
type WorkspaceEditClientCapabilities struct {
	// The client supports versioned document changes.
	DocumentChanges bool `json:"documentChanges,omitempty"`
}

// DidChangeConfigurationClientCapabilities defines capabilities for configuration changes.
type DidChangeConfigurationClientCapabilities struct {
	// Did change configuration notification supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DidChangeWatchedFilesClientCapabilities defines capabilities for watched file changes.
type DidChangeWatchedFilesClientCapabilities struct {
	// Did change watched files notification supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// WorkspaceSymbolClientCapabilities defines capabilities for workspace symbols.
type WorkspaceSymbolClientCapabilities struct {
	// Symbol request supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// ExecuteCommandClientCapabilities defines capabilities for command execution.
type ExecuteCommandClientCapabilities struct {
	// Execute command supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// TextDocumentClientCapabilities defines capabilities specific to text documents.
type TextDocumentClientCapabilities struct {
	// Capabilities specific to the `textDocument/synchronization`.
	Synchronization *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
	// Capabilities specific to the `textDocument/completion`.
	Completion *CompletionClientCapabilities `json:"completion,omitempty"`
	// Capabilities specific to the `textDocument/hover`.
	Hover *HoverClientCapabilities `json:"hover,omitempty"`
	// Capabilities specific to the `textDocument/signatureHelp`.
	SignatureHelp *SignatureHelpClientCapabilities `json:"signatureHelp,omitempty"`
	// Capabilities specific to the `textDocument/definition`.
	Definition *DefinitionClientCapabilities `json:"definition,omitempty"`
	// Capabilities specific to the `textDocument/references`.
	References *ReferenceClientCapabilities `json:"references,omitempty"`
	// Capabilities specific to the `textDocument/documentHighlight`.
	DocumentHighlight *DocumentHighlightClientCapabilities `json:"documentHighlight,omitempty"`
	// Capabilities specific to the `textDocument/documentSymbol`.
	DocumentSymbol *DocumentSymbolClientCapabilities `json:"documentSymbol,omitempty"`
	// Capabilities specific to the `textDocument/codeAction`.
	CodeAction *CodeActionClientCapabilities `json:"codeAction,omitempty"`
	// Capabilities specific to the `textDocument/codeLens`.
	CodeLens *CodeLensClientCapabilities `json:"codeLens,omitempty"`
	// Capabilities specific to the `textDocument/formatting`.
	Formatting *DocumentFormattingClientCapabilities `json:"formatting,omitempty"`
	// Capabilities specific to the `textDocument/rangeFormatting`.
	RangeFormatting *DocumentRangeFormattingClientCapabilities `json:"rangeFormatting,omitempty"`
	// Capabilities specific to the `textDocument/onTypeFormatting`.
	OnTypeFormatting *DocumentOnTypeFormattingClientCapabilities `json:"onTypeFormatting,omitempty"`
	// Capabilities specific to the `textDocument/rename`.
	Rename *RenameClientCapabilities `json:"rename,omitempty"`
	// Capabilities specific to the `textDocument/publishDiagnostics`.
	PublishDiagnostics *PublishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`
}

// TextDocumentSyncClientCapabilities defines synchronization capabilities.
type TextDocumentSyncClientCapabilities struct {
	// Whether text document synchronization supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	// The client supports sending will save notifications.
	WillSave bool `json:"willSave,omitempty"`
	// The client supports sending a will save request.
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
	// The client supports did save notifications.
	DidSave bool `json:"didSave,omitempty"`
}

// CompletionClientCapabilities defines completion capabilities.
type CompletionClientCapabilities struct {
	// Whether completion supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	// The client supports the following `CompletionItem` specific capabilities.
	CompletionItem *CompletionItemClientCapabilities `json:"completionItem,omitempty"`
}

// CompletionItemClientCapabilities defines completion item capabilities.
type CompletionItemClientCapabilities struct {
	// The client supports snippets as insert text.
	SnippetSupport bool `json:"snippetSupport,omitempty"`
	// The client supports commit characters on a completion item.
	CommitCharactersSupport bool `json:"commitCharactersSupport,omitempty"`
	// The client supports the following content formats for the documentation property.
	DocumentationFormat []MarkupKind `json:"documentationFormat,omitempty"`
	// The client supports the deprecated property on a completion item.
	DeprecatedSupport bool `json:"deprecatedSupport,omitempty"`
	// The client supports the preselect property on a completion item.
	PreselectSupport bool `json:"preselectSupport,omitempty"`
}

// HoverClientCapabilities defines hover capabilities.
type HoverClientCapabilities struct {
	// Whether hover supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	// The client supports the following content formats for the content property.
	ContentFormat []MarkupKind `json:"contentFormat,omitempty"`
}

// SignatureHelpClientCapabilities defines signature help capabilities.
type SignatureHelpClientCapabilities struct {
	// Whether signature help supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DefinitionClientCapabilities defines definition capabilities.
type DefinitionClientCapabilities struct {
	// Whether definition supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// ReferenceClientCapabilities defines reference capabilities.
type ReferenceClientCapabilities struct {
	// Whether references supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentHighlightClientCapabilities defines document highlight capabilities.
type DocumentHighlightClientCapabilities struct {
	// Whether document highlight supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentSymbolClientCapabilities defines document symbol capabilities.
type DocumentSymbolClientCapabilities struct {
	// Whether document symbol supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	// Specific capabilities for the `SymbolKind`.
	SymbolKind *SymbolKindCapabilities `json:"symbolKind,omitempty"`
}

// SymbolKindCapabilities defines symbol kind capabilities.
type SymbolKindCapabilities struct {
	// The symbol kind values the client supports.
	ValueSet []SymbolKind `json:"valueSet,omitempty"`
}

// CodeActionClientCapabilities defines code action capabilities.
type CodeActionClientCapabilities struct {
	// Whether code action supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// CodeLensClientCapabilities defines code lens capabilities.
type CodeLensClientCapabilities struct {
	// Whether code lens supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentFormattingClientCapabilities defines formatting capabilities.
type DocumentFormattingClientCapabilities struct {
	// Whether formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentRangeFormattingClientCapabilities defines range formatting capabilities.
type DocumentRangeFormattingClientCapabilities struct {
	// Whether range formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// DocumentOnTypeFormattingClientCapabilities defines on type formatting capabilities.
type DocumentOnTypeFormattingClientCapabilities struct {
	// Whether on type formatting supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
}

// RenameClientCapabilities defines rename capabilities.
type RenameClientCapabilities struct {
	// Whether rename supports dynamic registration.
	DynamicRegistration bool `json:"dynamicRegistration,omitempty"`
	// The client supports testing for validity of rename operations before execution.
	PrepareSupport bool `json:"prepareSupport,omitempty"`
}

// PublishDiagnosticsClientCapabilities defines publish diagnostics capabilities.
type PublishDiagnosticsClientCapabilities struct {
	// Whether the clients accepts diagnostics with related information.
	RelatedInformation bool `json:"relatedInformation,omitempty"`
	// Client supports a codeDescription property.
	CodeDescriptionSupport bool `json:"codeDescriptionSupport,omitempty"`
	// Whether the client interprets the version property of the `textDocument/publishDiagnostics` notification.
	VersionSupport bool `json:"versionSupport,omitempty"`
}

// WorkspaceFolder represents a workspace folder.
type WorkspaceFolder struct {
	// The associated URI for this workspace folder.
	URI DocumentURI `json:"uri"`
	// The name of the workspace folder.
	Name string `json:"name"`
}

// InitializeResult represents the result of an initialize request.
type InitializeResult struct {
	// The capabilities the language server provides.
	Capabilities ServerCapabilities `json:"capabilities"`
	// Information about the server.
	ServerInfo *ServerInfo `json:"serverInfo,omitempty"`
}

// ServerInfo contains information about the server.
type ServerInfo struct {
	// The name of the server as defined by the server.
	Name string `json:"name"`
	// The server's version as defined by the server.
	Version string `json:"version,omitempty"`
}

// ServerCapabilities defines the capabilities provided by a language server.
type ServerCapabilities struct {
	// Defines how text documents are synced.
	TextDocumentSync interface{} `json:"textDocumentSync,omitempty"`
	// The server provides completion support.
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`
	// The server provides hover support.
	HoverProvider bool `json:"hoverProvider,omitempty"`
	// The server provides signature help support.
	SignatureHelpProvider *SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`
	// The server provides go to definition support.
	DefinitionProvider bool `json:"definitionProvider,omitempty"`
	// The server provides find references support.
	ReferencesProvider bool `json:"referencesProvider,omitempty"`
	// The server provides document highlight support.
	DocumentHighlightProvider bool `json:"documentHighlightProvider,omitempty"`
	// The server provides document symbol support.
	DocumentSymbolProvider bool `json:"documentSymbolProvider,omitempty"`
	// The server provides workspace symbol support.
	WorkspaceSymbolProvider bool `json:"workspaceSymbolProvider,omitempty"`
	// The server provides code actions.
	CodeActionProvider interface{} `json:"codeActionProvider,omitempty"`
	// The server provides code lens.
	CodeLensProvider *CodeLensOptions `json:"codeLensProvider,omitempty"`
	// The server provides document formatting.
	DocumentFormattingProvider bool `json:"documentFormattingProvider,omitempty"`
	// The server provides document range formatting.
	DocumentRangeFormattingProvider bool `json:"documentRangeFormattingProvider,omitempty"`
	// The server provides document formatting on typing.
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	// The server provides rename support.
	RenameProvider interface{} `json:"renameProvider,omitempty"`
	// The server provides document link support.
	DocumentLinkProvider *DocumentLinkOptions `json:"documentLinkProvider,omitempty"`
	// The server provides execute command support.
	ExecuteCommandProvider *ExecuteCommandOptions `json:"executeCommandProvider,omitempty"`
	// Workspace specific server capabilities.
	Workspace *ServerWorkspaceCapabilities `json:"workspace,omitempty"`
	// Experimental server capabilities.
	Experimental interface{} `json:"experimental,omitempty"`
}

// TextDocumentSyncKind defines how the host (editor) should sync document changes to the language server.
type TextDocumentSyncKind int

const (
	// TextDocumentSyncKindNone means documents should not be synced at all.
	TextDocumentSyncKindNone TextDocumentSyncKind = 0
	// TextDocumentSyncKindFull means documents are synced by always sending the full content.
	TextDocumentSyncKindFull TextDocumentSyncKind = 1
	// TextDocumentSyncKindIncremental means documents are synced by sending the full content on open, and incremental updates.
	TextDocumentSyncKindIncremental TextDocumentSyncKind = 2
)

// TextDocumentSyncOptions represents text document sync options.
type TextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server.
	OpenClose bool `json:"openClose,omitempty"`
	// Change notifications are sent to the server.
	Change TextDocumentSyncKind `json:"change,omitempty"`
	// Will save notifications are sent to the server.
	WillSave bool `json:"willSave,omitempty"`
	// Will save wait until requests are sent to the server.
	WillSaveWaitUntil bool `json:"willSaveWaitUntil,omitempty"`
	// Save notifications are sent to the server.
	Save *SaveOptions `json:"save,omitempty"`
}

// SaveOptions represents save options.
type SaveOptions struct {
	// The client is supposed to include the content on save.
	IncludeText bool `json:"includeText,omitempty"`
}

// CompletionOptions represents completion options.
type CompletionOptions struct {
	// The server provides support to resolve additional information for a completion item.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
	// The characters that trigger completion automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// SignatureHelpOptions represents signature help options.
type SignatureHelpOptions struct {
	// The characters that trigger signature help automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
}

// CodeLensOptions represents code lens options.
type CodeLensOptions struct {
	// Code lens has a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

// DocumentOnTypeFormattingOptions represents document on type formatting options.
type DocumentOnTypeFormattingOptions struct {
	// A character on which formatting should be triggered.
	FirstTriggerCharacter string `json:"firstTriggerCharacter"`
	// More trigger characters.
	MoreTriggerCharacter []string `json:"moreTriggerCharacter,omitempty"`
}

// DocumentLinkOptions represents document link options.
type DocumentLinkOptions struct {
	// Document links have a resolve provider as well.
	ResolveProvider bool `json:"resolveProvider,omitempty"`
}

// ExecuteCommandOptions represents execute command options.
type ExecuteCommandOptions struct {
	// The commands to be executed on the server.
	Commands []string `json:"commands,omitempty"`
}

// ServerWorkspaceCapabilities represents workspace specific server capabilities.
type ServerWorkspaceCapabilities struct {
	// The server supports workspace folders.
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
}

// WorkspaceFoldersServerCapabilities represents workspace folders server capabilities.
type WorkspaceFoldersServerCapabilities struct {
	// The server has support for workspace folders.
	Supported bool `json:"supported,omitempty"`
	// Whether the server wants to receive workspace folder change notifications.
	ChangeNotifications interface{} `json:"changeNotifications,omitempty"`
}

// DidOpenTextDocumentParams represents the parameters sent in a textDocument/didOpen notification.
type DidOpenTextDocumentParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidCloseTextDocumentParams represents the parameters sent in a textDocument/didClose notification.
type DidCloseTextDocumentParams struct {
	// The document that was closed.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DidChangeTextDocumentParams represents the parameters sent in a textDocument/didChange notification.
type DidChangeTextDocumentParams struct {
	// The document that did change.
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`
	// The actual content changes.
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// DidSaveTextDocumentParams represents the parameters sent in a textDocument/didSave notification.
type DidSaveTextDocumentParams struct {
	// The document that was saved.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// Optional the content when saved.
	Text string `json:"text,omitempty"`
}

// PublishDiagnosticsParams represents the parameters sent in a textDocument/publishDiagnostics notification.
type PublishDiagnosticsParams struct {
	// The URI for which diagnostic information is reported.
	URI DocumentURI `json:"uri"`
	// Optional the version number of the document the diagnostics are published for.
	Version int32 `json:"version,omitempty"`
	// An array of diagnostic information items.
	Diagnostics []Diagnostic `json:"diagnostics"`
}

// CompletionParams represents the parameters for a textDocument/completion request.
type CompletionParams struct {
	TextDocumentPositionParams
	// The completion context.
	Context *CompletionContext `json:"context,omitempty"`
}

// CompletionContext contains additional information about the context in which a completion request is triggered.
type CompletionContext struct {
	// How the completion was triggered.
	TriggerKind CompletionTriggerKind `json:"triggerKind"`
	// The trigger character (a single character) that has trigger code complete.
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
}

// CompletionTriggerKind represents how a completion was triggered.
type CompletionTriggerKind int

const (
	// CompletionTriggerKindInvoked means completion was triggered by typing an identifier.
	CompletionTriggerKindInvoked CompletionTriggerKind = 1
	// CompletionTriggerKindTriggerCharacter means completion was triggered by a trigger character.
	CompletionTriggerKindTriggerCharacter CompletionTriggerKind = 2
	// CompletionTriggerKindTriggerForIncompleteCompletions means completion was re-triggered.
	CompletionTriggerKindTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// DocumentSymbolParams represents the parameters for a textDocument/documentSymbol request.
type DocumentSymbolParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// ReferenceParams represents the parameters for a textDocument/references request.
type ReferenceParams struct {
	TextDocumentPositionParams
	// Context carrying additional information.
	Context ReferenceContext `json:"context"`
}

// ReferenceContext represents the context for a references request.
type ReferenceContext struct {
	// Include the declaration of the current symbol.
	IncludeDeclaration bool `json:"includeDeclaration"`
}
