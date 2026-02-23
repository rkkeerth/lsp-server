package handler

import (
	"github.com/rkkeerth/lsp-server/internal/document"
	"github.com/rkkeerth/lsp-server/internal/protocol"
)

// TextDocument handles text document notifications
type TextDocument struct {
	store *document.Store
}

// NewTextDocument creates a new text document handler
func NewTextDocument(store *document.Store) *TextDocument {
	return &TextDocument{store: store}
}

// HandleDidOpen processes textDocument/didOpen
func (h *TextDocument) HandleDidOpen(params *protocol.DidOpenTextDocumentParams) {
	h.store.Open(
		params.TextDocument.URI,
		params.TextDocument.LanguageID,
		params.TextDocument.Version,
		params.TextDocument.Text,
	)
}

// HandleDidClose processes textDocument/didClose
func (h *TextDocument) HandleDidClose(params *protocol.DidCloseTextDocumentParams) {
	h.store.Close(params.TextDocument.URI)
}

// HandleDidChange processes textDocument/didChange
func (h *TextDocument) HandleDidChange(params *protocol.DidChangeTextDocumentParams) {
	if len(params.ContentChanges) > 0 {
		// Using full sync, so take the last content change
		lastChange := params.ContentChanges[len(params.ContentChanges)-1]
		h.store.Update(
			params.TextDocument.URI,
			params.TextDocument.Version,
			lastChange.Text,
		)
	}
}
