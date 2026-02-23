package server

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/rkkeerth/lsp-server/internal/document"
	"github.com/rkkeerth/lsp-server/internal/handler"
	"github.com/rkkeerth/lsp-server/internal/protocol"
)

// Server is the LSP server
type Server struct {
	reader     io.Reader
	writer     io.Writer
	logger     *log.Logger
	store      *document.Store
	lifecycle  *handler.Lifecycle
	textDoc    *handler.TextDocument
	completion *handler.Completion
	hover      *handler.Hover
}

// New creates a new LSP server
func New(reader io.Reader, writer io.Writer, logger *log.Logger) *Server {
	store := document.NewStore()
	return &Server{
		reader:     reader,
		writer:     writer,
		logger:     logger,
		store:      store,
		lifecycle:  handler.NewLifecycle(),
		textDoc:    handler.NewTextDocument(store),
		completion: handler.NewCompletion(store),
		hover:      handler.NewHover(store),
	}
}

// Run starts the server main loop
func (s *Server) Run() error {
	s.logger.Println("LSP server started")

	for {
		data, err := protocol.ReadMessage(s.reader)
		if err != nil {
			if err == io.EOF {
				s.logger.Println("Connection closed")
				return nil
			}
			s.logger.Printf("Error reading message: %v", err)
			continue
		}

		method := s.parseMethod(data)
		s.handleMessage(data, method)
	}
}

// parseMethod extracts the method from a JSON-RPC message
func (s *Server) parseMethod(data []byte) string {
	var msg struct {
		Method string `json:"method"`
	}
	json.Unmarshal(data, &msg)
	return msg.Method
}

// isRequest checks if the message is a request (has an id field)
func (s *Server) isRequest(data []byte) bool {
	var msg struct {
		ID interface{} `json:"id"`
	}
	json.Unmarshal(data, &msg)
	return msg.ID != nil
}

func (s *Server) handleMessage(data []byte, method string) {
	if s.isRequest(data) {
		var req protocol.Request
		if err := json.Unmarshal(data, &req); err != nil {
			s.logger.Printf("Error parsing request: %v", err)
			return
		}
		s.handleRequest(&req)
	} else {
		var notif protocol.Notification
		if err := json.Unmarshal(data, &notif); err != nil {
			s.logger.Printf("Error parsing notification: %v", err)
			return
		}
		s.handleNotification(&notif)
	}
}

func (s *Server) handleRequest(req *protocol.Request) {
	s.logger.Printf("Request: %s (id: %v)", req.Method, req.ID)

	var result interface{}
	var respErr *protocol.ResponseError

	switch req.Method {
	case protocol.MethodInitialize:
		var params protocol.InitializeParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			respErr = &protocol.ResponseError{Code: protocol.InvalidParams, Message: err.Error()}
		} else {
			result = s.lifecycle.HandleInitialize(&params)
		}

	case protocol.MethodShutdown:
		s.lifecycle.HandleShutdown()
		result = nil

	case protocol.MethodTextDocumentCompletion:
		var params protocol.CompletionParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			respErr = &protocol.ResponseError{Code: protocol.InvalidParams, Message: err.Error()}
		} else {
			result = s.completion.Handle(&params)
		}

	case protocol.MethodTextDocumentHover:
		var params protocol.HoverParams
		if err := json.Unmarshal(req.Params, &params); err != nil {
			respErr = &protocol.ResponseError{Code: protocol.InvalidParams, Message: err.Error()}
		} else {
			result = s.hover.Handle(&params)
		}

	default:
		respErr = &protocol.ResponseError{
			Code:    protocol.MethodNotFound,
			Message: "method not found: " + req.Method,
		}
	}

	s.sendResponse(req.ID, result, respErr)
}

func (s *Server) handleNotification(notif *protocol.Notification) {
	s.logger.Printf("Notification: %s", notif.Method)

	switch notif.Method {
	case protocol.MethodInitialized:
		s.lifecycle.HandleInitialized()

	case protocol.MethodExit:
		exitCode := 0
		if !s.lifecycle.ShutdownRequested() {
			exitCode = 1
		}
		os.Exit(exitCode)

	case protocol.MethodTextDocumentDidOpen:
		var params protocol.DidOpenTextDocumentParams
		if err := json.Unmarshal(notif.Params, &params); err != nil {
			s.logger.Printf("Error parsing didOpen params: %v", err)
			return
		}
		s.textDoc.HandleDidOpen(&params)

	case protocol.MethodTextDocumentDidClose:
		var params protocol.DidCloseTextDocumentParams
		if err := json.Unmarshal(notif.Params, &params); err != nil {
			s.logger.Printf("Error parsing didClose params: %v", err)
			return
		}
		s.textDoc.HandleDidClose(&params)

	case protocol.MethodTextDocumentDidChange:
		var params protocol.DidChangeTextDocumentParams
		if err := json.Unmarshal(notif.Params, &params); err != nil {
			s.logger.Printf("Error parsing didChange params: %v", err)
			return
		}
		s.textDoc.HandleDidChange(&params)
	}
}

func (s *Server) sendResponse(id interface{}, result interface{}, respErr *protocol.ResponseError) {
	resp := &protocol.Response{
		JSONRPC: "2.0",
		ID:      id,
	}

	if respErr != nil {
		resp.Error = respErr
	} else {
		resultBytes, err := json.Marshal(result)
		if err != nil {
			resp.Error = &protocol.ResponseError{
				Code:    protocol.InternalError,
				Message: err.Error(),
			}
		} else {
			resp.Result = resultBytes
		}
	}

	if err := protocol.WriteMessage(s.writer, resp); err != nil {
		s.logger.Printf("Error writing response: %v", err)
	}
}
