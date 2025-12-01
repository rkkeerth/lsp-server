package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rkkeerth/lsp-server/protocol"
	"github.com/sourcegraph/jsonrpc2"
)

// Server represents the LSP server
type Server struct {
	documents *DocumentStore
	logger    *log.Logger
}

// NewServer creates a new LSP server
func NewServer(logger *log.Logger) *Server {
	return &Server{
		documents: NewDocumentStore(),
		logger:    logger,
	}
}

// Handle processes incoming JSON-RPC requests
func (s *Server) Handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (interface{}, error) {
	s.logger.Printf("Received request: %s", req.Method)

	switch req.Method {
	case "initialize":
		return s.handleInitialize(ctx, req)
	case "initialized":
		return s.handleInitialized(ctx, req)
	case "shutdown":
		return s.handleShutdown(ctx, req)
	case "exit":
		return s.handleExit(ctx, req)
	case "textDocument/didOpen":
		return s.handleTextDocumentDidOpen(ctx, req)
	case "textDocument/didChange":
		return s.handleTextDocumentDidChange(ctx, req)
	case "textDocument/didClose":
		return s.handleTextDocumentDidClose(ctx, req)
	case "textDocument/hover":
		return s.handleTextDocumentHover(ctx, req)
	default:
		return nil, &jsonrpc2.Error{
			Code:    jsonrpc2.CodeMethodNotFound,
			Message: fmt.Sprintf("method not found: %s", req.Method),
		}
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	var params protocol.InitializeParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshaling initialize params: %v", err)
		return nil, err
	}

	s.logger.Printf("Initializing server for root URI: %s", params.RootURI)

	result := protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.Full,
			},
			HoverProvider: true,
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "basic-lsp-server",
			Version: "0.1.0",
		},
	}

	return result, nil
}

// handleInitialized handles the initialized notification
func (s *Server) handleInitialized(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	s.logger.Println("Server initialized")
	return nil, nil
}

// handleShutdown handles the shutdown request
func (s *Server) handleShutdown(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	s.logger.Println("Server shutting down")
	return nil, nil
}

// handleExit handles the exit notification
func (s *Server) handleExit(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	s.logger.Println("Server exiting")
	return nil, nil
}

// handleTextDocumentDidOpen handles the textDocument/didOpen notification
func (s *Server) handleTextDocumentDidOpen(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshaling didOpen params: %v", err)
		return nil, err
	}

	doc := params.TextDocument
	s.documents.Open(doc.URI, doc.LanguageID, doc.Version, doc.Text)
	s.logger.Printf("Opened document: %s (language: %s, version: %d)", doc.URI, doc.LanguageID, doc.Version)

	return nil, nil
}

// handleTextDocumentDidChange handles the textDocument/didChange notification
func (s *Server) handleTextDocumentDidChange(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshaling didChange params: %v", err)
		return nil, err
	}

	// For full document sync, we only need the first change event
	if len(params.ContentChanges) > 0 {
		change := params.ContentChanges[0]
		s.documents.Update(params.TextDocument.URI, params.TextDocument.Version, change.Text)
		s.logger.Printf("Updated document: %s (version: %d)", params.TextDocument.URI, params.TextDocument.Version)
	}

	return nil, nil
}

// handleTextDocumentDidClose handles the textDocument/didClose notification
func (s *Server) handleTextDocumentDidClose(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshaling didClose params: %v", err)
		return nil, err
	}

	s.documents.Close(params.TextDocument.URI)
	s.logger.Printf("Closed document: %s", params.TextDocument.URI)

	return nil, nil
}

// handleTextDocumentHover handles the textDocument/hover request
func (s *Server) handleTextDocumentHover(ctx context.Context, req *jsonrpc2.Request) (interface{}, error) {
	var params protocol.HoverParams
	if err := json.Unmarshal(*req.Params, &params); err != nil {
		s.logger.Printf("Error unmarshaling hover params: %v", err)
		return nil, err
	}

	doc, exists := s.documents.Get(params.TextDocument.URI)
	if !exists {
		s.logger.Printf("Document not found: %s", params.TextDocument.URI)
		return nil, nil
	}

	word := doc.GetWordAtPosition(params.Position.Line, params.Position.Character)
	if word == "" {
		return nil, nil
	}

	s.logger.Printf("Hover requested at %s:%d:%d, word: %s",
		params.TextDocument.URI, params.Position.Line, params.Position.Character, word)

	// Provide basic hover information
	hoverText := fmt.Sprintf("**%s**\n\nIdentifier in %s", word, doc.LanguageID)

	hover := protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.Markdown,
			Value: hoverText,
		},
	}

	return hover, nil
}
