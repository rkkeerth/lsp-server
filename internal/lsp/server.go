package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
)

// Server represents an LSP server
type Server struct {
	reader   *bufio.Reader
	writer   io.Writer
	mu       sync.RWMutex
	state    ServerState
	documents map[string]*Document
}

// ServerState represents the server's state
type ServerState int

const (
	StateUninitialized ServerState = iota
	StateInitializing
	StateInitialized
	StateShuttingDown
	StateShutdown
)

// Document represents an open text document
type Document struct {
	URI        string
	LanguageID string
	Version    int
	Content    string
}

// NewServer creates a new LSP server
func NewServer(reader io.Reader, writer io.Writer) *Server {
	return &Server{
		reader:    bufio.NewReader(reader),
		writer:    writer,
		state:     StateUninitialized,
		documents: make(map[string]*Document),
	}
}

// Start begins processing LSP messages
func (s *Server) Start() error {
	log.Println("LSP Server starting...")
	
	for {
		msg, err := s.readMessage()
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
				return nil
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		// Process message in a goroutine to handle concurrent requests
		go s.handleMessage(msg)
	}
}

// readMessage reads a single LSP message from the input stream
func (s *Server) readMessage() ([]byte, error) {
	// Read headers
	headers := make(map[string]string)
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = line[:len(line)-1] // Remove \n
		if len(line) > 0 && line[len(line)-1] == '\r' {
			line = line[:len(line)-1] // Remove \r
		}

		if line == "" {
			break // End of headers
		}

		// Parse header
		var key, value string
		if _, err := fmt.Sscanf(line, "%s %s", &key, &value); err != nil {
			continue
		}
		headers[key[:len(key)-1]] = value // Remove colon from key
	}

	// Get content length
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	var contentLength int
	if _, err := fmt.Sscanf(contentLengthStr, "%d", &contentLength); err != nil {
		return nil, fmt.Errorf("invalid Content-Length: %v", err)
	}

	// Read content
	content := make([]byte, contentLength)
	if _, err := io.ReadFull(s.reader, content); err != nil {
		return nil, err
	}

	return content, nil
}

// writeMessage writes an LSP message to the output stream
func (s *Server) writeMessage(msg interface{}) error {
	content, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	response := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
	
	if _, err := s.writer.Write([]byte(response)); err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}

// sendResponse sends a response message
func (s *Server) sendResponse(id interface{}, result interface{}) error {
	response := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	return s.writeMessage(response)
}

// sendError sends an error response
func (s *Server) sendError(id interface{}, code int, message string) error {
	response := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
		},
	}
	return s.writeMessage(response)
}

// handleMessage dispatches messages to appropriate handlers
func (s *Server) handleMessage(data []byte) {
	// Try to parse as a request first
	var request Request
	if err := json.Unmarshal(data, &request); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return
	}

	log.Printf("Received %s: %s", getMessageType(request.ID), request.Method)

	// Route to appropriate handler
	switch request.Method {
	case "initialize":
		s.handleInitialize(request)
	case "initialized":
		s.handleInitialized(request)
	case "shutdown":
		s.handleShutdown(request)
	case "exit":
		s.handleExit(request)
	case "textDocument/didOpen":
		s.handleTextDocumentDidOpen(request)
	case "textDocument/didChange":
		s.handleTextDocumentDidChange(request)
	case "textDocument/didClose":
		s.handleTextDocumentDidClose(request)
	default:
		if request.ID != nil {
			log.Printf("Method not found: %s", request.Method)
			s.sendError(request.ID, MethodNotFound, fmt.Sprintf("Method not found: %s", request.Method))
		} else {
			log.Printf("Ignoring unknown notification: %s", request.Method)
		}
	}
}

// getMessageType returns a string describing the message type
func getMessageType(id interface{}) string {
	if id == nil {
		return "notification"
	}
	return "request"
}

// GetDocument retrieves a document by URI
func (s *Server) GetDocument(uri string) (*Document, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	doc, ok := s.documents[uri]
	return doc, ok
}

// SetDocument stores or updates a document
func (s *Server) SetDocument(uri string, doc *Document) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents[uri] = doc
}

// DeleteDocument removes a document
func (s *Server) DeleteDocument(uri string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.documents, uri)
}

// GetState returns the current server state
func (s *Server) GetState() ServerState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// SetState updates the server state
func (s *Server) SetState(state ServerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = state
}
