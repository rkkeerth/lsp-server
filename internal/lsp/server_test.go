package lsp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestNewServer(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	writer := &bytes.Buffer{}
	
	server := NewServer(reader, writer)
	
	if server == nil {
		t.Fatal("NewServer returned nil")
	}
	
	if server.state != StateUninitialized {
		t.Errorf("Expected initial state to be Uninitialized, got %v", server.state)
	}
	
	if len(server.documents) != 0 {
		t.Errorf("Expected documents to be empty, got %d documents", len(server.documents))
	}
}

func TestServerState(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	writer := &bytes.Buffer{}
	server := NewServer(reader, writer)
	
	// Test initial state
	if server.GetState() != StateUninitialized {
		t.Errorf("Expected StateUninitialized, got %v", server.GetState())
	}
	
	// Test state change
	server.SetState(StateInitialized)
	if server.GetState() != StateInitialized {
		t.Errorf("Expected StateInitialized, got %v", server.GetState())
	}
}

func TestDocumentManagement(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	writer := &bytes.Buffer{}
	server := NewServer(reader, writer)
	
	// Test document doesn't exist initially
	_, exists := server.GetDocument("file:///test.txt")
	if exists {
		t.Error("Document should not exist initially")
	}
	
	// Test adding document
	doc := &Document{
		URI:        "file:///test.txt",
		LanguageID: "plaintext",
		Version:    1,
		Content:    "Hello, World!",
	}
	server.SetDocument(doc.URI, doc)
	
	// Test retrieving document
	retrieved, exists := server.GetDocument("file:///test.txt")
	if !exists {
		t.Fatal("Document should exist after adding")
	}
	
	if retrieved.URI != doc.URI {
		t.Errorf("Expected URI %s, got %s", doc.URI, retrieved.URI)
	}
	
	if retrieved.Content != doc.Content {
		t.Errorf("Expected content %s, got %s", doc.Content, retrieved.Content)
	}
	
	// Test deleting document
	server.DeleteDocument(doc.URI)
	_, exists = server.GetDocument("file:///test.txt")
	if exists {
		t.Error("Document should not exist after deletion")
	}
}

func TestWriteMessage(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	writer := &bytes.Buffer{}
	server := NewServer(reader, writer)
	
	msg := Response{
		JSONRPC: "2.0",
		ID:      1,
		Result:  "test",
	}
	
	err := server.writeMessage(msg)
	if err != nil {
		t.Fatalf("writeMessage failed: %v", err)
	}
	
	output := writer.String()
	
	// Check for Content-Length header
	if !bytes.Contains([]byte(output), []byte("Content-Length:")) {
		t.Error("Output should contain Content-Length header")
	}
	
	// Check for JSON content
	if !bytes.Contains([]byte(output), []byte(`"jsonrpc":"2.0"`)) {
		t.Error("Output should contain JSON-RPC version")
	}
}

func TestSendResponse(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	writer := &bytes.Buffer{}
	server := NewServer(reader, writer)
	
	err := server.sendResponse(1, map[string]string{"status": "ok"})
	if err != nil {
		t.Fatalf("sendResponse failed: %v", err)
	}
	
	output := writer.String()
	
	// Parse the response
	lines := bytes.Split([]byte(output), []byte("\r\n"))
	if len(lines) < 3 {
		t.Fatalf("Expected at least 3 lines (header, blank, content), got %d", len(lines))
	}
	
	// Check Content-Length
	headerLine := string(lines[0])
	if !bytes.Contains(lines[0], []byte("Content-Length:")) {
		t.Errorf("Expected Content-Length header, got: %s", headerLine)
	}
	
	// Find the JSON content (after the blank line)
	var jsonContent []byte
	for i, line := range lines {
		if len(line) == 0 && i+1 < len(lines) {
			jsonContent = lines[i+1]
			break
		}
	}
	
	// Parse JSON
	var response Response
	if err := json.Unmarshal(jsonContent, &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}
	
	if response.JSONRPC != "2.0" {
		t.Errorf("Expected jsonrpc 2.0, got %s", response.JSONRPC)
	}
	
	if response.ID != float64(1) { // JSON unmarshals numbers as float64
		t.Errorf("Expected ID 1, got %v", response.ID)
	}
}

func TestSendError(t *testing.T) {
	reader := bytes.NewReader([]byte{})
	writer := &bytes.Buffer{}
	server := NewServer(reader, writer)
	
	err := server.sendError(1, MethodNotFound, "Method not found")
	if err != nil {
		t.Fatalf("sendError failed: %v", err)
	}
	
	output := writer.String()
	
	// Find JSON content
	lines := bytes.Split([]byte(output), []byte("\r\n"))
	var jsonContent []byte
	for i, line := range lines {
		if len(line) == 0 && i+1 < len(lines) {
			jsonContent = lines[i+1]
			break
		}
	}
	
	// Parse JSON
	var response Response
	if err := json.Unmarshal(jsonContent, &response); err != nil {
		t.Fatalf("Failed to parse error response JSON: %v", err)
	}
	
	if response.Error == nil {
		t.Fatal("Expected error field to be present")
	}
	
	if response.Error.Code != MethodNotFound {
		t.Errorf("Expected error code %d, got %d", MethodNotFound, response.Error.Code)
	}
	
	if response.Error.Message != "Method not found" {
		t.Errorf("Expected message 'Method not found', got '%s'", response.Error.Message)
	}
}

func TestReadMessage(t *testing.T) {
	// Prepare a valid LSP message
	jsonContent := `{"jsonrpc":"2.0","id":1,"method":"test"}`
	message := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(jsonContent), jsonContent)
	
	reader := bytes.NewReader([]byte(message))
	writer := &bytes.Buffer{}
	server := NewServer(reader, writer)
	
	// Read the message
	data, err := server.readMessage()
	if err != nil {
		t.Fatalf("readMessage failed: %v", err)
	}
	
	// Parse the message
	var request Request
	if err := json.Unmarshal(data, &request); err != nil {
		t.Fatalf("Failed to parse message: %v", err)
	}
	
	if request.JSONRPC != "2.0" {
		t.Errorf("Expected jsonrpc 2.0, got %s", request.JSONRPC)
	}
	
	if request.Method != "test" {
		t.Errorf("Expected method 'test', got '%s'", request.Method)
	}
}

func TestGetMessageType(t *testing.T) {
	tests := []struct {
		name     string
		id       interface{}
		expected string
	}{
		{"Request with number ID", 1, "request"},
		{"Request with string ID", "abc", "request"},
		{"Notification", nil, "notification"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMessageType(tt.id)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
