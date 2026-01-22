package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
)

// Server represents the LSP server
type Server struct {
	conn      jsonrpc2.Conn
	documents *DocumentManager
	
	// Server capabilities
	initialized bool
	clientInfo  *protocol.ClientInfo
}

// NewServer creates a new LSP server instance
func NewServer() *Server {
	return &Server{
		documents: NewDocumentManager(),
	}
}

// Handle processes incoming JSON-RPC requests
func (s *Server) Handle(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	log.Printf("Received request: method=%s", req.Method())

	switch req.Method() {
	case "initialize":
		return s.handleInitialize(ctx, reply, req)
	case "initialized":
		return s.handleInitialized(ctx, reply, req)
	case "shutdown":
		return s.handleShutdown(ctx, reply, req)
	case "exit":
		return s.handleExit(ctx, reply, req)
	case "textDocument/didOpen":
		return s.handleTextDocumentDidOpen(ctx, reply, req)
	case "textDocument/didChange":
		return s.handleTextDocumentDidChange(ctx, reply, req)
	case "textDocument/didSave":
		return s.handleTextDocumentDidSave(ctx, reply, req)
	case "textDocument/didClose":
		return s.handleTextDocumentDidClose(ctx, reply, req)
	case "textDocument/completion":
		return s.handleTextDocumentCompletion(ctx, reply, req)
	case "textDocument/hover":
		return s.handleTextDocumentHover(ctx, reply, req)
	case "textDocument/definition":
		return s.handleTextDocumentDefinition(ctx, reply, req)
	default:
		log.Printf("Unhandled method: %s", req.Method())
		return reply(ctx, nil, jsonrpc2.ErrMethodNotFound)
	}
}

// handleInitialize handles the initialize request from the client
func (s *Server) handleInitialize(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.InitializeParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling initialize params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	s.clientInfo = params.ClientInfo
	log.Printf("Client: %s %s", s.clientInfo.Name, s.clientInfo.Version)

	// Define server capabilities
	result := protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
				Save: &protocol.SaveOptions{
					IncludeText: true,
				},
			},
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{".", ":", "("},
				ResolveProvider:   false,
			},
			HoverProvider: &protocol.HoverOptions{},
			DefinitionProvider: &protocol.DefinitionOptions{},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "example-lsp-server",
			Version: "0.1.0",
		},
	}

	return reply(ctx, result, nil)
}

// handleInitialized handles the initialized notification
func (s *Server) handleInitialized(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	s.initialized = true
	log.Println("Server initialized")
	return reply(ctx, nil, nil)
}

// handleShutdown handles the shutdown request
func (s *Server) handleShutdown(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	log.Println("Shutdown requested")
	s.initialized = false
	return reply(ctx, nil, nil)
}

// handleExit handles the exit notification
func (s *Server) handleExit(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	log.Println("Exit requested")
	return reply(ctx, nil, nil)
}

// handleTextDocumentDidOpen handles document open notifications
func (s *Server) handleTextDocumentDidOpen(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling didOpen params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	doc := params.TextDocument
	s.documents.Open(string(doc.URI), doc.Text, doc.Version)
	log.Printf("Opened document: %s (version %d)", doc.URI, doc.Version)

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidChange handles document change notifications
func (s *Server) handleTextDocumentDidChange(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling didChange params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	uri := string(params.TextDocument.URI)
	version := params.TextDocument.Version

	// For full text sync, we only have one change with the full text
	if len(params.ContentChanges) > 0 {
		if change, ok := params.ContentChanges[0].(protocol.TextDocumentContentChangeEvent); ok {
			s.documents.Update(uri, change.Text, version)
			log.Printf("Updated document: %s (version %d)", uri, version)
		}
	}

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidSave handles document save notifications
func (s *Server) handleTextDocumentDidSave(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidSaveTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling didSave params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	uri := string(params.TextDocument.URI)
	log.Printf("Saved document: %s", uri)

	// Optionally, perform additional actions on save
	// e.g., linting, formatting, etc.

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidClose handles document close notifications
func (s *Server) handleTextDocumentDidClose(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling didClose params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	uri := string(params.TextDocument.URI)
	s.documents.Close(uri)
	log.Printf("Closed document: %s", uri)

	return reply(ctx, nil, nil)
}

// handleTextDocumentCompletion handles completion requests
func (s *Server) handleTextDocumentCompletion(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.CompletionParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling completion params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	uri := string(params.TextDocument.URI)
	log.Printf("Completion requested for %s at %d:%d", uri, params.Position.Line, params.Position.Character)

	// Get document content
	doc, exists := s.documents.Get(uri)
	if !exists {
		return reply(ctx, nil, fmt.Errorf("document not found: %s", uri))
	}

	// Generate basic completion items
	completions := s.generateCompletions(doc, params.Position)

	result := protocol.CompletionList{
		IsIncomplete: false,
		Items:        completions,
	}

	return reply(ctx, result, nil)
}

// handleTextDocumentHover handles hover requests
func (s *Server) handleTextDocumentHover(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.HoverParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling hover params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	uri := string(params.TextDocument.URI)
	log.Printf("Hover requested for %s at %d:%d", uri, params.Position.Line, params.Position.Character)

	// Get document content
	doc, exists := s.documents.Get(uri)
	if !exists {
		return reply(ctx, nil, nil)
	}

	// Generate hover information
	hoverInfo := s.generateHover(doc, params.Position)

	return reply(ctx, hoverInfo, nil)
}

// handleTextDocumentDefinition handles go-to-definition requests
func (s *Server) handleTextDocumentDefinition(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DefinitionParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling definition params: %v", err)
		return reply(ctx, nil, jsonrpc2.NewError(jsonrpc2.InvalidParams, err.Error()))
	}

	uri := string(params.TextDocument.URI)
	log.Printf("Definition requested for %s at %d:%d", uri, params.Position.Line, params.Position.Character)

	// Get document content
	doc, exists := s.documents.Get(uri)
	if !exists {
		return reply(ctx, nil, nil)
	}

	// Find definition location
	location := s.findDefinition(doc, params.Position, uri)

	return reply(ctx, location, nil)
}
