// Package rpc implements JSON-RPC 2.0 protocol encoding and decoding
// with Content-Length headers as required by the LSP specification.
package rpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// BaseMessage represents a basic JSON-RPC 2.0 message
type BaseMessage struct {
	JSONRPC string `json:"jsonrpc"`
}

// Request represents a JSON-RPC 2.0 request message
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC 2.0 response message
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// Notification represents a JSON-RPC 2.0 notification message (no ID)
type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// ResponseError represents a JSON-RPC 2.0 error object
type ResponseError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Standard JSON-RPC error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// EncodeMessage encodes a message with the Content-Length header
// Format: "Content-Length: <length>\r\n\r\n<content>"
func EncodeMessage(msg interface{}) ([]byte, error) {
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message: %w", err)
	}
	return []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)), nil
}

// DecodeMessage reads and decodes a JSON-RPC message from the reader
// Returns the method name (for routing) and the raw content bytes
func DecodeMessage(reader *bufio.Reader) (string, []byte, error) {
	// Read the Content-Length header
	header, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Parse Content-Length
	contentLength, err := parseContentLength(header)
	if err != nil {
		return "", nil, err
	}

	// Read the empty line separator
	separator, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, fmt.Errorf("failed to read header separator: %w", err)
	}
	if separator != "\r\n" && separator != "\n" {
		return "", nil, fmt.Errorf("expected empty line after header, got: %q", separator)
	}

	// Read the content - use io.ReadFull to ensure we read exactly contentLength bytes
	content := make([]byte, contentLength)
	_, err = io.ReadFull(reader, content)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read content: %w", err)
	}

	// Extract the method name for routing
	method, err := extractMethod(content)
	if err != nil {
		return "", nil, err
	}

	return method, content, nil
}

// parseContentLength extracts the content length from the header line
func parseContentLength(header string) (int, error) {
	// Remove trailing whitespace
	header = trimRight(header)

	const prefix = "Content-Length: "
	if len(header) < len(prefix) || header[:len(prefix)] != prefix {
		return 0, fmt.Errorf("invalid header format: %q", header)
	}

	lengthStr := header[len(prefix):]
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		return 0, fmt.Errorf("invalid content length: %q", lengthStr)
	}

	return length, nil
}

// trimRight removes trailing \r and \n from a string
func trimRight(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\r' || s[len(s)-1] == '\n') {
		s = s[:len(s)-1]
	}
	return s
}

// extractMethod extracts the method name from a JSON-RPC message
func extractMethod(content []byte) (string, error) {
	var msg struct {
		Method string `json:"method"`
	}
	if err := json.Unmarshal(content, &msg); err != nil {
		return "", fmt.Errorf("failed to parse method: %w", err)
	}
	return msg.Method, nil
}
