// Package lsp provides the Language Server Protocol server implementation.
package lsp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/example/lsp-server/internal/protocol"
)

// Server represents an LSP server.
type Server struct {
	// Reader for incoming messages
	reader *bufio.Reader
	// Writer for outgoing messages
	writer io.Writer
	// Logger for server events
	logger *log.Logger
	// Handler for processing LSP requests
	handler *Handler
	// Context for the server
	ctx context.Context
	// Cancel function for graceful shutdown
	cancel context.CancelFunc
	// Mutex for writing responses
	writeMu sync.Mutex
	// Whether the server has been initialized
	initialized bool
	// Whether the server is shutting down
	shuttingDown bool
}

// NewServer creates a new LSP server.
func NewServer(reader io.Reader, writer io.Writer, logger *log.Logger) *Server {
	if logger == nil {
		logger = log.New(os.Stderr, "[LSP] ", log.LstdFlags)
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		reader: bufio.NewReader(reader),
		writer: writer,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}

	s.handler = NewHandler(s, logger)

	return s
}

// Run starts the server and processes messages until shutdown.
func (s *Server) Run() error {
	s.logger.Println("LSP server starting...")

	for {
		select {
		case <-s.ctx.Done():
			s.logger.Println("Server context cancelled, shutting down")
			return nil
		default:
		}

		msg, err := s.readMessage()
		if err != nil {
			if err == io.EOF {
				s.logger.Println("Connection closed")
				return nil
			}
			s.logger.Printf("Error reading message: %v", err)
			continue
		}

		go s.handleMessage(msg)
	}
}

// Shutdown gracefully shuts down the server.
func (s *Server) Shutdown() {
	s.shuttingDown = true
	s.cancel()
}

// readMessage reads a single LSP message from the input.
func (s *Server) readMessage() (json.RawMessage, error) {
	// Read headers
	var contentLength int
	for {
		line, err := s.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			break // End of headers
		}

		if strings.HasPrefix(line, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, err = strconv.Atoi(lengthStr)
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %v", err)
			}
		}
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	// Read content
	content := make([]byte, contentLength)
	_, err := io.ReadFull(s.reader, content)
	if err != nil {
		return nil, fmt.Errorf("error reading content: %v", err)
	}

	s.logger.Printf("Received: %s", string(content))

	return json.RawMessage(content), nil
}

// handleMessage processes a single LSP message.
func (s *Server) handleMessage(content json.RawMessage) {
	// First, try to determine if this is a request or notification
	var baseMsg struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      interface{}     `json:"id,omitempty"`
		Method  string          `json:"method,omitempty"`
		Params  json.RawMessage `json:"params,omitempty"`
	}

	if err := json.Unmarshal(content, &baseMsg); err != nil {
		s.logger.Printf("Error parsing message: %v", err)
		s.sendError(nil, protocol.ParseError, "Parse error", nil)
		return
	}

	// Check if this is a request (has ID) or notification (no ID)
	if baseMsg.ID != nil {
		// This is a request
		s.handleRequest(baseMsg.ID, baseMsg.Method, baseMsg.Params)
	} else {
		// This is a notification
		s.handleNotification(baseMsg.Method, baseMsg.Params)
	}
}

// handleRequest processes an LSP request and sends a response.
func (s *Server) handleRequest(id interface{}, method string, params json.RawMessage) {
	s.logger.Printf("Handling request: %s (id: %v)", method, id)

	// Check if server is initialized for non-initialize requests
	if !s.initialized && method != "initialize" {
		s.sendError(id, protocol.ServerNotInitialized, "Server not initialized", nil)
		return
	}

	result, err := s.handler.HandleRequest(s.ctx, method, params)
	if err != nil {
		s.logger.Printf("Error handling request %s: %v", method, err)
		if rpcErr, ok := err.(*RPCError); ok {
			s.sendError(id, rpcErr.Code, rpcErr.Message, rpcErr.Data)
		} else {
			s.sendError(id, protocol.InternalError, err.Error(), nil)
		}
		return
	}

	s.sendResult(id, result)
}

// handleNotification processes an LSP notification.
func (s *Server) handleNotification(method string, params json.RawMessage) {
	s.logger.Printf("Handling notification: %s", method)

	if err := s.handler.HandleNotification(s.ctx, method, params); err != nil {
		s.logger.Printf("Error handling notification %s: %v", method, err)
	}
}

// sendResult sends a successful response.
func (s *Server) sendResult(id interface{}, result interface{}) {
	resultJSON, err := json.Marshal(result)
	if err != nil {
		s.logger.Printf("Error marshaling result: %v", err)
		s.sendError(id, protocol.InternalError, "Error marshaling result", nil)
		return
	}

	response := protocol.ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Result:  resultJSON,
	}

	s.sendMessage(response)
}

// sendError sends an error response.
func (s *Server) sendError(id interface{}, code protocol.ErrorCode, message string, data interface{}) {
	response := protocol.ResponseMessage{
		JSONRPC: "2.0",
		ID:      id,
		Error: &protocol.ResponseError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	s.sendMessage(response)
}

// SendNotification sends a notification to the client.
func (s *Server) SendNotification(method string, params interface{}) error {
	paramsJSON, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("error marshaling params: %v", err)
	}

	notification := protocol.NotificationMessage{
		JSONRPC: "2.0",
		Method:  method,
		Params:  paramsJSON,
	}

	return s.sendMessage(notification)
}

// sendMessage writes a message to the output.
func (s *Server) sendMessage(msg interface{}) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()

	content, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	s.logger.Printf("Sending: %s", string(content))

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(content))

	if _, err := s.writer.Write([]byte(header)); err != nil {
		return fmt.Errorf("error writing header: %v", err)
	}

	if _, err := s.writer.Write(content); err != nil {
		return fmt.Errorf("error writing content: %v", err)
	}

	return nil
}

// SetInitialized marks the server as initialized.
func (s *Server) SetInitialized(initialized bool) {
	s.initialized = initialized
}

// IsInitialized returns whether the server has been initialized.
func (s *Server) IsInitialized() bool {
	return s.initialized
}

// RPCError represents a JSON-RPC error.
type RPCError struct {
	Code    protocol.ErrorCode
	Message string
	Data    interface{}
}

func (e *RPCError) Error() string {
	return e.Message
}

// NewRPCError creates a new RPC error.
func NewRPCError(code protocol.ErrorCode, message string) *RPCError {
	return &RPCError{
		Code:    code,
		Message: message,
	}
}
