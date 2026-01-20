// Package handler provides handlers for LSP text document operations.
package handler

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/example/lsp-server/internal/protocol"
)

// TextDocumentHandler handles text document related LSP operations.
type TextDocumentHandler struct {
	logger    *log.Logger
	documents map[protocol.DocumentURI]*Document
	mu        sync.RWMutex
}

// Document represents an open text document.
type Document struct {
	URI        protocol.DocumentURI
	LanguageID string
	Version    int32
	Content    string
}

// NewTextDocumentHandler creates a new text document handler.
func NewTextDocumentHandler(logger *log.Logger) *TextDocumentHandler {
	return &TextDocumentHandler{
		logger:    logger,
		documents: make(map[protocol.DocumentURI]*Document),
	}
}

// HandleDidOpen handles the textDocument/didOpen notification.
func (h *TextDocumentHandler) HandleDidOpen(ctx context.Context, params json.RawMessage) error {
	var p protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	h.documents[p.TextDocument.URI] = &Document{
		URI:        p.TextDocument.URI,
		LanguageID: p.TextDocument.LanguageID,
		Version:    p.TextDocument.Version,
		Content:    p.TextDocument.Text,
	}

	h.logger.Printf("Document opened: %s", p.TextDocument.URI)
	return nil
}

// HandleDidChange handles the textDocument/didChange notification.
func (h *TextDocumentHandler) HandleDidChange(ctx context.Context, params json.RawMessage) error {
	var p protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	doc, ok := h.documents[p.TextDocument.URI]
	if !ok {
		h.logger.Printf("Document not found: %s", p.TextDocument.URI)
		return nil
	}

	// Apply changes
	for _, change := range p.ContentChanges {
		if change.Range == nil {
			// Full document update
			doc.Content = change.Text
		} else {
			// Incremental update
			doc.Content = applyChange(doc.Content, change)
		}
	}
	doc.Version = p.TextDocument.Version

	h.logger.Printf("Document changed: %s (version %d)", p.TextDocument.URI, doc.Version)
	return nil
}

// HandleDidSave handles the textDocument/didSave notification.
func (h *TextDocumentHandler) HandleDidSave(ctx context.Context, params json.RawMessage) error {
	var p protocol.DidSaveTextDocumentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	h.logger.Printf("Document saved: %s", p.TextDocument.URI)
	return nil
}

// HandleDidClose handles the textDocument/didClose notification.
func (h *TextDocumentHandler) HandleDidClose(ctx context.Context, params json.RawMessage) error {
	var p protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.documents, p.TextDocument.URI)
	h.logger.Printf("Document closed: %s", p.TextDocument.URI)
	return nil
}

// HandleCompletion handles the textDocument/completion request.
func (h *TextDocumentHandler) HandleCompletion(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p protocol.CompletionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	h.logger.Printf("Completion requested at %s:%d:%d", p.TextDocument.URI, p.Position.Line, p.Position.Character)

	// Return example completion items
	// TODO: Implement actual completion logic based on your language
	items := []protocol.CompletionItem{
		{
			Label:      "example",
			Kind:       protocol.CompletionItemKindKeyword,
			Detail:     "An example completion item",
			InsertText: "example",
		},
		{
			Label:      "function",
			Kind:       protocol.CompletionItemKindFunction,
			Detail:     "A function completion",
			InsertText: "function() {\n\t$0\n}",
			InsertTextFormat: protocol.InsertTextFormatSnippet,
		},
		{
			Label:      "variable",
			Kind:       protocol.CompletionItemKindVariable,
			Detail:     "A variable completion",
			InsertText: "variable",
		},
	}

	return &protocol.CompletionList{
		IsIncomplete: false,
		Items:        items,
	}, nil
}

// HandleHover handles the textDocument/hover request.
func (h *TextDocumentHandler) HandleHover(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p protocol.TextDocumentPositionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	h.logger.Printf("Hover requested at %s:%d:%d", p.TextDocument.URI, p.Position.Line, p.Position.Character)

	// TODO: Implement actual hover logic based on your language
	// This is a placeholder that returns example hover information
	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "**Example Hover**\n\nThis is a placeholder hover message. Implement your language-specific hover logic here.",
		},
	}, nil
}

// HandleDefinition handles the textDocument/definition request.
func (h *TextDocumentHandler) HandleDefinition(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p protocol.TextDocumentPositionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	h.logger.Printf("Definition requested at %s:%d:%d", p.TextDocument.URI, p.Position.Line, p.Position.Character)

	// TODO: Implement actual go-to-definition logic
	// Return empty array when no definition is found
	return []protocol.Location{}, nil
}

// HandleReferences handles the textDocument/references request.
func (h *TextDocumentHandler) HandleReferences(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p protocol.ReferenceParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	h.logger.Printf("References requested at %s:%d:%d", p.TextDocument.URI, p.Position.Line, p.Position.Character)

	// TODO: Implement actual find-references logic
	// Return empty array when no references are found
	return []protocol.Location{}, nil
}

// HandleDocumentSymbol handles the textDocument/documentSymbol request.
func (h *TextDocumentHandler) HandleDocumentSymbol(ctx context.Context, params json.RawMessage) (interface{}, error) {
	var p protocol.DocumentSymbolParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}

	h.logger.Printf("Document symbols requested for %s", p.TextDocument.URI)

	// TODO: Implement actual document symbol extraction
	// Return empty array when no symbols are found
	return []protocol.DocumentSymbol{}, nil
}

// HandleFormatting handles the textDocument/formatting request.
func (h *TextDocumentHandler) HandleFormatting(ctx context.Context, params json.RawMessage) (interface{}, error) {
	h.logger.Println("Formatting requested")

	// TODO: Implement actual formatting logic
	// Return empty array when no edits are needed
	return []protocol.TextEdit{}, nil
}

// HandleCodeAction handles the textDocument/codeAction request.
func (h *TextDocumentHandler) HandleCodeAction(ctx context.Context, params json.RawMessage) (interface{}, error) {
	h.logger.Println("Code action requested")

	// TODO: Implement actual code action logic
	// Return empty array when no code actions are available
	return []interface{}{}, nil
}

// GetDocument retrieves a document by URI.
func (h *TextDocumentHandler) GetDocument(uri protocol.DocumentURI) (*Document, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	doc, ok := h.documents[uri]
	return doc, ok
}

// applyChange applies an incremental change to the document content.
func applyChange(content string, change protocol.TextDocumentContentChangeEvent) string {
	if change.Range == nil {
		return change.Text
	}

	lines := splitLines(content)
	
	startLine := int(change.Range.Start.Line)
	startChar := int(change.Range.Start.Character)
	endLine := int(change.Range.End.Line)
	endChar := int(change.Range.End.Character)

	// Ensure we have enough lines
	for len(lines) <= endLine {
		lines = append(lines, "")
	}

	// Build the new content
	var result string

	// Add content before the change
	for i := 0; i < startLine; i++ {
		result += lines[i] + "\n"
	}

	// Add start of the line before the change
	if startLine < len(lines) && startChar <= len(lines[startLine]) {
		result += lines[startLine][:startChar]
	}

	// Add the new text
	result += change.Text

	// Add the rest of the end line after the change
	if endLine < len(lines) && endChar <= len(lines[endLine]) {
		result += lines[endLine][endChar:]
	}

	// Add remaining lines
	for i := endLine + 1; i < len(lines); i++ {
		result += "\n" + lines[i]
	}

	return result
}

// splitLines splits content into lines while preserving empty lines.
func splitLines(content string) []string {
	var lines []string
	var current string
	
	for _, r := range content {
		if r == '\n' {
			lines = append(lines, current)
			current = ""
		} else if r != '\r' {
			current += string(r)
		}
	}
	
	// Add the last line if there's remaining content
	if current != "" || len(content) == 0 || (len(content) > 0 && content[len(content)-1] == '\n') {
		lines = append(lines, current)
	}
	
	return lines
}
