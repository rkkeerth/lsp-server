package jsonrpc

import (
	"encoding/json"
)

const (
	// JSONRPCVersion is the JSON-RPC version
	JSONRPCVersion = "2.0"
)

// Message represents a base JSON-RPC message
type Message struct {
	JSONRPC string `json:"jsonrpc"`
}

// Request represents a JSON-RPC request
type Request struct {
	Message
	ID     interface{}     `json:"id"`
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// Response represents a JSON-RPC response
type Response struct {
	Message
	ID     interface{}     `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *ResponseError  `json:"error,omitempty"`
}

// Notification represents a JSON-RPC notification (no ID)
type Notification struct {
	Message
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
}

// ResponseError represents a JSON-RPC error
type ResponseError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// Error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

// NewRequest creates a new JSON-RPC request
func NewRequest(id interface{}, method string, params interface{}) (*Request, error) {
	var paramsJSON json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		paramsJSON = data
	}

	return &Request{
		Message: Message{JSONRPC: JSONRPCVersion},
		ID:      id,
		Method:  method,
		Params:  paramsJSON,
	}, nil
}

// NewResponse creates a new JSON-RPC response
func NewResponse(id interface{}, result interface{}) (*Response, error) {
	var resultJSON json.RawMessage
	if result != nil {
		data, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		resultJSON = data
	}

	return &Response{
		Message: Message{JSONRPC: JSONRPCVersion},
		ID:      id,
		Result:  resultJSON,
	}, nil
}

// NewErrorResponse creates a new JSON-RPC error response
func NewErrorResponse(id interface{}, code int, message string) *Response {
	return &Response{
		Message: Message{JSONRPC: JSONRPCVersion},
		ID:      id,
		Error: &ResponseError{
			Code:    code,
			Message: message,
		},
	}
}

// NewNotification creates a new JSON-RPC notification
func NewNotification(method string, params interface{}) (*Notification, error) {
	var paramsJSON json.RawMessage
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		paramsJSON = data
	}

	return &Notification{
		Message: Message{JSONRPC: JSONRPCVersion},
		Method:  method,
		Params:  paramsJSON,
	}, nil
}
