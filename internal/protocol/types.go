// Package protocol defines the LSP (Language Server Protocol) types and structures.
package protocol

// Position represents a position in a text document.
type Position struct {
	// Line position in a document (zero-based).
	Line uint32 `json:"line"`
	// Character offset on a line in a document (zero-based).
	Character uint32 `json:"character"`
}

// Range represents a range in a text document.
type Range struct {
	// The range's start position.
	Start Position `json:"start"`
	// The range's end position.
	End Position `json:"end"`
}

// Location represents a location inside a resource, such as a line inside a text file.
type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

// DocumentURI is the URI of a document.
type DocumentURI string

// TextDocumentIdentifier identifies a text document.
type TextDocumentIdentifier struct {
	// The text document's URI.
	URI DocumentURI `json:"uri"`
}

// VersionedTextDocumentIdentifier is an identifier to denote a specific version of a text document.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	// The version number of this document.
	Version int32 `json:"version"`
}

// TextDocumentItem is an item to transfer a text document from the client to the server.
type TextDocumentItem struct {
	// The text document's URI.
	URI DocumentURI `json:"uri"`
	// The text document's language identifier.
	LanguageID string `json:"languageId"`
	// The version number of this document.
	Version int32 `json:"version"`
	// The content of the opened text document.
	Text string `json:"text"`
}

// TextDocumentPositionParams is a parameter literal used in requests to pass a text document and a position.
type TextDocumentPositionParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The position inside the text document.
	Position Position `json:"position"`
}

// TextEdit represents a textual edit applicable to a text document.
type TextEdit struct {
	// The range of the text document to be manipulated.
	Range Range `json:"range"`
	// The string to be inserted.
	NewText string `json:"newText"`
}

// TextDocumentContentChangeEvent represents an event describing a change to a text document.
type TextDocumentContentChangeEvent struct {
	// The range of the document that changed.
	Range *Range `json:"range,omitempty"`
	// The optional length of the range that got replaced.
	RangeLength *uint32 `json:"rangeLength,omitempty"`
	// The new text for the provided range or the full document.
	Text string `json:"text"`
}

// Diagnostic represents a diagnostic, such as a compiler error or warning.
type Diagnostic struct {
	// The range at which the message applies.
	Range Range `json:"range"`
	// The diagnostic's severity.
	Severity DiagnosticSeverity `json:"severity,omitempty"`
	// The diagnostic's code, which might appear in the user interface.
	Code interface{} `json:"code,omitempty"`
	// A human-readable string describing the source of this diagnostic.
	Source string `json:"source,omitempty"`
	// The diagnostic's message.
	Message string `json:"message"`
	// Additional metadata about the diagnostic.
	Tags []DiagnosticTag `json:"tags,omitempty"`
	// An array of related diagnostic information.
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
}

// DiagnosticSeverity represents the severity of a diagnostic.
type DiagnosticSeverity int

const (
	// DiagnosticSeverityError reports an error.
	DiagnosticSeverityError DiagnosticSeverity = 1
	// DiagnosticSeverityWarning reports a warning.
	DiagnosticSeverityWarning DiagnosticSeverity = 2
	// DiagnosticSeverityInformation reports an information.
	DiagnosticSeverityInformation DiagnosticSeverity = 3
	// DiagnosticSeverityHint reports a hint.
	DiagnosticSeverityHint DiagnosticSeverity = 4
)

// DiagnosticTag represents additional metadata about the type of a diagnostic.
type DiagnosticTag int

const (
	// DiagnosticTagUnnecessary marks the code as unnecessary.
	DiagnosticTagUnnecessary DiagnosticTag = 1
	// DiagnosticTagDeprecated marks the code as deprecated.
	DiagnosticTagDeprecated DiagnosticTag = 2
)

// DiagnosticRelatedInformation represents a related message and source code location for a diagnostic.
type DiagnosticRelatedInformation struct {
	// The location of this related diagnostic information.
	Location Location `json:"location"`
	// The message of this related diagnostic information.
	Message string `json:"message"`
}

// CompletionItem represents a completion item.
type CompletionItem struct {
	// The label of this completion item.
	Label string `json:"label"`
	// The kind of this completion item.
	Kind CompletionItemKind `json:"kind,omitempty"`
	// A human-readable string with additional information about this item.
	Detail string `json:"detail,omitempty"`
	// A human-readable string that represents a doc-comment.
	Documentation interface{} `json:"documentation,omitempty"`
	// A string that should be inserted into a document when selecting this completion.
	InsertText string `json:"insertText,omitempty"`
	// The format of the insert text.
	InsertTextFormat InsertTextFormat `json:"insertTextFormat,omitempty"`
	// An edit which is applied to a document when selecting this completion.
	TextEdit *TextEdit `json:"textEdit,omitempty"`
	// An optional array of additional text edits that are applied when selecting this completion.
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`
	// A data entry field that is preserved on a completion item between a completion and a completion resolve request.
	Data interface{} `json:"data,omitempty"`
}

// CompletionItemKind represents the kind of a completion entry.
type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

// InsertTextFormat defines whether the insert text in a completion item should be interpreted as plain text or a snippet.
type InsertTextFormat int

const (
	// InsertTextFormatPlainText indicates the insert text is a plain string.
	InsertTextFormatPlainText InsertTextFormat = 1
	// InsertTextFormatSnippet indicates the insert text is a snippet.
	InsertTextFormatSnippet InsertTextFormat = 2
)

// CompletionList represents a list of completion items.
type CompletionList struct {
	// This list is not complete. Further typing should result in recomputing this list.
	IsIncomplete bool `json:"isIncomplete"`
	// The completion items.
	Items []CompletionItem `json:"items"`
}

// Hover represents the result of a hover request.
type Hover struct {
	// The hover's content.
	Contents MarkupContent `json:"contents"`
	// An optional range.
	Range *Range `json:"range,omitempty"`
}

// MarkupContent represents a string value which content is interpreted base on its kind.
type MarkupContent struct {
	// The type of the markup.
	Kind MarkupKind `json:"kind"`
	// The content itself.
	Value string `json:"value"`
}

// MarkupKind describes the content type that a client supports.
type MarkupKind string

const (
	// MarkupKindPlainText indicates plain text is supported as a content format.
	MarkupKindPlainText MarkupKind = "plaintext"
	// MarkupKindMarkdown indicates markdown is supported as a content format.
	MarkupKindMarkdown MarkupKind = "markdown"
)

// SymbolKind represents a symbol kind.
type SymbolKind int

const (
	SymbolKindFile          SymbolKind = 1
	SymbolKindModule        SymbolKind = 2
	SymbolKindNamespace     SymbolKind = 3
	SymbolKindPackage       SymbolKind = 4
	SymbolKindClass         SymbolKind = 5
	SymbolKindMethod        SymbolKind = 6
	SymbolKindProperty      SymbolKind = 7
	SymbolKindField         SymbolKind = 8
	SymbolKindConstructor   SymbolKind = 9
	SymbolKindEnum          SymbolKind = 10
	SymbolKindInterface     SymbolKind = 11
	SymbolKindFunction      SymbolKind = 12
	SymbolKindVariable      SymbolKind = 13
	SymbolKindConstant      SymbolKind = 14
	SymbolKindString        SymbolKind = 15
	SymbolKindNumber        SymbolKind = 16
	SymbolKindBoolean       SymbolKind = 17
	SymbolKindArray         SymbolKind = 18
	SymbolKindObject        SymbolKind = 19
	SymbolKindKey           SymbolKind = 20
	SymbolKindNull          SymbolKind = 21
	SymbolKindEnumMember    SymbolKind = 22
	SymbolKindStruct        SymbolKind = 23
	SymbolKindEvent         SymbolKind = 24
	SymbolKindOperator      SymbolKind = 25
	SymbolKindTypeParameter SymbolKind = 26
)

// DocumentSymbol represents programming constructs like variables, classes, interfaces, etc.
type DocumentSymbol struct {
	// The name of this symbol.
	Name string `json:"name"`
	// More detail for this symbol, e.g. the signature of a function.
	Detail string `json:"detail,omitempty"`
	// The kind of this symbol.
	Kind SymbolKind `json:"kind"`
	// Tags for this document symbol.
	Tags []SymbolTag `json:"tags,omitempty"`
	// The range enclosing this symbol not including leading/trailing whitespace.
	Range Range `json:"range"`
	// The range that should be selected and revealed when this symbol is being picked.
	SelectionRange Range `json:"selectionRange"`
	// Children of this symbol, e.g. properties of a class.
	Children []DocumentSymbol `json:"children,omitempty"`
}

// SymbolTag represents extra annotations that tweak the rendering of a symbol.
type SymbolTag int

const (
	// SymbolTagDeprecated renders a symbol as obsolete, usually using a strike-out.
	SymbolTagDeprecated SymbolTag = 1
)

// SymbolInformation represents information about programming constructs.
type SymbolInformation struct {
	// The name of this symbol.
	Name string `json:"name"`
	// The kind of this symbol.
	Kind SymbolKind `json:"kind"`
	// Tags for this symbol.
	Tags []SymbolTag `json:"tags,omitempty"`
	// The location of this symbol.
	Location Location `json:"location"`
	// The name of the symbol containing this symbol.
	ContainerName string `json:"containerName,omitempty"`
}
