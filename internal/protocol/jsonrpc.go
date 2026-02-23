package protocol

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Request represents a JSON-RPC 2.0 request message.
// The id can be either an int or a string.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response message.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// Notification represents a JSON-RPC 2.0 notification message.
// Notifications do not have an id field.
type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// ResponseError represents a JSON-RPC 2.0 error object.
type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Error implements the error interface for ResponseError.
func (e *ResponseError) Error() string {
	return fmt.Sprintf("JSON-RPC error %d: %s", e.Code, e.Message)
}

// Standard JSON-RPC 2.0 error codes.
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// LSP specific error codes.
const (
	ServerNotInitialized = -32002
	UnknownErrorCode     = -32001
	RequestCancelled     = -32800
	ContentModified      = -32801
)

// ReadMessage reads a JSON-RPC message from the reader.
// It expects the Content-Length header followed by the JSON body.
func ReadMessage(reader io.Reader) ([]byte, error) {
	r := bufio.NewReader(reader)

	// Read headers until we find an empty line
	var contentLength int
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read header: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" {
			// End of headers
			break
		}

		if strings.HasPrefix(line, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			contentLength, err = strconv.Atoi(lengthStr)
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %w", err)
			}
		}
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	// Read the body
	body := make([]byte, contentLength)
	_, err := io.ReadFull(r, body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return body, nil
}

// WriteMessage writes a JSON-RPC message to the writer.
// It adds the Content-Length header before the JSON body.
func WriteMessage(writer io.Writer, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	_, err = writer.Write([]byte(header))
	if err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	_, err = writer.Write(body)
	if err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	return nil
}
