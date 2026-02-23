package handler

import (
	"github.com/rkkeerth/lsp-server/internal/protocol"
)

// Lifecycle handles server lifecycle methods
type Lifecycle struct {
	shutdownRequested bool
}

// NewLifecycle creates a new lifecycle handler
func NewLifecycle() *Lifecycle {
	return &Lifecycle{}
}

// HandleInitialize processes the initialize request
func (l *Lifecycle) HandleInitialize(params *protocol.InitializeParams) *protocol.InitializeResult {
	return &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
			},
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{".", ":"},
			},
			HoverProvider: true,
		},
	}
}

// HandleInitialized processes the initialized notification
func (l *Lifecycle) HandleInitialized() {
	// Server is fully initialized, can start background tasks here
}

// HandleShutdown processes the shutdown request
func (l *Lifecycle) HandleShutdown() {
	l.shutdownRequested = true
}

// ShutdownRequested returns whether shutdown was requested
func (l *Lifecycle) ShutdownRequested() bool {
	return l.shutdownRequested
}
