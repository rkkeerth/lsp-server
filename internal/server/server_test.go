package server

import (
	"encoding/json"
	"testing"

	"github.com/rkkeerth/lsp-server/internal/protocol"
)

func TestServerInitialize(t *testing.T) {
	server := NewServer()

	initParams := protocol.InitializeParams{
		ProcessID: intPtr(12345),
		RootURI:   stringPtr("file:///test"),
		Capabilities: protocol.ClientCapabilities{
			TextDocument: &protocol.TextDocumentClientCapabilities{
				Synchronization: &protocol.TextDocumentSyncClientCapabilities{
					DynamicRegistration: true,
				},
			},
		},
	}

	paramsJSON, err := json.Marshal(initParams)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	result, err := server.handleInitialize(paramsJSON)
	if err != nil {
		t.Fatalf("handleInitialize failed: %v", err)
	}

	initResult, ok := result.(protocol.InitializeResult)
	if !ok {
		t.Fatalf("Expected InitializeResult, got %T", result)
	}

	if initResult.Capabilities.TextDocumentSync == nil {
		t.Error("Expected TextDocumentSync capabilities to be set")
	}

	if initResult.ServerInfo == nil {
		t.Error("Expected ServerInfo to be set")
	}

	if initResult.ServerInfo.Name != "basic-lsp-server" {
		t.Errorf("Expected server name 'basic-lsp-server', got %s", initResult.ServerInfo.Name)
	}
}

func TestServerDidOpen(t *testing.T) {
	server := NewServer()

	didOpenParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.txt",
			LanguageID: "plaintext",
			Version:    1,
			Text:       "Hello, World!",
		},
	}

	paramsJSON, err := json.Marshal(didOpenParams)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	_, err = server.handleDidOpen(paramsJSON)
	if err != nil {
		t.Fatalf("handleDidOpen failed: %v", err)
	}

	doc, exists := server.GetDocument("file:///test.txt")
	if !exists {
		t.Fatal("Expected document to exist")
	}

	if doc.URI != "file:///test.txt" {
		t.Errorf("Expected URI 'file:///test.txt', got %s", doc.URI)
	}

	if doc.Text != "Hello, World!" {
		t.Errorf("Expected text 'Hello, World!', got %s", doc.Text)
	}

	if doc.Version != 1 {
		t.Errorf("Expected version 1, got %d", doc.Version)
	}
}

func TestServerDidChange(t *testing.T) {
	server := NewServer()

	// First open a document
	didOpenParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.txt",
			LanguageID: "plaintext",
			Version:    1,
			Text:       "Hello, World!",
		},
	}

	paramsJSON, err := json.Marshal(didOpenParams)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	_, err = server.handleDidOpen(paramsJSON)
	if err != nil {
		t.Fatalf("handleDidOpen failed: %v", err)
	}

	// Now change the document
	didChangeParams := protocol.DidChangeTextDocumentParams{
		TextDocument: protocol.VersionedTextDocumentIdentifier{
			URI:     "file:///test.txt",
			Version: 2,
		},
		ContentChanges: []protocol.TextDocumentContentChangeEvent{
			{
				Text: "Hello, LSP!",
			},
		},
	}

	changeParamsJSON, err := json.Marshal(didChangeParams)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	_, err = server.handleDidChange(changeParamsJSON)
	if err != nil {
		t.Fatalf("handleDidChange failed: %v", err)
	}

	doc, exists := server.GetDocument("file:///test.txt")
	if !exists {
		t.Fatal("Expected document to exist")
	}

	if doc.Text != "Hello, LSP!" {
		t.Errorf("Expected text 'Hello, LSP!', got %s", doc.Text)
	}

	if doc.Version != 2 {
		t.Errorf("Expected version 2, got %d", doc.Version)
	}
}

func TestServerDidClose(t *testing.T) {
	server := NewServer()

	// First open a document
	didOpenParams := protocol.DidOpenTextDocumentParams{
		TextDocument: protocol.TextDocumentItem{
			URI:        "file:///test.txt",
			LanguageID: "plaintext",
			Version:    1,
			Text:       "Hello, World!",
		},
	}

	paramsJSON, err := json.Marshal(didOpenParams)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	_, err = server.handleDidOpen(paramsJSON)
	if err != nil {
		t.Fatalf("handleDidOpen failed: %v", err)
	}

	// Now close the document
	didCloseParams := protocol.DidCloseTextDocumentParams{
		TextDocument: protocol.TextDocumentIdentifier{
			URI: "file:///test.txt",
		},
	}

	closeParamsJSON, err := json.Marshal(didCloseParams)
	if err != nil {
		t.Fatalf("Failed to marshal params: %v", err)
	}

	_, err = server.handleDidClose(closeParamsJSON)
	if err != nil {
		t.Fatalf("handleDidClose failed: %v", err)
	}

	_, exists := server.GetDocument("file:///test.txt")
	if exists {
		t.Error("Expected document to not exist after close")
	}
}

func TestServerShutdown(t *testing.T) {
	server := NewServer()

	_, err := server.handleShutdown(nil)
	if err != nil {
		t.Fatalf("handleShutdown failed: %v", err)
	}

	if !server.IsShutdown() {
		t.Error("Expected server to be shutdown")
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func stringPtr(s string) *string {
	return &s
}
