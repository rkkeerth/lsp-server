package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/rkkeerth/lsp-server/document"
	"github.com/rkkeerth/lsp-server/handlers"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
)

// Server represents the LSP server
type Server struct {
	logger          *zap.Logger
	conn            jsonrpc2.Conn
	documentManager *document.Manager
	handlers        *handlers.Handler
	capabilities    protocol.ServerCapabilities
	initialized     bool
	shutdown        bool
	mu              sync.RWMutex
}

// NewServer creates a new LSP server instance
func NewServer(logger *zap.Logger) *Server {
	docManager := document.NewManager()
	
	s := &Server{
		logger:          logger,
		documentManager: docManager,
		initialized:     false,
		shutdown:        false,
	}
	
	s.handlers = handlers.NewHandler(logger, docManager)
	s.setupCapabilities()
	
	return s
}

// setupCapabilities configures the server's LSP capabilities
func (s *Server) setupCapabilities() {
	s.capabilities = protocol.ServerCapabilities{
		TextDocumentSync: protocol.TextDocumentSyncOptions{
			OpenClose: true,
			Change:    protocol.TextDocumentSyncKindFull,
			Save: &protocol.SaveOptions{
				IncludeText: false,
			},
		},
		HoverProvider: true,
		DefinitionProvider: true,
		ReferencesProvider: true,
		DocumentSymbolProvider: true,
		WorkspaceSymbolProvider: true,
		CompletionProvider: &protocol.CompletionOptions{
			TriggerCharacters: []string{".", ":", ">"},
			ResolveProvider:   false,
		},
	}
}

// Run starts the LSP server with stdin/stdout communication
func (s *Server) Run(ctx context.Context, in io.ReadCloser, out io.WriteCloser) error {
	stream := jsonrpc2.NewStream(jsonrpc2.NewHeaderStream(in, out))
	conn := jsonrpc2.NewConn(stream)
	
	s.conn = conn
	
	conn.Go(ctx, s.handle)
	
	<-conn.Done()
	
	if err := conn.Err(); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	
	return nil
}

// handle routes incoming JSON-RPC requests to appropriate handlers
func (s *Server) handle(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	s.logger.Debug("Received request",
		zap.String("method", req.Method()),
		zap.ByteString("params", req.Params()),
	)

	// Check if server is shutdown
	s.mu.RLock()
	isShutdown := s.shutdown
	s.mu.RUnlock()
	
	if isShutdown && req.Method() != "exit" {
		return reply(ctx, nil, jsonrpc2.ErrServerNotInitialized)
	}

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
	case "textDocument/didClose":
		return s.handleTextDocumentDidClose(ctx, reply, req)
	case "textDocument/didSave":
		return s.handleTextDocumentDidSave(ctx, reply, req)
	case "textDocument/hover":
		return s.handleTextDocumentHover(ctx, reply, req)
	case "textDocument/definition":
		return s.handleTextDocumentDefinition(ctx, reply, req)
	case "textDocument/references":
		return s.handleTextDocumentReferences(ctx, reply, req)
	case "textDocument/documentSymbol":
		return s.handleTextDocumentSymbol(ctx, reply, req)
	case "workspace/symbol":
		return s.handleWorkspaceSymbol(ctx, reply, req)
	case "textDocument/completion":
		return s.handleTextDocumentCompletion(ctx, reply, req)
	default:
		s.logger.Warn("Unhandled method", zap.String("method", req.Method()))
		return reply(ctx, nil, jsonrpc2.ErrMethodNotFound)
	}
}

// handleInitialize processes the initialize request
func (s *Server) handleInitialize(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.InitializeParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal initialize params: %w", err))
	}

	s.logger.Info("Initialize request received",
		zap.String("rootURI", string(params.RootURI)),
		zap.String("clientName", params.ClientInfo.Name),
	)

	result := protocol.InitializeResult{
		Capabilities: s.capabilities,
		ServerInfo: &protocol.ServerInfo{
			Name:    "lsp-server",
			Version: "1.0.0",
		},
	}

	return reply(ctx, result, nil)
}

// handleInitialized processes the initialized notification
func (s *Server) handleInitialized(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	s.mu.Lock()
	s.initialized = true
	s.mu.Unlock()

	s.logger.Info("Server initialized")
	return reply(ctx, nil, nil)
}

// handleShutdown processes the shutdown request
func (s *Server) handleShutdown(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	s.mu.Lock()
	s.shutdown = true
	s.mu.Unlock()

	s.logger.Info("Server shutdown requested")
	return reply(ctx, nil, nil)
}

// handleExit processes the exit notification
func (s *Server) handleExit(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	s.logger.Info("Server exiting")
	return reply(ctx, nil, nil)
}

// handleTextDocumentDidOpen processes document open notifications
func (s *Server) handleTextDocumentDidOpen(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidOpenTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal didOpen params: %w", err))
	}

	s.documentManager.Open(params.TextDocument.URI, params.TextDocument.Text, params.TextDocument.Version)
	s.logger.Info("Document opened", zap.String("uri", string(params.TextDocument.URI)))

	// Send diagnostics for the opened document
	diagnostics := s.handlers.GetDiagnostics(params.TextDocument.URI)
	s.publishDiagnostics(ctx, params.TextDocument.URI, diagnostics)

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidChange processes document change notifications
func (s *Server) handleTextDocumentDidChange(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidChangeTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal didChange params: %w", err))
	}

	if len(params.ContentChanges) > 0 {
		// For full sync, we only care about the last change
		lastChange := params.ContentChanges[len(params.ContentChanges)-1]
		if textChange, ok := lastChange.(protocol.TextDocumentContentChangeEvent); ok {
			s.documentManager.Update(params.TextDocument.URI, textChange.Text, params.TextDocument.Version)
			s.logger.Debug("Document updated", zap.String("uri", string(params.TextDocument.URI)))

			// Send updated diagnostics
			diagnostics := s.handlers.GetDiagnostics(params.TextDocument.URI)
			s.publishDiagnostics(ctx, params.TextDocument.URI, diagnostics)
		}
	}

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidClose processes document close notifications
func (s *Server) handleTextDocumentDidClose(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidCloseTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal didClose params: %w", err))
	}

	s.documentManager.Close(params.TextDocument.URI)
	s.logger.Info("Document closed", zap.String("uri", string(params.TextDocument.URI)))

	return reply(ctx, nil, nil)
}

// handleTextDocumentDidSave processes document save notifications
func (s *Server) handleTextDocumentDidSave(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DidSaveTextDocumentParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal didSave params: %w", err))
	}

	s.logger.Info("Document saved", zap.String("uri", string(params.TextDocument.URI)))

	return reply(ctx, nil, nil)
}

// handleTextDocumentHover processes hover requests
func (s *Server) handleTextDocumentHover(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.HoverParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal hover params: %w", err))
	}

	hover := s.handlers.Hover(params.TextDocument.URI, params.Position)
	return reply(ctx, hover, nil)
}

// handleTextDocumentDefinition processes go-to-definition requests
func (s *Server) handleTextDocumentDefinition(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DefinitionParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal definition params: %w", err))
	}

	locations := s.handlers.Definition(params.TextDocument.URI, params.Position)
	return reply(ctx, locations, nil)
}

// handleTextDocumentReferences processes find-references requests
func (s *Server) handleTextDocumentReferences(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.ReferenceParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal references params: %w", err))
	}

	locations := s.handlers.References(params.TextDocument.URI, params.Position)
	return reply(ctx, locations, nil)
}

// handleTextDocumentSymbol processes document symbol requests
func (s *Server) handleTextDocumentSymbol(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.DocumentSymbolParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal document symbol params: %w", err))
	}

	symbols := s.handlers.DocumentSymbols(params.TextDocument.URI)
	return reply(ctx, symbols, nil)
}

// handleWorkspaceSymbol processes workspace symbol requests
func (s *Server) handleWorkspaceSymbol(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.WorkspaceSymbolParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal workspace symbol params: %w", err))
	}

	symbols := s.handlers.WorkspaceSymbols(params.Query)
	return reply(ctx, symbols, nil)
}

// handleTextDocumentCompletion processes code completion requests
func (s *Server) handleTextDocumentCompletion(ctx context.Context, reply jsonrpc2.Replier, req jsonrpc2.Request) error {
	var params protocol.CompletionParams
	if err := json.Unmarshal(req.Params(), &params); err != nil {
		return reply(ctx, nil, fmt.Errorf("failed to unmarshal completion params: %w", err))
	}

	completions := s.handlers.Completion(params.TextDocument.URI, params.Position)
	return reply(ctx, completions, nil)
}

// publishDiagnostics sends diagnostics to the client
func (s *Server) publishDiagnostics(ctx context.Context, uri protocol.DocumentURI, diagnostics []protocol.Diagnostic) {
	params := protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: diagnostics,
	}

	if err := s.conn.Notify(ctx, "textDocument/publishDiagnostics", params); err != nil {
		s.logger.Error("Failed to publish diagnostics", zap.Error(err))
	}
}
