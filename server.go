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
	documents *DocumentStore
	shutdown  bool
}

// NewServer creates a new LSP server instance
func NewServer() *Server {
	return &Server{
		documents: NewDocumentStore(),
		shutdown:  false,
	}
}

// Handle processes incoming JSON-RPC 2.0 requests and notifications
func (s *Server) Handle(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	log.Printf("Received request: %s", req.Method())

	switch req.Method() {
	case protocol.MethodInitialize:
		return s.handleInitialize(ctx, reply, req)
	case protocol.MethodInitialized:
		return s.handleInitialized(ctx, reply, req)
	case protocol.MethodShutdown:
		return s.handleShutdown(ctx, reply, req)
	case protocol.MethodExit:
		return s.handleExit(ctx, reply, req)
	case protocol.MethodTextDocumentDidOpen:
		return s.handleTextDocumentDidOpen(ctx, reply, req)
	case protocol.MethodTextDocumentDidChange:
		return s.handleTextDocumentDidChange(ctx, reply, req)
	case protocol.MethodTextDocumentDidClose:
		return s.handleTextDocumentDidClose(ctx, reply, req)
	case protocol.MethodTextDocumentHover:
		return s.handleTextDocumentHover(ctx, reply, req)
	case protocol.MethodTextDocumentCompletion:
		return s.handleTextDocumentCompletion(ctx, reply, req)
	default:
		log.Printf("Unhandled method: %s", req.Method())
		return reply(ctx, nil, jsonrpc2.ErrMethodNotFound)
	}
}

// handleInitialize processes the initialize request
func (s *Server) handleInitialize(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.InitializeParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		log.Printf("Error unmarshaling initialize params: %v", err)
		return reply(ctx, nil, err)
	}

	log.Printf("Initializing for workspace: %v", params.WorkspaceFolders)

	// Return server capabilities
	result := protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.TextDocumentSyncKindFull,
			},
			HoverProvider: &protocol.HoverOptions{
				WorkDoneProgressOptions: protocol.WorkDoneProgressOptions{},
			},
			CompletionProvider: &protocol.CompletionOptions{
				TriggerCharacters: []string{".", ":", ">"},
			},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "boilerplate-lsp-server",
			Version: "0.1.0",
		},
	}

	return reply(ctx, result, nil)
}

// handleInitialized processes the initialized notification
func (s *Server) handleInitialized(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	log.Println("Server initialized")
	return reply(ctx, nil, nil)
}

// handleShutdown processes the shutdown request
func (s *Server) handleShutdown(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	log.Println("Shutdown request received")
	s.shutdown = true
	return reply(ctx, nil, nil)
}

// handleExit processes the exit notification
func (s *Server) handleExit(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	log.Println("Exit notification received")
	if s.shutdown {
		// Clean shutdown
		s.conn.Close()
		return nil
	}
	// Abnormal exit
	log.Println("Exit without shutdown - abnormal termination")
	s.conn.Close()
	return nil
}
