package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/rkkeerth/lsp-server/lsp"
	"github.com/rkkeerth/lsp-server/rpc"
)

// Server handles LSP communication
type Server struct {
	state    *State
	reader   *bufio.Reader
	writer   io.Writer
	logger   *log.Logger
	shutdown bool
}

// NewServer creates a new LSP server
func NewServer(reader io.Reader, writer io.Writer, logger *log.Logger) *Server {
	return &Server{
		state:  NewState(),
		reader: bufio.NewReader(reader),
		writer: writer,
		logger: logger,
	}
}

// Run starts the main server loop
func (s *Server) Run() error {
	s.logger.Println("LSP server starting...")

	for {
		method, content, err := rpc.DecodeMessage(s.reader)
		if err != nil {
			if err == io.EOF {
				s.logger.Println("Client disconnected")
				return nil
			}
			s.logger.Printf("Error reading message: %v", err)
			continue
		}

		s.logger.Printf("Received: %s", method)

		// Check if this is the exit notification
		if method == lsp.MethodExit {
			s.logger.Println("Received exit notification, shutting down")
			return nil
		}

		// Handle the message
		response, err := s.handleMessage(method, content)
		if err != nil {
			s.logger.Printf("Error handling %s: %v", method, err)
			continue
		}

		// Send response if there is one (requests get responses, notifications don't)
		if response != nil {
			if err := s.sendResponse(response); err != nil {
				s.logger.Printf("Error sending response: %v", err)
			}
		}
	}
}

// handleMessage routes messages to the appropriate handler
func (s *Server) handleMessage(method string, content []byte) (*rpc.Response, error) {
	// Extract request ID if present (requests have IDs, notifications don't)
	var req struct {
		ID interface{} `json:"id"`
	}
	json.Unmarshal(content, &req)

	switch method {
	case lsp.MethodInitialize:
		var params lsp.InitializeParams
		if err := s.unmarshalParams(content, &params); err != nil {
			return s.errorResponse(req.ID, rpc.InvalidParams, err.Error()), nil
		}
		result, err := HandleInitialize(&params)
		if err != nil {
			return s.errorResponse(req.ID, rpc.InternalError, err.Error()), nil
		}
		return s.successResponse(req.ID, result), nil

	case lsp.MethodInitialized:
		var params lsp.InitializedParams
		s.unmarshalParams(content, &params)
		HandleInitialized(&params)
		return nil, nil // Notification - no response

	case lsp.MethodShutdown:
		s.shutdown = true
		result := HandleShutdown()
		return s.successResponse(req.ID, result), nil

	case lsp.MethodTextDocumentDidOpen:
		var params lsp.DidOpenTextDocumentParams
		if err := s.unmarshalParams(content, &params); err != nil {
			s.logger.Printf("Error parsing didOpen params: %v", err)
			return nil, nil
		}
		HandleTextDocumentDidOpen(s.state, &params)
		return nil, nil // Notification - no response

	case lsp.MethodTextDocumentDidChange:
		var params lsp.DidChangeTextDocumentParams
		if err := s.unmarshalParams(content, &params); err != nil {
			s.logger.Printf("Error parsing didChange params: %v", err)
			return nil, nil
		}
		HandleTextDocumentDidChange(s.state, &params)
		return nil, nil // Notification - no response

	case lsp.MethodTextDocumentDidClose:
		var params lsp.DidCloseTextDocumentParams
		if err := s.unmarshalParams(content, &params); err != nil {
			s.logger.Printf("Error parsing didClose params: %v", err)
			return nil, nil
		}
		HandleTextDocumentDidClose(s.state, &params)
		return nil, nil // Notification - no response

	default:
		s.logger.Printf("Unknown method: %s", method)
		if req.ID != nil {
			return s.errorResponse(req.ID, rpc.MethodNotFound, fmt.Sprintf("method not found: %s", method)), nil
		}
		return nil, nil
	}
}

// unmarshalParams extracts params from a JSON-RPC message
func (s *Server) unmarshalParams(content []byte, params interface{}) error {
	var msg struct {
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(content, &msg); err != nil {
		return err
	}
	if msg.Params == nil {
		return nil
	}
	return json.Unmarshal(msg.Params, params)
}

// successResponse creates a successful JSON-RPC response
func (s *Server) successResponse(id interface{}, result interface{}) *rpc.Response {
	return &rpc.Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// errorResponse creates an error JSON-RPC response
func (s *Server) errorResponse(id interface{}, code int, message string) *rpc.Response {
	return &rpc.Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &rpc.ResponseError{
			Code:    code,
			Message: message,
		},
	}
}

// sendResponse encodes and writes a response
func (s *Server) sendResponse(response *rpc.Response) error {
	encoded, err := rpc.EncodeMessage(response)
	if err != nil {
		return err
	}
	_, err = s.writer.Write(encoded)
	return err
}
