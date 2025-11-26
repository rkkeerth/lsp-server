package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Handler handles JSON-RPC requests
type Handler interface {
	Handle(method string, params json.RawMessage) (interface{}, error)
}

// RPC manages JSON-RPC communication
type RPC struct {
	reader  *bufio.Reader
	writer  io.Writer
	handler Handler
}

// NewRPC creates a new RPC instance
func NewRPC(reader io.Reader, writer io.Writer, handler Handler) *RPC {
	return &RPC{
		reader:  bufio.NewReader(reader),
		writer:  writer,
		handler: handler,
	}
}

// ReadMessage reads a JSON-RPC message from the reader
func (r *RPC) ReadMessage() ([]byte, error) {
	// Read headers
	headers := make(map[string]string)
	for {
		line, err := r.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			// Empty line indicates end of headers
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header: %s", line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	// Get content length
	contentLengthStr, ok := headers["Content-Length"]
	if !ok {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	contentLength, err := strconv.Atoi(contentLengthStr)
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Length: %s", contentLengthStr)
	}

	// Read content
	content := make([]byte, contentLength)
	_, err = io.ReadFull(r.reader, content)
	if err != nil {
		return nil, err
	}

	return content, nil
}

// WriteMessage writes a JSON-RPC message to the writer
func (r *RPC) WriteMessage(message interface{}) error {
	content, err := json.Marshal(message)
	if err != nil {
		return err
	}

	response := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
	_, err = r.writer.Write([]byte(response))
	return err
}

// ProcessMessage processes an incoming JSON-RPC message
func (r *RPC) ProcessMessage(content []byte) error {
	// Try to parse as a request first
	var req Request
	if err := json.Unmarshal(content, &req); err == nil && req.Method != "" {
		// Check if it's a request (has ID) or notification (no ID)
		if req.ID != nil {
			return r.handleRequest(&req)
		}
		return r.handleNotification(&Notification{
			Message: req.Message,
			Method:  req.Method,
			Params:  req.Params,
		})
	}

	return fmt.Errorf("invalid message format")
}

// handleRequest handles a JSON-RPC request
func (r *RPC) handleRequest(req *Request) error {
	result, err := r.handler.Handle(req.Method, req.Params)

	var response *Response
	if err != nil {
		response = NewErrorResponse(req.ID, InternalError, err.Error())
	} else {
		response, err = NewResponse(req.ID, result)
		if err != nil {
			response = NewErrorResponse(req.ID, InternalError, err.Error())
		}
	}

	return r.WriteMessage(response)
}

// handleNotification handles a JSON-RPC notification
func (r *RPC) handleNotification(notif *Notification) error {
	// Notifications don't expect a response
	_, err := r.handler.Handle(notif.Method, notif.Params)
	return err
}

// Run starts the RPC message loop
func (r *RPC) Run() error {
	for {
		content, err := r.ReadMessage()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err := r.ProcessMessage(content); err != nil {
			// Log error but continue processing
			fmt.Fprintf(io.Discard, "Error processing message: %v\n", err)
		}
	}
}
