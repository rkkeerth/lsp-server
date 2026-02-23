package handler

import (
	"github.com/rkkeerth/lsp-server/internal/document"
	"github.com/rkkeerth/lsp-server/internal/protocol"
)

// Hover handles hover requests
type Hover struct {
	store *document.Store
}

// NewHover creates a new hover handler
func NewHover(store *document.Store) *Hover {
	return &Hover{store: store}
}

// Handle processes textDocument/hover
func (h *Hover) Handle(params *protocol.HoverParams) *protocol.Hover {
	// Example placeholder hover
	// In a real implementation, analyze the document content and position
	// to provide context-aware hover information

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: "**Example Hover**\n\nThis is placeholder hover content.",
		},
	}
}
