// Package lsp provides the Language Server Protocol handler implementation.
package lsp

import (
	"context"
	"encoding/json"
	"log"

	"github.com/example/lsp-server/internal/handler"
	"github.com/example/lsp-server/internal/protocol"
)

// Handler processes LSP requests and notifications.
type Handler struct {
	server        *Server
	logger        *log.Logger
	textDocuments *handler.TextDocumentHandler
}

// NewHandler creates a new LSP handler.
func NewHandler(server *Server, logger *log.Logger) *Handler {
	return &Handler{
		server:        server,
		logger:        logger,
		textDocuments: handler.NewTextDocumentHandler(logger),
	}
}

// HandleRequest processes an LSP request and returns a result.
func (h *Handler) HandleRequest(ctx context.Context, method string, params json.RawMessage) (interface{}, error) {
	switch method {
	case "initialize":
		return h.handleInitialize(ctx, params)
	case "shutdown":
		return h.handleShutdown(ctx)
	case "textDocument/completion":
		return h.textDocuments.HandleCompletion(ctx, params)
	case "textDocument/hover":
		return h.textDocuments.HandleHover(ctx, params)
	case "textDocument/definition":
		return h.textDocuments.HandleDefinition(ctx, params)
	case "textDocument/references":
		return h.textDocuments.HandleReferences(ctx, params)
	case "textDocument/documentSymbol":
		return h.textDocuments.HandleDocumentSymbol(ctx, params)
	case "textDocument/formatting":
		return h.textDocuments.HandleFormatting(ctx, params)
	case "textDocument/codeAction":
		return h.textDocuments.HandleCodeAction(ctx, params)
	default:
		h.logger.Printf("Method not found: %s", method)
		return nil, NewRPCError(protocol.MethodNotFound, "Method not found: "+method)
	}
}

// HandleNotification processes an LSP notification.
func (h *Handler) HandleNotification(ctx context.Context, method string, params json.RawMessage) error {
	switch method {
	case "initialized":
		return h.handleInitialized(ctx)
	case "exit":
		return h.handleExit(ctx)
	case "textDocument/didOpen":
		return h.textDocuments.HandleDidOpen(ctx, params)
	case "textDocument/didChange":
		return h.textDocuments.HandleDidChange(ctx, params)
	case "textDocument/didSave":
		return h.textDocuments.HandleDidSave(ctx, params)
	case "textDocument/didClose":
		return h.textDocuments.HandleDidClose(ctx, params)
	case "$/cancelRequest":
		// Cancel requests are handled but we don't support cancellation yet
		return nil
	default:
		h.logger.Printf("Unknown notification: %s", method)
		return nil
	}
}

// handleInitialize handles the initialize request.
func (h *Handler) handleInitialize(ctx context.Context, params json.RawMessage) (*protocol.InitializeResult, error) {
	var initParams protocol.InitializeParams
	if err := json.Unmarshal(params, &initParams); err != nil {
		return nil, NewRPCError(protocol.InvalidParams, "Invalid initialize params")
	}

	h.logger.Printf("Initializing server for workspace: %s", initParams.RootURI)

	// Return server capabilities
	result := &protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindIncremental,
				Save: &protocol.SaveOptions{
					IncludeText: true,
				},
			},
			CompletionProvider: &protocol.CompletionOptions{
				ResolveProvider:   true,
				TriggerCharacters: []string{".", ":", "<", "/"},
			},
			HoverProvider:             true,
			DefinitionProvider:        true,
			ReferencesProvider:        true,
			DocumentSymbolProvider:    true,
			DocumentFormattingProvider: true,
			CodeActionProvider:        true,
			Workspace: &protocol.ServerWorkspaceCapabilities{
				WorkspaceFolders: &protocol.WorkspaceFoldersServerCapabilities{
					Supported:           true,
					ChangeNotifications: true,
				},
			},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "lsp-server-go",
			Version: "0.1.0",
		},
	}

	return result, nil
}

// handleInitialized handles the initialized notification.
func (h *Handler) handleInitialized(ctx context.Context) error {
	h.logger.Println("Server initialized")
	h.server.SetInitialized(true)
	return nil
}

// handleShutdown handles the shutdown request.
func (h *Handler) handleShutdown(ctx context.Context) (interface{}, error) {
	h.logger.Println("Shutdown request received")
	return nil, nil
}

// handleExit handles the exit notification.
func (h *Handler) handleExit(ctx context.Context) error {
	h.logger.Println("Exit notification received")
	h.server.Shutdown()
	return nil
}
