package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Transport handles JSON-RPC 2.0 message framing over stdin/stdout.
// Messages are framed with Content-Length headers as per LSP specification.
type Transport struct {
	reader *bufio.Reader
	writer io.Writer
}

// NewTransport creates a new Transport.
func NewTransport(reader io.Reader, writer io.Writer) *Transport {
	return &Transport{
		reader: bufio.NewReader(reader),
		writer: writer,
	}
}

// Read reads a JSON-RPC message from the transport.
// It parses the Content-Length header and reads the message body.
func (t *Transport) Read() (*Request, error) {
	// Read headers until we get an empty line
	var contentLength int
	for {
		line, err := t.reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read header: %w", err)
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			// End of headers
			break
		}

		// Parse header
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		if strings.ToLower(parts[0]) == "content-length" {
			var err error
			contentLength, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %w", err)
			}
		}
	}

	if contentLength == 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	// Read the message body
	body := make([]byte, contentLength)
	_, err := io.ReadFull(t.reader, body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	// Parse the JSON-RPC request
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("failed to parse request: %w", err)
	}

	return &req, nil
}

// Write writes a JSON-RPC response to the transport.
// It adds the Content-Length header before the message body.
func (t *Transport) Write(resp *Response) error {
	body, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	// Write headers
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err := t.writer.Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write body
	if _, err := t.writer.Write(body); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	return nil
}

// WriteNotification writes a JSON-RPC notification to the transport.
func (t *Transport) WriteNotification(method string, params interface{}) error {
	notif := struct {
		JSONRPC string      `json:"jsonrpc"`
		Method  string      `json:"method"`
		Params  interface{} `json:"params,omitempty"`
	}{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err := t.writer.Write([]byte(header)); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if _, err := t.writer.Write(body); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	return nil
}
