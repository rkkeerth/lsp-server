package server

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/rkkeerth/lsp-server/internal/protocol"
)

// Server represents an LSP server
type Server struct {
	initialized bool
	shutdown    bool
	documents   map[string]*Document
	mu          sync.RWMutex
}

// Document represents an opened text document
type Document struct {
	URI        string
	LanguageID string
	Version    int
	Text       string
}

// NewServer creates a new LSP server
func NewServer() *Server {
	return &Server{
		documents: make(map[string]*Document),
	}
}

// Handle handles incoming LSP requests and notifications
func (s *Server) Handle(method string, params json.RawMessage) (interface{}, error) {
	switch method {
	case "initialize":
		return s.handleInitialize(params)
	case "initialized":
		return s.handleInitialized(params)
	case "shutdown":
		return s.handleShutdown(params)
	case "exit":
		return s.handleExit(params)
	case "textDocument/didOpen":
		return s.handleDidOpen(params)
	case "textDocument/didChange":
		return s.handleDidChange(params)
	case "textDocument/didClose":
		return s.handleDidClose(params)
	default:
		return nil, fmt.Errorf("method not found: %s", method)
	}
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(params json.RawMessage) (interface{}, error) {
	var initParams protocol.InitializeParams
	if err := json.Unmarshal(params, &initParams); err != nil {
		return nil, fmt.Errorf("invalid initialize params: %v", err)
	}

	// Create server capabilities
	capabilities := protocol.ServerCapabilities{
		TextDocumentSync: &protocol.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    protocol.Full,
			Save: &protocol.SaveOptions{
				IncludeText: true,
			},
		},
		HoverProvider:      false,
		DefinitionProvider: false,
		ReferencesProvider: false,
	}

	result := protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.ServerInfo{
			Name:    "basic-lsp-server",
			Version: "0.1.0",
		},
	}

	return result, nil
}

// handleInitialized handles the initialized notification
func (s *Server) handleInitialized(params json.RawMessage) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.initialized = true
	return nil, nil
}

// handleShutdown handles the shutdown request
func (s *Server) handleShutdown(params json.RawMessage) (interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.shutdown = true
	return nil, nil
}

// handleExit handles the exit notification
func (s *Server) handleExit(params json.RawMessage) (interface{}, error) {
	// Exit is handled by the main function
	return nil, nil
}

// handleDidOpen handles the textDocument/didOpen notification
func (s *Server) handleDidOpen(params json.RawMessage) (interface{}, error) {
	var openParams protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(params, &openParams); err != nil {
		return nil, fmt.Errorf("invalid didOpen params: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	doc := &Document{
		URI:        openParams.TextDocument.URI,
		LanguageID: openParams.TextDocument.LanguageID,
		Version:    openParams.TextDocument.Version,
		Text:       openParams.TextDocument.Text,
	}

	s.documents[doc.URI] = doc

	return nil, nil
}

// handleDidChange handles the textDocument/didChange notification
func (s *Server) handleDidChange(params json.RawMessage) (interface{}, error) {
	var changeParams protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(params, &changeParams); err != nil {
		return nil, fmt.Errorf("invalid didChange params: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	doc, exists := s.documents[changeParams.TextDocument.URI]
	if !exists {
		return nil, fmt.Errorf("document not found: %s", changeParams.TextDocument.URI)
	}

	// Update document version
	doc.Version = changeParams.TextDocument.Version

	// Apply content changes (for full sync, we just replace the entire text)
	if len(changeParams.ContentChanges) > 0 {
		// For Full sync mode, the last change contains the full text
		doc.Text = changeParams.ContentChanges[len(changeParams.ContentChanges)-1].Text
	}

	return nil, nil
}

// handleDidClose handles the textDocument/didClose notification
func (s *Server) handleDidClose(params json.RawMessage) (interface{}, error) {
	var closeParams protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(params, &closeParams); err != nil {
		return nil, fmt.Errorf("invalid didClose params: %v", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.documents, closeParams.TextDocument.URI)

	return nil, nil
}

// GetDocument returns a document by URI
func (s *Server) GetDocument(uri string) (*Document, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doc, exists := s.documents[uri]
	return doc, exists
}

// IsInitialized returns whether the server is initialized
func (s *Server) IsInitialized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.initialized
}

// IsShutdown returns whether the server is shutdown
func (s *Server) IsShutdown() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.shutdown
}
