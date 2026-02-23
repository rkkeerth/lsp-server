package rpc

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestEncodeMessage(t *testing.T) {
	tests := []struct {
		name    string
		msg     interface{}
		wantErr bool
	}{
		{
			name: "simple message",
			msg: map[string]string{
				"jsonrpc": "2.0",
				"method":  "test",
			},
			wantErr: false,
		},
		{
			name: "response message",
			msg: Response{
				JSONRPC: "2.0",
				ID:      1,
				Result:  "ok",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := EncodeMessage(tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Check Content-Length header format
			result := string(encoded)
			if !strings.HasPrefix(result, "Content-Length: ") {
				t.Errorf("Missing Content-Length header: %s", result)
			}
			if !strings.Contains(result, "\r\n\r\n") {
				t.Errorf("Missing header separator: %s", result)
			}
		})
	}
}

func TestEncodeMessageContentLength(t *testing.T) {
	msg := map[string]string{"test": "value"}
	encoded, err := EncodeMessage(msg)
	if err != nil {
		t.Fatalf("EncodeMessage() error = %v", err)
	}

	// Parse out the content length
	result := string(encoded)
	parts := strings.SplitN(result, "\r\n\r\n", 2)
	if len(parts) != 2 {
		t.Fatalf("Invalid message format: %s", result)
	}

	// Extract declared length
	var declaredLen int
	_, err = fmt.Sscanf(parts[0], "Content-Length: %d", &declaredLen)
	if err != nil {
		t.Fatalf("Failed to parse Content-Length: %v", err)
	}

	// Verify actual content length matches
	if declaredLen != len(parts[1]) {
		t.Errorf("Content-Length mismatch: declared %d, actual %d", declaredLen, len(parts[1]))
	}
}

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantMethod string
		wantErr    bool
	}{
		{
			name:       "valid request",
			input:      "Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":1}",
			wantMethod: "test",
			wantErr:    false,
		},
		{
			name:       "valid notification",
			input:      "Content-Length: 35\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"notify\"}",
			wantMethod: "notify",
			wantErr:    false,
		},
		{
			name:       "invalid header",
			input:      "Invalid-Header: 10\r\n\r\n{\"test\":1}",
			wantMethod: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			method, _, err := DecodeMessage(reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if method != tt.wantMethod {
				t.Errorf("DecodeMessage() method = %v, want %v", method, tt.wantMethod)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	original := Request{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	// Encode
	encoded, err := EncodeMessage(original)
	if err != nil {
		t.Fatalf("EncodeMessage() error = %v", err)
	}

	// Decode
	reader := bufio.NewReader(bytes.NewReader(encoded))
	method, content, err := DecodeMessage(reader)
	if err != nil {
		t.Fatalf("DecodeMessage() error = %v", err)
	}

	if method != original.Method {
		t.Errorf("Method mismatch: got %v, want %v", method, original.Method)
	}

	// Verify content can be unmarshaled
	var decoded Request
	if err := json.Unmarshal(content, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal content: %v", err)
	}

	if decoded.Method != original.Method {
		t.Errorf("Decoded method mismatch: got %v, want %v", decoded.Method, original.Method)
	}
	// Note: JSON numbers unmarshal as float64, so we compare as float64
	if decoded.ID.(float64) != float64(original.ID.(int)) {
		t.Errorf("Decoded ID mismatch: got %v, want %v", decoded.ID, original.ID)
	}
}
