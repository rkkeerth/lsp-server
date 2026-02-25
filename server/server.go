// Package server provides the LSP server implementation.
package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"lsp-server/document"
	"lsp-server/jsonrpc"
	"lsp-server/protocol"
)

// Server is an LSP server.
type Server struct {
	transport   *jsonrpc.Transport
	documents   *document.Manager
	initialized bool
	shutdown    bool
	logger      *log.Logger
}

// New creates a new LSP server.
func New(reader io.Reader, writer io.Writer) *Server {
	return &Server{
		transport: jsonrpc.NewTransport(reader, writer),
		documents: document.NewManager(),
		logger:    log.New(os.Stderr, "[lsp] ", log.LstdFlags),
	}
}

// Run starts the server main loop.
func (s *Server) Run() error {
	for {
		req, err := s.transport.Read()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			s.logger.Printf("Error reading request: %v", err)
			return err
		}

		resp := s.Handle(req)
		if resp != nil {
			if err := s.transport.Write(resp); err != nil {
				s.logger.Printf("Error writing response: %v", err)
				return err
			}
		}
	}
}

// Handle dispatches a JSON-RPC request to the appropriate handler.
func (s *Server) Handle(req *jsonrpc.Request) *jsonrpc.Response {
	s.logger.Printf("Received: %s", req.Method)

	switch req.Method {
	case protocol.MethodInitialize:
		return s.handleInitialize(req)
	case protocol.MethodInitialized:
		s.handleInitialized(req)
		return nil // notification, no response
	case protocol.MethodShutdown:
		return s.handleShutdown(req)
	case protocol.MethodExit:
		s.handleExit(req)
		return nil // notification, no response (and we exit)
	case protocol.MethodDidOpen:
		s.handleDidOpen(req)
		return nil // notification, no response
	case protocol.MethodDidChange:
		s.handleDidChange(req)
		return nil // notification, no response
	case protocol.MethodDidClose:
		s.handleDidClose(req)
		return nil // notification, no response
	default:
		if req.IsNotification() {
			// Unknown notification, ignore
			s.logger.Printf("Unknown notification: %s", req.Method)
			return nil
		}
		// Unknown request, return method not found
		return jsonrpc.NewErrorResponse(req.ID, jsonrpc.MethodNotFound, fmt.Sprintf("method not found: %s", req.Method))
	}
}

func (s *Server) handleInitialize(req *jsonrpc.Request) *jsonrpc.Response {
	var params protocol.InitializeParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return jsonrpc.NewErrorResponse(req.ID, jsonrpc.InvalidParams, fmt.Sprintf("invalid params: %v", err))
	}

	result := protocol.InitializeResult{
		Capabilities: protocol.ServerCapabilities{
			TextDocumentSync: &protocol.TextDocumentSyncOptions{
				OpenClose: true,
				Change:    protocol.SyncFull,
				Save: &protocol.SaveOptions{
					IncludeText: false,
				},
			},
		},
		ServerInfo: &protocol.ServerInfo{
			Name:    "lsp-server",
			Version: "0.1.0",
		},
	}

	return jsonrpc.NewResponse(req.ID, result)
}

func (s *Server) handleInitialized(req *jsonrpc.Request) {
	s.initialized = true
	s.logger.Printf("Server initialized")
}

func (s *Server) handleShutdown(req *jsonrpc.Request) *jsonrpc.Response {
	s.shutdown = true
	s.logger.Printf("Shutdown requested")
	return jsonrpc.NewResponse(req.ID, nil)
}

func (s *Server) handleExit(req *jsonrpc.Request) {
	if s.shutdown {
		s.logger.Printf("Exiting with code 0")
		os.Exit(0)
	} else {
		s.logger.Printf("Exiting with code 1 (shutdown not called)")
		os.Exit(1)
	}
}

func (s *Server) handleDidOpen(req *jsonrpc.Request) {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.logger.Printf("Error parsing didOpen params: %v", err)
		return
	}

	s.documents.Open(
		params.TextDocument.URI,
		params.TextDocument.LanguageID,
		params.TextDocument.Version,
		params.TextDocument.Text,
	)
	s.logger.Printf("Opened document: %s", params.TextDocument.URI)
}

func (s *Server) handleDidChange(req *jsonrpc.Request) {
	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.logger.Printf("Error parsing didChange params: %v", err)
		return
	}

	// For full sync mode, we expect a single content change with the full text
	if len(params.ContentChanges) > 0 {
		// Use the last content change (full sync should have one change)
		content := params.ContentChanges[len(params.ContentChanges)-1].Text
		s.documents.Change(
			params.TextDocument.URI,
			params.TextDocument.Version,
			content,
		)
		s.logger.Printf("Changed document: %s (version %d)", params.TextDocument.URI, params.TextDocument.Version)
	}
}

func (s *Server) handleDidClose(req *jsonrpc.Request) {
	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.logger.Printf("Error parsing didClose params: %v", err)
		return
	}

	s.documents.Close(params.TextDocument.URI)
	s.logger.Printf("Closed document: %s", params.TextDocument.URI)
}
