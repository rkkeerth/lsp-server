package server

import (
	"github.com/rkkeerth/lsp-server/lsp"
)

// HandleInitialize processes the initialize request and returns server capabilities
func HandleInitialize(params *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    lsp.TextDocumentSyncKindFull,
				Save: &lsp.SaveOptions{
					IncludeText: true,
				},
			},
		},
		ServerInfo: &lsp.ServerInfo{
			Name:    "go-lsp-server",
			Version: "0.1.0",
		},
	}, nil
}

// HandleInitialized processes the initialized notification
func HandleInitialized(params *lsp.InitializedParams) {
	// No-op - server is ready
}

// HandleShutdown processes the shutdown request
func HandleShutdown() *lsp.ShutdownResult {
	return &lsp.ShutdownResult{}
}

// HandleTextDocumentDidOpen processes textDocument/didOpen notifications
func HandleTextDocumentDidOpen(state *State, params *lsp.DidOpenTextDocumentParams) {
	state.OpenDocument(params.TextDocument.URI, params.TextDocument.Text)
}

// HandleTextDocumentDidChange processes textDocument/didChange notifications
func HandleTextDocumentDidChange(state *State, params *lsp.DidChangeTextDocumentParams) {
	// Full sync: use the last content change (should be the complete document)
	if len(params.ContentChanges) > 0 {
		lastChange := params.ContentChanges[len(params.ContentChanges)-1]
		state.UpdateDocument(params.TextDocument.URI, lastChange.Text)
	}
}

// HandleTextDocumentDidClose processes textDocument/didClose notifications
func HandleTextDocumentDidClose(state *State, params *lsp.DidCloseTextDocumentParams) {
	state.CloseDocument(params.TextDocument.URI)
}
